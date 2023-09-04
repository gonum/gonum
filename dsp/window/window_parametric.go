// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package window

import "math"

// Gaussian can modify a sequence using the Gaussian window and return the
// result.
// See https://en.wikipedia.org/wiki/Window_function#Gaussian_window
// and https://www.recordingblogs.com/wiki/gaussian-window for details.
//
// The Gaussian window is an adjustable window.
//
// The sequence weights are
//
//	w[k] = exp(-0.5 * ((k-M)/(σ*M))² ), M = (N-1)/2,
//
// for k=0,1,...,N-1 where N is the length of the window.
//
// The properties of the window depend on the value of σ (sigma).
// It can be used as high or low resolution window, depending of the σ value.
//
// Spectral leakage parameters are summarized in the table:
//
//	       |  σ=0.3  |  σ=0.5 |  σ=1.2 |
//	-------|---------------------------|
//	ΔF_0   |   8     |   3.4  |   2.2  |
//	ΔF_0.5 |   1.82  |   1.2  |   0.94 |
//	K      |   4     |   1.7  |   1.1  |
//	ɣ_max  | -65     | -31.5  | -15.5  |
//	β      |  -8.52  |  -4.48 |  -0.96 |
type Gaussian struct {
	Sigma float64
}

// Transform applies the Gaussian transformation to seq in place, using the
// value of the receiver as the sigma parameter, and returning the result.
func (g Gaussian) Transform(seq []float64) []float64 {
	a := float64(len(seq)-1) / 2
	for i := range seq {
		x := -0.5 * math.Pow((float64(i)-a)/(g.Sigma*a), 2)
		seq[i] *= math.Exp(x)
	}
	return seq
}

// TransformComplex applies the Gaussian transformation to seq in place,
// using the value of the receiver as the sigma parameter, and returning
// the result.
func (g Gaussian) TransformComplex(seq []complex128) []complex128 {
	a := float64(len(seq)-1) / 2
	for i, v := range seq {
		x := -0.5 * math.Pow((float64(i)-a)/(g.Sigma*a), 2)
		w := math.Exp(x)
		seq[i] = complex(w*real(v), w*imag(v))
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
//
//	w[k] = 0.5 * (1 + cos(π*(|k - M| - αM)/((1-α) * M))), |k - M| ≥ αM
//	     = 1, |k - M| < αM
//
// with M = (N - 1)/2 for k=0,1,...,N-1 where N is the length of the window.
//
// Spectral leakage parameters are summarized in the table:
//
//	       |  α=0.3 |  α=0.5 |  α=0.7 |
//	-------|--------------------------|
//	ΔF_0   |   1.33 |   1.22 |   1.13 |
//	ΔF_0.5 |   1.28 |   1.16 |   1.04 |
//	K      |   0.67 |   0.61 |   0.57 |
//	ɣ_max  | -18.2  | -15.1  | -13.8  |
//	β      |  -1.41 |  -2.50 |  -3.74 |
type Tukey struct {
	Alpha float64
}

// Transform applies the Tukey transformation to seq in place, using the
// value of the receiver as the Alpha parameter, and returning the result.
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

// TransformComplex applies the Tukey transformation to seq in place, using
// the value of the receiver as the Alpha parameter, and returning the result.
func (t Tukey) TransformComplex(seq []complex128) []complex128 {
	switch {
	case t.Alpha <= 0:
		return RectangularComplex(seq)
	case t.Alpha >= 1:
		return HannComplex(seq)
	default:
		alphaL := t.Alpha * float64(len(seq)-1)
		width := int(0.5*alphaL) + 1
		for i, v := range seq[:width] {
			w := 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/alphaL))
			v = complex(w*real(v), w*imag(v))
			seq[i] = v
			seq[len(seq)-1-i] = v
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

// TransformTo applies the weights in the receiver to src placing the result
// in dst. If v is nil, TransformTo is a no-op, otherwise the length of v must
// match the length of src and dst.
func (v Values) TransformTo(dst, src []float64) {
	if v == nil {
		return
	}
	if len(v) != len(src) {
		panic("window: seq length mismatch")
	}
	if len(v) != len(dst) {
		panic("window: dst length mismatch")
	}
	for i, w := range v {
		dst[i] = w * src[i]
	}
}

// TransformComplex applies the weights in the receiver to seq in place,
// returning the result. If v is nil, TransformComplex is a no-op, otherwise
// the length of v must match the length of seq.
func (v Values) TransformComplex(seq []complex128) []complex128 {
	if v == nil {
		return seq
	}
	if len(v) != len(seq) {
		panic("window: length mismatch")
	}
	for i, w := range v {
		sv := seq[i]
		seq[i] = complex(w*real(sv), w*imag(sv))
	}
	return seq
}

// TransformComplexTo applies the weights in the receiver to src placing the
// result in dst. If v is nil, TransformComplexTo is a no-op, otherwise the
// length of v must match the length of src and dst.
func (v Values) TransformComplexTo(dst, src []complex128) {
	if v == nil {
		return
	}
	if len(v) != len(src) {
		panic("window: seq length mismatch")
	}
	if len(v) != len(dst) {
		panic("window: dst length mismatch")
	}
	for i, w := range v {
		sv := src[i]
		dst[i] = complex(w*real(sv), w*imag(sv))
	}
}
