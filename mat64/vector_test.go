package mat64

import (
	"math/rand"

	"github.com/gonum/blas/blas64"
	"github.com/gonum/floats"
	"gopkg.in/check.v1"
)

func (s *S) TestNewVector(c *check.C) {
	for i, test := range []struct {
		n      int
		data   []float64
		vector *Vector
	}{
		{
			n:    3,
			data: []float64{4, 5, 6},
			vector: &Vector{
				mat: blas64.Vector{
					Data: []float64{4, 5, 6},
					Inc:  1,
				},
				n: 3,
			},
		},
	} {
		v := NewVector(test.n, test.data)
		rows, cols := v.Dims()
		c.Check(rows, check.Equals, test.n, check.Commentf("Test %d", i))
		c.Check(cols, check.Equals, 1, check.Commentf("Test %d", i))
		c.Check(v, check.DeepEquals, test.vector, check.Commentf("Test %d", i))
		v2 := NewVector(test.n, nil)
		c.Check(v2.mat.Data, check.DeepEquals, []float64{0, 0, 0}, check.Commentf("Test %d", i))
	}
}

func (s *S) TestVectorAtSet(c *check.C) {
	for i, test := range []struct {
		vector *Vector
	}{
		{
			vector: &Vector{
				mat: blas64.Vector{
					Data: []float64{0, 1, 2},
					Inc:  1,
				},
				n: 3,
			},
		},
		{
			vector: &Vector{
				mat: blas64.Vector{
					Data: []float64{0, 10, 10, 1, 10, 10, 2},
					Inc:  3,
				},
				n: 3,
			},
		},
	} {
		v := test.vector
		n := test.vector.n
		c.Check(func() { v.At(n, 0) }, check.PanicMatches, ErrRowAccess.Error(), check.Commentf("Test %d", i))
		c.Check(func() { v.At(-1, 0) }, check.PanicMatches, ErrRowAccess.Error(), check.Commentf("Test %d", i))
		c.Check(func() { v.At(0, 1) }, check.PanicMatches, ErrColAccess.Error(), check.Commentf("Test %d", i))
		c.Check(func() { v.At(0, -1) }, check.PanicMatches, ErrColAccess.Error(), check.Commentf("Test %d", i))

		c.Check(v.At(0, 0), check.Equals, 0.0, check.Commentf("Test %d", i))
		c.Check(v.At(1, 0), check.Equals, 1.0, check.Commentf("Test %d", i))
		c.Check(v.At(n-1, 0), check.Equals, float64(n-1), check.Commentf("Test %d", i))

		c.Check(func() { v.Set(n, 0, 100) }, check.PanicMatches, ErrRowAccess.Error(), check.Commentf("Test %d", i))
		c.Check(func() { v.Set(-1, 0, 100) }, check.PanicMatches, ErrRowAccess.Error(), check.Commentf("Test %d", i))
		c.Check(func() { v.Set(0, 1, 100) }, check.PanicMatches, ErrColAccess.Error(), check.Commentf("Test %d", i))
		c.Check(func() { v.Set(0, -1, 100) }, check.PanicMatches, ErrColAccess.Error(), check.Commentf("Test %d", i))

		v.Set(0, 0, 100)
		c.Check(v.At(0, 0), check.Equals, 100.0, check.Commentf("Test %d", i))
		v.Set(2, 0, 101)
		c.Check(v.At(2, 0), check.Equals, 101.0, check.Commentf("Test %d", i))
	}
}

func (s *S) TestVectorMul(c *check.C) {

	for i, test := range []struct {
		m int
		n int
	}{
		{
			m: 10,
			n: 5,
		},
		{
			m: 5,
			n: 5,
		},
		{
			m: 5,
			n: 10,
		},
	} {
		vData := make([]float64, test.n)
		for i := range vData {
			vData[i] = rand.Float64()
		}
		vDataCopy := make([]float64, test.n)
		copy(vDataCopy, vData)
		v := NewVector(test.n, vData)
		aData := make([]float64, test.n*test.m)
		for i := range aData {
			aData[i] = rand.Float64()
		}
		a := NewDense(test.m, test.n, aData)
		var v2 Vector
		v2.MulVec(a, false, v)
		var v2M Dense
		v2M.Mul(a, v)
		same := floats.EqualApprox(v2.mat.Data, v2M.mat.Data, 1e-14)
		c.Check(same, check.Equals, true, check.Commentf("Test %d", i))

		var aT Dense
		aT.TCopy(a)
		v2.MulVec(&aT, true, v)
		same = floats.EqualApprox(v2.mat.Data, v2M.mat.Data, 1e-14)
		c.Check(same, check.Equals, true, check.Commentf("Test %d", i))

		/*
			v.MulVec(&aT, true, v)
			same = floats.EqualApprox(v.mat.Data, v2M.mat.Data, 1e-14)
			c.Check(same, check.Equals, true, check.Commentf("Test %d", i))
		*/
	}
}
