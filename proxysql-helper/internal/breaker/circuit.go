package breaker

import (
	"sync/atomic"
	"time"
)

type State int32

const (
	Closed State = iota
	Open
	HalfOpen
)

type CircuitBreaker struct {
	state       atomic.Int32
	failures    atomic.Int32
	lastFailure atomic.Int64 // unix nanos

	failureThreshold int32
	openCooldown     time.Duration
	maxHalfOpen      int32
	halfOpenCount    atomic.Int32
}

func New() *CircuitBreaker {
	cb := &CircuitBreaker{failureThreshold: 3, openCooldown: 3 * time.Second, maxHalfOpen: 1}
	cb.state.Store(int32(Closed))
	return cb
}

func (c *CircuitBreaker) Allow() bool {
	s := State(c.state.Load())
	switch s {
	case Closed:
		return true
	case Open:
		lf := time.Unix(0, c.lastFailure.Load())
		if time.Since(lf) >= c.openCooldown {
			if c.state.CompareAndSwap(int32(Open), int32(HalfOpen)) {
				c.halfOpenCount.Store(0)
			}
		}
		return false
	case HalfOpen:
		if c.halfOpenCount.Add(1) <= c.maxHalfOpen {
			return true
		}
		c.halfOpenCount.Add(-1)
		return false
	default:
		return false
	}
}

func (c *CircuitBreaker) OnSuccess() {
	if State(c.state.Load()) == HalfOpen {
		c.state.Store(int32(Closed))
		c.failures.Store(0)
		c.halfOpenCount.Store(0)
	}
}

func (c *CircuitBreaker) OnFailure() {
	c.lastFailure.Store(time.Now().UnixNano())
	s := State(c.state.Load())
	if s == HalfOpen {
		c.state.Store(int32(Open))
		c.halfOpenCount.Store(0)
		return
	}
	if c.failures.Add(1) >= c.failureThreshold {
		c.state.Store(int32(Open))
	}
}

func (c *CircuitBreaker) State() State { return State(c.state.Load()) }
