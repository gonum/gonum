// Copyright ©2023 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import "math"

// Dlapy3 computes sqrt(x**2+y**2+z**2) while minimizing the possibility of
// overflow.
//
// Dlapy3 is an internal routine. It is exported for testing purposes.
func (Implementation) Dlapy3(x, y, z float64) float64 {
	return math.Hypot(x, math.Hypot(y, z))
}
