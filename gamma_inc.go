// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathext

import "github.com/gonum/mathext/internal/cephes"

// GammaInc computes the incomplete Gamma integral.
//
//                   1    x   -t  a-1
// GammaInc(a,x) = -----  ∫  e   t    dt
//                  Γ(a)  0
//
// In this implementation both arguments must be positive.
// The integral is evaluated by either a power series or
// continued fraction expansion, depending on the relative
// values of a and x.
// See the following for more information:
// - http://mathworld.wolfram.com/IncompleteGammaFunction.html
// - https://en.wikipedia.org/wiki/Incomplete_gamma_function
func GammaInc(a, x float64) float64 {
	return cephes.Igam(a, x)
}

// GammaIncC computes the complemented incomplete Gamma integral.
//
// GammaIncC(a,x) = 1 - GammaInc(a,x)
//
//                   =   1    ∞   -t  a-1
//                     -----  ∫  e   t    dt
//                      Γ(a)  x
//
// In this implementation both arguments must be positive.
// The integral is evaluated by either a power series or
// continued fraction expansion, depending on the relative
// values of a and x.
func GammaIncC(a, x float64) float64 {
	return cephes.IgamC(a, x)
}

// GammaIncCInv returns x such that:
//
//  GammaIncC(a, x) = y
//
// for positive a and p between 0 and 1.
func GammaIncCInv(a, y float64) float64 {
	return cephes.IgamI(a, y)
}
