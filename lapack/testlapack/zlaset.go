// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"testing"

	"gonum.org/v1/gonum/blas"
)

type Zlaseter interface {
	Zlaset(uplo blas.Uplo, m, n int, alpha, beta complex128, a []complex128, lda int)
}

func ZlasetTest(t *testing.T, impl Zlaseter) {
	for ti, test := range []struct {
		m, n int
	}{
		{0, 0},
		{1, 1},
		{1, 10},
		{10, 1},
		{2, 2},
		{2, 10},
		{10, 2},
		{11, 11},
		{11, 50},
		{50, 11},
	} {
		m := test.m
		n := test.n
		for _, uplo := range []blas.Uplo{blas.Upper, blas.Lower, blas.All} {
			for _, extra := range []int{0, 7} {
				lda := n + extra
				if lda == 0 {
					lda = 1
				}
				// Allocate a, with out-of-range elements filled with a sentinel
				// so we can verify they are not touched.
				sentinel := complex(-7.5, 3.25)
				a := make([]complex128, max(1, m)*lda)
				for i := range a {
					a[i] = sentinel
				}
				alpha := complex(1, -2)
				beta := complex(3, 4)

				impl.Zlaset(uplo, m, n, alpha, beta, a, lda)

				prefix := fmt.Sprintf("Case #%v: m=%v,n=%v,uplo=%c,extra=%v",
					ti, m, n, uplo, extra)
				// Check out-of-range untouched.
				for i := 0; i < m; i++ {
					for j := n; j < lda; j++ {
						if a[i*lda+j] != sentinel {
							t.Errorf("%v: out-of-range write at (%d,%d)", prefix, i, j)
						}
					}
				}
				for i := 0; i < min(m, n); i++ {
					if a[i*lda+i] != beta {
						t.Errorf("%v: unexpected diagonal of A at %d: %v want %v", prefix, i, a[i*lda+i], beta)
					}
				}
				if uplo == blas.Upper || uplo == blas.All {
					for i := 0; i < m; i++ {
						for j := i + 1; j < n; j++ {
							if a[i*lda+j] != alpha {
								t.Errorf("%v: unexpected upper triangle at (%d,%d)", prefix, i, j)
							}
						}
					}
				}
				if uplo == blas.Lower || uplo == blas.All {
					for i := 1; i < m; i++ {
						for j := 0; j < min(i, n); j++ {
							if a[i*lda+j] != alpha {
								t.Errorf("%v: unexpected lower triangle at (%d,%d)", prefix, i, j)
							}
						}
					}
				}
				// Off-triangle region should be untouched if uplo restricts.
				if uplo == blas.Upper {
					for i := 1; i < m; i++ {
						for j := 0; j < min(i, n); j++ {
							if a[i*lda+j] != sentinel {
								t.Errorf("%v: lower triangle modified at (%d,%d)", prefix, i, j)
							}
						}
					}
				}
				if uplo == blas.Lower {
					for i := 0; i < m; i++ {
						for j := i + 1; j < n; j++ {
							if a[i*lda+j] != sentinel {
								t.Errorf("%v: upper triangle modified at (%d,%d)", prefix, i, j)
							}
						}
					}
				}
			}
		}
	}
}
