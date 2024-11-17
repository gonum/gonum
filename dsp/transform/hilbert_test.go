// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transform

import (
	"testing"

	"gonum.org/v1/gonum/cmplxs"
)

var hilbertAnalyticSignalTests = []struct {
	name string
	in   []float64
	want []complex128
}{
	{
		name: "zeros",
		in:   []float64{0, 0, 0, 0},
		want: []complex128{0, 0, 0, 0}},
	{
		name: "whole_components",
		in:   []float64{1, 2, 3, 4},
		want: []complex128{1 + 1i, 2 - 1i, 3 - 1i, 4 + 1i},
	},
	{
		name: "irrational_imaginary_components",
		in:   []float64{1, 2, 3, 4, 5},
		want: []complex128{
			1 + 1.7013016167i,
			2 - 1.3763819204i,
			3 - 0.6498393924i,
			4 - 1.3763819204i,
			5 + 1.7013016167i,
		},
	},
}

func TestHilbertAnalytic(t *testing.T) {
	const tol = 1e-10

	for _, test := range hilbertAnalyticSignalTests {
		t.Run(test.name, func(t *testing.T) {
			h := NewHilbert(len(test.in))
			if h.Len() != len(test.in) {
				t.Errorf("unexpected Hilbert transform length: got:%d, want:%d", h.Len(), len(test.in))
			}

			dst := make([]complex128, len(test.in))
			got := h.AnalyticSignal(dst, test.in)
			if !cmplxs.EqualApprox(got, test.want, tol) {
				t.Errorf("unexpected Hilbert transform result:\ngot: %v\nwant:%v", got, test.want)
			}
		})
	}
}
