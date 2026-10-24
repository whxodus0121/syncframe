package realtime

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"syncframe/internal/model"
)

type Presence struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}
type Signal struct {
	BoardID string    `json:"boardId"`
	Type    string    `json:"type"`
	Peer    *Presence `json:"peer,omitempty"`
	X       float64   `json:"x,omitempty"`
	Y       float64   `json:"y,omitempty"`
}
type Peer struct {
	Board     string
	User      Presence
	Wake      chan struct{}
	Ephemeral chan Signal
}
type Bus struct {
	Redis *redis.Client
	mu    sync.Mutex
	peers map[string]*Peer
}

func New(addr string) *Bus {
	return &Bus{Redis: redis.NewClient(&redis.Options{Addr: addr, DialTimeout: time.Second, ReadTimeout: time.Second, WriteTimeout: time.Second, MaxRetries: 0}), peers: map[string]*Peer{}}
}
func (b *Bus) Run(ctx context.Context) {
	sub := b.Redis.Subscribe(ctx, "syncframe:hints")
	defer sub.Close()
	ch := sub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			var s Signal
			if json.Unmarshal([]byte(msg.Payload), &s) == nil {
				b.deliver(s)
			}
		}
	}
}
func (b *Bus) deliver(s Signal) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, p := range b.peers {
		if p.Board != s.BoardID {
			continue
		}
		if s.Type == "cursor" {
			select {
			case p.Ephemeral <- s:
			default:
			}
		} else {
			select {
			case p.Wake <- struct{}{}:
			default:
			}
		}
	}
}
func (b *Bus) Notify(ctx context.Context, s Signal) {
	b.deliver(s)
	raw, _ := json.Marshal(s)
	_ = b.Redis.Publish(ctx, "syncframe:hints", raw).Err()
}
func (b *Bus) Add(board, name string) *Peer {
	id := model.ID()
	colors := []string{"#7963db", "#389b80", "#d39345", "#cf7194", "#4b9bb9", "#b17ac7"}
	p := &Peer{Board: board, User: Presence{ID: id, Name: name, Color: colors[int(id[0])%len(colors)]}, Wake: make(chan struct{}, 1), Ephemeral: make(chan Signal, 16)}
	b.mu.Lock()
	b.peers[p.User.ID] = p
	b.mu.Unlock()
	return p
}
func (b *Bus) Remove(ctx context.Context, p *Peer) {
	b.mu.Lock()
	delete(b.peers, p.User.ID)
	b.mu.Unlock()
	_ = b.Redis.Del(ctx, "syncframe:presence:"+p.Board+":"+p.User.ID).Err()
}
func (b *Bus) Presence(ctx context.Context, p *Peer) ([]Presence, bool) {
	raw, _ := json.Marshal(p.User)
	key := "syncframe:presence:" + p.Board + ":" + p.User.ID
	// Membership is bounded by TTL in a sorted set; connection payloads expire separately.
	now := time.Now().Unix()
	pipe := b.Redis.TxPipeline()
	pipe.Set(ctx, key, raw, 15*time.Second)
	pipe.ZAdd(ctx, "syncframe:members:"+p.Board, redis.Z{Score: float64(now + 15), Member: p.User.ID})
	pipe.ZRemRangeByScore(ctx, "syncframe:members:"+p.Board, "-inf", strconv.FormatInt(now, 10))
	pipe.Expire(ctx, "syncframe:members:"+p.Board, 30*time.Second)
	if _, err := pipe.Exec(ctx); err != nil {
		b.mu.Lock()
		defer b.mu.Unlock()
		out := []Presence{}
		for _, v := range b.peers {
			if v.Board == p.Board {
				out = append(out, v.User)
			}
		}
		return out, false
	}
	ids, err := b.Redis.ZRange(ctx, "syncframe:members:"+p.Board, 0, 199).Result()
	if err != nil {
		return []Presence{p.User}, false
	}
	out := []Presence{}
	for _, id := range ids {
		v, err := b.Redis.Get(ctx, "syncframe:presence:"+p.Board+":"+id).Bytes()
		if err != nil {
			continue
		}
		var u Presence
		if json.Unmarshal(v, &u) == nil {
			out = append(out, u)
		}
	}
	return out, true
}
