// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !noasm,!appengine

#include "textflag.h"

#define SIZE 8
#define SHIFT 3
#define SRC SI
#define DST DI
#define LEN CX
#define TAIL BX
#define INC_DST R8
#define INC3_DST R9
#define INC_SRC R10
#define INC3_SRC R11

// func Copy(n int, dst []float64, incDst int, src []float64, incSrc int)
TEXT ·Copy(SB), NOSPLIT, $0
	MOVQ n+0(FP), LEN
	CMPQ LEN, $0
	JEQ  ret
	MOVQ dst_base+8(FP), DST
	MOVQ src_base+40(FP), SRC
	MOVQ incDst+32(FP), INC_DST
	SHLQ $SHIFT, INC_DST
	MOVQ incSrc+64(FP), INC_SRC
	SHLQ $SHIFT, INC_SRC

	MOVQ LEN, TAIL
	SHRQ $2, LEN
	JZ   tail

	LEAQ        (INC_SRC)(INC_SRC*2), INC3_SRC
	LEAQ        (INC_DST)(INC_DST*2), INC3_DST
	PREFETCHNTA (SRC)(INC3_SRC*1)
	PREFETCHT0  (DST)(INC3_DST*1)

copy_loop:
	MOVSD (SRC), X0
	MOVSD (SRC)(INC_SRC*1), X1
	MOVSD (SRC)(INC_SRC*2), X2
	MOVSD (SRC)(INC3_SRC*1), X3
	MOVSD X0, (DST)
	MOVSD X1, (DST)(INC_DST*1)
	MOVSD X2, (DST)(INC_DST*2)
	MOVSD X3, (DST)(INC3_DST*1)

	LEAQ        (SRC)(INC_SRC*4), SRC
	LEAQ        (DST)(INC_DST*4), DST
	PREFETCHNTA (SRC)(INC3_SRC*1)
	PREFETCHT0  (DST)(INC3_DST*1)
	DECQ        LEN
	JNZ         copy_loop

tail:
	ANDQ $3, TAIL
	JZ   ret

tail_loop:
	MOVSD (SRC), X0
	MOVSD X0, (DST)
	ADDQ  INC_SRC, SRC
	ADDQ  INC_DST, DST
	DECQ  TAIL
	JNZ   tail_loop

ret:
	RET
