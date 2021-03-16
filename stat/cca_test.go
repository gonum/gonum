// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat_test

import (
	"testing"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

func TestCanonicalCorrelations(t *testing.T) {
tests:
	for i, test := range []struct {
		xdata     mat.Matrix
		ydata     mat.Matrix
		weights   []float64
		wantCorrs []float64
		wantpVecs *mat.Dense
		wantqVecs *mat.Dense
		wantphiVs *mat.Dense
		wantpsiVs *mat.Dense
		epsilon   float64
	}{
		// Test results verified using R.
		{ // Truncated iris data, Sepal vs Petal measurements.
			xdata: mat.NewDense(10, 2, []float64{
				5.1, 3.5,
				4.9, 3.0,
				4.7, 3.2,
				4.6, 3.1,
				5.0, 3.6,
				5.4, 3.9,
				4.6, 3.4,
				5.0, 3.4,
				4.4, 2.9,
				4.9, 3.1,
			}),
			ydata: mat.NewDense(10, 2, []float64{
				1.4, 0.2,
				1.4, 0.2,
				1.3, 0.2,
				1.5, 0.2,
				1.4, 0.2,
				1.7, 0.4,
				1.4, 0.3,
				1.5, 0.2,
				1.4, 0.2,
				1.5, 0.1,
			}),
			wantCorrs: []float64{0.7250624174504773, 0.5547679185730191},
			wantpVecs: mat.NewDense(2, 2, []float64{
				0.0765914610875867, 0.9970625597666721,
				0.9970625597666721, -0.0765914610875868,
			}),
			wantqVecs: mat.NewDense(2, 2, []float64{
				0.3075184850910837, 0.9515421069649439,
				0.9515421069649439, -0.3075184850910837,
			}),
			wantphiVs: mat.NewDense(2, 2, []float64{
				-1.9794877596804641, 5.2016325219025124,
				4.5211829944066553, -2.7263663170835697,
			}),
			wantpsiVs: mat.NewDense(2, 2, []float64{
				-0.0613084818030103, 10.8514169865438941,
				12.7209032660734298, -7.6793888180353775,
			}),
			epsilon: 1e-12,
		},
		// Test results compared to those results presented in examples by
		// Koch, Inge. Analysis of multivariate and high-dimensional data.
		// Vol. 32. Cambridge University Press, 2013. ISBN: 9780521887939
		{ // ASA Car Exposition Data of Ramos and Donoho (1983)
			// Displacement, Horsepower, Weight
			xdata: carData.Slice(0, 392, 0, 3),
			// Acceleration, MPG
			ydata:     carData.Slice(0, 392, 3, 5),
			wantCorrs: []float64{0.8782187384352336, 0.6328187219216761},
			wantpVecs: mat.NewDense(3, 2, []float64{
				0.3218296374829181, 0.3947540257657075,
				0.4162807660635797, 0.7573719053303306,
				0.8503740401982725, -0.5201509936144236,
			}),
			wantqVecs: mat.NewDense(2, 2, []float64{
				-0.5161984172278830, -0.8564690269072364,
				-0.8564690269072364, 0.5161984172278830,
			}),
			wantphiVs: mat.NewDense(3, 2, []float64{
				0.0025033152994308, 0.0047795464118615,
				0.0201923608080173, 0.0409150208725958,
				-0.0000247374128745, -0.0026766435161875,
			}),
			wantpsiVs: mat.NewDense(2, 2, []float64{
				-0.1666196759760772, -0.3637393866139658,
				-0.0915512109649727, 0.1077863777929168,
			}),
			epsilon: 1e-12,
		},
		// Test results compared to those results presented in examples by
		// Koch, Inge. Analysis of multivariate and high-dimensional data.
		// Vol. 32. Cambridge University Press, 2013. ISBN: 9780521887939
		{ // Boston Housing Data of Harrison and Rubinfeld (1978)
			// Per capita crime rate by town,
			// Proportion of non-retail business acres per town,
			// Nitric oxide concentration (parts per 10 million),
			// Weighted distances to Boston employment centres,
			// Index of accessibility to radial highways,
			// Pupil-teacher ratio by town, Proportion of blacks by town
			xdata: bostonData.Slice(0, 506, 0, 7),
			// Average number of rooms per dwelling,
			// Proportion of owner-occupied units built prior to 1940,
			// Full-value property-tax rate per $10000,
			// Median value of owner-occupied homes in $1000s
			ydata:     bostonData.Slice(0, 506, 7, 11),
			wantCorrs: []float64{0.9451239443886021, 0.6786622733370654, 0.5714338361583764, 0.2009739704710440},
			wantpVecs: mat.NewDense(7, 4, []float64{
				-0.2574391924541896, -0.015847751662118038, -0.21221699346310258, -0.09457338038947205,
				-0.48365944300184865, -0.3837101908138455, -0.14744483174159395, 0.6597324886718278,
				-0.08007763658732961, -0.34935567428092285, -0.3287336458109394, -0.2862040444334662,
				0.127758636038638, 0.7337427663667616, -0.4851134819036985, 0.22479648659701942,
				-0.6969432006136685, 0.43417487760028844, 0.360287288763638, 0.029066160862628414,
				-0.0990903250057202, -0.05034112154538474, -0.6384330631742202, 0.10223671362182897,
				0.42604599637650303, -0.032333435130815824, 0.22895275160308087, 0.6419232947608798,
			}),
			wantqVecs: mat.NewDense(4, 4, []float64{
				0.018166050236326788, 0.1583489460479047, 0.006672357764289544, -0.9871935400650647,
				-0.23476990459861324, -0.9483314614936598, 0.14624205056313114, -0.1554470767919039,
				-0.9700704038477144, 0.24060717410000537, 0.025183898422704167, 0.020913407435834964,
				0.05930006823184807, 0.13304600030976868, 0.9889057151969495, 0.029116149472076858,
			}),
			wantphiVs: mat.NewDense(7, 4, []float64{
				-0.002746223410819314, -0.009344451350088911, -0.04896439327142919, -0.015496718980582016,
				-0.042856445527953785, 0.024170870211944927, -0.036072347209397136, 0.18389832305881182,
				-1.2248435648802678, -5.603092136472504, -5.809414458379886, -4.792681219042103,
				-0.00436848250946508, 0.34241011649776265, -0.4469961215717922, 0.11501618143536857,
				-0.07415340695219563, 0.11931357949236807, 0.1115518305471455, 0.002163875832307984,
				-0.023327032310162924, -0.1046330818178401, -0.38530459750774165, -0.016092787010290065,
				0.00012930513878583387, -0.0004540746921447011, 0.0030296315865439264, 0.008189547797465318,
			}),
			wantpsiVs: mat.NewDense(4, 4, []float64{
				0.030159336201738367, 0.3002219289647159, -0.08782173775936601, -1.9583226531517122,
				-0.00654831040738931, -0.03922120867162458, 0.011757077620998818, -0.006111306448187141,
				-0.0052075523350125505, 0.004577020045295936, 0.0022762313289591976, 0.0008441873006823151,
				0.0020111735096325924, -0.0037352799829939247, 0.12925780716217938, 0.10377090563297825,
			}),
			epsilon: 1e-12,
		},
	} {
		var cc stat.CC
		var corrs []float64
		var pVecs, qVecs mat.Dense
		var phiVs, psiVs mat.Dense
		for j := 0; j < 2; j++ {
			err := cc.CanonicalCorrelations(test.xdata, test.ydata, test.weights)
			if err != nil {
				t.Errorf("%d use %d: unexpected error: %v", i, j, err)
				continue tests
			}

			corrs = cc.CorrsTo(corrs)
			cc.LeftTo(&pVecs, true)
			cc.RightTo(&qVecs, true)
			cc.LeftTo(&phiVs, false)
			cc.RightTo(&psiVs, false)

			if !floats.EqualApprox(corrs, test.wantCorrs, test.epsilon) {
				t.Errorf("%d use %d: unexpected variance result got:%v, want:%v",
					i, j, corrs, test.wantCorrs)
			}
			if !mat.EqualApprox(&pVecs, test.wantpVecs, test.epsilon) {
				t.Errorf("%d use %d: unexpected CCA result got:\n%v\nwant:\n%v",
					i, j, mat.Formatted(&pVecs), mat.Formatted(test.wantpVecs))
			}
			if !mat.EqualApprox(&qVecs, test.wantqVecs, test.epsilon) {
				t.Errorf("%d use %d: unexpected CCA result got:\n%v\nwant:\n%v",
					i, j, mat.Formatted(&qVecs), mat.Formatted(test.wantqVecs))
			}
			if !mat.EqualApprox(&phiVs, test.wantphiVs, test.epsilon) {
				t.Errorf("%d use %d: unexpected CCA result got:\n%v\nwant:\n%v",
					i, j, mat.Formatted(&phiVs), mat.Formatted(test.wantphiVs))
			}
			if !mat.EqualApprox(&psiVs, test.wantpsiVs, test.epsilon) {
				t.Errorf("%d use %d: unexpected CCA result got:\n%v\nwant:\n%v",
					i, j, mat.Formatted(&psiVs), mat.Formatted(test.wantpsiVs))
			}
		}
	}
}
