// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"math"
	"testing"
)

func TestTriangle(t *testing.T) {
	for i, dist := range []Triangle{
		{
			A: 0.0,
			C: 0.5,
			B: 1.0,
		},
		{
			A: 0.1,
			C: 0.2,
			B: 0.3,
		},
		{
			A: 1.0,
			C: 1.5,
			B: 2.0,
		},
	} {
		testFullDist(t, dist, i, true)
	}
}

func TestTriangleProb(t *testing.T) {
	pts := []univariateProbPoint{
		{
			loc:     0.5,
			prob:    0,
			cumProb: 0,
			logProb: math.Inf(-1),
		},
		{
			loc:     1,
			prob:    0,
			cumProb: 0,
			logProb: math.Inf(-1),
		},
		{
			loc:     2,
			prob:    1.0,
			cumProb: 0.5,
			logProb: 0,
		},
		{
			loc:     3,
			prob:    0,
			cumProb: 1,
			logProb: math.Inf(-1),
		},
	}
	testDistributionProbs(t, Triangle{A: 1, C: 2, B: 3}, "Standard 1,2,3 Triangle", pts)
}
