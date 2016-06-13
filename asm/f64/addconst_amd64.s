// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !noasm,!appengine

#include "textflag.h"

// func Addconst(alpha float64, x []float64)
TEXT ·AddConst(SB), NOSPLIT, $0
	MOVQ   x_base+8(FP), SI
	MOVQ   x_len+16(FP), CX
	CMPQ   CX, $0
	JE     ac_end
	MOVSD  alpha+0(FP), X4
	SHUFPD $0, X4, X4
	MOVUPS X4, X5
	XORQ   AX, AX
	MOVQ   CX, BX
	ANDQ   $7, BX
	SHRQ   $3, CX
	JZ     ac_tail_start

ac_loop:
	MOVUPS (SI)(AX*8), X0
	MOVUPS 16(SI)(AX*8), X1
	MOVUPS 32(SI)(AX*8), X2
	MOVUPS 48(SI)(AX*8), X3
	ADDPD  X4, X0
	ADDPD  X5, X1
	ADDPD  X4, X2
	ADDPD  X5, X3
	MOVUPS X0, (SI)(AX*8)
	MOVUPS X1, 16(SI)(AX*8)
	MOVUPS X2, 32(SI)(AX*8)
	MOVUPS X3, 48(SI)(AX*8)
	ADDQ   $8, AX
	LOOP   ac_loop
	CMPQ   BX, $0
	JE     ac_end

ac_tail_start:
	MOVQ BX, CX

ac_tail:
	MOVSD (SI)(AX*8), X0
	ADDSD X4, X0
	MOVSD X0, (SI)(AX*8)
	INCQ  AX
	LOOP  ac_tail

ac_end:
	RET
