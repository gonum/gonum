// Copyright ©2013 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats"
)

func TestSVD(t *testing.T) {
	t.Parallel()
	rnd := rand.New(rand.NewSource(1))
	// Hand coded tests
	for _, test := range []struct {
		a *Dense
		u *Dense
		v *Dense
		s []float64
	}{
		{
			a: NewDense(4, 2, []float64{2, 4, 1, 3, 0, 0, 0, 0}),
			u: NewDense(4, 2, []float64{
				-0.8174155604703632, -0.5760484367663209,
				-0.5760484367663209, 0.8174155604703633,
				0, 0,
				0, 0,
			}),
			v: NewDense(2, 2, []float64{
				-0.4045535848337571, -0.9145142956773044,
				-0.9145142956773044, 0.4045535848337571,
			}),
			s: []float64{5.464985704219041, 0.365966190626258},
		},
		{
			// Issue #5.
			a: NewDense(3, 11, []float64{
				1, 1, 0, 1, 0, 0, 0, 0, 0, 11, 1,
				1, 0, 0, 0, 0, 0, 1, 0, 0, 12, 2,
				1, 1, 0, 0, 0, 0, 0, 0, 1, 13, 3,
			}),
			u: NewDense(3, 3, []float64{
				-0.5224167862273765, 0.7864430360363114, 0.3295270133658976,
				-0.5739526766688285, -0.03852203026050301, -0.8179818935216693,
				-0.6306021141833781, -0.6164603833618163, 0.4715056408282468,
			}),
			v: NewDense(11, 3, []float64{
				-0.08123293141915189, 0.08528085505260324, -0.013165501690885152,
				-0.05423546426886932, 0.1102707844980355, 0.622210623111631,
				0, 0, 0,
				-0.0245733326078166, 0.510179651760153, 0.25596360803140994,
				0, 0, 0,
				0, 0, 0,
				-0.026997467150282436, -0.024989929445430496, -0.6353761248025164,
				0, 0, 0,
				-0.029662131661052707, -0.3999088672621176, 0.3662470150802212,
				-0.9798839760830571, 0.11328174160898856, -0.047702613241813366,
				-0.16755466189153964, -0.7395268089170608, 0.08395240366704032,
			}),
			s: []float64{21.259500881097434, 1.5415021616856566, 1.2873979074613628},
		},
	} {
		var svd SVD
		ok := svd.Factorize(test.a, SVDThin)
		if !ok {
			t.Errorf("SVD failed")
		}
		s, u, v := extractSVD(&svd)
		if !floats.EqualApprox(s, test.s, 1e-10) {
			t.Errorf("Singular value mismatch. Got %v, want %v.", s, test.s)
		}
		if !EqualApprox(u, test.u, 1e-10) {
			t.Errorf("U mismatch.\nGot:\n%v\nWant:\n%v", Formatted(u), Formatted(test.u))
		}
		if !EqualApprox(v, test.v, 1e-10) {
			t.Errorf("V mismatch.\nGot:\n%v\nWant:\n%v", Formatted(v), Formatted(test.v))
		}
		m, n := test.a.Dims()
		sigma := NewDense(min(m, n), min(m, n), nil)
		for i := 0; i < min(m, n); i++ {
			sigma.Set(i, i, s[i])
		}

		var ans Dense
		ans.Product(u, sigma, v.T())
		if !EqualApprox(test.a, &ans, 1e-10) {
			t.Errorf("A reconstruction mismatch.\nGot:\n%v\nWant:\n%v\n", Formatted(&ans), Formatted(test.a))
		}

		for _, kind := range []SVDKind{
			SVDThinU, SVDFullU, SVDThinV, SVDFullV,
		} {
			var svd SVD
			svd.Factorize(test.a, kind)
			if kind&SVDThinU == 0 && kind&SVDFullU == 0 {
				panicked, message := panics(func() {
					var dst Dense
					svd.UTo(&dst)
				})
				if !panicked {
					t.Error("expected panic with no U matrix requested")
					continue
				}
				want := "svd: u not computed during factorization"
				if message != want {
					t.Errorf("unexpected message: got:%q want:%q", message, want)
				}
			}
			if kind&SVDThinV == 0 && kind&SVDFullV == 0 {
				panicked, message := panics(func() {
					var dst Dense
					svd.VTo(&dst)
				})
				if !panicked {
					t.Error("expected panic with no V matrix requested")
					continue
				}
				want := "svd: v not computed during factorization"
				if message != want {
					t.Errorf("unexpected message: got:%q want:%q", message, want)
				}
			}
		}
	}

	for _, test := range []struct {
		m, n int
	}{
		{5, 5},
		{5, 3},
		{3, 5},
		{150, 150},
		{200, 150},
		{150, 200},
	} {
		m := test.m
		n := test.n
		for trial := 0; trial < 10; trial++ {
			a := NewDense(m, n, nil)
			for i := range a.mat.Data {
				a.mat.Data[i] = rnd.NormFloat64()
			}
			aCopy := DenseCopyOf(a)

			// Test Full decomposition.
			var svd SVD
			ok := svd.Factorize(a, SVDFull)
			if !ok {
				t.Errorf("SVD factorization failed")
			}
			if !Equal(a, aCopy) {
				t.Errorf("A changed during call to SVD with full")
			}
			s, u, v := extractSVD(&svd)
			sigma := NewDense(m, n, nil)
			for i := 0; i < min(m, n); i++ {
				sigma.Set(i, i, s[i])
			}
			var ansFull Dense
			ansFull.Product(u, sigma, v.T())
			if !EqualApprox(&ansFull, a, 1e-8) {
				t.Errorf("Answer mismatch when SVDFull")
			}

			// Test Thin decomposition.
			ok = svd.Factorize(a, SVDThin)
			if !ok {
				t.Errorf("SVD factorization failed")
			}
			if !Equal(a, aCopy) {
				t.Errorf("A changed during call to SVD with Thin")
			}
			sThin, u, v := extractSVD(&svd)
			if !floats.EqualApprox(s, sThin, 1e-8) {
				t.Errorf("Singular value mismatch between Full and Thin decomposition")
			}
			sigma = NewDense(min(m, n), min(m, n), nil)
			for i := 0; i < min(m, n); i++ {
				sigma.Set(i, i, sThin[i])
			}
			ansFull.Reset()
			ansFull.Product(u, sigma, v.T())
			if !EqualApprox(&ansFull, a, 1e-8) {
				t.Errorf("Answer mismatch when SVDFull")
			}

			// Test None decomposition.
			ok = svd.Factorize(a, SVDNone)
			if !ok {
				t.Errorf("SVD factorization failed")
			}
			if !Equal(a, aCopy) {
				t.Errorf("A changed during call to SVD with none")
			}
			sNone := make([]float64, min(m, n))
			svd.Values(sNone)
			if !floats.EqualApprox(s, sNone, 1e-8) {
				t.Errorf("Singular value mismatch between Full and None decomposition")
			}
		}
	}
}

func extractSVD(svd *SVD) (s []float64, u, v *Dense) {
	u = &Dense{}
	svd.UTo(u)
	v = &Dense{}
	svd.VTo(v)
	return svd.Values(nil), u, v
}

func TestSVDSolveTo(t *testing.T) {
	t.Parallel()
	rnd := rand.New(rand.NewSource(1))
	// Hand-coded cases.
	for i, test := range []struct {
		a      []float64
		m, n   int
		b      []float64
		bc     int
		rcond  float64
		want   []float64
		wm, wn int
	}{
		{
			a: []float64{6}, m: 1, n: 1,
			b: []float64{3}, bc: 1,
			want: []float64{0.5}, wm: 1, wn: 1,
		},
		{
			a: []float64{
				1, 0, 0,
				0, 1, 0,
				0, 0, 1,
			}, m: 3, n: 3,
			b: []float64{
				3,
				2,
				1,
			}, bc: 1,
			want: []float64{
				3,
				2,
				1,
			}, wm: 3, wn: 1,
		},
		{
			a: []float64{
				0.8147, 0.9134, 0.5528,
				0.9058, 0.6324, 0.8723,
				0.1270, 0.0975, 0.7612,
			}, m: 3, n: 3,
			b: []float64{
				0.278,
				0.547,
				0.958,
			}, bc: 1,
			want: []float64{
				-0.932687281002860,
				0.303963920182067,
				1.375216503507109,
			}, wm: 3, wn: 1,
		},
		{
			a: []float64{
				0.8147, 0.9134, 0.5528,
				0.9058, 0.6324, 0.8723,
			}, m: 2, n: 3,
			b: []float64{
				0.278,
				0.547,
			}, bc: 1,
			want: []float64{
				0.25919787248965376,
				-0.25560256266441034,
				0.5432324059702451,
			}, wm: 3, wn: 1,
		},
		{
			a: []float64{
				0.8147, 0.9134, 0.9,
				0.9058, 0.6324, 0.9,
				0.1270, 0.0975, 0.1,
				1.6, 2.8, -3.5,
			}, m: 4, n: 3,
			b: []float64{
				0.278,
				0.547,
				-0.958,
				1.452,
			}, bc: 1,
			want: []float64{
				0.820970340787782,
				-0.218604626527306,
				-0.212938815234215,
			}, wm: 3, wn: 1,
		},
		{
			a: []float64{
				0.8147, 0.9134, 0.231, -1.65,
				0.9058, 0.6324, 0.9, 0.72,
				0.1270, 0.0975, 0.1, 1.723,
				1.6, 2.8, -3.5, 0.987,
				7.231, 9.154, 1.823, 0.9,
			}, m: 5, n: 4,
			b: []float64{
				0.278, 8.635,
				0.547, 9.125,
				-0.958, -0.762,
				1.452, 1.444,
				1.999, -7.234,
			}, bc: 2,
			want: []float64{
				1.863006789511373, 44.467887791812750,
				-1.127270935407224, -34.073794226035126,
				-0.527926457947330, -8.032133759788573,
				-0.248621916204897, -2.366366415805275,
			}, wm: 4, wn: 2,
		},
		{
			// Test rank-deficient case compared with numpy.
			// >>> import numpy as np
			// >>> b = np.array([[-2.3181340317357653],
			// ...     [-0.7146777651358073],
			// ...     [1.8361340927945298],
			// ...     [-0.35699930593018775],
			// ...     [-1.6359508076249094]])
			// >>> A = np.array([[-1.7854591879711257, -0.42687285925779594, -0.12730256811265162],
			// ...     [-0.5728984211439724, -0.10093393134001777, -0.1181901192353067],
			// ...     [1.2484316018707418, 0.5646683943038734, -0.48229492403243485],
			// ...     [0.10174927665169475, -0.5805410929482445, 1.3054473231942054],
			// ...     [-1.134174808195733, -0.4732430202414438, 0.3528489486370508]])
			// >>> np.linalg.lstsq(A, b, rcond=None)
			// (array([[ 1.21208422],
			//        [ 0.41541503],
			//        [-0.18320349]]), array([], dtype=float64), 2, array([2.68451480e+00, 1.52593185e+00, 6.82840229e-17]))

			a: []float64{
				-1.7854591879711257, -0.42687285925779594, -0.12730256811265162,
				-0.5728984211439724, -0.10093393134001777, -0.1181901192353067,
				1.2484316018707418, 0.5646683943038734, -0.48229492403243485,
				0.10174927665169475, -0.5805410929482445, 1.3054473231942054,
				-1.134174808195733, -0.4732430202414438, 0.3528489486370508,
			}, m: 5, n: 3,
			b: []float64{
				-2.3181340317357653,
				-0.7146777651358073,
				1.8361340927945298,
				-0.35699930593018775,
				-1.6359508076249094,
			}, bc: 1,
			rcond: 1e-15,
			want: []float64{
				1.2120842180372118,
				0.4154150318658529,
				-0.1832034870198265,
			}, wm: 3, wn: 1,
		},
		{
			a: []float64{
				0, 0,
				0, 0,
			}, m: 2, n: 2,
			b: []float64{
				3,
				2,
			}, bc: 1,
		},
		{
			a: []float64{
				0, 0,
				0, 0,
				0, 0,
			}, m: 3, n: 2,
			b: []float64{
				3,
				2,
				1,
			}, bc: 1,
		},
		{
			a: []float64{
				0, 0, 0,
				0, 0, 0,
			}, m: 2, n: 3,
			b: []float64{
				3,
				2,
			}, bc: 1,
		},
	} {
		a := NewDense(test.m, test.n, test.a)
		b := NewDense(test.m, test.bc, test.b)

		var want *Dense
		if test.want != nil {
			want = NewDense(test.wm, test.wn, test.want)
		}

		var svd SVD
		ok := svd.Factorize(a, SVDFull)
		if !ok {
			t.Errorf("unexpected factorization failure for test %d", i)
			continue
		}

		var x Dense
		rank := svd.Rank(test.rcond)
		if rank == 0 {
			continue
		}
		svd.SolveTo(&x, b, rank)
		if !EqualApprox(&x, want, 1e-12) {
			t.Errorf("Solve answer mismatch. Want %v, got %v", want, x)
		}
	}

	// Random Cases.
	for i, test := range []struct {
		m, n, bc int
		rcond    float64
	}{
		{m: 5, n: 5, bc: 1},
		{m: 5, n: 10, bc: 1},
		{m: 10, n: 5, bc: 1},
		{m: 5, n: 5, bc: 7},
		{m: 5, n: 10, bc: 7},
		{m: 10, n: 5, bc: 7},
		{m: 5, n: 5, bc: 12},
		{m: 5, n: 10, bc: 12},
		{m: 10, n: 5, bc: 12},
	} {
		m := test.m
		n := test.n
		bc := test.bc
		a := NewDense(m, n, nil)
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				a.Set(i, j, rnd.Float64())
			}
		}
		br := m
		b := NewDense(br, bc, nil)
		for i := 0; i < br; i++ {
			for j := 0; j < bc; j++ {
				b.Set(i, j, rnd.Float64())
			}
		}

		var svd SVD
		ok := svd.Factorize(a, SVDFull)
		if !ok {
			t.Errorf("unexpected factorization failure for test %d", i)
			continue
		}

		var x Dense
		rank := svd.Rank(test.rcond)
		if rank == 0 {
			continue
		}
		svd.SolveTo(&x, b, rank)

		// Test that the normal equations hold.
		// Aᵀ * A * x = Aᵀ * b
		var tmp, lhs, rhs Dense
		tmp.Mul(a.T(), a)
		lhs.Mul(&tmp, &x)
		rhs.Mul(a.T(), b)
		if !EqualApprox(&lhs, &rhs, 1e-10) {
			t.Errorf("Normal equations do not hold.\nLHS: %v\n, RHS: %v\n", lhs, rhs)
		}
	}
}

func TestSVDSolveVecTo(t *testing.T) {
	t.Parallel()
	rnd := rand.New(rand.NewSource(1))
	// Hand-coded cases.
	for i, test := range []struct {
		a     []float64
		m, n  int
		b     []float64
		rcond float64
		want  []float64
	}{
		{
			a: []float64{6}, m: 1, n: 1,
			b:    []float64{3},
			want: []float64{0.5},
		},
		{
			a: []float64{
				1, 0, 0,
				0, 1, 0,
				0, 0, 1,
			}, m: 3, n: 3,
			b:    []float64{3, 2, 1},
			want: []float64{3, 2, 1},
		},
		{
			a: []float64{
				0.8147, 0.9134, 0.5528,
				0.9058, 0.6324, 0.8723,
				0.1270, 0.0975, 0.7612,
			}, m: 3, n: 3,
			b:    []float64{0.278, 0.547, 0.958},
			want: []float64{-0.932687281002860, 0.303963920182067, 1.375216503507109},
		},
		{
			a: []float64{
				0.8147, 0.9134, 0.5528,
				0.9058, 0.6324, 0.8723,
			}, m: 2, n: 3,
			b:    []float64{0.278, 0.547},
			want: []float64{0.25919787248965376, -0.25560256266441034, 0.5432324059702451},
		},
		{
			a: []float64{
				0.8147, 0.9134, 0.9,
				0.9058, 0.6324, 0.9,
				0.1270, 0.0975, 0.1,
				1.6, 2.8, -3.5,
			}, m: 4, n: 3,
			b:    []float64{0.278, 0.547, -0.958, 1.452},
			want: []float64{0.820970340787782, -0.218604626527306, -0.212938815234215},
		},
		{
			// Test rank-deficient case compared with numpy.
			// >>> import numpy as np
			// >>> b = np.array([[-2.3181340317357653],
			// ...     [-0.7146777651358073],
			// ...     [1.8361340927945298],
			// ...     [-0.35699930593018775],
			// ...     [-1.6359508076249094]])
			// >>> A = np.array([[-1.7854591879711257, -0.42687285925779594, -0.12730256811265162],
			// ...     [-0.5728984211439724, -0.10093393134001777, -0.1181901192353067],
			// ...     [1.2484316018707418, 0.5646683943038734, -0.48229492403243485],
			// ...     [0.10174927665169475, -0.5805410929482445, 1.3054473231942054],
			// ...     [-1.134174808195733, -0.4732430202414438, 0.3528489486370508]])
			// >>> np.linalg.lstsq(A, b, rcond=None)
			// (array([[ 1.21208422],
			//        [ 0.41541503],
			//        [-0.18320349]]), array([], dtype=float64), 2, array([2.68451480e+00, 1.52593185e+00, 6.82840229e-17]))

			a: []float64{
				-1.7854591879711257, -0.42687285925779594, -0.12730256811265162,
				-0.5728984211439724, -0.10093393134001777, -0.1181901192353067,
				1.2484316018707418, 0.5646683943038734, -0.48229492403243485,
				0.10174927665169475, -0.5805410929482445, 1.3054473231942054,
				-1.134174808195733, -0.4732430202414438, 0.3528489486370508,
			}, m: 5, n: 3,
			b:     []float64{-2.3181340317357653, -0.7146777651358073, 1.8361340927945298, -0.35699930593018775, -1.6359508076249094},
			rcond: 1e-15,
			want:  []float64{1.2120842180372118, 0.4154150318658529, -0.1832034870198265},
		},
		{
			a: []float64{
				0, 0,
				0, 0,
			}, m: 2, n: 2,
			b: []float64{3, 2},
		},
		{
			a: []float64{
				0, 0,
				0, 0,
				0, 0,
			}, m: 3, n: 2,
			b: []float64{3, 2, 1},
		},
		{
			a: []float64{
				0, 0, 0,
				0, 0, 0,
			}, m: 2, n: 3,
			b: []float64{3, 2},
		},
	} {
		a := NewDense(test.m, test.n, test.a)
		b := NewVecDense(len(test.b), test.b)

		var want *VecDense
		if test.want != nil {
			want = NewVecDense(len(test.want), test.want)
		}

		var svd SVD
		ok := svd.Factorize(a, SVDFull)
		if !ok {
			t.Errorf("unexpected factorization failure for test %d", i)
			continue
		}

		var x VecDense
		rank := svd.Rank(test.rcond)
		if rank == 0 {
			continue
		}
		svd.SolveVecTo(&x, b, rank)
		if !EqualApprox(&x, want, 1e-12) {
			t.Errorf("Solve answer mismatch. Want %v, got %v", want, x)
		}
	}

	// Random Cases.
	for i, test := range []struct {
		m, n  int
		rcond float64
	}{
		{m: 5, n: 5},
		{m: 5, n: 10},
		{m: 10, n: 5},
		{m: 5, n: 5},
		{m: 5, n: 10},
		{m: 10, n: 5},
		{m: 5, n: 5},
		{m: 5, n: 10},
		{m: 10, n: 5},
	} {
		m := test.m
		n := test.n
		a := NewDense(m, n, nil)
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				a.Set(i, j, rnd.Float64())
			}
		}
		br := m
		b := NewVecDense(br, nil)
		for i := 0; i < br; i++ {
			b.SetVec(i, rnd.Float64())
		}

		var svd SVD
		ok := svd.Factorize(a, SVDFull)
		if !ok {
			t.Errorf("unexpected factorization failure for test %d", i)
			continue
		}

		var x VecDense
		rank := svd.Rank(test.rcond)
		if rank == 0 {
			continue
		}
		svd.SolveVecTo(&x, b, rank)

		// Test that the normal equations hold.
		// Aᵀ * A * x = Aᵀ * b
		var tmp, lhs, rhs Dense
		tmp.Mul(a.T(), a)
		lhs.Mul(&tmp, &x)
		rhs.Mul(a.T(), b)
		if !EqualApprox(&lhs, &rhs, 1e-10) {
			t.Errorf("Normal equations do not hold.\nLHS: %v\n, RHS: %v\n", lhs, rhs)
		}
	}
}
