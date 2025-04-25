package umeyama

import (
	"errors"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
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
// The points require consistent indexing. This means
// that point i of X needs to correspond to point i of Y.
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

	rows, cols := rowsX, colsX

	// Calculate means.
	muX := mat.NewVecDense(rows, nil)
	muY := mat.NewVecDense(rows, nil)

	for i := 0; i < rows; i++ {
		muX.SetVec(i, stat.Mean(mat.Row(nil, i, X), nil))
		muY.SetVec(i, stat.Mean(mat.Row(nil, i, Y), nil))
	}

	// Calculate variance of X.
	varX := 0.0
	for j := 0; j < cols; j++ {
		for i := 0; i < rows; i++ {
			diff := X.At(i, j) - muX.AtVec(i)
			varX += diff * diff
		}
	}
	varX /= float64(cols)

	// Check for degenerate case.
	if varX < 1e-10 {
		return 0, nil, nil, errors.New(ErrMsgDegenerateInput)
	}

	// Calculate covariance matrix.
	covXY := mat.NewDense(rows, rows, nil)
	for j := 0; j < cols; j++ {
		for i := 0; i < rows; i++ {
			diffY := Y.At(i, j) - muY.AtVec(i)
			for k := 0; k < rows; k++ {
				diffX := X.At(k, j) - muX.AtVec(k)
				covXY.Set(i, k, covXY.At(i, k)+diffY*diffX/float64(cols))
			}
		}
	}

	// Singular Value Decomposition
	var svd mat.SVD
	ok := svd.Factorize(covXY, mat.SVDFull)
	if !ok {
		return 0, nil, nil, errors.New(ErrMsgSVDFailed)
	}

	// Get U, Σ, and V.
	u := mat.NewDense(rows, rows, nil)
	vt := mat.NewDense(rows, rows, nil)
	svd.UTo(u)
	svd.VTo(vt)

	// Transpose V to get VH.
	vh := mat.NewDense(rows, rows, nil)
	vh.Copy(vt.T())

	// Create S matrix (identity matrix).
	s := mat.NewDiagDense(rows, nil)
	for i := 0; i < rows; i++ {
		s.SetDiag(i, 1.0)
	}

	// Check determinants to ensure proper rotation matrix (not reflection).
	// This is a key contribution of Umeyama's paper - ensuring a proper rotation.
	uDet := mat.Det(u)
	vhDet := mat.Det(vh)
	if uDet*vhDet < 0 {
		s.SetDiag(rows-1, -1.0)
	}

	// Calculate scale factor c.
	c := 0.0
	singularValues := svd.Values(nil)
	for i := 0; i < rows; i++ {
		c += singularValues[i] * s.At(i, i)
	}
	c /= varX

	// Calculate rotation matrix R.
	tmp := mat.NewDense(rows, rows, nil)
	r := mat.NewDense(rows, rows, nil)
	tmp.Mul(u, s)
	r.Mul(tmp, vh)

	// Calculate translation vector t.
	t := mat.NewVecDense(rows, nil)
	rMuX := mat.NewVecDense(rows, nil)
	tmp2 := mat.NewDense(rows, 1, nil)
	for i := 0; i < rows; i++ {
		for j := 0; j < rows; j++ {
			tmp2.Set(i, 0, tmp2.At(i, 0)+r.At(i, j)*muX.AtVec(j))
		}
	}
	for i := 0; i < rows; i++ {
		rMuX.SetVec(i, tmp2.At(i, 0))
	}

	for i := 0; i < rows; i++ {
		t.SetVec(i, muY.AtVec(i)-c*rMuX.AtVec(i))
	}

	return c, r, t, nil
}
