// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmv

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/spatial/r1"
)

func TestUniformEntropy(t *testing.T) {
	for _, test := range []struct {
		Uniform *Uniform
		Entropy float64
	}{
		{
			NewUniform([]r1.Interval{{Min: 0, Max: 1}, {Min: 0, Max: 1}}, nil),
			0,
		},
		{
			NewUniform([]r1.Interval{{Min: -1, Max: 3}, {Min: 2, Max: 8}, {Min: -5, Max: -3}}, nil),
			math.Log(48),
		},
	} {
		ent := test.Uniform.Entropy()
		if math.Abs(ent-test.Entropy) > 1e-14 {
			t.Errorf("Entropy mismatch. Got %v, want %v", ent, test.Entropy)
		}
	}
}
