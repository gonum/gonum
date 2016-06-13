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
	ANDQ    $0x0F, BX
	JZ      add_no_trim

	// Align on 16-bit boundary
	MOVSD (DI)(AX*8), X0
	ADDSD (SI)(AX*8), X0
	MOVSD X0, (DI)(AX*8)
	INCQ  AX
	DECQ  CX
	JE    add_end

add_no_trim:
	MOVQ CX, BX
	ANDQ $7, BX
	SHRQ $3, CX
	JZ   add_tail_start

add_loop: // Loop unrolled 8x
	MOVUPS (SI)(AX*8), X0
	MOVUPS 16(SI)(AX*8), X1
	MOVUPS 32(SI)(AX*8), X2
	MOVUPS 48(SI)(AX*8), X3
	ADDPD  (DI)(AX*8), X0
	ADDPD  16(DI)(AX*8), X1
	ADDPD  32(DI)(AX*8), X2
	ADDPD  48(DI)(AX*8), X3
	MOVUPS X0, (DI)(AX*8)
	MOVUPS X1, 16(DI)(AX*8)
	MOVUPS X2, 32(DI)(AX*8)
	MOVUPS X3, 48(DI)(AX*8)
	ADDQ   $8, AX
	LOOP   add_loop
	CMPQ   BX, $0
	JE     add_end

add_tail_start:
	MOVQ BX, CX

add_tail:
	MOVSD (DI)(AX*8), X0
	ADDSD (SI)(AX*8), X0
	MOVSD X0, (DI)(AX*8)
	INCQ  AX
	LOOP  add_tail

add_end:
	RET
