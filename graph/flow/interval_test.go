// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flow

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"sort"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

// Ensure that Interval implements the flow.Graph interface.
var _ Graph = (*Interval)(nil)

func TestIntervals(t *testing.T) {
	golden := []struct {
		path string
		want [][]string
	}{
		{
			path: "testdata/allen.dot",
			want: [][]string{
				[]string{"1"},
				[]string{"2"},
				[]string{"3", "4", "5", "6"},
				[]string{"7", "8"},
			},
		},
		{
			path: "testdata/cifuentes.dot",
			want: [][]string{
				[]string{"B1", "B2", "B3", "B4", "B5"},
				[]string{"B6", "B7", "B8", "B9", "B10", "B11", "B12"},
				[]string{"B13", "B14", "B15"},
			},
		},
	}
	for _, gold := range golden {
		// Parse input.
		in, err := parseGraph(gold.path)
		if err != nil {
			t.Errorf("%q; unable to parse file; %v", gold.path, err)
			continue
		}
		// Locate intervals.
		is := Intervals(in)
		if len(gold.want) != len(is) {
			t.Errorf("%q: number of intervals mismatch; expected %d, got %d", gold.path, len(gold.want), len(is))
			continue
		}
		for i, want := range gold.want {
			var got []string
			for nodes := is[i].Nodes(); nodes.Next(); {
				n := nodes.Node()
				nn, ok := n.(dot.Node)
				if !ok {
					panic(fmt.Errorf("invalid node type; expected dot.Node, got %T", n))
				}
				got = append(got, nn.DOTID())
			}
			sort.Strings(got)
			sort.Strings(want)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("%q; output mismatch; expected `%s`, got `%s`", gold.path, want, got)
				continue
			}
		}
	}
}

// --- [ graph used for testing ] ---

// dotNode is the interface that groups the graph.Node and dot.Node interfaces.
type dotNode interface {
	graph.Node
	dot.Node
}

// parseGraph parses the control flow graph at the specified path. The node with
// the entry attribute is assigned as entry node of the control flow graph.
func parseGraph(dotPath string) (Graph, error) {
	data, err := ioutil.ReadFile(dotPath)
	if err != nil {
		return nil, err
	}
	dst := &CFG{
		directedBuilder: simple.NewDirectedGraph(),
	}
	if err := dot.Unmarshal(data, dst); err != nil {
		return nil, err
	}
	for nodes := dst.Nodes(); nodes.Next(); {
		if n, ok := nodes.Node().(*Node); ok {
			if _, ok := n.Attrs["entry"]; ok {
				return NewGraph(dst, n.ID()), nil
			}
		}
	}
	return nil, fmt.Errorf("unable to locate entry node of control flow graph %q", dotPath)
}

// directedBuilder is the interface that groups the graph.Directed and
// graph.Builder interfaces.
type directedBuilder interface {
	graph.Directed
	graph.Builder
}

// CFG is a control flow graph used for testing.
type CFG struct {
	directedBuilder
	entry graph.Node
}

// Entry returns the header node of the interval.
func (g *CFG) Entry() graph.Node {
	return g.entry
}

// NewNode returns a new Node with a unique arbitrary ID.
func (g *CFG) NewNode() graph.Node {
	return &Node{
		Node:  g.directedBuilder.NewNode(),
		Attrs: make(Attrs),
	}
}

// Node is a node of a control flow graph.
type Node struct {
	graph.Node
	// DOT ID of node.
	dotID string
	// DOT attributes of node.
	Attrs
}

// SetDOTID implements the dot.DOTIDSetter interface for Node.
func (n *Node) SetDOTID(id string) {
	n.dotID = id
}

// DOTID implements the dot.Node interface for Node.
func (n *Node) DOTID() string {
	return n.dotID
}

// Attrs is a set of key-value pair attributes used by graph.Node or graph.Edge.
type Attrs map[string]string

// Attributes implements encoding.Attributer for Attrs.
func (a Attrs) Attributes() []encoding.Attribute {
	attrs := make([]encoding.Attribute, 0, len(a))
	for key, val := range a {
		attr := encoding.Attribute{Key: key, Value: val}
		attrs = append(attrs, attr)
	}
	// Sort by key.
	less := func(i, j int) bool {
		if attrs[i].Key < attrs[j].Key {
			return true
		}
		return attrs[i].Value < attrs[j].Value
	}
	sort.Slice(attrs, less)
	return attrs
}

// SetAttribute implements encoding.AttributeSetter for Attrs.
func (a Attrs) SetAttribute(attr encoding.Attribute) error {
	a[attr.Key] = attr.Value
	return nil
}
