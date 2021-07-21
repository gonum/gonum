// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package curve defines space filling curves. A space filling curve is a curve
// whose range contains the entirety of a finite k-dimensional space. Space
// filling curves can be used to map between a linear space and a 2D, 3D, or 4D
// space.
//
// Hilbert Curves
//
// The Hilbert curve is a continuous, space filling curve first described by
// David Hilbert. The implementation of Hilbert2D is based on example code from
// the Wikipedia article
// (https://en.wikipedia.org/w/index.php?title=Hilbert_curve&oldid=1011599190).
// The implementation of Hilbert3D and Hilbert4D are extrapolated from
// Hilbert2D.
//
// For the first-order n-dimensional Hilbert curve, a spatial point V is mapped
// to a point on the curve D by XOR - each dimension of V is expected to be 0 or
// 1:
//
//     func map1stOrder(n int, v []int) (d int) {
//         for i := n - 1; i >= 0; i-- {
//             d = d<<1 | (d^v[i])&1
//         }
//         return d
//     }
//
// In a 2-space with the origin at the bottom left, this results in a ⊐ shape,
// wound counter clockwise.
//
// The runtime of Space and Curve scales as O(n∙k) where n is the dimension and
// k is the order. The length of the curve is 2^(n∙k).
//
// Limitations
//
// An n-dimensional, k-order Hilbert curve will not be fully usable if 2ⁿᵏ does
// not fit in a signed integer. Len will overflow if n∙k ≥ BitSize-1 (where
// BitSize is the size of an integer in bits). Curve may overflow if n∙k ≥
// BitSize, depending on V. Space will not overflow, but it cannot be called
// with values that do not fit in a signed integer, thus only a subset of the
// curve can be utilized.
package curve
