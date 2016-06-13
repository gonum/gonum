// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !noasm,!appengine

#include "textflag.h"

TEXT ·CumSum(SB), NOSPLIT, $0
	MOVQ    dst_base+0(FP), DI
	MOVQ    dst_len+8(FP), CX
	MOVQ    s_base+24(FP), SI
	CMPQ    s_len+32(FP), CX
	CMOVQLE s_len+32(FP), CX
	MOVQ    CX, ret_len+56(FP)
	CMPQ    CX, $0
	JE      cs_end
	XORQ    AX, AX
	PXOR    X5, X5
	MOVQ    CX, BX
	ANDQ    $3, BX
	SHRQ    $2, CX
	JZ      cs_tail_start

cs_loop:
	MOVUPS (SI)(AX*8), X0
	MOVUPS 16(SI)(AX*8), X2
	MOVAPS X0, X1
	MOVAPS X2, X3
	SHUFPD $1, X1, X1
	SHUFPD $1, X3, X3
	ADDPD  X0, X1
	ADDPD  X2, X3
	SHUFPD $2, X1, X0
	SHUFPD $3, X1, X1
	SHUFPD $2, X3, X2
	SHUFPD $3, X3, X3
	ADDPD  X5, X0
	ADDPD  X1, X5
	ADDPD  X5, X2
	MOVUPS X0, (DI)(AX*8)
	MOVUPS X2, 16(DI)(AX*8)
	ADDPD  X3, X5
	ADDQ   $4, AX
	LOOP   cs_loop
	CMPQ   BX, $0
	JE     cs_end

cs_tail_start:
	MOVQ BX, CX

cs_tail:
	ADDSD (SI)(AX*8), X5
	MOVSD X5, (DI)(AX*8)
	INCQ  AX
	LOOP  cs_tail

cs_end:
	MOVQ DI, ret_base+48(FP)
	MOVQ dst_cap+16(FP), SI
	MOVQ SI, ret_cap+64(FP)
	RET
