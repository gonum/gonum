// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

// Between returns statistics relating two nodes in a graph.
// The path type considered depends on the specific
// implementation e.g. shortest path, earliest-arrival path
// etc.
type Between interface {
	// Weight returns the path weight between uid and vid.
	Weight(uid, vid int64) float64
}