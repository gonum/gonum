// Copyright ©2025 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathext

import (
	"math"

	"gonum.org/v1/gonum/mathext/internal/cephes"
)

// Hypergeo returns the value of the Gaussian Hypergeometric function at z.
// For |z| < 1, this implementation follows the Cephes library.
// For |z| > 1, this implementation performs analytic continuation via relevant Hypergeometric identities.
//
// See https://en.wikipedia.org/wiki/Hypergeometric_function for more details.
func Hypergeo(a, b, c, z float64) float64 {
	if math.Abs(z) < 1 {
		return cephes.Hyp2f1(a, b, c, z)
	}

	// Function undefined between the 1 and Inf branch points.
	if z > 0 {
		return math.NaN()
	}

	// When a or b is a negative integer, Hypergeo reduces to the polynomial described in equation 15.4.1 in Abramowitz.
	// M. Abramowitz and I. A. Stegun 1965. Handbook of Mathematical Functions, New York: Dover.
	if isNonPosInt(a) {
		return eqn15_4_1(int(-a), b, c, z)
	}
	if isNonPosInt(b) {
		return eqn15_4_1(int(-b), a, c, z)
	}

	// Analytic continuation formula contains NaNs from Gamma(a-b) when a-b is an integer.
	// Fix this by making a and b different using equation 15.3.4 from Abramowitz.
	if isNonPosInt(a - b) {
		y := Hypergeo(a, c-b, c, z/(z-1))
		return math.Pow(1-z, -a) * y
	}
	if isNonPosInt(c-a) || isNonPosInt(c-b) {
		y := Hypergeo(c-a, c-b, c, z)
		return math.Pow(1-z, c-a-b) * y
	}

	// Analytic continuation based on https://www.johndcook.com/blog/2021/11/03/escaping-the-unit-disk/
	y1 := cephes.Hyp2f1(a, 1-c+a, 1-b+a, 1/z)
	y1 *= math.Gamma(c) / math.Gamma(b) * math.Gamma(b-a) / math.Gamma(c-a) * math.Pow(-z, -a)
	y2 := cephes.Hyp2f1(b, 1-c+b, 1-a+b, 1/z)
	y2 *= math.Gamma(c) / math.Gamma(a) * math.Gamma(a-b) / math.Gamma(c-b) * math.Pow(-z, -b)
	return y1 + y2
}

func eqn15_4_1(m int, b, c, z float64) float64 {
	// zn is (z^n / n!).
	var zn float64 = 1
	// sum is the sum of the evaluated polynomial.
	var sum float64 = 1
	for n := 1; n <= m; n++ {
		zn *= z / float64(n)
		sum += pochhammer(float64(-m), n) * pochhammer(b, n) / pochhammer(c, n) * zn
	}
	return sum
}

func pochhammer(x float64, n int) float64 {
	var y float64 = 1
	for k := range n {
		y *= x - float64(k)
	}
	return y
}

func isNonPosInt(x float64) bool {
	return (x <= 0) && (math.Round(x) == x)
}
