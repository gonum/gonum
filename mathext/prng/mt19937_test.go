// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prng

import (
	"testing"
	"time"

	"golang.org/x/exp/rand"
)

var _ rand.Source = (*MT19937)(nil)

// Random values in tests are produced by 40 iterations of the C code
// with or without an initial seed array.

func TestMT19937(t *testing.T) {
	t.Parallel()
	want := []uint32{
		3499211612, 581869302, 3890346734, 3586334585, 545404204,
		4161255391, 3922919429, 949333985, 2715962298, 1323567403,
		418932835, 2350294565, 1196140740, 809094426, 2348838239,
		4264392720, 4112460519, 4279768804, 4144164697, 4156218106,
		676943009, 3117454609, 4168664243, 4213834039, 4111000746,
		471852626, 2084672536, 3427838553, 3437178460, 1275731771,
		609397212, 20544909, 1811450929, 483031418, 3933054126,
		2747762695, 3402504553, 3772830893, 4120988587, 2163214728,
	}

	mt := NewMT19937()
	for i := range want {
		got := mt.Uint32()
		if got != want[i] {
			t.Errorf("unexpected random value at iteration %d: got:%d want:%d", i, got, want[i])
		}
	}
}

func TestMT19937SeedFromKeys(t *testing.T) {
	t.Parallel()
	want := []uint32{
		1067595299, 955945823, 477289528, 4107218783, 4228976476,
		3344332714, 3355579695, 227628506, 810200273, 2591290167,
		2560260675, 3242736208, 646746669, 1479517882, 4245472273,
		1143372638, 3863670494, 3221021970, 1773610557, 1138697238,
		1421897700, 1269916527, 2859934041, 1764463362, 3874892047,
		3965319921, 72549643, 2383988930, 2600218693, 3237492380,
		2792901476, 725331109, 605841842, 271258942, 715137098,
		3297999536, 1322965544, 4229579109, 1395091102, 3735697720,
	}

	mt := NewMT19937()
	mt.SeedFromKeys([]uint32{0x123, 0x234, 0x345, 0x456})
	for i := range want {
		got := mt.Uint32()
		if got != want[i] {
			t.Errorf("unexpected random value at iteration %d: got:%d want:%d", i, got, want[i])
		}
	}
}

func TestMT19937RoundTrip(t *testing.T) {
	t.Parallel()
	var src MT19937
	src.Seed(uint64(time.Now().Unix()))

	src.Uint64() // Step PRNG once to makes sure states are mixed.

	buf, err := src.MarshalBinary()
	if err != nil {
		t.Errorf("unexpected error marshaling state: %v", err)
	}

	var dst MT19937
	// Get dst into a non-zero state.
	dst.Seed(1)
	for i := 0; i < 10; i++ {
		dst.Uint64()
	}

	err = dst.UnmarshalBinary(buf)
	if err != nil {
		t.Errorf("unexpected error unmarshaling state: %v", err)
	}

	if dst != src {
		t.Errorf("mismatch between generator states: got:%+v want:%+v", dst, src)
	}
}
