// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathext

import (
	"math"
	"math/rand"
	"testing"
)

// Testing EllipticF (and CarlsonRF) using the addition theorems from http://dlmf.nist.gov/19.11.i
func TestEllipticF(t *testing.T) {
	const tol = 1.0e-14
	rng := rand.New(rand.NewSource(1))

	for test := 0; test < 100; test++ {
		alpha := rng.Float64() * math.Pi / 4
		beta := rng.Float64() * math.Pi / 4
		for mi := 0; mi < 9999; mi++ {
			m := float64(mi) / 10000
			Fa := EllipticF(alpha, m)
			Fb := EllipticF(beta, m)
			sina, cosa := math.Sincos(alpha)
			sinb, cosb := math.Sincos(beta)
			tan := (sina*math.Sqrt(1-m*sinb*sinb) + sinb*math.Sqrt(1-m*sina*sina)) / (cosa + cosb)
			gamma := 2 * math.Atan(tan)
			Fg := EllipticF(gamma, m)
			delta := math.Abs(Fa + Fb - Fg)
			if delta > tol {
				t.Fatalf("EllipticF test fail for m=%v, alpha=%v, beta=%v", m, alpha, beta)
			}
		}
	}
}

// Testing EllipticE (and CarlsonRF, CarlsonRD) using the addition theorems from http://dlmf.nist.gov/19.11.i
func TestEllipticE(t *testing.T) {
	const tol = 1.0e-14
	rng := rand.New(rand.NewSource(1))

	for test := 0; test < 100; test++ {
		alpha := rng.Float64() * math.Pi / 4
		beta := rng.Float64() * math.Pi / 4
		for mi := 0; mi < 9999; mi++ {
			m := float64(mi) / 10000
			Ea := EllipticE(alpha, m)
			Eb := EllipticE(beta, m)
			sina, cosa := math.Sincos(alpha)
			sinb, cosb := math.Sincos(beta)
			tan := (sina*math.Sqrt(1-m*sinb*sinb) + sinb*math.Sqrt(1-m*sina*sina)) / (cosa + cosb)
			gamma := 2 * math.Atan(tan)
			Eg := EllipticE(gamma, m)
			delta := math.Abs(Ea + Eb - Eg - m*sina*sinb*math.Sin(gamma))
			if delta > tol {
				t.Fatalf("EllipticE test fail for m=%v, alpha=%v, beta=%v", m, alpha, beta)
			}
		}
	}
}
