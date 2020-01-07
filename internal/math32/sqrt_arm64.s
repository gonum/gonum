// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine,!safe

#include "textflag.h"
// Don't insert stack check preamble.

// func Sqrt(x float32) float32
TEXT ·Sqrt(SB),NOSPLIT,$0
	FMOVD 	x+0(FP), F0
	FSQRTS F0, F0
	FMOVD F0, ret+8(FP)
	RET
