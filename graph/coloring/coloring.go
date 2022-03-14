// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// See Lewis, A Guide to Graph Colouring: Algorithms and Applications
// doi:10.1007/978-3-319-25730-3 for significant discussion of approaches.

package coloring

import (
	"errors"
	"sort"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/set"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/topo"
)

// ErrInvalidPartialColoring is returned when a partial coloring
// is provided for a graph with inadmissible color assignments.
var ErrInvalidPartialColoring = errors.New("coloring: invalid partial coloring")

// Sets returns the mapping from colors to sets of node IDs. Each set of
// node IDs is sorted by ascending value.
func Sets(colors map[int64]int) map[int][]int64 {
	sets := make(map[int][]int64)
	for id, c := range colors {
		sets[c] = append(sets[c], id)
	}
	for _, s := range sets {
		sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
	}
	return sets
}

// Dsatur returns an approximate minimal chromatic number of g and a
// corresponding vertex coloring using the heuristic Dsatur coloring algorithm.
// If a partial coloring is provided the coloring will be consistent with
// that partial coloring if possible. Otherwise Dsatur will return
// ErrInvalidPartialColoring.
// See Brélaz doi:10.1145/359094.359101 for details of the algorithm.
func Dsatur(g graph.Undirected, partial map[int64]int) (k int, colors map[int64]int, err error) {
	nodes := g.Nodes()
	n := nodes.Len()
	if n == 0 {
		return
	}
	partial, ok := newPartial(partial, g)
	if !ok {
		return -1, nil, ErrInvalidPartialColoring
	}
	order := bySaturationDegree(nodes, g, partial)
	order.heuristic = order.dsatur
	k, colors = greedyColoringOf(g, order, order.colors)
	return k, colors, nil
}

// Terminator is a cancellation-only context type. A context.Context
// may be used as a Terminator.
type Terminator interface {
	// Done returns a channel that is closed when work
	// should be terminated. Done may return nil if this
	// work can never be canceled.
	// Successive calls to Done should return the same value.
	Done() <-chan struct{}

	// If Done is not yet closed, Err returns nil.
	// If Done is closed, Err returns a non-nil error
	// explaining why.
	// After Err returns a non-nil error, successive
	// calls to Err should return the same error.
	Err() error
}

// DsaturExact returns the exact minimal chromatic number of g and a
// corresponding vertex coloring using the branch-and-bound Dsatur coloring
// algorithm of Brélaz. If the provided terminator is cancelled or times out
// before completion, the terminator's reason for termination will be returned
// along with a potentially sub-optimal chromatic number and coloring. If
// term is nil, DsaturExact will run to completion.
// See Brélaz doi:10.1145/359094.359101 for details of the algorithm.
func DsaturExact(term Terminator, g graph.Undirected) (k int, colors map[int64]int, err error) {
	// This is implemented essentially as described in algorithm 1 of
	// doi:10.1002/net.21716 with the exception that we obtain a
	// tighter upper bound by doing a single run of an approximate
	// Brélaz Dsatur coloring, using the result if the recurrence is
	// cancelled.
	// We also use the initial maximum clique as a starting point for
	// the exact search. If there is more than one maximum clique, we
	// need to ensure that we pick the one that will lead us down the
	// easiest branch of the search tree. This will be the maximum
	// clique with the lowest degree into the remainder of the graph.

	nodes := g.Nodes()
	n := nodes.Len()
	if n == 0 {
		return
	}

	lb, maxClique, cliques := maximumClique(g)
	if lb == n {
		return lb, colorClique(maxClique), nil
	}

	order := bySaturationDegree(nodes, g, make(map[int64]int))
	order.heuristic = order.dsatur
	// Find initial coloring via Dsatur heuristic.
	ub, initial := greedyColoringOf(g, order, order.colors)
	if lb == ub {
		return ub, initial, nil
	}

	selector := &order.saturationDegree
	cand := newDsaturColoring(order.nodes, bestMaximumClique(g, cliques))
	k, colors, err = dSaturExact(term, selector, cand, len(cand.colors), ub, nil)
	if colors == nil {
		return ub, initial, err
	}
	if k == lb {
		err = nil
	}
	return k, colors, err
}

// dSaturColoring is a partial graph coloring.
type dSaturColoring struct {
	colors    map[int64]int
	uncolored set.Int64s
}

// newDsaturColoring returns a dSaturColoring representing a partial coloring
// of a graph with the given nodes and colors.
func newDsaturColoring(nodes []graph.Node, colors map[int64]int) dSaturColoring {
	uncolored := make(set.Int64s)
	for _, v := range nodes {
		vid := v.ID()
		if _, ok := colors[vid]; !ok {
			uncolored.Add(vid)
		}
	}
	return dSaturColoring{
		colors:    colors,
		uncolored: uncolored,
	}
}

// color moves a node from the uncolored set to the colored set.
func (c dSaturColoring) color(id int64) {
	if !c.uncolored.Has(id) {
		if _, ok := c.colors[id]; ok {
			panic("coloring: coloring already colored node")
		}
		panic("coloring: coloring non-existent node")
	}
	// The node has its uncolored mark removed, but is
	// not explicitly colored until the dSaturExact
	// caller has completed its recursive exploration
	// of the feasible colors.
	c.uncolored.Remove(id)
}

// uncolor moves a node from the colored set to the uncolored set.
func (c dSaturColoring) uncolor(id int64) {
	if _, ok := c.colors[id]; !ok {
		if c.uncolored.Has(id) {
			panic("coloring: uncoloring already uncolored node")
		}
		panic("coloring: uncoloring non-existent node")
	}
	delete(c.colors, id)
	c.uncolored.Add(id)
}

// dSaturExact recursively searches for an exact minimum vertex coloring of the
// full graph in cand. If no chromatic number lower than ub is found, colors is
// returned as nil.
func dSaturExact(term Terminator, selector *saturationDegree, cand dSaturColoring, k, ub int, best map[int64]int) (newK int, colors map[int64]int, err error) {
	if len(cand.uncolored) == 0 {
		// In the published algorithm, this is guarded by k < ub,
		// but dSaturExact is never called with k >= ub; in the
		// initial call we have excluded cases where k == ub and
		// it cannot be greater, and in the recursive call, we
		// have already checked that k < ub.
		return k, clone(cand.colors), nil
	}

	if term != nil {
		select {
		case <-term.Done():
			if best == nil {
				return -1, nil, term.Err()
			}
			colors := make(set.Ints)
			for _, c := range best {
				colors.Add(c)
			}
			return colors.Count(), best, term.Err()
		default:
		}
	}

	// Select the next node.
	selector.reset(cand.colors)
	vid := selector.nodes[selector.dsatur()].ID()
	cand.color(vid)
	// If uncolor panics, we have failed to find a
	// feasible color. This should never happen.
	defer cand.uncolor(vid)

	// Keep the adjacent colors set as it will be
	// overwritten by child recurrences.
	adjColors := selector.adjColors[selector.indexOf[vid]]

	// Collect all feasible existing colors plus one, remembering it.
	feasible := make(set.Ints)
	for _, c := range cand.colors {
		if adjColors.Has(c) {
			continue
		}
		feasible.Add(c)
	}
	var newCol int
	for c := 0; c < ub; c++ {
		if feasible.Has(c) || adjColors.Has(c) {
			continue
		}
		feasible.Add(c)
		newCol = c
		break
	}

	// Recur over every feasible color.
	for c := range feasible {
		cand.colors[vid] = c
		effK := k
		if c == newCol {
			effK++
		}
		// In the published algorithm, the expression max(effK, lb) < ub is
		// used, but lb < ub always since it is not updated and dSaturExact
		// is not called if lb == ub, and it cannot be greater.
		if effK < ub {
			ub, best, err = dSaturExact(term, selector, cand, effK, ub, best)
			if err != nil {
				return ub, best, err
			}
		}
	}

	return ub, best, nil
}

// maximumClique returns a maximum clique in g and its order.
func maximumClique(g graph.Undirected) (k int, maxClique []graph.Node, cliques [][]graph.Node) {
	cliques = topo.BronKerbosch(g)
	for _, c := range topo.BronKerbosch(g) {
		if len(c) > len(maxClique) {
			maxClique = c
		}
	}
	return len(maxClique), maxClique, cliques
}

// bestMaximumClique returns the maximum clique in g with the lowest degree into
// the remainder of the graph.
func bestMaximumClique(g graph.Undirected, cliques [][]graph.Node) (colors map[int64]int) {
	switch len(cliques) {
	case 0:
		return nil
	case 1:
		return colorClique(cliques[0])
	}

	sort.Slice(cliques, func(i, j int) bool { return len(cliques[i]) > len(cliques[j]) })
	maxClique := cliques[0]
	minDegree := cliqueDegree(g, maxClique)
	for _, c := range cliques[1:] {
		if len(c) < len(maxClique) {
			break
		}
		d := cliqueDegree(g, c)
		if d < minDegree {
			minDegree = d
			maxClique = c
		}
	}

	return colorClique(maxClique)
}

// cliqueDegree returns the degree of the clique to nodes outside the clique.
func cliqueDegree(g graph.Undirected, clique []graph.Node) int {
	n := make(set.Int64s)
	for _, u := range clique {
		to := g.From(u.ID())
		for to.Next() {
			n.Add(to.Node().ID())
		}
	}
	return n.Count() - len(clique)
}

// colorClique returns a valid coloring for the given clique.
func colorClique(clique []graph.Node) map[int64]int {
	colors := make(map[int64]int, len(clique))
	for c, u := range clique {
		colors[u.ID()] = c
	}
	return colors
}

// Randomized returns an approximate minimal chromatic number of g and a
// corresponding vertex coloring using a greedy coloring algorithm with a
// random node ordering. If src is non-nil it will be used as the random
// source, otherwise the global random source will be used. If a partial
// coloring is provided the coloring will be consistent with that partial
// coloring if possible. Otherwise Randomized will return
// ErrInvalidPartialColoring.
func Randomized(g graph.Undirected, partial map[int64]int, src rand.Source) (k int, colors map[int64]int, err error) {
	nodes := g.Nodes()
	n := nodes.Len()
	if n == 0 {
		return
	}
	partial, ok := newPartial(partial, g)
	if !ok {
		return -1, nil, ErrInvalidPartialColoring
	}
	k, colors = greedyColoringOf(g, randomize(nodes, src), partial)
	return k, colors, nil
}

// randomize returns a graph.Node iterator that returns nodes in a random
// order.
func randomize(it graph.Nodes, src rand.Source) graph.Nodes {
	nodes := graph.NodesOf(it)
	var shuffle func(int, func(i, j int))
	if src == nil {
		shuffle = rand.Shuffle
	} else {
		shuffle = rand.New(src).Shuffle
	}
	shuffle(len(nodes), func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})
	return iterator.NewOrderedNodes(nodes)
}

// RecursiveLargestFirst returns an approximate minimal chromatic number
// of g and a corresponding vertex coloring using the Recursive Largest
// First coloring algorithm.
// See Leighton doi:10.6028/jres.084.024 for details of the algorithm.
func RecursiveLargestFirst(g graph.Undirected) (k int, colors map[int64]int) {
	it := g.Nodes()
	n := it.Len()
	if n == 0 {
		return
	}
	nodes := graph.NodesOf(it)
	colors = make(map[int64]int)

	// The names of variable here have been changed from the original PL-1
	// for clarity, but the correspondence is as follows:
	//  E   -> isolated
	//  F   -> boundary
	//  L   -> current
	//  COL -> k
	//  C   -> colors

	// Initialize the boundary vector to the node degrees.
	boundary := make([]int, len(nodes))
	indexOf := make(map[int64]int)
	for i, u := range nodes {
		uid := u.ID()
		indexOf[uid] = i
		boundary[i] = g.From(uid).Len()
	}
	deleteFrom := func(vec []int, idx int) {
		vec[idx] = -1
		to := g.From(nodes[idx].ID())
		for to.Next() {
			vec[indexOf[to.Node().ID()]]--
		}
	}
	isolated := make([]int, len(nodes))

	// If there are any uncolored nodes, initiate the assignment of the next color.
	// Incrementing color happens at the end of the loop in this implementation.
	var current int
	for j := 0; j < n; {
		// Reinitialize the isolated vector.
		copy(isolated, boundary)

		// Select the node in U₁ with maximal degree in U₁.
		for i := range nodes {
			if boundary[i] > boundary[current] {
				current = i
			}
		}

		// Color the node just selected and continue to
		// color nodes with color k until U₁ is empty.
		for isolated[current] >= 0 {
			// Color node and modify U₁ and U₂ accordingly.
			deleteFrom(isolated, current)
			deleteFrom(boundary, current)
			colors[nodes[current].ID()] = k
			j++
			to := g.From(nodes[current].ID())
			for to.Next() {
				i := indexOf[to.Node().ID()]
				if isolated[i] >= 0 {
					deleteFrom(isolated, i)
				}
			}

			// Find the first node in U₁, if any.
			for i := range nodes {
				if isolated[i] < 0 {
					continue
				}

				// If U₁ is not empty, select the next node for coloring.
				current = i
				currAvail := boundary[current] - isolated[current]
				for j := i; j < n; j++ {
					if isolated[j] < 0 {
						continue
					}
					nextAvail := boundary[j] - isolated[j]
					switch {
					case nextAvail > currAvail, nextAvail == currAvail && isolated[j] < isolated[current]:
						current = j
						currAvail = boundary[current] - isolated[current]
					}
				}
				break
			}
		}

		k++
	}

	return k, colors
}

// SanSegundo returns an approximate minimal chromatic number of g and a
// corresponding vertex coloring using the PASS rule with a single run of a
// greedy coloring algorithm. If a partial coloring is provided the coloring
// will be consistent with that partial coloring if possible. Otherwise
// SanSegundo will return ErrInvalidPartialColoring.
// See San Segundo doi:10.1016/j.cor.2011.10.008 for details of the algorithm.
func SanSegundo(g graph.Undirected, partial map[int64]int) (k int, colors map[int64]int, err error) {
	nodes := g.Nodes()
	n := nodes.Len()
	if n == 0 {
		return
	}
	partial, ok := newPartial(partial, g)
	if !ok {
		return -1, nil, ErrInvalidPartialColoring
	}
	order := bySaturationDegree(nodes, g, partial)
	order.heuristic = order.pass
	k, colors = greedyColoringOf(g, order, order.colors)
	return k, colors, nil
}

// WelshPowell returns an approximate minimal chromatic number of g and a
// corresponding vertex coloring using the Welsh and Powell coloring algorithm.
// If a partial coloring is provided the coloring will be consistent with that
// partial coloring if possible. Otherwise WelshPowell will return
// ErrInvalidPartialColoring.
// See Welsh and Powell doi:10.1093/comjnl/10.1.85 for details of the algorithm.
func WelshPowell(g graph.Undirected, partial map[int64]int) (k int, colors map[int64]int, err error) {
	nodes := g.Nodes()
	n := nodes.Len()
	if n == 0 {
		return
	}
	partial, ok := newPartial(partial, g)
	if !ok {
		return -1, nil, ErrInvalidPartialColoring
	}
	k, colors = greedyColoringOf(g, byDescendingDegree(nodes, g), partial)
	return k, colors, nil
}

// byDescendingDegree returns a graph.Node iterator that returns nodes
// in order of descending degree.
func byDescendingDegree(it graph.Nodes, g graph.Undirected) graph.Nodes {
	nodes := graph.NodesOf(it)
	n := byDescDegree{nodes: nodes, degrees: make([]int, len(nodes))}
	for i, u := range nodes {
		n.degrees[i] = g.From(u.ID()).Len()
	}
	sort.Sort(n)
	return iterator.NewOrderedNodes(nodes)
}

// byDescDegree sorts a slice of graph.Node descending by the corresponding
// value of the degrees slice.
type byDescDegree struct {
	nodes   []graph.Node
	degrees []int
}

func (n byDescDegree) Len() int           { return len(n.nodes) }
func (n byDescDegree) Less(i, j int) bool { return n.degrees[i] > n.degrees[j] }
func (n byDescDegree) Swap(i, j int) {
	n.nodes[i], n.nodes[j] = n.nodes[j], n.nodes[i]
	n.degrees[i], n.degrees[j] = n.degrees[j], n.degrees[i]
}

// newPartial returns a new valid partial coloring is valid for g. An empty
// partial coloring is valid. If the partial coloring is not valid, a nil map
// is returned, otherwise a new non-nil map is returned. If the input partial
// coloring is nil, a new map is created and returned.
func newPartial(partial map[int64]int, g graph.Undirected) (map[int64]int, bool) {
	if partial == nil {
		return make(map[int64]int), true
	}
	for id, c := range partial {
		if g.Node(id) == nil {
			return nil, false
		}
		to := g.From(id)
		for to.Next() {
			if oc, ok := partial[to.Node().ID()]; ok && c == oc {
				return nil, false
			}
		}
	}
	return clone(partial), true
}

func clone(colors map[int64]int) map[int64]int {
	new := make(map[int64]int, len(colors))
	for id, c := range colors {
		new[id] = c
	}
	return new
}

// greedyColoringOf returns the chromatic number and a graph coloring of g
// based on the sequential coloring of nodes given by order and starting from
// the given partial coloring.
func greedyColoringOf(g graph.Undirected, order graph.Nodes, partial map[int64]int) (k int, colors map[int64]int) {
	colors = partial
	constrained := false
	for _, c := range colors {
		if c > k {
			k = c
			constrained = true
		}
	}

	// Next nodes are chosen by the specified heuristic in order.
	for order.Next() {
		uid := order.Node().ID()
		used := colorsOf(g.From(uid), colors)
		if c, ok := colors[uid]; ok {
			if used.Has(c) {
				return -1, nil
			}
			continue
		}
		// Color the chosen vertex with the least possible
		// (lowest numbered) color.
		for c := 0; c <= k+1; c++ {
			if !used.Has(c) {
				colors[uid] = c
				if c > k {
					k = c
				}
				break
			}
		}
	}

	if !constrained {
		return k + 1, colors
	}
	seen := make(set.Ints)
	for _, c := range colors {
		seen.Add(c)
	}
	return seen.Count(), colors
}

// colorsOf returns all the colors in the coloring that are used by the
// given nodes.
func colorsOf(nodes graph.Nodes, coloring map[int64]int) set.Ints {
	c := make(set.Ints, nodes.Len())
	for nodes.Next() {
		used, ok := coloring[nodes.Node().ID()]
		if ok {
			c.Add(used)
		}
	}
	return c
}

// saturationDegreeIterator is a graph.Nodes iterator that returns nodes ordered
// by decreasing saturation degree.
type saturationDegreeIterator struct {
	// cnt is the number of nodes that
	// have been returned and curr is
	// the current selection.
	cnt, curr int

	// heuristic determines the
	// iterator's node selection
	// heuristic. It can be either
	// saturationDegree.dsatur or
	// saturationDegree.pass.
	heuristic func() int

	saturationDegree
}

// bySaturationDegree returns a new saturationDegreeIterator that will
// iterate over the node in it based on the given graph and partial coloring.
// The saturationDegreeIterator holds a reference to colors allowing
// greedyColoringOf to update its coloring.
func bySaturationDegree(it graph.Nodes, g graph.Undirected, colors map[int64]int) *saturationDegreeIterator {
	return &saturationDegreeIterator{
		cnt: -1, curr: -1,
		saturationDegree: newSaturationDegree(it, g, colors),
	}
}

// Len returns the number of elements remaining in the iterator.
// saturationDegreeIterator is an indeterminate iterator, so Len always
// returns -1. This is required to satisfy the graph.Iterator interface.
func (n *saturationDegreeIterator) Len() int { return -1 }

// Next advances the iterator to the next node and returns whether
// the next call to the Node method will return a valid Node.
func (n *saturationDegreeIterator) Next() bool {
	if uint(n.cnt)+1 < uint(len(n.nodes)) {
		n.cnt++
		switch n.cnt {
		case 0:
			max := -1
			for i, d := range n.degrees {
				if d > max {
					max = d
					n.curr = i
				}
			}
		default:
			prev := n.Node().ID()
			c := n.colors[prev]
			to := n.g.From(prev)
			for to.Next() {
				n.adjColors[n.indexOf[to.Node().ID()]].Add(c)
			}

			chosen := n.heuristic()
			if chosen < 0 || chosen == n.curr {
				return false
			}
			n.curr = chosen
		}
		return true
	}
	n.cnt = len(n.nodes)
	return false
}

// Node returns the current node.
func (n *saturationDegreeIterator) Node() graph.Node { return n.nodes[n.curr] }

// Reset implements the graph.Iterator interface. It should not be called.
func (n *saturationDegreeIterator) Reset() { panic("coloring: invalid call to Reset") }

// saturationDegree is a saturation degree node choice heuristic.
type saturationDegree struct {
	// nodes is the set of nodes being
	// iterated over.
	nodes []graph.Node

	// indexOf is a mapping between node
	// IDs and elements of degree and
	// adjColors. degrees holds the
	// degree of each node and adjColors
	// holds the current adjacent
	// colors of each node.
	indexOf   map[int64]int
	degrees   []int
	adjColors []set.Ints

	// g and colors are the graph coloring.
	// colors is held by both the iterator
	// and greedyColoringOf.
	g      graph.Undirected
	colors map[int64]int

	// work is a temporary workspace.
	work []int
}

// newSaturationDegree returns a saturationDegree heuristic based on the
// nodes in the given node iterator and graph, using the provided coloring.
func newSaturationDegree(it graph.Nodes, g graph.Undirected, colors map[int64]int) saturationDegree {
	nodes := graph.NodesOf(it)
	sd := saturationDegree{
		nodes:     nodes,
		indexOf:   make(map[int64]int, len(nodes)),
		degrees:   make([]int, len(nodes)),
		adjColors: make([]set.Ints, len(nodes)),
		g:         g,
		colors:    colors,
	}
	for i, u := range nodes {
		sd.degrees[i] = g.From(u.ID()).Len()
		sd.adjColors[i] = make(set.Ints)
		sd.indexOf[u.ID()] = i
	}
	for uid, c := range colors {
		to := g.From(uid)
		for to.Next() {
			sd.adjColors[sd.indexOf[to.Node().ID()]].Add(c)
		}
	}
	return sd
}

// reset re-initializes the saturation with the provided colors.
func (sd *saturationDegree) reset(colors map[int64]int) {
	sd.colors = colors
	for i := range sd.nodes {
		sd.adjColors[i] = make(set.Ints)
	}
	for uid, c := range colors {
		to := sd.g.From(uid)
		for to.Next() {
			sd.adjColors[sd.indexOf[to.Node().ID()]].Add(c)
		}
	}
}

// dsatur implements the Dsatur heuristic from Brélaz doi:10.1145/359094.359101.
func (sd *saturationDegree) dsatur() int {
	maxSat, maxDeg, chosen := -1, -1, -1
	for i, u := range sd.nodes {
		uid := u.ID()
		if _, ok := sd.colors[uid]; ok {
			continue
		}
		s := saturationDegreeOf(uid, sd.g, sd.colors)
		d := sd.degrees[i]
		switch {
		// Choose a vertex with a maximal saturation degree.
		// If there is an equality, choose any vertex of maximal
		// degree in the uncolored subgraph.
		case s > maxSat, s == maxSat && d > maxDeg:
			maxSat = s
			maxDeg = d
			chosen = i
		}
	}
	return chosen
}

// pass implements the PASS heuristic from San Segundo doi:10.1016/j.cor.2011.10.008.
func (sd *saturationDegree) pass() int {
	maxSat, chosen := -1, -1
	sd.work = sd.work[:0]
	for i, u := range sd.nodes {
		uid := u.ID()
		if _, ok := sd.colors[uid]; ok {
			continue
		}
		s := saturationDegreeOf(uid, sd.g, sd.colors)
		switch {
		case s > maxSat:
			maxSat = s
			sd.work = sd.work[:0]
			fallthrough
		case s == maxSat:
			sd.work = append(sd.work, i)
		}
	}
	maxAvail := -1
	for _, vs := range sd.work {
		var avail int
		for _, v := range sd.work {
			if v != vs {
				avail += sd.same(sd.adjColors[vs], sd.adjColors[v])
			}
		}
		switch {
		case avail > maxAvail, avail == maxAvail && sd.nodes[chosen].ID() < sd.nodes[vs].ID():
			maxAvail = avail
			chosen = vs
		}
	}
	return chosen
}

// same implements the same function from San Segundo doi:10.1016/j.cor.2011.10.008.
func (sd *saturationDegree) same(vi, vj set.Ints) int {
	valid := make(set.Ints)
	for _, c := range sd.colors {
		if !vi.Has(c) {
			valid.Add(c)
		}
	}
	for c := range vj {
		valid.Remove(c)
	}
	return valid.Count()
}

// saturationDegreeOf returns the saturation degree of the node corresponding to
// vid in g with the given coloring.
func saturationDegreeOf(vid int64, g graph.Undirected, colors map[int64]int) int {
	if _, ok := colors[vid]; ok {
		panic("coloring: saturation degree not defined for colored node")
	}
	adjColors := make(set.Ints)
	to := g.From(vid)
	for to.Next() {
		if c, ok := colors[to.Node().ID()]; ok {
			adjColors.Add(c)
		}
	}
	return adjColors.Count()
}
