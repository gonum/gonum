// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package set

// Set is a set of integer identifiers.
type Set map[int]struct{}

func NewSet() *Set {
	s := make(Set)
	return &s
}

func (s1 *Set) Clear() *Set {
	if len(*s1) == 0 {
		return s1
	}

	*s1 = make(Set)

	return s1
}

// Ensures a perfect copy from s1 to dst (meaning the sets will be equal)
func (dst *Set) Copy(src *Set) *Set {
	if src == dst {
		return dst
	}

	d := *dst
	if len(d) > 0 {
		d = make(Set, len(*src))
	}

	for e := range *src {
		d[e] = struct{}{}
	}
	*dst = d

	return dst
}

// If every element in s1 is also in s2 (and vice versa), the sets are deemed equal.
func Equal(s1, s2 *Set) bool {
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
func (dst *Set) Union(s1, s2 *Set) *Set {
	if s1 == s2 {
		return dst.Copy(s1)
	}

	if s1 != dst && s2 != dst {
		dst.Clear()
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
func (dst *Set) Intersection(s1, s2 *Set) *Set {
	var swap *Set

	if s1 == s2 {
		return dst.Copy(s1)
	}
	if s1 == dst {
		swap = s2
	} else if s2 == dst {
		swap = s1
	} else {
		dst.Clear()

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

// Takes the difference (-) of s1 and s2 and stores it in dst.
//
// The difference (-) between two sets, s1 and s2, is all the elements in s1 that are NOT also
// in s2.
//
//     {a,b,c} - {b,c,d} = {a}
//
// The difference between two identical sets is the empty set:
//
//     {a,b,c} - {a,b,c} = {}
//
// The difference between two sets with no overlapping elements is s1:
//
//     {a,b,c} - {d,e,f} = {a,b,c}
//
// Implementation note: if dst == s2 (meaning they have identical pointers), a temporary set must
// be used to store the data and then copied over, thus s2.Diff(s1,s2) has an extra allocation and
// may cause worse performance in some cases.
func (dst *Set) Diff(s1, s2 *Set) *Set {
	if s1 == s2 {
		return dst.Clear()
	}

	if s2 == dst {
		tmp := NewSet()

		tmp.Diff(s1, s2)
		*dst = *tmp
		return dst
	}

	s2m := *s2
	if s1 == dst {
		d := *dst
		for e := range d {
			if _, ok := s2m[e]; ok {
				delete(d, e)
			}
		}
	} else {
		dst.Clear()
		d := *dst
		for e := range *s1 {
			if _, ok := s2m[e]; !ok {
				d[e] = struct{}{}
			}
		}
	}

	return dst
}

// Returns true if s1 is an improper subset of s2.
//
// An improper subset occurs when every element in s1 is also in s2 OR s1 and s2 are equal:
//
//     {a,b,c}   SUBSET {a,b,c} = true
//     {a,b}     SUBSET {a,b,c} = true
//     {c,d}     SUBSET {a,b,c} = false
//     {a,b,c,d} SUBSET {a,b,c} = false
//
// Special case: The empty set is a subset of everything:
//
// 	   {} SUBSET {a,b} = true
//     {} SUBSET {}    = true
//
// In the case where one needs to test if s1 is smaller than s2, but not equal, use ProperSubset.
func Subset(s1, s2 *Set) bool {
	if len(*s1) > len(*s2) {
		return false
	}
	if s1 == s2 || len(*s1) == 0 {
		return true
	}

	s2m := *s2
	for e, _ := range *s1 {
		if _, ok := s2m[e]; !ok {
			return false
		}
	}

	return true
}

// Returns true if s1 is a proper subset of s2.
// A proper subset is when every element of s1 is in s2, but s1 is smaller than s2 (i.e. they are
// not equal):
//
//     {a,b,c}   PROPER SUBSET {a,b,c} = false
//     {a,b}     PROPER SUBSET {a,b,c} = true
//     {c,d}     PROPER SUBSET {a,b,c} = false
//     {a,b,c,d} PROPER SUBSET {a,b,c} = false
//
// Special case: The empty set is a proper subset of everything (except itself):
//
//      {} PROPER SUBSET {a,b} = true
//      {} PROPER SUBSET {}    = false
//
// When equality is allowed, use Subset.
func ProperSubset(s1, s2 *Set) bool {
	if len(*s1) >= len(*s2) { // implicitly tests if s1 and s2 are both the empty set
		return false
	} else if len(*s1) == 0 {
		return true
	} // We can eschew the s1 == s2 because if they are the same their lens are equal anyway

	s2m := *s2
	for e, _ := range *s1 {
		if _, ok := s2m[e]; !ok {
			return false
		}
	}

	return true
}

// Returns true if e is an element of s.
func (s *Set) Contains(e int) bool {
	_, ok := (*s)[e]
	return ok
}

// Adds the element e to s1.
func (s1 *Set) Add(e int) {
	(*s1)[e] = struct{}{}
}

func (s1 *Set) AddAll(es ...int) {
	for _, e := range es {
		s1.Add(e)
	}
}

// Removes the element e from s1.
func (s1 *Set) Remove(e int) {
	delete(*s1, e)
}

// Returns the number of elements in s1.
func (s1 *Set) Cardinality() int {
	return len(*s1)
}

func (s1 *Set) Elements() []int {
	els := make([]int, 0, len(*s1))
	for e, _ := range *s1 {
		els = append(els, e)
	}

	return els
}
