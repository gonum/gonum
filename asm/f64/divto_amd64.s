// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !noasm,!appengine

#include "textflag.h"

// func DivTo(dst, x, y []float64)
TEXT ·DivTo(SB), NOSPLIT, $0
	MOVQ    dst_base+0(FP), DI
	MOVQ    dst_len+8(FP), CX
	MOVQ    x_base+24(FP), SI
	MOVQ    y_base+48(FP), DX
	CMPQ    x_len+32(FP), CX
	CMOVQLE x_len+32(FP), CX
	CMPQ    y_len+56(FP), CX
	CMOVQLE y_len+56(FP), CX
	MOVQ    CX, ret_len+80(FP)
	CMPQ    CX, $0
	JE      div_end
	XORQ    AX, AX
	MOVQ    DI, BX
	ANDQ    $15, BX
	JZ      div_no_trim

	// Align on 16-bit boundary
	MOVSD (SI)(AX*8), X0
	DIVSD (DX)(AX*8), X0
	MOVSD X0, (DI)(AX*8)
	INCQ  AX
	DECQ  CX
	JZ    div_end        // */

div_no_trim:
	MOVQ CX, BX
	ANDQ $7, BX
	SHRQ $3, CX
	JZ   div_tail_start

div_loop: // Unroll 8x
	MOVUPS (SI)(AX*8), X0
	MOVUPS 16(SI)(AX*8), X1
	MOVUPS 32(SI)(AX*8), X2
	MOVUPS 48(SI)(AX*8), X3
	DIVPD  (DX)(AX*8), X0
	DIVPD  16(DX)(AX*8), X1
	DIVPD  32(DX)(AX*8), X2
	DIVPD  48(DX)(AX*8), X3
	MOVUPS X0, (DI)(AX*8)
	MOVUPS X1, 16(DI)(AX*8)
	MOVUPS X2, 32(DI)(AX*8)
	MOVUPS X3, 48(DI)(AX*8)
	ADDQ   $8, AX
	LOOP   div_loop
	CMPQ   CX, $0
	JE     div_end

div_tail_start:
	MOVQ BX, CX

div_tail:
	MOVSD (SI)(AX*8), X0
	DIVSD (DX)(AX*8), X0
	MOVSD X0, (DI)(AX*8)
	INCQ  AX
	LOOP  div_tail

div_end:
	MOVQ DI, ret_base+72(FP)
	MOVQ dst_cap+16(FP), DI
	MOVQ DI, ret_cap+88(FP)
	RET
