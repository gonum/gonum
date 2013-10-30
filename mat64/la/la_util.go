// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Based on the SingularValueDecomposition class from Jama 1.0.3.

package la

import (
	"github.com/gonum/matrix/mat64"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Det returns the determinant of the matrix a.
func Det(a mat64.Matrix) float64 {
	lu, _, sign := LUD(mat64.DenseCopyOf(a))
	return LUDet(lu, sign)
}

// Inverse returns the inverse or pseudoinverse of the matrix a.
func Inverse(a mat64.Matrix) *mat64.Dense {
	m, _ := a.Dims()
	d := make([]float64, m*m)
	for i := 0; i < m*m; i += m + 1 {
		d[i] = 1
	}
	eye, _ := mat64.NewDense(m, m, d)
	return Solve(a, eye)
}

// Solve returns a matrix x that satisfies ax = b.
func Solve(a, b mat64.Matrix) (x *mat64.Dense) {
	m, n := a.Dims()
	if m == n {
		lu, piv, _ := LUD(mat64.DenseCopyOf(a))
		return LUSolve(lu, mat64.DenseCopyOf(b), piv)
	}
	qr, rDiag := QRD(mat64.DenseCopyOf(a))
	return QRSolve(qr, mat64.DenseCopyOf(b), rDiag)
}
