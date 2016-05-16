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
)

type Dlaqr5er interface {
	Dlaqr5(wantt, wantz bool, kacc22 int, n, ktop, kbot, nshfts int, sr, si []float64, h []float64, ldh int, iloz, ihiz int, z []float64, ldz int, v []float64, ldv int, u []float64, ldu int, nh int, wh []float64, ldwh int, nv int, wv []float64, ldwv int)
}

type Dlaqr5test struct {
	WantT          bool
	N              int
	NShifts        int
	KTop, KBot     int
	ShiftR, ShiftI []float64
	H              []float64

	HWant []float64
	ZWant []float64
}

func Dlaqr5Test(t *testing.T, impl Dlaqr5er) {
	file, err := os.Open(filepath.FromSlash("../testlapack/testdata/dlaqr5data.json.gz"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	r, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	var tests []Dlaqr5test
	json.NewDecoder(r).Decode(&tests)
	for _, test := range tests {
		wantt := test.WantT
		n := test.N
		nshfts := test.NShifts
		ktop := test.KTop
		kbot := test.KBot
		sr := test.ShiftR
		si := test.ShiftI

		v := make([]float64, nshfts/2*3)
		u := make([]float64, (3*nshfts-3)*(3*nshfts-3))
		nh := n
		wh := make([]float64, (3*nshfts-3)*n)
		nv := n
		wv := make([]float64, n*(3*nshfts-3))

		for _, ldh := range []int{n, n + 1, n + 10} {
			h := make([]float64, (n-1)*ldh+n)

			for _, kacc22 := range []int{0, 1, 2} {
				copyMatrix(n, n, h, ldh, test.H)
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

				prefix := fmt.Sprintf("wantt=%v, n=%v, nshfts=%v, ktop=%v, kbot=%v, ldh=%v, kacc22=%v",
					wantt, n, nshfts, ktop, kbot, ldh, kacc22)
				if !equalApprox(n, n, h, ldh, test.HWant, 1e-13) {
					t.Errorf("Case %v: unexpected matrix H\nh    =%v\nhwant=%v", prefix, h, test.HWant)
				}
				if !equalApprox(n, n, z, ldh, test.ZWant, 1e-13) {
					t.Errorf("Case %v: unexpected matrix Z\nz    =%v\nzwant=%v", prefix, z, test.ZWant)
				}
			}
		}
	}
}
