// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distuv

import (
	"testing"

	"gonum.org/v1/gonum/floats"
)

func TestPoisson(t *testing.T) {
	const tol = 1e-16
	for i, tt := range []struct {
		k    int
		l    float64
		want float64
	}{
		{0, 1, 3.678794411714423e-01},
		{1, 1, 3.678794411714423e-01},
		{2, 1, 1.839397205857211e-01},
		{3, 1, 6.131324019524039e-02},
		{4, 1, 1.532831004881010e-02},
		{5, 1, 3.065662009762020e-03},
		{6, 1, 5.109436682936698e-04},
		{7, 1, 7.299195261338139e-05},
		{8, 1, 9.123994076672672e-06},
		{9, 1, 1.013777119630298e-06},

		{0, 2.5, 8.208499862389880e-02},
		{1, 2.5, 2.052124965597470e-01},
		{2, 2.5, 2.565156206996838e-01},
		{3, 2.5, 2.137630172497365e-01},
		{4, 2.5, 1.336018857810853e-01},
		{5, 2.5, 6.680094289054267e-02},
		{6, 2.5, 2.783372620439277e-02},
		{7, 2.5, 9.940616501568845e-03},
		{8, 2.5, 3.106442656740263e-03},
		{9, 2.5, 8.629007379834082e-04},
	} {
		p := Poisson{Lambda: tt.l}
		got := p.Prob(float64(tt.k))
		if !floats.EqualWithinAbs(got, tt.want, tol) {
			t.Errorf("test-%d: got=%e. want=%e\n", i, got, tt.want)
		}
	}
}
