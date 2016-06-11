// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !noasm,!appengine

#include "textflag.h"

// func Div(dst, s []float64)
TEXT ·Div(SB), NOSPLIT, $0
	MOVQ    dst_base+0(FP), DI
	MOVQ    dst_len+8(FP), CX
	MOVQ    s_base+24(FP), SI
	CMPQ    s_len+32(FP), CX
	CMOVQLE s_len+32(FP), CX
	CMPQ    CX, $0
	JE      div_end
	XORQ    AX, AX
	MOVQ    SI, BX
	ANDQ    $15, BX
	JZ      div_no_trim

	// Align on 16-bit boundary
	MOVSD (DI)(AX*8), X0
	DIVSD (SI)(AX*8), X0
	MOVSD X0, (DI)(AX*8)
	INCQ  AX
	DECQ  CX
	JZ    div_end

div_no_trim:
	MOVQ CX, BX
	ANDQ $7, BX
	SHRQ $3, CX
	JZ   div_tail_start

div_loop: // Loop unrolled 8x
	MOVUPS (SI)(AX*8), X0
	MOVUPS 16(SI)(AX*8), X1
	MOVUPS 32(SI)(AX*8), X2
	MOVUPS 48(SI)(AX*8), X3
	DIVPD  (DI)(AX*8), X0
	DIVPD  16(DI)(AX*8), X1
	DIVPD  32(DI)(AX*8), X2
	DIVPD  48(DI)(AX*8), X3
	MOVUPS X0, (DI)(AX*8)
	MOVUPS X1, 16(DI)(AX*8)
	MOVUPS X2, 32(DI)(AX*8)
	MOVUPS X3, 48(DI)(AX*8)
	ADDQ   $4, AX
	LOOP   div_loop
	CMPQ   BX, $0
	JE     div_end

div_tail_start:
	MOVQ BX, CX

div_tail:
	MOVSD (DI)(AX*8), X0
	DIVSD (SI)(AX*8), X0
	MOVSD X0, (DI)(AX*8)
	INCQ  AX
	LOOP  div_tail

div_end:
	RET

