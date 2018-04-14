// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fourier

import (
	"reflect"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats"
)

func TestFFT(t *testing.T) {
	const tol = 1e-10
	t.Run("NewFFT", func(t *testing.T) {
		for n := 1; n <= 200; n++ {
			fft := NewFFT(n)

			want := make([]float64, n)
			for i := range want {
				want[i] = rand.Float64()
			}

			coeff := fft.FFT(nil, want)
			got := fft.IFFT(nil, coeff)
			floats.Scale(1/float64(n), got)

			if !floats.EqualApprox(got, want, tol) {
				t.Errorf("unexpected result for ifft(fft(x)) for length %d", n)
			}
		}
	})
	t.Run("Reset FFT", func(t *testing.T) {
		fft := NewFFT(1000)
		for n := 1; n <= 2000; n++ {
			fft.Reset(n)

			want := make([]float64, n)
			for i := range want {
				want[i] = rand.Float64()
			}

			coeff := fft.FFT(nil, want)
			got := fft.IFFT(nil, coeff)
			floats.Scale(1/float64(n), got)

			if !floats.EqualApprox(got, want, tol) {
				t.Errorf("unexpected result for ifft(fft(x)) for length %d", n)
			}
		}
	})
	t.Run("known FFT", func(t *testing.T) {
		// Values confirmed with reference to numpy rfft.
		fft := NewFFT(1000)
		cases := []struct {
			in   []float64
			want []complex128
		}{
			{
				in:   []float64{1, 0, 1, 0, 1, 0, 1, 0},
				want: []complex128{4, 0, 0, 0, 4},
			},
			{
				in: []float64{1, 0, 1, 0, 1, 0, 1},
				want: []complex128{
					4,
					0.5 + 0.24078730940376442i,
					0.5 + 0.6269801688313512i,
					0.5 + 2.190643133767413i,
				},
			},
			{
				in: []float64{1, 0, 2, 0, 1, 0, 4, 0, 1, 0, 2, 0, 1, 0},
				want: []complex128{
					12,
					-2.301937735804838 - 1.108554787638881i,
					0.7469796037174659 + 0.9366827961047095i,
					-0.9450418679126271 - 4.140498958131061i,
					-0.9450418679126271 + 4.140498958131061i,
					0.7469796037174659 - 0.9366827961047095i,
					-2.301937735804838 + 1.108554787638881i,
					12,
				},
			},
		}
		for _, test := range cases {
			fft.Reset(len(test.in))
			got := fft.FFT(nil, test.in)
			if !equalApprox(got, test.want, tol) {
				t.Errorf("unexpected result for fft(%g):\ngot: %g\nwant:%g",
					test.in, got, test.want)
			}
		}
	})
	t.Run("Freq", func(t *testing.T) {
		var fft FFT
		cases := []struct {
			n    int
			want []float64
		}{
			{n: 1, want: []float64{0}},
			{n: 2, want: []float64{0, 0.5}},
			{n: 3, want: []float64{0, 1.0 / 3.0}},
			{n: 4, want: []float64{0, 0.25, 0.5}},
		}
		for _, test := range cases {
			fft.Reset(test.n)
			for i, want := range test.want {
				if got := fft.Freq(i); got != want {
					t.Errorf("unexpected result for freq(%d) for length %d: got:%v want:%v",
						i, test.n, got, want)
				}
			}
		}
	})
}

func TestCmplxFFT(t *testing.T) {
	const tol = 1e-12
	t.Run("NewFFT", func(t *testing.T) {
		for n := 1; n <= 200; n++ {
			fft := NewCmplxFFT(n)

			want := make([]complex128, n)
			for i := range want {
				want[i] = complex(rand.Float64(), rand.Float64())
			}

			coeff := fft.FFT(nil, want)
			got := fft.IFFT(nil, coeff)
			sf := complex(1/float64(n), 0)
			for i := range got {
				got[i] *= sf
			}

			if !equalApprox(got, want, tol) {
				t.Errorf("unexpected result for complex ifft(fft(x)) for length %d", n)
			}
		}
	})
	t.Run("Reset FFT", func(t *testing.T) {
		fft := NewCmplxFFT(1000)
		for n := 1; n <= 2000; n++ {
			fft.Reset(n)

			want := make([]complex128, n)
			for i := range want {
				want[i] = complex(rand.Float64(), rand.Float64())
			}

			coeff := fft.FFT(nil, want)
			got := fft.IFFT(nil, coeff)
			sf := complex(1/float64(n), 0)
			for i := range got {
				got[i] *= sf
			}

			if !equalApprox(got, want, tol) {
				t.Errorf("unexpected result for complex ifft(fft(x)) for length %d", n)
			}
		}
	})
	t.Run("Freq", func(t *testing.T) {
		var fft CmplxFFT
		cases := []struct {
			want []float64
		}{
			{want: []float64{0}},
			{want: []float64{0, -0.5}},
			{want: []float64{0, 1.0 / 3.0, -1.0 / 3.0}},
			{want: []float64{0, 0.25, -0.5, -0.25}},
		}
		for _, test := range cases {
			fft.Reset(len(test.want))
			for i, want := range test.want {
				if got := fft.Freq(i); got != want {
					t.Errorf("unexpected result for freq(%d) for length %d: got:%v want:%v",
						i, len(test.want), got, want)
				}
			}
		}
	})
	t.Run("Shift", func(t *testing.T) {
		var fft CmplxFFT
		cases := []struct {
			index []int
			want  []int
		}{
			{index: []int{0}, want: []int{0}},
			{index: []int{0, -1}, want: []int{-1, 0}},
			{index: []int{0, 1, -1}, want: []int{-1, 0, 1}},
			{index: []int{0, 1, -2, -1}, want: []int{-2, -1, 0, 1}},
			{index: []int{0, 1, 2, -2, -1}, want: []int{-2, -1, 0, 1, 2}},
		}
		for _, test := range cases {
			fft.Reset(len(test.index))
			got := make([]int, len(test.index))
			for i := range test.index {
				got[i] = test.index[fft.ShiftIdx(i)]
				su := fft.UnshiftIdx(fft.ShiftIdx(i))
				if su != i {
					t.Errorf("unexpected result for unshift(shift(%d)) with length %d:\ngot: %d\nwant:%d",
						i, len(test.index), su, i)
				}
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("unexpected result for shift(%d):\ngot: %d\nwant:%d",
					test.index, got, test.want)
			}
		}
	})
}

func equalApprox(a, b []complex128, tol float64) bool {
	if len(a) != len(b) {
		return false
	}
	ar := make([]float64, len(a))
	br := make([]float64, len(a))
	ai := make([]float64, len(a))
	bi := make([]float64, len(a))
	for i, cv := range a {
		ar[i] = real(cv)
		ai[i] = imag(cv)
	}
	for i, cv := range b {
		br[i] = real(cv)
		bi[i] = imag(cv)
	}
	return floats.EqualApprox(ar, br, tol) && floats.EqualApprox(ai, bi, tol)
}
