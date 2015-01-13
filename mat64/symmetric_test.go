package mat64

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"gopkg.in/check.v1"
)

func (s *S) TestNewSymmetric(c *check.C) {
	for i, test := range []struct {
		data []float64
		N    int
		mat  *Symmetric
	}{
		{
			data: []float64{
				1, 2, 3,
				4, 5, 6,
				7, 8, 9,
			},
			N: 3,
			mat: &Symmetric{blas64.Symmetric{
				N:      3,
				Stride: 3,
				Uplo:   blas.Upper,
				Data:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			}},
		},
	} {
		t := NewSymmetric(test.N, test.data)
		rows, cols := t.Dims()
		c.Check(rows, check.Equals, test.N, check.Commentf("Test %d", i))
		c.Check(cols, check.Equals, test.N, check.Commentf("Test %d", i))
		c.Check(t, check.DeepEquals, test.mat, check.Commentf("Test %d", i))

		m := NewDense(test.N, test.N, test.data)
		c.Check(t.mat.Data, check.DeepEquals, m.mat.Data, check.Commentf("Test %d", i))

		c.Check(func() { NewSymmetric(3, []float64{1, 2}) }, check.PanicMatches, ErrShape.Error())
	}
}

func (s *S) TestTriAtSet(c *check.C) {
	t := &Symmetric{blas64.Symmetric{
		N:      3,
		Stride: 3,
		Uplo:   blas.Upper,
		Data:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
	}}
	rows, cols := t.Dims()
	// Check At out of bounds
	c.Check(func() { t.At(rows, 0) }, check.PanicMatches, ErrRowAccess.Error(), check.Commentf("Test row out of bounds"))
	c.Check(func() { t.At(0, cols) }, check.PanicMatches, ErrColAccess.Error(), check.Commentf("Test col out of bounds"))
	c.Check(func() { t.At(rows+1, 0) }, check.PanicMatches, ErrRowAccess.Error(), check.Commentf("Test row out of bounds"))
	c.Check(func() { t.At(0, cols+1) }, check.PanicMatches, ErrColAccess.Error(), check.Commentf("Test col out of bounds"))
	c.Check(func() { t.At(-1, 0) }, check.PanicMatches, ErrRowAccess.Error(), check.Commentf("Test row out of bounds"))
	c.Check(func() { t.At(0, -1) }, check.PanicMatches, ErrColAccess.Error(), check.Commentf("Test col out of bounds"))

	// Check Set out of bounds
	c.Check(func() { t.SetSym(rows, 0, 1.2) }, check.PanicMatches, ErrRowAccess.Error(), check.Commentf("Test row out of bounds"))
	c.Check(func() { t.SetSym(0, cols, 1.2) }, check.PanicMatches, ErrColAccess.Error(), check.Commentf("Test col out of bounds"))
	c.Check(func() { t.SetSym(rows+1, 0, 1.2) }, check.PanicMatches, ErrRowAccess.Error(), check.Commentf("Test row out of bounds"))
	c.Check(func() { t.SetSym(0, cols+1, 1.2) }, check.PanicMatches, ErrColAccess.Error(), check.Commentf("Test col out of bounds"))
	c.Check(func() { t.SetSym(-1, 0, 1.2) }, check.PanicMatches, ErrRowAccess.Error(), check.Commentf("Test row out of bounds"))
	c.Check(func() { t.SetSym(0, -1, 1.2) }, check.PanicMatches, ErrColAccess.Error(), check.Commentf("Test col out of bounds"))

	c.Check(t.At(2, 1), check.Equals, 6.0)
	c.Check(t.At(1, 2), check.Equals, 6.0)
	t.SetSym(1, 2, 15)
	c.Check(t.At(2, 1), check.Equals, 15.0)
	c.Check(t.At(1, 2), check.Equals, 15.0)
	t.SetSym(2, 1, 12)
	c.Check(t.At(2, 1), check.Equals, 12.0)
	c.Check(t.At(1, 2), check.Equals, 12.0)
}
