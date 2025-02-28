// Copyright ©2025 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathext

import (
	"math"

	"gonum.org/v1/gonum/floats/scalar"
	"gonum.org/v1/gonum/mathext/internal/cephes"
)

// Hypergeo returns the value of the Gaussian Hypergeometric function at z.
// For |z| < 1, this implementation follows the Cephes library.
// For |z| > 1, this implementation performs analytic continuation via relevant Hypergeometric identities.
//
// See https://en.wikipedia.org/wiki/Hypergeometric_function for more details.
func Hypergeo(a, b, c, z float64) float64 {
	// Simplify discussion by ensuring |a| > |b|.
	if math.Abs(b) > math.Abs(a) {
		return Hypergeo(b, a, c, z)
	}
	// Fix numerical issues for large |a| and |c| using equations 15.2.10 and 15.2.12.
	// M. Abramowitz and I. A. Stegun 1965. Handbook of Mathematical Functions, New York: Dover.
	if a < -50 && isInt(c) {
		return eqn15_2_10(a, b, c, z)
	}
	if c < -50 && isPosInt(a) {
		return eqn15_2_12(a, b, c, z)
	}

	// Canonical case |z| < 1 supported by Cephes.
	if math.Abs(z) < 1 {
		return cephes.Hyp2f1(a, b, c, z)
	}

	// Function undefined between the 1 and Inf branch points.
	if z > 0 {
		return math.NaN()
	}

	// Special case for z=-1, use equation 15.3.4 in Abramowitz.
	if z == -1 {
		return eqn15_3_4(a, b, c, z)
	}

	// When a or b is a non-positive integer, Hypergeo reduces to the finite polynomial described in equation 15.4.1.
	if isNonPosInt(a) {
		return eqn15_4_1(int(-a), b, c, z)
	}
	if isNonPosInt(b) {
		return eqn15_4_1(int(-b), a, c, z)
	}

	// Analytic continuation formula contains NaNs from Gamma(a-b) when a-b is an integer.
	// Fix this by making a and b different via equations 15.3.3 and 15.3.4.
	if isNonPosInt(c-a) || isNonPosInt(c-b) {
		return eqn15_3_3(a, b, c, z)
	}
	if isNonPosInt(a - b) {
		return eqn15_3_4(a, b, c, z)
	}

	// Analytic continuation based on https://www.johndcook.com/blog/2021/11/03/escaping-the-unit-disk/
	y1 := cephes.Hyp2f1(a, 1-c+a, 1-b+a, 1/z)
	y1 *= math.Gamma(c) / math.Gamma(b) * math.Gamma(b-a) / math.Gamma(c-a) * math.Pow(-z, -a)
	y2 := cephes.Hyp2f1(b, 1-c+b, 1-a+b, 1/z)
	y2 *= math.Gamma(c) / math.Gamma(a) * math.Gamma(a-b) / math.Gamma(c-b) * math.Pow(-z, -b)
	return y1 + y2
}

func eqn15_2_10(a, b, c, x float64) float64 {
	t := a - math.Round(a)
	var f2 float64
	f1 := Hypergeo(t, b, c, x)
	var f0 float64
	if a < 0 {
		f0 = Hypergeo(t-1, b, c, x)
		t--
		for n := 1; n < int(math.Round(-a)); n++ {
			f2 = f1
			f1 = f0
			f0 = -(2*t-c-t*x+b*x)/(c-t)*f1 - t*(x-1)/(c-t)*f2
			t--
		}
	} else {
		f0 = Hypergeo(t+1, b, c, x)
		t++
		for n := 1; n < int(math.Round(a)); n++ {
			f2 = f1
			f1 = f0
			f0 = -((2*t-c-t*x+b*x)*f1 + (c-t)*f2) / (t * (x - 1))
			t++
		}
	}
	return f0
}

func eqn15_2_12(a, b, c, x float64) float64 {
	t := c - math.Round(c)
	var f2 float64
	f1 := Hypergeo(a, b, t, x)
	var f0 float64
	if c < 0 {
		f0 = Hypergeo(a, b, t-1, x)
		t--
		for n := 1; n < int(math.Round(-c)); n++ {
			f2 = f1
			f1 = f0
			f0 = -t*(t-1-(2*t-a-b-1)*x)/t/(t-1)/(x-1)*f1 - (t-a)*(t-b)*x/t/(t-1)/(x-1)*f2
			t--
		}
	} else {
		f0 = Hypergeo(a, b, t+1, x)
		t++
		for n := 1; n < int(math.Round(a)); n++ {
			f2 = f1
			f1 = f0
			f0 = -(t*(t-1-(2*t-a-b-1)*x)*f1 + t*(t-1)*(x-1)*f2) / ((t - a) * (t - b) * x)
			t++
		}
	}
	return f0
}

func eqn15_3_3(a, b, c, z float64) float64 {
	y := Hypergeo(c-a, c-b, c, z)
	return math.Pow(1-z, c-a-b) * y
}

func eqn15_3_4(a, b, c, z float64) float64 {
	y := Hypergeo(a, c-b, c, z/(z-1))
	return math.Pow(1-z, -a) * y
}

func eqn15_4_1(m int, b, c, z float64) float64 {
	// lzn is log(z^n / n!).
	var lzn float64
	var lznSign int = 1
	// sum is the sum of the evaluated polynomial.
	var sum float64 = 1
	for n := 1; n <= m; n++ {
		lzn += math.Log(math.Abs(z)) - math.Log(float64(n))
		if z < 0 {
			lznSign *= -1
		}

		var w float64
		var sign int = 1
		lphC, sn := lpochhammer(c, n)
		if sn == 0 { // divide by zero
			return math.NaN()
		}
		w -= lphC
		sign *= sn

		lphM, sn := lpochhammer(float64(-m), n)
		if sn == 0 {
			continue
		}
		w += lphM
		sign *= sn

		lphB, sn := lpochhammer(b, n)
		if sn == 0 {
			continue
		}
		w += lphB
		sign *= sn

		sum += float64(sign) * float64(lznSign) * math.Exp(w+lzn)
	}
	return sum
}

func lpochhammer(x float64, n int) (float64, int) {
	var y float64
	var sign int = 1
	for k := range n {
		xk := x - float64(k)
		if xk == 0 {
			return math.NaN(), 0
		}

		y += math.Log(math.Abs(xk))
		if xk < 0 {
			sign *= -1
		}
	}
	return y, sign
}

func isInt(x float64) bool {
	return scalar.EqualWithinAbs(math.Round(x), x, 1e-6)
}

func isPosInt(x float64) bool {
	return (x > 0) && isInt(x)
}

func isNonPosInt(x float64) bool {
	return (x <= 0) && isInt(x)
}
