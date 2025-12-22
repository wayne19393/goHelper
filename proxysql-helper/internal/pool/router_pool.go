package pool

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"proxysql-galera-app/internal/breaker"
)

type Node struct {
	Name    string
	DSN     string
	DB      *sql.DB
	Breaker *breaker.CircuitBreaker
	EWMA    atomic.Int64 // float64 bits
}

type RouterStrategy interface {
	PickNode(nodes []*Node) *Node
	Name() string
}

func EwmaToFloat(v int64) float64 { return math.Float64frombits(uint64(v)) }
func FloatToEwma(f float64) int64 { return int64(math.Float64bits(f)) }

type RouterPool struct {
	nodes   []*Node
	router  RouterStrategy
	retries int
	mu      sync.RWMutex
}

func NewRouterPool(endpoints []string, user, pass, dbname string, maxOpen, maxIdle int, router RouterStrategy) (*RouterPool, error) {
	nodes := make([]*Node, 0, len(endpoints))
	for i, ep := range endpoints {
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&multiStatements=true&timeout=2s&readTimeout=5s&writeTimeout=5s", user, pass, strings.TrimSpace(ep), dbname)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, err
		}
		db.SetMaxOpenConns(maxOpen)
		db.SetMaxIdleConns(maxIdle)
		db.SetConnMaxLifetime(30 * time.Minute)
		n := &Node{Name: fmt.Sprintf("px-%d", i+1), DSN: dsn, DB: db, Breaker: breaker.New()}
		n.EWMA.Store(FloatToEwma(50))
		nodes = append(nodes, n)
	}
	rp := &RouterPool{nodes: nodes, router: router, retries: 3}
	for _, n := range nodes {
		go rp.pinger(n)
	}
	return rp, nil
}

func (rp *RouterPool) pinger(n *Node) {
	t := time.NewTicker(2 * time.Second)
	defer t.Stop()
	alpha := 0.2
	for range t.C {
		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
		err := n.DB.PingContext(ctx)
		cancel()
		lat := float64(time.Since(start).Microseconds())
		curr := EwmaToFloat(n.EWMA.Load())
		if curr <= 0 {
			curr = lat
		}
		updated := (1-alpha)*curr + alpha*lat
		n.EWMA.Store(FloatToEwma(updated))
		if err != nil {
			n.Breaker.OnFailure()
		} else {
			n.Breaker.OnSuccess()
		}
	}
}

func (rp *RouterPool) WithConn(ctx context.Context, fn func(ctx context.Context, db *sql.DB) error) error {
	var lastErr error
	for attempt := 0; attempt < rp.retries; attempt++ {
		n := rp.pickNode()
		if n == nil {
			return errors.New("no available nodes")
		}
		if !n.Breaker.Allow() {
			lastErr = errors.New("circuit open")
			continue
		}
		subctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		err := fn(subctx, n.DB)
		cancel()
		if err == nil {
			n.Breaker.OnSuccess()
			return nil
		}
		if isRetryable(err) {
			n.Breaker.OnFailure()
			lastErr = err
			sleepBackoff(attempt)
			continue
		}
		return err
	}
	return lastErr
}

func (rp *RouterPool) WithTx(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
	return rp.WithConn(ctx, func(ctx context.Context, db *sql.DB) error {
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
		if err != nil {
			return err
		}
		err = fn(ctx, tx)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	})
}

func (rp *RouterPool) pickNode() *Node {
	rp.mu.RLock()
	defer rp.mu.RUnlock()
	healthy := make([]*Node, 0, len(rp.nodes))
	for _, n := range rp.nodes {
		if n.Breaker.State() != breaker.Open {
			healthy = append(healthy, n)
		}
	}
	if len(healthy) == 0 {
		healthy = rp.nodes
	}
	return rp.router.PickNode(healthy)
}

func isRetryable(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "timeout") || strings.Contains(msg, "deadlock") || strings.Contains(msg, "connection refused") || strings.Contains(msg, "connection reset") || strings.Contains(msg, "broken pipe")
}

func sleepBackoff(i int) { time.Sleep(time.Duration(200*(1<<i)) * time.Millisecond) }
