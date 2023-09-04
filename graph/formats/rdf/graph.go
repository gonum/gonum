// Copyright ©2022 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdf

import (
	"fmt"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/multi"
	"gonum.org/v1/gonum/graph/set/uid"
)

// Graph implements an RDF graph satisfying the graph.Graph and graph.Multigraph
// interfaces.
type Graph struct {
	nodes map[int64]graph.Node
	from  map[int64]map[int64]map[int64]graph.Line
	to    map[int64]map[int64]map[int64]graph.Line
	pred  map[int64]map[*Statement]bool

	termIDs map[string]int64
	ids     *uid.Set
}

// NewGraph returns a new empty Graph.
func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[int64]graph.Node),
		from:  make(map[int64]map[int64]map[int64]graph.Line),
		to:    make(map[int64]map[int64]map[int64]graph.Line),
		pred:  make(map[int64]map[*Statement]bool),

		termIDs: make(map[string]int64),
		ids:     uid.NewSet(),
	}
}

// addNode adds n to the graph. It panics if the added node ID matches an
// existing node ID.
func (g *Graph) addNode(n graph.Node) {
	if _, exists := g.nodes[n.ID()]; exists {
		panic(fmt.Sprintf("rdf: node ID collision: %d", n.ID()))
	}
	g.nodes[n.ID()] = n
	g.ids.Use(n.ID())
}

// AddStatement adds s to the graph. It panics if Term UIDs in the statement
// are not consistent with existing terms in the graph. Statements must not
// be altered while being held by the graph. If the UID fields of the terms
// in s are zero, they will be set to values consistent with the rest of the
// graph on return, mutating the parameter, otherwise the UIDs must match terms
// that already exist in the graph. The statement must be a valid RDF statement
// otherwise AddStatement will panic.
func (g *Graph) AddStatement(s *Statement) {
	_, _, kind, err := s.Predicate.Parts()
	if err != nil {
		panic(fmt.Errorf("rdf: error extracting predicate: %w", err))
	}
	if kind != IRI {
		panic(fmt.Errorf("rdf: predicate is not an IRI: %s", s.Predicate.Value))
	}

	_, _, kind, err = s.Subject.Parts()
	if err != nil {
		panic(fmt.Errorf("rdf: error extracting subject: %w", err))
	}
	switch kind {
	case IRI, Blank:
	default:
		panic(fmt.Errorf("rdf: subject is not an IRI or blank node: %s", s.Subject.Value))
	}

	_, _, kind, err = s.Object.Parts()
	if err != nil {
		panic(fmt.Errorf("rdf: error extracting object: %w", err))
	}
	if kind == Invalid {
		panic(fmt.Errorf("rdf: object is not a valid term: %s", s.Object.Value))
	}

	statements, ok := g.pred[s.Predicate.UID]
	if !ok {
		statements = make(map[*Statement]bool)
		g.pred[s.Predicate.UID] = statements
	}
	statements[s] = true
	g.addTerm(&s.Subject)
	g.addTerm(&s.Predicate)
	g.addTerm(&s.Object)
	g.setLine(s)
}

// addTerm adds t to the graph. It panics if the added node ID matches an existing node ID.
func (g *Graph) addTerm(t *Term) {
	if t.UID == 0 {
		id, ok := g.termIDs[t.Value]
		if ok {
			t.UID = id
			return
		}
		id = g.ids.NewID()
		g.ids.Use(id)
		t.UID = id
		g.termIDs[t.Value] = id
		return
	}

	id, ok := g.termIDs[t.Value]
	if !ok {
		g.termIDs[t.Value] = t.UID
	} else if id != t.UID {
		panic(fmt.Sprintf("rdf: term ID collision: term:%s new ID:%d old ID:%d", t.Value, t.UID, id))
	}
}

// AllStatements returns an iterator of the statements that make up the graph.
func (g *Graph) AllStatements() *Statements {
	return &Statements{eit: g.Edges()}
}

// Edge returns the edge from u to v if such an edge exists and nil otherwise.
// The node v must be directly reachable from u as defined by the From method.
// The returned graph.Edge is a multi.Edge if an edge exists.
func (g *Graph) Edge(uid, vid int64) graph.Edge {
	l := g.Lines(uid, vid)
	if l == graph.Empty {
		return nil
	}
	return multi.Edge{F: g.Node(uid), T: g.Node(vid), Lines: l}
}

// Edges returns all the edges in the graph. Each edge in the returned slice
// is a multi.Edge.
func (g *Graph) Edges() graph.Edges {
	if len(g.nodes) == 0 {
		return graph.Empty
	}
	var edges []graph.Edge
	for _, u := range g.nodes {
		for _, e := range g.from[u.ID()] {
			var lines []graph.Line
			for _, l := range e {
				lines = append(lines, l)
			}
			if len(lines) != 0 {
				edges = append(edges, multi.Edge{
					F:     g.Node(u.ID()),
					T:     g.Node(lines[0].To().ID()),
					Lines: iterator.NewOrderedLines(lines),
				})
			}
		}
	}
	if len(edges) == 0 {
		return graph.Empty
	}
	return iterator.NewOrderedEdges(edges)
}

// From returns all nodes in g that can be reached directly from n.
//
// The returned graph.Nodes is only valid until the next mutation of
// the receiver.
func (g *Graph) From(id int64) graph.Nodes {
	if len(g.from[id]) == 0 {
		return graph.Empty
	}
	return iterator.NewNodesByLines(g.nodes, g.from[id])
}

// FromSubject returns all nodes in g that can be reached directly from an
// RDF subject term.
//
// The returned graph.Nodes is only valid until the next mutation of
// the receiver.
func (g *Graph) FromSubject(t Term) graph.Nodes {
	return g.From(t.UID)
}

// HasEdgeBetween returns whether an edge exists between nodes x and y without
// considering direction.
func (g *Graph) HasEdgeBetween(xid, yid int64) bool {
	if _, ok := g.from[xid][yid]; ok {
		return true
	}
	_, ok := g.from[yid][xid]
	return ok
}

// HasEdgeFromTo returns whether an edge exists in the graph from u to v.
func (g *Graph) HasEdgeFromTo(uid, vid int64) bool {
	_, ok := g.from[uid][vid]
	return ok
}

// Lines returns the lines from u to v if such any such lines exists and nil otherwise.
// The node v must be directly reachable from u as defined by the From method.
func (g *Graph) Lines(uid, vid int64) graph.Lines {
	edge := g.from[uid][vid]
	if len(edge) == 0 {
		return graph.Empty
	}
	var lines []graph.Line
	for _, l := range edge {
		lines = append(lines, l)
	}
	return iterator.NewOrderedLines(lines)
}

// newLine returns a new Line from the source to the destination node.
// The returned Line will have a graph-unique ID.
// The Line's ID does not become valid in g until the Line is added to g.
func (g *Graph) newLine(from, to graph.Node) graph.Line {
	return multi.Line{F: from, T: to, UID: g.ids.NewID()}
}

// newNode returns a new unique Node to be added to g. The Node's ID does
// not become valid in g until the Node is added to g.
func (g *Graph) newNode() graph.Node {
	if len(g.nodes) == 0 {
		return multi.Node(0)
	}
	if int64(len(g.nodes)) == uid.Max {
		panic("rdf: cannot allocate node: no slot")
	}
	return multi.Node(g.ids.NewID())
}

// Node returns the node with the given ID if it exists in the graph,
// and nil otherwise.
func (g *Graph) Node(id int64) graph.Node {
	return g.nodes[id]
}

// TermFor returns the Term for the given text. The text must be
// an exact match for the Term's Value field.
func (g *Graph) TermFor(text string) (term Term, ok bool) {
	id, ok := g.termIDs[text]
	if !ok {
		return
	}
	n, ok := g.nodes[id]
	if !ok {
		var s map[*Statement]bool
		s, ok = g.pred[id]
		if !ok {
			return
		}
		for k := range s {
			return k.Predicate, true
		}
	}
	return n.(Term), true
}

// Nodes returns all the nodes in the graph.
//
// The returned graph.Nodes is only valid until the next mutation of
// the receiver.
func (g *Graph) Nodes() graph.Nodes {
	if len(g.nodes) == 0 {
		return graph.Empty
	}
	return iterator.NewNodes(g.nodes)
}

// Predicates returns a slice of all the predicates used in the graph.
func (g *Graph) Predicates() []Term {
	p := make([]Term, len(g.pred))
	i := 0
	for _, statements := range g.pred {
		for s := range statements {
			p[i] = s.Predicate
			i++
			break
		}
	}
	return p
}

// removeLine removes the line with the given end point and line IDs from
// the graph, leaving the terminal nodes. If the line does not exist it is
// a no-op.
func (g *Graph) removeLine(fid, tid, id int64) {
	if _, ok := g.nodes[fid]; !ok {
		return
	}
	if _, ok := g.nodes[tid]; !ok {
		return
	}

	delete(g.from[fid][tid], id)
	if len(g.from[fid][tid]) == 0 {
		delete(g.from[fid], tid)
	}
	delete(g.to[tid][fid], id)
	if len(g.to[tid][fid]) == 0 {
		delete(g.to[tid], fid)
	}

	g.ids.Release(id)
}

// removeNode removes the node with the given ID from the graph, as well as
// any edges attached to it. If the node is not in the graph it is a no-op.
func (g *Graph) removeNode(id int64) {
	if _, ok := g.nodes[id]; !ok {
		return
	}
	delete(g.nodes, id)

	for from := range g.from[id] {
		delete(g.to[from], id)
	}
	delete(g.from, id)

	for to := range g.to[id] {
		delete(g.from[to], id)
	}
	delete(g.to, id)

	g.ids.Release(id)
}

// RemoveStatement removes s from the graph, leaving the terminal nodes if they
// are part of another statement. If the statement does not exist in g it is a no-op.
func (g *Graph) RemoveStatement(s *Statement) {
	if !g.pred[s.Predicate.UID][s] {
		return
	}

	// Remove the connection.
	g.removeLine(s.Subject.UID, s.Object.UID, s.Predicate.UID)
	statements := g.pred[s.Predicate.UID]
	delete(statements, s)
	if len(statements) == 0 {
		delete(g.pred, s.Predicate.UID)
		if len(g.from[s.Predicate.UID]) == 0 {
			g.ids.Release(s.Predicate.UID)
			delete(g.termIDs, s.Predicate.Value)
		}
	}

	// Remove any orphan terms.
	if g.From(s.Subject.UID).Len() == 0 && g.To(s.Subject.UID).Len() == 0 {
		g.removeNode(s.Subject.UID)
		delete(g.termIDs, s.Subject.Value)
	}
	if g.From(s.Object.UID).Len() == 0 && g.To(s.Object.UID).Len() == 0 {
		g.removeNode(s.Object.UID)
		delete(g.termIDs, s.Object.Value)
	}
}

// RemoveTerm removes t and any statements referencing t from the graph. If
// the term is a predicate, all statements with the predicate are removed. If
// the term does not exist it is a no-op.
func (g *Graph) RemoveTerm(t Term) {
	// Remove any predicates.
	if statements, ok := g.pred[t.UID]; ok {
		for s := range statements {
			g.RemoveStatement(s)
		}
	}

	// Quick return.
	_, nok := g.nodes[t.UID]
	_, fok := g.from[t.UID]
	_, tok := g.to[t.UID]
	if !nok && !fok && !tok {
		return
	}

	// Remove any statements that impinge on the term.
	to := g.From(t.UID)
	for to.Next() {
		lines := g.Lines(t.UID, to.Node().ID())
		for lines.Next() {
			g.RemoveStatement(lines.Line().(*Statement))
		}
	}
	from := g.To(t.UID)
	if from.Next() {
		lines := g.Lines(from.Node().ID(), t.UID)
		for lines.Next() {
			g.RemoveStatement(lines.Line().(*Statement))
		}
	}

	// Remove the node.
	g.removeNode(t.UID)
	delete(g.termIDs, t.Value)
}

// setLine adds l, a line from one node to another. If the nodes do not exist,
// they are added, and are set to the nodes of the line otherwise.
func (g *Graph) setLine(l graph.Line) {
	var (
		from = l.From()
		fid  = from.ID()
		to   = l.To()
		tid  = to.ID()
		lid  = l.ID()
	)

	if _, ok := g.nodes[fid]; !ok {
		g.addNode(from)
	} else {
		g.nodes[fid] = from
	}
	if _, ok := g.nodes[tid]; !ok {
		g.addNode(to)
	} else {
		g.nodes[tid] = to
	}

	switch {
	case g.from[fid] == nil:
		g.from[fid] = map[int64]map[int64]graph.Line{tid: {lid: l}}
	case g.from[fid][tid] == nil:
		g.from[fid][tid] = map[int64]graph.Line{lid: l}
	default:
		g.from[fid][tid][lid] = l
	}
	switch {
	case g.to[tid] == nil:
		g.to[tid] = map[int64]map[int64]graph.Line{fid: {lid: l}}
	case g.to[tid][fid] == nil:
		g.to[tid][fid] = map[int64]graph.Line{lid: l}
	default:
		g.to[tid][fid][lid] = l
	}

	g.ids.Use(lid)
}

// Statements returns an iterator of the statements that connect the subject
// term node u to the object term node v.
func (g *Graph) Statements(uid, vid int64) *Statements {
	return &Statements{lit: g.Lines(uid, vid)}
}

// To returns all nodes in g that can reach directly to n.
//
// The returned graph.Nodes is only valid until the next mutation of
// the receiver.
func (g *Graph) To(id int64) graph.Nodes {
	if len(g.to[id]) == 0 {
		return graph.Empty
	}
	return iterator.NewNodesByLines(g.nodes, g.to[id])
}

// ToObject returns all nodes in g that can reach directly to an RDF object
// term.
//
// The returned graph.Nodes is only valid until the next mutation of
// the receiver.
func (g *Graph) ToObject(t Term) graph.Nodes {
	return g.To(t.UID)
}

// Statements is an RDF statement iterator.
type Statements struct {
	eit graph.Edges
	lit graph.Lines
}

// Next returns whether the iterator holds any additional statements.
func (s *Statements) Next() bool {
	if s.lit != nil && s.lit.Next() {
		return true
	}
	if s.eit == nil || !s.eit.Next() {
		return false
	}
	s.lit = s.eit.Edge().(multi.Edge).Lines
	return s.lit.Next()
}

// Statement returns the current statement.
func (s *Statements) Statement() *Statement {
	return s.lit.Line().(*Statement)
}

// ConnectedByAny is a helper function to for simplifying graph traversal
// conditions.
func ConnectedByAny(e graph.Edge, with func(*Statement) bool) bool {
	switch e := e.(type) {
	case *Statement:
		return with(e)
	case graph.Lines:
		it := e
		for it.Next() {
			s, ok := it.Line().(*Statement)
			if !ok {
				continue
			}
			ok = with(s)
			if ok {
				return true
			}
		}
	}
	return false
}
