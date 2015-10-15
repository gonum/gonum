// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package native

// Dgebd2 reduces an m×n matrix A to upper or lower bidiagonal form by an orthogonal
// transformation.
//  Q^T * A * P = B
// if m >= n, B is upper diagonal, otherwise B is lower bidiagonal.
// d is the diagonal, len = min(m,n)
// e is the off-diagonal len = min(m,n)-1
func (impl Implementation) Dgebd2(m, n int, a []float64, lda int, d, e, tauQ, tauP, work []float64) {
	checkMatrix(m, n, a, lda)
	if len(d) < min(m, n) {
		panic("lapack: insufficient d")
	}
	if len(e) < min(m, n)-1 {
		panic("lapack: insufficient e")
	}
	if m > n {
		for i := 0; i < n; i++ {
			impl.Dlarfg(m-i, a[i*lda+i], a[min(i+1, m-1)*lda+i:], 1)
		}
	}
}
