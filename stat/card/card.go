// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate ./generate_64bit.sh

package card

import "math"

const (
	w32 = 32
	w64 = 64
)

func alpha(m uint64) float64 {
	if m < 128 {
		return alphaValues[m]
	}
	return 0.7213 / (1 + 1.079/float64(m))
}

var alphaValues = [...]float64{
	16: 0.673,
	32: 0.697,
	64: 0.709,
}

func linearCounting(m, v float64) float64 {
	return m * (math.Log(m) - math.Log(v))
}

func max(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}

func min(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}
