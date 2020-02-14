// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package window_test

import (
	"fmt"
	"math"
	"math/cmplx"

	"gonum.org/v1/gonum/dsp/window"
	"gonum.org/v1/gonum/fourier"
)

func Example() {
	// The input sequence is a 2.5 period of the Sin function.
	src := make([]float64, 20)
	k := 5.0 * math.Pi / float64(len(src)-1)
	for i := range src {
		src[i] = math.Sin(k * float64(i))
	}

	// Initialize an FFT and perform the analysis.
	fft := fourier.NewFFT(len(src))
	coeff := fft.Coefficients(nil, src)

	// The result shows, that width of the main lobe with center
	// between frequencies 0.1 and 0.15 is small,
	// but height of the side lobes is pretty big.
	fmt.Println("Rectangular window (or no window):")
	for i, c := range coeff {
		fmt.Printf("freq=%.4f\tcycles/period, magnitude=%.4f,\tphase=%.4f\n",
			fft.Freq(i), cmplx.Abs(c), cmplx.Phase(c))
	}

	// Initialize an FFT and perform the analysis on sequence
	// transformed by the Hamming window function.
	fft = fourier.NewFFT(len(src))
	coeff = fft.Coefficients(nil, window.Hamming(src))

	// The result shows, that width of the main lobe is wider,
	// but height of the side lobes is less.
	fmt.Println("Hamming window:")
	// The magnitude of all bins has been decreased on the β, [dB] value.
	// Generally, to perform analysis,  amplification could be omitted,
	// but to make a comparable data, result should be amplified.
	// result should be amplified on -β value of the Hamming window function,
	// which is +5.37 dB. -β = 20 log_10(amplifier), so...
	amplifier := math.Pow(10.0, 5.37/20.0)
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
	// freq=0.0000	cycles/period, magnitude=0.0542,	phase=3.1416
	// freq=0.0500	cycles/period, magnitude=0.8458,	phase=-2.9845
	// freq=0.1000	cycles/period, magnitude=7.1519,	phase=0.3142
	// freq=0.1500	cycles/period, magnitude=8.5907,	phase=-2.6704
	// freq=0.2000	cycles/period, magnitude=2.0804,	phase=0.6283
	// freq=0.2500	cycles/period, magnitude=0.0816,	phase=0.7854
	// freq=0.3000	cycles/period, magnitude=0.0156,	phase=-2.1991
	// freq=0.3500	cycles/period, magnitude=0.0224,	phase=-2.0420
	// freq=0.4000	cycles/period, magnitude=0.0163,	phase=-1.8850
	// freq=0.4500	cycles/period, magnitude=0.0083,	phase=-1.7279
	// freq=0.5000	cycles/period, magnitude=0.0000,	phase=0.0000
}

func ExampleHamming() {
	src := []float64{1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000,
		1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000}

	// Window function change data in place. So, if input data need to stay unchanged, it should be copied.
	srcCpy := append([]float64(nil), src...)
	// Applay Window function to srcCpy.
	dst := window.Hamming(srcCpy)

	//src data has not changed, because it has been copied.
	fmt.Printf("src:    %f\n", src)
	//srcCpy data has been changed by window function.
	fmt.Printf("srcCpy: %f\n", srcCpy)
	//dst contains the slice with the same data as srcDst with underlying array at the same address.
	fmt.Printf("dst:    %f\n", dst)

	// Output:
	//
	// src:    [1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000]
	// srcCpy: [0.080000 0.104924 0.176995 0.288404 0.427077 0.577986 0.724780 0.851550 0.944558 0.993726 0.993726 0.944558 0.851550 0.724780 0.577986 0.427077 0.288404 0.176995 0.104924 0.080000]
	// dst:    [0.080000 0.104924 0.176995 0.288404 0.427077 0.577986 0.724780 0.851550 0.944558 0.993726 0.993726 0.944558 0.851550 0.724780 0.577986 0.427077 0.288404 0.176995 0.104924 0.080000]
}
