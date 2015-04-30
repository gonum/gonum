// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

// A set is a set of integer identifiers.
type IntSet map[int]struct{}

// The simple accessor methods for Set are provided to allow ease of
// implementation change should the need arise.

// add inserts an element into the set.
func (s IntSet) Add(e int) {
	s[e] = struct{}{}
}

// has reports the existence of the element in the set.
func (s IntSet) Has(e int) bool {
	_, ok := s[e]
	return ok
}

// remove deletes the specified element from the set.
func (s IntSet) Remove(e int) {
	delete(s, e)
}

// count reports the number of elements stored in the set.
func (s IntSet) Count() int {
	return len(s)
}
