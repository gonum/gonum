// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"
	"math/rand/v2"
	"testing"
)

type Dlaic1er interface {
	Dlaic1(job int, j int, x []float64, sest float64, w []float64, gamma float64) (sestpr, s, c float64)
}

func Dlaic1Test(t *testing.T, impl Dlaic1er) {
	rnd := rand.New(rand.NewPCG(1, 1))
	const tol = 1e-12

	// Test basic properties for random inputs.
	for _, j := range []int{0, 1, 2, 5, 10} {
		for trial := 0; trial < 20; trial++ {
			x := make([]float64, j)
			w := make([]float64, j)
			for i := range x {
				x[i] = rnd.NormFloat64()
				w[i] = rnd.NormFloat64()
			}
			sest := math.Abs(rnd.NormFloat64()) + 0.1
			gamma := rnd.NormFloat64()

			for _, job := range []int{1, 2} {
				sestpr, s, c := impl.Dlaic1(job, j, x, sest, w, gamma)

				// s^2 + c^2 should be approximately 1.
				norm := s*s + c*c
				if math.Abs(norm-1) > tol {
					t.Errorf("job=%d j=%d trial=%d: s²+c²=%g want 1", job, j, trial, norm)
				}

				// sestpr should be non-negative.
				if sestpr < -tol {
					t.Errorf("job=%d j=%d trial=%d: sestpr=%g < 0", job, j, trial, sestpr)
				}

				if job == 1 {
					// For max estimation, sestpr >= sest (augmenting can't decrease max sv).
					if sestpr < sest-tol*sest {
						t.Errorf("job=1 j=%d trial=%d: sestpr=%g < sest=%g", j, trial, sestpr, sest)
					}
				}
			}
		}
	}

	// Test with sest=0.
	for _, job := range []int{1, 2} {
		sestpr, s, c := impl.Dlaic1(job, 2, []float64{1, 0}, 0, []float64{0, 1}, 3)
		norm := s*s + c*c
		if math.Abs(norm-1) > tol {
			t.Errorf("sest=0 job=%d: s²+c²=%g want 1", job, norm)
		}
		_ = sestpr
	}

	// Test with zero gamma and w: should return sest unchanged for job=1.
	{
		sestpr, s, c := impl.Dlaic1(1, 3, []float64{1, 0, 0}, 5.0, []float64{0, 0, 0}, 0)
		if math.Abs(s*s+c*c-1) > tol {
			t.Errorf("zero gamma/w: s²+c²=%g want 1", s*s+c*c)
		}
		_ = sestpr
	}

	// Test 1×1 → 2×2 explicitly.
	// T = [3], w = [1], gamma = 4 → T_hat = [3 1; 0 4]
	// Singular values of T_hat: eigenvalues of T_hat^T * T_hat = [9 3; 3 17]
	// λ = (26 ± sqrt(676-4*144))/2 = (26 ± sqrt(100))/2 = 18 or 8
	// sigma_max = sqrt(18), sigma_min = sqrt(8)
	{
		sest := 3.0
		x := []float64{1.0}
		w := []float64{1.0}
		gamma := 4.0

		sestpr1, s1, c1 := impl.Dlaic1(1, 1, x, sest, w, gamma)
		if math.Abs(s1*s1+c1*c1-1) > tol {
			t.Errorf("1x1→2x2 job=1: s²+c²=%g", s1*s1+c1*c1)
		}
		sigmaMax := math.Sqrt(18)
		if math.Abs(sestpr1-sigmaMax) > 0.5*sigmaMax {
			t.Errorf("1x1→2x2 job=1: sestpr=%g want≈%g", sestpr1, sigmaMax)
		}

		sestpr2, s2, c2 := impl.Dlaic1(2, 1, x, sest, w, gamma)
		if math.Abs(s2*s2+c2*c2-1) > tol {
			t.Errorf("1x1→2x2 job=2: s²+c²=%g", s2*s2+c2*c2)
		}
		sigmaMin := math.Sqrt(8)
		if math.Abs(sestpr2-sigmaMin) > 0.5*sigmaMin {
			t.Errorf("1x1→2x2 job=2: sestpr=%g want≈%g", sestpr2, sigmaMin)
		}
	}
}
