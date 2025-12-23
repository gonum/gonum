// Copyright ©2025 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scalar

import (
	"math"
	"testing"
)

func TestBrentMin(t *testing.T) {
	tests := []struct {
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
			wantX:      2.0,
			wantStatus: Converged,
			distTol:    1e-7,
		},
		{
			name: "Cosine: cos(x) in [2, 4]",
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
			name: "Function: x * cos(x) in [0, 5]",
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
			name: "Strict Tolerance: (x-2)^2 + 1",
			f: func(x float64) float64 {
				return (x-2)*(x-2) + 1
			},
			min:        0,
			max:        5,
			settings:   &Settings{Tol: 1e-10}, // Stricter tolerance
			wantX:      2.0,
			wantStatus: Converged,
			distTol:    1e-9, // Expect higher precision
		},
		{
			name: "Iteration Limit Check",
			f: func(x float64) float64 {
				return (x-2)*(x-2) + 1
			},
			min: 0,
			max: 5,
			// Set a very low iteration limit to ensure it stops early
			settings:   &Settings{MaxIterations: 2},
			wantX:      2.0,            // Ignored in check
			wantStatus: IterationLimit, // Should hit limit
			distTol:    10.0,           // Ignored in check
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := BrentMin(tc.f, tc.min, tc.max, tc.settings)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if res.Status != tc.wantStatus {
				t.Errorf("status mismatch: got %v, want %v", res.Status, tc.wantStatus)
			}

			if tc.wantStatus == Converged {
				dist := math.Abs(res.X - tc.wantX)
				if dist > tc.distTol {
					t.Errorf("x mismatch:\n got  %.15f\n want %.15f\n diff %.2e > tol %.2e",
						res.X, tc.wantX, dist, tc.distTol)
				}
			}
		})
	}
}
