// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math"

	"gonum.org/v1/gonum/blas/cblas128"
)

// Zlarfg generates a complex elementary reflector H of order n such that
//
//	H * (alpha) = (beta)
//	    (    x)   (   0)
//	H^H * H = I
//
// where alpha and beta are complex scalars with beta real, and x is an (n-1)-element
// complex vector. H is represented in the form
//
//	H = I - tau * (1; v) * (1; v)^H
//
// where tau is a complex scalar and v is a complex (n-1)-element vector. If the
// elements of x are all zero and imag(alpha) = 0, then tau = 0 and H is taken to
// be the identity matrix.
//
// Otherwise 1 <= real(tau) <= 2 and abs(tau-1) <= 1.
//
// On entry, x contains x. On exit, x is overwritten with v. alpha is returned
// as beta (a real number).
//
// Zlarfg is an internal routine. It is exported for testing purposes.
func (impl Implementation) Zlarfg(n int, alpha complex128, x []complex128, incX int) (beta float64, tau complex128) {
	switch {
	case n < 0:
		panic(nLT0)
	case incX <= 0:
		panic(badIncX)
	}

	if n <= 0 {
		return real(alpha), 0
	}

	// For n >= 2, x has length >= 1+(n-2)*incX; for n == 1, x is unused.
	if n > 1 && len(x) < 1+(n-2)*incX {
		panic(shortX)
	}

	bi := cblas128.Implementation()

	var xnorm float64
	if n > 1 {
		xnorm = bi.Dznrm2(n-1, x, incX)
	}
	alphr := real(alpha)
	alphi := imag(alpha)

	if xnorm == 0 && alphi == 0 {
		// H = I.
		return alphr, 0
	}

	// beta = -sign(alphr) * sqrt(alphr^2 + alphi^2 + xnorm^2).
	beta = -math.Copysign(dlapy3(alphr, alphi, xnorm), alphr)
	safmin := dlamchS / dlamchE
	rsafmn := 1 / safmin
	knt := 0
	if math.Abs(beta) < safmin {
		// xnorm, beta may be inaccurate; scale x and recompute them.
		for {
			knt++
			if n > 1 {
				bi.Zdscal(n-1, rsafmn, x, incX)
			}
			beta *= rsafmn
			alphi *= rsafmn
			alphr *= rsafmn
			if math.Abs(beta) >= safmin {
				break
			}
		}
		// New beta is at most 1, at least safmin.
		if n > 1 {
			xnorm = bi.Dznrm2(n-1, x, incX)
		}
		alpha = complex(alphr, alphi)
		beta = -math.Copysign(dlapy3(alphr, alphi, xnorm), alphr)
	}

	// tau = (beta - alpha) / beta = complex((beta-alphr)/beta, -alphi/beta).
	tau = complex((beta-alphr)/beta, -alphi/beta)

	// Scale x: v = x / (alpha - beta).
	inv := 1 / complex(alphr-beta, alphi)
	if n > 1 {
		bi.Zscal(n-1, inv, x, incX)
	}

	// If beta was scaled, unscale it.
	for j := 0; j < knt; j++ {
		beta *= safmin
	}
	return beta, tau
}

// dlapy3 returns sqrt(x^2 + y^2 + z^2), avoiding unnecessary overflow/underflow.
func dlapy3(x, y, z float64) float64 {
	xabs := math.Abs(x)
	yabs := math.Abs(y)
	zabs := math.Abs(z)
	w := math.Max(xabs, math.Max(yabs, zabs))
	if w == 0 || math.IsInf(w, 1) {
		return xabs + yabs + zabs
	}
	xw := xabs / w
	yw := yabs / w
	zw := zabs / w
	return w * math.Sqrt(xw*xw+yw*yw+zw*zw)
}
