// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mmarket provides a type to read Matrix Market format files.
package mmarket

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"gonum.org/v1/exp/linsolve/internal/triplet"
)

var (
	errBadFormat   = errors.New("mmarket: bad file format")
	errUnsupported = errors.New("mmarket: matrix type not supported")
)

type Reader struct {
	s *bufio.Scanner
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		s: bufio.NewScanner(r),
	}
}

// Read a real matrix in coordinate format and return its triplet representation.
func (r *Reader) Read() (*triplet.Matrix, error) {
	r.s.Scan()
	if err := r.s.Err(); err != nil {
		return nil, err
	}
	header := strings.Fields(r.s.Text())
	if header[0] != "%%MatrixMarket" {
		return nil, errBadFormat
	}
	if header[1] != "matrix" {
		return nil, errBadFormat
	}
	if header[2] != "coordinate" {
		return nil, errBadFormat
	}
	if header[3] != "real" {
		return nil, errUnsupported
	}
	sym := header[4] == "symmetric"

	var nr, nc, nnz int
	for r.s.Scan() {
		line := r.s.Text()
		if line[0] == '%' {
			continue
		}
		n, err := fmt.Sscan(line, &nr, &nc, &nnz)
		if err != nil {
			return nil, err
		}
		if n != 3 {
			return nil, errBadFormat
		}
		break
	}
	if err := r.s.Err(); err != nil {
		return nil, err
	}

	if sym && nr != nc {
		return nil, errBadFormat
	}

	m := triplet.NewMatrix(nr, nc)
	for i := 0; i < nnz; i++ {
		if !r.s.Scan() {
			return nil, errBadFormat
		}
		var (
			i, j int
			v    float64
		)
		n, err := fmt.Sscan(r.s.Text(), &i, &j, &v)
		if err != nil {
			return nil, err
		}
		if n != 3 {
			return nil, errBadFormat
		}
		if i < 1 || nr < i {
			return nil, errBadFormat
		}
		if j < 1 || nc < j {
			return nil, errBadFormat
		}
		m.Append(i-1, j-1, v)
		if sym && i != j {
			m.Append(j-1, i-1, v)
		}
	}
	return m, nil
}
