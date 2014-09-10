package dbw

import "github.com/gonum/blas"

import (
	"errors"
	"fmt"
)

func Allocate(dims ...int) []float64 {
	if len(dims) == 0 {
		return nil
	}
	n := 1
	for _, v := range dims {
		n *= v
	}
	return make([]float64, n)
}

type General struct {
	Rows, Cols int
	Stride     int
	Data       []float64
}

func NewGeneral(m, n int, data []float64) General {
	var A General
	if data == nil {
		data = make([]float64, m*n)
	}
	A = General{m, n, n, data}
	must(A.Check())
	return A
}

func (A General) At(i, j int) float64 {
	return A.Data[A.Index(i, j)]
}

func (A General) Set(i, j int, v float64) {
	A.Data[A.Index(i, j)] = v
}

func (A General) Dims() (int, int) {
	return A.Rows, A.Cols
}

func (A General) Index(i, j int) int {
	return i*A.Stride + j
}

func (A General) Check() error {
	if A.Cols < 0 {
		return errors.New("blas: n < 0")
	}
	if A.Rows < 0 {
		return errors.New("blas: m < 0")
	}
	if A.Stride < 1 {
		return errors.New("blas: illegal stride")
	}
	if A.Stride < A.Cols {
		return errors.New("blas: illegal stride")
	}
	if (A.Rows-1)*A.Stride+A.Cols > len(A.Data) {
		return errors.New("blas: insufficient amount of data")
	}
	return nil
}

func (A General) Row(i int) Vector {
	if i >= A.Rows || i < 0 {
		panic("blas: index out of range")
	}
	return Vector{A.Data[A.Stride*i:], A.Cols, 1}
}

func (A General) Col(i int) Vector {
	if i >= A.Cols || i < 0 {
		panic("blas: index out of range")
	}

	return Vector{A.Data[i:], A.Rows, A.Stride}
}

func (A General) Sub(i, j, r, c int) General {
	must(A.Check())
	if i >= A.Rows || i < 0 {
		panic("blas: index out of range")
	}
	if j >= A.Cols || i < 0 {
		panic("blas: index out of range")
	}
	if r < 0 || c < 0 {
		panic("blas: r < 0 or c < 0")
	}
	return General{r, c, A.Stride, A.Data[A.Index(i, j):]}
}

type GeneralBand struct {
	General
	KL, KU int
}

type Triangular struct {
	Data   []float64
	N      int
	Stride int
	Uplo   blas.Uplo
	Diag   blas.Diag
}

type TriangularBand struct {
	Data   []float64
	N, K   int
	Stride int
	Uplo   blas.Uplo
	Diag   blas.Diag
}

type TriangularPacked struct {
	Data []float64
	N    int
	Uplo blas.Uplo
	Diag blas.Diag
}

type Symmetric struct {
	Data      []float64
	N, Stride int
	Uplo      blas.Uplo
}

type SymmetricBand struct {
	Data         []float64
	N, K, Stride int
	Uplo         blas.Uplo
}

type SymmetricPacked struct {
	Data []float64
	N    int
	Uplo blas.Uplo
}

type Vector struct {
	Data []float64
	N    int
	Inc  int
}

func NewVector(v []float64) Vector {
	return Vector{v, len(v), 1}
}

func (v Vector) Slice(l, r int) Vector {
	if l < 0 || r > v.N {
		panic("blas: index out of range")
	}
	if r > l {
		panic(fmt.Sprintf("blas: invalid slice index:", r, ">", l))
	}
	return Vector{v.Data[l*v.Inc:], r - l, v.Inc}
}

func (v Vector) Check() error {
	if v.N < 0 {
		return errors.New("blas: n < 0")
	}
	if v.Inc == 0 {
		return errors.New("blas: zero x index increment")
	}
	if (v.N-1)*v.Inc >= len(v.Data) {
		return errors.New("blas: index out of range")
	}
	return nil
}

func Ge2Tr(A General, d blas.Diag, ul blas.Uplo) Triangular {
	n := A.Rows
	if A.Cols < n {
		n = A.Cols
	}
	return Triangular{A.Data, n, A.Stride, ul, d}
}

func Ge2Sy(A General, ul blas.Uplo) Symmetric {
	n := A.Rows
	if A.Cols < n {
		n = A.Cols
	}
	return Symmetric{A.Data, n, A.Stride, ul}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
