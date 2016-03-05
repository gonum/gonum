// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"math"

	"github.com/gonum/floats"
	"github.com/gonum/matrix"
	"github.com/gonum/matrix/mat64"
)

// PrincipalComponents returns the principal component direction vectors and
// the column variances of the principal component scores, vecs * a, computed
// using the singular value decomposition of the input. The input a is an n×d
// matrix where each row is an observation and each column represents a variable.
//
// PrincipalComponents centers the variables but does not scale the variance.
//
// The slice weights is used to weight the observations. If weights is nil,
// each weight is considered to have a value of one, otherwise the length of
// weights must match the number of observations or PrincipalComponents will
// panic.
//
// On successful completion, the principal component direction vectors are
// returned in vecs as a d×min(n, d) matrix, and the variances are returned in
// vars as a min(n, d)-long slice in descending sort order.
//
// If no singular value decomposition is possible, vecs and vars are returned
// nil and ok is returned false.
func PrincipalComponents(a mat64.Matrix, weights []float64) (vecs *mat64.Dense, vars []float64, ok bool) {
	n, d := a.Dims()
	if weights != nil && len(weights) != n {
		panic("stat: len(weights) != observations")
	}

	centered := mat64.NewDense(n, d, nil)
	col := make([]float64, n)
	for j := 0; j < d; j++ {
		mat64.Col(col, j, a)
		floats.AddConst(-Mean(col, weights), col)
		centered.SetCol(j, col)
	}
	for i, w := range weights {
		floats.Scale(math.Sqrt(w), centered.RawRowView(i))
	}

	kind := matrix.SVDFull
	if n > d {
		kind = matrix.SVDThin
	}
	var svd mat64.SVD
	ok = svd.Factorize(centered, kind)
	if !ok {
		return nil, nil, false
	}

	vecs = &mat64.Dense{}
	vecs.VFromSVD(&svd)
	if n < d {
		// Don't retain columns that are not valid direction vectors.
		vecs.Clone(vecs.View(0, 0, d, n))
	}
	vars = svd.Values(nil)
	var f float64
	if weights == nil {
		f = 1 / float64(n-1)
	} else {
		f = 1 / (floats.Sum(weights) - 1)
	}
	for i, v := range vars {
		vars[i] = f * v * v
	}
	return vecs, vars, true
}
