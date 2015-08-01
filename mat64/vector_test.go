package mat64

import (
	"github.com/gonum/blas/blas64"
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

		c.Check(func() { v.SetVec(n, 100) }, check.PanicMatches, ErrVectorAccess.Error(), check.Commentf("Test %d", i))
		c.Check(func() { v.SetVec(-1, 100) }, check.PanicMatches, ErrVectorAccess.Error(), check.Commentf("Test %d", i))

		v.SetVec(0, 100)
		c.Check(v.At(0, 0), check.Equals, 100.0, check.Commentf("Test %d", i))
		v.SetVec(2, 101)
		c.Check(v.At(2, 0), check.Equals, 101.0, check.Commentf("Test %d", i))
	}
}

func (s *S) TestVectorMul(c *check.C) {
	method := func(receiver, a, b Matrix) {
		type mulVecer interface {
			MulVec(a Matrix, b *Vector)
		}
		rd := receiver.(mulVecer)
		rd.MulVec(a, b.(*Vector))
	}
	denseComparison := func(receiver, a, b *Dense) {
		receiver.Mul(a, b)
	}
	legalSizeMulVec := func(ar, ac, br, bc int) bool {
		var legal bool
		if bc != 1 {
			legal = false
		} else {
			legal = ac == br
		}
		return legal
	}
	testTwoInput(c, "MulVec", &Vector{}, method, denseComparison, legalTypesNotVecVec, legalSizeMulVec)
}

func (s *S) TestVectorAdd(c *check.C) {
	for i, test := range []struct {
		a, b *Vector
		want *Vector
	}{
		{
			a:    NewVector(3, []float64{0, 1, 2}),
			b:    NewVector(3, []float64{0, 2, 3}),
			want: NewVector(3, []float64{0, 3, 5}),
		},
		{
			a:    NewVector(3, []float64{0, 1, 2}),
			b:    NewDense(3, 1, []float64{0, 2, 3}).ColView(0),
			want: NewVector(3, []float64{0, 3, 5}),
		},
		{
			a:    NewDense(3, 1, []float64{0, 1, 2}).ColView(0),
			b:    NewDense(3, 1, []float64{0, 2, 3}).ColView(0),
			want: NewVector(3, []float64{0, 3, 5}),
		},
	} {
		var v Vector
		v.AddVec(test.a, test.b)
		c.Check(v.RawVector(), check.DeepEquals, test.want.RawVector(), check.Commentf("Test %d", i))
	}
}

func (s *S) TestVectorSub(c *check.C) {
	for i, test := range []struct {
		a, b *Vector
		want *Vector
	}{
		{
			a:    NewVector(3, []float64{0, 1, 2}),
			b:    NewVector(3, []float64{0, 0.5, 1}),
			want: NewVector(3, []float64{0, 0.5, 1}),
		},
		{
			a:    NewVector(3, []float64{0, 1, 2}),
			b:    NewDense(3, 1, []float64{0, 0.5, 1}).ColView(0),
			want: NewVector(3, []float64{0, 0.5, 1}),
		},
		{
			a:    NewDense(3, 1, []float64{0, 1, 2}).ColView(0),
			b:    NewDense(3, 1, []float64{0, 0.5, 1}).ColView(0),
			want: NewVector(3, []float64{0, 0.5, 1}),
		},
	} {
		var v Vector
		v.SubVec(test.a, test.b)
		c.Check(v.RawVector(), check.DeepEquals, test.want.RawVector(), check.Commentf("Test %d", i))
	}
}

func (s *S) TestVectorMulElem(c *check.C) {
	for i, test := range []struct {
		a, b *Vector
		want *Vector
	}{
		{
			a:    NewVector(3, []float64{0, 1, 2}),
			b:    NewVector(3, []float64{0, 2, 3}),
			want: NewVector(3, []float64{0, 2, 6}),
		},
		{
			a:    NewVector(3, []float64{0, 1, 2}),
			b:    NewDense(3, 1, []float64{0, 2, 3}).ColView(0),
			want: NewVector(3, []float64{0, 2, 6}),
		},
		{
			a:    NewDense(3, 1, []float64{0, 1, 2}).ColView(0),
			b:    NewDense(3, 1, []float64{0, 2, 3}).ColView(0),
			want: NewVector(3, []float64{0, 2, 6}),
		},
	} {
		var v Vector
		v.MulElemVec(test.a, test.b)
		c.Check(v.RawVector(), check.DeepEquals, test.want.RawVector(), check.Commentf("Test %d", i))
	}
}

func (s *S) TestVectorDivElem(c *check.C) {
	for i, test := range []struct {
		a, b *Vector
		want *Vector
	}{
		{
			a:    NewVector(3, []float64{0.5, 1, 2}),
			b:    NewVector(3, []float64{0.5, 0.5, 1}),
			want: NewVector(3, []float64{1, 2, 2}),
		},
		{
			a:    NewVector(3, []float64{0.5, 1, 2}),
			b:    NewDense(3, 1, []float64{0.5, 0.5, 1}).ColView(0),
			want: NewVector(3, []float64{1, 2, 2}),
		},
		{
			a:    NewDense(3, 1, []float64{0.5, 1, 2}).ColView(0),
			b:    NewDense(3, 1, []float64{0.5, 0.5, 1}).ColView(0),
			want: NewVector(3, []float64{1, 2, 2}),
		},
	} {
		var v Vector
		v.DivElemVec(test.a, test.b)
		c.Check(v.RawVector(), check.DeepEquals, test.want.RawVector(), check.Commentf("Test %d", i))
	}
}
