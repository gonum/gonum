// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package curve defines space filling curves. A space filling curve is a curve
// whose range contains the entirety of a finite k-dimensional space. Space
// filling curves can be used to map between a linear space and a 2D, 3D, or 4D
// space.
//
// # Hilbert Curves
//
// The Hilbert curve is a continuous, space filling curve first described by
// David Hilbert. The implementation of Hilbert2D is based on example code from
// the [Hilbert curve (Wikipedia)]. The implementation of Hilbert3D and
// Hilbert4D are extrapolated from Hilbert2D.
//
// Technically, a Hilbert curve is a continuous, fractal, space-filling curve of
// infinite length, constructed as the limit of a series of piecewise linear
// curves. We refer to the kth piecewise linear curve as the k-order Hilbert
// curve. The first-order 2D Hilbert curve looks like ⊐. The k-order
// n-dimensional curve can be constructed from 2ⁿ copies of the (k-1) order
// curve. These (finite) Hilbert curves are iterative/recursive, and the order
// refers to the iteration number.
//
// The runtime of Space and Curve scales as O(n∙k) where n is the dimension and
// k is the order. The length of the curve is 2^(n∙k).
//
// # Limitations
//
// An n-dimensional, k-order Hilbert curve will not be fully usable if 2ⁿᵏ
// overflows int (which is dependent on architecture). Len will overflow if n∙k
// ≥ bits.UintSize-1. Curve will overflow if n∙k > bits.UintSize-1 for some
// values of v. Space will not overflow, but it cannot be called with values
// that do not fit in a signed integer, thus only a subset of the curve can be
// utilized.
//
// [Hilbert curve (Wikipedia)]: https://en.wikipedia.org/w/index.php?title=Hilbert_curve&oldid=1011599190
package curve
