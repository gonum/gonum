// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !appengine

package mat64

import "unsafe"

// offset returns the number of float64 values b[0] is after a[0].
func offset(a, b []float64) int {
	// This block must be atomic with respect to GC moves.
	// At this stage this is true, because the GC does not
	// move.
	a0 := uintptr(unsafe.Pointer(&a[0]))
	b0 := uintptr(unsafe.Pointer(&b[0]))

	if a0 == b0 {
		return 0
	}
	if a0 < b0 {
		return int((b0 - a0) / unsafe.Sizeof(float64(0)))
	}
	return -int((a0 - b0) / unsafe.Sizeof(float64(0)))
}
