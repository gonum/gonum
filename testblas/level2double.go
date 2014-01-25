package testblas

import (
	"testing"

	"github.com/gonum/blas"
)

// throwPanic will throw unexpected panics if true, or will just report them as errors if false
const throwPanic = true

type DoubleMatTwoVecCase struct {
	Name   string
	m      int
	n      int
	A      [][]float64
	o      blas.Order
	tA     blas.Transpose
	x      []float64
	incX   int
	y      []float64
	incY   int
	lda    int
	xCopy  []float64
	yCopy  []float64
	Panics bool

	DgemvCases []DgemvCase
}

type DgemvCase struct {
	alpha float64
	beta  float64
	ans   []float64
}

var DoubleMatTwoVecCases []DoubleMatTwoVecCase = []DoubleMatTwoVecCase{
	{
		Name: "M_gt_N_Inc1_RowMajor_NoTrans",
		o:    blas.RowMajor,
		tA:   blas.NoTrans,
		m:    5,
		n:    3,
		A: [][]float64{
			{4.1, 6.2, 8.1},
			{9.6, 3.5, 9.1},
			{10, 7, 3},
			{1, 1, 2},
			{9, 2, 5},
		},
		incX:   1,
		incY:   1,
		x:      []float64{1, 2, 3},
		y:      []float64{7, 8, 9, 10, 11},
		lda:    3,
		Panics: false,

		DgemvCases: []DgemvCase{
			{
				alpha: 0,
				beta:  0,
				ans:   []float64{0, 0, 0, 0, 0},
			},
			{
				alpha: 0,
				beta:  1,
				ans:   []float64{7, 8, 9, 10, 11},
			},
			{
				alpha: 1,
				beta:  0,
				ans:   []float64{40.8, 43.9, 33, 9, 28},
			},
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{284.4, 303.2, 210, 12, 158},
			},
		},
	},
	{
		Name: "M_gt_N_Inc1_RowMajor_Trans",
		o:    blas.RowMajor,
		tA:   blas.Trans,
		m:    5,
		n:    3,
		A: [][]float64{
			{4.1, 6.2, 8.1},
			{9.6, 3.5, 9.1},
			{10, 7, 3},
			{1, 1, 2},
			{9, 2, 5},
		},
		incX:   1,
		incY:   1,
		x:      []float64{1, 2, 3, -4, 5},
		y:      []float64{7, 8, 9},
		lda:    3,
		Panics: false,

		DgemvCases: []DgemvCase{
			{
				alpha: 0,
				beta:  0,
				ans:   []float64{0, 0, 0},
			},
			{
				alpha: 0,
				beta:  1,
				ans:   []float64{7, 8, 9},
			},
			{
				alpha: 1,
				beta:  0,
				ans:   []float64{94.3, 40.2, 52.3},
			},
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{712.4, 273.6, 364.4},
			},
		},
	},
	{
		Name: "M_eq_N_Inc1_RowMajor_NoTrans",
		o:    blas.RowMajor,
		tA:   blas.NoTrans,
		m:    3,
		n:    3,
		A: [][]float64{
			{4.1, 6.2, 8.1},
			{9.6, 3.5, 9.1},
			{10, 7, 3},
		},
		incX:   1,
		incY:   1,
		x:      []float64{1, 2, 3},
		y:      []float64{7, 2, 2},
		lda:    3,
		Panics: false,

		DgemvCases: []DgemvCase{
			{
				alpha: 0,
				beta:  0,
				ans:   []float64{0, 0, 0},
			},
			{
				alpha: 0,
				beta:  1,
				ans:   []float64{7, 2, 2},
			},
			{
				alpha: 1,
				beta:  0,
				ans:   []float64{40.8, 43.9, 33},
			},
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{40.8*8 - 6*7, 43.9*8 - 6*2, 33*8 - 6*2},
			},
		},
	},
	{
		Name: "M_eq_N_Inc1_RowMajor_Trans",
		o:    blas.RowMajor,
		tA:   blas.Trans,
		m:    3,
		n:    3,
		A: [][]float64{
			{4.1, 6.2, 8.1},
			{9.6, 3.5, 9.1},
			{10, 7, 3},
		},
		incX:   1,
		incY:   1,
		x:      []float64{1, 2, 3},
		y:      []float64{7, 2, 2},
		lda:    3,
		Panics: false,

		DgemvCases: []DgemvCase{
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{384.4, 261.6, 270.4},
			},
		},
	},
	{
		Name: "M_lt_N_Inc1_RowMajor_NoTrans",
		o:    blas.RowMajor,
		tA:   blas.NoTrans,
		m:    3,
		n:    5,
		A: [][]float64{
			{4.1, 6.2, 8.1, 10, 7},
			{9.6, 3.5, 9.1, -2, 9},
			{10, 7, 3, 1, -5},
		},
		incX:   1,
		incY:   1,
		x:      []float64{1, 2, 3, -7.6, 8.1},
		y:      []float64{7, 2, 2},
		lda:    5,
		Panics: false,

		DgemvCases: []DgemvCase{
			{
				alpha: 0,
				beta:  0,
				ans:   []float64{0, 0, 0},
			},
			{
				alpha: 0,
				beta:  1,
				ans:   []float64{7, 2, 2},
			},
			{
				alpha: 1,
				beta:  0,
				ans:   []float64{21.5, 132, -15.1},
			},

			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{21.5*8 - 6*7, 132*8 - 6*2, -15.1*8 - 6*2},
			},
		},
	},
	{
		Name: "M_lt_N_Inc1_RowMajor_Trans",
		o:    blas.RowMajor,
		tA:   blas.Trans,
		m:    3,
		n:    5,
		A: [][]float64{
			{4.1, 6.2, 8.1, 10, 7},
			{9.6, 3.5, 9.1, -2, 9},
			{10, 7, 3, 1, -5},
		},
		incX:   1,
		incY:   1,
		x:      []float64{1, 2, 3},
		y:      []float64{7, 2, 2, -3, 5},
		lda:    5,
		Panics: false,

		DgemvCases: []DgemvCase{
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{384.4, 261.6, 270.4, 90, 50},
			},
		},
	},
	{
		Name: "M_gt_N_IncNot1_RowMajor_NoTrans",
		o:    blas.RowMajor,
		tA:   blas.NoTrans,
		m:    5,
		n:    3,

		A: [][]float64{
			{4.1, 6.2, 8.1},
			{9.6, 3.5, 9.1},
			{10, 7, 3},
			{1, 1, 2},
			{9, 2, 5},
		},
		incX:   2,
		incY:   3,
		x:      []float64{1, 15, 2, 150, 3},
		y:      []float64{7, 2, 6, 8, -4, -5, 9, 1, 1, 10, 19, 22, 11},
		lda:    3,
		Panics: false,
		DgemvCases: []DgemvCase{
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{284.4, 2, 6, 303.2, -4, -5, 210, 1, 1, 12, 19, 22, 158},
			},
		},
	},
	{
		Name: "M_gt_N_IncNot1_RowMajor_Trans",
		o:    blas.RowMajor,
		tA:   blas.Trans,
		m:    5,
		n:    3,

		A: [][]float64{
			{4.1, 6.2, 8.1},
			{9.6, 3.5, 9.1},
			{10, 7, 3},
			{1, 1, 2},
			{9, 2, 5},
		},
		incX:   2,
		incY:   3,
		x:      []float64{1, 15, 2, 150, 3, 8, -3, 6, 5},
		y:      []float64{7, 2, 6, 8, -4, -5, 9},
		lda:    3,
		Panics: false,
		DgemvCases: []DgemvCase{
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{720.4, 2, 6, 281.6, -4, -5, 380.4},
			},
		},
	},
	{
		Name: "M_eq_N_IncNot1_RowMajor_NoTrans",
		o:    blas.RowMajor,
		tA:   blas.NoTrans,
		m:    3,
		n:    3,
		A: [][]float64{
			{4.1, 6.2, 8.1},
			{9.6, 3.5, 9.1},
			{10, 7, 3},
		},
		incX:   2,
		incY:   3,
		x:      []float64{1, 15, 2, 150, 3},
		y:      []float64{7, 2, 6, 8, -4, -5, 9},
		lda:    3,
		Panics: false,

		DgemvCases: []DgemvCase{
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{284.4, 2, 6, 303.2, -4, -5, 210},
			},
		},
	},
	{
		Name: "M_eq_N_IncNot1_RowMajor_Trans",
		o:    blas.RowMajor,
		tA:   blas.Trans,
		m:    3,
		n:    3,
		A: [][]float64{
			{4.1, 6.2, 8.1},
			{9.6, 3.5, 9.1},
			{10, 7, 3},
		},
		incX:   2,
		incY:   3,
		x:      []float64{1, 15, 2, 150, 3},
		y:      []float64{7, 2, 6, 8, -4, -5, 9},
		lda:    3,
		Panics: false,

		DgemvCases: []DgemvCase{
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{384.4, 2, 6, 225.6, -4, -5, 228.4},
			},
		},
	},
	{
		Name: "M_lt_N_IncNot1_RowMajor_NoTrans",
		o:    blas.RowMajor,
		tA:   blas.NoTrans,
		m:    3,
		n:    5,
		A: [][]float64{
			{4.1, 6.2, 8.1, 10, 11},
			{9.6, 3.5, 9.1, -3, -2},
			{10, 7, 3, -7, -4},
		},
		incX:   2,
		incY:   3,
		x:      []float64{1, 15, 2, 150, 3, -2, -4, 8, -9},
		y:      []float64{7, 2, 6, 8, -4, -5, 9},
		lda:    5,
		Panics: false,

		DgemvCases: []DgemvCase{
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{-827.6, 2, 6, 543.2, -4, -5, 722},
			},
		},
	},
	{
		Name: "M_lt_N_IncNot1_RowMajor_Trans",
		o:    blas.RowMajor,
		tA:   blas.Trans,
		m:    3,
		n:    5,
		A: [][]float64{
			{4.1, 6.2, 8.1, 10, 11},
			{9.6, 3.5, 9.1, -3, -2},
			{10, 7, 3, -7, -4},
		},
		incX:   2,
		incY:   3,
		x:      []float64{1, 15, 2, 150, 3},
		y:      []float64{7, 2, 6, 8, -4, -5, 9, -4, -1, -9, 1, 1, 2},
		lda:    5,
		Panics: false,

		DgemvCases: []DgemvCase{
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{384.4, 2, 6, 225.6, -4, -5, 228.4, -4, -1, -82, 1, 1, -52},
			},
		},
	},

	// TODO: A can be longer than mxn. Add cases where it is longer
	// TODO: x and y can also be longer. Add tests for these
	// TODO: Add column major
	// TODO: Add tests for all the bad inputs
	// TODO: Add tests for dimension mismatch
	// TODO: Add negative increments
	// TODO: Add places with a "submatrix view", where lda != m
}

func sliceOfSliceCopy(a [][]float64) [][]float64 {
	n := make([][]float64, len(a))
	for i := range a {
		n[i] = make([]float64, len(a[i]))
		copy(n[i], a[i])
	}
	return n
}

func sliceCopy(a []float64) []float64 {
	n := make([]float64, len(a))
	copy(n, a)
	return n
}

func flatten(a [][]float64, o blas.Order) []float64 {
	if len(a) == 0 {
		return nil
	}
	m := len(a)
	n := len(a[0])
	s := make([]float64, m*n)
	if o == blas.RowMajor {
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				s[i*n+j] = a[i][j]
			}
		}
		return s
	}
	if o == blas.ColMajor {
		for j := 0; j < n; j++ {
			for i := 0; i < m; i++ {
				s[j*m+i] = a[i][j]
			}
		}
		return s
	}
	return nil
}

type Dgemver interface {
	Dgemv(o blas.Order, tA blas.Transpose, m, n int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
}

func DgemvTest(t *testing.T, blasser Dgemver) {
	for _, test := range DoubleMatTwoVecCases {
		for i, cas := range test.DgemvCases {
			x := sliceCopy(test.x)
			y := sliceCopy(test.y)
			a := sliceOfSliceCopy(test.A)
			aFlat := flatten(a, test.o)
			f := func() {
				blasser.Dgemv(test.o, test.tA, test.m, test.n, cas.alpha, aFlat, test.lda, x, test.incX, cas.beta, y, test.incY)
			}
			if panics(f) {
				if !test.Panics {
					t.Errorf("Test %v case %v unexpected panic", test.Name, i)
					if throwPanic {
						blasser.Dgemv(test.o, test.tA, test.m, test.n, cas.alpha, aFlat, test.lda, x, test.incX, cas.beta, y, test.incY)
					}
				}
				continue
			}
			// Check that x and a are unchanged
			if !dSliceEqual(x, test.x) {
				t.Errorf("Test %v, case %v x modified during call", test.Name, i)
			}
			aFlat2 := flatten(sliceOfSliceCopy(test.A), test.o)
			if !dSliceEqual(aFlat2, aFlat) {
				t.Errorf("Test %v, case %v a modified during call", test.Name, i)
			}

			// Check that the answer matches
			if !dSliceTolEqual(cas.ans, y) {
				t.Errorf("Test %v, case %v answer mismatch: Expected %v, Found %v", test.Name, i, cas.ans, y)
			}
		}
	}
}
