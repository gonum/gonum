// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !noasm,!gccgo,!safe

#include "textflag.h"

// func L1Norm(x []float64) float64
TEXT ·L1Norm(SB), NOSPLIT, $0
	MOVQ x_base+0(FP), SI
	MOVQ x_len+8(FP), CX
	XORQ AX, AX
	PXOR X0, X0
	PXOR X2, X2
	PXOR X4, X4
	PXOR X6, X6
	CMPQ CX, $0
	JE   absum_end
	// Clear each input sign before addition. Computing max(sum+x, sum-x)
	// instead can select NaN from Inf-Inf when sum and x are both +Inf.
	PCMPEQL X12, X12
	PSRLQ   $1, X12 // two copies of 0x7fffffffffffffff
	MOVQ CX, BX
	ANDQ $7, BX
	SHRQ $3, CX
	JZ   absum_tail_start

absum_loop:
	MOVUPS (SI)(AX*8), X8
	MOVUPS 16(SI)(AX*8), X9
	MOVUPS 32(SI)(AX*8), X10
	MOVUPS 48(SI)(AX*8), X11
	ANDPD X12, X8
	ANDPD X12, X9
	ANDPD X12, X10
	ANDPD X12, X11
	ADDPD X8, X0
	ADDPD X9, X2
	ADDPD X10, X4
	ADDPD X11, X6
	ADDQ $8, AX
	LOOP absum_loop

	// Preserve the original lane and horizontal addition order.
	ADDPD X2, X0
	ADDPD X4, X6
	ADDPD X6, X0
	MOVAPS X0, X1
	SHUFPD $0x3, X0, X0
	ADDSD X1, X0
	CMPQ BX, $0
	JE absum_end

absum_tail_start:
	MOVQ BX, CX

absum_tail:
	MOVSD (SI)(AX*8), X8
	ANDPD X12, X8
	ADDSD X8, X0
	INCQ AX
	LOOP absum_tail

absum_end:
	MOVSD X0, sum+24(FP)
	RET
