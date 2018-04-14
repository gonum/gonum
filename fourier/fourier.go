// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fourier

import "gonum.org/v1/gonum/fourier/internal/fftpack"

// FFT implements Fast Fourier Transform and its inverse for real sequences.
type FFT struct {
	work []float64
	ifac [15]int

	// real temporarily store complex data as
	// pairs of real values to allow passing to
	// the backing code. The length of real
	// must always be half the length of work.
	real []float64
}

// NewFFT returns an FFT initialized for work on sequences of length n.
func NewFFT(n int) *FFT {
	var t FFT
	t.Reset(n)
	return &t
}

// Len returns the length of the acceptable input.
func (t *FFT) Len() int { return len(t.real) }

// Reset reinitializes the FFT for work on sequences of length n.
func (t *FFT) Reset(n int) {
	if 2*n <= cap(t.work) {
		t.work = t.work[:2*n]
		t.real = t.real[:n]
	} else {
		t.work = make([]float64, 2*n)
		t.real = make([]float64, n)
	}
	fftpack.Rffti(n, t.work, t.ifac[:])
}

// FFT computes the Fourier coefficients of the input sequence, seq,
// placing the result in dst and returning it. This transform is
// unnormalized; a call to FFT followed by a call of IFFT will
// multiply the input sequence by the length of the sequence.
//
// If the length of seq is not t.Len(), FFT will panic.
// If dst is nil, a new slice is allocated and returned. If dst is not nil and
// the length of dst does not equal t.Len()/2+1, FFT will panic.
func (t *FFT) FFT(dst []complex128, seq []float64) []complex128 {
	if len(seq) != t.Len() {
		panic("fourier: sequence length mismatch")
	}
	if dst == nil {
		dst = make([]complex128, t.Len()/2+1)
	} else if len(dst) != t.Len()/2+1 {
		panic("fourier: destination length mismatch")
	}
	copy(t.real, seq)
	fftpack.Rfftf(len(t.real), t.real, t.work, t.ifac[:])
	dst[0] = complex(t.real[0], 0)
	if len(seq) < 2 {
		return dst
	}
	if len(seq)%2 == 1 {
		dst[len(dst)-1] = complex(t.real[len(t.real)-2], t.real[len(t.real)-1])
	} else {
		dst[len(dst)-1] = complex(t.real[len(t.real)-1], 0)
	}
	for i := 1; i < len(dst)-1; i++ {
		dst[i] = complex(t.real[2*i-1], t.real[2*i])
	}
	return dst
}

// IFFT computes the real perodic sequence from the Fourier coefficients,
// coeff, placing the result in dst and returning it. This transform is
// unnormalized; a call to FFT followed by a call of IFFT will
// multiply the input sequence by the length of the sequence.
//
// If the length of coeff is not t.Len()/2+1, IFFT will panic.
// If dst is nil, a new slice is allocated and returned. If dst is not nil and
// the length of dst does not equal the length of coeff, IFFT will panic.
func (t *FFT) IFFT(dst []float64, coeff []complex128) []float64 {
	if len(coeff) != t.Len()/2+1 {
		panic("fourier: coefficients length mismatch")
	}
	if dst == nil {
		dst = make([]float64, t.Len())
	} else if len(dst) != t.Len() {
		panic("fourier: destination length mismatch")
	}
	dst[0] = real(coeff[0])
	if len(dst) < 2 {
		return dst
	}
	nf := coeff[len(coeff)-1]
	if len(dst)%2 == 1 {
		dst[len(dst)-2] = real(nf)
		dst[len(dst)-1] = imag(nf)
	} else {
		dst[len(dst)-1] = real(nf)
	}

	for i, cv := range coeff[1 : len(coeff)-1] {
		dst[2*i+1] = real(cv)
		dst[2*i+2] = imag(cv)
	}
	fftpack.Rfftb(len(dst), dst, t.work, t.ifac[:])
	return dst
}

// Freq returns the relative frequency center for coefficient i.
// Freq will panic if i is negative or greater than or equal to t.Len().
func (t *FFT) Freq(i int) float64 {
	if i < 0 || t.Len() <= i {
		panic("fourier: index out of range")
	}
	step := 1 / float64(t.Len())
	return step * float64(i)
}

// CmplxFFT implements Fast Fourier Transform and its inverse for complex sequences.
type CmplxFFT struct {
	work []float64
	ifac [15]int

	// real temporarily store complex data as
	// pairs of real values to allow passing to
	// the backing code. The length of real
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
	fftpack.Cffti(n, t.work, t.ifac[:])
}

// FFT computes the Fourier coefficients of a complex input sequence,
// seq, placing the result in dst and returning it. This transform is
// unnormalized; a call to FFT followed by a call of IFFT will
// multiply the input sequence by the length of the sequence.
//
// If the length of seq is not t.Len(), FFT will panic.
// If dst is nil, a new slice is allocated and returned. If dst is not nil and
// the length of dst does not equal the length of seq, FFT will panic.
// It is safe to use the same slice for dst and seq.
func (t *CmplxFFT) FFT(dst, seq []complex128) []complex128 {
	if len(seq) != t.Len() {
		panic("fourier: sequence length mismatch")
	}
	if dst == nil {
		dst = make([]complex128, len(seq))
	} else if len(dst) != len(seq) {
		panic("fourier: destination length mismatch")
	}
	for i, cv := range seq {
		t.real[2*i] = real(cv)
		t.real[2*i+1] = imag(cv)
	}
	fftpack.Cfftf(len(dst), t.real, t.work, t.ifac[:])
	for i := range dst {
		dst[i] = complex(t.real[2*i], t.real[2*i+1])
	}
	return dst
}

// IFFT computes the complex perodic sequence from the Fourier coefficients,
// coeff, placing the result in dst and returning it. This transform is
// unnormalized; a call to FFT followed by a call of IFFT will
// multiply the input sequence by the length of the sequence.
//
// If the length of coeff is not t.Len(), IFFT will panic.
// If dst is nil, a new slice is allocated and returned. If dst is not nil and
// the length of dst does not equal the length of coeff, IFFT will panic.
// It is safe to use the same slice for dst and coeff.
func (t *CmplxFFT) IFFT(dst, coeff []complex128) []complex128 {
	if len(coeff) != t.Len() {
		panic("fourier: coefficients length mismatch")
	}
	if dst == nil {
		dst = make([]complex128, len(coeff))
	} else if len(dst) != len(coeff) {
		panic("fourier: destination length mismatch")
	}
	for i, cv := range coeff {
		t.real[2*i] = real(cv)
		t.real[2*i+1] = imag(cv)
	}
	fftpack.Cfftb(len(dst), t.real, t.work, t.ifac[:])
	for i := range dst {
		dst[i] = complex(t.real[2*i], t.real[2*i+1])
	}
	return dst
}

// Freq returns the relative frequency center for coefficient i.
// Freq will panic if i is negative or greater than or equal to t.Len().
func (t *CmplxFFT) Freq(i int) float64 {
	if i < 0 || t.Len() <= i {
		panic("fourier: index out of range")
	}
	step := 1 / float64(t.Len())
	if i < (t.Len()-1)/2+1 {
		return step * float64(i)
	}
	return step * float64(i-t.Len())
}

// ShiftIdx returns returns a shifted index into a slice of
// coefficients returned by the CmplxFFT so that indexing
// into the coefficients places the zero frequency component
// at the center of the spectrum. ShiftIdx will panic if i is
// negative or greater than or equal to t.Len().
func (t *CmplxFFT) ShiftIdx(i int) int {
	if i < 0 || t.Len() <= i {
		panic("fourier: index out of range")
	}
	h := t.Len() / 2
	if i < h {
		return i + (t.Len()+1)/2
	}
	return i - h
}

// UnshiftIdx returns inverse of ShiftIdx. UnshiftIdx will panic if i is
// negative or greater than or equal to t.Len().
func (t *CmplxFFT) UnshiftIdx(i int) int {
	if i < 0 || t.Len() <= i {
		panic("fourier: index out of range")
	}
	h := (t.Len() + 1) / 2
	if i < h {
		return i + t.Len()/2
	}
	return i - h
}
