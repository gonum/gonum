// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fftpack

// swap returns c with the real and imaginary parts swapped.
func swap(c complex128) complex128 {
	return complex(imag(c), real(c))
}

// scale scales the complex number c by f.
func scale(f float64, c complex128) complex128 {
	return complex(f*real(c), f*imag(c))
}
