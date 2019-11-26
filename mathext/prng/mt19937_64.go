// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Original C program copyright Takuji Nishimura and Makoto Matsumoto 2004.
// http://www.math.sci.hiroshima-u.ac.jp/~m-mat/MT/VERSIONS/C-LANG/mt19937-64.c

package prng

import (
	"encoding/binary"
	"io"
)

const (
	mt19937_64NN        = 312
	mt19937_64MM        = 156
	mt19937_64MatrixA   = 0xB5026F5AA96619E9
	mt19937_64UpperMask = 0xFFFFFFFF80000000
	mt19937_64LowerMask = 0x7FFFFFFF
)

// MT19937_64 implements the 64 bit Mersenne Twister PRNG. MT19937_64
// is the 64 bit version of MT19937, it has the same sized state, but
// generates a different sequence.
// See https://en.wikipedia.org/wiki/Mersenne_Twister.
type MT19937_64 struct {
	mt  [mt19937_64NN]uint64
	mti uint64
}

// NewMT19937_64 returns a new MT19937_64 PRNG. The returned PRNG will
// use the default seed 5489 unless the Seed method is called with
// another value.
func NewMT19937_64() *MT19937_64 {
	return &MT19937_64{mti: mt19937_64NN + 1}
}

// Seed uses the provided seed value to initialize the generator to a
// deterministic state.
func (src *MT19937_64) Seed(seed uint64) {
	src.mt[0] = seed
	for src.mti = 1; src.mti < mt19937_64NN; src.mti++ {
		src.mt[src.mti] = (6364136223846793005*(src.mt[src.mti-1]^(src.mt[src.mti-1]>>62)) + src.mti)
	}
}

// SeedFromKeys uses the provided seed key value to initialize the
// generator to a deterministic state. It is provided for compatibility
// with C implementations.
func (src *MT19937_64) SeedFromKeys(keys []uint64) {
	src.Seed(19650218)
	i := uint64(1)
	j := uint64(0)
	k := uint64(mt19937_64NN)
	if k <= uint64(len(keys)) {
		k = uint64(len(keys))
	}
	for ; k != 0; k-- {
		src.mt[i] = (src.mt[i] ^ ((src.mt[i-1] ^ (src.mt[i-1] >> 62)) * 3935559000370003845)) + keys[j] + j // Non linear.
		i++
		j++
		if i >= mt19937_64NN {
			src.mt[0] = src.mt[mt19937_64NN-1]
			i = 1
		}
		if j >= uint64(len(keys)) {
			j = 0
		}
	}
	for k = mt19937_64NN - 1; k != 0; k-- {
		src.mt[i] = (src.mt[i] ^ ((src.mt[i-1] ^ (src.mt[i-1] >> 62)) * 2862933555777941757)) - i // Non linear.
		i++
		if i >= mt19937_64NN {
			src.mt[0] = src.mt[mt19937_64NN-1]
			i = 1
		}
	}

	src.mt[0] = 1 << 63 /* MSB is 1; assuring non-zero initial array */
}

// Uint64 returns a pseudo-random 64-bit unsigned integer as a uint64.
func (src *MT19937_64) Uint64() uint64 {
	mag01 := [2]uint64{0, mt19937_64MatrixA}

	var x uint64
	if src.mti >= mt19937_64NN { // Generate mt19937_64NN words at one time.
		if src.mti == mt19937_64NN+1 {
			// If Seed() has not been called
			// a default initial seed is used.
			src.Seed(5489)
		}

		var i int
		for ; i < mt19937_64NN-mt19937_64MM; i++ {
			x = (src.mt[i] & mt19937_64UpperMask) | (src.mt[i+1] & mt19937_64LowerMask)
			src.mt[i] = src.mt[i+mt19937_64MM] ^ (x >> 1) ^ mag01[(int)(x&0x1)]
		}
		for ; i < mt19937_64NN-1; i++ {
			x = (src.mt[i] & mt19937_64UpperMask) | (src.mt[i+1] & mt19937_64LowerMask)
			src.mt[i] = src.mt[i+(mt19937_64MM-mt19937_64NN)] ^ (x >> 1) ^ mag01[(int)(x&0x1)]
		}
		x = (src.mt[mt19937_64NN-1] & mt19937_64UpperMask) | (src.mt[0] & mt19937_64LowerMask)
		src.mt[mt19937_64NN-1] = src.mt[mt19937_64MM-1] ^ (x >> 1) ^ mag01[(int)(x&0x1)]

		src.mti = 0
	}

	x = src.mt[src.mti]
	src.mti++

	// Tempering.
	x ^= (x >> 29) & 0x5555555555555555
	x ^= (x << 17) & 0x71D67FFFEDA60000
	x ^= (x << 37) & 0xFFF7EEE000000000
	x ^= (x >> 43)

	return x
}

// MarshalBinary returns the binary representation of the current state of the generator.
func (src *MT19937_64) MarshalBinary() ([]byte, error) {
	var buf [(mt19937_64NN + 1) * 8]byte
	for i := 0; i < mt19937_64NN; i++ {
		binary.BigEndian.PutUint64(buf[i*8:(i+1)*8], src.mt[i])
	}
	binary.BigEndian.PutUint64(buf[mt19937_64NN*8:], src.mti)
	return buf[:], nil
}

// UnmarshalBinary sets the state of the generator to the state represented in data.
func (src *MT19937_64) UnmarshalBinary(data []byte) error {
	if len(data) < (mt19937_64NN+1)*8 {
		return io.ErrUnexpectedEOF
	}
	for i := 0; i < mt19937_64NN; i++ {
		src.mt[i] = binary.BigEndian.Uint64(data[i*8 : (i+1)*8])
	}
	src.mti = binary.BigEndian.Uint64(data[mt19937_64NN*8:])
	return nil
}
