// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transform

import (
	"math/cmplx"

	"gonum.org/v1/gonum/dsp/fourier"
)

// Hilbert implements an approximate Hilbert transform that allows calculation
// of an approximate analytical signal of a real signal, and determine the
// real envelope of a signal.
//
// The underlying implementation uses a discrete Fourier transform and inverse
// and discrete Fourier transform, so usual methods for speeding up these
// transforms are likely to apply here as well.
type Hilbert struct {
	fft  *fourier.CmplxFFT
	work []complex128
}

// NewHilbert returns a new Hilbert transformer for signals of size N.
func NewHilbert(n int) *Hilbert {
	return &Hilbert{
		fft:  fourier.NewCmplxFFT(n),
		work: make([]complex128, n),
	}
}

// Retrieve the length N of signals that are valid input for this Hilbert transform.
func (h *Hilbert) Len() int {
	return len(h.work)
}

// Compute the analytical signal of a real signal, and stores the result in the dst slice,
// returning it.
// If the dst slice is nil, a new slice will be created and returned. The dst slice
// must be the same length as the input signal.
func (h *Hilbert) AnalyticSignal(signal []float64, dst []complex128) []complex128 {
	if dst == nil {
		dst = make([]complex128, len(signal))
	}

	for i, v := range signal {
		h.work[i] = complex(v, 0)
	}

	// forward fft of signal
	coeff := h.fft.Coefficients(dst, h.work)
	for i := range h.work {
		h.work[i] = 0
	}

	// multiply positive frequencies by 2, zero out negative frequencies
	// however, leave dc unchanged (and nyquist when n%2 == 0)
	h.work[0] = coeff[0]
	for i, d := range coeff[1 : len(coeff)/2+1] {
		h.work[i+1] = d * 2
	}
	if len(coeff)%2 == 0 {
		h.work[len(coeff)/2] = coeff[len(coeff)/2]
	}

	// normalize the results for now
	// TODO: is this important for analytic signal or envelope?
	unnorm := h.fft.Sequence(dst, h.work)
	for i, u := range unnorm {
		unnorm[i] = u / complex(float64(len(unnorm)), 0)
	}
	return unnorm
}

// Compute the positive envelope of a real signal, and stores the result in the dst slice,
// returning it.
// If the dst slice is nil, a new slice will be created and returned. The dst slice
// must be the same length as the input signal.
func (h *Hilbert) Envelope(signal []float64, dst []float64) []float64 {
	if dst == nil {
		dst = make([]float64, len(signal))
	}

	analytic := h.AnalyticSignal(signal, nil)
	for i, a := range analytic {
		dst[i] = cmplx.Abs(a)
	}
	return dst
}
