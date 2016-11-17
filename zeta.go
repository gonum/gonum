// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathext

import "github.com/zeroviscosity/mathext/internal/cephes"

// Zeta computes the Riemann zeta function of two arguments
func Zeta(x, q float64) float64 {
	return cephes.Zeta(x, q)
}
