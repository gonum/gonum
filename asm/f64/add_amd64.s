// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !noasm,!appengine

#include "textflag.h"

// func Add(dst, s []float64)
TEXT ·Add(SB), NOSPLIT, $0
	MOVQ    dst_base+0(FP), DI
	MOVQ    dst_len+8(FP), DX
	MOVQ    s_base+24(FP), SI
	CMPQ    s_len+32(FP), DX
	CMOVQLE s_len+32(FP), DX
	CMPQ    DX, $0
	JE      add_end
	XORQ    AX, AX
	CMPQ    DX, $4
	JL      add_tail

add_loop:
	MOVUPS (SI)(AX*8), X0
	MOVUPS 16(SI)(AX*8), X1
	ADDPD  (DI)(AX*8), X0
	MOVUPS X0, (DI)(AX*8)
	ADDPD  16(DI)(AX*8), X1
	MOVUPS X1, 16(DI)(AX*8)
	ADDQ   $4, AX
	SUBQ   $4, DX
	CMPQ   DX, $4
	JGE    add_loop
	CMPQ   DX, $0
	JE     add_end

add_tail:
	MOVSD (DI)(AX*8), X0
	ADDSD (SI)(AX*8), X0
	MOVSD X0, (DI)(AX*8)
	INCQ  AX
	DECQ  DX
	JNZ   add_tail

add_end:
	RET
