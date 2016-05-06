// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !noasm,!appengine

#include "textflag.h"

// func AbsSumInc(x []float64, n, incX int) (sum float64)
TEXT ·AbsSumInc(SB), NOSPLIT, $0
	MOVQ x_base+0(FP), SI
	MOVQ n+24(FP), DX
	MOVQ incX+32(FP), CX
	SHLQ $3, CX
	XORQ AX, AX
	PXOR X0, X0
	PXOR X1, X1
	PXOR X4, X4
	CMPQ DX, $0
	JE   absum_end

absum_loop:
	MOVSD (SI), X2
	ADDSD X2, X0
	ADDSD X1, X2
	MOVSD (SI)(CX*1), X3
	ADDSD X3, X4
	ADDSD X5, X3
	MAXSD X2, X0
	MAXSD X3, X4
	MOVSD X0, X1
	MOVSD X4, X5
	ADDQ  CX, SI
	ADDQ  AX, SI
	CMPQ  AX, DX
	JL    absum_loop
	ADDSD X4, X0

absum_end:
	MOVSD  X1, X0
	SHUFPD $0, X0, X0
	MAXSD  X0, X1
	MOVSD  X1, sum+40(FP)
	RET
