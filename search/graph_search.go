// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

import (
	"container/heap"
	"fmt"
	"math"
	"sort"

	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
	"github.com/gonum/graph/internal"
	"github.com/gonum/graph/traverse"
)

// Returns an ordered list consisting of the nodes between start and goal. The path will be the
// shortest path assuming the function heuristicCost is admissible. The second return value is the
// cost, and the third is the number of nodes expanded while searching (useful info for tuning
// heuristics). Negative Costs will cause bad things to happen, as well as negative heuristic
// estimates.
//
// A heuristic is admissible if, for any node in the graph, the heuristic estimate of the cost
// between the node and the goal is less than or set to the true cost.
//
// Performance may be improved by providing a consistent heuristic (though one is not needed to
// find the optimal path), a heuristic is consistent if its value for a given node is less than
// (or equal to) the actual cost of reaching its neighbors + the heuristic estimate for the
// neighbor itself. You can force consistency by making your HeuristicCost function return
// max(NonConsistentHeuristicCost(neighbor,goal), NonConsistentHeuristicCost(self,goal) -
// Cost(self,neighbor)). If there are multiple neighbors, take the max of all of them.
//
// Cost and HeuristicCost take precedence for evaluating cost/heuristic distance. If one is not
// present (i.e. nil) the function will check the graph's interface for the respective interface:
// Coster for Cost and HeuristicCoster for HeuristicCost. If the correct one is present, it will
// use the graph's function for evaluation.
//
// Finally, if neither the argument nor the interface is present, the function will assume
// UniformCost for Cost and NullHeuristic for HeuristicCost.
//
// To run Uniform Cost Search, run A* with the NullHeuristic.
//
// To run Breadth First Search, run A* with both the NullHeuristic and UniformCost (or any cost
// function that returns a uniform positive value.)
func AStar(start, goal graph.Node, g graph.Graph, weight graph.CostFunc, heuristic graph.HeuristicCostFunc) (path []graph.Node, pathCost float64, nodesExpanded int) {
	if weight == nil {
		if g, ok := g.(graph.Coster); ok {
			weight = g.Cost
		} else {
			weight = UniformCost
		}
	}
	if heuristic == nil {
		if g, ok := g.(graph.HeuristicCoster); ok {
			heuristic = g.HeuristicCost
		} else {
			heuristic = NullHeuristic
		}
	}

	closedSet := make(map[int]internalNode)
	openSet := &aStarPriorityQueue{nodes: make([]internalNode, 0), indexList: make(map[int]int)}
	heap.Init(openSet)
	node := internalNode{start, 0, heuristic(start, goal)}
	heap.Push(openSet, node)
	predecessor := make(map[int]graph.Node)

	for openSet.Len() != 0 {
		curr := heap.Pop(openSet).(internalNode)

		nodesExpanded += 1

		if curr.ID() == goal.ID() {
			return rebuildPath(predecessor, goal), curr.gscore, nodesExpanded
		}

		closedSet[curr.ID()] = curr

		for _, neighbor := range g.From(curr.Node) {
			if _, ok := closedSet[neighbor.ID()]; ok {
				continue
			}

			g := curr.gscore + weight(g.Edge(curr.Node, neighbor))

			if existing, exists := openSet.Find(neighbor.ID()); !exists {
				predecessor[neighbor.ID()] = curr
				node = internalNode{neighbor, g, g + heuristic(neighbor, goal)}
				heap.Push(openSet, node)
			} else if g < existing.gscore {
				predecessor[neighbor.ID()] = curr
				openSet.Fix(neighbor.ID(), g, g+heuristic(neighbor, goal))
			}
		}
	}

	return nil, 0, nodesExpanded
}

// An admissible, consistent heuristic that won't speed up computation time at all.
func NullHeuristic(_, _ graph.Node) float64 {
	return 0
}

// Assumes all edges in the graph have the same weight (including edges that don't exist!)
func UniformCost(e graph.Edge) float64 {
	if e == nil {
		return math.Inf(1)
	}
	return 1
}

/* Simple operations */

// Copies a graph into the destination; maintaining all node IDs. The destination
// need not be empty, though overlapping node IDs and conflicting edges will overwrite
// existing data.
func CopyUndirectedGraph(dst graph.MutableGraph, src graph.Graph) {
	var weight graph.CostFunc
	if g, ok := src.(graph.Coster); ok {
		weight = g.Cost
	} else {
		weight = UniformCost
	}

	for _, node := range src.Nodes() {
		succs := src.From(node)
		dst.AddNode(node)
		for _, succ := range succs {
			edge := src.Edge(node, succ)
			dst.AddUndirectedEdge(edge, weight(edge))
		}
	}

}

// Copies a graph into the destination; maintaining all node IDs. The destination
// need not be empty, though overlapping node IDs and conflicting edges will overwrite
// existing data.
func CopyDirectedGraph(dst graph.MutableDirectedGraph, src graph.Graph) {
	var weight graph.CostFunc
	if g, ok := src.(graph.Coster); ok {
		weight = g.Cost
	} else {
		weight = UniformCost
	}

	for _, node := range src.Nodes() {
		succs := src.From(node)
		dst.AddNode(node)
		for _, succ := range succs {
			edge := src.Edge(node, succ)
			dst.AddDirectedEdge(edge, weight(edge))
		}
	}

}

/* Basic Graph tests */

// Unorderable is an error containing sets of unorderable graph.Nodes.
type Unorderable [][]graph.Node

// Error satisfies the error interface.
func (e Unorderable) Error() string {
	const maxNodes = 10
	var n int
	for _, c := range e {
		n += len(c)
	}
	if n > maxNodes {
		// Don't return errors that are too long.
		return fmt.Sprintf("search: no topological ordering: %d nodes in %d cyclic components", n, len(e))
	}
	return fmt.Sprintf("search: no topological ordering: cyclic components: %v", [][]graph.Node(e))
}

// Sort performs a topological sort of the directed graph g returning the 'from' to 'to'
// sort order. If a topological ordering is not possible, an Unorderable error is returned
// listing cyclic components in g with each cyclic component's members sorted by ID. When
// an Unorderable error is returned, each cyclic component's topological position within
// the sorted nodes is marked with a nil graph.Node.
func Sort(g graph.Directed) (sorted []graph.Node, err error) {
	sccs := TarjanSCC(g)
	sorted = make([]graph.Node, 0, len(sccs))
	var sc Unorderable
	for _, s := range sccs {
		if len(s) != 1 {
			sort.Sort(byID(s))
			sc = append(sc, s)
			sorted = append(sorted, nil)
			continue
		}
		sorted = append(sorted, s[0])
	}
	if sc != nil {
		for i, j := 0, len(sc)-1; i < j; i, j = i+1, j-1 {
			sc[i], sc[j] = sc[j], sc[i]
		}
		err = sc
	}
	reverse(sorted)
	return sorted, err
}

// TarjanSCC returns the strongly connected components of the graph g using Tarjan's algorithm.
//
// A strongly connected component of a graph is a set of vertices where it's possible to reach any
// vertex in the set from any other (meaning there's a cycle between them.)
//
// Generally speaking, a directed graph where the number of strongly connected components is equal
// to the number of nodes is acyclic, unless you count reflexive edges as a cycle (which requires
// only a little extra testing.)
//
func TarjanSCC(g graph.Directed) [][]graph.Node {
	nodes := g.Nodes()
	t := tarjan{
		succ: g.From,

		indexTable: make(map[int]int, len(nodes)),
		lowLink:    make(map[int]int, len(nodes)),
		onStack:    make(internal.IntSet, len(nodes)),
	}
	for _, v := range nodes {
		if t.indexTable[v.ID()] == 0 {
			t.strongconnect(v)
		}
	}
	return t.sccs
}

// tarjan implements Tarjan's strongly connected component finding
// algorithm. The implementation is from the pseudocode at
//
// http://en.wikipedia.org/wiki/Tarjan%27s_strongly_connected_components_algorithm?oldid=642744644
//
type tarjan struct {
	succ func(graph.Node) []graph.Node

	index      int
	indexTable map[int]int
	lowLink    map[int]int
	onStack    internal.IntSet

	stack []graph.Node

	sccs [][]graph.Node
}

// strongconnect is the strongconnect function described in the
// wikipedia article.
func (t *tarjan) strongconnect(v graph.Node) {
	vID := v.ID()

	// Set the depth index for v to the smallest unused index.
	t.index++
	t.indexTable[vID] = t.index
	t.lowLink[vID] = t.index
	t.stack = append(t.stack, v)
	t.onStack.Add(vID)

	// Consider successors of v.
	for _, w := range t.succ(v) {
		wID := w.ID()
		if t.indexTable[wID] == 0 {
			// Successor w has not yet been visited; recur on it.
			t.strongconnect(w)
			t.lowLink[vID] = min(t.lowLink[vID], t.lowLink[wID])
		} else if t.onStack.Has(wID) {
			// Successor w is in stack s and hence in the current SCC.
			t.lowLink[vID] = min(t.lowLink[vID], t.indexTable[wID])
		}
	}

	// If v is a root node, pop the stack and generate an SCC.
	if t.lowLink[vID] == t.indexTable[vID] {
		// Start a new strongly connected component.
		var (
			scc []graph.Node
			w   graph.Node
		)
		for {
			w, t.stack = t.stack[len(t.stack)-1], t.stack[:len(t.stack)-1]
			t.onStack.Remove(w.ID())
			// Add w to current strongly connected component.
			scc = append(scc, w)
			if w.ID() == vID {
				break
			}
		}
		// Output the current strongly connected component.
		t.sccs = append(t.sccs, scc)
	}
}

// IsPath returns true for a connected path within a graph.
//
// IsPath returns true if, starting at path[0] and ending at path[len(path)-1], all nodes between
// are valid neighbors. That is, for each element path[i], path[i+1] is a valid successor.
//
// As special cases, IsPath returns true for a nil or zero length path, and for a path of length 1
// (only one node) but only if the node listed in path exists within the graph.
//
// Graph must be non-nil.
func IsPath(path []graph.Node, g graph.Graph) bool {
	var canReach func(u, v graph.Node) bool
	switch g := g.(type) {
	case graph.Directed:
		canReach = func(u, v graph.Node) bool {
			return g.EdgeFromTo(u, v) != nil
		}
	default:
		canReach = g.HasEdge
	}

	if path == nil || len(path) == 0 {
		return true
	} else if len(path) == 1 {
		return g.Has(path[0])
	}

	for i := 0; i < len(path)-1; i++ {
		if !canReach(path[i], path[i+1]) {
			return false
		}
	}

	return true
}

/* Implements minimum-spanning tree algorithms;
puts the resulting minimum spanning tree in the dst graph */

// Generates a minimum spanning tree with sets.
//
// As with other algorithms that use Cost, the order of precedence is
// Argument > Interface > UniformCost.
//
// The destination must be empty (or at least disjoint with the node IDs of the input)
func Prim(dst graph.MutableGraph, g graph.EdgeListGraph, weight graph.CostFunc) {
	if weight == nil {
		if g, ok := g.(graph.Coster); ok {
			weight = g.Cost
		} else {
			weight = UniformCost
		}
	}

	nlist := g.Nodes()

	if nlist == nil || len(nlist) == 0 {
		return
	}

	dst.AddNode(nlist[0])
	remainingNodes := make(internal.IntSet)
	for _, node := range nlist[1:] {
		remainingNodes.Add(node.ID())
	}

	edgeList := g.Edges()
	for remainingNodes.Count() != 0 {
		var edges []concrete.WeightedEdge
		for _, edge := range edgeList {
			if (dst.Has(edge.From()) && remainingNodes.Has(edge.To().ID())) ||
				(dst.Has(edge.To()) && remainingNodes.Has(edge.From().ID())) {

				edges = append(edges, concrete.WeightedEdge{Edge: edge, Cost: weight(edge)})
			}
		}

		sort.Sort(byWeight(edges))
		myEdge := edges[0]

		dst.AddUndirectedEdge(myEdge.Edge, myEdge.Cost)
		remainingNodes.Remove(myEdge.Edge.From().ID())
	}

}

// Generates a minimum spanning tree for a graph using discrete.DisjointSet.
//
// As with other algorithms with Cost, the precedence goes Argument > Interface > UniformCost.
//
// The destination must be empty (or at least disjoint with the node IDs of the input)
func Kruskal(dst graph.MutableGraph, g graph.EdgeListGraph, weight graph.CostFunc) {
	if weight == nil {
		if g, ok := g.(graph.Coster); ok {
			weight = g.Cost
		} else {
			weight = UniformCost
		}
	}

	edgeList := g.Edges()
	edges := make([]concrete.WeightedEdge, 0, len(edgeList))
	for _, edge := range edgeList {
		edges = append(edges, concrete.WeightedEdge{Edge: edge, Cost: weight(edge)})
	}

	sort.Sort(byWeight(edges))

	ds := newDisjointSet()
	for _, node := range g.Nodes() {
		ds.makeSet(node.ID())
	}

	for _, edge := range edges {
		// The disjoint set doesn't really care for which is head and which is tail so this
		// should work fine without checking both ways
		if s1, s2 := ds.find(edge.Edge.From().ID()), ds.find(edge.Edge.To().ID()); s1 != s2 {
			ds.union(s1, s2)
			dst.AddUndirectedEdge(edge.Edge, edge.Cost)
		}
	}
}

/* Control flow graph stuff */

// A dominates B if and only if the only path through B travels through A.
//
// This returns all possible dominators for all nodes, it does not prune for strict dominators,
// immediate dominators etc.
//
func Dominators(start graph.Node, g graph.Graph) map[int]Set {
	allNodes := make(Set)
	nlist := g.Nodes()
	dominators := make(map[int]Set, len(nlist))
	for _, node := range nlist {
		allNodes.add(node)
	}

	var to func(graph.Node) []graph.Node
	switch g := g.(type) {
	case graph.Directed:
		to = g.To
	default:
		to = g.From
	}

	for _, node := range nlist {
		dominators[node.ID()] = make(Set)
		if node.ID() == start.ID() {
			dominators[node.ID()].add(start)
		} else {
			dominators[node.ID()].copy(allNodes)
		}
	}

	for somethingChanged := true; somethingChanged; {
		somethingChanged = false
		for _, node := range nlist {
			if node.ID() == start.ID() {
				continue
			}
			preds := to(node)
			if len(preds) == 0 {
				continue
			}
			tmp := make(Set).copy(dominators[preds[0].ID()])
			for _, pred := range preds[1:] {
				tmp.intersect(tmp, dominators[pred.ID()])
			}

			dom := make(Set)
			dom.add(node)

			dom.union(dom, tmp)
			if !equal(dom, dominators[node.ID()]) {
				dominators[node.ID()] = dom
				somethingChanged = true
			}
		}
	}

	return dominators
}

// A Postdominates B if and only if all paths from B travel through A.
//
// This returns all possible post-dominators for all nodes, it does not prune for strict
// postdominators, immediate postdominators etc.
func PostDominators(end graph.Node, g graph.Graph) map[int]Set {
	allNodes := make(Set)
	nlist := g.Nodes()
	dominators := make(map[int]Set, len(nlist))
	for _, node := range nlist {
		allNodes.add(node)
	}

	for _, node := range nlist {
		dominators[node.ID()] = make(Set)
		if node.ID() == end.ID() {
			dominators[node.ID()].add(end)
		} else {
			dominators[node.ID()].copy(allNodes)
		}
	}

	for somethingChanged := true; somethingChanged; {
		somethingChanged = false
		for _, node := range nlist {
			if node.ID() == end.ID() {
				continue
			}
			succs := g.From(node)
			if len(succs) == 0 {
				continue
			}
			tmp := make(Set).copy(dominators[succs[0].ID()])
			for _, succ := range succs[1:] {
				tmp.intersect(tmp, dominators[succ.ID()])
			}

			dom := make(Set)
			dom.add(node)

			dom.union(dom, tmp)
			if !equal(dom, dominators[node.ID()]) {
				dominators[node.ID()] = dom
				somethingChanged = true
			}
		}
	}

	return dominators
}

// VertexOrdering returns the vertex ordering and the k-cores of
// the undirected graph g.
func VertexOrdering(g graph.Graph) (order []graph.Node, cores [][]graph.Node) {
	nodes := g.Nodes()

	// The algorithm used here is essentially as described at
	// http://en.wikipedia.org/w/index.php?title=Degeneracy_%28graph_theory%29&oldid=640308710

	// Initialize an output list L.
	var l []graph.Node

	// Compute a number d_v for each vertex v in G,
	// the number of neighbors of v that are not already in L.
	// Initially, these numbers are just the degrees of the vertices.
	dv := make(map[int]int, len(nodes))
	var (
		maxDegree  int
		neighbours = make(map[int][]graph.Node)
	)
	for _, n := range nodes {
		adj := g.From(n)
		neighbours[n.ID()] = adj
		dv[n.ID()] = len(adj)
		if len(adj) > maxDegree {
			maxDegree = len(adj)
		}
	}

	// Initialize an array D such that D[i] contains a list of the
	// vertices v that are not already in L for which d_v = i.
	d := make([][]graph.Node, maxDegree+1)
	for _, n := range nodes {
		deg := dv[n.ID()]
		d[deg] = append(d[deg], n)
	}

	// Initialize k to 0.
	k := 0
	// Repeat n times:
	s := []int{0}
	for _ = range nodes { // TODO(kortschak): Remove blank assignment when go1.3.3 is no longer supported.
		// Scan the array cells D[0], D[1], ... until
		// finding an i for which D[i] is nonempty.
		var (
			i  int
			di []graph.Node
		)
		for i, di = range d {
			if len(di) != 0 {
				break
			}
		}

		// Set k to max(k,i).
		if i > k {
			k = i
			s = append(s, make([]int, k-len(s)+1)...)
		}

		// Select a vertex v from D[i]. Add v to the
		// beginning of L and remove it from D[i].
		var v graph.Node
		v, d[i] = di[len(di)-1], di[:len(di)-1]
		l = append(l, v)
		s[k]++
		delete(dv, v.ID())

		// For each neighbor w of v not already in L,
		// subtract one from d_w and move w to the
		// cell of D corresponding to the new value of d_w.
		for _, w := range neighbours[v.ID()] {
			dw, ok := dv[w.ID()]
			if !ok {
				continue
			}
			for i, n := range d[dw] {
				if n.ID() == w.ID() {
					d[dw][i], d[dw] = d[dw][len(d[dw])-1], d[dw][:len(d[dw])-1]
					dw--
					d[dw] = append(d[dw], w)
					break
				}
			}
			dv[w.ID()] = dw
		}
	}

	for i, j := 0, len(l)-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}
	cores = make([][]graph.Node, len(s))
	offset := len(l)
	for i, n := range s {
		cores[i] = l[offset-n : offset]
		offset -= n
	}
	return l, cores
}

// BronKerbosch returns the set of maximal cliques of the undirected graph g.
func BronKerbosch(g graph.Graph) [][]graph.Node {
	nodes := g.Nodes()

	// The algorithm used here is essentially BronKerbosch3 as described at
	// http://en.wikipedia.org/w/index.php?title=Bron%E2%80%93Kerbosch_algorithm&oldid=656805858

	p := make(Set, len(nodes))
	for _, n := range nodes {
		p.add(n)
	}
	x := make(Set)
	var bk bronKerbosch
	order, _ := VertexOrdering(g)
	for _, v := range order {
		neighbours := g.From(v)
		nv := make(Set, len(neighbours))
		for _, n := range neighbours {
			nv.add(n)
		}
		bk.maximalCliquePivot(g, []graph.Node{v}, make(Set).intersect(p, nv), make(Set).intersect(x, nv))
		p.remove(v)
		x.add(v)
	}
	return bk
}

type bronKerbosch [][]graph.Node

func (bk *bronKerbosch) maximalCliquePivot(g graph.Graph, r []graph.Node, p, x Set) {
	if len(p) == 0 && len(x) == 0 {
		*bk = append(*bk, r)
		return
	}

	neighbours := bk.choosePivotFrom(g, p, x)
	nu := make(Set, len(neighbours))
	for _, n := range neighbours {
		nu.add(n)
	}
	for _, v := range p {
		if nu.has(v) {
			continue
		}
		neighbours := g.From(v)
		nv := make(Set, len(neighbours))
		for _, n := range neighbours {
			nv.add(n)
		}

		var found bool
		for _, n := range r {
			if n.ID() == v.ID() {
				found = true
				break
			}
		}
		var sr []graph.Node
		if !found {
			sr = append(r[:len(r):len(r)], v)
		}

		bk.maximalCliquePivot(g, sr, make(Set).intersect(p, nv), make(Set).intersect(x, nv))
		p.remove(v)
		x.add(v)
	}
}

func (*bronKerbosch) choosePivotFrom(g graph.Graph, p, x Set) (neighbors []graph.Node) {
	// TODO(kortschak): Investigate the impact of pivot choice that maximises
	// |p ⋂ neighbours(u)| as a function of input size. Until then, leave as
	// compile time option.
	if !tomitaTanakaTakahashi {
		for _, n := range p {
			return g.From(n)
		}
		for _, n := range x {
			return g.From(n)
		}
		panic("bronKerbosch: empty set")
	}

	var (
		max   = -1
		pivot graph.Node
	)
	maxNeighbors := func(s Set) {
	outer:
		for _, u := range s {
			nb := g.From(u)
			c := len(nb)
			if c <= max {
				continue
			}
			for n := range nb {
				if _, ok := p[n]; ok {
					continue
				}
				c--
				if c <= max {
					continue outer
				}
			}
			max = c
			pivot = u
			neighbors = nb
		}
	}
	maxNeighbors(p)
	maxNeighbors(x)
	if pivot == nil {
		panic("bronKerbosch: empty set")
	}
	return neighbors
}

// ConnectedComponents returns the connected components of the graph g. All
// edges are treated as undirected.
func ConnectedComponents(g graph.Undirected) [][]graph.Node {
	var (
		w  traverse.DepthFirst
		c  []graph.Node
		cc [][]graph.Node
	)
	during := func(n graph.Node) {
		c = append(c, n)
	}
	after := func() {
		cc = append(cc, []graph.Node(nil))
		cc[len(cc)-1] = append(cc[len(cc)-1], c...)
		c = c[:0]
	}
	w.WalkAll(g, nil, after, during)

	return cc
}
