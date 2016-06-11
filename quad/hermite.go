// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package quad

import (
	"math"

	"github.com/gonum/floats"
	"github.com/gonum/mathext/airy"
)

// Hermite generates sample locations and weights for performing quadrature with
// with a squared-exponential weight
//  int_-inf^inf e^(-x^2) f(x) dx .
type Hermite struct{}

func (h Hermite) FixedLocations(x, weight []float64, min, max float64) {
	// TODO(btracey): Implement the case where x > 20, x < 200 so that we don't
	// need to store all of that data.

	// References:
	// Algorithm:
	// G. H. Golub and J. A. Welsch, "Calculation of Gauss quadrature rules",
	// Math. Comp. 23:221-230, 1969.
	// A. Glaser, X. Liu and V. Rokhlin, "A fast algorithm for the
	// calculation of the roots of special functions", SIAM Journal
	// on Scientific Computing", 29(4):1420-1438:, 2007.
	// A. Townsend, T. Trogdon, and S.Olver, Fast computation of Gauss quadrature
	// nodes and weights on the whole real line, IMA J. Numer. Anal., 36: 337–358,
	// 2016. http://arxiv.org/abs/1410.5286
	//
	// Algorithm adapted from Chubfun http://www.chebfun.org/.

	if len(x) != len(weight) {
		panic("hermite: slice length mismatch")
	}
	if min >= max {
		panic("hermite: min >= max")
	}
	if !math.IsInf(min, -1) || !math.IsInf(max, 1) {
		panic("hermite: non-infinite bound")
	}
	h.locations(x, weight)
}

func (h Hermite) locations(x, weights []float64) {
	n := len(x)
	switch {
	case 0 < n && n <= 200:
		copy(x, xCacheHermite[n-1])
		copy(weights, wCacheHermite[n-1])
	case n > 200:
		xasy, weightsasy := h.locationsAsy(n)
		copy(x, xasy)
		copy(weights, weightsasy)
	}
}

// Algorithm adapted from Chebfun http://www.chebfun.org/. Specific code
// https://github.com/chebfun/chebfun/blob/development/hermpts.m.

// Original Copyright Notice:

/*
Copyright (c) 2015, The Chancellor, Masters and Scholars of the University
of Oxford, and the Chebfun Developers. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
    * Neither the name of the University of Oxford nor the names of its
      contributors may be used to endorse or promote products derived from
      this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

func (h Hermite) locationsAsy(n int) (x, w []float64) {
	// A. Townsend, T. Trogdon, and S.Olver, Fast computation of Gauss quadrature
	// nodes and weights the whole real line, IMA J. Numer. Anal.,
	// 36: 337–358, 2016. http://arxiv.org/abs/1410.5286
	xa, wa := h.locationsAsy0(n)
	if n%2 == 1 {
		for i := len(xa) - 1; i >= 0; i-- {
			x = append(x, -xa[i])
			w = append(w, wa[i])
		}
		for i := 1; i < len(xa); i++ {
			x = append(x, xa[i])
			w = append(w, wa[i])
		}
	} else {
		lxa := len(xa)
		x = make([]float64, 2*lxa)
		for i, v := range xa {
			x[lxa-1-i] = -v
			x[lxa+i] = v
		}
		lwa := len(wa)
		w = make([]float64, 2*lwa)
		for i, v := range wa {
			w[lwa-1-i] = v
			w[lwa+i] = v
		}
	}
	sumW := floats.Sum(w)
	c := math.SqrtPi / sumW
	floats.Scale(c, w)
	return x, w
}

func (h Hermite) locationsAsy0(n int) (x, w []float64) {
	eps := math.Nextafter(1, math.Inf(1)) - 1
	x0pts := make([]float64, n/2+n%2)
	// Compute Hermite nodes and weights using asymptotic formula.
	h.hermiteInitialGuesses(x0pts, n)
	theta0 := x0pts
	for i, x0 := range theta0 {
		t0 := x0 / math.Sqrt(2*float64(n)+1)
		theta0[i] = math.Acos(t0)
	}
	dts := make([]float64, len(theta0))
	var val, dval []float64
	for k := 0; k < 20; k++ {
		val, dval = h.hermpolyAsyAiry(n, theta0)
		for i, t0 := range theta0 {
			dt := -val[i] / (math.Sqrt2 * math.Sqrt(2*float64(n)+1) * dval[i] * math.Sin(t0))
			theta0[i] = t0 - dt
			dts[i] = dt
		}
		if floats.Norm(dts, math.Inf(1)) < math.Sqrt(eps)/10 {
			break
		}
	}

	x = make([]float64, len(theta0))
	w = make([]float64, len(theta0))
	for i, t := range theta0 {
		t0 := math.Cos(t)
		xi := math.Sqrt(2*float64(n)+1) * t0
		x[i] = xi
		ders := xi*val[i] + math.Sqrt2*dval[i]
		w[i] = math.Exp(-xi*xi) / (ders * ders)
	}
	return x, w
}

// hermpolyAsyAiry evaluates the Hermite polynomials using the Airy asymptotic
// formula in theta-space.
func (h Hermite) hermpolyAsyAiry(n int, theta []float64) (valVec, dvalVec []float64) {
	valVec = make([]float64, len(theta))
	dvalVec = make([]float64, len(theta))
	musq := 2*float64(n) + 1

	for i, t := range theta {
		cosT := math.Cos(t)
		sinT := math.Sin(t)
		sin2T := 2 * cosT * sinT
		eta := 0.5*t - 0.25*sin2T
		chi := -math.Pow(3*eta/2, 2.0/3)
		phi := math.Pow(-chi/(sinT*sinT), 1.0/4)
		cnst := 2 * math.SqrtPi * math.Pow(musq, 1.0/6) * phi
		airy0 := real(airy.Ai(complex(math.Pow(musq, 2.0/3)*chi, 0)))
		airy1 := real(airy.AiDeriv(complex(math.Pow(musq, 2.0/3)*chi, 0)))
		// Terms in 12.10.43:
		const (
			a0 = 1.0
			b0 = 1.0
			a1 = 15.0 / 144
			b1 = -7.0 / 5 * a1
			a2 = 5.0 * 7 * 9 * 11.0 / 2.0 / 144.0 / 144.0
			b2 = -13.0 / 11 * a2
			a3 = 7.0 * 9 * 11 * 13 * 15 * 17 / 6.0 / 144.0 / 144.0 / 144.0
			b3 = -19.0 / 17 * a3
		)

		// u polynomials in 12.10.9.
		u0 := 1.0
		u1 := (cosT*cosT*cosT - 6*cosT) / 24.0
		u2 := (-9*cosT*cosT*cosT*cosT + 249*cosT*cosT + 145) / 1152.0
		u3 := (-4042*math.Pow(cosT, 9) + 18189*math.Pow(cosT, 7) -
			28287*math.Pow(cosT, 5) - 151995*math.Pow(cosT, 3) -
			259290*cosT) / 414720.0

		// First term.
		A0 := 1.0
		val := A0 * airy0

		// Second term.
		B0 := -(a0*math.Pow(phi, 6)*u1 + a1*u0) / (chi * chi)
		val += B0 * airy1 / math.Pow(musq, 4.0/3)

		// Third term.
		A1 := (b0*math.Pow(phi, 12)*u2 + b1*math.Pow(phi, 6)*u1 + b2*u0) / (chi * chi * chi)
		val += A1 * airy0 / (musq * musq)

		// Fourth term.
		B1 := -(math.Pow(phi, 18)*u3 + a1*math.Pow(phi, 12)*u2 +
			a2*math.Pow(phi, 6)*u1 + a3*u0) / math.Pow(chi, 5)

		val += B1 * airy1 / math.Pow(musq, 4.0/3+2)
		val *= cnst

		// Derivative.
		eta = 0.5*t - 0.25*sin2T
		chi = -math.Pow(3*eta/2, 2.0/3)
		phi = math.Pow(-chi/(sinT*sinT), 1.0/4)
		cnst = math.Sqrt2 * math.SqrtPi * math.Pow(musq, 1.0/3) / phi

		// v polynomials in 12.10.10.
		v0 := 1.0
		v1 := (cosT*cosT*cosT + 6*cosT) / 24
		v2 := (15*cosT*cosT*cosT*cosT - 327*cosT*cosT - 143) / 1152
		v3 := (259290*cosT + 238425*cosT*cosT*cosT - 36387*math.Pow(cosT, 5) +
			18189*math.Pow(cosT, 7) - 4042*math.Pow(cosT, 9)) / 414720

		// First term.
		C0 := -(b0*math.Pow(phi, 6)*v1 + b1*v0) / chi
		dval := C0 * airy0 / math.Pow(musq, 2.0/3)

		// Second term.
		D0 := a0 * v0
		dval += D0 * airy1

		// Third term.
		C1 := -(math.Pow(phi, 18)*v3 + b1*math.Pow(phi, 12)*v2 +
			b2*math.Pow(phi, 6)*v1 + b3*v0) / math.Pow(chi, 4)
		dval += C1 * airy0 / math.Pow(musq, 2.0/3+2)

		// Fourth term.
		D1 := (a0*math.Pow(phi, 12)*v2 + a1*math.Pow(phi, 6)*v1 + a2*v0) / math.Pow(chi, 3)
		dval += D1 * airy1 / (musq * musq)
		dval *= cnst

		valVec[i] = val
		dvalVec[i] = dval
	}
	return valVec, dvalVec
}

// hermiteInitialGuesses returns a set of initial guesses for the hermite
// quadrature locations. The results are stored in-place into guessses. Guesses
// has length of ceil(n/2).
func (h Hermite) hermiteInitialGuesses(guesses []float64, n int) {
	// Initial guesses for Hermite zeros.
	if len(guesses) != n/2+n%2 {
		panic("hermite: bad guesses length")
	}

	// There are two different formulas for the initial guesses of the hermite
	// quadrature locations. The first uses the Gatteschi formula and is good
	// near x = sqrt(n+0.5)
	//  [1] L. Gatteschi, Asymptotics and bounds for the zeros of Laguerre
	//  polynomials: a survey, J. Comput. Appl. Math., 144 (2002), pp. 7-27.
	// The second is the Tricomi initial guesses, good near x = 0. This is
	// equation 2.1 in [1] and is originally from
	//  [2] F. G. Tricomi, Sugli zeri delle funzioni di cui si conosce una
	//  rappresentazione asintotica, Ann. Mat. Pura Appl. 26 (1947), pp. 283-300.

	// If the number of points is odd, there is a quadrature point at 1, which
	// has an initial guess of 0.
	if n%2 == 1 {
		guesses[0] = 0
		guesses = guesses[1:]
	}

	m := n / 2
	a := -0.5
	if n%2 == 1 {
		a = 0.5
	}
	nu := 4*float64(m) + 2*a + 2

	// Find the split between Gatteschi guesses and Tricomi guesses.
	p := 0.4985 + math.SmallestNonzeroFloat64
	pidx := int(math.Floor(p * float64(n)))

	// Use the Tricomi initial guesses in the first half where x is nearer to zero.
	// Note: zeros of besselj(+/-.5,x) are integer and half-integer multiples of pi.
	for i := 0; i < pidx; i++ {
		rhs := math.Pi * (4*float64(m) - 4*(float64(i)+1) + 3) / nu
		tnk := math.Pi / 2
		for k := 0; k < 7; k++ {
			val := tnk - math.Sin(tnk) - rhs
			dval := 1 - math.Cos(tnk)
			dTnk := val / dval
			tnk -= dTnk
			if math.Abs(dTnk) < 1e-14 {
				break
			}
		}
		vc := math.Cos(tnk / 2)
		t := vc * vc
		guesses[i] = math.Sqrt(nu*t - (5.0/(4.0*(1-t)*(1-t))-1.0/(1-t)-1+3*a*a)/3/nu)
	}

	// Use Gatteschi guesses in the second half where x is nearer to sqrt(n+0.5)
	for i := 0; i < m-pidx; i++ {
		var ar float64
		if i < len(airyRtsExact) {
			ar = airyRtsExact[i]
		} else {
			t := 3.0 / 8 * math.Pi * (4*(float64(i)+1) - 1)
			ar = math.Pow(t, 2.0/3) * (1 +
				5.0/48*math.Pow(t, -2) -
				5.0/36*math.Pow(t, -4) +
				77125.0/82944*math.Pow(t, -6) -
				108056875.0/6967296*math.Pow(t, -8) +
				162375596875.0/334430208*math.Pow(t, -10))
		}
		r := nu + math.Pow(2, 2.0/3)*ar*math.Pow(nu, 1.0/3) +
			0.2*math.Pow(2, 4.0/3)*ar*ar*math.Pow(nu, -1.0/3) +
			(11.0/35-a*a-12.0/175*ar*ar*ar)/nu +
			(16.0/1575*ar+92.0/7875*math.Pow(ar, 4))*math.Pow(2, 2.0/3)*math.Pow(nu, -5.0/3) -
			(15152.0/3031875*math.Pow(ar, 5)+1088.0/121275*ar*ar)*math.Pow(2, 1.0/3)*math.Pow(nu, -7.0/3)
		if r < 0 {
			ar = 0
		} else {
			ar = math.Sqrt(r)
		}
		guesses[m-1-i] = ar
	}
}

// airyRtsExact are the first airy roots.
var airyRtsExact = []float64{
	-2.338107410459762,
	-4.087949444130970,
	-5.520559828095555,
	-6.786708090071765,
	-7.944133587120863,
	-9.022650853340979,
	-10.040174341558084,
	-11.008524303733260,
	-11.936015563236262,
	-12.828776752865757,
}
