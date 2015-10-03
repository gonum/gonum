// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build appengine

package mat64

import "reflect"

var sizeOfFloat64 = reflect.TypeOf(float64(0)).Size()

// offset returns the number of float64 values b[0] is after a[0].
func offset(a, b []float64) int {
	// This block must be atomic with respect to GC moves.
	// At this stage this is true, because the GC does not
	// move.
	a0 := reflect.ValueOf(a).Index(0).UnsafeAddr()
	b0 := reflect.ValueOf(b).Index(0).UnsafeAddr()

	if a0 == b0 {
		return 0
	}
	if a0 < b0 {
		return int((b0 - a0) / sizeOfFloat64)
	}
	return -int((a0 - b0) / sizeOfFloat64)
}
