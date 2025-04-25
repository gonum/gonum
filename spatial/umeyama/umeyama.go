package umeyama

import (
	"errors"

	"gonum.org/v1/gonum/mat"
)

const (
	// ErrMsgSVDFailed is the error message for SVD factorization failure
	ErrMsgSVDFailed = "umeyama: SVD factorization failed"
	// ErrMsgDegenerateInput is the error message for degenerate input data
	ErrMsgDegenerateInput = "umeyama: variance of X is too small"
)

// Umeyama estimates the similarity transformation parameters between two matrices X and Y.
//
// This is an implementation of the algorithm presented in:
// "Least-Squares Estimation of Transformation Parameters Between Two Point Patterns"
// by Shinji Umeyama, IEEE Transactions on Pattern Analysis and Machine Intelligence,
// Vol. 13, No. 4, April 1991, which can be found here: https://doi.org/10.1109/34.88573
//
// The algorithm finds the optimal similarity transformation [c, R, t] ∈ Sim(m)
// that minimizes the mean squared error between two point patterns.
//
// The transformation relates the point sets as:
// Y ≈ c * R * X + t
//
// The dimensions of X and Y must be equal. The function will panic if they are not.
// The points require consistent indexing. This means that point i of X needs to correspond
// to point i of Y.
//
// In this implementation, rows represent points and columns represent dimensions.
//
// Umeyama returns the scale factor c, the rotation matrix r and the translation vector t.
// If a computation fails, Umeyama will return an error.
func Umeyama(X, Y *mat.Dense) (float64, *mat.Dense, *mat.VecDense, error) {
	rowsX, colsX := X.Dims()
	rowsY, colsY := Y.Dims()

	// Check dimensions.
	if rowsX != rowsY || colsX != colsY {
		panic("umeyama: dimensions of X and Y do not match")
	}

	n := rowsX // number of points
	m := colsX // number of dimensions

	// Calculate means.
	muX := mat.NewVecDense(m, nil)
	muY := mat.NewVecDense(m, nil)

	for i := 0; i < m; i++ {
		var sumX, sumY float64
		for j := 0; j < n; j++ {
			sumX += X.At(j, i)
			sumY += Y.At(j, i)
		}
		muX.SetVec(i, sumX/float64(n))
		muY.SetVec(i, sumY/float64(n))
	}

	// Calculate variance of X.
	var varX float64
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			diff := X.At(i, j) - muX.AtVec(j)
			varX += diff * diff
		}
	}
	varX /= float64(n)

	// Check for degenerate case.
	if varX < 1e-10 {
		return 0, nil, nil, errors.New(ErrMsgDegenerateInput)
	}

	// Calculate covariance matrix.
	covXY := mat.NewDense(m, m, nil)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			diffY := Y.At(i, j) - muY.AtVec(j)
			for k := 0; k < m; k++ {
				diffX := X.At(i, k) - muX.AtVec(k)
				covXY.Set(j, k, covXY.At(j, k)+diffY*diffX/float64(n))
			}
		}
	}

	// Singular Value Decomposition
	var svd mat.SVD
	if !svd.Factorize(covXY, mat.SVDFull) {
		return 0, nil, nil, errors.New(ErrMsgSVDFailed)
	}

	// Get U and V.
	u := mat.NewDense(m, m, nil)
	v := mat.NewDense(m, m, nil)
	svd.UTo(u)
	svd.VTo(v)

	// Create S matrix (identity matrix).
	s := mat.NewDiagDense(m, nil)
	for i := 0; i < m; i++ {
		s.SetDiag(i, 1.0)
	}

	// Check determinants to ensure proper rotation matrix (not reflection).
	if mat.Det(u)*mat.Det(v) < 0 {
		s.SetDiag(m-1, -1.0)
	}

	// Calculate scale factor c.
	var c float64
	singularValues := svd.Values(nil)
	for i := 0; i < m; i++ {
		c += singularValues[i] * s.At(i, i)
	}
	c /= varX

	// Calculate rotation matrix R.
	r := mat.NewDense(m, m, nil)
	tmp := mat.NewDense(m, m, nil)
	tmp.Mul(u, s)
	r.Mul(tmp, v.T())

	// Calculate translation vector t.
	t := mat.NewVecDense(m, nil)
	rMuX := mat.NewVecDense(m, nil)
	rMuX.MulVec(r, muX)

	for i := 0; i < m; i++ {
		t.SetVec(i, muY.AtVec(i)-c*rMuX.AtVec(i))
	}

	return c, r, t, nil
}
