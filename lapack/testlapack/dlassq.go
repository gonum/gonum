// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats"
)

type Dlassqer interface {
	Dlassq(n int, x []float64, incx int, scale, ssq float64) (float64, float64)
}

func DlassqTest(t *testing.T, impl Dlassqer) {
	values := []float64{
		0,
		2 * safmin,
		smlnum,
		ulp,
		1,
		1 / ulp,
		bignum,
		safmax,
		math.Inf(1),
		math.NaN(),
	}

	rnd := rand.New(rand.NewSource(1))
	for _, n := range []int{0, 1, 2, 3, 4, 5, 10, 20, 30, 40} {
		for _, incx := range []int{1, 3} {
			for cas := 0; cas < 3; cas++ {
				for _, v0 := range values {
					if v0 > 1 {
						v0 *= 0.5
					}
					for _, v1 := range values {
						if v1 > 1 {
							v1 = 0.5 * v1 / math.Sqrt(float64(n+1))
						}
						dlassqTest(t, impl, rnd, n, incx, cas, v0, v1)
					}
				}
			}
		}
	}
}

func dlassqTest(t *testing.T, impl Dlassqer, rnd *rand.Rand, n, incx, cas int, v0, v1 float64) {
	const (
		rogue = 1234.5678
		tol   = 1e-14
	)

	name := fmt.Sprintf("n=%v,incx=%v,cas=%v,v0=%v,v1=%v", n, incx, cas, v0, v1)

	// Generate n random values in (-1,1].
	work := make([]float64, n)
	for i := range work {
		work[i] = 1 - 2*rnd.Float64()
	}

	// Compute the sum of squares by an unscaled algorithm.
	var workssq float64
	for _, wi := range work {
		workssq += wi * wi
	}

	// Construct a test (n+1)-vector z as
	//  z[0] = known value
	//  z[1:n+1] = scaling factor * work[:]
	z := make([]float64, n+1)
	z[0] = math.Sqrt(v0)
	for i, wi := range work {
		z[i+1] = v1 * wi
	}
	// Set initial scale and ssq corresponding to z[0].
	var scale, ssq float64
	switch cas {
	case 0:
		scale = 1
		ssq = v0
	case 1:
		scale = math.Sqrt(v0)
		ssq = 1
	case 2:
		if v0 < 1 {
			scale = 1.0 / 3
			ssq = 9 * v0
		} else {
			scale = 3
			ssq = v0 / 9
		}
	default:
		panic("bad cas")
	}

	const (
		dtsml = 0x1p-511
		dtbig = 0x1p486
		dssml = 0x1p537
		dsbig = 0x1p-538
	)
	if scale*math.Sqrt(ssq) > dtbig && scale < math.Sqrt(math.SmallestNonzeroFloat64/ulp)/dsbig {
		// The scaled sum is big but the scale itself is small.
		return
	}
	if scale*math.Sqrt(ssq) < dtsml && scale > math.Sqrt(math.MaxFloat64)/dssml {
		// The scaled sum is small but the scale itself is big.
		return
	}

	// Compute the expected value of the sum of squares using the (n+1)-vector z.
	z0 := z[0]
	var z1n float64
	if n >= 1 {
		z1n = v1 * math.Sqrt(workssq)
	}
	zmin := math.Min(z0, z1n)
	zmax := math.Max(z0, z1n)
	var nrmWant float64
	switch {
	case math.IsNaN(z0) || math.IsNaN(z1n):
		nrmWant = math.NaN()
	case zmin == zmax:
		nrmWant = math.Sqrt2 * zmax
	case zmax == 0:
		nrmWant = 0
	default:
		nrmWant = zmax * math.Sqrt(1+(zmin/zmax)*(zmin/zmax))
	}

	// Allocate input slice for Dlassq and fill it with z[1:].
	x := make([]float64, max(0, 1+(n-1)*incx))
	for i := range x {
		x[i] = rogue
	}
	for i, zi := range z[1:] {
		x[i*incx] = zi
	}
	xCopy := make([]float64, len(x))
	copy(xCopy, x)

	scaleGot, ssqGot := impl.Dlassq(n, x, incx, scale, ssq)
	nrmGot := scaleGot * math.Sqrt(ssqGot)

	if !floats.Same(x, xCopy) {
		t.Fatalf("%v: unexpected modification of x", name)
	}

	// Check the result.
	switch {
	case math.IsNaN(nrmGot) || math.IsNaN(nrmWant):
		if !math.IsNaN(nrmGot) {
			t.Errorf("%v: expected NaN; got %v", name, nrmGot)
		}
		if !math.IsNaN(nrmWant) {
			t.Errorf("%v: unexpected NaN; want %v", name, nrmWant)
		}
	case nrmGot == nrmWant:
	case nrmWant == 0:
		if nrmGot > tol {
			t.Errorf("%v: unexpected result; got %v, want 0", name, nrmGot)
		}
	default:
		diff := math.Abs(nrmGot-nrmWant) / nrmWant / math.Max(1, float64(n))
		if math.IsNaN(diff) || diff > tol {
			t.Errorf("%v: unexpected result; got %v, want %v", name, nrmGot, nrmWant)
		}
	}
}
