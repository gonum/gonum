// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package window_test

import (
	"fmt"
	"math"
	"math/cmplx"

	"gonum.org/v1/gonum/dsp/fourier"
	"gonum.org/v1/gonum/dsp/window"
)

func Example() {
	// The input sequence is a 2.5 period of the Sin function.
	src := make([]float64, 20)
	k := 5 * math.Pi / float64(len(src)-1)
	for i := range src {
		src[i] = math.Sin(k * float64(i))
	}

	// Initialize an FFT and perform the analysis.
	fft := fourier.NewFFT(len(src))
	coeff := fft.Coefficients(nil, src)

	// The result shows that width of the main lobe with center
	// between frequencies 0.1 and 0.15 is small, but that the
	// height of the side lobes is large.
	fmt.Println("Rectangular window (or no window):")
	for i, c := range coeff {
		fmt.Printf("freq=%.4f\tcycles/period, magnitude=%.4f,\tphase=%.4f\n",
			fft.Freq(i), cmplx.Abs(c), cmplx.Phase(c))
	}

	// Initialize an FFT and perform the analysis on a sequence
	// transformed by the Hamming window function.
	fft = fourier.NewFFT(len(src))
	coeff = fft.Coefficients(nil, window.Hamming(src))

	// The result shows that width of the main lobe is wider,
	// but height of the side lobes is lower.
	fmt.Println("Hamming window:")
	// The magnitude of all bins has been decreased by β.
	// Generally in an analysis amplification may be omitted, but to
	// make a comparable data, the result should be amplified by -β
	// of the window function — +5.37 dB for the Hamming window.
	//  -β = 20 log_10(amplifier).
	amplifier := math.Pow(10, 5.37/20.0)
	for i, c := range coeff {
		fmt.Printf("freq=%.4f\tcycles/period, magnitude=%.4f,\tphase=%.4f\n",
			fft.Freq(i), amplifier*cmplx.Abs(c), cmplx.Phase(c))
	}
	// Output:
	//
	// Rectangular window (or no window):
	// freq=0.0000	cycles/period, magnitude=2.2798,	phase=0.0000
	// freq=0.0500	cycles/period, magnitude=2.6542,	phase=0.1571
	// freq=0.1000	cycles/period, magnitude=5.3115,	phase=0.3142
	// freq=0.1500	cycles/period, magnitude=7.3247,	phase=-2.6704
	// freq=0.2000	cycles/period, magnitude=1.6163,	phase=-2.5133
	// freq=0.2500	cycles/period, magnitude=0.7681,	phase=-2.3562
	// freq=0.3000	cycles/period, magnitude=0.4385,	phase=-2.1991
	// freq=0.3500	cycles/period, magnitude=0.2640,	phase=-2.0420
	// freq=0.4000	cycles/period, magnitude=0.1530,	phase=-1.8850
	// freq=0.4500	cycles/period, magnitude=0.0707,	phase=-1.7279
	// freq=0.5000	cycles/period, magnitude=0.0000,	phase=0.0000
	// Hamming window:
	// freq=0.0000	cycles/period, magnitude=0.0218,	phase=3.1416
	// freq=0.0500	cycles/period, magnitude=0.8022,	phase=-2.9845
	// freq=0.1000	cycles/period, magnitude=7.1723,	phase=0.3142
	// freq=0.1500	cycles/period, magnitude=8.6285,	phase=-2.6704
	// freq=0.2000	cycles/period, magnitude=2.0420,	phase=0.6283
	// freq=0.2500	cycles/period, magnitude=0.0702,	phase=0.7854
	// freq=0.3000	cycles/period, magnitude=0.0217,	phase=-2.1991
	// freq=0.3500	cycles/period, magnitude=0.0259,	phase=-2.0420
	// freq=0.4000	cycles/period, magnitude=0.0184,	phase=-1.8850
	// freq=0.4500	cycles/period, magnitude=0.0092,	phase=-1.7279
	// freq=0.5000	cycles/period, magnitude=0.0000,	phase=0.0000
}

func ExampleHamming() {
	src := []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

	// Window functions change data in place. So, if input data
	// needs to stay unchanged, it must be copied.
	srcCpy := append([]float64(nil), src...)
	// Apply window function to srcCpy.
	dst := window.Hamming(srcCpy)

	// src is unchanged.
	fmt.Printf("src:    %f\n", src)
	// srcCpy is altered.
	fmt.Printf("srcCpy: %f\n", srcCpy)
	// dst mirrors the srcCpy slice.
	fmt.Printf("dst:    %f\n", dst)

	// Output:
	//
	// src:    [1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000]
	// srcCpy: [0.086957 0.111692 0.183218 0.293785 0.431409 0.581178 0.726861 0.852672 0.944977 0.993774 0.993774 0.944977 0.852672 0.726861 0.581178 0.431409 0.293785 0.183218 0.111692 0.086957]
	// dst:    [0.086957 0.111692 0.183218 0.293785 0.431409 0.581178 0.726861 0.852672 0.944977 0.993774 0.993774 0.944977 0.852672 0.726861 0.581178 0.431409 0.293785 0.183218 0.111692 0.086957]
}

func ExampleValues() {
	src := []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

	// Create a Sine Window lookup table.
	sine := window.NewValues(window.Sine, len(src))

	// Apply the transformation to the src.
	fmt.Printf("dst: %f\n", sine.Transform(src))

	// Output:
	//
	// dst: [0.000000 0.164595 0.324699 0.475947 0.614213 0.735724 0.837166 0.915773 0.969400 0.996584 0.996584 0.969400 0.915773 0.837166 0.735724 0.614213 0.475947 0.324699 0.164595 0.000000]
}
