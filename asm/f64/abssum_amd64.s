// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !noasm,!appengine

#include "textflag.h"

// func AbsSum(x []float64) float64
TEXT ·AbsSum(SB), NOSPLIT, $0
	MOVQ x_base+0(FP), SI
	MOVQ x_len+8(FP), DX
	XORQ AX, AX
	PXOR X0, X0
	PXOR X1, X1
	CMPQ DX, $1
	JL   absum_end
	JE   absum_tail

absum_loop:
	MOVUPS (SI)(AX*8), X2
	ADDPD  X2, X0
	ADDPD  X1, X2
	MAXPD  X2, X0
	MOVAPS X0, X1
	ADDQ   $2, AX
	CMPQ   AX, DX
	JL     absum_loop
	JE     absum_end

absum_tail:
	MOVSD (SI)(AX*8), X2
	ADDSD X2, X0
	ADDSD X1, X2
	MAXSD X2, X0

absum_end:
	MOVAPS X0, X1
	SHUFPD $0, X0, X0
	MAXSD  X0, X1
	MOVSD  X1, sum+24(FP)
	RET
