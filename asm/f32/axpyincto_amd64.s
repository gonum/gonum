// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

#include "textflag.h"

// func AxpyIncTo(dst []float32, incDst, idst uintptr, alpha float32, x, y []float32, n, incX, incY, ix, iy uintptr)
TEXT ·AxpyIncTo(SB), NOSPLIT, $0
	MOVQ  n+96(FP), CX
	CMPQ  CX, $0             // if n==0, return
	JE    axpyi_end
	MOVQ  dst_base+0(FP), DI
	MOVQ  x_base+48(FP), SI
	MOVQ  y_base+72(FP), DX
	MOVQ  ix+120(FP), R8     // Load the first index
	MOVQ  iy+128(FP), R9
	MOVQ  idst+32(FP), R10
	LEAQ  (SI)(R8*4), SI     // Calculate address of first index
	LEAQ  (DX)(R9*4), DX
	LEAQ  (DI)(R10*4), DI
	MOVQ  incX+104(FP), R8   // Incrementors*4 for easy iteration (ADDQ)
	SHLQ  $2, R8
	MOVQ  incY+112(FP), R9
	SHLQ  $2, R9
	MOVQ  incDst+24(FP), R10
	SHLQ  $2, R10
	MOVSS alpha+40(FP), X0
	MOVSS X0, X1
	MOVQ  CX, BX
	ANDQ  $3, BX
	SHRQ  $2, CX
	JZ    axpyi_tail_start

axpyi_loop:
	MOVSS (SI), X2
	MOVSS (SI)(R8*1), X3
	LEAQ  (SI)(R8*2), SI
	MOVSS (SI), X4
	MOVSS (SI)(R8*1), X5
	MULSS X1, X2          // (a*x)
	MULSS X0, X3
	MULSS X1, X4
	MULSS X0, X5
	ADDSS (DX), X2        // (a*x+y)
	ADDSS (DX)(R9*1), X3
	LEAQ  (DX)(R9*2), DX
	ADDSS (DX), X4
	ADDSS (DX)(R9*1), X5
	MOVSS X2, (DI)        // Write result back to dst
	MOVSS X3, (DI)(R10*1)
	LEAQ  (DI)(R10*2), DI
	MOVSS X4, (DI)
	MOVSS X5, (DI)(R10*1)
	LEAQ  (SI)(R8*2), SI  // Increment addresses
	LEAQ  (DX)(R9*2), DX
	LEAQ  (DI)(R10*2), DI
	LOOP  axpyi_loop
	CMPQ  BX, $0
	JE    axpyi_end

axpyi_tail_start:
	MOVQ BX, CX

axpyi_tail:
	MOVSS (SI), X2
	MULSS X1, X2
	ADDSS (DX), X2
	MOVSS X2, (DI)
	LEAQ  (SI)(R8*1), SI
	LEAQ  (DX)(R9*1), DX
	LEAQ  (DI)(R10*1), DI
	LOOP  axpyi_tail

axpyi_end:
	RET

