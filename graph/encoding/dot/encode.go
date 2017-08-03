// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dot

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/internal/ordered"
)

// Node is a DOT graph node.
type Node interface {
	// DOTID returns a DOT node ID.
	//
	// An ID is one of the following:
	//
	//  - a string of alphabetic ([a-zA-Z\x80-\xff]) characters, underscores ('_').
	//    digits ([0-9]), not beginning with a digit.
	//  - a numeral [-]?(.[0-9]+ | [0-9]+(.[0-9]*)?).
	//  - a double-quoted string ("...") possibly containing escaped quotes (\").
	//  - an HTML string (<...>).
	DOTID() string
}

// Attributers are graph.Graph values that specify top-level DOT
// attributes.
type Attributers interface {
	DOTAttributers() (graph, node, edge encoding.Attributer)
}

// Porter defines the behavior of graph.Edge values that can specify
// connection ports for their end points. The returned port corresponds
// to the the DOT node port to be used by the edge, compass corresponds
// to DOT compass point to which the edge will be aimed.
type Porter interface {
	FromPort() (port, compass string)
	ToPort() (port, compass string)
}

// Structurer represents a graph.Graph that can define subgraphs.
type Structurer interface {
	Structure() []Graph
}

// Graph wraps named graph.Graph values.
type Graph interface {
	graph.Graph
	DOTID() string
}

// Subgrapher wraps graph.Node values that represent subgraphs.
type Subgrapher interface {
	Subgraph() graph.Graph
}

// Marshal returns the DOT encoding for the graph g, applying the prefix
// and indent to the encoding. Name is used to specify the graph name. If
// name is empty and g implements Graph, the returned string from DOTID
// will be used. If strict is true the output bytes will be prefixed with
// the DOT "strict" keyword.
//
// Graph serialization will work for a graph.Graph without modification,
// however, advanced GraphViz DOT features provided by Marshal depend on
// implementation of the Node, Attributer, Porter, Attributers, Structurer,
// Subgrapher and Graph interfaces.
func Marshal(g graph.Graph, name, prefix, indent string, strict bool) ([]byte, error) {
	var p printer
	p.indent = indent
	p.prefix = prefix
	p.visited = make(map[edge]bool)
	if strict {
		p.buf.WriteString("strict ")
	}
	err := p.print(g, name, false, false)
	if err != nil {
		return nil, err
	}
	return p.buf.Bytes(), nil
}

type printer struct {
	buf bytes.Buffer

	prefix string
	indent string
	depth  int

	visited map[edge]bool

	err error
}

type edge struct {
	inGraph  string
	from, to int64
}

func (p *printer) print(g graph.Graph, name string, needsIndent, isSubgraph bool) error {
	nodes := g.Nodes()
	sort.Sort(ordered.ByID(nodes))

	p.buf.WriteString(p.prefix)
	if needsIndent {
		for i := 0; i < p.depth; i++ {
			p.buf.WriteString(p.indent)
		}
	}
	_, isDirected := g.(graph.Directed)
	if isSubgraph {
		p.buf.WriteString("sub")
	} else if isDirected {
		p.buf.WriteString("di")
	}
	p.buf.WriteString("graph")

	if name == "" {
		if g, ok := g.(Graph); ok {
			name = g.DOTID()
		}
	}
	if name != "" {
		p.buf.WriteByte(' ')
		p.buf.WriteString(name)
	}

	p.openBlock(" {")
	if a, ok := g.(Attributers); ok {
		p.writeAttributeComplex(a)
	}
	if s, ok := g.(Structurer); ok {
		for _, g := range s.Structure() {
			_, subIsDirected := g.(graph.Directed)
			if subIsDirected != isDirected {
				return errors.New("dot: mismatched graph type")
			}
			p.buf.WriteByte('\n')
			p.print(g, g.DOTID(), true, true)
		}
	}

	havePrintedNodeHeader := false
	for _, n := range nodes {
		if s, ok := n.(Subgrapher); ok {
			// If the node is not linked to any other node
			// the graph needs to be written now.
			if len(g.From(n)) == 0 {
				g := s.Subgraph()
				_, subIsDirected := g.(graph.Directed)
				if subIsDirected != isDirected {
					return errors.New("dot: mismatched graph type")
				}
				if !havePrintedNodeHeader {
					p.newline()
					p.buf.WriteString("// Node definitions.")
					havePrintedNodeHeader = true
				}
				p.newline()
				p.print(g, graphID(g, n), false, true)
			}
			continue
		}
		if !havePrintedNodeHeader {
			p.newline()
			p.buf.WriteString("// Node definitions.")
			havePrintedNodeHeader = true
		}
		p.newline()
		p.writeNode(n)
		if a, ok := n.(encoding.Attributer); ok {
			p.writeAttributeList(a)
		}
		p.buf.WriteByte(';')
	}

	havePrintedEdgeHeader := false
	for _, n := range nodes {
		to := g.From(n)
		sort.Sort(ordered.ByID(to))
		for _, t := range to {
			if isDirected {
				if p.visited[edge{inGraph: name, from: n.ID(), to: t.ID()}] {
					continue
				}
				p.visited[edge{inGraph: name, from: n.ID(), to: t.ID()}] = true
			} else {
				if p.visited[edge{inGraph: name, from: n.ID(), to: t.ID()}] {
					continue
				}
				p.visited[edge{inGraph: name, from: n.ID(), to: t.ID()}] = true
				p.visited[edge{inGraph: name, from: t.ID(), to: n.ID()}] = true
			}

			if !havePrintedEdgeHeader {
				p.buf.WriteByte('\n')
				p.buf.WriteString(strings.TrimRight(p.prefix, " \t\n")) // Trim whitespace suffix.
				p.newline()
				p.buf.WriteString("// Edge definitions.")
				havePrintedEdgeHeader = true
			}
			p.newline()

			if s, ok := n.(Subgrapher); ok {
				g := s.Subgraph()
				_, subIsDirected := g.(graph.Directed)
				if subIsDirected != isDirected {
					return errors.New("dot: mismatched graph type")
				}
				p.print(g, graphID(g, n), false, true)
			} else {
				p.writeNode(n)
			}
			e, edgeIsPorter := g.Edge(n, t).(Porter)
			if edgeIsPorter {
				p.writePorts(e.FromPort())
			}

			if isDirected {
				p.buf.WriteString(" -> ")
			} else {
				p.buf.WriteString(" -- ")
			}

			if s, ok := t.(Subgrapher); ok {
				g := s.Subgraph()
				_, subIsDirected := g.(graph.Directed)
				if subIsDirected != isDirected {
					return errors.New("dot: mismatched graph type")
				}
				p.print(g, graphID(g, t), false, true)
			} else {
				p.writeNode(t)
			}
			if edgeIsPorter {
				p.writePorts(e.ToPort())
			}

			if a, ok := g.Edge(n, t).(encoding.Attributer); ok {
				p.writeAttributeList(a)
			}

			p.buf.WriteByte(';')
		}
	}
	p.closeBlock("}")

	return nil
}

func (p *printer) writeNode(n graph.Node) {
	p.buf.WriteString(nodeID(n))
}

func (p *printer) writePorts(port, cp string) {
	if port != "" {
		p.buf.WriteByte(':')
		p.buf.WriteString(port)
	}
	if cp != "" {
		p.buf.WriteByte(':')
		p.buf.WriteString(cp)
	}
}

func nodeID(n graph.Node) string {
	switch n := n.(type) {
	case Node:
		return n.DOTID()
	default:
		return fmt.Sprint(n.ID())
	}
}

func graphID(g graph.Graph, n graph.Node) string {
	switch g := g.(type) {
	case Node:
		return g.DOTID()
	default:
		return nodeID(n)
	}
}

func (p *printer) writeAttributeList(a encoding.Attributer) {
	attributes := a.Attributes()
	switch len(attributes) {
	case 0:
	case 1:
		p.buf.WriteString(" [")
		p.buf.WriteString(attributes[0].Key)
		p.buf.WriteByte('=')
		p.buf.WriteString(attributes[0].Value)
		p.buf.WriteString("]")
	default:
		p.openBlock(" [")
		for _, att := range attributes {
			p.newline()
			p.buf.WriteString(att.Key)
			p.buf.WriteByte('=')
			p.buf.WriteString(att.Value)
		}
		p.closeBlock("]")
	}
}

var attType = []string{"graph", "node", "edge"}

func (p *printer) writeAttributeComplex(ca Attributers) {
	g, n, e := ca.DOTAttributers()
	haveWrittenBlock := false
	for i, a := range []encoding.Attributer{g, n, e} {
		attributes := a.Attributes()
		if len(attributes) == 0 {
			continue
		}
		if haveWrittenBlock {
			p.buf.WriteByte(';')
		}
		p.newline()
		p.buf.WriteString(attType[i])
		p.openBlock(" [")
		for _, att := range attributes {
			p.newline()
			p.buf.WriteString(att.Key)
			p.buf.WriteByte('=')
			p.buf.WriteString(att.Value)
		}
		p.closeBlock("]")
		haveWrittenBlock = true
	}
	if haveWrittenBlock {
		p.buf.WriteString(";\n")
	}
}

func (p *printer) newline() {
	p.buf.WriteByte('\n')
	p.buf.WriteString(p.prefix)
	for i := 0; i < p.depth; i++ {
		p.buf.WriteString(p.indent)
	}
}

func (p *printer) openBlock(b string) {
	p.buf.WriteString(b)
	p.depth++
}

func (p *printer) closeBlock(b string) {
	p.depth--
	p.newline()
	p.buf.WriteString(b)
}
