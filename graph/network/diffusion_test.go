// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package network

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/graph/simple"
)

var heatDiffusionTests = []struct {
	g []set
	h map[int64]float64
	t float64

	wantTol float64
	want    map[bool]map[int64]float64
}{
	{
		g: grid(5),
		h: map[int64]float64{0: 1},
		t: 0.1,

		wantTol: 1e-9,
		want: map[bool]map[int64]float64{
			false: {
				A: 0.826684055, B: 0.078548060, C: 0.003858840, D: 0.000127487, E: 0.000003233,
				F: 0.078548060, G: 0.007463308, H: 0.000366651, I: 0.000012113, J: 0.000000307,
				K: 0.003858840, L: 0.000366651, M: 0.000018012, N: 0.000000595, O: 0.000000015,
				P: 0.000127487, Q: 0.000012113, R: 0.000000595, S: 0.000000020, T: 0.000000000,
				U: 0.000003233, V: 0.000000307, W: 0.000000015, X: 0.000000000, Y: 0.000000000,
			},
			true: {
				A: 0.9063462486, B: 0.0369774705, C: 0.0006161414, D: 0.0000068453, E: 0.0000000699,
				F: 0.0369774705, G: 0.0010670895, H: 0.0000148186, I: 0.0000001420, J: 0.0000000014,
				K: 0.0006161414, L: 0.0000148186, M: 0.0000001852, N: 0.0000000016, O: 0.0000000000,
				P: 0.0000068453, Q: 0.0000001420, R: 0.0000000016, S: 0.0000000000, T: 0.0000000000,
				U: 0.0000000699, V: 0.0000000014, W: 0.0000000000, X: 0.0000000000, Y: 0.0000000000,
			},
		},
	},
	{
		g: grid(5),
		h: map[int64]float64{0: 1},
		t: 1,

		wantTol: 1e-9,
		want: map[bool]map[int64]float64{
			false: {
				A: 0.2743435076, B: 0.1615920872, C: 0.0639346641, D: 0.0188054933, E: 0.0051023569,
				F: 0.1615920872, G: 0.0951799548, H: 0.0376583937, I: 0.0110766934, J: 0.0030053582,
				K: 0.0639346641, L: 0.0376583937, M: 0.0148997194, N: 0.0043825455, O: 0.0011890840,
				P: 0.0188054933, Q: 0.0110766934, R: 0.0043825455, S: 0.0012890649, T: 0.0003497525,
				U: 0.0051023569, V: 0.0030053582, W: 0.0011890840, X: 0.0003497525, Y: 0.0000948958,
			},
			true: {
				A: 0.4323917545, B: 0.1660487336, C: 0.0270298904, D: 0.0029720194, E: 0.0003007247,
				F: 0.1660487336, G: 0.0463974679, H: 0.0063556078, I: 0.0006056850, J: 0.0000589574,
				K: 0.0270298904, L: 0.0063556078, M: 0.0007860810, N: 0.0000691647, O: 0.0000065586,
				P: 0.0029720194, Q: 0.0006056850, R: 0.0000691647, S: 0.0000057466, T: 0.0000005475,
				U: 0.0003007247, V: 0.0000589574, W: 0.0000065586, X: 0.0000005475, Y: 0.0000000555,
			},
		},
	},
	{
		g: grid(5),
		h: map[int64]float64{0: 1},
		t: 10,

		wantTol: 1e-9,
		want: map[bool]map[int64]float64{
			false: {
				A: 0.0432408511, B: 0.0425986522, C: 0.0415977802, D: 0.0405588482, E: 0.0399403788,
				F: 0.0425986522, G: 0.0420083007, H: 0.0409532810, I: 0.0399982373, J: 0.0393463013,
				K: 0.0415977802, L: 0.0409532810, M: 0.0400339958, N: 0.0389913353, O: 0.0384232854,
				P: 0.0405588482, Q: 0.0399982373, R: 0.0389913353, S: 0.0380844049, T: 0.0374622025,
				U: 0.0399403788, V: 0.0393463013, W: 0.0384232854, X: 0.0374622025, Y: 0.0368918429,
			},
			true: {
				A: 0.0532814862, B: 0.0594280160, C: 0.0462076361, D: 0.0330529557, E: 0.0211688130,
				F: 0.0594280160, G: 0.0612529898, H: 0.0462850376, I: 0.0319891593, J: 0.0213123519,
				K: 0.0462076361, L: 0.0462850376, M: 0.0340410963, N: 0.0229646704, O: 0.0152763556,
				P: 0.0330529557, Q: 0.0319891593, R: 0.0229646704, S: 0.0153031853, T: 0.0103681461,
				U: 0.0211688130, V: 0.0213123519, W: 0.0152763556, X: 0.0103681461, Y: 0.0068893147,
			},
		},
	},
	{
		g: grid(5),
		h: func() map[int64]float64 {
			m := make(map[int64]float64, 25)
			for i := int64(A); i <= Y; i++ {
				m[i] = 1
			}
			return m
		}(),
		t: 0.01, // FIXME(kortschak): Low t used due to instability in mat.Exp.

		wantTol: 1e-2, // FIXME(kortschak): High tolerance used due to instability in mat.Exp.
		want: map[bool]map[int64]float64{
			false: {
				A: 1, B: 1, C: 1, D: 1, E: 1,
				F: 1, G: 1, H: 1, I: 1, J: 1,
				K: 1, L: 1, M: 1, N: 1, O: 1,
				P: 1, Q: 1, R: 1, S: 1, T: 1,
				U: 1, V: 1, W: 1, X: 1, Y: 1,
			},
			true: {
				A: 1, B: 1, C: 1, D: 1, E: 1,
				F: 1, G: 1, H: 1, I: 1, J: 1,
				K: 1, L: 1, M: 1, N: 1, O: 1,
				P: 1, Q: 1, R: 1, S: 1, T: 1,
				U: 1, V: 1, W: 1, X: 1, Y: 1,
			},
		},
	},
}

func grid(d int) []set {
	dim := int64(d)
	s := make([]set, dim*dim)
	for i := range s {
		s[i] = make(set)
	}
	for i := int64(0); i < dim*dim; i++ {
		if i%dim != 0 {
			s[i][i-1] = struct{}{}
		}
		if i/dim != 0 {
			s[i][i-dim] = struct{}{}
		}
	}
	return s
}

func TestHeatDiffusion(t *testing.T) {
	for i, test := range heatDiffusionTests {
		g := simple.NewUndirectedGraph()
		for u, e := range test.g {
			// Add nodes that are not defined by an edge.
			if !g.Has(simple.Node(u)) {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
			}
		}
		for _, normalize := range []bool{false, true} {
			var wantTemp float64
			h := make(map[int64]float64)
			for k, v := range test.h {
				h[k] = v
				wantTemp += v
			}
			got := HeatDiffusion(g, h, test.t, normalize)
			prec := 1 - int(math.Log10(test.wantTol))
			for n := range test.g {
				if !floats.EqualWithinAbsOrRel(got[int64(n)], test.want[normalize][int64(n)], test.wantTol, test.wantTol) {
					t.Errorf("unexpected HeatDiffusion result for test %d with normalize=%t:\ngot: %v\nwant:%v",
						i, normalize, orderedFloats(got, prec), orderedFloats(test.want[normalize], prec))
					break
				}
			}

			if normalize {
				continue
			}

			var gotTemp float64
			for _, v := range got {
				gotTemp += v
			}
			gotTemp /= float64(len(got))
			wantTemp /= float64(len(got))
			if !floats.EqualWithinAbsOrRel(gotTemp, wantTemp, test.wantTol, test.wantTol) {
				t.Errorf("unexpected total heat for test %d with normalize=%t: got:%v want:%v",
					i, normalize, gotTemp, wantTemp)
			}
		}
	}
}
