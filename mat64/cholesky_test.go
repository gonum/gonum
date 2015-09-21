// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"math"
	"math/rand"
	"testing"

	"gopkg.in/check.v1"
)

func (s *S) TestCholesky(c *check.C) {
	for _, t := range []struct {
		a *SymDense

		cond   float64
		want   *TriDense
		posdef bool
	}{
		{
			a: NewSymDense(3, []float64{
				4, 1, 1,
				0, 2, 3,
				0, 0, 6,
			}),
			cond: 37,
			want: NewTriDense(3, true, []float64{
				2, 0.5, 0.5,
				0, 1.3228756555322954, 2.0788046015507495,
				0, 0, 1.195228609334394,
			}),
			posdef: true,
		},
	} {
		_, n := t.a.Dims()
		// Try with a new cholesky struct
		var chol Cholesky
		ok := chol.Factorize(t.a)
		c.Check(ok, check.Equals, t.posdef)
		fc := DenseCopyOf(chol.chol)
		c.Check(Equal(fc, t.want), check.Equals, true)
		if math.Abs(t.cond-chol.cond) > 1e-13 {
			c.Errorf("Condition number mismatch: Want %v, got %v", t.cond, chol.cond)
		}
		var U TriDense
		U.UFromCholesky(&chol)
		aCopy := DenseCopyOf(t.a)
		var a Dense
		a.Mul(U.TTri(), &U)
		c.Check(EqualApprox(&a, aCopy, 1e-14), check.Equals, true)

		var L TriDense
		L.LFromCholesky(&chol)
		a.Mul(&L, L.TTri())
		c.Check(EqualApprox(&a, aCopy, 1e-14), check.Equals, true)

		// Try with a cholesky struct that is too small
		cholSmall := &Cholesky{
			chol: NewTriDense(n-1, true, nil),
		}
		for i := range cholSmall.chol.mat.Data {
			cholSmall.chol.mat.Data[i] = rand.Float64()
		}
		ok = cholSmall.Factorize(t.a)
		c.Check(ok, check.Equals, t.posdef)
		c.Check(Equal(fc, t.want), check.Equals, true)
		if math.Abs(t.cond-cholSmall.cond) > 1e-13 {
			c.Errorf("Condition number mismatch: Want %v, got %v", t.cond, chol.cond)
		}

		// Try with a cholesky struct that is the right size.
		cholCorrect := &Cholesky{
			chol: NewTriDense(n, true, nil),
		}
		for i := range cholCorrect.chol.mat.Data {
			cholCorrect.chol.mat.Data[i] = rand.Float64()
		}
		ok = cholCorrect.Factorize(t.a)
		c.Check(ok, check.Equals, t.posdef)
		c.Check(Equal(fc, t.want), check.Equals, true)
		if math.Abs(t.cond-cholCorrect.cond) > 1e-13 {
			c.Errorf("Condition number mismatch: Want %v, got %v", t.cond, chol.cond)
		}

		// Try with a cholesky struct that is too large
		cholLarge := &Cholesky{
			chol: NewTriDense(n+1, true, nil),
		}
		for i := range cholLarge.chol.mat.Data {
			cholLarge.chol.mat.Data[i] = rand.Float64()
		}
		ok = cholLarge.Factorize(t.a)
		c.Check(ok, check.Equals, t.posdef)
		c.Check(Equal(fc, t.want), check.Equals, true)
		if math.Abs(t.cond-cholLarge.cond) > 1e-13 {
			c.Errorf("Condition number mismatch: Want %v, got %v", t.cond, chol.cond)
		}
	}
}

func (s *S) TestCholeskySolve(c *check.C) {
	for _, t := range []struct {
		a   *SymDense
		b   *Dense
		ans *Dense
	}{
		{
			a: NewSymDense(2, []float64{
				1, 0,
				0, 1,
			}),
			b:   NewDense(2, 1, []float64{5, 6}),
			ans: NewDense(2, 1, []float64{5, 6}),
		},
		{
			a: NewSymDense(3, []float64{
				53, 59, 37,
				0, 83, 71,
				37, 71, 101,
			}),
			b:   NewDense(3, 1, []float64{5, 6, 7}),
			ans: NewDense(3, 1, []float64{0.20745069393718094, -0.17421475529583694, 0.11577794010226464}),
		},
	} {
		var chol Cholesky
		ok := chol.Factorize(t.a)
		c.Assert(ok, check.Equals, true)

		var x Dense
		x.SolveCholesky(&chol, t.b)
		c.Check(EqualApprox(&x, t.ans, 1e-12), check.Equals, true)

		var ans Dense
		ans.Mul(t.a, &x)
		c.Check(EqualApprox(&ans, t.b, 1e-12), check.Equals, true)
	}
}

func (s *S) TestCholeskySolveVec(c *check.C) {
	for _, t := range []struct {
		a   *SymDense
		b   *Vector
		ans *Vector
	}{
		{
			a: NewSymDense(2, []float64{
				1, 0,
				0, 1,
			}),
			b:   NewVector(2, []float64{5, 6}),
			ans: NewVector(2, []float64{5, 6}),
		},
		{
			a: NewSymDense(3, []float64{
				53, 59, 37,
				0, 83, 71,
				0, 0, 101,
			}),
			b:   NewVector(3, []float64{5, 6, 7}),
			ans: NewVector(3, []float64{0.20745069393718094, -0.17421475529583694, 0.11577794010226464}),
		},
	} {
		var chol Cholesky
		ok := chol.Factorize(t.a)
		c.Assert(ok, check.Equals, true)

		var x Vector
		x.SolveCholeskyVec(&chol, t.b)
		c.Check(EqualApprox(&x, t.ans, 1e-12), check.Equals, true)

		var ans Vector
		ans.MulVec(t.a, &x)
		c.Check(EqualApprox(&ans, t.b, 1e-12), check.Equals, true)
	}
}

func BenchmarkCholeskySmall(b *testing.B) {
	benchmarkCholesky(b, 2)
}

func BenchmarkCholeskyMedium(b *testing.B) {
	benchmarkCholesky(b, Med)
}

func BenchmarkCholeskyLarge(b *testing.B) {
	benchmarkCholesky(b, Lg)
}

func benchmarkCholesky(b *testing.B, n int) {
	base := make([]float64, n*n)
	for i := range base {
		base[i] = rand.Float64()
	}
	bm := NewDense(n, n, base)
	bm.Mul(bm.T(), bm)
	am := NewSymDense(n, bm.mat.Data)

	var chol Cholesky
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ok := chol.Factorize(am)
		if !ok {
			panic("not pos def")
		}
	}
}
