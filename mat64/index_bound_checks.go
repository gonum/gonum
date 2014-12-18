// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file must be kept in sync with index_no_bound_checks.go.

//+build bounds

package mat64

func (m *Dense) At(r, c int) float64 {
	return m.at(r, c)
}

func (m *Dense) at(r, c int) float64 {
	if r >= m.mat.Rows || r < 0 {
		panic("index error: row access out of bounds")
	}
	if c >= m.mat.Cols || c < 0 {
		panic("index error: column access out of bounds")
	}
	return m.mat.Data[r*m.mat.Stride+c]
}

func (m *Dense) Set(r, c int, v float64) {
	m.set(r, c, v)
}

func (m *Dense) set(r, c int, v float64) {
	if r >= m.mat.Rows || r < 0 {
		panic("index error: row access out of bounds")
	}
	if c >= m.mat.Cols || c < 0 {
		panic("index error: column access out of bounds")
	}
	m.mat.Data[r*m.mat.Stride+c] = v
}
