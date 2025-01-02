// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fourier

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"
	"unsafe"

	"gonum.org/v1/gonum/cmplxs"
)

func TestCoefficients(t *testing.T) {
	const tol = 1e-8

	src := rand.NewPCG(1, 1)
	for n := 4; n < 1<<20; n <<= 1 {
		for i := 0; i < 10; i++ {
			t.Run(fmt.Sprintf("Radix2/%d", n), func(t *testing.T) {
				d := randComplexes(n, src)
				fft := NewCmplxFFT(n)
				want := fft.Coefficients(nil, d)
				CoefficientsRadix2(d)
				got := d
				if !cmplxs.EqualApprox(got, want, tol) {
					t.Errorf("unexpected result for n=%d |got-want|^2=%g", n, cmplxs.Distance(got, want, 2))
				}

				want = fft.Sequence(nil, got)
				scale(1/float64(n), want)

				SequenceRadix2(got)
				scale(1/float64(n), got)

				if !cmplxs.EqualApprox(got, want, tol) {
					t.Errorf("unexpected ifft result for n=%d |got-want|^2=%g", n, cmplxs.Distance(got, want, 2))
				}
			})
			if bits.Len(uint(n))&0x1 == 0 {
				continue
			}
			t.Run(fmt.Sprintf("Radix4/%d", n), func(t *testing.T) {
				d := randComplexes(n, src)
				fft := NewCmplxFFT(n)
				want := fft.Coefficients(nil, d)
				CoefficientsRadix4(d)
				got := d
				if !cmplxs.EqualApprox(got, want, tol) {
					t.Errorf("unexpected fft result for n=%d |got-want|^2=%g", n, cmplxs.Distance(got, want, 2))
				}

				want = fft.Sequence(nil, got)
				scale(1/float64(n), want)

				SequenceRadix4(got)
				scale(1/float64(n), got)

				if !cmplxs.EqualApprox(got, want, tol) {
					t.Errorf("unexpected ifft result for n=%d |got-want|^2=%g", n, cmplxs.Distance(got, want, 2))
				}
			})
		}
	}
}

func TestSequence(t *testing.T) {
	const tol = 1e-10

	src := rand.NewPCG(1, 1)
	for n := 4; n < 1<<20; n <<= 1 {
		for i := 0; i < 10; i++ {
			t.Run(fmt.Sprintf("Radix2/%d", n), func(t *testing.T) {
				d := randComplexes(n, src)
				want := make([]complex128, n)
				copy(want, d)
				SequenceRadix2(CoefficientsRadix2(d))
				got := d

				scale(1/float64(n), got)

				if !cmplxs.EqualApprox(got, want, tol) {
					t.Errorf("unexpected result for ifft(fft(d)) n=%d |got-want|^2=%g", n, cmplxs.Distance(got, want, 2))
				}
			})
			if bits.Len(uint(n))&0x1 == 0 {
				continue
			}
			t.Run(fmt.Sprintf("Radix4/%d", n), func(t *testing.T) {
				d := randComplexes(n, src)
				want := make([]complex128, n)
				copy(want, d)
				SequenceRadix4(CoefficientsRadix4(d))
				got := d

				scale(1/float64(n), got)

				if !cmplxs.EqualApprox(got, want, tol) {
					t.Errorf("unexpected result for ifft(fft(d)) n=%d |got-want|^2=%g", n, cmplxs.Distance(got, want, 2))
				}
			})
		}
	}
}

func scale(f float64, c []complex128) {
	for i, v := range c {
		c[i] = complex(f*real(v), f*imag(v))
	}
}

func TestBitReversePermute(t *testing.T) {
	for n := 2; n <= 1024; n <<= 1 {
		x := make([]complex128, n)
		for i := range x {
			x[i] = complex(float64(i), float64(i))
		}
		bitReversePermute(x)
		for i, got := range x {
			j := bits.Reverse(uint(i)) >> bits.LeadingZeros(uint(n-1))
			want := complex(float64(j), float64(j))
			if got != want {
				t.Errorf("unexpected value at %d: got:%f want:%f", i, got, want)
			}
		}
	}
}

func TestPadRadix2(t *testing.T) {
	for n := 1; n <= 1025; n++ {
		x := make([]complex128, n)
		y := PadRadix2(x)
		if bits.OnesCount(uint(len(y))) != 1 {
			t.Errorf("unexpected length of padded slice: not a power of 2: %d", len(y))
		}
		if len(x) == len(y) && &y[0] != &x[0] {
			t.Errorf("unexpected new allocation for power of 2 input length: len(x)=%d", n)
		}
		if len(y) < len(x) {
			t.Errorf("unexpected short result: len(y)=%d < len(x)=%d", len(y), len(x))
		}
	}
}

func TestTrimRadix2(t *testing.T) {
	for n := 1; n <= 1025; n++ {
		x := make([]complex128, n)
		y, r := TrimRadix2(x)
		if bits.OnesCount(uint(len(y))) != 1 {
			t.Errorf("unexpected length of trimmed slice: not a power of 2: %d", len(y))
		}
		if len(y)+len(r) != len(x) {
			t.Errorf("unexpected total result: len(y)=%d + len(r)%d != len(x)=%d", len(y), len(r), len(x))
		}
		if len(x) == len(y) && &y[0] != &x[0] {
			t.Errorf("unexpected new allocation for power of 2 input length: len(x)=%d", n)
		}
		if len(y) > len(x) {
			t.Errorf("unexpected long result: len(y)=%d > len(x)=%d", len(y), len(x))
		}
	}
}

func TestBitPairReversePermute(t *testing.T) {
	for n := 4; n <= 1024; n <<= 2 {
		x := make([]complex128, n)
		for i := range x {
			x[i] = complex(float64(i), float64(i))
		}
		bitPairReversePermute(x)
		for i, got := range x {
			j := reversePairs(uint(i)) >> bits.LeadingZeros(uint(n-1))
			want := complex(float64(j), float64(j))
			if got != want {
				t.Errorf("unexpected value at %d: got:%f want:%f", i, got, want)
			}
		}
	}
}

func TestReversePairs(t *testing.T) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for i := 0; i < 1000; i++ {
		x := uint(rnd.Uint64())
		got := reversePairs(x)
		want := naiveReversePairs(x)
		if got != want {
			t.Errorf("unexpected bit-pair reversal for 0b%064b:\ngot: 0b%064b\nwant:0b%064b", x, got, want)
		}
	}
}

// naiveReversePairs does a bit-pair reversal by formatting as a base-4 string,
// reversing the digits of the formatted number and then re-parsing the value.
func naiveReversePairs(x uint) uint {
	bits := int(unsafe.Sizeof(uint(0)) * 8)

	// Format the number as a quaternary, padded with zeros.
	// We avoid the leftpad issue by doing it ourselves.
	b := strconv.AppendUint(bytes.Repeat([]byte("0"), bits/2), uint64(x), 4)
	b = b[len(b)-bits/2:]

	// Reverse the quits.
	slices.Reverse(b)

	// Re-parse the reversed number.
	y, err := strconv.ParseUint(string(b), 4, 64)
	if err != nil {
		panic(fmt.Sprintf("unexpected parse error: %v", err))
	}
	return uint(y)
}

func TestPadRadix4(t *testing.T) {
	for n := 1; n <= 1025; n++ {
		x := make([]complex128, n)
		y := PadRadix4(x)
		if bits.OnesCount(uint(len(y))) != 1 || bits.Len(uint(len(y)))&0x1 == 0 {
			t.Errorf("unexpected length of padded slice: not a power of 4: %d", len(y))
		}
		if len(x) == len(y) && &y[0] != &x[0] {
			t.Errorf("unexpected new allocation for power of 2 input length: len(x)=%d", n)
		}
		if len(y) < len(x) {
			t.Errorf("unexpected short result: len(y)=%d < len(x)=%d", len(y), len(x))
		}
	}
}

func TestTrimRadix4(t *testing.T) {
	for n := 1; n <= 1025; n++ {
		x := make([]complex128, n)
		y, r := TrimRadix4(x)
		if bits.OnesCount(uint(len(y))) != 1 || bits.Len(uint(len(y)))&0x1 == 0 {
			t.Errorf("unexpected length of trimmed slice: not a power of 4: %d", len(y))
		}
		if len(y)+len(r) != len(x) {
			t.Errorf("unexpected total result: len(y)=%d + len(r)%d != len(x)=%d", len(y), len(r), len(x))
		}
		if len(x) == len(y) && &y[0] != &x[0] {
			t.Errorf("unexpected new allocation for power of 2 input length: len(x)=%d", n)
		}
		if len(y) > len(x) {
			t.Errorf("unexpected long result: len(y)=%d > len(x)=%d", len(y), len(x))
		}
	}
}

func BenchmarkCoefficients(b *testing.B) {
	for n := 16; n < 1<<24; n <<= 3 {
		d := randComplexes(n, rand.NewPCG(1, 1))
		b.Run(fmt.Sprintf("Radix2/%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				CoefficientsRadix2(d)
			}
		})
		if bits.Len(uint(n))&0x1 == 0 {
			continue
		}
		b.Run(fmt.Sprintf("Radix4/%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				CoefficientsRadix4(d)
			}
		})
	}
}

func BenchmarkSequence(b *testing.B) {
	for n := 16; n < 1<<24; n <<= 3 {
		d := randComplexes(n, rand.NewPCG(1, 1))
		b.Run(fmt.Sprintf("Radix2/%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				SequenceRadix2(d)
			}
		})
		if bits.Len(uint(n))&0x1 == 0 {
			continue
		}
		b.Run(fmt.Sprintf("Radix4/%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				SequenceRadix4(d)
			}
		})
	}
}
