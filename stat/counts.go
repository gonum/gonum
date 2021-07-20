// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import "gonum.org/v1/gonum/mat"

// CountsToDist returns a probability distribution from a discrete vector of
// counts.
//   p(x) = c_x / N
// where N = \sum c_x. The function returns mat.Vector
func CountsToDist(c []int) mat.Vector {
	var n float64
	r := make([]float64, len(c), len(c))

	for _, v := range c {
		n += float64(v)
	}

	for i, v := range c {
		r[i] = float64(v) / n
	}

	p := mat.NewVecDense(len(r), r)

	return p
}

// Counts returns an integer vector with number of occurrences of each number in
// the data vector from 0 to max, where max is the highest number in the data.
func Counts(d []int) []int {
	var max int

	for _, v := range d {
		if v > max {
			max = v
		}
	}

	max++

	c := make([]int, max, max)

	for _, v := range d {
		c[v]++
	}

	return c
}
