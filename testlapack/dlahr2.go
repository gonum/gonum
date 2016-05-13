// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/gonum/floats"
)

type Dlahr2er interface {
	Dlahr2(n, k, nb int, a []float64, lda int, tau, t []float64, ldt int, y []float64, ldy int)
}

type Dlahr2test struct {
	N, K, NB int
	A        []float64

	AWant   []float64
	TWant   []float64
	YWant   []float64
	TauWant []float64
}

func Dlahr2Test(t *testing.T, impl Dlahr2er) {
	// Go runs tests from the source directory, so unfortunately we need to
	// include the "../testlapack" part.
	file, err := os.Open(filepath.FromSlash("../testlapack/testdata/dlahr2data.json.gz"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	r, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	var tests []Dlahr2test
	json.NewDecoder(r).Decode(&tests)
	for _, test := range tests {
		tau := make([]float64, len(test.TauWant))
		for _, ldex := range []int{0, 1, 20} {
			n := test.N
			k := test.K
			nb := test.NB

			lda := n - k + 1 + ldex
			a := make([]float64, (n-1)*lda+n-k+1)
			copyMatrix(n, n-k+1, a, lda, test.A)

			ldt := nb + ldex
			tmat := make([]float64, (nb-1)*ldt+nb)

			ldy := nb + ldex
			y := make([]float64, (n-1)*ldy+nb)

			impl.Dlahr2(n, k, nb, a, lda, tau, tmat, ldt, y, ldy)

			prefix := fmt.Sprintf("n=%v, k=%v, nb=%v, ldex=%v", n, k, nb, ldex)
			if !equalApprox(n, n-k+1, a, lda, test.AWant, 1e-14) {
				t.Errorf("Case %v: unexpected matrix A\n got=%v\nwant=%v", prefix, a, test.AWant)
			}
			if !equalApproxTriangular(true, nb, tmat, ldt, test.TWant, 1e-14) {
				t.Errorf("Case %v: unexpected matrix T\n got=%v\nwant=%v", prefix, tmat, test.TWant)
			}
			if !equalApprox(n, nb, y, ldy, test.YWant, 1e-14) {
				t.Errorf("Case %v: unexpected matrix Y\n got=%v\nwant=%v", prefix, y, test.YWant)
			}
			if !floats.EqualApprox(tau, test.TauWant, 1e-14) {
				t.Errorf("Case %v: unexpected slice tau\n got=%v\nwant=%v", prefix, tau, test.TauWant)
			}
		}
	}
}
