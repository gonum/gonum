// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathext

import "github.com/gonum/mathext/internal/cephes"

// IncGamma computes the incomplete Gamma integral.
//
//                   1    x   -t  a-1
// IncGamma(a,x) = -----  ∫  e   t    dt
//                  Γ(a)  0
//
// In this implementation both arguments must be positive.
// The integral is evaluated by either a power series or
// continued fraction expansion, depending on the relative
// values of a and x.
// See the following for more information:
// - http://mathworld.wolfram.com/IncompleteGammaFunction.html
// - https://en.wikipedia.org/wiki/Incomplete_gamma_function
func IncGamma(a, x float64) float64 {
	return cephes.Igam(a, x)
}

// IncGammaComp computes the complemented incomplete Gamma integral.
//
// IncGammaComp(a,x) = 1 - IncGamma(a,x)
//
//                   =   1    ∞   -t  a-1
//                     -----  ∫  e   t    dt
//                      Γ(a)  x
//
// In this implementation both arguments must be positive.
// The integral is evaluated by either a power series or
// continued fraction expansion, depending on the relative
// values of a and x.
func IncGammaComp(a, x float64) float64 {
	return cephes.IgamC(a, x)
}
