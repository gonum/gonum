// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"testing"

	"gonum.org/v1/gonum/graph/path/internal/testgraphs"
	"gonum.org/v1/gonum/graph"
)


func TestYenKSP(t *testing.T) {
	for _, test := range testgraphs.YenShortestPathTests {
		g := test.Graph()
		for _, e := range test.Edges {
			g.SetWeightedEdge(e)
		}

		paths := YenKShortestPath(g.(graph.Graph), test.K, test.Query.From(), test.Query.To())
		for i := 0; i < len(test.WantPaths); i++ {
			expected := test.WantPaths[i]
			path := paths[i]
			if (len(expected) == len(path)) {
				for n := 0; n < len(path); n++ {
					if (expected[n] != path[n].ID()) {
						t.Errorf("ERROR: path #%d expected: %d, got: %d", i+1, expected[n], path[n].ID())
					}
				}
			}
		}

	}
}
