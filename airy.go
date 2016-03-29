// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathext

import "github.com/gonum/mathext/internal/amos"

// Airy returns the Airy function at z. The Airy function here, Ai(y), is one
// of the solutions to
//  y'' - y*z = 0.
// See http://mathworld.wolfram.com/AiryFunctions.html for more detailed information.
func Airy(z complex128) complex128 {
	id := 0
	kode := 1
	air, aii, _ := amos.Zairy(real(z), imag(z), id, kode)
	return complex(air, aii)
}

// AiryDeriv computes the derivative of the Airy function at z. The Airy function
// here, Ai(y), is one of the solutions to
//  y'' - y*z = 0.
// See http://mathworld.wolfram.com/AiryFunctions.html for more detailed information.
func AiryDeriv(z complex128) complex128 {
	id := 1
	kode := 1
	air, aii, _ := amos.Zairy(real(z), imag(z), id, kode)
	return complex(air, aii)
}
