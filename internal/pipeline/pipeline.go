package pipeline

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
	"syncframe/internal/model"
	"syncframe/internal/store"
)

const Topic = "syncframe.edits.v1"

func Writer(broker string) *kafka.Writer {
	return &kafka.Writer{Addr: kafka.TCP(broker), Topic: Topic, Balancer: &kafka.Hash{}, RequiredAcks: kafka.RequireAll, AllowAutoTopicCreation: true, ReadTimeout: 3 * time.Second, WriteTimeout: 3 * time.Second, MaxAttempts: 1, BatchTimeout: 20 * time.Millisecond}
}

// SQL lock is held only by the asynchronous worker, never by editing transactions.
// Crash after Kafka ack but before SQL commit intentionally causes redelivery.
func PublishOne(ctx context.Context, s *store.Store, w *kafka.Writer) (bool, error) {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return false, err
	}
	defer tx.Rollback()
	var id int64
	var board string
	var payload []byte
	err = tx.QueryRowContext(ctx, "SELECT id,board_id,payload FROM outbox WHERE published_at IS NULL ORDER BY id LIMIT 1 FOR UPDATE SKIP LOCKED").Scan(&id, &board, &payload)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	if err = w.WriteMessages(ctx, kafka.Message{Key: []byte(board), Value: payload}); err != nil {
		return false, err
	}
	if _, err = tx.ExecContext(ctx, "UPDATE outbox SET published_at=CURRENT_TIMESTAMP(6) WHERE id=?", id); err != nil {
		return false, err
	}
	return true, tx.Commit()
}
func RunPublisher(ctx context.Context, s *store.Store, broker string) {
	w := Writer(broker)
	defer w.Close()
	for ctx.Err() == nil {
		call, cancel := context.WithTimeout(ctx, 5*time.Second)
		worked, err := PublishOne(call, s, w)
		cancel()
		if err != nil && ctx.Err() == nil {
			slog.Warn("audit publish delayed", "error", err)
		}
		if !worked || err != nil {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
			}
		}
	}
}
func RunConsumer(ctx context.Context, s *store.Store, broker string) {
	r := kafka.NewReader(kafka.ReaderConfig{Brokers: []string{broker}, Topic: Topic, GroupID: "syncframe-audit-v1", MinBytes: 1, MaxBytes: 1e6, MaxWait: time.Second, CommitInterval: 0, StartOffset: kafka.FirstOffset})
	defer r.Close()
	for ctx.Err() == nil {
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			slog.Warn("audit fetch delayed", "error", err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
			}
			continue
		}
		var e model.Event
		if err = json.Unmarshal(msg.Value, &e); err != nil {
			slog.Error("invalid audit event; consumer halted", "error", err)
			return
		}
		for ctx.Err() == nil {
			call, cancel := context.WithTimeout(ctx, 5*time.Second)
			err = s.Audit(call, e)
			cancel()
			if err == nil {
				break
			}
			slog.Warn("audit SQL retry", "error", err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
			}
		}
		if ctx.Err() != nil {
			return
		}
		if err = r.CommitMessages(ctx, msg); err != nil && ctx.Err() == nil {
			slog.Warn("audit offset commit delayed", "error", err)
		}
	}
}
