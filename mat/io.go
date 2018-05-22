// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

// version is the current on-disk codec version.
const version uint64 = 0x1

// maxLen is the biggest slice/array len one can create on a 32/64b platform.
const maxLen = int64(int(^uint(0) >> 1))

var (
	headerSize  = binary.Size(storage{})
	sizeInt64   = binary.Size(int64(0))
	sizeFloat64 = binary.Size(float64(0))

	errWrongType = errors.New("mat: wrong data type")

	errTooBig    = errors.New("mat: resulting data slice too big")
	errTooSmall  = errors.New("mat: input slice too small")
	errBadBuffer = errors.New("mat: data buffer size mismatch")
	errBadSize   = errors.New("mat: invalid dimension")
)

// MarshalBinary encodes the receiver into a binary form and returns the result.
//
// Dense is little-endian encoded as follows:
//   0 -  7  Version = 1          (uint64)
//   8       'G'                  (byte)
//   9       'F'                  (byte)
//  10       'A'                  (byte)
//  11       0                    (byte)
//  12 - 19  number of rows       (int64)
//  20 - 27  number of columns    (int64)
//  28 - 35  0                    (int64)
//  36 - 43  0                    (int64)
//  44 - ..  matrix data elements (float64)
//           [0,0] [0,1] ... [0,ncols-1]
//           [1,0] [1,1] ... [1,ncols-1]
//           ...
//           [nrows-1,0] ... [nrows-1,ncols-1]
func (m Dense) MarshalBinary() ([]byte, error) {
	bufLen := int64(headerSize) + int64(m.mat.Rows)*int64(m.mat.Cols)*int64(sizeFloat64)
	if bufLen <= 0 {
		// bufLen is too big and has wrapped around.
		return nil, errTooBig
	}

	b := make([]byte, bufLen)
	buf := bytes.NewBuffer(b[:0])
	err := binary.Write(buf, binary.LittleEndian, storage{
		Form: 'G', Packing: 'F', Uplo: 'A',
		Rows: int64(m.mat.Rows), Cols: int64(m.mat.Cols),
		Version: version,
	})
	if err != nil {
		return buf.Bytes(), err
	}
	p := headerSize
	r, c := m.Dims()
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			binary.LittleEndian.PutUint64(b[p:p+sizeFloat64], math.Float64bits(m.at(i, j)))
			p += sizeFloat64
		}
	}

	return b, nil
}

// MarshalBinaryTo encodes the receiver into a binary form and writes it into w.
// MarshalBinaryTo returns the number of bytes written into w and an error, if any.
//
// See MarshalBinary for the on-disk layout.
func (m Dense) MarshalBinaryTo(w io.Writer) (int, error) {
	buf := bytes.NewBuffer(make([]byte, 0, headerSize))
	err := binary.Write(buf, binary.LittleEndian, storage{
		Form: 'G', Packing: 'F', Uplo: 'A',
		Rows: int64(m.mat.Rows), Cols: int64(m.mat.Cols),
		Version: version,
	})
	if err != nil {
		return 0, err
	}
	n, err := w.Write(buf.Bytes())
	if err != nil {
		return n, err
	}

	r, c := m.Dims()
	var b [8]byte
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			binary.LittleEndian.PutUint64(b[:], math.Float64bits(m.at(i, j)))
			nn, err := w.Write(b[:])
			n += nn
			if err != nil {
				return n, err
			}
		}
	}

	return n, nil
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
	if !m.IsZero() {
		panic("mat: unmarshal into non-zero matrix")
	}

	if len(data) < headerSize {
		return errTooSmall
	}

	var s storage
	binary.Read(bytes.NewReader(data[:headerSize]), binary.LittleEndian, &s)
	if s.Version != version {
		return fmt.Errorf("mat: incorrect version: %d", s.Version)
	}
	rows := s.Rows
	cols := s.Cols
	s.Version = 0
	s.Rows = 0
	s.Cols = 0
	if (s != storage{Form: 'G', Packing: 'F', Uplo: 'A'}) {
		return errWrongType
	}
	if rows < 0 || cols < 0 {
		return errBadSize
	}
	size := rows * cols
	if size == 0 {
		return ErrZeroLength
	}
	if int(size) < 0 || size > maxLen {
		return errTooBig
	}
	if len(data) != headerSize+int(rows*cols)*sizeFloat64 {
		return errBadBuffer
	}

	p := headerSize
	m.reuseAs(int(rows), int(cols))
	for i := range m.mat.Data {
		m.mat.Data[i] = math.Float64frombits(binary.LittleEndian.Uint64(data[p : p+sizeFloat64]))
		p += sizeFloat64
	}

	return nil
}

// UnmarshalBinaryFrom decodes the binary form into the receiver and returns
// the number of bytes read and an error if any.
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
func (m *Dense) UnmarshalBinaryFrom(r io.Reader) (int, error) {
	if !m.IsZero() {
		panic("mat: unmarshal into non-zero matrix")
	}

	var s storage
	buf := make([]byte, headerSize)
	n, err := readFull(r, buf)
	if err != nil {
		return n, err
	}
	err = binary.Read(bytes.NewReader(buf), binary.LittleEndian, &s)
	if err != nil {
		return n, err
	}
	if s.Version != version {
		return n, fmt.Errorf("mat: incorrect version: %d", s.Version)
	}
	rows := s.Rows
	cols := s.Cols
	s.Version = 0
	s.Rows = 0
	s.Cols = 0
	if (s != storage{Form: 'G', Packing: 'F', Uplo: 'A'}) {
		return n, errWrongType
	}
	if rows < 0 || cols < 0 {
		return n, errBadSize
	}
	size := rows * cols
	if size == 0 {
		return n, ErrZeroLength
	}
	if int(size) < 0 || size > maxLen {
		return n, errTooBig
	}

	m.reuseAs(int(rows), int(cols))
	var b [8]byte
	for i := range m.mat.Data {
		nn, err := readFull(r, b[:])
		n += nn
		if err != nil {
			if err == io.EOF {
				return n, io.ErrUnexpectedEOF
			}
			return n, err
		}
		m.mat.Data[i] = math.Float64frombits(binary.LittleEndian.Uint64(b[:]))
	}

	return n, nil
}

// MarshalBinary encodes the receiver into a binary form and returns the result.
//
// VecDense is little-endian encoded as follows:
//
//   0 -  7  Version = 1            (uint64)
//   8       'G'                    (byte)
//   9       'F'                    (byte)
//  10       'A'                    (byte)
//  11       0                      (byte)
//  12 - 19  number of elements     (int64)
//  20 - 27  1                      (int64)
//  28 - 35  0                      (int64)
//  36 - 43  0                      (int64)
//  44 - ..  vector's data elements (float64)
func (v VecDense) MarshalBinary() ([]byte, error) {
	bufLen := int64(headerSize) + int64(v.n)*int64(sizeFloat64)
	if bufLen <= 0 {
		// bufLen is too big and has wrapped around.
		return nil, errTooBig
	}

	b := make([]byte, bufLen)
	buf := bytes.NewBuffer(b[:0])
	err := binary.Write(buf, binary.LittleEndian, storage{
		Form: 'G', Packing: 'F', Uplo: 'A',
		Rows: int64(v.n), Cols: 1,
		Version: version,
	})
	if err != nil {
		return buf.Bytes(), err
	}

	p := headerSize
	for i := 0; i < v.n; i++ {
		binary.LittleEndian.PutUint64(b[p:p+sizeFloat64], math.Float64bits(v.at(i)))
		p += sizeFloat64
	}

	return b, nil
}

// MarshalBinaryTo encodes the receiver into a binary form, writes it to w and
// returns the number of bytes written and an error if any.
//
// See MarshalBainry for the on-disk format.
func (v VecDense) MarshalBinaryTo(w io.Writer) (int, error) {
	buf := bytes.NewBuffer(make([]byte, 0, headerSize))
	err := binary.Write(buf, binary.LittleEndian, storage{
		Form: 'G', Packing: 'F', Uplo: 'A',
		Rows: int64(v.n), Cols: 1,
		Version: version,
	})
	if err != nil {
		return 0, err
	}
	n, err := w.Write(buf.Bytes())
	if err != nil {
		return n, err
	}

	var b [8]byte
	for i := 0; i < v.n; i++ {
		binary.LittleEndian.PutUint64(b[:], math.Float64bits(v.at(i)))
		nn, err := w.Write(b[:])
		n += nn
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

// UnmarshalBinary decodes the binary form into the receiver.
// It panics if the receiver is a non-zero VecDense.
//
// See MarshalBinary for the on-disk layout.
//
// Limited checks on the validity of the binary input are performed:
//  - matrix.ErrShape is returned if the number of rows is negative,
//  - an error is returned if the resulting VecDense is too
//  big for the current architecture (e.g. a 16GB vector written by a
//  64b application and read back from a 32b application.)
// UnmarshalBinary does not limit the size of the unmarshaled vector, and so
// it should not be used on untrusted data.
func (v *VecDense) UnmarshalBinary(data []byte) error {
	if !v.IsZero() {
		panic("mat: unmarshal into non-zero vector")
	}

	if len(data) < headerSize {
		return errTooSmall
	}

	var s storage
	binary.Read(bytes.NewReader(data[:headerSize]), binary.LittleEndian, &s)
	if s.Version != version {
		return fmt.Errorf("mat: incorrect version: %d", s.Version)
	}
	if s.Cols != 1 {
		return ErrShape
	}
	n := s.Rows
	s.Version = 0
	s.Rows = 0
	s.Cols = 0
	if (s != storage{Form: 'G', Packing: 'F', Uplo: 'A'}) {
		return errWrongType
	}
	if n == 0 {
		return ErrZeroLength
	}
	if n < 0 {
		return errBadSize
	}
	if int64(maxLen) < n {
		return errTooBig
	}
	if len(data) != headerSize+int(n)*sizeFloat64 {
		return errBadBuffer
	}

	p := headerSize
	v.reuseAs(int(n))
	for i := range v.mat.Data {
		v.mat.Data[i] = math.Float64frombits(binary.LittleEndian.Uint64(data[p : p+sizeFloat64]))
		p += sizeFloat64
	}

	return nil
}

// UnmarshalBinaryFrom decodes the binary form into the receiver, from the
// io.Reader and returns the number of bytes read and an error if any.
// It panics if the receiver is a non-zero VecDense.
//
// See MarshalBinary for the on-disk layout.
// See UnmarshalBinary for the list of sanity checks performed on the input.
func (v *VecDense) UnmarshalBinaryFrom(r io.Reader) (int, error) {
	if !v.IsZero() {
		panic("mat: unmarshal into non-zero vector")
	}

	var s storage
	buf := make([]byte, headerSize)
	n, err := readFull(r, buf)
	if err != nil {
		return n, err
	}
	err = binary.Read(bytes.NewReader(buf), binary.LittleEndian, &s)
	if err != nil {
		return n, err
	}
	if s.Version != version {
		return n, fmt.Errorf("mat: incorrect version: %d", s.Version)
	}
	if s.Cols != 1 {
		return n, ErrShape
	}
	l := s.Rows
	s.Version = 0
	s.Rows = 0
	s.Cols = 0
	if (s != storage{Form: 'G', Packing: 'F', Uplo: 'A'}) {
		return n, errWrongType
	}
	if l == 0 {
		return n, ErrZeroLength
	}
	if l < 0 {
		return n, errBadSize
	}
	if int64(maxLen) < l {
		return n, errTooBig
	}

	v.reuseAs(int(l))
	var b [8]byte
	for i := range v.mat.Data {
		nn, err := readFull(r, b[:])
		n += nn
		if err != nil {
			if err == io.EOF {
				return n, io.ErrUnexpectedEOF
			}
			return n, err
		}
		v.mat.Data[i] = math.Float64frombits(binary.LittleEndian.Uint64(b[:]))
	}

	return n, nil
}

// readFull reads from r into buf until it has read len(buf).
// It returns the number of bytes copied and an error if fewer bytes were read.
// If an EOF happens after reading fewer than len(buf) bytes, io.ErrUnexpectedEOF is returned.
func readFull(r io.Reader, buf []byte) (int, error) {
	var n int
	var err error
	for n < len(buf) && err == nil {
		var nn int
		nn, err = r.Read(buf[n:])
		n += nn
	}
	if n == len(buf) {
		return n, nil
	}
	if err == io.EOF {
		return n, io.ErrUnexpectedEOF
	}
	return n, err
}

// storage is the internal representation of the storage format of a
// serialised matrix.
type storage struct {
	Version uint64 // Keep this first.
	Form    byte   // [GST]
	Packing byte   // [BPF]
	Uplo    byte   // [AUL]
	Unit    bool
	Rows    int64
	Cols    int64
	KU      int64
	KL      int64
}
