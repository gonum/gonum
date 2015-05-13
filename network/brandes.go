// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package network

import (
	"github.com/gonum/graph"
	"github.com/gonum/graph/internal"
)

// Betweenness returns the betweenness centrality for nodes in the unweighted graph g.
func Betweenness(g graph.Graph) map[int]float64 {
	// Brandes' algorithm for finding betweenness centrality for nodes in
	// and unweighted graph:
	//
	// http://www.inf.uni-konstanz.de/algo/publications/b-fabc-01.pdf

	// TODO(kortschak): Consider using the parallel algorithm when
	// GOMAXPROCS != 0.
	//
	// http://htor.inf.ethz.ch/publications/img/edmonds-hoefler-lumsdaine-bc.pdf

	// Also note special case for sparse networks:
	// http://wwwold.iit.cnr.it/staff/marco.pellegrini/papiri/asonam-final.pdf

	nodes := g.NodeList()
	cb := make(map[int]float64, len(nodes))
	for _, n := range nodes {
		cb[n.ID()] = 0
	}

	var (
		stack internal.NodeStack
		p     = make(map[int][]graph.Node, len(nodes))
		sigma = make(map[int]float64, len(nodes))
		d     = make(map[int]int, len(nodes))
		delta = make(map[int]float64, len(nodes))
		queue internal.NodeQueue
	)
	for _, s := range nodes {
		stack = stack[:0]

		for _, w := range nodes {
			p[w.ID()] = p[w.ID()][:0]
		}

		for _, t := range nodes {
			sigma[t.ID()] = 0
			d[t.ID()] = -1
		}
		sigma[s.ID()] = 1
		d[s.ID()] = 0

		queue.Enqueue(s)
		for queue.Len() != 0 {
			v := queue.Dequeue()
			stack.Push(v)
			for _, w := range g.Neighbors(v) {
				// w found for the first time?
				if d[w.ID()] < 0 {
					queue.Enqueue(w)
					d[w.ID()] = d[v.ID()] + 1
				}
				// shortest path to w via v?
				if d[w.ID()] == d[v.ID()]+1 {
					sigma[w.ID()] += sigma[v.ID()]
					p[w.ID()] = append(p[w.ID()], v)
				}
			}
		}

		for _, v := range nodes {
			delta[v.ID()] = 0
		}
		// S returns vertices in order of non-increasing distance from s
		for stack.Len() != 0 {
			w := stack.Pop()
			for _, v := range p[w.ID()] {
				delta[v.ID()] += sigma[v.ID()] / sigma[w.ID()] * (1 + delta[w.ID()])
			}
			if w.ID() != s.ID() {
				cb[w.ID()] += delta[w.ID()]
			}
		}
	}

	return cb
}
