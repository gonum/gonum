// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

// TriKind represents the triangularity of the matrix.
type TriKind bool

const (
	// Upper specifies an upper triangular matrix.
	Upper TriKind = true
	// Lower specifies a lower triangular matrix.
	Lower TriKind = false
)

// GSVDKind specifies the treatment of singular vectors during a GSVD
// factorization.
type GSVDKind int

const (
	// GSVDU specifies that the U singular vectors should be computed during
	// the decomposition.
	GSVDU GSVDKind = 1 << iota
	// GSVDV specifies that the V singular vectors should be computed during
	// the decomposition.
	GSVDV
	// GSVDQ specifies that the Q singular vectors should be computed during
	// the decomposition.
	GSVDQ

	// GSVDNone specifies that no singular vector should be computed during
	// the decomposition.
	GSVDNone
)
