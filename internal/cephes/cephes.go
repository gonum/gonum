// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package cephes implements functions originally in the Netlib code by Stephen Mosher
package cephes

/*
Additional copyright information:

Code in this package is adapted from the Cephes library (http://www.netlib.org/cephes/).
There is no explicit licence on Netlib, but the author has agreed to a BSD release.
See https://github.com/deepmind/torch-cephes/blob/master/LICENSE.txt and
https://lists.debian.org/debian-legal/2004/12/msg00295.html
*/

var (
	badParamOutOfBounds = "cephes: parameter out of bounds"
)

const (
	machEp = 1.11022302462515654042e-16 // 2^-53
	maxLog = 7.09782712893383996732e2   // log(2^127)
	minLog = -7.451332191019412076235e2 // log(2^-128)
)
