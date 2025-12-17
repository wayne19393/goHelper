package router

import (
	"math/rand"
	"sync/atomic"

	"proxysql-galera-app/internal/pool"
)

type RoundRobin struct{ idx atomic.Int32 }

func (r *RoundRobin) PickNode(nodes []*pool.Node) *pool.Node {
	if len(nodes) == 0 {
		return nil
	}
	i := int(r.idx.Add(1)-1) % len(nodes)
	return nodes[i]
}
func (r *RoundRobin) Name() string { return "round_robin" }

type Random struct{}

func (Random) PickNode(nodes []*pool.Node) *pool.Node {
	if len(nodes) == 0 {
		return nil
	}
	return nodes[rand.Intn(len(nodes))]
}
func (Random) Name() string { return "random" }

type LowestLatency struct{}

func (LowestLatency) PickNode(nodes []*pool.Node) *pool.Node {
	if len(nodes) == 0 {
		return nil
	}
	best := nodes[0]
	bestEWMA := pool.EwmaToFloat(best.EWMA.Load())
	for _, n := range nodes[1:] {
		if e := pool.EwmaToFloat(n.EWMA.Load()); e < bestEWMA {
			best = n
			bestEWMA = e
		}
	}
	return best
}
func (LowestLatency) Name() string { return "lowest_latency" }
