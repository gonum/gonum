package umeyama

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestUmeyama(t *testing.T) {
	tol := 1e-10

	testCases := []struct {
		name      string
		src       *mat.Dense
		dst       *mat.Dense
		wantScale float64
		wantRot   *mat.Dense
		wantTrans *mat.VecDense
		tolerance float64
	}{
		{
			name:      "2D_case_from_paper",
			src:       mat.NewDense(2, 3, []float64{0, 1, 0, 0, 0, 2}),
			dst:       mat.NewDense(2, 3, []float64{0, -1, 0, 0, 0, 2}),
			wantScale: 0.7211102550927978,
			wantRot:   mat.NewDense(2, 2, []float64{0.8320502943378437, 0.554700196225229, -0.554700196225229, 0.8320502943378436}),
			wantTrans: mat.NewVecDense(2, []float64{-0.7999999999999998, 0.4}),
			tolerance: tol,
		},
		{
			name:      "2D_identity",
			src:       mat.NewDense(2, 3, []float64{0, 1, 2, 0, 1, 2}),
			dst:       mat.NewDense(2, 3, []float64{0, 1, 2, 0, 1, 2}),
			wantScale: 1.0,
			wantRot:   mat.NewDense(2, 2, []float64{1.0, 1.2116883882008518e-16, 1.2116883882008518e-16, 1.0}),
			wantTrans: mat.NewVecDense(2, []float64{-2.220446049250313e-16, -2.220446049250313e-16}),
			tolerance: tol,
		},
		{
			name:      "2D_rotation_90deg",
			src:       mat.NewDense(2, 3, []float64{0, 1, 1, 0, 0, 1}),
			dst:       mat.NewDense(2, 3, []float64{0, 0, -1, 0, 1, 1}),
			wantScale: 0.9999999999999999,
			wantRot:   mat.NewDense(2, 2, []float64{-5.613347976343136e-17, -0.9999999999999998, 0.9999999999999998, -2.9040269150165053e-16}),
			wantTrans: mat.NewVecDense(2, []float64{-5.551115123125783e-17, 3.3306690738754696e-16}),
			tolerance: tol,
		},
		{
			name:      "2D_scale_2x",
			src:       mat.NewDense(2, 3, []float64{0, 1, 2, 0, 1, 2}),
			dst:       mat.NewDense(2, 3, []float64{0, 2, 4, 0, 2, 4}),
			wantScale: 2.0,
			wantRot:   mat.NewDense(2, 2, []float64{1.0, 1.2116883882008518e-16, 1.2116883882008518e-16, 1.0}),
			wantTrans: mat.NewVecDense(2, []float64{-4.440892098500626e-16, -4.440892098500626e-16}),
			tolerance: tol,
		},
		{
			name:      "2D_translation",
			src:       mat.NewDense(2, 3, []float64{0, 1, 2, 0, 1, 2}),
			dst:       mat.NewDense(2, 3, []float64{3, 4, 5, 3, 4, 5}),
			wantScale: 1.0,
			wantRot:   mat.NewDense(2, 2, []float64{1.0, 1.2116883882008518e-16, 1.2116883882008518e-16, 1.0}),
			wantTrans: mat.NewVecDense(2, []float64{3.0, 3.0}),
			tolerance: tol,
		},
		{
			name:      "3D_case",
			src:       mat.NewDense(3, 3, []float64{0, 1, 2, 0, 0, 5, 1, 3, 8}),
			dst:       mat.NewDense(3, 3, []float64{1, 0, 1, 2, 1, 7, 4, 6, 11}),
			wantScale: 1.0205423989219404,
			wantRot:   mat.NewDense(3, 3, []float64{0.5699453289954445, 0.5900767342443888, -0.5718144538744644, -0.5030534073108366, 0.8008235178014148, 0.324990711758234, 0.6496919203355019, 0.10242627123762431, 0.7532657350571071}),
			wantTrans: mat.NewVecDense(3, []float64{1.4155929948174535, 1.1579295387121973, 3.0877861136679647}),
			tolerance: tol,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scale, rotation, translation, err := Umeyama(tc.src, tc.dst)
			if err != nil {
				t.Fatalf("Umeyama returned error: %v", err)
			}

			// Check scale
			if math.Abs(scale-tc.wantScale) > tc.tolerance {
				t.Errorf("Scale = %v, want %v", scale, tc.wantScale)
			}

			// Check rotation
			d, _ := tc.wantRot.Dims()
			for i := 0; i < d; i++ {
				for j := 0; j < d; j++ {
					if math.Abs(rotation.At(i, j)-tc.wantRot.At(i, j)) > tc.tolerance {
						t.Errorf("Rotation[%d,%d] = %v, want %v", i, j, rotation.At(i, j), tc.wantRot.At(i, j))
					}
				}
			}

			// Check translation
			for i := 0; i < d; i++ {
				if math.Abs(translation.AtVec(i)-tc.wantTrans.AtVec(i)) > tc.tolerance {
					t.Errorf("Translation[%d] = %v, want %v", i, translation.AtVec(i), tc.wantTrans.AtVec(i))
				}
			}
		})
	}
}
