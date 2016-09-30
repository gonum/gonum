// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gen

import (
	"math"
	"math/rand"

	"github.com/gonum/graph"
	"github.com/gonum/graph/simple"
	"github.com/gonum/graph/topo"
)

type UniformGenerator interface {
	Modify(g GraphBuilder, from, to graph.Node)
	Start() graph.Directed
}

// UniformDAG generates a graph uniformly at random over the set of directed
// acyclic graphs. It uses the algorithm described in
//  Ide, Jaime S., and Fabio G. Cozman. "Random generation of Bayesian networks."
//  Brazilian symposium on artificial intelligence. Springer Berlin Heidelberg, 2002.
type UniformDAG struct {
	N       int
	EdgeGen func(from, to graph.Node) graph.Edge
}

func (u UniformDAG) Start() graph.Directed {
	if u.N == 0 {
		panic("uniformdag: number of edges must be positive")
	}
	// Starting graph is the line graph (A -> B -> C ...)
	g := simple.NewDirectedGraph(0, math.NaN())
	for i := simple.Node(0); i < simple.Node(u.N); i++ {
		n := simple.Node(i)
		g.AddNode(n)
		if i != 0 {
			edge := simple.Edge{
				i - 1, i, 1,
			}
			g.SetEdge(edge)
		}
	}
	return g
}

func (u UniformDAG) Modify(g DirectedMutator, from, to graph.Node) {
	if g.HasEdgeFromTo(from, to) {
		// If there an edge from -> to, delete it if it keeps the graph connected.
		// The graph is still connected if there is an undirected path between
		// from and to after delection.
		e := g.Edge(from, to)
		g.RemoveEdge(e)
		if !topo.PathExistsIn(graph.Undirect{G: g}, from, to) {
			g.SetEdge(e)
		}
	} else {
		// If there is no edge from -> to, add it if the graph stays acyclic.
		// This can be added if we cannot get from j to i.
		if !topo.PathExistsIn(g, to, from) {
			var e graph.Edge
			if u.EdgeGen == nil {
				e = u.EdgeGen(from, to)
			} else {
				e = simple.Edge{from, to, 1}
			}
			g.SetEdge(e)
		}
	}
}

// UniformEdges generates the edges of a graph uniformly at random across the
// class of graphs specified by class.
//
// UniformEdges runs a Markov chain starting with the graph specified by start
// and running a number of steps equal to steps. Each step in the markov chain
// proposes an edge from from to to. The class function modifies the graph
// according to the generation procedure. If start is nil, the initial graph
// will be that specified in class.DefaultStart().
//
// class cannot be an arbitrary function, it must obey specific properties in
// order for the generation probability to be uniform over the class of possible
// graphs. class must ensure that the implicit Markov chain is irreducible,
// aperiodic, and, importantly, doubly stochastic. Furthermore, it must be that
// the starting graph
//
// step specifies the number of steps in the Markov chain. step should be
// sufficiently high to allow for convergence of the Markov chain.
//
// src specifies the random number generator. If src is nil, math/rand is used.
//
// More information can be found in
//  Ide, Jaime S., and Fabio G. Cozman. "Random generation of Bayesian networks."
//  Brazilian symposium on artificial intelligence. Springer Berlin Heidelberg, 2002.
func UniformEdges(dst DirectedMutator, start graph.Directed, class UniformGenerator, steps int, src *rand.Rand) {
	if steps <= 0 {
		panic("gen: number of steps must be positive")
	}
	if start == nil {
		start = class.Start()
	}
	var rnd func(int) int
	if src == nil {
		rnd = rand.Intn
	} else {
		rnd = src.Intn
	}

	// Copy the starting graph into the graph builder.
	nodes := start.Nodes()
	for _, node := range nodes {
		dst.AddNode(node)
		tos := start.From(node)
		for _, to := range tos {
			edge := start.Edge(node, to)
			dst.SetEdge(edge)
		}
	}
	nNodes := len(nodes)

	// Run the MarkovChain, choosing an edge at random in each step.
	for step := 0; step < steps; step++ {
		i := rand.Intn(nNodes)
		j := rand.Intn(nNodes - 1)
		if j >= i {
			j++
		}
		class.Modify(dst, nodes[i], nodes[j])
	}
}
