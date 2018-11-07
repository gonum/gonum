// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

type TemporalLine interface {
	Line

	// Interval returns the edge starting time and ending time. The
	// edge traversal time is the difference between the two.
	Interval() (start, end uint64)
}