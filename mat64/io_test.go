// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import "testing"

func TestDenseRW(t *testing.T) {
	for i, test := range []*Dense{
		NewDense(0, 0, []float64{}),
		NewDense(2, 2, []float64{1, 2, 3, 4}),
		NewDense(2, 3, []float64{1, 2, 3, 4, 5, 6}),
		NewDense(3, 2, []float64{1, 2, 3, 4, 5, 6}),
		NewDense(3, 3, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}),
		NewDense(3, 3, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}).View(0, 0, 2, 2).(*Dense),
		NewDense(3, 3, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}).View(1, 1, 2, 2).(*Dense),
		NewDense(3, 3, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}).View(0, 1, 3, 2).(*Dense),
	} {
		buf, err := test.MarshalBinary()
		if err != nil {
			t.Errorf("error encoding test #%d: %v\n", i, err)
		}

		nrows, ncols := test.Dims()
		sz := nrows*ncols*sizeFloat64 + 2*sizeInt64
		if len(buf) != sz {
			t.Errorf("encoded size test #%d: want=%d got=%d\n", i, sz, len(buf))
		}

		var got Dense
		err = got.UnmarshalBinary(buf)
		if err != nil {
			t.Errorf("error decoding test #%d: %v\n", i, err)
		}

		if !Equal(&got, test) {
			t.Errorf("r/w test #%d failed\nwant=%#v\n got=%#v\n", i, test, &got)
		}
	}
}
