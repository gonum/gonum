// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !noasm,!appengine

#include "textflag.h"

TEXT ·CumProd(SB), NOSPLIT, $0
	MOVQ    dst_base+0(FP), DI
	MOVQ    dst_len+8(FP), DX
	MOVQ    s_base+24(FP), SI
	CMPQ    s_len+32(FP), DX
	CMOVQLE s_len+32(FP), DX
	MOVQ    DX, ret_len+56(FP)
	CMPQ    DX, $0
	JE      cs_end
	MOVSD   (SI), X2
	MOVSD   X2, (DI)
	SHUFPD  $0, X2, X2
	MOVQ    $1, AX
	SUBQ    $1, DX
	JZ      cs_end
	SUBQ    $2, DX
	JL      cs_tail

cs_loop:
	PXOR   X1, X1
	MOVUPS (SI)(AX*8), X0
	MOVAPS X0, X1
	SHUFPD $1, X1, X1
	MULPD  X0, X1
	SHUFPD $2, X1, X0
	SHUFPD $3, X1, X1
	MULPD  X2, X0
	MOVUPS X0, (DI)(AX*8)
	MULPD  X1, X2
	ADDQ   $2, AX
	SUBQ   $2, DX
	JGE    cs_loop
	ADDQ   $2, DX
	JZ     cs_end

cs_tail:
	MULSD (SI)(AX*8), X2
	MOVSD X2, (DI)(AX*8)

cs_end:
	MOVQ DI, ret_base+48(FP)
	MOVQ dst_cap+16(FP), DI
	MOVQ DI, ret_cap+64(FP)
	RET
