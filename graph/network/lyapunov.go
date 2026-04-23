// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package network

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/mat"
)

// LyapunovSolve is the main entry point for graph consensus.
func LyapunovSolve(g graph.Graph, initialState *mat.VecDense) *mat.VecDense {
	n := initialState.Len()

	// Case 1: Empty graph or single node
	if n <= 1 {
		return initialState
	}

	// Case 2: Simple Undirected & Unweighted Graph (Arithmetic Consensus)
	// If the graph is undirected and not weighted, the equilibrium is the average.
	if _, isUndirected := g.(graph.Undirected); isUndirected {
		if _, isWeighted := g.(graph.Weighted); !isWeighted {
			return solveArithmeticAverage(initialState)
		}
	}

	// Case 3: Complex case (Directed, Weighted, or Constrained)
	// We use the robust Lyapunov minimization.
	gamma := 0.1
	beta := 0.5
	iterations := 1000

	// We type assert to WeightedDirected as it is the most general interface
	// supported by our unnormalizedLaplacianMatrix function.
	return minimizeLyapunov(g.(graph.WeightedDirected), initialState, gamma, beta, iterations)
}

// solveArithmeticAverage computes the consensus in linear time O(n).
// This is used as a Fast Path when the final state is mathematically known.
func solveArithmeticAverage(initialState *mat.VecDense) *mat.VecDense {
	n := initialState.Len()
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += initialState.AtVec(i)
	}
	avg := sum / float64(n)

	data := make([]float64, n)
	for i := range data {
		data[i] = avg
	}
	return mat.NewVecDense(n, data)
}

// unnormalizedLaplacianMatrix computes L = D - A for weighted directed graphs.
// D is the out-degree diagonal matrix and A is the weighted adjacency matrix.
func unnormalizedLaplacianMatrix(g graph.WeightedDirected) *mat.Dense {
	nodes := graph.NodesOf(g.Nodes())
	n := len(nodes)
	L := mat.NewDense(n, n, nil)

	// Map node IDs to matrix indices (0 to n-1)
	idToIndex := make(map[int64]int)
	for i, node := range nodes {
		idToIndex[node.ID()] = i
	}

	for i, u := range nodes {
		var outDegree float64

		// Iterate over successors of node u
		toNeighbors := g.From(u.ID())
		for toNeighbors.Next() {
			v := toNeighbors.Node()
			j := idToIndex[v.ID()]

			// Retrieve the actual weight of the edge u -> v
			weight := g.WeightedEdge(u.ID(), v.ID()).Weight()

			// L_ij = -weight
			L.Set(i, j, -weight)

			// Accumulate out-degree for the diagonal
			outDegree += weight
		}

		// L_ii = weighted out-degree
		L.Set(i, i, outDegree)
	}

	return L
}

// Lyapunovfct computes the energy of the system given the Laplacian L and state x.
// The energy (stability) is defined by the quadratic form x^T * L * x.
func Lyapunovfct(L *mat.Dense, x *mat.VecDense) float64 {
	n, _ := L.Dims()

	// lx = L * x
	lx := mat.NewVecDense(n, nil)
	lx.MulVec(L, x)

	// Result = x^T * lx
	return mat.Dot(x, lx)
}

// minimizeLyapunov performs gradient descent with backtracking line search
// to find the state x that minimizes the graph's energy.
func minimizeLyapunov(g graph.WeightedDirected, initialState *mat.VecDense, initialGamma float64, beta float64, iterations int) *mat.VecDense {
	L := unnormalizedLaplacianMatrix(g)
	epsilon := 1e-6
	n := initialState.Len()

	// Initialize state x with a copy of initialState
	x := mat.NewVecDense(n, nil)
	x.CopyVec(initialState)

	// Pre-allocate workspaces to optimize memory
	grad := mat.NewVecDense(n, nil)
	step := mat.NewVecDense(n, nil)
	next := mat.NewVecDense(n, nil)

	for i := 0; i < iterations; i++ {
		// Compute gradient: ∇V(x) = L * x
		grad.MulVec(L, x)
		gradnorm := mat.Norm(grad, 2)

		if gradnorm < epsilon {
			break
		}

		// Reset gamma for backtracking at each main iteration
		gamma := initialGamma

		// Backtracking Line Search to find optimal step size
		for {
			step.ScaleVec(gamma, grad)
			next.SubVec(x, step)

			if Lyapunovfct(L, next) <= Lyapunovfct(L, x)-0.5*gamma*gradnorm*gradnorm {
				break
			}

			gamma *= beta
			if gamma < 1e-12 {
				break
			}
		}

		// Update state x
		x.CopyVec(next)
	}
	return x
}
