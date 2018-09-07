// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package stat

import (
	"testing"
)

func TestCounts(t *testing.T) {
	t.Log("Testing Counts")
	d := []int{0, 1, 1, 2, 2, 2, 3, 3, 3, 3}
	c := Counts(d)

	if len(c) != 4 {
		t.Errorf("Counts should return a slice of length 4 and not %d", len(c))
	}

	for i := 0; i < len(c); i++ {
		if (i + 1) != c[i] {
			t.Errorf("c[%d] should be 1 but it is %d", (i + 1), c[i])
		}
	}
}
