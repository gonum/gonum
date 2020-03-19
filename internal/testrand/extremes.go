// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testrand

import "math"

const (
	maxUint = ^uint(0)
	maxInt  = int(maxUint >> 1)
)

var (
	extremeFloat64Unit = [...]float64{
		0,
		math.SmallestNonzeroFloat64,
		0.5,
		1 - math.SmallestNonzeroFloat64,
		1,
	}

	extremeFloat64Norm = [...]float64{
		-math.MaxFloat64,
		-math.MaxFloat64 / 2,
		-1,
		-math.SmallestNonzeroFloat64,
		0,
		math.SmallestNonzeroFloat64,
		1,
		math.MaxFloat64 / 2,
		math.MaxFloat64,
	}

	extremeFloat64Exp = [...]float64{
		0,
		math.SmallestNonzeroFloat64,
		1,
		math.MaxFloat64 / 2,
		math.MaxFloat64,
	}

	extremeFloat32Unit = [...]float32{
		0,
		math.SmallestNonzeroFloat32,
		0.5,
		1 - math.SmallestNonzeroFloat32,
		1,
	}

	extremeInt = [...]int{
		0,
		1,
		maxInt / 2,
		maxInt - 1,
		maxInt,
	}

	extremeInt31 = [...]int32{
		0,
		1,
		math.MaxInt32 / 2,
		math.MaxInt32 - 1,
		math.MaxInt32,
	}

	extremeInt63 = [...]int64{
		0,
		1,
		math.MaxInt64 / 2,
		math.MaxInt64 - 1,
		math.MaxInt64,
	}

	extremeUint32 = [...]uint32{
		0,
		1,
		math.MaxUint32 / 2,
		math.MaxUint32 - 1,
		math.MaxUint32,
	}

	extremeUint64 = [...]uint64{
		0,
		1,
		math.MaxUint64 / 2,
		math.MaxUint64 - 1,
		math.MaxUint64,
	}
)
