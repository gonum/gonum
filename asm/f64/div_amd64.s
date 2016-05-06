// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !noasm,!appengine

#include "textflag.h"

// func Daddconst(dst, s []float64)
TEXT ·Div(SB), NOSPLIT, $0
	MOVQ    dst_base+0(FP), DI
	MOVQ    dst_len+8(FP), DX
	MOVQ    s_base+24(FP), SI
	CMPQ    s_len+32(FP), DX
	CMOVQLE s_len+32(FP), DX
	CMPQ    DX, $0
	JE      div_end
	XORQ    AX, AX
	CMPQ    DX, $4
	JL      div_tail

div_loop: // Unroll 4x
	MOVUPS (DI)(AX*8), X0
	DIVPD  (SI)(AX*8), X0
	MOVUPS X0, (DI)(AX*8)
	MOVUPS 16(DI)(AX*8), X0
	DIVPD  16(SI)(AX*8), X0
	MOVUPS X0, 16(DI)(AX*8)
	ADDQ   $4, AX
	SUBQ   $4, DX
	CMPQ   DX, $4
	JGE    div_loop
	CMPQ   DX, $0
	JE     div_end

div_tail:
	MOVSD (DI)(AX*8), X0
	DIVSD (SI)(AX*8), X0
	MOVSD X0, (DI)(AX*8)
	INCQ  AX
	DECQ  DX
	JNZ   div_tail

div_end:
	RET

