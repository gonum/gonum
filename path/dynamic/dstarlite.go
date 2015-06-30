// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dynamic

import (
	"container/heap"
	"fmt"
	"math"

	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
	"github.com/gonum/graph/path"
)

// DStarLite implements the D* Lite dynamic re-planning path search algorithm.
//
//  doi:10.1109/tro.2004.838026 and ISBN:0-262-51129-0 pp476-483
//
type DStarLite struct {
	s, t *dStarLiteNode
	last *dStarLiteNode

	model       WorldModel
	queue       dStarLiteQueue
	keyModifier float64

	weight    graph.WeightFunc
	heuristic path.Heuristic
}

// WorldModel is a mutable weighted directed graph that returns nodes identified
// by id number.
type WorldModel interface {
	graph.MutableDirected
	graph.Weighter
	Node(id int) graph.Node
}

// NewDStarLite returns a new DStarLite planner for the path from s to t in g using the
// heuristic h. The world model, m, is used to store shortest path information during path
// planning. The world model must be an empty graph when NewDStarLite is called.
//
// If h is nil, the DStarLite will use the g.HeuristicCost method if g implements
// path.HeuristicCoster, falling back to path.NullHeuristic otherwise. If the graph does not
// implement graph.Weighter, graph.UniformCost is used. NewDStarLite will panic if g has
// a negative edge weight.
func NewDStarLite(s, t graph.Node, g graph.Graph, h path.Heuristic, m WorldModel) *DStarLite {
	/*
	   procedure Initialize()
	   {02”} U = ∅;
	   {03”} k_m = 0;
	   {04”} for all s ∈ S rhs(s) = g(s) = ∞;
	   {05”} rhs(s_goal) = 0;
	   {06”} U.Insert(s_goal, [h(s_start, s_goal); 0]);
	*/

	d := &DStarLite{
		s: &dStarLiteNode{Node: s, rhs: math.Inf(1), g: math.Inf(1), key: badKey},
		t: &dStarLiteNode{Node: t, g: math.Inf(1), key: badKey}, // badKey is overwritten below.

		model: m,

		heuristic: h,
	}

	/*
		procedure Main()
		{29”} s_last = s_start;
		{30”} Initialize();
	*/
	d.last = d.s

	if g, ok := g.(graph.Weighter); ok {
		d.weight = g.Weight
	} else {
		d.weight = graph.UniformCost
	}
	if d.heuristic == nil {
		if g, ok := g.(path.HeuristicCoster); ok {
			d.heuristic = g.HeuristicCost
		} else {
			d.heuristic = path.NullHeuristic
		}
	}

	d.queue.indexOf = make(map[int]int)
	d.queue.insert(d.t, key{d.heuristic(s, t), 0})

	for _, n := range g.Nodes() {
		switch n.ID() {
		case d.s.ID():
			d.model.AddNode(d.s)
		case d.t.ID():
			d.model.AddNode(d.t)
		default:
			d.model.AddNode(&dStarLiteNode{Node: n, rhs: math.Inf(1), g: math.Inf(1), key: badKey})
		}
	}
	for _, u := range d.model.Nodes() {
		for _, v := range g.From(u) {
			w := d.weight(g.Edge(u, v))
			if w < 0 {
				panic("D* Lite: negative edge weight")
			}
			d.model.SetEdge(concrete.Edge{F: u, T: d.model.Node(v.ID())}, w)
		}
	}

	/*
		procedure Main()
		{31”} ComputeShortestPath();
	*/
	d.findShortestPath()

	return d
}

// keyFor is the CalculateKey procedure in the D* Lite papers.
func (d *DStarLite) keyFor(s *dStarLiteNode) key {
	/*
	   procedure CalculateKey(s)
	   {01”} return [min(g(s), rhs(s)) + h(s_start, s) + k_m; min(g(s), rhs(s))];
	*/
	k := key{1: math.Min(s.g, s.rhs)}
	k[0] = k[1] + d.heuristic(d.s.Node, s.Node) + d.keyModifier
	return k
}

// update is the UpdateVertex procedure in the D* Lite papers.
func (d *DStarLite) update(u *dStarLiteNode) {
	/*
	   procedure UpdateVertex(u)
	   {07”} if (g(u) != rhs(u) AND u ∈ U) U.Update(u,CalculateKey(u));
	   {08”} else if (g(u) != rhs(u) AND u /∈ U) U.Insert(u,CalculateKey(u));
	   {09”} else if (g(u) = rhs(u) AND u ∈ U) U.Remove(u);
	*/
	uid := u.ID()
	inQueue := d.queue.has(uid)
	switch {
	case inQueue && u.g != u.rhs:
		d.queue.update(uid, d.keyFor(u))
	case !inQueue && u.g != u.rhs:
		d.queue.insert(u, d.keyFor(u))
	case inQueue && u.g == u.rhs:
		d.queue.remove(uid)
	}
}

// findShortestPath is the ComputeShortestPath procedure in the D* Lite papers.
func (d *DStarLite) findShortestPath() {
	/*
	   procedure ComputeShortestPath()
	   {10”} while (U.TopKey() < CalculateKey(s_start) OR rhs(s_start) > g(s_start))
	   {11”} u = U.Top();
	   {12”} k_old = U.TopKey();
	   {13”} k_new = CalculateKey(u);
	   {14”} if(k_old < k_new)
	   {15”}   U.Update(u, k_new);
	   {16”} else if (g(u) > rhs(u))
	   {17”}   g(u) = rhs(u);
	   {18”}   U.Remove(u);
	   {19”}   for all s ∈ Pred(u)
	   {20”}     if (s != s_goal) rhs(s) = min(rhs(s), c(s, u) + g(u));
	   {21”}     UpdateVertex(s);
	   {22”} else
	   {23”}   g_old = g(u);
	   {24”}   g(u) = ∞;
	   {25”}   for all s ∈ Pred(u) ∪ {u}
	   {26”}     if (rhs(s) = c(s, u) + g_old)
	   {27”}       if (s != s_goal) rhs(s) = min s'∈Succ(s)(c(s, s') + g(s'));
	   {28”}     UpdateVertex(s);
	*/
	for d.queue.Len() != 0 { // We use d.queue.Len since d.queue does not return an infinite key when empty.
		u := d.queue.top()
		if !u.key.less(d.keyFor(d.s)) && d.s.rhs <= d.s.g {
			break
		}
		switch kNew := d.keyFor(u); {
		case u.key.less(kNew):
			d.queue.update(u.ID(), kNew)
		case u.g > u.rhs:
			u.g = u.rhs
			d.queue.remove(u.ID())
			for _, _s := range d.model.To(u) {
				s := _s.(*dStarLiteNode)
				if s.ID() != d.t.ID() {
					s.rhs = math.Min(s.rhs, d.model.Weight(d.model.Edge(s, u))+u.g)
				}
				d.update(s)
			}
		default:
			gOld := u.g
			u.g = math.Inf(1)
			for _, _s := range append(d.model.To(u), u) {
				s := _s.(*dStarLiteNode)
				// This is necessary since we are potentially
				// asking for the weight of u to u, but we
				// never return a self-edge.
				var w float64
				if s.ID() != u.ID() {
					w = d.model.Weight(d.model.Edge(s, u))
				}
				if s.rhs == w+gOld {
					if s.ID() != d.t.ID() {
						s.rhs = math.Inf(1)
						for _, t := range d.model.From(s) {
							s.rhs = math.Min(s.rhs, d.model.Weight(d.model.Edge(s, t))+t.(*dStarLiteNode).g)
						}
					}
				}
				d.update(s)
			}
		}
	}
}

// Step performs one movement step along the best path towards the goal.
// It returns false if no further progression toward the goal can be
// achieved, either because the goal has been reached or because there
// is no path.
func (d *DStarLite) Step() bool {
	/*
	   procedure Main()
	   {32”} while (s_start != s_goal)
	   {33”} // if (rhs(s_start) = ∞) then there is no known path
	   {34”}   s_start = argmin s'∈Succ(s_start)(c(s_start, s') + g(s'));
	*/
	if d.s.ID() == d.t.ID() {
		return false
	}
	if math.IsInf(d.s.rhs, 1) {
		return false
	}

	// We use rhs comparison to break ties
	// between coequally weighted nodes.
	rhs := math.Inf(1)
	min := math.Inf(1)

	var next *dStarLiteNode
	for _, _s := range d.model.From(d.s) {
		s := _s.(*dStarLiteNode)
		w := d.model.Weight(d.model.Edge(d.s, s)) + s.g
		if w < min || (w == min && s.rhs < rhs) {
			next = s
			min = w
			rhs = s.rhs
		}
	}
	d.s = next

	/*
	   procedure Main()
	   {35”}   Move to s_start;
	*/
	return true
}

// MoveTo moves to n in the world graph.
func (d *DStarLite) MoveTo(n graph.Node) {
	d.last = d.s
	d.s = d.model.Node(n.ID()).(*dStarLiteNode)
	d.keyModifier += d.heuristic(d.last, d.s)
}

// UpdateWorld updates or adds edges in the world graph. UpdateWorld will
// panic if changes include a negative edge weight.
func (d *DStarLite) UpdateWorld(changes []graph.Edge) {
	/*
	   procedure Main()
	   {36”}   Scan graph for changed edge costs;
	   {37”}   if any edge costs changed
	   {38”}     k_m = k_m + h(s_last, s_start);
	   {39”}     s_last = s_start;
	   {40”}     for all directed edges (u, v) with changed edge costs
	   {41”}       c_old = c(u, v);
	   {42”}       Update the edge cost c(u, v);
	   {43”}       if (c_old > c(u, v))
	   {44”}         if (u != s_goal) rhs(u) = min(rhs(u), c(u, v) + g(v));
	   {45”}       else if (rhs(u) = c_old + g(v))
	   {46”}         if (u != s_goal) rhs(u) = min s'∈Succ(u)(c(u, s') + g(s'));
	   {47”}       UpdateVertex(u);
	   {48”}     ComputeShortestPath()
	*/
	if len(changes) == 0 {
		return
	}
	d.keyModifier += d.heuristic(d.last, d.s)
	d.last = d.s
	for _, e := range changes {
		c := d.weight(e)
		if c < 0 {
			panic("D* Lite: negative edge weight")
		}
		cOld := d.model.Weight(e)
		u := d.worldNodeFor(e.From())
		v := d.worldNodeFor(e.To())
		d.model.SetEdge(concrete.Edge{F: u, T: v}, c)
		if cOld > c {
			if u.ID() != d.t.ID() {
				u.rhs = math.Min(u.rhs, c+v.g)
			}
		} else if u.rhs == cOld+v.g {
			if u.ID() != d.t.ID() {
				u.rhs = math.Inf(1)
				for _, t := range d.model.From(u) {
					u.rhs = math.Min(u.rhs, d.model.Weight(d.model.Edge(u, t))+t.(*dStarLiteNode).g)
				}
			}
		}
		d.update(u)
	}
	d.findShortestPath()
}

func (d *DStarLite) worldNodeFor(n graph.Node) *dStarLiteNode {
	switch w := d.model.Node(n.ID()).(type) {
	case *dStarLiteNode:
		return w
	case graph.Node:
		panic(fmt.Sprintf("D* Lite: illegal world model node type: %T", w))
	default:
		return &dStarLiteNode{Node: n, rhs: math.Inf(1), g: math.Inf(1), key: badKey}
	}
}

// Here returns the current location.
func (d *DStarLite) Here() graph.Node {
	return d.s.Node
}

// Path returns the path from the current location to the goal and the
// weight of the path.
func (d *DStarLite) Path() (p []graph.Node, weight float64) {
	u := d.s
	p = []graph.Node{u.Node}
	for u.ID() != d.t.ID() {
		if math.IsInf(u.rhs, 1) {
			return nil, math.Inf(1)
		}

		// We use stored rhs comparison to break
		// ties between calculated rhs-coequal nodes.
		rhsMin := math.Inf(1)
		min := math.Inf(1)
		var (
			next *dStarLiteNode
			cost float64
		)
		for _, _v := range d.model.From(u) {
			v := _v.(*dStarLiteNode)
			w := d.model.Weight(d.model.Edge(u, v))
			if rhs := w + v.g; rhs < min || (rhs == min && v.rhs < rhsMin) {
				next = v
				min = rhs
				rhsMin = v.rhs
				cost = w
			}
		}
		if next == nil {
			return nil, math.NaN()
		}
		u = next
		weight += cost
		p = append(p, u.Node)
	}
	return p, weight
}

/*
The pseudocode uses the following functions to manage the priority
queue:

      * U.Top() returns a vertex with the smallest priority of all
        vertices in priority queue U.
      * U.TopKey() returns the smallest priority of all vertices in
        priority queue U. (If is empty, then U.TopKey() returns [∞;∞].)
      * U.Pop() deletes the vertex with the smallest priority in
        priority queue U and returns the vertex.
      * U.Insert(s, k) inserts vertex s into priority queue with
        priority k.
      * U.Update(s, k) changes the priority of vertex s in priority
        queue U to k. (It does nothing if the current priority of vertex
        s already equals k.)
      * Finally, U.Remove(s) removes vertex s from priority queue U.
*/

// key is a D* Lite priority queue key.
type key [2]float64

var badKey = key{math.NaN(), math.NaN()}

// less returns whether k is less than other. From ISBN:0-262-51129-0 pp476-483:
//
//  k ≤ k' iff k₁ < k'₁ OR (k₁ == k'₁ AND k₂ ≤ k'₂)
//
func (k key) less(other key) bool {
	if k != k || other != other {
		panic("D* Lite: poisoned key")
	}
	return k[0] < other[0] || (k[0] == other[0] && k[1] < other[1])
}

// dStarLiteNode adds D* Lite accounting to a graph.Node.
type dStarLiteNode struct {
	graph.Node
	key key
	rhs float64
	g   float64
}

// dStarLiteQueue is a D* Lite priority queue.
type dStarLiteQueue struct {
	indexOf map[int]int
	nodes   []*dStarLiteNode
}

func (q *dStarLiteQueue) Less(i, j int) bool {
	return q.nodes[i].key.less(q.nodes[j].key)
}

func (q *dStarLiteQueue) Swap(i, j int) {
	q.indexOf[q.nodes[i].ID()] = j
	q.indexOf[q.nodes[j].ID()] = i
	q.nodes[i], q.nodes[j] = q.nodes[j], q.nodes[i]
}

func (q *dStarLiteQueue) Len() int {
	return len(q.nodes)
}

func (q *dStarLiteQueue) Push(x interface{}) {
	n := x.(*dStarLiteNode)
	q.indexOf[n.ID()] = len(q.nodes)
	q.nodes = append(q.nodes, n)
}

func (q *dStarLiteQueue) Pop() interface{} {
	n := q.nodes[len(q.nodes)-1]
	q.nodes = q.nodes[:len(q.nodes)-1]
	delete(q.indexOf, n.ID())
	return n
}

// has returns whether the node identified by id is in the queue.
func (q *dStarLiteQueue) has(id int) bool {
	_, ok := q.indexOf[id]
	return ok
}

// top returns the top node in the queue. Note that instead of
// returning a key [∞;∞] when q is empty, the caller checks for
// an empty queue by calling q.Len.
func (q *dStarLiteQueue) top() *dStarLiteNode {
	return q.nodes[0]
}

// insert puts the node u into the queue with the key k.
func (q *dStarLiteQueue) insert(u *dStarLiteNode, k key) {
	u.key = k
	heap.Push(q, u)
}

// update updates the node in the queue identified by id with the key k.
func (q *dStarLiteQueue) update(id int, k key) {
	i, ok := q.indexOf[id]
	if !ok {
		return
	}
	q.nodes[i].key = k
	heap.Fix(q, i)
}

// remove removes the node identified by id from the queue.
func (q *dStarLiteQueue) remove(id int) {
	i, ok := q.indexOf[id]
	if !ok {
		return
	}
	q.nodes[i].key = badKey
	heap.Remove(q, i)
}
