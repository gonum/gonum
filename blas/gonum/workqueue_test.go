// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import "testing"

func TestBlockWorkQueue(t *testing.T) {
	// exhaustive test for 4 items
	for m := 0; m < blockSize*2+1; m++ {
		for n := 0; n < blockSize*2+1; n++ {
			testBlockWorkQueue(t, m, n)
		}
	}
}

func testBlockWorkQueue(t *testing.T, m, n int) {
	var q blockWorkQueue
	q.Reset(m, n)
	for i := 0; i < m; i += blockSize {
		for j := 0; j < n; j += blockSize {
			ei, ej, ok := q.Next()
			if !ok {
				t.Fatalf("[%d:%d:%d] expected <%d, %d>", m, n, blockSize, i, j)
			}
			if i != ei || j != ej {
				t.Fatalf("[%d:%d:%d] expected <%d, %d> got <%d, %d>", m, n, blockSize, i, j, ei, ej)
			}
		}
	}
}
