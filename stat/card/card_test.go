// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package card

import (
	"encoding"
	"fmt"
	"hash"
	"hash/fnv"
	"io"
	"strconv"
	"strings"
	"sync"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats/scalar"
)

// exact is an exact cardinality accumulator.
type exact map[string]struct{}

func (e exact) Write(b []byte) (int, error) {
	if _, exists := e[string(b)]; exists {
		return len(b), nil
	}
	e[string(b)] = struct{}{}
	return len(b), nil
}

func (e exact) Count() float64 {
	return float64(len(e))
}

type counter interface {
	io.Writer
	Count() float64
}

var counterTests = []struct {
	name    string
	count   float64
	counter func() counter
	tol     float64
}{
	{name: "exact-1e5", count: 1e5, counter: func() counter { return make(exact) }, tol: 0},

	{name: "HyperLogLog32-0-10-FNV-1a", count: 0, counter: func() counter { return mustCounter(NewHyperLogLog32(10, fnv.New32a())) }, tol: 0},
	{name: "HyperLogLog64-0-10-FNV-1a", count: 0, counter: func() counter { return mustCounter(NewHyperLogLog64(10, fnv.New64a())) }, tol: 0},
	{name: "HyperLogLog32-10-14-FNV-1a", count: 10, counter: func() counter { return mustCounter(NewHyperLogLog32(14, fnv.New32a())) }, tol: 0.0005},
	{name: "HyperLogLog32-1e3-4-FNV-1a", count: 1e3, counter: func() counter { return mustCounter(NewHyperLogLog32(4, fnv.New32a())) }, tol: 0.1},
	{name: "HyperLogLog32-1e4-6-FNV-1a", count: 1e4, counter: func() counter { return mustCounter(NewHyperLogLog32(6, fnv.New32a())) }, tol: 0.06},
	{name: "HyperLogLog32-1e7-8-FNV-1a", count: 1e7, counter: func() counter { return mustCounter(NewHyperLogLog32(8, fnv.New32a())) }, tol: 0.03},
	{name: "HyperLogLog64-1e7-8-FNV-1a", count: 1e7, counter: func() counter { return mustCounter(NewHyperLogLog64(8, fnv.New64a())) }, tol: 0.07},
	{name: "HyperLogLog32-1e7-10-FNV-1a", count: 1e7, counter: func() counter { return mustCounter(NewHyperLogLog32(10, fnv.New32a())) }, tol: 0.06},
	{name: "HyperLogLog64-1e7-10-FNV-1a", count: 1e7, counter: func() counter { return mustCounter(NewHyperLogLog64(10, fnv.New64a())) }, tol: 0.02},
	{name: "HyperLogLog32-1e7-14-FNV-1a", count: 1e7, counter: func() counter { return mustCounter(NewHyperLogLog32(14, fnv.New32a())) }, tol: 0.005},
	{name: "HyperLogLog64-1e7-14-FNV-1a", count: 1e7, counter: func() counter { return mustCounter(NewHyperLogLog64(14, fnv.New64a())) }, tol: 0.002},
	{name: "HyperLogLog32-1e7-16-FNV-1a", count: 1e7, counter: func() counter { return mustCounter(NewHyperLogLog32(16, fnv.New32a())) }, tol: 0.005},
	{name: "HyperLogLog64-1e7-16-FNV-1a", count: 1e7, counter: func() counter { return mustCounter(NewHyperLogLog64(16, fnv.New64a())) }, tol: 0.002},
	{name: "HyperLogLog64-1e7-20-FNV-1a", count: 1e7, counter: func() counter { return mustCounter(NewHyperLogLog64(20, fnv.New64a())) }, tol: 0.001},
	{name: "HyperLogLog64-1e3-20-FNV-1a", count: 1e3, counter: func() counter { return mustCounter(NewHyperLogLog64(20, fnv.New64a())) }, tol: 0.001},
}

func mustCounter(c counter, err error) counter {
	if err != nil {
		panic(fmt.Sprintf("bad test: %v", err))
	}
	return c
}

func TestCounters(t *testing.T) {
	t.Parallel()

	for _, test := range counterTests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			rnd := rand.New(rand.NewSource(1))
			var dst []byte
			c := test.counter()
			for i := 0; i < int(test.count); i++ {
				dst = strconv.AppendUint(dst[:0], rnd.Uint64(), 16)
				dst = append(dst, '-')
				dst = strconv.AppendUint(dst, uint64(i), 16)
				n, err := c.Write(dst)
				if n != len(dst) {
					t.Errorf("unexpected number of bytes written for %s: got:%d want:%d",
						test.name, n, len(dst))
					break
				}
				if err != nil {
					t.Errorf("unexpected error for %s: %v", test.name, err)
					break
				}
			}

			if got := c.Count(); !scalar.EqualWithinRel(got, test.count, test.tol) {
				t.Errorf("unexpected count for %s: got:%.0f want:%.0f", test.name, got, test.count)
			}
		})
	}
}

func TestUnion(t *testing.T) {
	t.Parallel()

	for _, test := range counterTests {
		if strings.HasPrefix(test.name, "exact") {
			continue
		}
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			rnd := rand.New(rand.NewSource(1))
			var dst []byte
			var cs [2]counter
			for j := range cs {
				cs[j] = test.counter()
				for i := 0; i < int(test.count); i++ {
					dst = strconv.AppendUint(dst[:0], rnd.Uint64(), 16)
					dst = append(dst, '-')
					dst = strconv.AppendUint(dst, uint64(i), 16)
					n, err := cs[j].Write(dst)
					if n != len(dst) {
						t.Errorf("unexpected number of bytes written for %s: got:%d want:%d",
							test.name, n, len(dst))
						break
					}
					if err != nil {
						t.Errorf("unexpected error for %s: %v", test.name, err)
						break
					}
				}
			}

			u := test.counter()
			var err error
			switch u := u.(type) {
			case *HyperLogLog32:
				err = u.Union(cs[0].(*HyperLogLog32), cs[1].(*HyperLogLog32))
			case *HyperLogLog64:
				err = u.Union(cs[0].(*HyperLogLog64), cs[1].(*HyperLogLog64))
			}
			if err != nil {
				t.Errorf("unexpected error from Union call: %v", err)
			}
			if got := u.Count(); !scalar.EqualWithinRel(got, 2*test.count, 2*test.tol) {
				t.Errorf("unexpected count for %s: got:%.0f want:%.0f", test.name, got, 2*test.count)
			}
		})
	}
}

type resetCounter interface {
	counter
	Reset()
}

var counterResetTests = []struct {
	name         string
	count        int
	resetCounter func() resetCounter
}{
	{name: "HyperLogLog32-1e3-4-FNV-1a", count: 1e3, resetCounter: func() resetCounter { return mustResetCounter(NewHyperLogLog32(4, fnv.New32a())) }},
	{name: "HyperLogLog64-1e3-4-FNV-1a", count: 1e3, resetCounter: func() resetCounter { return mustResetCounter(NewHyperLogLog64(4, fnv.New64a())) }},
	{name: "HyperLogLog32-1e4-6-FNV-1a", count: 1e4, resetCounter: func() resetCounter { return mustResetCounter(NewHyperLogLog32(6, fnv.New32a())) }},
	{name: "HyperLogLog64-1e4-6-FNV-1a", count: 1e4, resetCounter: func() resetCounter { return mustResetCounter(NewHyperLogLog64(6, fnv.New64a())) }},
}

func mustResetCounter(c resetCounter, err error) resetCounter {
	if err != nil {
		panic(fmt.Sprintf("bad test: %v", err))
	}
	return c
}

func TestResetCounters(t *testing.T) {
	t.Parallel()

	var dst []byte
	for _, test := range counterResetTests {
		c := test.resetCounter()
		var counts [2]float64
		for k := range counts {
			rnd := rand.New(rand.NewSource(1))
			for i := 0; i < int(test.count); i++ {
				dst = strconv.AppendUint(dst[:0], rnd.Uint64(), 16)
				dst = append(dst, '-')
				dst = strconv.AppendUint(dst, uint64(i), 16)
				n, err := c.Write(dst)
				if n != len(dst) {
					t.Errorf("unexpected number of bytes written for %s: got:%d want:%d",
						test.name, n, len(dst))
					break
				}
				if err != nil {
					t.Errorf("unexpected error for %s: %v", test.name, err)
					break
				}
			}
			counts[k] = c.Count()
			c.Reset()
		}

		if counts[0] != counts[1] {
			t.Errorf("unexpected counts for %s after reset: got:%.0f", test.name, counts)
		}
	}
}

type counterEncoder interface {
	counter
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

var counterEncoderTests = []struct {
	name           string
	count          int
	src, dst, zdst func() counterEncoder
}{
	{
		name: "HyperLogLog32-4-4-FNV-1a", count: 1e3,
		src:  func() counterEncoder { return mustCounterEncoder(NewHyperLogLog32(4, fnv.New32a())) },
		dst:  func() counterEncoder { return mustCounterEncoder(NewHyperLogLog32(4, fnv.New32a())) },
		zdst: func() counterEncoder { return &HyperLogLog32{} },
	},
	{
		name: "HyperLogLog32-4-8-FNV-1a", count: 1e3,
		src:  func() counterEncoder { return mustCounterEncoder(NewHyperLogLog32(4, fnv.New32a())) },
		dst:  func() counterEncoder { return mustCounterEncoder(NewHyperLogLog32(8, fnv.New32a())) },
		zdst: func() counterEncoder { return &HyperLogLog32{} },
	},
	{
		name: "HyperLogLog32-8-4-FNV-1a", count: 1e3,
		src:  func() counterEncoder { return mustCounterEncoder(NewHyperLogLog32(8, fnv.New32a())) },
		dst:  func() counterEncoder { return mustCounterEncoder(NewHyperLogLog32(4, fnv.New32a())) },
		zdst: func() counterEncoder { return &HyperLogLog32{} },
	},
	{
		name: "HyperLogLog64-4-4-FNV-1a", count: 1e3,
		src:  func() counterEncoder { return mustCounterEncoder(NewHyperLogLog64(4, fnv.New64a())) },
		dst:  func() counterEncoder { return mustCounterEncoder(NewHyperLogLog64(4, fnv.New64a())) },
		zdst: func() counterEncoder { return &HyperLogLog64{} },
	},
	{
		name: "HyperLogLog64-4-8-FNV-1a", count: 1e3,
		src:  func() counterEncoder { return mustCounterEncoder(NewHyperLogLog64(4, fnv.New64a())) },
		dst:  func() counterEncoder { return mustCounterEncoder(NewHyperLogLog64(8, fnv.New64a())) },
		zdst: func() counterEncoder { return &HyperLogLog64{} },
	},
	{
		name: "HyperLogLog64-8-4-FNV-1a", count: 1e3,
		src:  func() counterEncoder { return mustCounterEncoder(NewHyperLogLog64(8, fnv.New64a())) },
		dst:  func() counterEncoder { return mustCounterEncoder(NewHyperLogLog64(4, fnv.New64a())) },
		zdst: func() counterEncoder { return &HyperLogLog64{} },
	},
}

func mustCounterEncoder(c counterEncoder, err error) counterEncoder {
	if err != nil {
		panic(fmt.Sprintf("bad test: %v", err))
	}
	return c
}

func TestBinaryEncoding(t *testing.T) {
	t.Parallel()

	RegisterHash(fnv.New32a)
	RegisterHash(fnv.New64a)
	defer func() {
		hashes = sync.Map{}
	}()
	for _, test := range counterEncoderTests {
		rnd := rand.New(rand.NewSource(1))
		src := test.src()
		for i := 0; i < int(test.count); i++ {
			buf := strconv.AppendUint(nil, rnd.Uint64(), 16)
			buf = append(buf, '-')
			buf = strconv.AppendUint(buf, uint64(i), 16)
			n, err := src.Write(buf)
			if n != len(buf) {
				t.Errorf("unexpected number of bytes written for %s: got:%d want:%d",
					test.name, n, len(buf))
				break
			}
			if err != nil {
				t.Errorf("unexpected error for %s: %v", test.name, err)
				break
			}
		}

		buf, err := src.MarshalBinary()
		if err != nil {
			t.Errorf("unexpected error marshaling binary for %s: %v", test.name, err)
			continue
		}
		dst := test.dst()
		err = dst.UnmarshalBinary(buf)
		if err != nil {
			t.Errorf("unexpected error unmarshaling binary for %s: %v", test.name, err)
			continue
		}
		zdst := test.zdst()
		err = zdst.UnmarshalBinary(buf)
		if err != nil {
			t.Errorf("unexpected error unmarshaling binary into zero receiver for %s: %v", test.name, err)
			continue
		}
		gotSrc := src.Count()
		gotDst := dst.Count()
		gotZdst := zdst.Count()

		if gotSrc != gotDst {
			t.Errorf("unexpected count for %s: got:%.0f want:%.0f", test.name, gotDst, gotSrc)
		}
		if gotSrc != gotZdst {
			t.Errorf("unexpected count for %s into zero receiver: got:%.0f want:%.0f", test.name, gotZdst, gotSrc)
		}
	}
}

var invalidRegisterTests = []struct {
	fn     interface{}
	panics bool
}{
	{fn: int(0), panics: true},
	{fn: func() {}, panics: true},
	{fn: func(int) {}, panics: true},
	{fn: func() int { return 0 }, panics: true},
	{fn: func() hash.Hash { return fnv.New32a() }, panics: true},
	{fn: func() hash.Hash32 { return fnv.New32a() }, panics: false},
	{fn: func() hash.Hash { return fnv.New64a() }, panics: true},
	{fn: func() hash.Hash64 { return fnv.New64a() }, panics: false},
}

func TestRegisterInvalid(t *testing.T) {
	t.Parallel()

	for _, test := range invalidRegisterTests {
		var r interface{}
		func() {
			defer func() {
				r = recover()
			}()
			RegisterHash(test.fn)
		}()
		panicked := r != nil
		if panicked != test.panics {
			if panicked {
				t.Errorf("unexpected panic for %T", test.fn)
			} else {
				t.Errorf("expected panic for %T", test.fn)
			}
		}
	}
}

var rhoQTests = []struct {
	bits uint
	q    uint8
	want uint8
}{
	{bits: 0xff, q: 8, want: 1},
	{bits: 0xfe, q: 8, want: 1},
	{bits: 0x0f, q: 8, want: 5},
	{bits: 0x1f, q: 8, want: 4},
	{bits: 0x00, q: 8, want: 9},
}

func TestRhoQ(t *testing.T) {
	t.Parallel()

	for _, test := range rhoQTests {
		got := rho32q(uint32(test.bits), test.q)
		if got != test.want {
			t.Errorf("unexpected rho32q for %0*b: got:%d want:%d", test.q, test.bits, got, test.want)
		}
		got = rho64q(uint64(test.bits), test.q)
		if got != test.want {
			t.Errorf("unexpected rho64q for %0*b: got:%d want:%d", test.q, test.bits, got, test.want)
		}
	}
}

var counterBenchmarks = []struct {
	name    string
	count   int
	counter func() counter
}{
	{name: "exact-1e6", count: 1e6, counter: func() counter { return make(exact) }},
	{name: "HyperLogLog32-1e6-8-FNV-1a", count: 1e6, counter: func() counter { return mustCounter(NewHyperLogLog32(8, fnv.New32a())) }},
	{name: "HyperLogLog64-1e6-8-FNV-1a", count: 1e6, counter: func() counter { return mustCounter(NewHyperLogLog64(8, fnv.New64a())) }},
	{name: "HyperLogLog32-1e6-16-FNV-1a", count: 1e6, counter: func() counter { return mustCounter(NewHyperLogLog32(16, fnv.New32a())) }},
	{name: "HyperLogLog64-1e6-16-FNV-1a", count: 1e6, counter: func() counter { return mustCounter(NewHyperLogLog64(16, fnv.New64a())) }},
}

func BenchmarkCounters(b *testing.B) {
	for _, bench := range counterBenchmarks {
		c := bench.counter()
		rnd := rand.New(rand.NewSource(1))
		var dst []byte
		b.Run(bench.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := 0; j < int(bench.count); j++ {
					dst = strconv.AppendUint(dst[:0], rnd.Uint64(), 16)
					dst = append(dst, '-')
					dst = strconv.AppendUint(dst, uint64(j), 16)
					_, _ = c.Write(dst)
				}
			}
			_ = c.Count()
		})
	}
}
