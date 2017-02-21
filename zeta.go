// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathext

import "github.com/gonum/mathext/internal/cephes"

// Zeta computes the Riemann zeta function of two arguments.
//
//             ∞      -x
// Zeta(x,q) = ∑ (k+q)
//            k=0
//
// where x > 1 and q is not a negative integer or zero.
// See the following for more information:
// - http://mathworld.wolfram.com/HurwitzZetaFunction.html
// - https://en.wikipedia.org/wiki/Multiple_zeta_function#Two_parameters_case
func Zeta(x, q float64) float64 {
	return cephes.Zeta(x, q)
}
