// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

// intSet is a set of integer identifiers.
type intSet map[int]struct{}

// Returns true if e is an element of s.
func (s intSet) Contains(e int) bool {
	_, ok := s[e]
	return ok
}

// Adds the element e to s1.
func (s intSet) Add(e int) {
	s[e] = struct{}{}
}

// Removes the element e from s1.
func (s intSet) Remove(e int) {
	delete(s, e)
}

func (s intSet) Cardinality() int {
	return len(s)
}
