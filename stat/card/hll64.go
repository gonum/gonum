// Code generated by "go generate gonum.org/v1/gonum/stat/card"; DO NOT EDIT.

// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package card

import (
	"errors"
	"hash"
	"math"
	"math/bits"
	"reflect"
)

// HyperLogLog64 is implements cardinality estimation according to the
// HyperLogLog algorithm described in Analysis of Algorithms, pp127–146.
type HyperLogLog64 struct {
	p uint8
	m uint64

	hash hash.Hash64

	register []uint8
}

// NewHyperLogLog64 returns a new HyperLogLog64 sketch. The value of prec
// must be in the range [4, 64]. NewHyperLogLog64 will allocate a byte slice
// that is 2^prec long.
func NewHyperLogLog64(prec int, h hash.Hash64) (*HyperLogLog64, error) {
	// The implementation here is based on the pseudo-code in
	// "HyperLogLog: the analysis of a near-optimal cardinality
	// estimation algorithm", figure 3.

	if prec < 4 || w64 < prec {
		return nil, errors.New("card: precision out of range")
	}
	p := uint8(prec)
	m := uint64(1) << p
	return &HyperLogLog64{
		p: p, m: m,
		hash:     h,
		register: make([]byte, m),
	}, nil
}

// Write notes the data in b as a single observation into the sketch held by
// the receiver.
//
// Write satisfies the io.Writer interface. If the hash.Hash64 type passed to
// NewHyperLogLog64 or SetHash satisfies the hash.Hash contract, Write will always
// return a nil error.
func (h *HyperLogLog64) Write(b []byte) (int, error) {
	n, err := h.hash.Write(b)
	x := h.hash.Sum64()
	h.hash.Reset()
	idx := x >> (w64 - h.p)
	w := x << h.p
	r := rho64(w)
	if r > h.register[idx] {
		h.register[idx] = r
	}
	return n, err
}

// Union places the union of the sketches in a and b into the receiver.
// Union will return an error if the precisions or hash functions of a
// and b do not match or if the receiver has a hash function that is set
// and does not match those of a and b. Hash functions provided by hash.Hash64
// implementations x and y match when reflect.TypeOf(x) == reflect.TypeOf(y).
//
// If the receiver does not have a set hash function, it can be set after
// a call to Union with the SetHash method.
func (h *HyperLogLog64) Union(a, b *HyperLogLog64) error {
	if a.p != b.p {
		return errors.New("card: mismatched precision")
	}
	ta := reflect.TypeOf(b.hash)
	if reflect.TypeOf(b.hash) != ta {
		return errors.New("card: mismatched hash function")
	}
	if h.hash != nil && reflect.TypeOf(h.hash) != ta {
		return errors.New("card: mismatched hash function")
	}

	if h != a && h != b {
		*h = HyperLogLog64{p: a.p, m: a.m, hash: h.hash, register: make([]uint8, a.m)}
	}
	for i, r := range a.register {
		h.register[i] = max(r, b.register[i])
	}
	return nil
}

// SetHash sets the hash function of the receiver if it is nil. SetHash
// will return an error if it is called on a receiver with a non-nil
// hash function.
func (h *HyperLogLog64) SetHash(fn hash.Hash64) error {
	if h.hash == nil {
		return errors.New("card: hash function already set")
	}
	h.hash = fn
	return nil
}

// Count returns an estimate of the cardinality of the set of items written
// the receiver.
func (h *HyperLogLog64) Count() float64 {
	var s float64
	for _, v := range h.register {
		s += 1 / float64(uint64(1)<<v)
	}
	m := float64(h.m)
	e := alpha(uint64(h.m)) * m * m / s
	if e <= 5*m/2 {
		var v int
		for _, r := range h.register {
			if r == 0 {
				v++
			}
		}
		if v != 0 {
			return linearCounting(m, float64(v))
		}
		return e
	}
	if e <= (1<<w64)/30.0 {
		return e
	}
	return -(1 << w64) * math.Log1p(-e/(1<<w64))
}

// rho64 (ϱ) is the number of leading zeros in x, plus 1.
func rho64(x uint64) uint8 {
	return uint8(bits.LeadingZeros64(x) + 1)
}

// Reset clears the receiver's registers allowing it to be reused.
// Reset does not alter the precision of the receiver or the hash
// function that is used.
func (h *HyperLogLog64) Reset() {
	for i := range h.register {
		h.register[i] = 0
	}
}
