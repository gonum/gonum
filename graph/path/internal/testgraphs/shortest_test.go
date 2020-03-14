// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testgraphs

import "testing"

func TestVerifyShortestPathTests(t *testing.T) {
	t.Parallel()
	for _, test := range ShortestPathTests {
		if len(test.WantPaths) != 1 && test.HasUniquePath {
			t.Fatalf("%q: bad shortest path test: non-unique paths marked unique", test.Name)
		}
	}
}
