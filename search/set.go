// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

import (
	"unsafe"
)

// same determines whether two sets are backed by the same store. In the
// current implementation using hash maps it makes use of the fact that
// hash maps (at least in the gc implementation) are passed as a pointer
// to a runtime Hmap struct.
//
// A map is not seen by the runtime as a pointer though, so we cannot
// directly compare the sets converted to unsafe.Pointer and need to take
// the sets' addressed and dereference them as pointers to some comparable
// type.
func same(s1, s2 Set) bool {
	return *(*uintptr)(unsafe.Pointer(&s1)) == *(*uintptr)(unsafe.Pointer(&s2))
}

// A set is a set of integer identifiers.
type Set map[int]struct{}

// The simple accessor methods for Set are provided to allow ease of
// implementation change should the need arise.

// add inserts an element into the set.
func (s Set) add(e int) {
	s[e] = struct{}{}
}

// has reports the existence of the element in the set.
func (s Set) has(e int) bool {
	_, ok := s[e]
	return ok
}

// remove delete the specified element from the set.
func (s Set) remove(e int) {
	delete(s, e)
}

// count reports the number of elements stored in the set.
func (s Set) count() int {
	return len(s)
}

// elements returns a newly created slice containing all elements of the
// set. The order of elements in the slice is unspecified.
func (s Set) elements() []int {
	els := make([]int, 0, len(s))
	for e, _ := range s {
		els = append(els, e)
	}

	return els
}

// clear returns an empty set, possibly using the same backing store.
// clear is not provided as a method since there is no way to replace
// the calling value if clearing is performed by a make(set). clear
// should never be called without keeping the returned value.
func clear(s Set) Set {
	if len(s) == 0 {
		return s
	}

	return make(Set)
}

// copy performs a perfect copy from s1 to dst (meaning the sets will
// be equal).
func (dst Set) copy(src Set) Set {
	if same(src, dst) {
		return dst
	}

	if len(dst) > 0 {
		dst = make(Set, len(src))
	}

	for e := range src {
		dst[e] = struct{}{}
	}

	return dst
}

// equal reports set equality between the parameters. Sets are equal if
// and only if they have the same elements.
func equal(s1, s2 Set) bool {
	if same(s1, s2) {
		return true
	}

	if len(s1) != len(s2) {
		return false
	}

	for e := range s1 {
		if _, ok := s2[e]; !ok {
			return false
		}
	}

	return true
}

// union takes the union of s1 and s2, and stores it in dst.
//
// The union of two sets, s1 and s2, is the set containing all the
// elements of each, for instance:
//
//     {a,b,c} UNION {d,e,f} = {a,b,c,d,e,f}
//
// Since sets may not have repetition, unions of two sets that overlap
// do not contain repeat elements, that is:
//
//     {a,b,c} UNION {b,c,d} = {a,b,c,d}
//
func (dst Set) union(s1, s2 Set) Set {
	if same(s1, s2) {
		return dst.copy(s1)
	}

	if !same(s1, dst) && !same(s2, dst) {
		dst = clear(dst)
	}

	if !same(dst, s1) {
		for e := range s1 {
			dst[e] = struct{}{}
		}
	}

	if !same(dst, s2) {
		for e := range s2 {
			dst[e] = struct{}{}
		}
	}

	return dst
}

// intersect takes the intersection of s1 and s2, and stores it in dst.
//
// The intersection of two sets, s1 and s2, is the set containing all
// the elements shared between the two sets, for instance:
//
//     {a,b,c} INTERSECT {b,c,d} = {b,c}
//
// The intersection between a set and itself is itself, and thus
// effectively a copy operation:
//
//     {a,b,c} INTERSECT {a,b,c} = {a,b,c}
//
// The intersection between two sets that share no elements is the empty
// set:
//
//     {a,b,c} INTERSECT {d,e,f} = {}
//
func (dst Set) intersect(s1, s2 Set) Set {
	var swap Set

	if same(s1, s2) {
		return dst.copy(s1)
	}
	if same(s1, dst) {
		swap = s2
	} else if same(s2, dst) {
		swap = s1
	} else {
		dst = clear(dst)

		if len(s1) > len(s2) {
			s1, s2 = s2, s1
		}

		for e := range s1 {
			if _, ok := s2[e]; ok {
				dst[e] = struct{}{}
			}
		}

		return dst
	}

	for e := range dst {
		if _, ok := swap[e]; !ok {
			delete(dst, e)
		}
	}

	return dst
}
