// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transform

import (
	"fmt"
	"math"
	"slices"
	"testing"
)

func TestHilbertAnalytic(t *testing.T) {
	testCases := []struct {
		input []float64
	}{
		{[]float64{0, 0, 0, 0}},
		{[]float64{1, 2, 3, 4}},
		{[]float64{1, 2, 3, 4, 5}},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("tc%d", i), func(t *testing.T) {
			transformer := NewHilbert(len(tc.input))
			if transformer.Len() != len(tc.input) {
				t.Errorf("expected Hilbert transform length to be %d, got %d", len(tc.input), transformer.Len())
			}

			dst := make([]complex128, len(tc.input))
			analytic := transformer.AnalyticSignal(tc.input, dst)
			result := make([]float64, len(tc.input))
			for i, c := range analytic {
				result[i] = math.Round(real(c))
			}
			if !slices.Equal(tc.input, result) {
				t.Errorf("expected Hilbert transform result %v, got %v", tc.input, result)
			}
		})
	}
}
