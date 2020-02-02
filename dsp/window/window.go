// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package window

import (
	"math"
)

// Rectangle - modifies the seq in place by Rectangle window and returns the seq.
//
// Rectangle window is a base high-resolution window,
// which result is correspond to a selection of limited length sequence of values without any modification
// (so, it doesn't modifies the seq at all).
// Rectangle window has the lowest width of the main lobe and largest level of the side lobes.
//
// Indicators of quality:
//  ΔF_0	=  2
//  ΔF_0.5	=  0.89
//  K	=  1
//  ɣ_max	= -13
//  β	=  0
func Rectangle(seq []float64) []float64 {
	return seq
}

// Sin - modifies the seq in place by Sin window and returns the seq.
//
// Sin window is a high-resolution window.
//
// Indicators of quality:
//  ΔF_0	=  3
//  ΔF_0.5	=  1.23
//  K	=  1.5
//  ɣ_max	= -23
//  β	= -3.93
func Sin(seq []float64) []float64 {
	N := len(seq)

	k := math.Pi / float64(N-1)
	for n := range seq {
		seq[n] *= math.Sin(k * float64(n))
	}
	return seq
}

// Lanczos - modifies the seq in place by Lanczos window and returns the seq.
//
// Lanczos window is a high-resolution window.
//
// Indicators of quality:
//  ΔF_0	=  3.24
//  ΔF_0.5	=  1.3
//  K	=  1.62
//  ɣ_max	= -26.4
//  β	= -4.6
func Lanczos(seq []float64) []float64 {
	N := len(seq)

	var x float64
	k := 2.0 / float64(N-1)
	for n := range seq {
		x = math.Pi * (k*float64(n) - 1.0)
		seq[n] *= math.Sin(x) / (x)
	}
	return seq
}

// Bartlett - modifies the seq in place by Bartlett window and returns the seq.
//
// Bartlett window is a high-resolution window.
//
// Indicators of quality:
//  ΔF_0	=  4
//  ΔF_0.5	=  1.33
//  K	=  2
//  ɣ_max	= -26.5
//  β	= -6
func Bartlett(seq []float64) []float64 {
	N := len(seq)

	A := float64(N-1) / 2.0
	for n := range seq {
		seq[n] *= 1.0 - math.Abs(float64(n)/A-1.0)
	}
	return seq
}

// Hann - modifies the seq in place by Hann window and returns the seq.
//
// Hann window is a high-resolution window.
//  ΔF_0	=  4
//  ΔF_0.5	=  1.5
//  K	=  2
//  ɣ_max	= -31.5
//  β	= -6
func Hann(seq []float64) []float64 {
	N := len(seq)

	k := 2.0 * math.Pi / float64(N-1)
	for n := range seq {
		seq[n] *= 0.5 * (1 - math.Cos(k*float64(n)))
	}
	return seq
}

// BartlettHann - modifies the seq in place by Bartlett-Hann window and returns the seq.
//
// Bartlett-Hann window is a high-resolution window.
//
// Indicators of quality:
//  ΔF_0	=  4
//  ΔF_0.5	=  1.45
//  K	=  2
//  ɣ_max	= -35.9
//  β	= -6
func BartlettHann(seq []float64) []float64 {
	N := len(seq)

	a0, a1, a2 := 0.62, 0.48, 0.38
	k := 2.0 * math.Pi / float64(N-1)
	for n := range seq {
		seq[n] *= a0 - a1*math.Abs(float64(n)/float64(N-1)-0.5) - a2*math.Cos(k*float64(n))
	}
	return seq
}

// Hamming -  modifies the seq in place by Hamming window and returns the seq.
//
// Hamming window is a high-resolution window. Among K=2 windows it has a highest ɣ_max.
//
// Indicators of quality:
//  ΔF_0	=  4
//  ΔF_0.5	=  1.33
//  K	=  2
//  ɣ_max	= -42
//  β	= -5.37
func Hamming(seq []float64) []float64 {
	N := len(seq)

	a0, a1 := 0.54, 0.46
	k := 2.0 * math.Pi / float64(N-1)
	for n := range seq {
		seq[n] *= a0 - a1*math.Cos(k*float64(n))
	}
	return seq
}

// Blackman - modifies the seq in place by Blackman window and returns the seq.
//
// Blackman window is a high-resolution window.
//
// Indicators of quality:
//  ΔF_0	=  6
//  ΔF_0.5	=  1.7
//  K	=  3
//  ɣ_max	= -58
//  β	= -7.54
func Blackman(seq []float64) []float64 {
	N := len(seq)

	a0, a1, a2 := 0.42, 0.5, 0.08
	var x float64
	k := 2.0 * math.Pi / float64(N-1)
	for n := range seq {
		x = k * float64(n)
		seq[n] *= a0 - a1*math.Cos(x) + a2*math.Cos(2.0*x)
	}
	return seq
}

// BlackmanHarris - modifies the seq in place by Blackman-Harris window and returns the seq.
//
// Blackman-Harris window is a low-resolution window.
//
// Indicators of quality:
//  ΔF_0	=  8
//  ΔF_0.5	=  1.97
//  K	=  4
//  ɣ_max	= -92
//  β	= -8.91
func BlackmanHarris(seq []float64) []float64 {
	N := len(seq)

	a0, a1, a2, a3 := 0.35875, 0.48829, 0.14128, 0.01168
	var x float64
	k := 2.0 * math.Pi / float64(N-1)
	for n := range seq {
		x = k * float64(n)
		seq[n] *= a0 - a1*math.Cos(x) + a2*math.Cos(2.0*x) - a3*math.Cos(3.0*x)
	}
	return seq
}

// Nuttall - modifies the seq in place by Nuttall window and returns the seq.
//
// Nuttall window is a low-resolution window.
//
// Indicators of quality:
//  ΔF_0	=  8
//  ΔF_0.5	=  1.98
//  K	=  4
//  ɣ_max	= -93
//  β	= -9
func Nuttall(seq []float64) []float64 {
	N := len(seq)

	a0, a1, a2, a3 := 0.355768, 0.487396, 0.144232, 0.012604
	var x float64
	k := 2.0 * math.Pi / float64(N-1)
	for n := range seq {
		x = k * float64(n)
		seq[n] *= a0 - a1*math.Cos(x) + a2*math.Cos(2.0*x) - a3*math.Cos(3.0*x)
	}
	return seq
}

// BlackmanNuttall - modifies the seq in place by Blackman-Nuttall window and returns the seq.
//
// Blackman-Nuttall window is a low-resolution window.
//
// Indicators of quality:
//  ΔF_0	=  8
//  ΔF_0.5	=  1.94
//  K	=  4
//  ɣ_max	= -98
//  β	= -8.8
func BlackmanNuttall(seq []float64) []float64 {
	N := len(seq)

	a0, a1, a2, a3 := 0.3635819, 0.4891775, 0.1365995, 0.0106411
	var x float64
	k := 2.0 * math.Pi / float64(N-1)
	for n := range seq {
		x = k * float64(n)
		seq[n] *= a0 - a1*math.Cos(x) + a2*math.Cos(2.0*x) - a3*math.Cos(3.0*x)
	}
	return seq
}

// FlatTop - modifies the seq in place by Flat Top window and returns the seq.
//
// Flat Top window is a low-resolution window.
//
// Indicators of quality:
//  ΔF_0	=  10
//  ΔF_0.5	=  3.86
//  K	=  5
//  ɣ_max	= -69
//  β	=  0
func FlatTop(seq []float64) []float64 {
	N := len(seq)

	a0, a1, a2, a3, a4 := 1.0, 1.93, 1.29, 0.388, 0.032
	var x float64
	k := 2.0 * math.Pi / float64(N-1)
	for n := range seq {
		x = k * float64(n)
		seq[n] *= a0 - a1*math.Cos(x) + a2*math.Cos(2.0*x) - a3*math.Cos(3.0*x) + a4*math.Cos(4.0*x)
	}
	return seq
}

// Gauss - modifies the seq in place by Gauss window and returns the seq.
//
// Gauss window is a adjustable window. The properties of window depends on sigma argument. It can be used as high or low resolution window, depends of a sigma value.
//
// Indicators of quality:
//  	   sigma=0.3	 sigma=0.5	 sigma=1.2
//  ΔF_0	=  8		 3.4		 2.2
//  ΔF_0.5	=  1.82		 1.2		 0.94
//  K	=  4		 1.7		 1.1
//  ɣ_max	= -65		-31.5		-15.5
//  β	= -8.52		-4.48		-0.96
func Gauss(seq []float64, sigma float64) []float64 {
	N := len(seq)

	var x float64
	A := float64(N-1) / 2.0
	for n := range seq {
		x = -0.5 * math.Pow((float64(n)-A)/(sigma*A), 2)
		seq[n] *= math.Exp(x)
	}
	return seq
}
