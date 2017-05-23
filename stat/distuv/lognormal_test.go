// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import "testing"

func TestLognormal(t *testing.T) {
	for i, dist := range []LogNormal{
		{
			Mu:    0.1,
			Sigma: 0.3,
		},
		{
			Mu:    0.01,
			Sigma: 0.01,
		},
		{
			Mu:    2,
			Sigma: 0.01,
		},
	} {
		testFullDist(t, dist, i, true)
	}
}
