// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gen

import (
	"math"
	"math/rand"
	"sort"

	"github.com/gonum/graph"
	"github.com/gonum/graph/simple"
)

// RandomDAG generates a random directed acyclic graph uniformly at random
// over the space of DAGs. Random DAG uses the algorithm in
//  Ide, Jaime S., and Fabio G. Cozman. "Random generation of Bayesian networks."
//  Brazilian symposium on artificial intelligence. Springer Berlin Heidelberg, 2002.
// nodes specifies the number of nodes in the network, and steps is the number of
// steps to take in the Markov chain.
func RandomDAG(nodes, steps int) *simple.DirectedGraph {
	// Construct the line DAG to start.
	g := simple.NewDirectedGraph(0, math.NaN())
	for i := simple.Node(0); i < simple.Node(nodes); i++ {
		n := simple.Node(i)
		g.AddNode(n)
		if i != 0 {
			edge := simple.Edge{
				i - 1, i, 1,
			}
			g.SetEdge(edge)
		}
	}

	for step := 0; step < steps; step++ {
		// Pick two random nodes
		i := simple.Node(rand.Intn(nodes))
		j := simple.Node(rand.Intn(nodes - 1))
		if j >= i {
			j++
		}
		if g.HasEdgeFromTo(i, j) {
			// If there is an edge between them, delete it if it keeps the graph
			// connection. The graph is still connected if there is a path between
			// i and j after deletion.
			e := g.Edge(i, j)
			g.RemoveEdge(e)
			u := graph.Undirect{G: g}
			if !hasPath(u, i, j) {
				g.SetEdge(e)
			}
		} else {
			// There is no edge between i and j. Add an edge between i and j
			// if the graph stays acyclic. This edge can be added if we cannot
			// get from j to i.
			if !hasPath(g, j, i) {
				g.SetEdge(simple.Edge{i, j, 1})
			}
		}
	}
	edges := g.Edges()
	sort.Sort(EdgeSorter(edges))
	return g
}
