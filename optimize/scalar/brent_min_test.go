// Copyright ©2025 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scalar

import (
	"math"
	"testing"
)

var brentMinTests = []struct {
	name       string
	f          func(float64) float64
	min, max   float64
	settings   *Settings
	wantX      float64
	wantStatus Status
	distTol    float64 // Tolerance for the distance from wantX
}{
	{
		name: "parabola",
		f: func(x float64) float64 {
			return (x-2)*(x-2) + 1
		},
		min:        0,
		max:        5,
		settings:   nil, // Default settings
		wantX:      2,
		wantStatus: Converged,
		distTol:    1e-7,
	},
	{
		name: "cosine",
		f: func(x float64) float64 {
			return math.Cos(x)
		},
		min:        2,
		max:        4,
		settings:   nil,
		wantX:      math.Pi, // Minimum is at π
		wantStatus: Converged,
		distTol:    1e-7,
	},
	{
		name: "x_cos_x",
		f: func(x float64) float64 {
			return x * math.Cos(x)
		},
		min:        0,
		max:        5,
		settings:   nil,
		wantX:      3.425618459481728, // Approximate solution
		wantStatus: Converged,
		distTol:    1e-7,
	},
	{
		name: "strict_tolerance",
		f: func(x float64) float64 {
			return (x-2)*(x-2) + 1
		},
		min:        0,
		max:        5,
		settings:   &Settings{Tol: 1e-10}, // Stricter tolerance
		wantX:      2,
		wantStatus: Converged,
		distTol:    1e-9, // Expect higher precision
	},
	{
		name: "iteration_limit_check",
		f: func(x float64) float64 {
			return (x-2)*(x-2) + 1
		},
		min: 0,
		max: 5,
		// Set a very low iteration limit to ensure it stops early
		settings:   &Settings{MaxIterations: 2},
		wantX:      2,              // Ignored in check
		wantStatus: IterationLimit, // Should hit limit
		distTol:    10,             // Ignored in check
	},
}

func TestBrentMin(t *testing.T) {
	for _, test := range brentMinTests {
		t.Run(test.name, func(t *testing.T) {
			res, err := BrentMin(test.f, test.min, test.max, test.settings)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if res.Status != test.wantStatus {
				t.Errorf("status mismatch: got %v, want %v", res.Status, test.wantStatus)
			}

			if test.wantStatus == Converged {
				dist := math.Abs(res.X - test.wantX)
				if dist > test.distTol {
					t.Errorf("x mismatch:\n got  %.15f\n want %.15f\n diff %.2e > tol %.2e",
						res.X, test.wantX, dist, test.distTol)
				}
			}
		})
	}
}
