// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fourier

import (
	"math"
	"math/bits"
	"math/cmplx"
)

// CoefficientsRadix2 computes the Fourier coefficients of the input
// sequence, converting the time series in seq into the frequency spectrum,
// in place and returning it. This transform is unnormalized; a call to
// CoefficientsRadix2 followed by a call of SequenceRadix2 will multiply the
// input sequence by the length of the sequence.
//
// CoefficientsRadix2 does not allocate, requiring that FFT twiddle factors
// be calculated lazily. For performance reasons, this is done by successive
// multiplication, so numerical accuracies can accumulate for large inputs.
// If accuracy is needed, the FFT or CmplxFFT types should be used.
//
// If the length of seq is not an integer power of 2, CoefficientsRadix2 will
// panic.
func CoefficientsRadix2(seq []complex128) []complex128 {
	x := seq
	switch len(x) {
	default:
		if bits.OnesCount(uint(len(x))) != 1 {
			panic("fourier: radix-2 fft called with non-power 2 length")
		}

	case 0, 1:
		return x

	case 2:
		x[0], x[1] =
			x[0]+x[1],
			x[0]-x[1]
		return x

	case 4:
		t := x[1] + x[3]
		u := x[2]
		v := negi(x[1] - x[3])
		x[0], x[1], x[2], x[3] =
			x[0]+u+t,
			x[0]-u+v,
			x[0]+u-t,
			x[0]-u-v
		return x
	}

	bitReversePermute(x)

	for k := 0; k < len(x); k += 4 {
		t := x[k+2] + x[k+3]
		u := x[k+1]
		v := negi(x[k+2] - x[k+3])
		x[k], x[k+1], x[k+2], x[k+3] =
			x[k]+u+t,
			x[k]-u+v,
			x[k]+u-t,
			x[k]-u-v
	}

	for m := 4; m < len(x); m *= 2 {
		f := swap(complex(math.Sincos(-math.Pi / float64(m))))
		for k := 0; k < len(x); k += 2 * m {
			w := 1 + 0i
			for j := 0; j < m; j++ {
				i := j + k

				u := w * x[i+m]
				x[i], x[i+m] =
					x[i]+u,
					x[i]-u

				w *= f
			}
		}
	}

	return x
}

// bitReversePermute performs a bit-reversal permutation on x.
func bitReversePermute(x []complex128) {
	if len(x) < 2 || bits.OnesCount(uint(len(x))) != 1 {
		panic("fourier: invalid bitReversePermute call")
	}
	lz := bits.LeadingZeros(uint(len(x) - 1))
	i := 0
	for ; i < len(x)/2; i++ {
		j := int(bits.Reverse(uint(i)) >> lz)
		if i < j {
			x[i], x[j] = x[j], x[i]
		}
	}
	for i++; i < len(x); i += 2 {
		j := int(bits.Reverse(uint(i)) >> lz)
		if i < j {
			x[i], x[j] = x[j], x[i]
		}
	}
}

// SequenceRadix2 computes the real periodic sequence from the Fourier
// coefficients, converting the frequency spectrum in coeff into a time
// series, in place and returning it. This transform is unnormalized; a
// call to CoefficientsRadix2 followed by a call of SequenceRadix2 will
// multiply the input sequence by the length of the sequence.
//
// SequenceRadix2 does not allocate, requiring that FFT twiddle factors
// be calculated lazily. For performance reasons, this is done by successive
// multiplication, so numerical accuracies can accumulate for large inputs.
// If accuracy is needed, the FFT or CmplxFFT types should be used.
//
// If the length of coeff is not an integer power of 2, SequenceRadix2
// will panic.
func SequenceRadix2(coeff []complex128) []complex128 {
	x := coeff
	for i, j := 1, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}

	CoefficientsRadix2(x)
	return x
}

// PadRadix2 returns the values in x in a slice that is an integer
// power of 2 long. If x already has an integer power of 2 length
// it is returned unaltered.
func PadRadix2(x []complex128) []complex128 {
	if len(x) == 0 {
		return x
	}
	b := bits.Len(uint(len(x)))
	if len(x) == 1<<(b-1) {
		return x
	}
	p := make([]complex128, 1<<b)
	copy(p, x)
	return p
}

// TrimRadix2 returns the largest slice of x that is has an integer
// power of 2 length, and a slice holding the remaining elements.
func TrimRadix2(x []complex128) (even, remains []complex128) {
	if len(x) == 0 {
		return x, nil
	}
	n := 1 << (bits.Len(uint(len(x))) - 1)
	return x[:n], x[n:]
}

// CoefficientsRadix4 computes the Fourier coefficients of the input
// sequence, converting the time series in seq into the frequency spectrum,
// in place and returning it. This transform is unnormalized; a call to
// CoefficientsRadix4 followed by a call of SequenceRadix4 will multiply the
// input sequence by the length of the sequence.
//
// CoefficientsRadix4 does not allocate, requiring that FFT twiddle factors
// be calculated lazily. For performance reasons, this is done by successive
// multiplication, so numerical accuracies can accumulate for large inputs.
// If accuracy is needed, the FFT or CmplxFFT types should be used.
//
// If the length of seq is not an integer power of 4, CoefficientsRadix4 will
// panic.
func CoefficientsRadix4(seq []complex128) []complex128 {
	x := seq
	switch len(x) {
	default:
		if bits.OnesCount(uint(len(x))) != 1 || bits.TrailingZeros(uint(len(x)))&0x1 != 0 {
			panic("fourier: radix-4 fft called with non-power 4 length")
		}

	case 0, 1:
		return x

	case 4:
		t := x[1] + x[3]
		u := x[2]
		v := negi(x[1] - x[3])
		x[0], x[1], x[2], x[3] =
			x[0]+u+t,
			x[0]-u+v,
			x[0]+u-t,
			x[0]-u-v
		return x
	}

	bitPairReversePermute(x)

	for k := 0; k < len(x); k += 4 {
		t := x[k+1] + x[k+3]
		u := x[k+2]
		v := negi(x[k+1] - x[k+3])
		x[k], x[k+1], x[k+2], x[k+3] =
			x[k]+u+t,
			x[k]-u+v,
			x[k]+u-t,
			x[k]-u-v
	}

	for m := 4; m < len(x); m *= 4 {
		f := swap(complex(math.Sincos((-math.Pi / 2) / float64(m))))
		for k := 0; k < len(x); k += m * 4 {
			w := 1 + 0i
			w2 := w
			w3 := w2
			for j := 0; j < m; j++ {
				i := j + k

				t := x[i+m]*w + x[i+3*m]*w3
				u := x[i+2*m] * w2
				v := negi(x[i+m]*w - x[i+3*m]*w3)
				x[i], x[i+m], x[i+2*m], x[i+3*m] =
					x[i]+u+t,
					x[i]-u+v,
					x[i]+u-t,
					x[i]-u-v

				wt := f
				w *= wt
				wt *= f
				w2 *= wt
				wt *= f
				w3 *= wt
			}
		}
	}

	return x
}

// bitPairReversePermute performs a bit pair-reversal permutation on x.
func bitPairReversePermute(x []complex128) {
	if len(x) < 4 || bits.OnesCount(uint(len(x))) != 1 || bits.TrailingZeros(uint(len(x)))&0x1 != 0 {
		panic("fourier: invalid bitPairReversePermute call")
	}
	lz := bits.LeadingZeros(uint(len(x) - 1))
	i := 0
	for ; i < 3*len(x)/4; i++ {
		j := int(reversePairs(uint(i)) >> lz)
		if i < j {
			x[i], x[j] = x[j], x[i]
		}
	}
	for i++; i < len(x); i += 2 {
		j := int(reversePairs(uint(i)) >> lz)
		if i < j {
			x[i], x[j] = x[j], x[i]
		}
	}
}

// SequenceRadix4 computes the real periodic sequence from the Fourier
// coefficients, converting the frequency spectrum in coeff into a time
// series, in place and returning it. This transform is unnormalized; a
// call to CoefficientsRadix4 followed by a call of SequenceRadix4 will
// multiply the input sequence by the length of the sequence.
//
// SequenceRadix4 does not allocate, requiring that FFT twiddle factors
// be calculated lazily. For performance reasons, this is done by successive
// multiplication, so numerical accuracies can accumulate for large inputs.
// If accuracy is needed, the FFT or CmplxFFT types should be used.
//
// If the length of coeff is not an integer power of 4, SequenceRadix4
// will panic.
func SequenceRadix4(coeff []complex128) []complex128 {
	x := coeff
	for i, j := 1, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}

	CoefficientsRadix4(x)
	return x
}

// PadRadix4 returns the values in x in a slice that is an integer
// power of 4 long. If x already has an integer power of 4 length
// it is returned unaltered.
func PadRadix4(x []complex128) []complex128 {
	if len(x) == 0 {
		return x
	}
	b := bits.Len(uint(len(x)))
	if len(x) == 1<<(b-1) && b&0x1 == 1 {
		return x
	}
	p := make([]complex128, 1<<((b+1)&^0x1))
	copy(p, x)
	return p
}

// TrimRadix4 returns the largest slice of x that is has an integer
// power of 4 length, and a slice holding the remaining elements.
func TrimRadix4(x []complex128) (even, remains []complex128) {
	if len(x) == 0 {
		return x, nil
	}
	n := 1 << ((bits.Len(uint(len(x))) - 1) &^ 0x1)
	return x[:n], x[n:]
}

// reversePairs returns the value of x with its bit pairs in reversed order.
func reversePairs(x uint) uint {
	if bits.UintSize == 32 {
		return uint(reversePairs32(uint32(x)))
	}
	return uint(reversePairs64(uint64(x)))
}

const (
	m1 = 0x3333333333333333
	m2 = 0x0f0f0f0f0f0f0f0f
)

// reversePairs32 returns the value of x with its bit pairs in reversed order.
func reversePairs32(x uint32) uint32 {
	const m = 1<<32 - 1
	x = x>>2&(m1&m) | x&(m1&m)<<2
	x = x>>4&(m2&m) | x&(m2&m)<<4
	return bits.ReverseBytes32(x)
}

// reversePairs64 returns the value of x with its bit pairs in reversed order.
func reversePairs64(x uint64) uint64 {
	const m = 1<<64 - 1
	x = x>>2&(m1&m) | x&(m1&m)<<2
	x = x>>4&(m2&m) | x&(m2&m)<<4
	return bits.ReverseBytes64(x)
}

func negi(c complex128) complex128 {
	return cmplx.Conj(swap(c))
}

func swap(c complex128) complex128 {
	return complex(imag(c), real(c))
}
