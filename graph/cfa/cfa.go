// Package cfa provides control flow analysis functions.
package cfa

import "gonum.org/v1/gonum/graph"

// Graph is a control flow graph, a directed graph rooted at an entry node.
type Graph interface {
	// Entry returns the entry node of the control flow graph.
	Entry() graph.Node
	graph.Directed
}
