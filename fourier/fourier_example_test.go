// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fourier_test

import (
	"fmt"
	"math/cmplx"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/fourier"
)

func ExampleFFT_FFT() {
	// Period is a set of samples over a given period.
	period := []float64{1, 0, 2, 0, 4, 0, 2, 0}

	// Initialize an FFT and perform the analysis.
	fft := fourier.NewFFT(len(period))
	coeff := fft.FFT(nil, period)

	// Prepare the frequency tables.
	freqs := floats.Span(make([]float64, len(coeff)), 0, float64(len(coeff)-1))
	floats.Scale(1/float64(len(period)), freqs)

	for i, f := range freqs {
		fmt.Printf("freq=%v cycles/period, magnitude=%v, phase=%.4g\n",
			f, cmplx.Abs(coeff[i]), cmplx.Phase(coeff[i]))
	}

	// Output:
	//
	// freq=0 cycles/period, magnitude=9, phase=0
	// freq=0.125 cycles/period, magnitude=3, phase=3.142
	// freq=0.25 cycles/period, magnitude=1, phase=-0
	// freq=0.375 cycles/period, magnitude=3, phase=3.142
	// freq=0.5 cycles/period, magnitude=9, phase=0
}

func ExampleCmplxFFT_FFT() {
	// Period is a set of samples over a given period.
	period := []complex128{1, 0, 2, 0, 4, 0, 2, 0}

	// Initialize a complex FFT and perform the analysis.
	fft := fourier.NewCmplxFFT(len(period))
	coeff := fft.FFT(nil, period)

	// Prepare the frequency tables.
	freqs := make([]float64, len(coeff))
	n := (len(coeff) - 1) / 2
	floats.Span(freqs[:n+1], 0, float64(n))
	floats.Span(freqs[n+1:], -float64(len(coeff)/2), -1)
	floats.Scale(1/float64(len(period)), freqs)

	for i := range freqs {
		// Center the spectrum.
		i = fft.Shift(i)

		fmt.Printf("freq=%v cycles/period, magnitude=%v, phase=%.4g\n",
			freqs[i], cmplx.Abs(coeff[i]), cmplx.Phase(coeff[i]))
	}

	// Output:
	//
	// freq=-0.5 cycles/period, magnitude=9, phase=0
	// freq=-0.375 cycles/period, magnitude=3, phase=3.142
	// freq=-0.25 cycles/period, magnitude=1, phase=0
	// freq=-0.125 cycles/period, magnitude=3, phase=3.142
	// freq=0 cycles/period, magnitude=9, phase=0
	// freq=0.125 cycles/period, magnitude=3, phase=3.142
	// freq=0.25 cycles/period, magnitude=1, phase=0
	// freq=0.375 cycles/period, magnitude=3, phase=3.142
}
