// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/temporal"
	"math"
)

// Earliest stores the earliest-arrival paths
type Earliest struct {
	from  graph.Node
	at    uint64
	until uint64
	nodes map[int64]struct{
		earliest uint64
		via      []graph.Node
	}
}

func (e *Earliest) set(v graph.Node, t uint64, p []graph.Node) {
	e.nodes[v.ID()] = struct{
		earliest uint64
		via      []graph.Node
	}{
		t,
		p,
	}
}

func (e *Earliest) From() graph.Node { return e.from }

func (e *Earliest) At() uint64 { return e.at }

func (e *Earliest) Until() uint64 { return e.until }

func (e *Earliest) To(uid int64) (path []graph.Node, duration uint64) {
	eu, ok := e.nodes[uid]
	if !ok {
		return nil, ^uint64(0)
	}
	duration = eu.earliest
	u := temporal.Node(uid)
	path = append(eu.via, &u)
	return
}

func (e *Earliest) Len() int {
	return len(e.nodes)
}

func (e *Earliest) Weight(uid, vid int64) float64 {
	if uid != e.from.ID() {
		return math.Inf(1)
	}
	ev, ok := e.nodes[vid]
	if !ok {
		return math.Inf(1)
	}

	return float64(ev.earliest)
}

// EarliestArrivalFrom computes the earliest-arrival paths to all nodes reachable from u where
// from <= t <=until. It is an implementation of Algorithm 1. described in
// https://www.vldb.org/pvldb/vol7/p721-wu.pdf
func EarliestArrivalFrom(g graph.LineStreamer, from graph.Node, at uint64, until uint64) Earliest {
	earliest := Earliest{
		from:  from,
		at:    at,
		until: until,
		nodes: make(map[int64]struct{
			earliest uint64
			via      []graph.Node
		}),
	}
	earliest.set(from, at, []graph.Node{})
	s := g.LineStream()
	for s.Next() {
		l := s.Line()
		u := l.From()
		uid := u.ID()
		eu, ok := earliest.nodes[uid]
		if !ok {
			continue
		}
		var ets, ete uint64
		if tl, ok := l.(graph.TemporalLine); ok {
			ets, ete = tl.Interval()
		} else {
			// When not dealing with temporal edges degnerates to case
			// that gives reachability of nodes in the static graph
			// from u
			ets, ete = uint64(at), uint64(at)
		}
		if eu.earliest <= ets && ete <= until {
			v := l.To()
			vid := v.ID()
			ev, ok := earliest.nodes[vid]
			if !ok || ete < ev.earliest {
				earliest.set(v, ete, append(eu.via, u))
			}
		} else if ets >= until {
			break
		}
	}
	return earliest
}