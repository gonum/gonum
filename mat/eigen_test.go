// Copyright ©2013 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"math"
	"math/rand/v2"
	"sort"
	"testing"

	"gonum.org/v1/gonum/floats"
)

func TestEigen(t *testing.T) {
	t.Parallel()
	for i, test := range []struct {
		a *Dense

		values []complex128
		left   *CDense
		right  *CDense
	}{
		{
			a: NewDense(3, 3, []float64{
				1, 0, 0,
				0, 1, 0,
				0, 0, 1,
			}),
			values: []complex128{1, 1, 1},
			left: NewCDense(3, 3, []complex128{
				1, 0, 0,
				0, 1, 0,
				0, 0, 1,
			}),
			right: NewCDense(3, 3, []complex128{
				1, 0, 0,
				0, 1, 0,
				0, 0, 1,
			}),
		},
		{
			// Values compared with numpy.
			a: NewDense(4, 4, []float64{
				0.9025, 0.025, 0.475, 0.0475,
				0.0475, 0.475, 0.475, 0.0025,
				0.0475, 0.025, 0.025, 0.9025,
				0.0025, 0.475, 0.025, 0.0475,
			}),
			values: []complex128{1, 0.7300317046114154, -0.1400158523057075 + 0.452854925738716i, -0.1400158523057075 - 0.452854925738716i},
			left: NewCDense(4, 4, []complex128{
				0.5, -0.3135167160788313, -0.02058121780136903 + 0.004580939300127051i, -0.02058121780136903 - 0.004580939300127051i,
				0.5, 0.7842199280224781, 0.37551026954193356 - 0.2924634904103879i, 0.37551026954193356 + 0.2924634904103879i,
				0.5, 0.33202200780783525, 0.16052616322784943 + 0.3881393645202527i, 0.16052616322784943 - 0.3881393645202527i,
				0.5, 0.42008065840123954, -0.7723935249234155, -0.7723935249234155,
			}),
			right: NewCDense(4, 4, []complex128{
				0.9476399565969628, -0.8637347682162745, -0.2688989440320280 - 0.1282234938321029i, -0.2688989440320280 + 0.1282234938321029i,
				0.2394935907064427, 0.3457075153704627, -0.3621360383713332 - 0.2583198964498771i, -0.3621360383713332 + 0.2583198964498771i,
				0.1692743801716332, 0.2706851011641580, 0.7426369401030960, 0.7426369401030960,
				0.1263626404003607, 0.2473421516816520, -0.1116019576997347 + 0.3865433902819795i, -0.1116019576997347 - 0.3865433902819795i,
			}),
		},
	} {
		var e1, e2, e3, e4 Eigen
		ok := e1.Factorize(test.a, EigenBoth)
		if !ok {
			panic("bad factorization")
		}
		e2.Factorize(test.a, EigenRight)
		e3.Factorize(test.a, EigenLeft)
		e4.Factorize(test.a, EigenNone)

		v1 := e1.Values(nil)
		if !cmplxEqualTol(v1, test.values, 1e-14) {
			t.Errorf("eigenvalue mismatch. Case %v", i)
		}
		var left CDense
		e1.LeftVectorsTo(&left)
		if !CEqualApprox(&left, test.left, 1e-14) {
			t.Errorf("left eigenvector mismatch. Case %v", i)
		}
		var right CDense
		e1.VectorsTo(&right)
		if !CEqualApprox(&right, test.right, 1e-14) {
			t.Errorf("right eigenvector mismatch. Case %v", i)
		}

		// Check that the eigenvectors and values are the same in all combinations.
		if !cmplxEqual(v1, e2.Values(nil)) {
			t.Errorf("eigenvector mismatch. Case %v", i)
		}
		if !cmplxEqual(v1, e3.Values(nil)) {
			t.Errorf("eigenvector mismatch. Case %v", i)
		}
		if !cmplxEqual(v1, e4.Values(nil)) {
			t.Errorf("eigenvector mismatch. Case %v", i)
		}
		var right2 CDense
		e2.VectorsTo(&right2)
		if !CEqual(&right, &right2) {
			t.Errorf("right eigenvector mismatch. Case %v", i)
		}
		var left3 CDense
		e3.LeftVectorsTo(&left3)
		if !CEqual(&left, &left3) {
			t.Errorf("left eigenvector mismatch. Case %v", i)
		}

		// TODO(btracey): Also add in a test for correctness when #308 is
		// resolved and we have a CMat.Mul().
	}
}

func cmplxEqual(v1, v2 []complex128) bool {
	for i, v := range v1 {
		if v != v2[i] {
			return false
		}
	}
	return true
}

func cmplxEqualTol(v1, v2 []complex128, tol float64) bool {
	for i, v := range v1 {
		if !cEqualWithinAbsOrRel(v, v2[i], tol, tol) {
			return false
		}
	}
	return true
}

func TestEigenSym(t *testing.T) {
	t.Parallel()
	const tol = 1e-14
	// Hand coded tests with results from lapack.
	for cas, test := range []struct {
		mat *SymDense

		values  []float64
		vectors *Dense
	}{
		{
			mat:    NewSymDense(3, []float64{8, 2, 4, 2, 6, 10, 4, 10, 5}),
			values: []float64{-4.707679201365891, 6.294580208480216, 17.413098992885672},
			vectors: NewDense(3, 3, []float64{
				-0.127343483135656, -0.902414161226903, -0.411621572466779,
				-0.664177720955769, 0.385801900032553, -0.640331827193739,
				0.736648893495999, 0.191847792659746, -0.648492738712395,
			}),
		},
	} {
		var es EigenSym
		ok := es.Factorize(test.mat, true)
		if !ok {
			t.Errorf("case %d: bad test", cas)
			continue
		}
		if !floats.EqualApprox(test.values, es.values, tol) {
			t.Errorf("case %d: eigenvalue mismatch", cas)
		}
		if !EqualApprox(test.vectors, es.vectors, tol) {
			t.Errorf("case %d: eigenvector mismatch", cas)
		}

		var es2 EigenSym
		es2.Factorize(test.mat, false)
		if !floats.EqualApprox(es2.values, es.values, tol) {
			t.Errorf("case %d: eigenvalue mismatch when no vectors computed", cas)
		}
	}

	// Randomized tests
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, n := range []int{1, 2, 3, 5, 10, 70} {
		for cas := 0; cas < 10; cas++ {
			a := make([]float64, n*n)
			for i := range a {
				a[i] = rnd.NormFloat64()
			}
			s := NewSymDense(n, a)
			var es EigenSym
			ok := es.Factorize(s, true)
			if !ok {
				t.Errorf("n=%d,cas=%d: bad test", n, cas)
				continue
			}

			// Check that A and EigenSym are equal as Matrix.
			if !EqualApprox(s, &es, tol*float64(n)) {
				t.Errorf("n=%d,cas=%d: A and EigenSym are not equal as Matrix", n, cas)
			}
			if !EqualApprox(s.T(), es.T(), tol*float64(n)) {
				t.Errorf("n=%d,cas=%d: Aᵀ and EigenSymᵀ are not equal as Matrix", n, cas)
			}

			// Check that the eigenvectors are orthonormal.
			if !isOrthonormal(es.vectors, 1e-8) {
				t.Errorf("n=%d,cas=%d: eigenvectors not orthonormal", n, cas)
			}

			// Check that the eigenvalues are actually eigenvalues.
			for i := 0; i < n; i++ {
				v := NewVecDense(n, Col(nil, i, es.vectors))
				var m VecDense
				m.MulVec(s, v)

				var scal VecDense
				scal.ScaleVec(es.values[i], v)

				if !EqualApprox(&m, &scal, 1e-8) {
					t.Errorf("n=%d,cas=%d: eigenvalue %d does not match", n, cas, i)
				}
			}

			// Check that A = Q * D * Qᵀ using the Raw methods.
			var got Dense
			got.Product(es.RawQ(), NewDiagDense(n, es.RawValues()), es.RawQ().T())
			if !EqualApprox(s, &got, tol*float64(n)) {
				var diff Dense
				diff.Sub(s, &got)
				diff.Apply(func(i, j int, v float64) float64 { return math.Abs(diff.At(i, j)) }, &diff)
				t.Errorf("n=%d,cas=%d: A not reconstructed from Q*D*Qᵀ\n|diff|=%v", n, cas,
					Formatted(&diff, Prefix("       ")))
			}

			// Check that the eigenvalues are in ascending order.
			if !sort.Float64sAreSorted(es.values) {
				t.Errorf("n=%d,cas=%d: eigenvalues not ascending", n, cas)
			}
		}
	}
}
