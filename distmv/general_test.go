// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmv

import (
	"math"
	"testing"
)

type prober interface {
	Prob(x []float64) float64
	LogProb(x []float64) float64
}

type probCase struct {
	dist    prober
	loc     []float64
	logProb float64
}

func testProbability(t *testing.T, cases []probCase) {
	for _, test := range cases {
		logProb := test.dist.LogProb(test.loc)
		if math.Abs(logProb-test.logProb) > 1e-14 {
			t.Errorf("LogProb mismatch: want: %v, got: %v", test.logProb, logProb)
		}
		prob := test.dist.Prob(test.loc)
		if math.Abs(prob-math.Exp(test.logProb)) > 1e-14 {
			t.Errorf("Prob mismatch: want: %v, got: %v", math.Exp(test.logProb), prob)
		}
	}
}
