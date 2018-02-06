// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fourier

import "fmt"

type twoArrayZ struct {
	i, j    int
	jStride int
	data    []float64
}

func newTwoArrayZ(i, j int, data []float64) twoArrayZ {
	if len(data) < i*j {
		panic(fmt.Sprintf("short data: len(data)=%d, i=%d, j=%d", len(data), i, j))
	}
	return twoArrayZ{
		i:       i,
		j:       j,
		jStride: i,
		data:    data[:i*j],
	}
}

func (a twoArrayZ) at(i, j int) float64 {
	if i < 0 || a.i <= i || j < 0 || a.j <= j {
		panic(fmt.Sprintf("out of bounds at(%d, %d): bounds i=%d, j=%d", i, j, a.i, a.j))
	}
	return a.data[i+a.jStride*j]
}

func (a twoArrayZ) set(i, j int, v float64) {
	if i < 0 || a.i <= i || j < 0 || a.j <= j {
		panic(fmt.Sprintf("out of bounds set(%d, %d): bounds i=%d, j=%d", i, j, a.i, a.j))
	}
	a.data[i+a.jStride*j] = v
}

type threeArrayZ struct {
	i, j, k          int
	jStride, kStride int
	data             []float64
}

func newThreeArrayZ(i, j, k int, data []float64) threeArrayZ {
	if len(data) < i*j*k {
		panic(fmt.Sprintf("short data: len(data)=%d, i=%d, j=%d, k=%d", len(data), i, j, k))
	}
	return threeArrayZ{
		i:       i,
		j:       j,
		k:       k,
		jStride: i,
		kStride: i * j,
		data:    data[:i*j*k],
	}
}

func (a threeArrayZ) at(i, j, k int) float64 {
	if i < 0 || a.i <= i || j < 0 || a.j <= j || k < 0 || a.k <= k {
		panic(fmt.Sprintf("out of bounds at(%d, %d, %d): bounds i=%d, j=%d, k=%d", i, j, k, a.i, a.j, a.k))
	}
	return a.data[i+a.jStride*j+a.kStride*k]
}

func (a threeArrayZ) set(i, j, k int, v float64) {
	if i < 0 || a.i <= i || j < 0 || a.j <= j || k < 0 || a.k <= k {
		panic(fmt.Sprintf("out of bounds set(%d, %d, %d): bounds i=%d, j=%d, k=%d", i, j, k, a.i, a.j, a.k))
	}
	a.data[i+a.jStride*j+a.kStride*k] = v
}
