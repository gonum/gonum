package lapack

import "github.com/gonum/blas"

const None = 'N'

type Job byte

const (
	All       (Job) = 'A'
	Slim      (Job) = 'S'
	Overwrite (Job) = 'O'
)

type CompSV byte

const (
	Compact  (CompSV) = 'P'
	Explicit (CompSV) = 'I'
)

// Float64 defines the float64 interface for the Lapack function. This interface
// contains the functions needed in the gonum suite.
type Float64 interface {
	Dpotrf(ul blas.Uplo, n int, a []float64, lda int) (ok bool)
}

type Complex128 interface{}
