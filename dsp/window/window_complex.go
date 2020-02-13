// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package window

import "math"

// RectangleComplex modifies the seq in place by Rectangle window and returns the seq.
//
// Rectangle window is a base high-resolution window
// (https://www.recordingblogs.com/wiki/rectangular-window):
//
// w[k] = 1,
// where N is the length of the window and k=0,1,...,N-1.
//
// The result is correspond to a selection of limited length sequence of values
// without any modification (so, it doesn't modifies the seq at all).
// Rectangle window has the lowest width of the main lobe and largest level of the side lobes.
//
// Spectral leakage parameters: ΔF_0 = 2, ΔF_0.5 = 0.89, K = 1, ɣ_max = -13, β =  0.
func RectangleComplex(seq []complex128) []complex128 {
	return seq
}

// SineComplex modifies the seq in place by Sine window and returns the seq.
//
// Sine window is a high-resolution window
// (https://www.recordingblogs.com/wiki/sine-window):
//
// w[k] = sin(π*k/(N-1)),
// where N is the length of the window and k=0,1,...,N-1.
//
// Spectral leakage parameters: ΔF_0 = 3, ΔF_0.5 = 1.23, K = 1.5, ɣ_max = -23, β = -3.93.
func SineComplex(seq []complex128) []complex128 {
	k := math.Pi / float64(len(seq)-1)
	for i := range seq {
		seq[i] *= complex(math.Sin(k*float64(i)), 0.0)
	}
	return seq
}

// LanczosComplex modifies the seq in place by Lanczos window and returns the seq.
//
// Lanczos window is a high-resolution window
// (https://www.recordingblogs.com/wiki/lanczos-window):
//
// w[k] = sinc( 2*k/(N-1) - 1),
// where N is the length of the window and k=0,1,...,N-1.
//
// Spectral leakage parameters: ΔF_0 = 3.24, ΔF_0.5 = 1.3, K = 1.62, ɣ_max = -26.4, β = -4.6.
func LanczosComplex(seq []complex128) []complex128 {
	k := 2.0 / float64(len(seq)-1)
	for i := range seq {
		x := math.Pi * (k*float64(i) - 1.0)
		seq[i] *= complex(math.Sin(x)/(x), 0.0)
	}
	return seq
}

// BartlettComplex modifies the seq in place by Bartlett window and returns the seq.
//
// Bartlett window is a high-resolution window
// (https://www.recordingblogs.com/wiki/triangular-window):
//
// w[k] = 1 - |k/A -1|, A=(N-1)/2,
// where N is the length of the window and k=0,1,...,N-1.
//
// Spectral leakage parameters: ΔF_0 = 4, ΔF_0.5 = 1.33, K = 2, ɣ_max = -26.5, β = -6.
func BartlettComplex(seq []complex128) []complex128 {
	a := float64(len(seq)-1) / 2.0
	for i := range seq {
		seq[i] *= complex(1.0-math.Abs(float64(i)/a-1.0), 0.0)
	}
	return seq
}

// HannComplex modifies the seq in place by Hann window and returns the seq.
//
// Hann window is a high-resolution window
// (https://www.recordingblogs.com/wiki/hann-window):
//
// w[k] = 0.5*(1 - cos(2*π*k/(N-1))),
// where N is the length of the window and k=0,1,...,N-1.
//
// Spectral leakage parameters: ΔF_0 = 4, ΔF_0.5 = 1.5, K = 2, ɣ_max = -31.5, β = -6.
func HannComplex(seq []complex128) []complex128 {
	k := 2.0 * math.Pi / float64(len(seq)-1)
	for i := range seq {
		seq[i] *= complex(0.5*(1-math.Cos(k*float64(i))), 0.0)
	}
	return seq
}

// BartlettHannComplex modifies the seq in place by Bartlett-Hann window and returns the seq.
//
// Bartlett-Hann window is a high-resolution window
// (https://www.recordingblogs.com/wiki/bartlett-hann-window):
//
// w[k] = 0.62 - 0.48*|k/(n-1)-0.5| - 0.38*cos(2*π*k/(N-1)),
// where N is the length of the window and k=0,1,...,N-1.
//
// Spectral leakage parameters: ΔF_0 = 4, ΔF_0.5 = 1.45, K = 2, ɣ_max = -35.9, β = -6.
func BartlettHannComplex(seq []complex128) []complex128 {
	const (
		a0 = 0.62
		a1 = 0.48
		a2 = 0.38
	)

	k := 2.0 * math.Pi / float64(len(seq)-1)
	for i := range seq {
		seq[i] *= complex(a0-a1*math.Abs(float64(i)/float64(len(seq)-1)-0.5)-a2*math.Cos(k*float64(i)), 0.0)
	}
	return seq
}

// HammingComplex -  modifies the seq in place by Hamming window and returns the seq.
//
// Hamming window is a high-resolution window
// (https://www.recordingblogs.com/wiki/hamming-window):
//
// w[k] = 0.54 - 0.46*cos(2*π*k/(N-1)),
// where N is the length of the window and k=0,1,...,N-1.
// Among K=2 windows it has a highest ɣ_max.
//
// Spectral leakage parameters: ΔF_0 = 4, ΔF_0.5 = 1.33, K = 2, ɣ_max = -42, β = -5.37.
func HammingComplex(seq []complex128) []complex128 {
	const (
		a0 = 0.54
		a1 = 0.46
	)

	k := 2.0 * math.Pi / float64(len(seq)-1)
	for i := range seq {
		seq[i] *= complex(a0-a1*math.Cos(k*float64(i)), 0.0)
	}
	return seq
}

// BlackmanComplex modifies the seq in place by Blackman window and returns the seq.
//
// Blackman window is a high-resolution window
// (https://www.recordingblogs.com/wiki/blackman-window):
//
// w[k] = 0.42 - 0.5*cos(2*π*k/(N-1)) + 0.08*cos(4*π*k/(N-1)),
// where N is the length of the window and k=0,1,...,N-1.
//
// Spectral leakage parameters: ΔF_0 = 6, ΔF_0.5 = 1.7, K = 3, ɣ_max = -58, β = -7.54.
func BlackmanComplex(seq []complex128) []complex128 {
	const (
		a0 = 0.42
		a1 = 0.5
		a2 = 0.08
	)

	k := 2.0 * math.Pi / float64(len(seq)-1)
	for i := range seq {
		x := k * float64(i)
		seq[i] *= complex(a0-a1*math.Cos(x)+a2*math.Cos(2.0*x), 0.0)
	}
	return seq
}

// BlackmanHarrisComplex modifies the seq in place by Blackman-Harris window and returns the seq.
//
// Blackman-Harris window is a low-resolution window
// (https://www.recordingblogs.com/wiki/blackman-harris-window):
//
// w[k] = 0.35875 - 0.48829*cos(2*π*k/(N-1)) +
// 0.14128*cos(4*π*k/(N-1)) - 0.01168*cos(6*π*k/(N-1)),
// where N is the length of the window and k=0,1,...,N-1.
//
// Spectral leakage parameters:  ΔF_0 = 8, ΔF_0.5 = 1.97, K = 4, ɣ_max = -92, β = -8.91.
func BlackmanHarrisComplex(seq []complex128) []complex128 {
	const (
		a0 = 0.35875
		a1 = 0.48829
		a2 = 0.14128
		a3 = 0.01168
	)

	k := 2.0 * math.Pi / float64(len(seq)-1)
	for i := range seq {
		x := k * float64(i)
		seq[i] *= complex(a0-a1*math.Cos(x)+a2*math.Cos(2.0*x)-a3*math.Cos(3.0*x), 0.0)
	}
	return seq
}

// NuttallComplex modifies the seq in place by Nuttall window and returns the seq.
//
// Nuttall window is a low-resolution window
// (https://www.recordingblogs.com/wiki/nuttall-window):
//
// w[k] = 0.355768 - 0.487396*cos(2*π*k/(N-1)) + 0.144232*cos(4*π*k/(N-1)) -
// 0.012604*cos(6*π*k/(N-1)),
// where N is the length of the window and k=0,1,...,N-1.
//
// Spectral leakage parameters: ΔF_0 = 8, ΔF_0.5 = 1.98, K = 4, ɣ_max = -93, β = -9.
func NuttallComplex(seq []complex128) []complex128 {
	const (
		a0 = 0.355768
		a1 = 0.487396
		a2 = 0.144232
		a3 = 0.012604
	)

	k := 2.0 * math.Pi / float64(len(seq)-1)
	for i := range seq {
		x := k * float64(i)
		seq[i] *= complex(a0-a1*math.Cos(x)+a2*math.Cos(2.0*x)-a3*math.Cos(3.0*x), 0.0)
	}
	return seq
}

// BlackmanNuttallComplex modifies the seq in place by Blackman-Nuttall window and returns the seq.
//
// Blackman-Nuttall window is a low-resolution window
// (https://www.recordingblogs.com/wiki/blackman-nuttall-window):
//
// w[k] = 0.3635819 - 0.4891775*cos(2*π*k/(N-1)) + 0.1365995*cos(4*π*k/(N-1)) -
// 0.0106411*cos(6*π*k/(N-1)),
// where N is the length of the window and k=0,1,...,N-1.
//
// Spectral leakage parameters: ΔF_0 = 8, ΔF_0.5 = 1.94, K = 4, ɣ_max = -98, β = -8.8.
func BlackmanNuttallComplex(seq []complex128) []complex128 {
	const (
		a0 = 0.3635819
		a1 = 0.4891775
		a2 = 0.1365995
		a3 = 0.0106411
	)

	k := 2.0 * math.Pi / float64(len(seq)-1)
	for i := range seq {
		x := k * float64(i)
		seq[i] *= complex(a0-a1*math.Cos(x)+a2*math.Cos(2.0*x)-a3*math.Cos(3.0*x), 0.0)
	}
	return seq
}

// FlatTopComplex modifies the seq in place by Flat Top window and returns the seq.
//
// Flat Top window is a low-resolution window
// (https://www.recordingblogs.com/wiki/flat-top-window):
//
// w[k] = 0.21557895 - 0.41663158*cos(2*π*k/(N-1)) +
// 0.277263158*cos(4*π*k/(N-1)) - 0.083578947*cos(6*π*k/(N-1)) +
// 0.006947368*cos(4*π*k/(N-1)),
// where N is the length of the window and k=0,1,...,N-1.
//
// Spectral leakage parameters: ΔF_0 = 10, ΔF_0.5 = 3.72, K = 5, ɣ_max = -93.0, β = -13.34.
func FlatTopComplex(seq []complex128) []complex128 {
	const (
		a0 = 1.0
		a1 = 1.93
		a2 = 1.29
		a3 = 0.388
		a4 = 0.032
	)

	k := 2.0 * math.Pi / float64(len(seq)-1)
	for i := range seq {
		x := k * float64(i)
		seq[i] *= complex(a0-a1*math.Cos(x)+a2*math.Cos(2.0*x)-a3*math.Cos(3.0*x)+a4*math.Cos(4.0*x), 0.0)
	}
	return seq
}

// GaussComplex modifies the seq in place by Gauss window and returns the seq.
//
// Gauss window is a adjustable window
// (https://www.recordingblogs.com/wiki/gaussian-window):
//
// w[k] = exp(-0.5 * ((k-M)/(σ*M))² ), M = (N-1)/2,
// where N is the length of the window and k=0,1,...,N-1.
// The properties of window depends on σ (sigma) argument.
// It can be used as high or low resolution window, depends of a σ value.
//
// Spectral leakage parameters are summarized in the table:
//         |  σ=0.3 |  σ=0.5 |  σ=1.2 |
//  -------|--------------------------|
//  ΔF_0   |  8     |  3.4   |  2.2   |
//  ΔF_0.5 |  1.82  |  1.2   |  0.94  |
//  K      |  4     |  1.7   |  1.1   |
//  ɣ_max  | -65    | -31.5  | -15.5  |
//  β      | -8.52  | -4.48  | -0.96  |
func GaussComplex(seq []complex128, sigma float64) []complex128 {
	a := float64(len(seq)-1) / 2.0
	for i := range seq {
		x := -0.5 * math.Pow((float64(i)-a)/(sigma*a), 2)
		seq[i] *= complex(math.Exp(x), 0.0)
	}
	return seq
}
