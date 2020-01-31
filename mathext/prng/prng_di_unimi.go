// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// PRNGs from Dipartimento di Informatica Università degli Studi di Milano.
// David Blackman and Sebastiano Vigna licensed under CC0 1.0
// http://creativecommons.org/publicdomain/zero/1.0/

package prng

import (
	"encoding/binary"
	"io"
	"math/bits"
)

// SplitMix64 is the splitmix64 PRNG from http://prng.di.unimi.it/splitmix64.c.
// The zero value is usable directly. SplitMix64 is primarily provided to support
// seeding the xoshiro PRNGs.
type SplitMix64 struct {
	state uint64
}

// NewSplitMix64 returns a new pseudo-random splitmix64 source seeded
// with the given value.
func NewSplitMix64(seed uint64) *SplitMix64 {
	var src SplitMix64
	src.Seed(seed)
	return &src
}

// Seed uses the provided seed value to initialize the generator to a
// deterministic state.
func (src *SplitMix64) Seed(seed uint64) {
	src.state = seed
}

// Uint64 returns a pseudo-random 64-bit unsigned integer as a uint64.
func (src *SplitMix64) Uint64() uint64 {
	src.state += 0x9e3779b97f4a7c15
	z := src.state
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}

// MarshalBinary returns the binary representation of the current state of the generator.
func (src *SplitMix64) MarshalBinary() ([]byte, error) {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], src.state)
	return buf[:], nil
}

// UnmarshalBinary sets the state of the generator to the state represented in data.
func (src *SplitMix64) UnmarshalBinary(data []byte) error {
	if len(data) < 8 {
		return io.ErrUnexpectedEOF
	}
	src.state = binary.BigEndian.Uint64(data)
	return nil
}

// Xoshiro256plus is the xoshiro256+ 1.0 PRNG from http://prng.di.unimi.it/xoshiro256plus.c.
// The xoshiro PRNGs are described in http://vigna.di.unimi.it/ftp/papers/ScrambledLinear.pdf
// and http://prng.di.unimi.it/.
// A Xoshiro256plus value is only valid if returned by NewXoshiro256plus.
type Xoshiro256plus struct {
	state [4]uint64
}

// NewXoshiro256plus returns a new pseudo-random xoshiro256+ source
// seeded with the given value.
func NewXoshiro256plus(seed uint64) *Xoshiro256plus {
	var src Xoshiro256plus
	src.Seed(seed)
	return &src
}

// Seed uses the provided seed value to initialize the generator to a
// deterministic state.
func (src *Xoshiro256plus) Seed(seed uint64) {
	var boot SplitMix64
	boot.Seed(seed)
	for i := range src.state {
		src.state[i] = boot.Uint64()
	}
}

// Uint64 returns a pseudo-random 64-bit unsigned integer as a uint64.
func (src *Xoshiro256plus) Uint64() uint64 {
	result := src.state[0] + src.state[3]

	t := src.state[1] << 17

	src.state[2] ^= src.state[0]
	src.state[3] ^= src.state[1]
	src.state[1] ^= src.state[2]
	src.state[0] ^= src.state[3]

	src.state[2] ^= t

	src.state[3] = bits.RotateLeft64(src.state[3], 45)

	return result
}

// MarshalBinary returns the binary representation of the current state of the generator.
func (src *Xoshiro256plus) MarshalBinary() ([]byte, error) {
	var buf [32]byte
	for i := 0; i < 4; i++ {
		binary.BigEndian.PutUint64(buf[i*8:(i+1)*8], src.state[i])
	}
	return buf[:], nil
}

// UnmarshalBinary sets the state of the generator to the state represented in data.
func (src *Xoshiro256plus) UnmarshalBinary(data []byte) error {
	if len(data) < 32 {
		return io.ErrUnexpectedEOF
	}
	for i := 0; i < 4; i++ {
		src.state[i] = binary.BigEndian.Uint64(data[i*8 : (i+1)*8])
	}
	return nil
}

// Xoshiro256plusplus is the xoshiro256++ 1.0 PRNG from http://prng.di.unimi.it/xoshiro256plusplus.c.
// The xoshiro PRNGs are described in http://vigna.di.unimi.it/ftp/papers/ScrambledLinear.pdf
// and http://prng.di.unimi.it/.
// A Xoshiro256plusplus value is only valid if returned by NewXoshiro256plusplus.
type Xoshiro256plusplus struct {
	state [4]uint64
}

// NewXoshiro256plusplus returns a new pseudo-random xoshiro256++ source
// seeded with the given value.
func NewXoshiro256plusplus(seed uint64) *Xoshiro256plusplus {
	var src Xoshiro256plusplus
	src.Seed(seed)
	return &src
}

// Seed uses the provided seed value to initialize the generator to a
// deterministic state.
func (src *Xoshiro256plusplus) Seed(seed uint64) {
	var boot SplitMix64
	boot.Seed(seed)
	for i := range src.state {
		src.state[i] = boot.Uint64()
	}
}

// Uint64 returns a pseudo-random 64-bit unsigned integer as a uint64.
func (src *Xoshiro256plusplus) Uint64() uint64 {
	result := bits.RotateLeft64(src.state[0]+src.state[3], 23) + src.state[0]

	t := src.state[1] << 17

	src.state[2] ^= src.state[0]
	src.state[3] ^= src.state[1]
	src.state[1] ^= src.state[2]
	src.state[0] ^= src.state[3]

	src.state[2] ^= t

	src.state[3] = bits.RotateLeft64(src.state[3], 45)

	return result
}

// MarshalBinary returns the binary representation of the current state of the generator.
func (src *Xoshiro256plusplus) MarshalBinary() ([]byte, error) {
	var buf [32]byte
	for i := 0; i < 4; i++ {
		binary.BigEndian.PutUint64(buf[i*8:(i+1)*8], src.state[i])
	}
	return buf[:], nil
}

// UnmarshalBinary sets the state of the generator to the state represented in data.
func (src *Xoshiro256plusplus) UnmarshalBinary(data []byte) error {
	if len(data) < 32 {
		return io.ErrUnexpectedEOF
	}
	for i := 0; i < 4; i++ {
		src.state[i] = binary.BigEndian.Uint64(data[i*8 : (i+1)*8])
	}
	return nil
}

// Xoshiro256starstar is the xoshiro256** 1.0 PRNG from http://prng.di.unimi.it/xoshiro256starstar.c.
// The xoshiro PRNGs are described in http://vigna.di.unimi.it/ftp/papers/ScrambledLinear.pdf
// and http://prng.di.unimi.it/.
// A Xoshiro256starstar value is only valid if returned by NewXoshiro256starstar.
type Xoshiro256starstar struct {
	state [4]uint64
}

// NewXoshiro256starstar returns a new pseudo-random xoshiro256** source
// seeded with the given value.
func NewXoshiro256starstar(seed uint64) *Xoshiro256starstar {
	var src Xoshiro256starstar
	src.Seed(seed)
	return &src
}

// Seed uses the provided seed value to initialize the generator to a
// deterministic state.
func (src *Xoshiro256starstar) Seed(seed uint64) {
	var boot SplitMix64
	boot.Seed(seed)
	for i := range src.state {
		src.state[i] = boot.Uint64()
	}
}

// Uint64 returns a pseudo-random 64-bit unsigned integer as a uint64.
func (src *Xoshiro256starstar) Uint64() uint64 {
	result := bits.RotateLeft64(src.state[1]*5, 7) * 9

	t := src.state[1] << 17

	src.state[2] ^= src.state[0]
	src.state[3] ^= src.state[1]
	src.state[1] ^= src.state[2]
	src.state[0] ^= src.state[3]

	src.state[2] ^= t

	src.state[3] = bits.RotateLeft64(src.state[3], 45)

	return result
}

// MarshalBinary returns the binary representation of the current state of the generator.
func (src *Xoshiro256starstar) MarshalBinary() ([]byte, error) {
	var buf [32]byte
	for i := 0; i < 4; i++ {
		binary.BigEndian.PutUint64(buf[i*8:(i+1)*8], src.state[i])
	}
	return buf[:], nil
}

// UnmarshalBinary sets the state of the generator to the state represented in data.
func (src *Xoshiro256starstar) UnmarshalBinary(data []byte) error {
	if len(data) < 32 {
		return io.ErrUnexpectedEOF
	}
	for i := 0; i < 4; i++ {
		src.state[i] = binary.BigEndian.Uint64(data[i*8 : (i+1)*8])
	}
	return nil
}
