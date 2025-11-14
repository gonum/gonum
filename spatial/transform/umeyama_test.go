// Copyright ©2025 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transform

import (
	"testing"

	"gonum.org/v1/gonum/floats/scalar"
	"gonum.org/v1/gonum/mat"
)

var umeyamaTests = []struct {
	name      string
	from      *mat.Dense
	to        *mat.Dense
	wantScale float64
	wantRot   *mat.Dense
	wantTrans *mat.VecDense
}{
	{
		name: "2D_case_from_paper",
		from: mat.NewDense(3, 2, []float64{
			0, 0,
			1, 0,
			0, 2,
		}),
		to: mat.NewDense(3, 2, []float64{
			0, 0,
			-1, 0,
			0, 2,
		}),
		wantScale: 0.7211102550927978,
		wantRot: mat.NewDense(2, 2, []float64{
			0.8320502943378437, 0.554700196225229,
			-0.554700196225229, 0.8320502943378436,
		}),
		wantTrans: mat.NewVecDense(2, []float64{
			-0.8,
			0.4,
		}),
	},
	{
		name: "2D_identity",
		from: mat.NewDense(3, 2, []float64{
			0, 0,
			1, 1,
			2, 2,
		}),
		to: mat.NewDense(3, 2, []float64{
			0, 0,
			1, 1,
			2, 2,
		}),
		wantScale: 1,
		wantRot: mat.NewDense(2, 2, []float64{
			1, 0,
			0, 1,
		}),
		wantTrans: mat.NewVecDense(2, []float64{
			0,
			0,
		}),
	},
	{
		name: "2D_rotation_90deg",
		from: mat.NewDense(3, 2, []float64{
			0, 0,
			1, 0,
			1, 1,
		}),
		to: mat.NewDense(3, 2, []float64{
			0, 0,
			0, 1,
			-1, 1,
		}),
		wantScale: 1,
		wantRot: mat.NewDense(2, 2, []float64{
			0, -1,
			1, -0,
		}),
		wantTrans: mat.NewVecDense(2, []float64{
			0,
			0,
		}),
	},
	{
		name: "2D_scale_2x",
		from: mat.NewDense(3, 2, []float64{
			0, 0,
			1, 1,
			2, 2,
		}),
		to: mat.NewDense(3, 2, []float64{
			0, 0,
			2, 2,
			4, 4,
		}),
		wantScale: 2,
		wantRot: mat.NewDense(2, 2, []float64{
			1, 0,
			0, 1,
		}),
		wantTrans: mat.NewVecDense(2, []float64{
			0,
			0,
		}),
	},
	{
		name: "2D_translation",
		from: mat.NewDense(3, 2, []float64{
			0, 0,
			1, 1,
			2, 2,
		}),
		to: mat.NewDense(3, 2, []float64{
			3, 3,
			4, 4,
			5, 5,
		}),
		wantScale: 1,
		wantRot: mat.NewDense(2, 2, []float64{
			1, 0,
			0, 1,
		}),
		wantTrans: mat.NewVecDense(2, []float64{
			3,
			3,
		}),
	},
	{
		name: "3D_case",
		from: mat.NewDense(3, 3, []float64{
			0, 0, 1,
			1, 0, 3,
			2, 5, 8,
		}),
		to: mat.NewDense(3, 3, []float64{
			1, 2, 4,
			0, 1, 6,
			1, 7, 11,
		}),
		wantScale: 1.0205423989219404,
		wantRot: mat.NewDense(3, 3, []float64{
			0.5699453289954445, 0.5900767342443888, -0.5718144538744644,
			-0.5030534073108366, 0.8008235178014148, 0.324990711758234,
			0.6496919203355019, 0.10242627123762431, 0.7532657350571071,
		}),
		wantTrans: mat.NewVecDense(3, []float64{
			1.4155929948174535,
			1.1579295387121973,
			3.0877861136679647,
		}),
	},
}

func TestUmeyama(t *testing.T) {
	tol := 1e-14

	for _, test := range umeyamaTests {
		t.Run(test.name, func(t *testing.T) {
			scale, rotation, translation, err := Umeyama(test.from, test.to, -1)
			if err != nil {
				t.Fatalf("Umeyama returned error: %v", err)
			}

			// Check scale
			if !scalar.EqualWithinAbs(scale, test.wantScale, tol) {
				t.Errorf("Scale = %v, want %v", scale, test.wantScale)
			}

			// Check rotation
			var rDiff mat.Dense
			rDiff.Sub(rotation, test.wantRot)
			diff := rDiff.Norm(1)

			if diff > tol {
				t.Errorf("unexpected rotation matrix, |R_got-R_want| = %v", diff)
			}

			// Check translation
			var tDiff mat.VecDense
			tDiff.SubVec(translation, test.wantTrans)
			diff = tDiff.Norm(1)

			if diff > tol {
				t.Errorf("unexpected translation vector, |t_got-t_want| = %v", diff)
			}
		})
	}
}
