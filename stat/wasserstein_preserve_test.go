// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat_test

import (
	"math"
	"slices"
	"testing"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

func TestWassersteinDistancePreservesInputs(t *testing.T) {
	for _, nd := range []bool{false, true} {
		name := "1D"
		if nd {
			name = "ND"
		}
		t.Run(name, func(t *testing.T) {
			for _, test := range []struct {
				name  string
				setup func() (p, q, pw, qw []float64)
				want  float64
			}{
				{
					name: "count weights",
					setup: func() (p, q, pw, qw []float64) {
						return []float64{0, 10}, []float64{0, 20}, []float64{2, 2}, []float64{3, 3}
					},
					want: 5, // Half the mass moves from 10 to 20.
				},
				{
					name: "shared weights",
					setup: func() (p, q, pw, qw []float64) {
						weights := []float64{2, 2}
						return []float64{0, 10}, []float64{0, 20}, weights, weights
					},
					want: 5,
				},
				{
					name: "overlapping weights",
					setup: func() (p, q, pw, qw []float64) {
						weights := []float64{1, 1, 1}
						return []float64{0, 10}, []float64{0, 20}, weights[:2], weights[1:]
					},
					want: 5,
				},
				{
					name: "points and weights share storage",
					setup: func() (p, q, pw, qw []float64) {
						p, q = []float64{1, 3}, []float64{2, 4}
						return p, q, p, q
					},
					want: 1, // CDF differences over [1,2], [2,3], [3,4] sum to 1/4+1/12+2/3.
				},
				{
					name: "nonuniform weights",
					setup: func() (p, q, pw, qw []float64) {
						return []float64{0, 10}, []float64{0, 20}, []float64{3, 1}, []float64{2, 2}
					},
					want: 7.5, // One quarter moves 0 to 20 and one quarter moves 10 to 20.
				},
				{
					name: "nil p weights",
					setup: func() (p, q, pw, qw []float64) {
						return []float64{0, 10}, []float64{0, 20}, nil, []float64{3, 3}
					},
					want: 5,
				},
				{
					name: "nil q weights",
					setup: func() (p, q, pw, qw []float64) {
						return []float64{0, 10}, []float64{0, 20}, []float64{2, 2}, nil
					},
					want: 5,
				},
				{
					name: "single points",
					setup: func() (p, q, pw, qw []float64) {
						return []float64{1}, []float64{4}, []float64{2}, []float64{3}
					},
					want: 3,
				},
				{
					name: "identical distributions",
					setup: func() (p, q, pw, qw []float64) {
						return []float64{0, 10}, []float64{0, 10}, []float64{2, 2}, []float64{2, 2}
					},
				},
			} {
				t.Run(test.name, func(t *testing.T) {
					p, q, pw, qw := test.setup()
					pOrig, qOrig := slices.Clone(p), slices.Clone(q)
					pwOrig, qwOrig := slices.Clone(pw), slices.Clone(qw)
					variance := stat.Variance(p, pw)
					var got float64
					if nd {
						var err error
						got, err = stat.WassersteinDistanceND(mat.NewDense(len(p), 1, p), mat.NewDense(len(q), 1, q), pw, qw, 1e-12)
						if err != nil {
							t.Fatal(err)
						}
					} else {
						got = stat.WassersteinDistance(p, q, pw, qw)
					}
					if math.IsNaN(got) || math.Abs(got-test.want) > 1e-10 {
						t.Errorf("distance=%g, want %g", got, test.want)
					}
					if !slices.Equal(p, pOrig) || !slices.Equal(q, qOrig) || !slices.Equal(pw, pwOrig) || !slices.Equal(qw, qwOrig) {
						t.Errorf("inputs changed: p=%v q=%v pw=%v qw=%v; want %v %v %v %v", p, q, pw, qw, pOrig, qOrig, pwOrig, qwOrig)
					}
					if test.name == "count weights" {
						// These are sample counts: the unbiased variance is 100/(4-1).
						if variance != 100.0/3 || stat.Variance(p, pw) != variance {
							t.Errorf("variance changed after reusing count weights: before=%g after=%g", variance, stat.Variance(p, pw))
						}
					}
				})
			}
		})
	}
}

func TestWassersteinDistancePreservesWeightsOnPanic(t *testing.T) {
	for _, nd := range []bool{false, true} {
		name := "1D"
		if nd {
			name = "ND"
		}
		t.Run(name, func(t *testing.T) {
			for _, badWeight := range []float64{0, -1} {
				p, q := []float64{0, 10}, []float64{0, 20}
				pw, qw := []float64{2, 2}, []float64{1, badWeight}
				func() {
					defer func() {
						if got := recover(); got != "stat: all weights must be positive" {
							t.Errorf("unexpected panic: %v", got)
						}
					}()
					if nd {
						stat.WassersteinDistanceND(mat.NewDense(2, 1, p), mat.NewDense(2, 1, q), pw, qw, 1e-12)
					} else {
						stat.WassersteinDistance(p, q, pw, qw)
					}
				}()
				if !slices.Equal(pw, []float64{2, 2}) || !slices.Equal(qw, []float64{1, badWeight}) {
					t.Errorf("weights changed before panic: p=%v q=%v", pw, qw)
				}
			}
		})
	}
}
