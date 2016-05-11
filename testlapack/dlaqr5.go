// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"testing"

	"github.com/gonum/floats"
)

type Dlaqr5er interface {
	Dlaqr5(wantt, wantz bool, kacc22 int, n, ktop, kbot, nshfts int, sr, si []float64, h []float64, ldh int, iloz, ihiz int, z []float64, ldz int, v []float64, ldv int, u []float64, ldu int, nh int, wh []float64, ldwh int, nv int, wv []float64, ldwv int)
}

func Dlaqr5Test(t *testing.T, impl Dlaqr5er) {
	r, err := zip.OpenReader("../internal/testdata/dlaqr5test/dlaqr5data.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		tc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		wantt, n, nshfts, ktop, kbot, sr, si, hOrig, hwant, zwant := readCase(tc)
		tc.Close()

		v := make([]float64, nshfts/2*3)
		u := make([]float64, (3*nshfts-3)*(3*nshfts-3))
		nh := n
		wh := make([]float64, (3*nshfts-3)*n)
		nv := n
		wv := make([]float64, n*(3*nshfts-3))

		for _, ldh := range []int{n, n + 1, n + 10} {
			h := make([]float64, (n-1)*ldh+n)
			for _, kacc22 := range []int{0, 1, 2} {
				copySquareMatrix(n, h, ldh, hOrig)
				z := eye(n, ldh)

				impl.Dlaqr5(wantt, true, kacc22,
					n, ktop, kbot,
					nshfts, sr, si,
					h, ldh,
					0, n-1, z, ldh,
					v, 3,
					u, 3*nshfts-3,
					nh, wh, nh,
					nv, wv, 3*nshfts-3)

				if !equalApprox(n, h, ldh, hwant, 1e-13) {
					t.Errorf("Case %v, kacc22=%v: unexpected matrix H\nh    =%v\nhwant=%v", f.Name, kacc22, h, hwant)
				}
				if !equalApprox(n, z, ldh, zwant, 1e-13) {
					t.Errorf("Case %v, kacc22=%v: unexpected matrix Z\nz    =%v\nzwant=%v", f.Name, kacc22, z, zwant)
				}
			}
		}
	}
}

// copySquareMatrix copies a square matrix src of order and stride n into a
// square matrix dst of order n and stride ldh.
func copySquareMatrix(n int, dst []float64, ldh int, src []float64) {
	for i := 0; i < n; i++ {
		copy(dst[i*ldh:i*ldh+n], src[i*n:i*n+n])
	}
}

// equalApprox returns whether the square matrices A and B of order n are
// approximately equal within given tolerance.
func equalApprox(n int, a []float64, lda int, b []float64, tol float64) bool {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if !floats.EqualWithinAbsOrRel(a[i*lda+j], b[i*n+j], tol, tol) {
				return false
			}
		}
	}
	return true
}

// readCase reads and returns test data written by internal/testdata/dlaqr5test/main.go.
func readCase(r io.Reader) (wantt bool, n, nshfts, ktop, kbot int, sr, si []float64, h, hwant, zwant []float64) {
	_, err := fmt.Fscanln(r, &wantt, &n, &nshfts, &ktop, &kbot)
	if err != nil {
		log.Fatal(err)
	}

	sr = make([]float64, nshfts)
	si = make([]float64, nshfts)
	h = make([]float64, n*n)
	hwant = make([]float64, n*n)
	zwant = make([]float64, n*n)

	for i := range sr {
		_, err = fmt.Fscanln(r, &sr[i])
		if err != nil {
			log.Fatal(err)
		}
	}
	for i := range si {
		_, err = fmt.Fscanln(r, &si[i])
		if err != nil {
			log.Fatal(err)
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			_, err = fmt.Fscanln(r, &h[i*n+j])
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			_, err = fmt.Fscanln(r, &hwant[i*n+j])
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			_, err = fmt.Fscanln(r, &zwant[i*n+j])
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return wantt, n, nshfts, ktop, kbot, sr, si, h, hwant, zwant
}

// eye returns an identity matrix of order n and stride ld.
func eye(n, ld int) []float64 {
	m := make([]float64, (n-1)*ld+n)
	for i := 0; i < n; i++ {
		m[i*ld+i] = 1
	}
	return m
}
