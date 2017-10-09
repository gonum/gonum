// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import "testing"

func TestBinomial(t *testing.T) {
	for i, dist := range []Binomial{
		{
			P: 0.5,
			N: 3,
		},
		{
			P: 0.9,
			N: 3,
		},
		{
			P: 0.2,
			N: 3,
		},
		{
			P: 0.5,
			N: 5,
		},
		{
			P: 0.9,
			N: 5,
		},
		{
			P: 0.2,
			N: 5,
		},
		{
			P: 0.5,
			N: 10,
		},
		{
			P: 0.9,
			N: 10,
		},
		{
			P: 0.2,
			N: 10,
		},
	} {
		testFullDist(t, dist, i, false)
	}
}
