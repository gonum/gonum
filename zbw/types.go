package zbw

import (
	"errors"
	"fmt"

	"github.com/gonum/blas"
)

func Allocate(dims ...int) []complex128 {
	if len(dims) == 0 {
		return nil
	}
	n := 1
	for _, v := range dims {
		n *= v
	}
	return make([]complex128, n)
}

type General struct {
	Rows, Cols int
	Stride     int
	Data       []complex128
}

func NewGeneral(m, n int, data []complex128) General {
	var A General
	if data == nil {
		data = make([]complex128, m*n)
	}
	A = General{m, n, n, data}
	must(A.Check())
	return A
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
	Data   []complex128
	N      int
	Stride int
	Uplo   blas.Uplo
	Diag   blas.Diag
}

type TriangularBand struct {
	Data   []complex128
	N, K   int
	Stride int
	Uplo   blas.Uplo
	Diag   blas.Diag
}

type TriangularPacked struct {
	Data []complex128
	N    int
	Uplo blas.Uplo
	Diag blas.Diag
}

type Symmetric struct {
	Data      []complex128
	N, Stride int
	Uplo      blas.Uplo
}

type Hermitian struct {
	Data      []complex128
	N, Stride int
	Uplo      blas.Uplo
}

type HermitianBand struct {
	Data         []complex128
	N, K, Stride int
	Uplo         blas.Uplo
}

type HermitianPacked struct {
	Data []complex128
	N    int
	Uplo blas.Uplo
}

type Vector struct {
	Data []complex128
	N    int
	Inc  int
}

func NewVector(v []complex128) Vector {
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

func Ge2He(A General, ul blas.Uplo) Hermitian {
	n := A.Rows
	if A.Cols < n {
		n = A.Cols
	}
	return Hermitian{A.Data, n, A.Stride, ul}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func Real2Cmplx(r []float64, c []complex128) {
	if len(r) != len(c) {
		panic("length missmatch")
	}
	for ix, v := range r {
		c[ix] = complex(v, 0)
	}
}
