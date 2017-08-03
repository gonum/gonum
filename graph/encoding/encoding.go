// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package encoding provides a common graph encoding API.
package encoding // import "gonum.org/v1/gonum/graph/encoding"

import "gonum.org/v1/gonum/graph"

// Builder is a graph that can have user-defined nodes and edges added.
type Builder interface {
	graph.Graph
	graph.Builder
	// NewEdge adds a new edge from the source to the destination node to the
	// graph, or returns the existing edge if already present.
	NewEdge(from, to graph.Node) graph.Edge
}

// UnmarshalerAttr is implemented by types that can unmarshal a graph
// attribute description of themselves.
type UnmarshalerAttr interface {
	// UnmarshalAttr decodes a single attribute.
	UnmarshalAttr(attr Attribute) error
}

// Attributer defines graph.Node or graph.Edge values that can
// specify graph attributes.
type Attributer interface {
	Attributes() []Attribute
}

// Attribute is an encoded key value attribute pair use in graph encoding.
type Attribute struct {
	Key, Value string
}
