// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathext

import "github.com/zeroviscosity/mathext/internal/cephes"

// IncGamma computes the incomplete Gamma integral
// In this implementation both arguments must be positive.
// The integral is evaluated by either a power series or
// continued fraction expansion, depending on the relative
// values of a and x.
func IncGamma(a, x float64) float64 {
	return cephes.Igam(a, x)
}

// IncGammaComp computes the complemented incomplete Gamma integral
// In this implementation both arguments must be positive.
// The integral is evaluated by either a power series or
// continued fraction expansion, depending on the relative
// values of a and x.
func IncGammaComp(a, x float64) float64 {
	return cephes.IgamC(a, x)
}
