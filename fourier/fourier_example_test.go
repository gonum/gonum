// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fourier_test

import (
	"fmt"
	"math/cmplx"

	"gonum.org/v1/gonum/fourier"
)

func ExampleFFT_FFT() {
	// Period is a set of samples over a given period.
	period := []float64{1, 0, 2, 0, 4, 0, 2, 0}

	// Initialize an FFT and perform the analysis.
	fft := fourier.NewFFT(len(period))
	coeff := fft.FFT(nil, period)

	for i, c := range coeff {
		fmt.Printf("freq=%v cycles/period, magnitude=%v, phase=%.4g\n",
			fft.Freq(i), cmplx.Abs(c), cmplx.Phase(c))
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

	for i := range coeff {
		// Center the spectrum.
		i = fft.Shift(i)

		fmt.Printf("freq=%v cycles/period, magnitude=%v, phase=%.4g\n",
			fft.Freq(i), cmplx.Abs(coeff[i]), cmplx.Phase(coeff[i]))
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
