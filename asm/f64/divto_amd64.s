// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !noasm,!appengine

#include "textflag.h"

// func DivTo(dst, x, y []float64)
TEXT ·DivTo(SB), NOSPLIT, $0
	MOVQ    dst_base+0(FP), DI
	MOVQ    dst_len+8(FP), DX
	MOVQ    x_base+24(FP), SI
	MOVQ    y_base+48(FP), BX
	CMPQ    x_len+32(FP), DX
	CMOVQLE x_len+32(FP), DX
	CMPQ    y_len+56(FP), DX
	CMOVQLE y_len+56(FP), DX
	MOVQ    DX, ret_len+80(FP)
	CMPQ    DX, $0
	JE      div_end
	XORQ    AX, AX
	CMPQ    DX, $4
	JL      div_tail

div_loop: // Unroll 4x
	MOVUPS (SI)(AX*8), X0
	DIVPD  (BX)(AX*8), X0
	MOVUPS X0, (DI)(AX*8)
	MOVUPS 16(SI)(AX*8), X0
	DIVPD  16(BX)(AX*8), X0
	MOVUPS X0, 16(DI)(AX*8)
	ADDQ   $4, AX
	SUBQ   $4, DX
	CMPQ   DX, $4
	JGE    div_loop
	CMPQ   DX, $0
	JE     div_end

div_tail:
	MOVSD (SI)(AX*8), X0
	DIVSD (BX)(AX*8), X0
	MOVSD X0, (DI)(AX*8)
	INCQ  AX
	DECQ  DX
	JNZ   div_tail

div_end:
	MOVQ DI, ret_base+72(FP)
	MOVQ dst_cap+16(FP), DI
	MOVQ DI, ret_cap+88(FP)
	RET
