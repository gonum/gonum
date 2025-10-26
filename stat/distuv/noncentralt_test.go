// Copyright ©2025 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"fmt"
	"math"
	"math/rand/v2"
	"sort"
	"testing"

	"gonum.org/v1/gonum/floats/scalar"
)

func TestNoncentralTRand(t *testing.T) {
	t.Parallel()
	const (
		n    = 3e5
		bins = 50
	)
	rsrc := rand.NewPCG(1, 1)
	tests := []struct {
		dist NoncentralT
		tol  float64
	}{
		{dist: NoncentralT{Nu: 58, Mu: 1.9364916731037085, Src: rsrc}, tol: 1e-2},
		{dist: NoncentralT{Nu: 198, Mu: 7.0710678118654755, Src: rsrc}, tol: 1e-2},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			x := make([]float64, n)
			generateSamples(x, test.dist)
			sort.Float64s(x)

			testRandLogProbContinuous(t, i, math.Inf(-1), x, test.dist, test.tol, bins)
		})
	}
}

func TestNoncentralTMoments(t *testing.T) {
	t.Parallel()
	tests := []struct {
		dist     NoncentralT
		mean     float64
		variance float64
		tol      float64
	}{
		{dist: NoncentralT{Nu: 20, Mu: 23}, mean: 23.909902487370097, variance: 17.205451933342147, tol: 1e-13},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			mean := test.dist.Mean()
			if !scalar.EqualWithinRel(mean, test.mean, test.tol) {
				t.Errorf("{%f, %f}.Mean: got %f want %f", test.dist.Nu, test.dist.Mu, mean, test.mean)
			}
			variance := test.dist.Variance()
			if !scalar.EqualWithinRel(variance, test.variance, test.tol) {
				t.Errorf("{%f, %f}.Variance: got %f want %f", test.dist.Nu, test.dist.Mu, variance, test.variance)
			}
		})
	}
}

func TestNoncentralTProb(t *testing.T) {
	t.Parallel()
	tests := []struct {
		dist NoncentralT
		x    float64
		p    float64
		tol  float64
		abs  float64
	}{
		// Based on https://github.com/wch/r-source/blob/trunk/tests/reg-tests-2.R
		{dist: NoncentralT{Nu: 10, Mu: 0}, x: 1.8, p: 0.08311638965387959, tol: 5e-12},
		{dist: NoncentralT{Nu: 10, Mu: 0.0001}, x: 1.8, p: 0.0831297211093226, tol: 5e-12},
		{dist: NoncentralT{Nu: 10, Mu: 1}, x: 1.8, p: 0.2665039326310101, tol: 5e-12},
		{dist: NoncentralT{Nu: 10, Mu: -0.0001}, x: 1.8, p: 0.08310305960463463, tol: 5e-12},
		{dist: NoncentralT{Nu: 10, Mu: -1}, x: 1.8, p: 0.01074629362995775, tol: 5e-12},

		// Based on https://github.com/boostorg/math/blob/develop/test/test_nc_t.hpp
		{dist: NoncentralT{Nu: 8, Mu: 12}, x: 12, p: 1.23532971542589493e-1, tol: 5e-12},
		{dist: NoncentralT{Nu: 126, Mu: -2}, x: -4, p: 5.79793228936581470e-2, tol: 2e-11},
		{dist: NoncentralT{Nu: 126, Mu: 2}, x: 4, p: 5.79793228936581470e-2, tol: 2e-11},
		{dist: NoncentralT{Nu: 126, Mu: 2}, x: 0, p: 5.38839489063995713e-2, tol: 5e-12},
		{dist: NoncentralT{Nu: 8, Mu: 16}, x: 0, p: 9.94670846108541165e-57, tol: 5e-12},
		{dist: NoncentralT{Nu: 8, Mu: 16}, x: 0x1p-1022, p: 9.94670846108541165e-57, tol: 5e-12},
		{dist: NoncentralT{Nu: 8, Mu: 16}, x: -0x1p-1022, p: 9.94670846108541165e-57, tol: 5e-12},
		{dist: NoncentralT{Nu: 8, Mu: 16}, x: -0.125, p: 1.40958893993909266e-57, tol: 5e-12, abs: 1e-16},
		{dist: NoncentralT{Nu: 8, Mu: 16}, x: -1e-16, p: 9.94670846108539523e-57, tol: 5e-12},
		{dist: NoncentralT{Nu: 8, Mu: 16}, x: 1e-16, p: 9.94670846108542807e-57, tol: 5e-12},
		{dist: NoncentralT{Nu: 8, Mu: 16}, x: 0.125, p: 8.78741270305645722e-56, tol: 5e-4},

		// Custom tests, wanted values from the R language version 4.4.2.
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: -0.3930852906078905, p: 0.02631178912036441, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: 0.6553966734339551, p: 0.1756830868490319, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: 1.498837237776003, p: 0.3591385592167584, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: 1.944945574355751, p: 0.3908781127965009, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: 2.369093423208461, p: 0.3534073785987947, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: 3.290148485967974, p: 0.1615374126723982, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: 4.469416864177883, p: 0.02264358202648499, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: 0, p: 0.0609166161931202, tol: 5e-12},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			p := test.dist.Prob(test.x)
			if !scalar.EqualWithinAbsOrRel(p, test.p, test.abs, test.tol) {
				t.Errorf("{%f, %f}.Prob(%f): got %g want %g", test.dist.Nu, test.dist.Mu, test.x, p, test.p)
			}
		})
	}
}

func TestNoncentralTCDF(t *testing.T) {
	t.Parallel()
	tests := []struct {
		dist NoncentralT
		x    float64
		cdf  float64
		tol  float64
		abs  float64
	}{
		// Based on https://github.com/wch/r-source/blob/trunk/tests/reg-tests-2.R
		{dist: NoncentralT{Nu: 10, Mu: 0}, x: 1.8, cdf: 0.9489738784326605, tol: 5e-12},
		{dist: NoncentralT{Nu: 10, Mu: 0.0001}, x: 1.8, cdf: 0.948964072175642, tol: 5e-12},
		{dist: NoncentralT{Nu: 10, Mu: 1}, x: 1.8, cdf: 0.7584267206837773, tol: 5e-12},
		{dist: NoncentralT{Nu: 10, Mu: -0.0001}, x: 1.8, cdf: 0.9489836831935839, tol: 5e-12},
		{dist: NoncentralT{Nu: 10, Mu: -1}, x: 1.8, cdf: 0.9949471996805094, tol: 5e-12},

		// Based on https://github.com/boostorg/math/blob/develop/test/scipy_issue_14901.cpp
		{dist: NoncentralT{Nu: 2, Mu: 2}, x: 0.05, cdf: 0.02528206132724582, tol: 5e-12},

		{dist: NoncentralT{Nu: 1, Mu: 3}, x: 0.05, cdf: 0.00154456589169420, tol: 5e-11},

		// Based on https://github.com/boostorg/math/blob/develop/test/scipy_issue_17916_nct.cpp
		{dist: NoncentralT{Nu: 2, Mu: 482023264}, x: 2, cdf: 0, tol: 5e-12},

		// Based on https://github.com/boostorg/math/blob/develop/test/test_nc_t.hpp
		{dist: NoncentralT{Nu: 3, Mu: 1}, x: 2.34, cdf: 0.801888999613917, tol: 5e-12},
		{dist: NoncentralT{Nu: 126, Mu: -2}, x: -4.33, cdf: 1.252846196792878e-2, tol: 5e-11},
		{dist: NoncentralT{Nu: 20, Mu: 23}, x: 23, cdf: 0.460134400391924, tol: 5e-12},
		{dist: NoncentralT{Nu: 20, Mu: 33}, x: 34, cdf: 0.532008386378725, tol: 5e-12},
		{dist: NoncentralT{Nu: 12, Mu: 38}, x: 39, cdf: 0.495868184917805, tol: 5e-12},
		{dist: NoncentralT{Nu: 12, Mu: 39}, x: 39, cdf: 0.446304024668836, tol: 5e-2},
		{dist: NoncentralT{Nu: 200, Mu: 38}, x: 39, cdf: 0.666194209961795, tol: 5e-11},
		{dist: NoncentralT{Nu: 200, Mu: 42}, x: 40, cdf: 0.179292265426085, tol: 2e-3},
		{dist: NoncentralT{Nu: 2, Mu: 4}, x: 5, cdf: 0.532020698669953, tol: 5e-12},
		{dist: NoncentralT{Nu: 8, Mu: 16}, x: 0, cdf: 6.388754400538087e-58, tol: 5e-12},
		{dist: NoncentralT{Nu: 8, Mu: 16}, x: 0x1p-1022, cdf: 6.388754400538087e-58, tol: 5e-12},
		{dist: NoncentralT{Nu: 8, Mu: 16}, x: -0x1p-1022, cdf: 6.388754400538087e-58, tol: 5e-12, abs: 1e-16},
		{dist: NoncentralT{Nu: 8, Mu: 16}, x: -0.125, cdf: 1.018937769092816e-58, tol: 5e-12, abs: 1e-16},
		{dist: NoncentralT{Nu: 8, Mu: 16}, x: -1e-16, cdf: 6.388754400538077e-58, tol: 5e-12, abs: 1e-16},
		{dist: NoncentralT{Nu: 8, Mu: 16}, x: 1e-16, cdf: 6.388754400538097e-58, tol: 5e-12},
		{dist: NoncentralT{Nu: 8, Mu: 16}, x: 0.125, cdf: 5.029904883914148e-57, tol: 1e-4},
		{dist: NoncentralT{Nu: 8, Mu: 8.5}, x: -1, cdf: 6.174794808375702e-20, tol: 2e-5, abs: 1e-16},

		// Custom tests, wanted values from the R language version 4.4.2.
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: -0.3930852906078905, cdf: 0.01, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: 0.6553966734339551, cdf: 0.1, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: 1.498837237776003, cdf: 0.33, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: 1.944945574355751, cdf: 0.5, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: 2.369093423208461, cdf: 0.66, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: 3.290148485967974, cdf: 0.9, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: 4.469416864177883, cdf: 0.99, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Mu: 3.535534}, x: 1.206434338714708, cdf: 0.01, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Mu: 3.535534}, x: 2.248666397791987, cdf: 0.1, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Mu: 3.535534}, x: 3.094270376888891, cdf: 0.33, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Mu: 3.535534}, x: 3.54004668094423, cdf: 0.5, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Mu: 3.535534}, x: 3.961130093860616, cdf: 0.66, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Mu: 3.535534}, x: 4.860775618489509, cdf: 0.9, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Mu: 3.535534}, x: 5.97083270110298, cdf: 0.99, tol: 5e-12},
		{dist: NoncentralT{Nu: 10, Mu: 38.5}, x: 40, cdf: 0.5073909173564686, tol: 5e-3},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			cdf := test.dist.CDF(test.x)
			if !scalar.EqualWithinAbsOrRel(cdf, test.cdf, test.abs, test.tol) {
				t.Errorf("{Nu: %f, Mu: %f}.CDF(%f): got %g want %g", test.dist.Nu, test.dist.Mu, test.x, cdf, test.cdf)
			}
		})
	}
}

func TestNoncentralTQuantile(t *testing.T) {
	t.Parallel()
	tests := []struct {
		dist NoncentralT
		x    float64
		cdf  float64
		tol  float64
		abs  float64
	}{
		// Based on https://github.com/wch/r-source/blob/trunk/tests/reg-tests-2.R
		{dist: NoncentralT{Nu: 10, Mu: 0}, x: 1.812461122811676, cdf: 0.95, tol: 5e-12},
		{dist: NoncentralT{Nu: 10, Mu: 0.0001}, x: 1.812579296911650, cdf: 0.95, tol: 5e-12},
		{dist: NoncentralT{Nu: 10, Mu: 1}, x: 3.041741814971971, cdf: 0.95, tol: 5e-12},
		{dist: NoncentralT{Nu: 10, Mu: -0.0001}, x: 1.8123429496892811, cdf: 0.95, tol: 5e-12},
		{dist: NoncentralT{Nu: 10, Mu: -1}, x: 0.6797901881827499, cdf: 0.95, tol: 5e-12},

		// Based on https://github.com/boostorg/math/blob/develop/test/test_nc_t.hpp
		{dist: NoncentralT{Nu: 3, Mu: 1}, x: 2.34, cdf: 0.801888999613917, tol: 5e-12},
		{dist: NoncentralT{Nu: 126, Mu: -2}, x: -4.33, cdf: 1.252846196792878e-2, tol: 5e-11},
		{dist: NoncentralT{Nu: 20, Mu: 23}, x: 23, cdf: 0.460134400391924, tol: 5e-12},
		{dist: NoncentralT{Nu: 20, Mu: 33}, x: 34, cdf: 0.532008386378725, tol: 5e-12},
		{dist: NoncentralT{Nu: 12, Mu: 38}, x: 39, cdf: 0.495868184917805, tol: 5e-12},
		{dist: NoncentralT{Nu: 12, Mu: 39}, x: 39, cdf: 0.446304024668836, tol: 5e-2},
		{dist: NoncentralT{Nu: 200, Mu: 38}, x: 39, cdf: 0.666194209961795, tol: 5e-11},
		{dist: NoncentralT{Nu: 200, Mu: 42}, x: 40, cdf: 0.179292265426085, tol: 2e-3},
		{dist: NoncentralT{Nu: 2, Mu: 4}, x: 5, cdf: 0.532020698669953, tol: 5e-12},

		// Custom tests, wanted values from the R language version 4.4.2.
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: -0.3930852906078905, cdf: 0.01, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: 0.6553966734339551, cdf: 0.1, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: 1.498837237776003, cdf: 0.33, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: 1.944945574355751, cdf: 0.5, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: 2.369093423208461, cdf: 0.66, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: 3.290148485967974, cdf: 0.9, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Mu: 1.936492}, x: 4.469416864177883, cdf: 0.99, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Mu: 3.535534}, x: 1.206434338714708, cdf: 0.01, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Mu: 3.535534}, x: 2.248666397791987, cdf: 0.1, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Mu: 3.535534}, x: 3.094270376888891, cdf: 0.33, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Mu: 3.535534}, x: 3.54004668094423, cdf: 0.5, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Mu: 3.535534}, x: 3.961130093860616, cdf: 0.66, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Mu: 3.535534}, x: 4.860775618489509, cdf: 0.9, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Mu: 3.535534}, x: 5.97083270110298, cdf: 0.99, tol: 5e-12},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			x := test.dist.Quantile(test.cdf)
			if !scalar.EqualWithinAbsOrRel(x, test.x, test.abs, test.tol) {
				t.Errorf("{Nu: %f, Mu: %f}.Quantile(%f): got %g want %g", test.dist.Nu, test.dist.Mu, test.cdf, x, test.x)
			}
		})
	}
}
