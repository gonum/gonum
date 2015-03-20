package mat64

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
)

var (
	triangular *Triangular
	_          Matrix = triangular
)

// Triangular represents an upper or lower triangular matrix.
type Triangular struct {
	mat blas64.Triangular
}

// NewTriangular constructs an n x n triangular matrix. The constructed matrix
// is upper triangular if upper == true and lower triangular otherwise.
// If len(mat) == n * n, mat will be used to hold the underlying data, if
// mat == nil, new data will be allocated, and will panic if neither of these
// cases is true.
// The underlying data representation is the same as that of a Dense matrix,
// except the values of the entries in the opposite half are completely ignored.
func NewTriangular(n int, upper bool, mat []float64) *Triangular {
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
	return &Triangular{blas64.Triangular{
		N:      n,
		Stride: n,
		Data:   mat,
		Uplo:   uplo,
		Diag:   blas.NonUnit,
	}}
}

func (t *Triangular) Dims() (r, c int) {
	return t.mat.N, t.mat.N
}

func (t *Triangular) RawTriangular() blas64.Triangular {
	return t.mat
}

func (t *Triangular) isZero() bool {
	// It must be the case that t.Dims() returns
	// zeros in this case. See comment in Reset().
	return t.mat.Stride == 0
}

// Reset zeros the dimensions of the matrix so that it can be reused as the
// receiver of a dimensionally restricted operation.
//
// See the Reseter interface for more information.
func (t *Triangular) Reset() {
	// No change of Stride, N to 0 may
	// be made unless both are set to 0.
	t.mat.N, t.mat.Stride = 0, 0
	// Defensively zero Uplo to ensure
	// it is set correctly later.
	t.mat.Uplo = 0
	t.mat.Data = t.mat.Data[:0]
}
