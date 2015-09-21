package distmv

import (
	"testing"

	"github.com/gonum/floats"
	"github.com/gonum/matrix/mat64"
	"github.com/gonum/stat"
)

type mvTest struct {
	Mu      []float64
	Sigma   *mat64.SymDense
	Loc     []float64
	Logprob float64
	Prob    float64
}

func TestNormProbs(t *testing.T) {
	dist1, ok := NewNormal([]float64{0, 0}, mat64.NewSymDense(2, []float64{1, 0, 0, 1}), nil)
	if !ok {
		t.Errorf("bad test")
	}
	dist2, ok := NewNormal([]float64{6, 7}, mat64.NewSymDense(2, []float64{8, 2, 0, 4}), nil)
	if !ok {
		t.Errorf("bad test")
	}
	testProbability(t, []probCase{
		{
			dist:    dist1,
			loc:     []float64{0, 0},
			logProb: -1.837877066409345,
		},
		{
			dist:    dist2,
			loc:     []float64{6, 7},
			logProb: -3.503979321496947,
		},
		{
			dist:    dist2,
			loc:     []float64{1, 2},
			logProb: -7.075407892925519,
		},
	})
}

func TestNormRand(t *testing.T) {
	for _, test := range []struct {
		mean []float64
		cov  []float64
	}{
		{
			mean: []float64{0, 0},
			cov: []float64{
				1, 0,
				0, 1,
			},
		},
		{
			mean: []float64{0, 0},
			cov: []float64{
				1, 0.9,
				0.9, 1,
			},
		},
		{
			mean: []float64{6, 7},
			cov: []float64{
				5, 0.9,
				0.9, 2,
			},
		},
	} {
		dim := len(test.mean)
		cov := mat64.NewSymDense(dim, test.cov)
		n, ok := NewNormal(test.mean, cov, nil)
		if !ok {
			t.Errorf("bad covariance matrix")
		}

		nSamples := 1000000
		samps := mat64.NewDense(nSamples, dim, nil)
		for i := 0; i < nSamples; i++ {
			n.Rand(samps.RawRowView(i))
		}
		estMean := make([]float64, dim)
		for i := range estMean {
			estMean[i] = stat.Mean(samps.Col(nil, i), nil)
		}
		if !floats.EqualApprox(estMean, test.mean, 1e-2) {
			t.Errorf("Mean mismatch: want: %v, got %v", test.mean, estMean)
		}
		estCov := stat.CovarianceMatrix(nil, samps, nil)
		if !mat64.EqualApprox(estCov, cov, 1e-2) {
			t.Errorf("Cov mismatch: want: %v, got %v", cov, estCov)
		}
	}
}
