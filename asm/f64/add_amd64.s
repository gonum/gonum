// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !noasm,!appengine

#include "textflag.h"

// func Add(dst, s []float64)
TEXT ·Add(SB), NOSPLIT, $0
	MOVQ    dst_base+0(FP), DI
	MOVQ    dst_len+8(FP), CX
	MOVQ    s_base+24(FP), SI
	CMPQ    s_len+32(FP), CX
	CMOVQLE s_len+32(FP), CX
	CMPQ    CX, $0
	JE      add_end
	XORQ    AX, AX
	MOVQ    DI, BX
	ANDQ    $15, BX
	JZ      add_no_trim
	MOVSD   (DI)(AX*8), X0
	ADDSD   (SI)(AX*8), X0
	MOVSD   X0, (DI)(AX*8)
	INCQ    AX
	DECQ    CX
	JE      add_end

add_no_trim:
	MOVQ CX, BX
	ANDQ $3, BX
	SHRQ $2, CX
	JZ   add_tail_start

add_loop:
	MOVUPS (SI)(AX*8), X0
	MOVUPS 16(SI)(AX*8), X1
	ADDPD  (DI)(AX*8), X0
	MOVUPS X0, (DI)(AX*8)
	ADDPD  16(DI)(AX*8), X1
	MOVUPS X1, 16(DI)(AX*8)
	ADDQ   $4, AX
	LOOPNE add_loop
	CMPQ   BX, $0
	JE     add_end

add_tail_start:
	MOVQ BX, CX

add_tail:
	MOVSD  (DI)(AX*8), X0
	ADDSD  (SI)(AX*8), X0
	MOVSD  X0, (DI)(AX*8)
	INCQ   AX
	LOOPNE add_tail

add_end:
	RET
