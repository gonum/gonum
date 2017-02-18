// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"
	"math/rand"
)

// Dlatm1 computes the entries of dst as specified by mode, cond and rsign.
//
// mode describes how dst will be computed:
//  |mode| == 1: dst[0] = 1 and dst[1:n] = 1/cond
//  |mode| == 2: dst[:n-1] = 1/cond and dst[n-1] = 1
//  |mode| == 3: dst[i] = cond^{-i/(n-1)}, i=0,...,n-1
//  |mode| == 4: dst[i] = 1 - i*(1-1/cond)/(n-1)
//  |mode| == 5: dst[i] = random number in the range (1/cond, 1) such that
//                    their logarithms are uniformly distributed
//  |mode| == 6: dst[i] = random number from the distribution given by dist
// If mode is negative, the order of the elements of dst will be reversed.
// For other values of mode Dlatm1 will panic.
//
// If rsign is true and mode is not ±6, each entry of dst will be multiplied by 1
// or -1 with probability 0.5
//
// dist specifies the type of distribution to be used when mode == ±6:
//  dist == 1: Uniform[0,1)
//  dist == 2: Uniform[-1,1)
//  dist == 3: Normal(0,1)
// For other values of dist Dlatm1 will panic.
//
// rnd is used as a source of random numbers.
func Dlatm1(dst []float64, mode int, cond float64, rsign bool, dist int, rnd *rand.Rand) {
	amode := mode
	if amode < 0 {
		amode = -amode
	}
	if amode < 1 || 6 < amode {
		panic("testlapack: invalid mode")
	}
	if cond < 1 {
		panic("testlapack: cond < 1")
	}
	if amode == 6 && (dist < 1 || 3 < dist) {
		panic("testlapack: invalid dist")
	}

	n := len(dst)
	if n == 0 {
		return
	}

	switch amode {
	case 1:
		dst[0] = 1
		for i := 1; i < n; i++ {
			dst[i] = 1 / cond
		}
	case 2:
		for i := 0; i < n-1; i++ {
			dst[i] = 1
		}
		dst[n-1] = 1 / cond
	case 3:
		dst[0] = 1
		if n > 1 {
			alpha := math.Pow(cond, -1/float64(n-1))
			for i := 1; i < n; i++ {
				dst[i] = math.Pow(alpha, float64(i))
			}
		}
	case 4:
		dst[0] = 1
		if n > 1 {
			condInv := 1 / cond
			alpha := (1 - condInv) / float64(n-1)
			for i := 1; i < n; i++ {
				dst[i] = float64(n-i-1)*alpha + condInv
			}
		}
	case 5:
		alpha := math.Log(1 / cond)
		for i := range dst {
			dst[i] = math.Exp(alpha * rnd.Float64())
		}
	case 6:
		switch dist {
		case 1:
			for i := range dst {
				dst[i] = rnd.Float64()
			}
		case 2:
			for i := range dst {
				dst[i] = 2*rnd.Float64() - 1
			}
		case 3:
			for i := range dst {
				dst[i] = rnd.NormFloat64()
			}
		}
	}

	if rsign && amode != 6 {
		for i, v := range dst {
			if rnd.Float64() < 0.5 {
				dst[i] = -v
			}
		}
	}

	if mode < 0 {
		for i := 0; i < n/2; i++ {
			dst[i], dst[n-i-1] = dst[n-i-1], dst[i]
		}
	}
}
