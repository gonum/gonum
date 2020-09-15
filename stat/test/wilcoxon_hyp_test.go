// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestWilcoxonSignedRankTest(t *testing.T) {
	tests := []struct {
		x           []float64
		y           []float64
		exactPValue bool
		testStatistic float64
		pValue float64
	}{
		{
			x:           []float64{1.83, 0.50, 1.62, 2.48, 1.68, 1.88, 1.55, 3.06, 1.30},
			y:           []float64{0.878, 0.647, 0.598, 2.050, 1.060, 1.290, 1.060, 3.140, 1.290},
			exactPValue: true,
			pValue:        0.0390625,
			testStatistic: 40,
		},
		{
			x:           []float64{2.0, 3.6, 2.6, 2.6, 7.3, 3.4, 14.9, 6.6, 2.3, 2.0, 6.8, 8.5},
			y:           []float64{3.5, 5.7, 2.9, 2.4, 9.9, 3.3, 16.7, 6.0, 3.8, 4.0, 9.1, 20.9},
			exactPValue: false,
			pValue:        0.010757133119789652,
			testStatistic: 7,
		},
	}
	for index, tt := range tests {
		t.Run(strconv.Itoa(index), func(t *testing.T) {
			gotPValue, gotTestStatistic  := WilcoxonSignedRankTest(tt.x, tt.y, tt.exactPValue)
			if !strings.EqualFold(fmt.Sprintf("%0.5f", gotPValue), fmt.Sprintf("%0.5f", tt.pValue)) || !strings.EqualFold(fmt.Sprintf("%0.5f", gotTestStatistic), fmt.Sprintf("%0.5f", tt.testStatistic)){
				t.Errorf("WilcoxonSignedRankTest() = %v, %v want %v, %v", gotPValue, gotTestStatistic, tt.pValue, tt.testStatistic)
			}
		})
	}
}
