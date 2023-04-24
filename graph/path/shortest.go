// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"math"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats/scalar"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/ordered"
	"gonum.org/v1/gonum/graph/internal/set"
	"gonum.org/v1/gonum/mat"
)

// Shortest is a shortest-path tree created by the BellmanFordFrom, DijkstraFrom
// or AStar single-source shortest path functions.
type Shortest struct {
	// from holds the source node given to
	// the function that returned the
	// Shortest value.
	from graph.Node

	// nodes hold the nodes of the analysed
	// graph.
	nodes []graph.Node
	// indexOf contains a mapping between
	// the id-dense representation of the
	// graph and the potentially id-sparse
	// nodes held in nodes.
	indexOf map[int64]int

	// dist and next represent the shortest
	// paths between nodes.
	//
	// Indices into dist and next are
	// mapped through indexOf.
	//
	// dist contains the distances
	// from the from node for each
	// node in the graph.
	dist []float64
	// next contains the shortest-path
	// tree of the graph. The index is a
	// linear mapping of to-dense-id.
	next []int

	// hasNegativeCycle indicates
	// whether the Shortest includes
	// a negative cycle. This should
	// be set by the function that
	// returned the Shortest value.
	hasNegativeCycle bool

	// negCosts holds negative costs
	// between pairs of nodes to report
	// negative cycles.
	// negCosts must be initialised by
	// routines that can handle negative
	// edge weights.
	negCosts map[negEdge]float64
}

// newShortestFrom returns a shortest path tree for paths from u
// initialised with the given nodes. The nodes held by the returned
// Shortest may be lazily added.
func newShortestFrom(u graph.Node, nodes []graph.Node) Shortest {
	indexOf := make(map[int64]int, len(nodes))
	uid := u.ID()
	for i, n := range nodes {
		indexOf[n.ID()] = i
		if n.ID() == uid {
			u = n
		}
	}

	p := Shortest{
		from: u,

		nodes:   nodes,
		indexOf: indexOf,

		dist: make([]float64, len(nodes)),
		next: make([]int, len(nodes)),
	}
	for i := range nodes {
		p.dist[i] = math.Inf(1)
		p.next[i] = -1
	}
	p.dist[indexOf[uid]] = 0

	return p
}

// add adds a node to the Shortest, initialising its stored index and returning, and
// setting the distance and position as unconnected. add will panic if the node is
// already present.
func (p *Shortest) add(u graph.Node) int {
	uid := u.ID()
	if _, exists := p.indexOf[uid]; exists {
		panic("shortest: adding existing node")
	}
	idx := len(p.nodes)
	p.indexOf[uid] = idx
	p.nodes = append(p.nodes, u)
	p.dist = append(p.dist, math.Inf(1))
	p.next = append(p.next, -1)
	return idx
}

// set sets the weight of the path from the node in p.nodes indexed by mid to the node
// indexed by to.
func (p Shortest) set(to int, weight float64, mid int) {
	p.dist[to] = weight
	p.next[to] = mid
	if weight < 0 {
		e := negEdge{from: mid, to: to}
		c, ok := p.negCosts[e]
		if !ok {
			p.negCosts[e] = weight
		} else if weight < c {
			// The only ways that we can have a new weight that is
			// lower than the previous weight is if either the edge
			// has already been traversed in a negative cycle, or
			// the edge is reachable from a negative cycle.
			// Either way the reported path is returned with a
			// negative infinite path weight.
			p.negCosts[e] = math.Inf(-1)
		}
	}
}

// From returns the starting node of the paths held by the Shortest.
func (p Shortest) From() graph.Node { return p.from }

// WeightTo returns the weight of the minimum path to v. If the path to v includes
// a negative cycle, the returned weight will not reflect the true path weight.
func (p Shortest) WeightTo(vid int64) float64 {
	to, toOK := p.indexOf[vid]
	if !toOK {
		return math.Inf(1)
	}
	return p.dist[to]
}

// To returns a shortest path to v and the weight of the path. If the path
// to v includes a negative cycle, one pass through the cycle will be included
// in path, but any path leading into the negative cycle will be lost, and
// weight will be returned as -Inf.
func (p Shortest) To(vid int64) (path []graph.Node, weight float64) {
	to, toOK := p.indexOf[vid]
	if !toOK || math.IsInf(p.dist[to], 1) {
		return nil, math.Inf(1)
	}
	from := p.indexOf[p.from.ID()]
	path = []graph.Node{p.nodes[to]}
	weight = math.Inf(1)
	if p.hasNegativeCycle {
		seen := make(set.Ints)
		seen.Add(from)
		for to != from {
			next := p.next[to]
			if math.IsInf(p.negCosts[negEdge{from: next, to: to}], -1) {
				weight = math.Inf(-1)
			}
			if seen.Has(to) {
				break
			}
			seen.Add(to)
			path = append(path, p.nodes[next])
			to = next
		}
	} else {
		n := len(p.nodes)
		for to != from {
			to = p.next[to]
			path = append(path, p.nodes[to])
			if n < 0 {
				panic("path: unexpected negative cycle")
			}
			n--
		}
	}
	ordered.Reverse(path)
	return path, math.Min(weight, p.dist[p.indexOf[vid]])
}

// ShortestAlts is a shortest-path tree created by the BellmanFordAllFrom or DijkstraAllFrom
// single-source shortest path functions.
type ShortestAlts struct {
	// from holds the source node given to
	// the function that returned the
	// ShortestAlts value.
	from graph.Node

	// nodes hold the nodes of the analysed
	// graph.
	nodes []graph.Node
	// indexOf contains a mapping between
	// the id-dense representation of the
	// graph and the potentially id-sparse
	// nodes held in nodes.
	indexOf map[int64]int

	// dist and next represent the shortest
	// paths between nodes.
	//
	// Indices into dist and next are
	// mapped through indexOf.
	//
	// dist contains the distances
	// from the from node for each
	// node in the graph.
	dist []float64
	// next contains the shortest-path
	// tree of the graph. The index is a
	// linear mapping of to-dense-id.
	next [][]int

	// hasNegativeCycle indicates
	// whether the ShortestAlts includes
	// a negative cycle. This should
	// be set by the function that
	// returned the ShortestAlts value.
	hasNegativeCycle bool

	// negCosts holds negative costs
	// between pairs of nodes to report
	// negative cycles.
	// negCosts must be initialised by
	// routines that can handle negative
	// edge weights.
	negCosts map[negEdge]float64
}

// newShortestAltsFrom returns a shortest path tree for all paths from u
// initialised with the given nodes. The nodes held by the returned
// Shortest may be lazily added.
func newShortestAltsFrom(u graph.Node, nodes []graph.Node) ShortestAlts {
	indexOf := make(map[int64]int, len(nodes))
	uid := u.ID()
	for i, n := range nodes {
		indexOf[n.ID()] = i
		if n.ID() == uid {
			u = n
		}
	}

	p := ShortestAlts{
		from: u,

		nodes:   nodes,
		indexOf: indexOf,

		dist: make([]float64, len(nodes)),
		next: make([][]int, len(nodes)),
	}
	for i := range nodes {
		p.dist[i] = math.Inf(1)
		p.next[i] = nil
	}
	p.dist[indexOf[uid]] = 0

	return p
}

// add adds a node to the ShortestAlts, initialising its stored index and returning, and
// setting the distance and position as unconnected. add will panic if the node is
// already present.
func (p *ShortestAlts) add(u graph.Node) int {
	uid := u.ID()
	if _, exists := p.indexOf[uid]; exists {
		panic("shortest: adding existing node")
	}
	idx := len(p.nodes)
	p.indexOf[uid] = idx
	p.nodes = append(p.nodes, u)
	p.dist = append(p.dist, math.Inf(1))
	p.next = append(p.next, nil)
	return idx
}

// set sets the weight of the path from the node in p.nodes indexed by mid to the node
// indexed by to.
func (p ShortestAlts) set(to int, weight float64, mid int) {
	p.dist[to] = weight
	p.next[to] = []int{mid}
	if weight < 0 {
		e := negEdge{from: mid, to: to}
		c, ok := p.negCosts[e]
		if !ok {
			p.negCosts[e] = weight
		} else if weight < c {
			// The only ways that we can have a new weight that is
			// lower than the previous weight is if either the edge
			// has already been traversed in a negative cycle, or
			// the edge is reachable from a negative cycle.
			// Either way the reported path is returned with a
			// negative infinite path weight.
			p.negCosts[e] = math.Inf(-1)
		}
	}
}

// addPath adds a new path from the node in p.nodes indexed by mid to the node indexed
// by to. The weight of the path is expected to be the same as already existing paths
// between these nodes, but no check is made for this.
func (p ShortestAlts) addPath(to, mid int) {
	// These are likely to be rare, so just loop over collisions.
	for _, v := range p.next[to] {
		if mid == v {
			return
		}
	}
	p.next[to] = append(p.next[to], mid)
}

// From returns the starting node of the paths held by the ShortestAlts.
func (p ShortestAlts) From() graph.Node { return p.from }

// WeightTo returns the weight of the minimum path to v. If the path to v includes
// a negative cycle, the returned weight will not reflect the true path weight.
func (p ShortestAlts) WeightTo(vid int64) float64 {
	to, toOK := p.indexOf[vid]
	if !toOK {
		return math.Inf(1)
	}
	return p.dist[to]
}

// To returns a shortest path to v and the weight of the path. If more than
// one shortest path exists between u and v, a randomly chosen path will be
// returned and unique is returned false. If a cycle with zero weight exists
// in the path, it will not be included, but unique will be returned false.
// If the path to v includes a negative cycle, one pass through the cycle will
// be included in path, but any path leading into the negative cycle will be
// lost, and weight will be returned as -Inf.
func (p ShortestAlts) To(vid int64) (path []graph.Node, weight float64, unique bool) {
	to, toOK := p.indexOf[vid]
	if !toOK || math.IsInf(p.dist[to], 1) {
		return nil, math.Inf(1), false
	}
	from := p.indexOf[p.from.ID()]
	unique = true
	path = []graph.Node{p.nodes[to]}
	if p.hasNegativeCycle {
		weight = math.Inf(1)
		seen := make(set.Ints)
		seen.Add(from)
		for to != from {
			c := p.next[to]
			var next int
			if len(c) != 1 {
				unique = false
				next = c[rand.Intn(len(c))]
			} else {
				next = c[0]
			}
			if math.IsInf(p.negCosts[negEdge{from: next, to: to}], -1) {
				weight = math.Inf(-1)
				unique = false
			}
			if seen.Has(to) {
				break
			}
			seen.Add(to)
			path = append(path, p.nodes[next])
			to = next
		}
		weight = math.Min(weight, p.dist[p.indexOf[vid]])
	} else {
		seen := make([]int, len(p.nodes))
		for i := range seen {
			seen[i] = -1
		}
		seen[to] = 0

		var next int
		for from != to {
			c := p.next[to]
			if len(c) != 1 {
				unique = false
				next = c[rand.Intn(len(c))]
			} else {
				next = c[0]
			}
			if seen[next] >= 0 {
				path = path[:seen[next]]
			}
			seen[next] = len(path)
			path = append(path, p.nodes[next])
			to = next
		}
		weight = p.dist[p.indexOf[vid]]
	}

	ordered.Reverse(path)
	return path, weight, unique
}

// AllTo returns all shortest paths to v and the weight of the paths. Paths
// containing zero-weight cycles are not returned. If a negative cycle exists between
// u and v, paths is returned nil and weight is returned as -Inf.
func (p ShortestAlts) AllTo(vid int64) (paths [][]graph.Node, weight float64) {
	from := p.indexOf[p.from.ID()]
	to, toOK := p.indexOf[vid]
	if !toOK || len(p.next[to]) == 0 {
		if p.from.ID() == vid {
			return [][]graph.Node{{p.nodes[from]}}, 0
		}
		return nil, math.Inf(1)
	}

	_, weight, unique := p.To(vid)
	if math.IsInf(weight, -1) && !unique {
		return nil, math.Inf(-1)
	}

	seen := make([]bool, len(p.nodes))
	p.allTo(from, to, seen, []graph.Node{p.nodes[to]}, func(path []graph.Node) {
		paths = append(paths, append([]graph.Node(nil), path...))
	})
	weight = p.dist[to]

	return paths, weight
}

// AllToFunc calls fn on all shortest paths to v. Paths containing zero-weight
// cycles are not considered. If a negative cycle exists between u and v, no
// path is considered. The fn closure must not retain the path parameter.
func (p ShortestAlts) AllToFunc(vid int64, fn func(path []graph.Node)) {
	from := p.indexOf[p.from.ID()]
	to, toOK := p.indexOf[vid]
	if !toOK || len(p.next[to]) == 0 {
		if p.from.ID() == vid {
			fn([]graph.Node{p.nodes[from]})
		}
		return
	}

	_, weight, unique := p.To(vid)
	if math.IsInf(weight, -1) && !unique {
		return
	}

	seen := make([]bool, len(p.nodes))
	p.allTo(from, to, seen, []graph.Node{p.nodes[to]}, fn)
}

// allTo recursively constructs a slice of paths extending from the node
// indexed into p.nodes by from to the node indexed by to. len(seen) must match
// the number of nodes held by the receiver. The path parameter is the current
// working path and the results passed to fn.
func (p ShortestAlts) allTo(from, to int, seen []bool, path []graph.Node, fn func(path []graph.Node)) {
	seen[to] = true
	if from == to {
		if path == nil {
			return
		}
		ordered.Reverse(path)
		fn(path)
		ordered.Reverse(path)
		return
	}
	first := true
	var seenWork []bool
	for _, to := range p.next[to] {
		if seen[to] {
			continue
		}
		if first {
			p := make([]graph.Node, len(path), len(path)+1)
			copy(p, path)
			path = p
			seenWork = make([]bool, len(seen))
			first = false
		}
		copy(seenWork, seen)
		p.allTo(from, to, seenWork, append(path, p.nodes[to]), fn)
	}
}

// negEdge is a key into the negative costs map used by Shortest and ShortestAlts.
type negEdge struct{ from, to int }

// AllShortest is a shortest-path tree created by the DijkstraAllPaths, FloydWarshall
// or JohnsonAllPaths all-pairs shortest paths functions.
type AllShortest struct {
	// nodes hold the nodes of the analysed
	// graph.
	nodes []graph.Node
	// indexOf contains a mapping between
	// the id-dense representation of the
	// graph and the potentially id-sparse
	// nodes held in nodes.
	indexOf map[int64]int

	// dist, next and forward represent
	// the shortest paths between nodes.
	//
	// Indices into dist and next are
	// mapped through indexOf.
	//
	// dist contains the pairwise
	// distances between nodes.
	//
	// Internally, edges on negative
	// cycles are given a special NaN
	// weight, NaN(0xdefaced).
	// This is returned to the user as
	// -Inf. This approach allows -Inf
	// weight edges on simple paths to be
	// distinguished from -Inf weight
	// paths that contain negative cycles.
	// The distinction is visible to the
	// user through whether then path
	// returned with a -Inf weight is
	// nil or contains a set of nodes.
	dist *mat.Dense
	// next contains the shortest-path
	// tree of the graph. The first index
	// is a linear mapping of from-dense-id
	// and to-dense-id, to-major with a
	// stride equal to len(nodes); the
	// slice indexed to is the list of
	// intermediates leading from the 'from'
	// node to the 'to' node represented
	// by dense id.
	// The interpretation of next is
	// dependent on the state of forward.
	next [][]int
	// forward indicates the direction of
	// path reconstruction. Forward
	// reconstruction is used for Floyd-
	// Warshall and reverse is used for
	// Dijkstra.
	forward bool
}

var (
	// defaced is NaN(0xdefaced) used as a marker for -Inf weight edges
	// within paths containing negative cycles. Routines marking these
	// edges should use this value.
	defaced = scalar.NaNWith(0xdefaced)
	// defacedBits is the bit pattern we look for in AllShortest to
	// identify the edges.
	defacedBits = math.Float64bits(defaced)
)

// newAllShortest returns an all-pairs shortest path forest for paths with the
// given nodes. The forward flag indicates whether the path reconstruction is
// performed in the forward (Floyd-Warshall) or reverse (Dijkstra/Johnson's) order.
func newAllShortest(nodes []graph.Node, forward bool) AllShortest {
	if len(nodes) == 0 {
		return AllShortest{}
	}
	indexOf := make(map[int64]int, len(nodes))
	for i, n := range nodes {
		indexOf[n.ID()] = i
	}
	dist := make([]float64, len(nodes)*len(nodes))
	for i := range dist {
		dist[i] = math.Inf(1)
	}
	return AllShortest{
		nodes:   nodes,
		indexOf: indexOf,

		dist:    mat.NewDense(len(nodes), len(nodes), dist),
		next:    make([][]int, len(nodes)*len(nodes)),
		forward: forward,
	}
}

// at returns a slice of node indexes into p.nodes for nodes that are mid points
// between nodes indexed by from and to.
func (p AllShortest) at(from, to int) (mid []int) {
	return p.next[from+to*len(p.nodes)]
}

// set sets the weights of paths between node indexes into p.nodes for from and to
// passing through the nodes indexed by mid.
func (p AllShortest) set(from, to int, weight float64, mid ...int) {
	p.dist.Set(from, to, weight)
	p.next[from+to*len(p.nodes)] = append(p.next[from+to*len(p.nodes)][:0], mid...)
}

// add adds paths between node indexed in p.nodes by from and to passing through
// the nodes indexed by mid.
func (p AllShortest) add(from, to int, mid ...int) {
loop: // These are likely to be rare, so just loop over collisions.
	for _, k := range mid {
		for _, v := range p.next[from+to*len(p.nodes)] {
			if k == v {
				continue loop
			}
		}
		p.next[from+to*len(p.nodes)] = append(p.next[from+to*len(p.nodes)], k)
	}
}

// Weight returns the weight of the minimum path between u and v.
func (p AllShortest) Weight(uid, vid int64) float64 {
	from, fromOK := p.indexOf[uid]
	to, toOK := p.indexOf[vid]
	if !fromOK || !toOK {
		return math.Inf(1)
	}
	w := p.dist.At(from, to)
	if math.Float64bits(w) == defacedBits {
		return math.Inf(-1)
	}
	return w
}

// Between returns a shortest path from u to v and the weight of the path. If more than
// one shortest path exists between u and v, a randomly chosen path will be returned and
// unique is returned false. If a cycle with zero weight exists in the path, it will not
// be included, but unique will be returned false. If a negative cycle exists on the path
// from u to v, path will be returned nil, weight will be -Inf and unique will be false.
func (p AllShortest) Between(uid, vid int64) (path []graph.Node, weight float64, unique bool) {
	from, fromOK := p.indexOf[uid]
	to, toOK := p.indexOf[vid]
	if !fromOK || !toOK || len(p.at(from, to)) == 0 {
		if uid == vid {
			if !fromOK {
				return []graph.Node{node(uid)}, 0, true
			}
			return []graph.Node{p.nodes[from]}, 0, true
		}
		return nil, math.Inf(1), false
	}

	weight = p.dist.At(from, to)
	if math.Float64bits(weight) == defacedBits {
		return nil, math.Inf(-1), false
	}

	seen := make([]int, len(p.nodes))
	for i := range seen {
		seen[i] = -1
	}
	var n graph.Node
	if p.forward {
		n = p.nodes[from]
		seen[from] = 0
	} else {
		n = p.nodes[to]
		seen[to] = 0
	}

	path = []graph.Node{n}
	unique = true

	var next int
	for from != to {
		c := p.at(from, to)
		if len(c) != 1 {
			unique = false
			next = c[rand.Intn(len(c))]
		} else {
			next = c[0]
		}
		if seen[next] >= 0 {
			path = path[:seen[next]]
		}
		seen[next] = len(path)
		path = append(path, p.nodes[next])
		if p.forward {
			from = next
		} else {
			to = next
		}
	}
	if !p.forward {
		ordered.Reverse(path)
	}

	return path, weight, unique
}

// AllBetween returns all shortest paths from u to v and the weight of the paths. Paths
// containing zero-weight cycles are not returned. If a negative cycle exists between
// u and v, paths is returned nil and weight is returned as -Inf.
func (p AllShortest) AllBetween(uid, vid int64) (paths [][]graph.Node, weight float64) {
	from, fromOK := p.indexOf[uid]
	to, toOK := p.indexOf[vid]
	if !fromOK || !toOK || len(p.at(from, to)) == 0 {
		if uid == vid {
			if !fromOK {
				return [][]graph.Node{{node(uid)}}, 0
			}
			return [][]graph.Node{{p.nodes[from]}}, 0
		}
		return nil, math.Inf(1)
	}

	weight = p.dist.At(from, to)
	if math.Float64bits(weight) == defacedBits {
		return nil, math.Inf(-1)
	}

	var n graph.Node
	if p.forward {
		n = p.nodes[from]
	} else {
		n = p.nodes[to]
	}
	seen := make([]bool, len(p.nodes))
	p.allBetween(from, to, seen, []graph.Node{n}, func(path []graph.Node) {
		paths = append(paths, append([]graph.Node(nil), path...))
	})

	return paths, weight
}

// AllBetweenFunc calls fn on all shortest paths from u to v. Paths containing
// zero-weight cycles are not considered. If a negative cycle exists between u
// and v, no path is considered. The fn closure must not retain the path
// parameter.
func (p AllShortest) AllBetweenFunc(uid, vid int64, fn func(path []graph.Node)) {
	from, fromOK := p.indexOf[uid]
	to, toOK := p.indexOf[vid]
	if !fromOK || !toOK || len(p.at(from, to)) == 0 {
		if uid == vid {
			if !fromOK {
				fn([]graph.Node{node(uid)})
				return
			}
			fn([]graph.Node{p.nodes[from]})
			return
		}
		return
	}

	if math.Float64bits(p.dist.At(from, to)) == defacedBits {
		return
	}

	var n graph.Node
	if p.forward {
		n = p.nodes[from]
	} else {
		n = p.nodes[to]
	}
	seen := make([]bool, len(p.nodes))
	p.allBetween(from, to, seen, []graph.Node{n}, fn)
}

// allBetween recursively constructs a set of paths extending from the node
// indexed into p.nodes by from to the node indexed by to. len(seen) must match
// the number of nodes held by the receiver. The path parameter is the current
// working path and the results passed to fn.
func (p AllShortest) allBetween(from, to int, seen []bool, path []graph.Node, fn func([]graph.Node)) {
	if p.forward {
		seen[from] = true
	} else {
		seen[to] = true
	}
	if from == to {
		if path == nil {
			return
		}
		if !p.forward {
			ordered.Reverse(path)
		}
		fn(path)
		if !p.forward {
			ordered.Reverse(path)
		}
		return
	}
	first := true
	var seenWork []bool
	for _, n := range p.at(from, to) {
		if seen[n] {
			continue
		}
		if first {
			p := make([]graph.Node, len(path), len(path)+1)
			copy(p, path)
			path = p
			seenWork = make([]bool, len(seen))
			first = false
		}
		if p.forward {
			from = n
		} else {
			to = n
		}
		copy(seenWork, seen)
		p.allBetween(from, to, seenWork, append(path, p.nodes[n]), fn)
	}
}

type node int64

func (n node) ID() int64 { return int64(n) }
