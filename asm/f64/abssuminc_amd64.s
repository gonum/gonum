// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !noasm,!appengine

#include "textflag.h"

// func AbsSumInc(x []float64, n, incX int) (sum float64)
TEXT ·AbsSumInc(SB), NOSPLIT, $0
	MOVQ  x_base+0(FP), SI
	MOVQ  n+24(FP), CX
	MOVQ  incX+32(FP), AX
	SHLQ  $3, AX
	MOVQ  AX, DX
	IMULQ $3, DX
	PXOR  X0, X0
	PXOR  X1, X1
	PXOR  X2, X2
	PXOR  X3, X3
	PXOR  X4, X4
	PXOR  X5, X5
	PXOR  X6, X6
	PXOR  X7, X7
	CMPQ  CX, $0
	JE    absum_end
	MOVQ  CX, BX
	ANDQ  $7, BX
	SHRQ  $3, CX
	JZ    absum_tail_start

absum_loop:
	MOVSD  (SI), X8
	MOVSD  (SI)(AX*1), X9
	MOVSD  (SI)(AX*2), X10
	MOVSD  (SI)(DX*1), X11
	LEAQ   (SI)(AX*4), SI
	MOVHPD (SI), X8
	MOVHPD (SI)(AX*1), X9
	MOVHPD (SI)(AX*2), X10
	MOVHPD (SI)(DX*1), X11
	ADDPD  X8, X0
	ADDPD  X9, X2
	ADDPD  X10, X4
	ADDPD  X11, X6
	SUBPD  X8, X1
	SUBPD  X9, X3
	SUBPD  X10, X5
	SUBPD  X11, X7
	MAXPD  X1, X0
	MAXPD  X3, X2
	MAXPD  X5, X4
	MAXPD  X7, X6
	MOVAPS X0, X1
	MOVAPS X2, X3
	MOVAPS X4, X5
	MOVAPS X6, X7
	LEAQ   (SI)(AX*4), SI
	LOOP   absum_loop
	ADDPD  X3, X0
	ADDPD  X5, X7
	ADDPD  X7, X0
	MOVAPS X0, X1
	SHUFPD $0x3, X0, X0
	ADDSD  X1, X0
	MOVSD  X0, X1
	CMPQ   BX, $0
	JE     absum_end

absum_tail_start:
	MOVQ  BX, CX
	XORPS X8, X8

absum_tail:
	MOVSD (SI), X8
	ADDSD X8, X0
	SUBSD X8, X1
	MAXSD X1, X0
	MOVSD X0, X1
	ADDQ  AX, SI
	LOOP  absum_tail

absum_end:
	MOVSD X1, sum+40(FP)
	RET
