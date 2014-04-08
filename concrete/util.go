// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
