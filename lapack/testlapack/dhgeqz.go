// Copyright ©2023 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

/*
#cgo CFLAGS: -I/home/pato/src/ongoing/lapack/SRC/  -I/home/pato/src/ongoing/lapack/BLAS/SRC
#cgo LDFLAGS: -L/home/pato/src/ongoing/lapack -llapack -lrefblas -lgfortran -lm -ltmglib


void dhgeqz_(char * JOB, char * COMPQ, char * COMPZ, int * N, int * ILO, int * IHI, double * H, int * LDH, double * T, int * LDT, double * ALPHAR, double * ALPHAI, double * BETA, double * Q, int * LDQ, double * Z, int * LDZ, double * WORK, int * LWORK, int * INFO);
*/
import "C"

type Dhgeqzer interface {
	Dhgeqz(job lapack.SchurJob, compq, compz lapack.SchurComp, n, ilo, ihi int,
		h []float64, ldh int, t []float64, ldt int, alphar, alphai, beta,
		q []float64, ldq int, z []float64, ldz int, work []float64, workspaceQuery bool) (info int)
}

func DhgeqzTest(t *testing.T, impl Dhgeqzer) {
	src := uint64(8878)
	rnd := rand.New(rand.NewSource(src))
	const ldaAdd = 0
	compvec := []lapack.SchurComp{lapack.SchurNone, lapack.SchurHess, lapack.SchurOrig}
	for _, compq := range compvec {
		for _, compz := range compvec {
			for _, n := range []int{2, 3, 4, 16} {
				minLDA := max(1, n)
				for _, ldh := range []int{minLDA, n + ldaAdd} {
					for _, ldt := range []int{minLDA, n + ldaAdd} {
						for _, ldq := range []int{minLDA, n + ldaAdd} {
							for _, ldz := range []int{minLDA, n + ldaAdd} {
								for ilo := 0; ilo < n; ilo++ {
									for ihi := ilo; ihi < n; ihi++ {
										for cas := 0; cas < 1; cas++ {
											testDhgeqz(t, rnd, impl, lapack.EigenvaluesAndSchur, compq, compz, n, ilo, ihi, ldh, ldt, ldq, ldz)
											testDhgeqz(t, rnd, impl, lapack.EigenvaluesOnly, compq, compz, n, ilo, ihi, ldh, ldt, ldq, ldz)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func testDhgeqz(t *testing.T, rnd *rand.Rand, impl Dhgeqzer, job lapack.SchurJob, compq, compz lapack.SchurComp, n, ilo, ihi, ldh, ldt, ldq, ldz int) {
	name := fmt.Sprintf("Case job=%q, compq=%q, compz=%q, n=%v, ilo=%v, ihi=%v, ldh=%v, ldt=%v, ldq=%v, ldz=%v",
		job, compq, compz, n, ilo, ihi, ldh, ldt, ldq, ldz)
	generalFromComp := func(comp lapack.SchurComp, n, ld int, rnd *rand.Rand) blas64.General {
		switch comp {
		case lapack.SchurNone:
			return blas64.General{Stride: 1}
		case lapack.SchurHess:
			return nanGeneral(n, n, ld)
		case lapack.SchurOrig:
			return randomOrthogonal(n, rnd)
		default:
			panic("bad comp")
		}
	}
	hg := randomHessenberg(n, ldh, rnd)
	tg := blockedUpperTriGeneral(n, n, 0, n, ldt, false, rnd)

	alphar := make([]float64, n)
	alphai := make([]float64, n)
	beta := make([]float64, n)
	q := generalFromComp(compq, n, ldq, rnd)
	z := generalFromComp(compz, n, ldz, rnd)

	// Query workspace needed.
	var query [1]float64
	impl.Dhgeqz(job, compq, compz, n, ilo, ihi, hg.Data, hg.Stride, tg.Data, tg.Stride, alphar, alphai, beta, q.Data, q.Stride, z.Data, z.Stride, query[:], true)
	lwork := int(query[0])
	if lwork < 1 {
		t.Fatal(name, "bad lwork")
	}

	work := nanSlice(lwork)
	info := impl.Dhgeqz(job, compq, compz, n, ilo, ihi, hg.Data, hg.Stride, tg.Data, tg.Stride, alphar, alphai, beta, q.Data, q.Stride, z.Data, z.Stride, work, false)
	if info >= 0 {
		t.Error(name, "got nonzero info", info)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if math.IsNaN(hg.Data[i*hg.Stride+j]) {
				t.Fatal("H has NaN(s)")
			}
			if math.IsNaN(tg.Data[i*tg.Stride+j]) {
				t.Fatal("T has NaN(s)")
			}
			if compq == lapack.SchurHess || compq == lapack.SchurOrig {
				if math.IsNaN(q.Data[i*q.Stride+j]) {
					t.Fatal("Q has NaN(s)")
				}
			}
			if compz == lapack.SchurHess || compz == lapack.SchurOrig {
				if math.IsNaN(z.Data[i*z.Stride+j]) {
					t.Fatal("Z has NaN(s)")
				}
			}
		}
	}
}
