// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package network

const (
	A = iota
	B
	C
	D
	E
	F
	G
	H
	I
	J
	K
)

// set is an integer set.
type set map[int64]struct{}

func linksTo(i ...int64) set {
	if len(i) == 0 {
		return nil
	}
	s := make(set)
	for _, v := range i {
		s[v] = struct{}{}
	}
	return s
}
