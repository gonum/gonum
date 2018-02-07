// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fourier

// FFT implements Fast Fourier Transform and its inverse for real sequences.
type FFT struct {
	work []float64
	ifac [15]int
}

// NewFFT returns an FFT initialized for work on sequences of length n.
func NewFFT(n int) *FFT {
	var t FFT
	t.Reset(n)
	return &t
}

// Len returns the length of the acceptable input.
func (t *FFT) Len() int { return len(t.work) / 2 }

// Reset reinitializes the FFT for work on sequences of length n.
func (t *FFT) Reset(n int) {
	if 2*n <= cap(t.work) {
		t.work = t.work[:2*n]
	} else {
		t.work = make([]float64, 2*n)
	}
	rffti(n, t.work, t.ifac[:])
}

// FFT computes the Fourier coefficients of the input sequence, seq,
// placing the result in dst and returning it. This transform is
// unnormalized since a call to FFT followed by a call of IFFT will
// multiply the input sequence by the length of the sequence.
//
// If the length of seq is not t.Len(), FFT will panic.
// If dst is nil, a new slice is allocated and returned. If dst is not nil and
// the length of dst does not equal the length of seq, FFT will panic.
func (t *FFT) FFT(dst, seq []float64) []float64 {
	if len(seq) != t.Len() {
		panic("fourier: sequence length mismatch")
	}
	if dst == nil {
		dst = make([]float64, len(seq))
	} else if len(seq) != len(dst) {
		panic("fourier: destination length mismatch")
	}
	copy(dst, seq)
	rfftf(len(dst), dst, t.work, t.ifac[:])
	return dst
}

// IFFT computes the real perodic sequence from the Fourier coefficients,
// coeff, placing the result in dst and returning it. This transform is
// unnormalized since a call to FFT followed by a call of IFFT will
// multiply the input sequence by the length of the sequence.
//
// If the length of coeff is not t.Len(), IFFT will panic.
// If dst is nil, a new slice is allocated and returned. If dst is not nil and
// the length of dst does not equal the length of coeff, IFFT will panic.
func (t *FFT) IFFT(dst, coeff []float64) []float64 {
	if len(coeff) != t.Len() {
		panic("fourier: coefficients length mismatch")
	}
	if dst == nil {
		dst = make([]float64, len(coeff))
	} else if len(coeff) != len(dst) {
		panic("fourier: destination length mismatch")
	}
	copy(dst, coeff)
	rfftb(len(dst), dst, t.work, t.ifac[:])
	return dst
}

// CmplxFFT implements Fast Fourier Transform and its inverse for real sequences.
type CmplxFFT struct {
	work []float64
	ifac [15]int

	// real temporarily store complex data as
	// pairs of real values to allow passing to
	// the backing code. The length of conv
	// must always be half the length of work.
	real []float64
}

// NewCmplxFFT returns an CmplxFFT initialized for work on sequences of length n.
func NewCmplxFFT(n int) *CmplxFFT {
	var t CmplxFFT
	t.Reset(n)
	return &t
}

// Len returns the length of the acceptable input.
func (t *CmplxFFT) Len() int { return len(t.work) / 4 }

// Reset reinitializes the FFT for work on sequences of length n.
func (t *CmplxFFT) Reset(n int) {
	if 4*n <= cap(t.work) {
		t.work = t.work[:4*n]
		t.real = t.real[:2*n]
	} else {
		t.work = make([]float64, 4*n)
		t.real = make([]float64, 2*n)
	}
	cffti(n, t.work, t.ifac[:])
}

// FFT computes the Fourier coefficients of a complex input sequence,
// seq, placing the result in dst and returning it. This transform is
// unnormalized since a call to FFT followed by a call of IFFT will
// multiply the input sequence by the length of the sequence.
//
// If the length of seq is not t.Len(), FFT will panic.
// If dst is nil, a new slice is allocated and returned. If dst is not nil and
// the length of dst does not equal the length of seq, FFT will panic.
func (t *CmplxFFT) FFT(dst, seq []complex128) []complex128 {
	if len(seq) != t.Len() {
		panic("fourier: sequence length mismatch")
	}
	if dst == nil {
		dst = make([]complex128, len(seq))
	} else if len(seq) != len(dst) {
		panic("fourier: destination length mismatch")
	}
	for i, cv := range seq {
		t.real[2*i] = real(cv)
		t.real[2*i+1] = imag(cv)
	}
	cfftf(len(dst), t.real, t.work, t.ifac[:])
	for i := range dst {
		dst[i] = complex(t.real[2*i], t.real[2*i+1])
	}
	return dst
}

// IFFT computes the complex perodic sequence from the Fourier coefficients,
// coeff, placing the result in dst and returning it. This transform is
// unnormalized since a call to FFT followed by a call of IFFT will
// multiply the input sequence by the length of the sequence.
//
// If the length of coeff is not t.Len(), IFFT will panic.
// If dst is nil, a new slice is allocated and returned. If dst is not nil and
// the length of dst does not equal the length of coeff, IFFT will panic.
func (t *CmplxFFT) IFFT(dst, coeff []complex128) []complex128 {
	if len(coeff) != t.Len() {
		panic("fourier: coefficients length mismatch")
	}
	if dst == nil {
		dst = make([]complex128, len(coeff))
	} else if len(coeff) != len(dst) {
		panic("fourier: destination length mismatch")
	}
	for i, cv := range coeff {
		t.real[2*i] = real(cv)
		t.real[2*i+1] = imag(cv)
	}
	cfftb(len(dst), t.real, t.work, t.ifac[:])
	for i := range dst {
		dst[i] = complex(t.real[2*i], t.real[2*i+1])
	}
	return dst
}
