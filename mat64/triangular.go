package mat64

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
)

var (
	triDense *TriDense
	_        Matrix        = triDense
	_        Triangular    = triDense
	_        RawTriangular = triDense
)

// TriDense represents an upper or lower triangular matrix in dense storage
// format.
type TriDense struct {
	mat blas64.Triangular
}

type Triangular interface {
	Matrix
	// Triangular returns the number of rows/columns in the matrix and if it is
	// an upper triangular matrix.
	Triangle() (n int, upper bool)

	// TTri is the equivalent of the T() method in the Matrix interface but
	// guarantees the transpose is of triangular type.
	TTri() Triangular
}

type RawTriangular interface {
	RawTriangular() blas64.Triangular
}

var (
	_ Matrix     = TransposeTri{}
	_ Triangular = TransposeTri{}
)

// TransposeTri is a type for performing an implicit transpose of a Triangular
// matrix. It implements the Triangular interface, returning values from the
// transpose of the matrix within.
type TransposeTri struct {
	Triangular Triangular
}

// At returns the value of the element at row i and column j of the transposed
// matrix, that is, row j and column i of the Triangular field.
func (t TransposeTri) At(i, j int) float64 {
	return t.Triangular.At(j, i)
}

// Dims returns the dimensions of the transposed matrix. Triangular matrices are
// square and thus this is the same size as the original Triangular.
func (t TransposeTri) Dims() (r, c int) {
	c, r = t.Triangular.Dims()
	return r, c
}

// T performs an implicit transpose by returning the Triangular field.
func (t TransposeTri) T() Matrix {
	return t.Triangular
}

// Triangle returns the number of rows/columns in the matrix and if it is
// an upper triangular matrix.
func (t TransposeTri) Triangle() (int, bool) {
	n, upper := t.Triangular.Triangle()
	return n, !upper
}

// TTri performs an implicit transpose by returning the Triangular field.
func (t TransposeTri) TTri() Triangular {
	return t.Triangular
}

// Untranspose returns the Triangular field.
func (t TransposeTri) Untranspose() Matrix {
	return t.Triangular
}

// NewTriangular constructs an n x n triangular matrix. The constructed matrix
// is upper triangular if upper == true and lower triangular otherwise.
// If len(mat) == n * n, mat will be used to hold the underlying data, if
// mat == nil, new data will be allocated, and will panic if neither of these
// cases is true.
// The underlying data representation is the same as that of a Dense matrix,
// except the values of the entries in the opposite half are completely ignored.
func NewTriDense(n int, upper bool, mat []float64) *TriDense {
	if n < 0 {
		panic("mat64: negative dimension")
	}
	if mat != nil && len(mat) != n*n {
		panic(ErrShape)
	}
	if mat == nil {
		mat = make([]float64, n*n)
	}
	uplo := blas.Lower
	if upper {
		uplo = blas.Upper
	}
	return &TriDense{blas64.Triangular{
		N:      n,
		Stride: n,
		Data:   mat,
		Uplo:   uplo,
		Diag:   blas.NonUnit,
	}}
}

func (t *TriDense) Dims() (r, c int) {
	return t.mat.N, t.mat.N
}

func (t *TriDense) Triangle() (n int, upper bool) {
	return t.mat.N, t.mat.Uplo == blas.Upper
}

// T performs an implicit transpose by returning the receiver inside a Transpose.
func (t *TriDense) T() Matrix {
	return Transpose{t}
}

// TTri performs an implicit transpose by returning the receiver inside a TransposeTri.
func (t *TriDense) TTri() Triangular {
	return TransposeTri{t}
}

func (t *TriDense) RawTriangular() blas64.Triangular {
	return t.mat
}

func (t *TriDense) isZero() bool {
	// It must be the case that t.Dims() returns
	// zeros in this case. See comment in Reset().
	return t.mat.Stride == 0
}

// reuseAS resizes a zero receiver to an n×n triangular matrix with the given
// orientation. If the receiver is non-zero, reuseAs checks that the receiver
// is the correct size and orientation.
func (t *TriDense) reuseAs(n int, ul blas.Uplo) {
	if t.isZero() {
		t.mat = blas64.Triangular{
			N:      n,
			Stride: n,
			Diag:   blas.NonUnit,
			Data:   use(t.mat.Data, n*n),
			Uplo:   ul,
		}
		return
	}
	if t.mat.N != n || t.mat.Uplo != ul {
		panic(ErrShape)
	}
}

// Reset zeros the dimensions of the matrix so that it can be reused as the
// receiver of a dimensionally restricted operation.
//
// See the Reseter interface for more information.
func (t *TriDense) Reset() {
	// No change of Stride, N to 0 may
	// be made unless both are set to 0.
	t.mat.N, t.mat.Stride = 0, 0
	// Defensively zero Uplo to ensure
	// it is set correctly later.
	t.mat.Uplo = 0
	t.mat.Data = t.mat.Data[:0]
}

// getBlasTriangular transforms t into a blas64.Triangular. If t is a RawTriangular,
// the direct matrix representation is returned, otherwise t is copied into one.
func getBlasTriangular(t Triangular) blas64.Triangular {
	n, upper := t.Triangle()
	rt, ok := t.(RawTriangular)
	if ok {
		return rt.RawTriangular()
	}
	ta := blas64.Triangular{
		N:      n,
		Stride: n,
		Diag:   blas.NonUnit,
		Data:   make([]float64, n*n),
	}
	if upper {
		ta.Uplo = blas.Upper
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				ta.Data[i*n+j] = t.At(i, j)
			}
		}
		return ta
	}
	ta.Uplo = blas.Lower
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			ta.Data[i*n+j] = t.At(i, j)
		}
	}
	return ta
}

// copySymIntoTriangle copies a symmetric matrix into a TriDense
func copySymIntoTriangle(t *TriDense, s Symmetric) {
	n, upper := t.Triangle()
	ns := s.Symmetric()
	if n != ns {
		panic("mat64: triangle size mismatch")
	}
	ts := t.mat.Stride
	if rs, ok := s.(RawSymmetricer); ok {
		sd := rs.RawSymmetric()
		ss := sd.Stride
		if upper {
			if sd.Uplo == blas.Upper {
				for i := 0; i < n; i++ {
					copy(t.mat.Data[i*ts+i:i*ts+n], sd.Data[i*ss+i:i*ss+n])
				}
				return
			}
			for i := 0; i < n; i++ {
				for j := i; j < n; j++ {
					t.mat.Data[i*ts+j] = sd.Data[j*ss+i]
				}
				return
			}
		}
		if sd.Uplo == blas.Upper {
			for i := 0; i < n; i++ {
				for j := 0; j <= i; j++ {
					t.mat.Data[i*ts+j] = sd.Data[j*ss+i]
				}
			}
			return
		}
		for i := 0; i < n; i++ {
			copy(t.mat.Data[i*ts:i*ts+i+1], sd.Data[i*ss:i*ss+i+1])
		}
		return
	}
	if upper {
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				t.mat.Data[i*ts+j] = s.At(i, j)
			}
		}
		return
	}
	for i := 0; i < n; i++ {
		for j := 0; j <= i; j++ {
			t.mat.Data[i*ts+j] = s.At(i, j)
		}
	}
}
