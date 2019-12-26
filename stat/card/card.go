// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate ./generate_64bit.sh

package card

import (
	"fmt"
	"hash"
	"math"
	"reflect"
	"sync"
)

const (
	w32 = 32
	w64 = 64
)

func alpha(m uint64) float64 {
	if m < 128 {
		return alphaValues[m]
	}
	return 0.7213 / (1 + 1.079/float64(m))
}

var alphaValues = [...]float64{
	16: 0.673,
	32: 0.697,
	64: 0.709,
}

func linearCounting(m, v float64) float64 {
	return m * (math.Log(m) - math.Log(v))
}

func max(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}

func min(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}

func typeNameOf(v interface{}) string {
	t := reflect.TypeOf(v)
	var prefix string
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		prefix = "*"
	}
	if t.PkgPath() == "" {
		return prefix + t.Name()
	}
	return prefix + t.PkgPath() + "." + t.Name()
}

// hashes holds registered hashes.
var hashes sync.Map // map[string]userType

type userType struct {
	fn  reflect.Value // Holds a func() hash.Hash{32,64}.
	typ reflect.Type  // Type of the returned hash implementation.
}

// RegisterHash registers a function that returns a new hash.Hash32 or hash.Hash64
// to the name of the type implementing the interface. The value of fn must be a
// func() hash.Hash32 or func() hash.Hash64, otherwise RegisterHash will panic.
// RegisterHash will panic if there is not a unique mapping from the name to the
// returned type.
func RegisterHash(fn interface{}) {
	const invalidType = "card: must register func() hash.Hash32 or func() hash.Hash64"

	rf := reflect.ValueOf(fn)
	rt := rf.Type()
	if rf.Kind() != reflect.Func {
		panic(invalidType)
	}
	if rt.NumIn() != 0 {
		panic(invalidType)
	}
	if rt.NumOut() != 1 {
		panic(invalidType)
	}
	h := rf.Call(nil)[0].Interface()
	var name string
	var h32 hash.Hash32
	var h64 hash.Hash64
	switch rf.Type().Out(0) {
	case reflect.TypeOf(&h32).Elem(), reflect.TypeOf(&h64).Elem():
		name = typeNameOf(h)
	default:
		panic(invalidType)
	}
	user := userType{fn: rf, typ: reflect.TypeOf(h)}
	ut, dup := hashes.LoadOrStore(name, user)
	stored := ut.(userType)
	if dup && stored.typ != user.typ {
		panic(fmt.Sprintf("card: registering duplicate types for %q: %s != %s", name, stored.typ, user.typ))
	}
}
