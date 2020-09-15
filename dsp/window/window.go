// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package window

import "math"

// Rectangular modifies seq in place by the Rectangular window and returns the result.
// See https://en.wikipedia.org/wiki/Window_function#Rectangular_window and
// https://www.recordingblogs.com/wiki/rectangular-window for details.
//
// The rectangular window has the lowest width of the main lobe and largest
// level of the side lobes. The result corresponds to a selection of
// limited length sequence of values without any modification.
//
// The sequence weights are
//  w[k] = 1,
// for k=0,1,...,N-1 where N is the length of the window.
//
// Spectral leakage parameters: ΔF_0 = 2, ΔF_0.5 = 0.89, K = 1, ɣ_max = -13, β = 0.
func Rectangular(seq []float64) []float64 {
	return seq
}

// Sine modifies seq in place by the Sine window and returns the result.
// See https://en.wikipedia.org/wiki/Window_function#Sine_window and
// https://www.recordingblogs.com/wiki/sine-window for details.
//
// Sine window is a high-resolution window.
//
// The sequence weights are
//  w[k] = sin(π*k/(N-1)),
// for k=0,1,...,N-1 where N is the length of the window.
//
// Spectral leakage parameters: ΔF_0 = 3, ΔF_0.5 = 1.23, K = 1.5, ɣ_max = -23, β = -3.93.
func Sine(seq []float64) []float64 {
	k := math.Pi / float64(len(seq)-1)
	for i := range seq {
		seq[i] *= math.Sin(k * float64(i))
	}
	return seq
}

// Lanczos modifies seq in place by the Lanczos window and returns the result.
// See https://en.wikipedia.org/wiki/Window_function#Lanczos_window and
// https://www.recordingblogs.com/wiki/lanczos-window for details.
//
// The Lanczos window is a high-resolution window.
//
// The sequence weights are
//  w[k] = sinc(2*k/(N-1) - 1),
// for k=0,1,...,N-1 where N is the length of the window.
//
// Spectral leakage parameters: ΔF_0 = 3.24, ΔF_0.5 = 1.3, K = 1.62, ɣ_max = -26.4, β = -4.6.
func Lanczos(seq []float64) []float64 {
	k := 2 / float64(len(seq)-1)
	for i := range seq {
		x := math.Pi * (k*float64(i) - 1)
		if x == 0 {
			// Avoid NaN.
			continue
		}
		seq[i] *= math.Sin(x) / x
	}
	return seq
}

// Triangular modifies seq in place by the Triangular window and returns the result.
// See https://en.wikipedia.org/wiki/Window_function#Triangular_window and
// https://www.recordingblogs.com/wiki/triangular-window for details.
//
// The Triangular window is a high-resolution window.
//
// The sequence weights are
//  w[k] = 1 - |k/A -1|, A=(N-1)/2,
// for k=0,1,...,N-1 where N is the length of the window.
//
// Spectral leakage parameters: ΔF_0 = 4, ΔF_0.5 = 1.33, K = 2, ɣ_max = -26.5, β = -6.
func Triangular(seq []float64) []float64 {
	a := float64(len(seq)-1) / 2
	for i := range seq {
		seq[i] *= 1 - math.Abs(float64(i)/a-1)
	}
	return seq
}

// Hann modifies seq in place by the Hann window and returns the result.
// See https://en.wikipedia.org/wiki/Window_function#Hann_and_Hamming_windows
// and https://www.recordingblogs.com/wiki/hann-window for details.
//
// The Hann window is a high-resolution window.
//
// The sequence weights are
//  w[k] = 0.5*(1 - cos(2*π*k/(N-1))),
// for k=0,1,...,N-1 where N is the length of the window.
//
// Spectral leakage parameters: ΔF_0 = 4, ΔF_0.5 = 1.5, K = 2, ɣ_max = -31.5, β = -6.
func Hann(seq []float64) []float64 {
	k := 2 * math.Pi / float64(len(seq)-1)
	for i := range seq {
		seq[i] *= 0.5 * (1 - math.Cos(k*float64(i)))
	}
	return seq
}

// BartlettHann modifies seq in place by the Bartlett-Hann window and returns result.
// See https://en.wikipedia.org/wiki/Window_function#Bartlett%E2%80%93Hann_window
// and https://www.recordingblogs.com/wiki/bartlett-hann-window for details.
//
// The Bartlett-Hann window is a high-resolution window.
//
// The sequence weights are
//  w[k] = 0.62 - 0.48*|k/(N-1)-0.5| - 0.38*cos(2*π*k/(N-1)),
// for k=0,1,...,N-1 where N is the length of the window.
//
// Spectral leakage parameters: ΔF_0 = 4, ΔF_0.5 = 1.45, K = 2, ɣ_max = -35.9, β = -6.
func BartlettHann(seq []float64) []float64 {
	const (
		a0 = 0.62
		a1 = 0.48
		a2 = 0.38
	)

	k := 2 * math.Pi / float64(len(seq)-1)
	for i := range seq {
		seq[i] *= a0 - a1*math.Abs(float64(i)/float64(len(seq)-1)-0.5) - a2*math.Cos(k*float64(i))
	}
	return seq
}

// Hamming modifies seq in place by the Hamming window and returns the result.
// See https://en.wikipedia.org/wiki/Window_function#Hann_and_Hamming_windows
// and https://www.recordingblogs.com/wiki/hamming-window for details.
//
// The Hamming window is a high-resolution window. Among K=2 windows it has
// the highest ɣ_max.
//
// The sequence weights are
//  w[k] = 25/46 - 21/46 * cos(2*π*k/(N-1)),
// for k=0,1,...,N-1 where N is the length of the window.
//
// Spectral leakage parameters: ΔF_0 = 4, ΔF_0.5 = 1.33, K = 2, ɣ_max = -42, β = -5.37.
func Hamming(seq []float64) []float64 {
	const (
		a0 = 25.0 / 46.0
		a1 = 1 - a0
	)

	k := 2 * math.Pi / float64(len(seq)-1)
	for i := range seq {
		seq[i] *= a0 - a1*math.Cos(k*float64(i))
	}
	return seq
}

// Blackman modifies seq in place by the Blackman window and returns the result.
// See https://en.wikipedia.org/wiki/Window_function#Blackman_window and
// https://www.recordingblogs.com/wiki/blackman-window for details.
//
// The Blackman window is a high-resolution window.
//
// The sequence weights are
//  w[k] = 0.42 - 0.5*cos(2*π*k/(N-1)) + 0.08*cos(4*π*k/(N-1)),
// for k=0,1,...,N-1 where N is the length of the window.
//
// Spectral leakage parameters: ΔF_0 = 6, ΔF_0.5 = 1.7, K = 3, ɣ_max = -58, β = -7.54.
func Blackman(seq []float64) []float64 {
	const (
		a0 = 0.42
		a1 = 0.5
		a2 = 0.08
	)

	k := 2 * math.Pi / float64(len(seq)-1)
	for i := range seq {
		x := k * float64(i)
		seq[i] *= a0 - a1*math.Cos(x) + a2*math.Cos(2*x)
	}
	return seq
}

// BlackmanHarris modifies seq in place by the Blackman-Harris window and returns the result.
// See https://en.wikipedia.org/wiki/Window_function#Blackman%E2%80%93Harris_window
// and https://www.recordingblogs.com/wiki/blackman-harris-window for details.
//
// The Blackman-Harris window is a low-resolution window.
//
// The sequence weights are
//  w[k] = 0.35875 - 0.48829*cos(2*π*k/(N-1)) +
//         0.14128*cos(4*π*k/(N-1)) - 0.01168*cos(6*π*k/(N-1)),
// for k=0,1,...,N-1 where N is the length of the window.
//
// Spectral leakage parameters:  ΔF_0 = 8, ΔF_0.5 = 1.97, K = 4, ɣ_max = -92, β = -8.91.
func BlackmanHarris(seq []float64) []float64 {
	const (
		a0 = 0.35875
		a1 = 0.48829
		a2 = 0.14128
		a3 = 0.01168
	)

	k := 2 * math.Pi / float64(len(seq)-1)
	for i := range seq {
		x := k * float64(i)
		seq[i] *= a0 - a1*math.Cos(x) + a2*math.Cos(2*x) - a3*math.Cos(3*x)
	}
	return seq
}

// Nuttall modifies seq in place by the Nuttall window and returns the result.
// See https://en.wikipedia.org/wiki/Window_function#Nuttall_window,_continuous_first_derivative
// and https://www.recordingblogs.com/wiki/nuttall-window for details.
//
// The Nuttall window is a low-resolution window.
//
// The sequence weights are
//  w[k] = 0.355768 - 0.487396*cos(2*π*k/(N-1)) + 0.144232*cos(4*π*k/(N-1)) -
//         0.012604*cos(6*π*k/(N-1)),
// for k=0,1,...,N-1 where N is the length of the window.
//
// Spectral leakage parameters: ΔF_0 = 8, ΔF_0.5 = 1.98, K = 4, ɣ_max = -93, β = -9.
func Nuttall(seq []float64) []float64 {
	const (
		a0 = 0.355768
		a1 = 0.487396
		a2 = 0.144232
		a3 = 0.012604
	)

	k := 2 * math.Pi / float64(len(seq)-1)
	for i := range seq {
		x := k * float64(i)
		seq[i] *= a0 - a1*math.Cos(x) + a2*math.Cos(2*x) - a3*math.Cos(3*x)
	}
	return seq
}

// BlackmanNuttall modifies seq in place by the Blackman-Nuttall window and returns the result.
// See https://en.wikipedia.org/wiki/Window_function#Blackman%E2%80%93Nuttall_window
// and https://www.recordingblogs.com/wiki/blackman-nuttall-window for details.
//
// The Blackman-Nuttall window is a low-resolution window.
//
// The sequence weights are
//  w[k] = 0.3635819 - 0.4891775*cos(2*π*k/(N-1)) + 0.1365995*cos(4*π*k/(N-1)) -
//         0.0106411*cos(6*π*k/(N-1)),
// for k=0,1,...,N-1 where N is the length of the window.
//
// Spectral leakage parameters: ΔF_0 = 8, ΔF_0.5 = 1.94, K = 4, ɣ_max = -98, β = -8.8.
func BlackmanNuttall(seq []float64) []float64 {
	const (
		a0 = 0.3635819
		a1 = 0.4891775
		a2 = 0.1365995
		a3 = 0.0106411
	)

	k := 2 * math.Pi / float64(len(seq)-1)
	for i := range seq {
		x := k * float64(i)
		seq[i] *= a0 - a1*math.Cos(x) + a2*math.Cos(2*x) - a3*math.Cos(3*x)
	}
	return seq
}

// FlatTop modifies seq in place by the Flat Top window and returns the result.
// See https://en.wikipedia.org/wiki/Window_function#Flat_top_window and
// https://www.recordingblogs.com/wiki/flat-top-window for details.
//
// The Flat Top window is a low-resolution window.
//
// The sequence weights are
//  w[k] = 0.21557895 - 0.41663158*cos(2*π*k/(N-1)) +
//         0.277263158*cos(4*π*k/(N-1)) - 0.083578947*cos(6*π*k/(N-1)) +
//         0.006947368*cos(4*π*k/(N-1)),
// for k=0,1,...,N-1 where N is the length of the window.
//
// Spectral leakage parameters: ΔF_0 = 10, ΔF_0.5 = 3.72, K = 5, ɣ_max = -93.0, β = -13.34.
func FlatTop(seq []float64) []float64 {
	const (
		a0 = 0.21557895
		a1 = 0.41663158
		a2 = 0.277263158
		a3 = 0.083578947
		a4 = 0.006947368
	)

	k := 2 * math.Pi / float64(len(seq)-1)
	for i := range seq {
		x := k * float64(i)
		seq[i] *= a0 - a1*math.Cos(x) + a2*math.Cos(2*x) - a3*math.Cos(3*x) + a4*math.Cos(4*x)
	}
	return seq
}

// Gaussian can modify a sequence by the Gaussian window and return the result.
// See https://en.wikipedia.org/wiki/Window_function#Gaussian_window
// and https://www.recordingblogs.com/wiki/gaussian-window for details.
//
// The Gaussian window is an adjustable window.
//
// The sequence weights are
//  w[k] = exp(-0.5 * ((k-M)/(σ*M))² ), M = (N-1)/2,
// for k=0,1,...,N-1 where N is the length of the window.
//
// The properties of the window depend on the value of σ (sigma).
// It can be used as high or low resolution window, depending of the σ value.
//
// Spectral leakage parameters are summarized in the table:
//         |  σ=0.3  |  σ=0.5 |  σ=1.2 |
//  -------|---------------------------|
//  ΔF_0   |   8     |   3.4  |   2.2  |
//  ΔF_0.5 |   1.82  |   1.2  |   0.94 |
//  K      |   4     |   1.7  |   1.1  |
//  ɣ_max  | -65     | -31.5  | -15.5  |
//  β      |  -8.52  |  -4.48 |  -0.96 |
type Gaussian struct {
	Sigma float64
}

// Transform applies the Gaussian transformation to seq in place, using the value
// of the receiver as the sigma parameter, and returning the result.
func (g Gaussian) Transform(seq []float64) []float64 {
	a := float64(len(seq)-1) / 2
	for i := range seq {
		x := -0.5 * math.Pow((float64(i)-a)/(g.Sigma*a), 2)
		seq[i] *= math.Exp(x)
	}
	return seq
}

// Tukey can modify a sequence using the Tukey window and return the result.
// See https://en.wikipedia.org/wiki/Window_function#Tukey_window
// and https://www.recordingblogs.com/wiki/tukey-window for details.
//
// The Tukey window is an adjustable window.
//
// The sequence weights are
//  w[k] = 0.5 * (1 + cos(π*(|k - M| - αM)/((1-α) * M))), |k - M| ≥ αM
//       = 1, |k - M| < αM
// with M = (N - 1)/2 for k=0,1,...,N-1 where N is the length of the window.
//
// Spectral leakage parameters are summarized in the table:
//         |  α=0.3 |  α=0.5 |  α=0.7 |
//  -------|--------------------------|
//  ΔF_0   |   1.33 |   1.22 |   1.13 |
//  ΔF_0.5 |   1.28 |   1.16 |   1.04 |
//  K      |   0.67 |   0.61 |   0.57 |
//  ɣ_max  | -18.2  | -15.1  | -13.8  |
//  β      |  -1.41 |  -2.50 |  -3.74 |
type Tukey struct {
	Alpha float64
}

// Transform applies the Tukey transformation to seq in place, using the value
// of the receiver as the Alpha parameter, and returning the result
func (t Tukey) Transform(seq []float64) []float64 {
	switch {
	case t.Alpha <= 0:
		return Rectangular(seq)
	case t.Alpha >= 1:
		return Hann(seq)
	default:
		alphaL := t.Alpha * float64(len(seq)-1)
		width := int(0.5*alphaL) + 1
		for i := range seq[:width] {
			w := 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/alphaL))
			seq[i] *= w
			seq[len(seq)-1-i] *= w
		}
		return seq
	}
}

// Values is an arbitrary real window function.
type Values []float64

// NewValues returns a Values of length n with weights corresponding to the
// provided window function.
func NewValues(window func([]float64) []float64, n int) Values {
	v := make(Values, n)
	for i := range v {
		v[i] = 1
	}
	return window(v)
}

// Transform applies the weights in the receiver to seq in place, returning the
// result. If v is nil, Transform is a no-op, otherwise the length of v must
// match the length of seq.
func (v Values) Transform(seq []float64) []float64 {
	if v == nil {
		return seq
	}
	if len(v) != len(seq) {
		panic("window: length mismatch")
	}
	for i, w := range v {
		seq[i] *= w
	}
	return seq
}
