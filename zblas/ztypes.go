package blas

import (
	"errors"
)

type GeneralCmplx struct {
	Order      Order
	Rows, Cols int
	Stride     int
	Data       []complex128
}

func NewGeneralCmplx(o Order, m, n int, data []complex128) GeneralCmplx {
	var A GeneralCmplx
	if o == RowMajor {
		A = GeneralCmplx{o, m, n, n, data}
	} else {
		A = GeneralCmplx{o, m, n, m, data}
	}
	must(A.Check())
	return A
}

func (A GeneralCmplx) Index(i, j int) int {
	if A.Order == RowMajor {
		return i*A.Stride + j
	} else {
		return i + j*A.Stride
	}
}

func (A GeneralCmplx) Check() error {
	if A.Cols < 0 {
		return errors.New("blas: n < 0")
	}
	if A.Rows < 0 {
		return errors.New("blas: m < 0")
	}
	if A.Stride < 1 {
		return errors.New("blas: illegal stride")
	}
	if A.Order == ColMajor {
		if A.Stride < A.Rows {
			return errors.New("blas: illegal stride")
		}
		if (A.Cols-1)*A.Stride+A.Rows > len(A.Data) {
			return errors.New("blas: insufficient amount of data")
		}
	} else if A.Order == RowMajor {
		if A.Stride < A.Cols {
			return errors.New("blas: illegal stride")
		}
		if (A.Rows-1)*A.Stride+A.Cols > len(A.Data) {
			return errors.New("blas: insufficient amount of data")
		}
	} else {
		return errors.New("blas: illegal order")
	}
	return nil
}

func (A GeneralCmplx) Row(i int) VectorCmplx {
	if i >= A.Rows || i < 0 {
		panic("blas: index out of range")
	}
	if A.Order == RowMajor {
		return VectorCmplx{A.Data[A.Stride*i:], A.Cols, 1}
	} else if A.Order == ColMajor {
		return VectorCmplx{A.Data[i:], A.Cols, A.Stride}
	}
	panic("blas: illegal order")
}

func (A GeneralCmplx) Col(i int) VectorCmplx {
	if i >= A.Cols || i < 0 {
		panic("blas: index out of range")
	}
	if A.Order == RowMajor {
		return VectorCmplx{A.Data[i:], A.Rows, A.Stride}
	} else if A.Order == ColMajor {
		return VectorCmplx{A.Data[A.Stride*i:], A.Rows, 1}
	}
	panic("blas: illegal order")
}

func (A GeneralCmplx) Sub(i, j, r, c int) GeneralCmplx {
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
	return GeneralCmplx{A.Order, r, c, A.Stride, A.Data[A.Index(i, j):]}
}

type GeneralCmplxBand struct {
	Order Order
	GeneralCmplx
	KL, KU int
}

type TriangularCmplx struct {
	Order  Order
	Data   []complex128
	N      int
	Stride int
	Uplo   Uplo
	Diag   Diag
}

type TriangularCmplxBand struct {
	Order  Order
	Data   []complex128
	N, K   int
	Stride int
	Uplo   Uplo
	Diag   Diag
}

type TriangularCmplxPacked struct {
	Order Order
	Data  []complex128
	N     int
	Uplo  Uplo
	Diag  Diag
}

type SymmetricCmplx struct {
	Order     Order
	Data      []complex128
	N, Stride int
	Uplo      Uplo
}

type Hermitian struct {
	Order     Order
	Data      []complex128
	N, Stride int
	Uplo      Uplo
}

type HermitianBand struct {
	Order        Order
	Data         []complex128
	N, K, Stride int
	Uplo         Uplo
}

type HermitianPacked struct {
	Order Order
	Data  []complex128
	N     int
	Uplo  Uplo
}

type VectorCmplx struct {
	Data []complex128
	N    int
	Inc  int
}

func NewVectorCmplx(v []complex128) VectorCmplx {
	return VectorCmplx{v, len(v), 1}
}

func (v VectorCmplx) Check() error {
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

func Gec2Trc(A GeneralCmplx, d Diag, ul Uplo) TriangularCmplx {
	n := A.Rows
	if A.Cols < n {
		n = A.Cols
	}
	return TriangularCmplx{A.Order, A.Data, n, A.Stride, ul, d}
}

func Gec2He(A GeneralCmplx, ul Uplo) Hermitian {
	n := A.Rows
	if A.Cols < n {
		n = A.Cols
	}
	return Hermitian{A.Order, A.Data, n, A.Stride, ul}
}
