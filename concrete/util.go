package concrete

import (
	"github.com/gonum/graph"
)

type nodeSorter []graph.Node

func (ns nodeSorter) Less(i, j int) bool {
	return ns[i].ID() < ns[j].ID()
}

func (ns nodeSorter) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}

func (ns nodeSorter) Len() int {
	return len(ns)
}
