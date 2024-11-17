package transform

import (
	"fmt"
	"math/cmplx"
)

func ExampleHilbert_AnalyticSignal() {
	// Samples is a set of real amplitudes that make up a signal.
	samples := []float64{1, 0, 2, 0, 4, 0, 2, 0}

	// Initialize a Hilbert transform and 'demodulate' to get the
	// analytic signal.
	// The result is the complex I/Q (In-Phase / Quadrature) demodulation
	// of the input signal.
	h := NewHilbert(len(samples))
	iqSamples := h.AnalyticSignal(nil, samples)

	// We can compute the instantaneous amplitude of the signal
	// (or 'envelope') using absolute value. Analyzing the envelope
	// is an easy way to measure changes in amplitude over time in a
	// signal.
	envelope := make([]float64, len(samples))
	for ind, iq := range iqSamples {
		envelope[ind] = cmplx.Abs(iq)
	}

	// We can also compute the instantaneous phase of each part of the
	// signal using the 4-quadrant arc-tangent. With multiple samples,
	// the instantaneous phase can be used to estimate instantaneous
	// frequency of a signal.
	phase := make([]float64, len(samples))
	for ind, iq := range iqSamples {
		phase[ind] = cmplx.Phase(iq)
	}

	for i, iq := range iqSamples {
		fmt.Printf("ind=%d -> I=%.4f, Q=%.4f, envelope=%.4f, phase=%.4f\n",
			i, real(iq), imag(iq), envelope[i], phase[i])
	}

	// Output:
	//
	// ind=0 -> I=1, Q=0, envelope=1, phase=0
	// ind=1 -> I=-9.571e-17, Q=-0.8107, envelope=0.8107, phase=-1.571
	// ind=2 -> I=2, Q=0, envelope=2, phase=0
	// ind=3 -> I=-1.263e-16, Q=-1.311, envelope=1.311, phase=-1.571
	// ind=4 -> I=4, Q=0, envelope=4, phase=0
	// ind=5 -> I=1.263e-16, Q=1.311, envelope=1.311, phase=1.571
	// ind=6 -> I=2, Q=0, envelope=2, phase=0
	// ind=7 -> I=9.571e-17, Q=0.8107, envelope=0.8107, phase=1.571
}