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
	// The input sequence is a Sin function with 2.5 pereod per all sequence length
	src := make([]float64, 20)
	k := 5.0 * math.Pi / float64(len(src)-1)
	for i := range src {
		src[i] = math.Sin(k * float64(i))
	}

	// Initialize an FFT and perform the analysis.
	fft := fourier.NewFFT(len(src))
	coeff := fft.Coefficients(nil, src)

	// The result shows, that width of the main lobe with center between frequencies 0.1000 and 0.1500 is small, but height of the side lobes is pretty big.
	fmt.Println("Rectangle window (or no window):")
	for i, c := range coeff {
		fmt.Printf("freq=%.4f\tcycles/period, magnitude=%.4f,\tphase=%.4f\n",
			fft.Freq(i), cmplx.Abs(c), cmplx.Phase(c))
	}

	// Initialize an FFT and perform the analysis on sequence transformed by the Hamming window function.
	fft = fourier.NewFFT(len(src))
	coeff = fft.Coefficients(nil, window.Hamming(src))

	// The result shows, that width of the main lobe is wider, but height of the side lobes is less.
	fmt.Println("Hamming window:")
	// The magnitude of all bins has been decreased on the β, [dB] value.
	// Generally, to perform analysis,  amplification could be omitted, but to make a comparable data, result should be amplified.
	// result should be amplified on -β value of the Hamming window function, which is +5.37 dB.
	// -β = 20 log_10(amplifier), so...
	amplifier := math.Pow(10.0, 5.37/20.0)
	for i, c := range coeff {
		fmt.Printf("freq=%.4f\tcycles/period, magnitude=%.4f,\tphase=%.4f\n",
			fft.Freq(i), amplifier*cmplx.Abs(c), cmplx.Phase(c))
	}
	// Output:
	//
	// Rectangle window (or no window):
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

func ExampleRectangle() {
	// source data sequence.
	src := []float64{1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000,
		1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000}

	// Window function change data in place. So, if input data need to stay unchanged, it should be copied.
	srcCpy := append([]float64(nil), src...)
	// Applay Window function to srcCpy (Rectangle - basically does nothing).
	dst := window.Rectangle(srcCpy)

	//src data has not changed, because it has been copied.
	fmt.Printf("src:    %f\n", src)
	//srcCpy data has been changed by window function, wich is not true in this specific case, because Rectangle window does nothing.
	fmt.Printf("srcCpy: %f\n", srcCpy)
	//dst contains the slice with the same data as srcDst with underlying array at the same address.
	fmt.Printf("dst:    %f\n", dst)

	// Output:
	//
	// src:    [1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000]
	// srcCpy: [1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000]
	// dst:    [1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000]
}

func ExampleSine() {
	// source data sequence.
	src := []float64{1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000,
		1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000}

	// Window function change data in place. So, if input data need to stay unchanged, it should be copied.
	srcCpy := append([]float64(nil), src...)
	// Applay Window function to srcCpy.
	dst := window.Sine(srcCpy)

	//src data has not changed, because it has been copied.
	fmt.Printf("src:    %f\n", src)
	//srcCpy data has been changed by window function.
	fmt.Printf("srcCpy: %f\n", srcCpy)
	//dst contains the slice with the same data as srcDst with underlying array at the same address.
	fmt.Printf("dst:    %f\n", dst)

	// Output:
	//
	// src:    [1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000]
	// srcCpy: [0.000000 0.164595 0.324699 0.475947 0.614213 0.735724 0.837166 0.915773 0.969400 0.996584 0.996584 0.969400 0.915773 0.837166 0.735724 0.614213 0.475947 0.324699 0.164595 0.000000]
	// dst:    [0.000000 0.164595 0.324699 0.475947 0.614213 0.735724 0.837166 0.915773 0.969400 0.996584 0.996584 0.969400 0.915773 0.837166 0.735724 0.614213 0.475947 0.324699 0.164595 0.000000]
}

func ExampleLanczos() {
	// source data sequence.
	src := []float64{1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000,
		1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000}

	// Window function change data in place. So, if input data need to stay unchanged, it should be copied.
	srcCpy := append([]float64(nil), src...)
	// Applay Window function to srcCpy.
	dst := window.Lanczos(srcCpy)

	//src data has not changed, because it has been copied.
	fmt.Printf("src:    %f\n", src)
	//srcCpy data has been changed by window function.
	fmt.Printf("srcCpy: %f\n", srcCpy)
	//dst contains the slice with the same data as srcDst with underlying array at the same address.
	fmt.Printf("dst:    %f\n", dst)

	// Output:
	//
	// src:    [1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000]
	// srcCpy: [0.000000 0.115514 0.247646 0.389468 0.532984 0.669692 0.791213 0.889915 0.959492 0.995450 0.995450 0.959492 0.889915 0.791213 0.669692 0.532984 0.389468 0.247646 0.115514 0.000000]
	// dst:    [0.000000 0.115514 0.247646 0.389468 0.532984 0.669692 0.791213 0.889915 0.959492 0.995450 0.995450 0.959492 0.889915 0.791213 0.669692 0.532984 0.389468 0.247646 0.115514 0.000000]
}

func ExampleBartlett() {
	// source data sequence.
	src := []float64{1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000,
		1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000}

	// Window function change data in place. So, if input data need to stay unchanged, it should be copied.
	srcCpy := append([]float64(nil), src...)
	// Applay Window function to srcCpy.
	dst := window.Bartlett(srcCpy)

	//src data has not changed, because it has been copied.
	fmt.Printf("src:    %f\n", src)
	//srcCpy data has been changed by window function.
	fmt.Printf("srcCpy: %f\n", srcCpy)
	//dst contains the slice with the same data as srcDst with underlying array at the same address.
	fmt.Printf("dst:    %f\n", dst)

	// Output:
	//
	// src:    [1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000]
	// srcCpy: [0.000000 0.105263 0.210526 0.315789 0.421053 0.526316 0.631579 0.736842 0.842105 0.947368 0.947368 0.842105 0.736842 0.631579 0.526316 0.421053 0.315789 0.210526 0.105263 0.000000]
	// dst:    [0.000000 0.105263 0.210526 0.315789 0.421053 0.526316 0.631579 0.736842 0.842105 0.947368 0.947368 0.842105 0.736842 0.631579 0.526316 0.421053 0.315789 0.210526 0.105263 0.000000]
}

func ExampleHann() {
	// source data sequence.
	src := []float64{1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000,
		1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000}

	// Window function change data in place. So, if input data need to stay unchanged, it should be copied.
	srcCpy := append([]float64(nil), src...)
	// Applay Window function to srcCpy.
	dst := window.Hann(srcCpy)

	//src data has not changed, because it has been copied.
	fmt.Printf("src:    %f\n", src)
	//srcCpy data has been changed by window function.
	fmt.Printf("srcCpy: %f\n", srcCpy)
	//dst contains the slice with the same data as srcDst with underlying array at the same address.
	fmt.Printf("dst:    %f\n", dst)

	// Output:
	//
	// src:    [1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000]
	// srcCpy: [0.000000 0.027091 0.105430 0.226526 0.377257 0.541290 0.700848 0.838641 0.939737 0.993181 0.993181 0.939737 0.838641 0.700848 0.541290 0.377257 0.226526 0.105430 0.027091 0.000000]
	// dst:    [0.000000 0.027091 0.105430 0.226526 0.377257 0.541290 0.700848 0.838641 0.939737 0.993181 0.993181 0.939737 0.838641 0.700848 0.541290 0.377257 0.226526 0.105430 0.027091 0.000000]
}

func ExampleBartlettHann() {
	// source data sequence.
	src := []float64{1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000,
		1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000}

	// Window function change data in place. So, if input data need to stay unchanged, it should be copied.
	srcCpy := append([]float64(nil), src...)
	// Applay Window function to srcCpy.
	dst := window.BartlettHann(srcCpy)

	//src data has not changed, because it has been copied.
	fmt.Printf("src:    %f\n", src)
	//srcCpy data has been changed by window function.
	fmt.Printf("srcCpy: %f\n", srcCpy)
	//dst contains the slice with the same data as srcDst with underlying array at the same address.
	fmt.Printf("dst:    %f\n", dst)

	// Output:
	//
	// src:    [1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000]
	// srcCpy: [0.000000 0.045853 0.130653 0.247949 0.387768 0.537696 0.684223 0.814209 0.916305 0.982186 0.982186 0.916305 0.814209 0.684223 0.537696 0.387768 0.247949 0.130653 0.045853 0.000000]
	// dst:    [0.000000 0.045853 0.130653 0.247949 0.387768 0.537696 0.684223 0.814209 0.916305 0.982186 0.982186 0.916305 0.814209 0.684223 0.537696 0.387768 0.247949 0.130653 0.045853 0.000000]
}

func ExampleHamming() {
	// source data sequence.
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

func ExampleBlackman() {
	// source data sequence.
	src := []float64{1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000,
		1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000}

	// Window function change data in place. So, if input data need to stay unchanged, it should be copied.
	srcCpy := append([]float64(nil), src...)
	// Applay Window function to srcCpy.
	dst := window.Blackman(srcCpy)

	//src data has not changed, because it has been copied.
	fmt.Printf("src:    %f\n", src)
	//srcCpy data has been changed by window function.
	fmt.Printf("srcCpy: %f\n", srcCpy)
	//dst contains the slice with the same data as srcDst with underlying array at the same address.
	fmt.Printf("dst:    %f\n", dst)

	// Output:
	//
	// src:    [1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000]
	// srcCpy: [-0.000000 0.010223 0.045069 0.114390 0.226899 0.382381 0.566665 0.752034 0.903493 0.988846 0.988846 0.903493 0.752034 0.566665 0.382381 0.226899 0.114390 0.045069 0.010223 -0.000000]
	// dst:    [-0.000000 0.010223 0.045069 0.114390 0.226899 0.382381 0.566665 0.752034 0.903493 0.988846 0.988846 0.903493 0.752034 0.566665 0.382381 0.226899 0.114390 0.045069 0.010223 -0.000000]
}

func ExampleBlackmanHarris() {
	// source data sequence.
	src := []float64{1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000,
		1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000}

	// Window function change data in place. So, if input data need to stay unchanged, it should be copied.
	srcCpy := append([]float64(nil), src...)
	// Applay Window function to srcCpy.
	dst := window.BlackmanHarris(srcCpy)

	//src data has not changed, because it has been copied.
	fmt.Printf("src:    %f\n", src)
	//srcCpy data has been changed by window function.
	fmt.Printf("srcCpy: %f\n", srcCpy)
	//dst contains the slice with the same data as srcDst with underlying array at the same address.
	fmt.Printf("dst:    %f\n", dst)

	// Output:
	//
	// src:    [1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000]
	// srcCpy: [0.000060 0.002018 0.012795 0.046450 0.122540 0.256852 0.448160 0.668576 0.866426 0.984278 0.984278 0.866426 0.668576 0.448160 0.256852 0.122540 0.046450 0.012795 0.002018 0.000060]
	// dst:    [0.000060 0.002018 0.012795 0.046450 0.122540 0.256852 0.448160 0.668576 0.866426 0.984278 0.984278 0.866426 0.668576 0.448160 0.256852 0.122540 0.046450 0.012795 0.002018 0.000060]
}

func ExampleNuttall() {
	// source data sequence.
	src := []float64{1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000,
		1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000}

	// Window function change data in place. So, if input data need to stay unchanged, it should be copied.
	srcCpy := append([]float64(nil), src...)
	// Applay Window function to srcCpy.
	dst := window.Nuttall(srcCpy)

	//src data has not changed, because it has been copied.
	fmt.Printf("src:    %f\n", src)
	//srcCpy data has been changed by window function.
	fmt.Printf("srcCpy: %f\n", srcCpy)
	//dst contains the slice with the same data as srcDst with underlying array at the same address.
	fmt.Printf("dst:    %f\n", dst)

	// Output:
	//
	// src:    [1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000]
	// srcCpy: [-0.000000 0.001706 0.011614 0.043682 0.117808 0.250658 0.441946 0.664015 0.864348 0.984019 0.984019 0.864348 0.664015 0.441946 0.250658 0.117808 0.043682 0.011614 0.001706 -0.000000]
	// dst:    [-0.000000 0.001706 0.011614 0.043682 0.117808 0.250658 0.441946 0.664015 0.864348 0.984019 0.984019 0.864348 0.664015 0.441946 0.250658 0.117808 0.043682 0.011614 0.001706 -0.000000]
}

func ExampleBlackmanNuttall() {
	// source data sequence.
	src := []float64{1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000,
		1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000}

	// Window function change data in place. So, if input data need to stay unchanged, it should be copied.
	srcCpy := append([]float64(nil), src...)
	// Applay Window function to srcCpy.
	dst := window.BlackmanNuttall(srcCpy)

	//src data has not changed, because it has been copied.
	fmt.Printf("src:    %f\n", src)
	//srcCpy data has been changed by window function.
	fmt.Printf("srcCpy: %f\n", srcCpy)
	//dst contains the slice with the same data as srcDst with underlying array at the same address.
	fmt.Printf("dst:    %f\n", dst)

	// Output:
	//
	// src:    [1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000]
	// srcCpy: [0.000363 0.002885 0.015360 0.051652 0.130567 0.266629 0.457501 0.675215 0.869392 0.984644 0.984644 0.869392 0.675215 0.457501 0.266629 0.130567 0.051652 0.015360 0.002885 0.000363]
	// dst:    [0.000363 0.002885 0.015360 0.051652 0.130567 0.266629 0.457501 0.675215 0.869392 0.984644 0.984644 0.869392 0.675215 0.457501 0.266629 0.130567 0.051652 0.015360 0.002885 0.000363]
}

func ExampleFlatTop() {
	// source data sequence.
	src := []float64{1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000,
		1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000}

	// Window function change data in place. So, if input data need to stay unchanged, it should be copied.
	srcCpy := append([]float64(nil), src...)
	// Applay Window function to srcCpy.
	dst := window.FlatTop(srcCpy)

	//src data has not changed, because it has been copied.
	fmt.Printf("src:    %f\n", src)
	//srcCpy data has been changed by window function.
	fmt.Printf("srcCpy: %f\n", srcCpy)
	//dst contains the slice with the same data as srcDst with underlying array at the same address.
	fmt.Printf("dst:    %f\n", dst)

	// Output:
	//
	// src:    [1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000]
	// srcCpy: [-0.000421 -0.003687 -0.017675 -0.045939 -0.070137 -0.037444 0.115529 0.402051 0.737755 0.967756 0.967756 0.737755 0.402051 0.115529 -0.037444 -0.070137 -0.045939 -0.017675 -0.003687 -0.000421]
	// dst:    [-0.000421 -0.003687 -0.017675 -0.045939 -0.070137 -0.037444 0.115529 0.402051 0.737755 0.967756 0.967756 0.737755 0.402051 0.115529 -0.037444 -0.070137 -0.045939 -0.017675 -0.003687 -0.000421]
}

func ExampleGauss() {
	// source data sequence.
	src := []float64{1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000,
		1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000, 1.000000}

	// Window function change data in place. So, if input data need to stay unchanged, it should be copied.
	srcCpy := append([]float64(nil), src...)
	// Applay Window function to srcCpy.
	dst := window.Gauss(srcCpy, 0.3)

	//src data has not changed, because it has been copied.
	fmt.Printf("src:    %f\n", src)
	//srcCpy data has been changed by window function.
	fmt.Printf("srcCpy: %f\n", srcCpy)
	//dst contains the slice with the same data as srcDst with underlying array at the same address.
	fmt.Printf("dst:    %f\n", dst)

	// Output:
	//
	// src:    [1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000]
	// srcCpy: [0.003866 0.011708 0.031348 0.074214 0.155344 0.287499 0.470444 0.680632 0.870660 0.984728 0.984728 0.870660 0.680632 0.470444 0.287499 0.155344 0.074214 0.031348 0.011708 0.003866]
	// dst:    [0.003866 0.011708 0.031348 0.074214 0.155344 0.287499 0.470444 0.680632 0.870660 0.984728 0.984728 0.870660 0.680632 0.470444 0.287499 0.155344 0.074214 0.031348 0.011708 0.003866]
}
