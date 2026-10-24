package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syncframe/internal/httpapi"
	"syncframe/internal/pipeline"
	"syncframe/internal/realtime"
	"syncframe/internal/store"
	"syscall"
	"time"
)

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	s, err := store.Open(env("MYSQL_DSN", "syncframe:syncframe@tcp(127.0.0.1:33060)/syncframe?parseTime=true&timeout=3s&readTimeout=5s&writeTimeout=5s"))
	if err != nil {
		slog.Error("database configuration", "error", err)
		os.Exit(1)
	}
	defer s.DB.Close()
	call, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	for {
		err = s.DB.PingContext(call)
		if err == nil {
			err = s.Migrate(call)
			break
		}
		select {
		case <-call.Done():
			slog.Error("database unavailable", "error", err)
			os.Exit(1)
		case <-time.After(time.Second):
		}
	}
	if err != nil {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}
	bus := realtime.New(env("REDIS_ADDR", "127.0.0.1:63790"))
	defer bus.Redis.Close()
	api := httpapi.New(s, bus, env("APP_ORIGIN", ""), env("WEB_DIR", "web/dist"))
	server := &http.Server{Addr: env("ADDR", ":8080"), Handler: api.Handler(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 16384}
	var wg sync.WaitGroup
	start := func(f func()) { wg.Add(1); go func() { defer wg.Done(); f() }() }
	start(func() { bus.Run(ctx) })
	if os.Getenv("KAFKA_DISABLED") != "1" {
		broker := env("KAFKA_BROKER", "127.0.0.1:19092")
		start(func() { pipeline.RunPublisher(ctx, s, broker) })
		start(func() { pipeline.RunConsumer(ctx, s, broker) })
	}
	start(func() {
		slog.Info("SYNCFRAME listening", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server stopped", "error", err)
			stop()
		}
	})
	<-ctx.Done()
	api.Close()
	shutdown, c := context.WithTimeout(context.Background(), 10*time.Second)
	defer c()
	_ = server.Shutdown(shutdown)
	wg.Wait()
}
