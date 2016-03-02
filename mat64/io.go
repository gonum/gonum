// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	// maxLen is the biggest slice/array len one can create on a 32/64b platform.
	maxLen = int64(int(^uint(0) >> 1))
)

var (
	sizeInt64   = binary.Size(int64(0))
	sizeFloat64 = binary.Size(float64(0))

	errTooBig    = errors.New("mat64: resulting data slice too big")
	errBadBuffer = errors.New("mat64: data buffer size mismatch")
	errBadSize   = errors.New("mat64: invalid dimension")
)

// MarshalBinary encodes the receiver into a binary form and returns the result.
//
// Dense is little-endian encoded as follows:
//   0 -  7  number of rows    (int64)
//   8 - 15  number of columns (int64)
//  16 - ..  matrix data elements (float64)
//           [0,0] [0,1] ... [0,ncols-1]
//           [1,0] [1,1] ... [1,ncols-1]
//           ...
//           [nrows-1,0] ... [nrows-1,ncols-1]
func (m Dense) MarshalBinary() ([]byte, error) {
	bufLen := int64(m.mat.Rows)*int64(m.mat.Cols)*int64(sizeFloat64) + 2*int64(sizeInt64)
	if bufLen <= 0 {
		// bufLen is too big and has wrapped around.
		return nil, errTooBig
	}

	buf := bytes.NewBuffer(make([]byte, 0, bufLen))
	err := binary.Write(buf, binary.LittleEndian, int64(m.mat.Rows))
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, int64(m.mat.Cols))
	if err != nil {
		return nil, err
	}

	r, c := m.Dims()
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			err = binary.Write(buf, binary.LittleEndian, m.at(i, j))
			if err != nil {
				return nil, err
			}
		}
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary decodes the binary form into the receiver.
// It panics if the receiver is a non-zero Dense matrix.
//
// See MarshalBinary for the on-disk layout.
//
// Limited checks on the validity of the binary input are performed:
//  - matrix.ErrShape is returned if the number of rows or columns is negative,
//  - an error is returned if the resulting Dense matrix is too
//  big for the current architecture (e.g. a 16GB matrix written by a
//  64b application and read back from a 32b application.)
// UnmarshalBinary does not limit the size of the unmarshaled matrix, and so
// it should not be used on untrusted data.
func (m *Dense) UnmarshalBinary(data []byte) error {
	if !m.isZero() {
		panic("mat64: unmarshal into non-zero matrix")
	}

	buf := bytes.NewReader(data)
	var rows int64
	err := binary.Read(buf, binary.LittleEndian, &rows)
	if err != nil {
		return err
	}
	var cols int64
	err = binary.Read(buf, binary.LittleEndian, &cols)
	if err != nil {
		return err
	}

	if rows < 0 || cols < 0 {
		return errBadSize
	}

	size := rows * cols
	if int(size) < 0 || size > maxLen {
		return errTooBig
	}

	if len(data) != int(size)*sizeFloat64+2*sizeInt64 {
		return errBadBuffer
	}

	m.mat.Rows = int(rows)
	m.mat.Cols = int(cols)
	m.mat.Stride = int(cols)
	m.capRows = int(rows)
	m.capCols = int(cols)
	m.mat.Data = use(m.mat.Data, int(size))

	return binary.Read(buf, binary.LittleEndian, m.mat.Data)
}

// MarshalBinary encodes the receiver into a binary form and returns the result.
//
// Vector is little-endian encoded as follows:
//   0 -  7  number of elements     (int64)
//   8 - ..  vector's data elements (float64)
func (v Vector) MarshalBinary() ([]byte, error) {
	bufLen := int64(sizeInt64) + int64(v.n)*int64(sizeFloat64)
	if bufLen <= 0 {
		// bufLen is too big and has wrapped around.
		return nil, errTooBig
	}

	buf := bytes.NewBuffer(make([]byte, 0, bufLen))
	err := binary.Write(buf, binary.LittleEndian, int64(v.n))
	if err != nil {
		return nil, err
	}

	for i := 0; i < v.n; i++ {
		err = binary.Write(buf, binary.LittleEndian, v.at(i))
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary decodes the binary form into the receiver.
// It panics if the receiver is a non-zero Vector.
//
// See MarshalBinary for the on-disk layout.
//
// Limited checks on the validity of the binary input are performed:
//  - matrix.ErrShape is returned if the number of rows is negative,
//  - an error is returned if the resulting Vector is too
//  big for the current architecture (e.g. a 16GB vector written by a
//  64b application and read back from a 32b application.)
// UnmarshalBinary does not limit the size of the unmarshaled vector, and so
// it should not be used on untrusted data.
func (v *Vector) UnmarshalBinary(data []byte) error {
	if !v.isZero() {
		panic("mat64: unmarshal into non-zero vector")
	}

	buf := bytes.NewReader(data)
	var n int64
	err := binary.Read(buf, binary.LittleEndian, &n)
	if err != nil {
		return err
	}

	if n < 0 {
		return errBadSize
	}
	if n > maxLen {
		return errTooBig
	}
	if len(data) != int(n)*sizeFloat64+sizeInt64 {
		return errBadBuffer
	}

	v.n = int(n)
	v.mat.Inc = 1
	v.mat.Data = use(v.mat.Data, v.n)

	return binary.Read(buf, binary.LittleEndian, v.mat.Data)
}
