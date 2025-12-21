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
		distTol    float64 // 許容する x の誤差 (wantX との距離)
	}{
		{
			name: "Parabola: (x-2)^2 + 1",
			f: func(x float64) float64 {
				return (x-2)*(x-2) + 1
			},
			min:        0,
			max:        5,
			settings:   nil, // デフォルト設定
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
			wantX:      math.Pi, // 最小値は π
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
			wantX:      3.425618459481728, // 近似解
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
			settings:   &Settings{Tol: 1e-10}, // 厳しめの許容誤差
			wantX:      2.0,
			wantStatus: Converged,
			distTol:    1e-9, // 結果もより高精度であることを期待
		},
		{
			name: "Iteration Limit Check",
			f: func(x float64) float64 {
				return (x-2)*(x-2) + 1
			},
			min: 0,
			max: 5,
			// 反復回数を極端に少なく設定し、エラー終了することを確認
			settings:   &Settings{MaxIterations: 2},
			wantX:      2.0,            // ※チェックしないが型合わせのため
			wantStatus: IterationLimit, // 制限に達するはず
			distTol:    10.0,           // 誤差チェックは無視
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := BrentMin(tc.f, tc.min, tc.max, tc.settings)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// ステータスのチェック
			if res.Status != tc.wantStatus {
				t.Errorf("status mismatch: got %v, want %v", res.Status, tc.wantStatus)
			}

			// 正常収束が期待されるケースでは、値の精度もチェック
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
