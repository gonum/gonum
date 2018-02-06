// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fourier

import "fmt"

// The types in array.go implement Fortran-like arrays for bootstrapping
// the implementation of the FFT functions translated from FFTPACK; they
// are column-major one-based.
// TODO(kortschak): Make types zero-based and remove linear array types.
// TODO(kortschak): Make a bounds-checking and non-bounds-checking version
// for testing/production.

type oneArray []float64

func (a oneArray) at(i int) float64 {
	i--
	if i < 0 || len(a) <= i {
		panic(fmt.Sprintf("out of bounds at(%d): bounds i=%d", i, len(a)))
	}
	return a[i]
}

func (a oneArray) set(i int, v float64) {
	i--
	if i < 0 || len(a) <= i {
		panic(fmt.Sprintf("out of bounds set(%d): bounds i=%d", i, len(a)))
	}
	a[i] = v
}

type twoArray struct {
	i, j    int
	jStride int
	data    []float64
}

func newTwoArray(i, j int, data []float64) twoArray {
	if len(data) < i*j {
		panic(fmt.Sprintf("short data: len(data)=%d, i=%d, j=%d", len(data), i, j))
	}
	return twoArray{
		i:       i,
		j:       j,
		jStride: i,
		data:    data[:i*j],
	}
}

func (a twoArray) at(i, j int) float64 {
	i--
	j--
	if i < 0 || a.i <= i || j < 0 || a.j <= j {
		panic(fmt.Sprintf("out of bounds at(%d, %d): bounds i=%d, j=%d", i, j, a.i, a.j))
	}
	return a.data[i+a.jStride*j]
}

func (a twoArray) set(i, j int, v float64) {
	i--
	j--
	if i < 0 || a.i <= i || j < 0 || a.j <= j {
		panic(fmt.Sprintf("out of bounds set(%d, %d): bounds i=%d, j=%d", i, j, a.i, a.j))
	}
	a.data[i+a.jStride*j] = v
}

type threeArray struct {
	i, j, k          int
	jStride, kStride int
	data             []float64
}

func newThreeArray(i, j, k int, data []float64) threeArray {
	if len(data) < i*j*k {
		panic(fmt.Sprintf("short data: len(data)=%d, i=%d, j=%d, k=%d", len(data), i, j, k))
	}
	return threeArray{
		i:       i,
		j:       j,
		k:       k,
		jStride: i,
		kStride: i * j,
		data:    data[:i*j*k],
	}
}

func (a threeArray) at(i, j, k int) float64 {
	i--
	j--
	k--
	if i < 0 || a.i <= i || j < 0 || a.j <= j || k < 0 || a.k <= k {
		panic(fmt.Sprintf("out of bounds at(%d, %d, %d): bounds i=%d, j=%d, k=%d", i, j, k, a.i, a.j, a.k))
	}
	return a.data[i+a.jStride*j+a.kStride*k]
}

func (a threeArray) set(i, j, k int, v float64) {
	i--
	j--
	k--
	if i < 0 || a.i <= i || j < 0 || a.j <= j || k < 0 || a.k <= k {
		panic(fmt.Sprintf("out of bounds set(%d, %d, %d): bounds i=%d, j=%d, k=%d", i, j, k, a.i, a.j, a.k))
	}
	a.data[i+a.jStride*j+a.kStride*k] = v
}
