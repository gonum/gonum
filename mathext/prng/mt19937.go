// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Original C program copyright Takuji Nishimura and Makoto Matsumoto 2002.
// http://www.math.sci.hiroshima-u.ac.jp/~m-mat/MT/MT2002/CODES/mt19937ar.c

package prng

import (
	"encoding/binary"
	"io"
)

const (
	mt19937N         = 624
	mt19937M         = 397
	mt19937matrixA   = 0x9908b0df
	mt19937UpperMask = 0x80000000
	mt19937LowerMask = 0x7fffffff
)

// MT19937 implements the 32 bit Mersenne Twister PRNG. MT19937
// is the default PRNG for a wide variety of programming systems.
// See https://en.wikipedia.org/wiki/Mersenne_Twister.
type MT19937 struct {
	mt  [mt19937N]uint32
	mti uint32
}

// NewMT19937 returns a new MT19937 PRNG. The returned PRNG will
// use the default seed 5489 unless the Seed method is called with
// another value.
func NewMT19937() *MT19937 {
	return &MT19937{mti: mt19937N + 1}
}

// Seed uses the provided seed value to initialize the generator to a
// deterministic state. Only the lower 32 bits of seed are used to seed
// the PRNG.
func (src *MT19937) Seed(seed uint64) {
	src.mt[0] = uint32(seed)
	for src.mti = 1; src.mti < mt19937N; src.mti++ {
		src.mt[src.mti] = (1812433253*(src.mt[src.mti-1]^(src.mt[src.mti-1]>>30)) + src.mti)
	}
}

// SeedFromKeys uses the provided seed key value to initialize the
// generator to a deterministic state. It is provided for compatibility
// with C implementations.
func (src *MT19937) SeedFromKeys(keys []uint32) {
	src.Seed(19650218)
	i := uint32(1)
	j := uint32(0)
	k := uint32(mt19937N)
	if k <= uint32(len(keys)) {
		k = uint32(len(keys))
	}
	for ; k != 0; k-- {
		src.mt[i] = (src.mt[i] ^ ((src.mt[i-1] ^ (src.mt[i-1] >> 30)) * 1664525)) + keys[j] + j // Non linear.
		i++
		j++
		if i >= mt19937N {
			src.mt[0] = src.mt[mt19937N-1]
			i = 1
		}
		if j >= uint32(len(keys)) {
			j = 0
		}
	}
	for k = mt19937N - 1; k != 0; k-- {
		src.mt[i] = (src.mt[i] ^ ((src.mt[i-1] ^ (src.mt[i-1] >> 30)) * 1566083941)) - i // Non linear.
		i++
		if i >= mt19937N {
			src.mt[0] = src.mt[mt19937N-1]
			i = 1
		}
	}
	src.mt[0] = 0x80000000 // MSB is 1; assuring non-zero initial array.
}

// Uint32 returns a pseudo-random 32-bit unsigned integer as a uint32.
func (src *MT19937) Uint32() uint32 {
	mag01 := [2]uint32{0x0, mt19937matrixA}

	var y uint32
	if src.mti >= mt19937N { // Generate mt19937N words at one time.
		if src.mti == mt19937N+1 {
			// If Seed() has not been called
			// a default initial seed is used.
			src.Seed(5489)
		}

		var kk int
		for ; kk < mt19937N-mt19937M; kk++ {
			y = (src.mt[kk] & mt19937UpperMask) | (src.mt[kk+1] & mt19937LowerMask)
			src.mt[kk] = src.mt[kk+mt19937M] ^ (y >> 1) ^ mag01[y&0x1]
		}
		for ; kk < mt19937N-1; kk++ {
			y = (src.mt[kk] & mt19937UpperMask) | (src.mt[kk+1] & mt19937LowerMask)
			src.mt[kk] = src.mt[kk+(mt19937M-mt19937N)] ^ (y >> 1) ^ mag01[y&0x1]
		}
		y = (src.mt[mt19937N-1] & mt19937UpperMask) | (src.mt[0] & mt19937LowerMask)
		src.mt[mt19937N-1] = src.mt[mt19937M-1] ^ (y >> 1) ^ mag01[y&0x1]

		src.mti = 0
	}

	y = src.mt[src.mti]
	src.mti++

	// Tempering.
	y ^= (y >> 11)
	y ^= (y << 7) & 0x9d2c5680
	y ^= (y << 15) & 0xefc60000
	y ^= (y >> 18)

	return y
}

// Uint64 returns a pseudo-random 64-bit unsigned integer as a uint64.
// It makes use of two calls to Uint32 placing the first result in the
// upper bits and the second result in the lower bits of the returned
// value.
func (src *MT19937) Uint64() uint64 {
	h := uint64(src.Uint32())
	l := uint64(src.Uint32())
	return h<<32 | l
}

// MarshalBinary returns the binary representation of the current state of the generator.
func (src *MT19937) MarshalBinary() ([]byte, error) {
	var buf [(mt19937N + 1) * 4]byte
	for i := 0; i < mt19937N; i++ {
		binary.BigEndian.PutUint32(buf[i*4:(i+1)*4], src.mt[i])
	}
	binary.BigEndian.PutUint32(buf[mt19937N*4:], src.mti)
	return buf[:], nil
}

// UnmarshalBinary sets the state of the generator to the state represented in data.
func (src *MT19937) UnmarshalBinary(data []byte) error {
	if len(data) < (mt19937N+1)*4 {
		return io.ErrUnexpectedEOF
	}
	for i := 0; i < mt19937N; i++ {
		src.mt[i] = binary.BigEndian.Uint32(data[i*4 : (i+1)*4])
	}
	src.mti = binary.BigEndian.Uint32(data[mt19937N*4:])
	return nil
}
