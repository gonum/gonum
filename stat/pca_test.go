// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

func TestPrincipalComponents(t *testing.T) {
	// Threshold for detecting zero variances.
	const epsilon = 1e-15
tests:
	for i, test := range []struct {
		data     mat.Matrix
		weights  []float64
		wantVecs *mat.Dense
		wantVars []float64
		epsilon  float64
	}{
		// Test results verified using R.
		{
			data: mat.NewDense(3, 3, []float64{
				1, 2, 3,
				4, 5, 6,
				7, 8, 9,
			}),
			wantVecs: mat.NewDense(3, 3, []float64{
				0.5773502691896258, 0.8164965809277261, 0,
				0.577350269189626, -0.4082482904638632, -0.7071067811865476,
				0.5773502691896258, -0.4082482904638631, 0.7071067811865475,
			}),
			wantVars: []float64{27, 0, 0},
			epsilon:  1e-12,
		},
		{ // Truncated iris data.
			data: mat.NewDense(10, 4, []float64{
				5.1, 3.5, 1.4, 0.2,
				4.9, 3.0, 1.4, 0.2,
				4.7, 3.2, 1.3, 0.2,
				4.6, 3.1, 1.5, 0.2,
				5.0, 3.6, 1.4, 0.2,
				5.4, 3.9, 1.7, 0.4,
				4.6, 3.4, 1.4, 0.3,
				5.0, 3.4, 1.5, 0.2,
				4.4, 2.9, 1.4, 0.2,
				4.9, 3.1, 1.5, 0.1,
			}),
			wantVecs: mat.NewDense(4, 4, []float64{
				-0.6681110197952722, 0.7064764857539533, -0.14026590216895132, -0.18666578956412125,
				-0.7166344774801547, -0.6427036135482664, -0.135650285905254, 0.23444848208629923,
				-0.164411275166307, 0.11898477441068218, 0.9136367900709548, 0.35224901970831746,
				-0.11415613655453069, -0.2714141920887426, 0.35664028439226514, -0.8866286823515034,
			}),
			wantVars: []float64{0.1665786313282786, 0.02065509475412993, 0.007944620317765855, 0.0019327647109368329},
			epsilon:  1e-12,
		},
		{ // Truncated iris data to form wide matrix.
			data: mat.NewDense(3, 4, []float64{
				5.1, 3.5, 1.4, 0.2,
				4.9, 3.0, 1.4, 0.2,
				4.7, 3.2, 1.3, 0.2,
			}),
			wantVecs: mat.NewDense(4, 3, []float64{
				-0.5705187254552365, -0.7505979435049239, 0.08084520834544455,
				-0.8166537769529318, 0.5615147645527523, -0.032338083338177705,
				-0.08709186238359454, -0.3482870890450082, -0.22636658336724505,
				0, 0, -0.9701425001453315,
			}),
			wantVars: []float64{0.0844692361537822, 0.022197430512884326, 0},
			epsilon:  1e-12,
		},
		{ // Truncated iris data transposed to check for operation on fat input.
			data: mat.NewDense(10, 4, []float64{
				5.1, 3.5, 1.4, 0.2,
				4.9, 3.0, 1.4, 0.2,
				4.7, 3.2, 1.3, 0.2,
				4.6, 3.1, 1.5, 0.2,
				5.0, 3.6, 1.4, 0.2,
				5.4, 3.9, 1.7, 0.4,
				4.6, 3.4, 1.4, 0.3,
				5.0, 3.4, 1.5, 0.2,
				4.4, 2.9, 1.4, 0.2,
				4.9, 3.1, 1.5, 0.1,
			}).T(),
			wantVecs: mat.NewDense(10, 4, []float64{
				-0.3366602459946619, -0.1373634006401213, 0.3465102523547623, -0.10290179303893479,
				-0.31381852053861975, 0.5197145790632827, 0.5567296129086686, -0.15923062170153618,
				-0.30857197637565165, -0.07670930360819002, 0.36159923003337235, 0.3342301027853355,
				-0.29527124351656137, 0.16885455995353074, -0.5056204762881208, 0.32580913261444344,
				-0.3327611073694004, -0.39365834489416474, 0.04900050959307464, 0.46812879383236555,
				-0.34445484362044815, -0.2985206914561878, -0.1009714701361799, -0.16803618186050803,
				-0.2986246350957691, -0.4222037823717799, -0.11838613462182519, -0.580283530375069,
				-0.325911246223126, 0.024366468758217238, -0.12082035131864265, 0.16756027181337868,
				-0.2814284432361538, 0.240812316260054, -0.24061437569068145, -0.365034616264623,
				-0.31906138507685167, 0.4423912824105986, -0.2906412122303604, 0.027551046870337714,
			}),
			wantVars: []float64{41.8851906634233, 0.07762619213464989, 0.010516477775373585, 0},
			epsilon:  1e-12,
		},
		{ // Truncated iris data unitary weights.
			data: mat.NewDense(10, 4, []float64{
				5.1, 3.5, 1.4, 0.2,
				4.9, 3.0, 1.4, 0.2,
				4.7, 3.2, 1.3, 0.2,
				4.6, 3.1, 1.5, 0.2,
				5.0, 3.6, 1.4, 0.2,
				5.4, 3.9, 1.7, 0.4,
				4.6, 3.4, 1.4, 0.3,
				5.0, 3.4, 1.5, 0.2,
				4.4, 2.9, 1.4, 0.2,
				4.9, 3.1, 1.5, 0.1,
			}),
			weights: []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			wantVecs: mat.NewDense(4, 4, []float64{
				-0.6681110197952722, 0.7064764857539533, -0.14026590216895132, -0.18666578956412125,
				-0.7166344774801547, -0.6427036135482664, -0.135650285905254, 0.23444848208629923,
				-0.164411275166307, 0.11898477441068218, 0.9136367900709548, 0.35224901970831746,
				-0.11415613655453069, -0.2714141920887426, 0.35664028439226514, -0.8866286823515034,
			}),
			wantVars: []float64{0.1665786313282786, 0.02065509475412993, 0.007944620317765855, 0.0019327647109368329},
			epsilon:  1e-12,
		},
		{ // Truncated iris data non-unitary weights.
			data: mat.NewDense(10, 4, []float64{
				5.1, 3.5, 1.4, 0.2,
				4.9, 3.0, 1.4, 0.2,
				4.7, 3.2, 1.3, 0.2,
				4.6, 3.1, 1.5, 0.2,
				5.0, 3.6, 1.4, 0.2,
				5.4, 3.9, 1.7, 0.4,
				4.6, 3.4, 1.4, 0.3,
				5.0, 3.4, 1.5, 0.2,
				4.4, 2.9, 1.4, 0.2,
				4.9, 3.1, 1.5, 0.1,
			}),
			weights: []float64{2, 3, 1, 1, 1, 1, 1, 1, 1, 2},
			wantVecs: mat.NewDense(4, 4, []float64{
				-0.618936145422414, 0.763069301531647, 0.124857741232537, 0.138035623677211,
				-0.763958271606519, -0.603881770702898, 0.118267155321333, -0.194184052457746,
				-0.143552119754944, 0.090014599564871, -0.942209377020044, -0.289018426115945,
				-0.112599271966947, -0.212012782487076, -0.287515067921680, 0.927203898682805,
			}),
			wantVars: []float64{0.129621985550623, 0.022417487771598, 0.006454461065715, 0.002495076601075},
			epsilon:  1e-12,
		},
	} {
		var pc PC
		var vecs *mat.Dense
		var vars []float64
		for j := 0; j < 2; j++ {
			ok := pc.PrincipalComponents(test.data, test.weights)
			vecs = pc.VectorsTo(vecs)
			vars = pc.VarsTo(vars)
			if !ok {
				t.Errorf("unexpected SVD failure for test %d use %d", i, j)
				continue tests
			}

			// Find the number of non-zero variances to handle
			// non-uniqueness in SVD result (issue #21).
			nnz := len(vars)
			for k, v := range vars {
				if math.Abs(v) < epsilon {
					nnz = k
					break
				}
			}
			r, c := vecs.Dims()
			if !mat.EqualApprox(vecs.Slice(0, r, 0, nnz), test.wantVecs.Slice(0, r, 0, nnz), test.epsilon) {
				t.Errorf("%d use %d: unexpected PCA result got:\n%v\nwant:\n%v",
					i, j, mat.Formatted(vecs), mat.Formatted(test.wantVecs))
			}
			if !approxEqual(vars, test.wantVars, test.epsilon) {
				t.Errorf("%d use %d: unexpected variance result got:%v, want:%v",
					i, j, vars, test.wantVars)
			}

			// Check that the set of principal vectors is
			// orthonormal by comparing V^T*V to the identity matrix.
			I := mat.NewDiagDense(c, nil)
			for k := 0; k < c; k++ {
				I.SetDiag(k, 1)
			}
			var vv mat.Dense
			vv.Mul(vecs.T(), vecs)
			if !mat.EqualApprox(&vv, I, test.epsilon) {
				t.Errorf("%d use %d: vectors not orthonormal\n%v", i, j, mat.Formatted(I))
			}
		}
	}
}

func approxEqual(a, b []float64, epsilon float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if !floats.EqualWithinAbsOrRel(v, b[i], epsilon, epsilon) {
			return false
		}
	}
	return true
}
