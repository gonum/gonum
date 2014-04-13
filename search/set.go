// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

// simpleSet is a set of integer identifiers.
type simpleSet map[int]struct{}

// Returns true if e is an element of s.
func (s simpleSet) has(e int) bool {
	_, ok := s[e]
	return ok
}

// Adds the element e to s1.
func (s simpleSet) add(e int) {
	s[e] = struct{}{}
}

// Removes the element e from s1.
func (s simpleSet) remove(e int) {
	delete(s, e)
}

func (s simpleSet) count() int {
	return len(s)
}

// Set is a set of integer identifiers.
type set map[int]struct{}

func newSet() *set {
	s := make(set)
	return &s
}

func (s1 *set) clear() *set {
	if len(*s1) == 0 {
		return s1
	}

	*s1 = make(set)

	return s1
}

// Ensures a perfect copy from s1 to dst (meaning the sets will be equal)
func (dst *set) copy(src *set) *set {
	if src == dst {
		return dst
	}

	d := *dst
	if len(d) > 0 {
		d = make(set, len(*src))
	}

	for e := range *src {
		d[e] = struct{}{}
	}
	*dst = d

	return dst
}

// If every element in s1 is also in s2 (and vice versa), the sets are deemed equal.
func equal(s1, s2 *set) bool {
	if s1 == s2 {
		return true
	}
	s2m := *s2
	if len(*s1) != len(s2m) {
		return false
	}

	for e := range *s1 {
		if _, ok := s2m[e]; !ok {
			return false
		}
	}

	return true
}

// Takes the union of s1 and s2, and stores it in dst.
//
// The union of two sets, s1 and s2, is the set containing all the elements of each, for instance:
//
//     {a,b,c} UNION {d,e,f} = {a,b,c,d,e,f}
//
// Since sets may not have repetition, unions of two sets that overlap do not contain repeat
// elements, that is:
//
//     {a,b,c} UNION {b,c,d} = {a,b,c,d}
func (dst *set) union(s1, s2 *set) *set {
	if s1 == s2 {
		return dst.copy(s1)
	}

	if s1 != dst && s2 != dst {
		dst.clear()
	}

	d := *dst
	if dst != s1 {
		for e := range *s1 {
			d[e] = struct{}{}
		}
	}

	if dst != s2 {
		for e := range *s2 {
			d[e] = struct{}{}
		}
	}

	return dst
}

// Takes the intersection of s1 and s2, and stores it in dst
//
// The intersection of two sets, s1 and s2, is the set containing all the elements shared between
// the two sets, for instance:
//
//     {a,b,c} INTERSECT {b,c,d} = {b,c}
//
// The intersection between a set and itself is itself, and thus effectively a copy operation:
//
//     {a,b,c} INTERSECT {a,b,c} = {a,b,c}
//
// The intersection between two sets that share no elements is the empty set:
//
//     {a,b,c} INTERSECT {d,e,f} = {}
func (dst *set) intersect(s1, s2 *set) *set {
	var swap *set

	if s1 == s2 {
		return dst.copy(s1)
	}
	if s1 == dst {
		swap = s2
	} else if s2 == dst {
		swap = s1
	} else {
		dst.clear()

		if len(*s1) > len(*s2) {
			s1, s2 = s2, s1
		}
		s2m := *s2
		d := *dst
		for e := range *s1 {
			if _, ok := s2m[e]; ok {
				d[e] = struct{}{}
			}
		}

		return dst
	}

	d := *dst
	s := *swap
	for e := range d {
		if _, ok := s[e]; !ok {
			delete(d, e)
		}
	}

	return dst
}

func (s *set) add(e int) {
	(*s)[e] = struct{}{}
}

func (s *set) elements() []int {
	els := make([]int, 0, len(*s))
	for e, _ := range *s {
		els = append(els, e)
	}

	return els
}
