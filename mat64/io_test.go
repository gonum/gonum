// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"bytes"
	"encoding"
	"math"
	"testing"

	"github.com/gonum/blas/blas64"
)

var (
	_ encoding.BinaryMarshaler   = (*Dense)(nil)
	_ encoding.BinaryUnmarshaler = (*Dense)(nil)
	_ encoding.BinaryMarshaler   = (*Vector)(nil)
	_ encoding.BinaryUnmarshaler = (*Vector)(nil)
)

var denseData = []struct {
	raw  []byte
	want *Dense
	eq   func(got, want Matrix) bool
}{
	{
		raw:  []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"),
		want: NewDense(0, 0, []float64{}),
		eq:   Equal,
	},
	{
		raw:  []byte("\x02\x00\x00\x00\x00\x00\x00\x00\x02\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf0?\x00\x00\x00\x00\x00\x00\x00@\x00\x00\x00\x00\x00\x00\b@\x00\x00\x00\x00\x00\x00\x10@"),
		want: NewDense(2, 2, []float64{1, 2, 3, 4}),
		eq:   Equal,
	},
	{
		raw:  []byte("\x02\x00\x00\x00\x00\x00\x00\x00\x03\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf0?\x00\x00\x00\x00\x00\x00\x00@\x00\x00\x00\x00\x00\x00\b@\x00\x00\x00\x00\x00\x00\x10@\x00\x00\x00\x00\x00\x00\x14@\x00\x00\x00\x00\x00\x00\x18@"),
		want: NewDense(2, 3, []float64{1, 2, 3, 4, 5, 6}),
		eq:   Equal,
	},
	{
		raw:  []byte("\x03\x00\x00\x00\x00\x00\x00\x00\x02\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf0?\x00\x00\x00\x00\x00\x00\x00@\x00\x00\x00\x00\x00\x00\b@\x00\x00\x00\x00\x00\x00\x10@\x00\x00\x00\x00\x00\x00\x14@\x00\x00\x00\x00\x00\x00\x18@"),
		want: NewDense(3, 2, []float64{1, 2, 3, 4, 5, 6}),
		eq:   Equal,
	},
	{
		raw:  []byte("\x03\x00\x00\x00\x00\x00\x00\x00\x03\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf0?\x00\x00\x00\x00\x00\x00\x00@\x00\x00\x00\x00\x00\x00\b@\x00\x00\x00\x00\x00\x00\x10@\x00\x00\x00\x00\x00\x00\x14@\x00\x00\x00\x00\x00\x00\x18@\x00\x00\x00\x00\x00\x00\x1c@\x00\x00\x00\x00\x00\x00 @\x00\x00\x00\x00\x00\x00\"@"),
		want: NewDense(3, 3, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}),
		eq:   Equal,
	},
	{
		raw:  []byte("\x02\x00\x00\x00\x00\x00\x00\x00\x02\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf0?\x00\x00\x00\x00\x00\x00\x00@\x00\x00\x00\x00\x00\x00\x10@\x00\x00\x00\x00\x00\x00\x14@"),
		want: NewDense(3, 3, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}).View(0, 0, 2, 2).(*Dense),
		eq:   Equal,
	},
	{
		raw:  []byte("\x02\x00\x00\x00\x00\x00\x00\x00\x02\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x14@\x00\x00\x00\x00\x00\x00\x18@\x00\x00\x00\x00\x00\x00 @\x00\x00\x00\x00\x00\x00\"@"),
		want: NewDense(3, 3, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}).View(1, 1, 2, 2).(*Dense),
		eq:   Equal,
	},
	{
		raw:  []byte("\x03\x00\x00\x00\x00\x00\x00\x00\x02\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00@\x00\x00\x00\x00\x00\x00\b@\x00\x00\x00\x00\x00\x00\x14@\x00\x00\x00\x00\x00\x00\x18@\x00\x00\x00\x00\x00\x00 @\x00\x00\x00\x00\x00\x00\"@"),
		want: NewDense(3, 3, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}).View(0, 1, 3, 2).(*Dense),
		eq:   Equal,
	},
	{
		raw:  []byte("\x01\x00\x00\x00\x00\x00\x00\x00\x04\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf0\xff\x00\x00\x00\x00\x00\x00\xf0\u007f\x01\x00\x00\x00\x00\x00\xf8\u007f"),
		want: NewDense(1, 4, []float64{0, math.Inf(-1), math.Inf(+1), math.NaN()}),
		eq: func(got, want Matrix) bool {
			for _, v := range []bool{
				got.At(0, 0) == 0,
				math.IsInf(got.At(0, 1), -1),
				math.IsInf(got.At(0, 2), +1),
				math.IsNaN(got.At(0, 3)),
			} {
				if !v {
					return false
				}
			}
			return true
		},
	},
}

func TestDenseMarshal(t *testing.T) {
	for i, test := range denseData {
		buf, err := test.want.MarshalBinary()
		if err != nil {
			t.Errorf("error encoding test-%d: %v\n", i, err)
			continue
		}

		nrows, ncols := test.want.Dims()
		sz := nrows*ncols*sizeFloat64 + 2*sizeInt64
		if len(buf) != sz {
			t.Errorf("encoded size test-%d: want=%d got=%d\n", i, sz, len(buf))
		}

		if !bytes.Equal(buf, test.raw) {
			t.Errorf("error encoding test-%d: bytes mismatch.\n got=%q\nwant=%q\n",
				i,
				string(buf),
				string(test.raw),
			)
			continue
		}
	}
}

func TestDenseUnmarshal(t *testing.T) {
	for i, test := range denseData {
		var v Dense
		err := v.UnmarshalBinary(test.raw)
		if err != nil {
			t.Errorf("error decoding test-%d: %v\n", i, err)
			continue
		}
		if !test.eq(&v, test.want) {
			t.Errorf("error decoding test-%d: values differ.\n got=%v\nwant=%v\n",
				i,
				&v,
				test.want,
			)
		}
	}
}

func TestDenseIORoundTrip(t *testing.T) {
	for i, test := range denseData {
		buf, err := test.want.MarshalBinary()
		if err != nil {
			t.Errorf("error encoding test #%d: %v\n", i, err)
		}

		var got Dense
		err = got.UnmarshalBinary(buf)
		if err != nil {
			t.Errorf("error decoding test #%d: %v\n", i, err)
		}

		if !test.eq(&got, test.want) {
			t.Errorf("r/w test #%d failed\nwant=%#v\n got=%#v\n", i, test.want, &got)
		}
	}
}

var vectorData = []struct {
	raw  []byte
	want *Vector
	eq   func(got, want Matrix) bool
}{
	{
		raw:  []byte("\x00\x00\x00\x00\x00\x00\x00\x00"),
		want: NewVector(0, []float64{}),
		eq:   Equal,
	},
	{
		raw:  []byte("\x04\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf0?\x00\x00\x00\x00\x00\x00\x00@\x00\x00\x00\x00\x00\x00\b@\x00\x00\x00\x00\x00\x00\x10@"),
		want: NewVector(4, []float64{1, 2, 3, 4}),
		eq:   Equal,
	},
	{
		raw:  []byte("\x06\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf0?\x00\x00\x00\x00\x00\x00\x00@\x00\x00\x00\x00\x00\x00\b@\x00\x00\x00\x00\x00\x00\x10@\x00\x00\x00\x00\x00\x00\x14@\x00\x00\x00\x00\x00\x00\x18@"),
		want: NewVector(6, []float64{1, 2, 3, 4, 5, 6}),
		eq:   Equal,
	},
	{
		raw:  []byte("\t\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf0?\x00\x00\x00\x00\x00\x00\x00@\x00\x00\x00\x00\x00\x00\b@\x00\x00\x00\x00\x00\x00\x10@\x00\x00\x00\x00\x00\x00\x14@\x00\x00\x00\x00\x00\x00\x18@\x00\x00\x00\x00\x00\x00\x1c@\x00\x00\x00\x00\x00\x00 @\x00\x00\x00\x00\x00\x00\"@"),
		want: NewVector(9, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}),
		eq:   Equal,
	},
	{
		raw:  []byte("\x03\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf0?\x00\x00\x00\x00\x00\x00\x00@\x00\x00\x00\x00\x00\x00\b@"),
		want: NewVector(9, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}).ViewVec(0, 3),
		eq:   Equal,
	},
	{
		raw:  []byte("\x03\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00@\x00\x00\x00\x00\x00\x00\b@\x00\x00\x00\x00\x00\x00\x10@"),
		want: NewVector(9, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}).ViewVec(1, 3),
		eq:   Equal,
	},
	{
		raw:  []byte("\b\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf0?\x00\x00\x00\x00\x00\x00\x00@\x00\x00\x00\x00\x00\x00\b@\x00\x00\x00\x00\x00\x00\x10@\x00\x00\x00\x00\x00\x00\x14@\x00\x00\x00\x00\x00\x00\x18@\x00\x00\x00\x00\x00\x00\x1c@\x00\x00\x00\x00\x00\x00 @"),
		want: NewVector(9, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}).ViewVec(0, 8),
		eq:   Equal,
	},
	{
		raw: []byte("\x03\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\b@\x00\x00\x00\x00\x00\x00\x18@"),
		want: &Vector{
			mat: blas64.Vector{
				Data: []float64{0, 1, 2, 3, 4, 5, 6},
				Inc:  3,
			},
			n: 3,
		},
		eq: Equal,
	},
	{
		raw:  []byte("\x04\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xf0\xff\x00\x00\x00\x00\x00\x00\xf0\u007f\x01\x00\x00\x00\x00\x00\xf8\u007f"),
		want: NewVector(4, []float64{0, math.Inf(-1), math.Inf(+1), math.NaN()}),
		eq: func(got, want Matrix) bool {
			for _, v := range []bool{
				got.At(0, 0) == 0,
				math.IsInf(got.At(1, 0), -1),
				math.IsInf(got.At(2, 0), +1),
				math.IsNaN(got.At(3, 0)),
			} {
				if !v {
					return false
				}
			}
			return true
		},
	},
}

func TestVectorMarshal(t *testing.T) {
	for i, test := range vectorData {
		buf, err := test.want.MarshalBinary()
		if err != nil {
			t.Errorf("error encoding test-%d: %v\n", i, err)
			continue
		}

		nrows, ncols := test.want.Dims()
		sz := nrows*ncols*sizeFloat64 + sizeInt64
		if len(buf) != sz {
			t.Errorf("encoded size test-%d: want=%d got=%d\n", i, sz, len(buf))
		}

		if !bytes.Equal(buf, test.raw) {
			t.Errorf("error encoding test-%d: bytes mismatch.\n got=%q\nwant=%q\n",
				i,
				string(buf),
				string(test.raw),
			)
			continue
		}
	}
}

func TestVectorUnmarshal(t *testing.T) {
	for i, test := range vectorData {
		var v Vector
		err := v.UnmarshalBinary(test.raw)
		if err != nil {
			t.Errorf("error decoding test-%d: %v\n", i, err)
			continue
		}
		if !test.eq(&v, test.want) {
			t.Errorf("error decoding test-%d: values differ.\n got=%v\nwant=%v\n",
				i,
				&v,
				test.want,
			)
		}
	}
}

func TestVectorIORoundTrip(t *testing.T) {
	for i, test := range vectorData {
		buf, err := test.want.MarshalBinary()
		if err != nil {
			t.Errorf("error encoding test #%d: %v\n", i, err)
		}

		var got Vector
		err = got.UnmarshalBinary(buf)
		if err != nil {
			t.Errorf("error decoding test #%d: %v\n", i, err)
		}
		if !test.eq(&got, test.want) {
			t.Errorf("r/w test #%d failed\n got=%#v\nwant=%#v\n", i, &got, test.want)
		}
	}
}
