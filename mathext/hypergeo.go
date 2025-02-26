// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathext

import (
	"math"

	"gonum.org/v1/gonum/mathext/internal/cephes"
)

// Hypergeo returns the value of the Gaussian Hypergeometric function at z.
// For |z| < 1, this implementation follows the Cephes library.
// For |z| > 1, this implementation perform analytic continuation via relevant Hypergeometric identities.
// See https://en.wikipedia.org/wiki/Hypergeometric_function for more details.
func Hypergeo(a float64, b float64, c float64, z float64) float64 {
	if math.Abs(z) < 1 {
		return cephes.Hyp2f1(a, b, c, z)
	}

	// Function undefined between the 1 and Inf branch points.
	if z > 0 {
		return math.NaN()
	}

	// Analytic continuation formula contains infinities from Gamma(a-b) when a == b.
	// Fix this by making a and b different using equation 15.3.4 from Abramowitz.
	// M. Abramowitz and I. A. Stegun 1965. Handbook of Mathematical Functions, New York: Dover.
	if a == b {
		y := Hypergeo(a, c-b, c, z/(z-1))
		return math.Pow(1-z, -a) * y
	}

	// Analytic continuation based on https://www.johndcook.com/blog/2021/11/03/escaping-the-unit-disk/
	y1 := cephes.Hyp2f1(a, 1-c+a, 1-b+a, 1/z)
	y1 *= math.Gamma(c) / math.Gamma(b) * math.Gamma(b-a) / math.Gamma(c-a) * math.Pow(-z, -a)
	y2 := cephes.Hyp2f1(b, 1-c+b, 1-a+b, 1/z)
	y2 *= math.Gamma(c) / math.Gamma(a) * math.Gamma(a-b) / math.Gamma(c-b) * math.Pow(-z, -b)
	return y1 + y2
}
