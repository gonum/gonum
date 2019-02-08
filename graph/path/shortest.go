// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"math"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/internal/ordered"
	"gonum.org/v1/gonum/graph/internal/set"
	"gonum.org/v1/gonum/mat"
)

// Shortest is a shortest-path tree created by the BellmanFordFrom or DijkstraFrom
// single-source shortest path functions.
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
}

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

func (p Shortest) set(to int, weight float64, mid int) {
	p.dist[to] = weight
	p.next[to] = mid
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
// in path and weight will be returned as -Inf.
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
			if seen.Has(to) {
				weight = math.Inf(-1)
				break
			}
			seen.Add(to)
			path = append(path, p.nodes[p.next[to]])
			to = p.next[to]
		}
	} else {
		n := len(p.nodes)
		for to != from {
			path = append(path, p.nodes[p.next[to]])
			to = p.next[to]
			if n < 0 {
				panic("path: unexpected negative cycle")
			}
			n--
		}
	}
	ordered.Reverse(path)
	return path, math.Min(weight, p.dist[p.indexOf[vid]])
}

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
	defaced = floats.NaNWith(0xdefaced)
	// defacedBits is the bit pattern we look for in AllShortest to
	// identify the edges.
	defacedBits = math.Float64bits(defaced)
)

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

func (p AllShortest) at(from, to int) (mid []int) {
	return p.next[from+to*len(p.nodes)]
}

func (p AllShortest) set(from, to int, weight float64, mid ...int) {
	p.dist.Set(from, to, weight)
	p.next[from+to*len(p.nodes)] = append(p.next[from+to*len(p.nodes)][:0], mid...)
}

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
	paths = p.allBetween(from, to, seen, []graph.Node{n}, nil)

	return paths, weight
}

func (p AllShortest) allBetween(from, to int, seen []bool, path []graph.Node, paths [][]graph.Node) [][]graph.Node {
	if p.forward {
		seen[from] = true
	} else {
		seen[to] = true
	}
	if from == to {
		if path == nil {
			return paths
		}
		if !p.forward {
			ordered.Reverse(path)
		}
		return append(paths, path)
	}
	first := true
	for _, n := range p.at(from, to) {
		if seen[n] {
			continue
		}
		if first {
			path = append([]graph.Node(nil), path...)
			first = false
		}
		if p.forward {
			from = n
		} else {
			to = n
		}
		paths = p.allBetween(from, to, append([]bool(nil), seen...), append(path, p.nodes[n]), paths)
	}
	return paths
}
