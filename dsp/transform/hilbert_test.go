// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transform

import (
	"fmt"
	"testing"

	"gonum.org/v1/gonum/cmplxs"
)

func TestHilbertAnalytic(t *testing.T) {
	const tol = 1e-10
	testCases := []struct {
		in   []float64
		want []complex128
	}{
		{[]float64{0, 0, 0, 0}, []complex128{0, 0, 0, 0}},
		{[]float64{1, 2, 3, 4}, []complex128{1 + 1i, 2 - 1i, 3 - 1i, 4 + 1i}},
		{
			[]float64{1, 2, 3, 4, 5},
			[]complex128{
				1 + 1.7013016167i,
				2 - 1.3763819204i,
				3 - 0.6498393924i,
				4 - 1.3763819204i,
				5 + 1.7013016167i,
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("known Hilbert %d", i), func(t *testing.T) {
			h := NewHilbert(len(tc.in))
			if h.Len() != len(tc.in) {
				t.Errorf("unexpected Hilbert transform length: got:%d, want:%d", h.Len(), len(tc.in))
			}

			dst := make([]complex128, len(tc.in))
			analytic := h.AnalyticSignal(tc.in, dst)
			if !cmplxs.EqualApprox(tc.want, analytic, tol) {
				t.Errorf("expected Hilbert transform result %v, got %v", tc.want, analytic)
			}
		})
	}
}
