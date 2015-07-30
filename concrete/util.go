// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package concrete

import "math"

// The math package only provides explicitly sized max
// values. This ensures we get the max for the actual
// type int.
const maxInt int = int(^uint(0) >> 1)

var inf = math.Inf(1)

func isSame(a, b float64) bool {
	return a == b || (math.IsNaN(a) && math.IsNaN(b))
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
