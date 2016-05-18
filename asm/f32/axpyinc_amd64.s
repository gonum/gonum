// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

#include "textflag.h"

// func AxpyInc(alpha float32, x, y []float32, n, incX, incY, ix, iy uintptr)
TEXT ·AxpyInc(SB), NOSPLIT, $0
	MOVQ    n+56(FP), CX
	CMPQ    CX, $0
	JLE     saxyi_end
	MOVQ    x_base+8(FP), SI
	MOVQ    y_base+32(FP), DI
	MOVQ    ix+80(FP), AX
	MOVQ    iy+88(FP), BX
	LEAQ    (SI)(AX*4), SI
	LEAQ    (DI)(BX*4), DI
	MOVQ    incX+64(FP), AX
	MOVQ    incY+72(FP), BX
	IMULQ   $4, AX
	IMULQ   $4, BX
	MOVSS   alpha+0(FP), X0
	MOVSS   X0, X2
	XORQ    R9, R9
	SHRQ    $1, CX
	CMOVQCS AX, R9
	JZ      saxyi_odd

saxyi_loop:
	MOVSS  (SI), X1
	MOVSS  (SI)(AX*1), X3
	MULSS  X0, X1
	MULSS  X2, X3
	ADDSS  (DI), X1
	ADDSS  (DI)(BX*1), X3
	MOVSS  X1, (DI)
	MOVSS  X3, (DI)(AX*1)
	LEAQ   (SI)(AX*2), SI
	LEAQ   (DI)(BX*2), DI
	LOOPNE saxyi_loop
	CMPQ   R9, $0
	JE     saxyi_end

saxyi_odd:
	// Trim odd n
	MOVSS (SI), X1
	MULSS X0, X1
	ADDSS (DI), X1
	MOVSS X1, (DI)
	ADDQ  AX, SI
	ADDQ  BX, DI

saxyi_end:
	RET
