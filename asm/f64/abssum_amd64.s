// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !noasm,!appengine

#include "textflag.h"

// func AbsSum(x []float64) float64
TEXT ·AbsSum(SB), NOSPLIT, $0
	MOVQ x_base+0(FP), SI
	MOVQ x_len+8(FP), CX
	XORQ AX, AX
	PXOR X0, X0
	PXOR X1, X1
	PXOR X2, X2
	PXOR X3, X3
	PXOR X4, X4
	PXOR X5, X5
	PXOR X6, X6
	PXOR X7, X7
	CMPQ CX, $0
	JE   absum_end
	MOVQ CX, BX
	ANDQ $7, BX
	SHRQ $3, CX
	JZ   absum_tail_start

absum_loop:
	MOVUPS (SI)(AX*8), X8
	MOVUPS 16(SI)(AX*8), X9
	MOVUPS 32(SI)(AX*8), X10
	MOVUPS 48(SI)(AX*8), X11
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
	ADDQ   $8, AX
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
	MOVSD (SI)(AX*8), X8
	ADDSD X8, X0
	SUBSD X8, X1
	MAXSD X1, X0
	MOVSD X0, X1
	INCQ  AX
	LOOP  absum_tail

absum_end:
	MOVSD X1, sum+24(FP)
	RET
