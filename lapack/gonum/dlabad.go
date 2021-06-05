// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import "math"

// DLABAD takes as input the values computed by DLAMCH for underflow and
// overflow, and returns the square root of each of these values if the
// log of LARGE is sufficiently large.  This subroutine is intended to
// identify machines with a large exponent range, such as the Crays, and
// redefine the underflow and overflow limits to be the square roots of
// the values computed by DLAMCH.  This subroutine is needed because
// DLAMCH does not compensate for poor arithmetic in the upper half of
// the exponent range, as is found on a Cray.
func (impl Implementation) Dlabad(small, large float64) (float64, float64) {
	if math.Log10(large) > 2000. {
		return math.Sqrt(small), math.Sqrt(large)
	}
	return small, large
}
