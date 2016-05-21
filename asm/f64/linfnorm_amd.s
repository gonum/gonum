// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"
	
//func LinfNorm(s, t []float64) float64
TEXT ·LinfNorm(SB), NOSPLIT, $0
	MOVQ    s_base+0(FP), DI
	MOVQ    t_base+24(FP), SI
	MOVQ    s_len+8(FP), DX
	CMPQ    t_len+32(FP), DX
	CMOVQLE t_len+32(FP), DX
	PXOR    X3, X3
	XORQ	AX, AX
	CMPQ    DX, $1
	JL      l1_end
	SUBQ	$1, DX
	JE	l1_tail
l1_loop:
	MOVUPS  (SI)(AX*8), X0
	MOVUPS  (DI)(AX*8), X1
	MOVAPS  X0, X2
	SUBPD   X1, X0
	SUBPD   X2, X1
	MAXPD   X1, X0
	MAXPD   X0, X3
	ADDQ	$2, AX
	CMPQ    AX, DX
	JL	l1_loop
	JG      l1_end
l1_tail:
	PXOR    X0 ,X0
	PXOR    X1 ,X1
	MOVSD   (SI)(AX*8), X0
	MOVSD   (DI)(AX*8), X1
	MOVAPD  X0, X2
	SUBSD   X1, X0
	SUBSD   X2, X1
	MAXSD   X1, X0
	MAXSD   X0, X3
l1_end:
	MOVAPS  X3, X2
	SHUFPD  $1, X2, X2
	MAXSD   X3, X2
	MOVSD   X2, ret+48(FP)
	RET
