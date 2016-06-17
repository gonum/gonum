// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

#include "textflag.h"

// func AxpyIncTo(dst []float32, incDst, idst uintptr, alpha float32, x, y []float32, n, incX, incY, ix, iy uintptr)
TEXT ·AxpyIncTo(SB), NOSPLIT, $0
	MOVQ  n+96(FP), CX       // CX := n
	CMPQ  CX, $0             // if n==0 { return }
	JLE   axpyi_end
	MOVQ  dst_base+0(FP), DI // DI := &dst
	MOVQ  x_base+48(FP), SI  // SI := &x
	MOVQ  y_base+72(FP), DX  // DX := &y
	MOVQ  ix+120(FP), R8     // Load the first index
	MOVQ  iy+128(FP), R9
	MOVQ  idst+32(FP), R10
	LEAQ  (SI)(R8*4), SI     // SI = &(x[ix])
	LEAQ  (DX)(R9*4), DX
	LEAQ  (DI)(R10*4), DI
	MOVQ  incX+104(FP), R8   // Incrementors*4 for easy iteration (ADDQ)
	SHLQ  $2, R8
	MOVQ  incY+112(FP), R9
	SHLQ  $2, R9
	MOVQ  incDst+24(FP), R10
	SHLQ  $2, R10
	MOVSS alpha+40(FP), X0   // X0 = alpha
	MOVSS X0, X1
	MOVQ  CX, BX
	ANDQ  $3, BX             // BX = CX % 4
	SHRQ  $2, CX             // CX = floor( CX / 4 )
	JZ    axpyi_tail_start   // if CX == 0 { goto axpyi_tail_start }

axpyi_loop: // Loop unrolled 4x   do {
	MOVSS (SI), X2        // X_i = x[i]
	MOVSS (SI)(R8*1), X3
	LEAQ  (SI)(R8*2), SI  // SI = &(x[i+2)
	MOVSS (SI), X4
	MOVSS (SI)(R8*1), X5
	MULSS X1, X2          // X_i *= a
	MULSS X0, X3
	MULSS X1, X4
	MULSS X0, X5
	ADDSS (DX), X2        // X_i += y[i]
	ADDSS (DX)(R9*1), X3
	LEAQ  (DX)(R9*2), DX
	ADDSS (DX), X4
	ADDSS (DX)(R9*1), X5
	MOVSS X2, (DI)        // dst[i] = X_i
	MOVSS X3, (DI)(R10*1)
	LEAQ  (DI)(R10*2), DI
	MOVSS X4, (DI)
	MOVSS X5, (DI)(R10*1)
	LEAQ  (SI)(R8*2), SI  // Increment addresses
	LEAQ  (DX)(R9*2), DX
	LEAQ  (DI)(R10*2), DI
	LOOP  axpyi_loop      // } while --CX > 0
	CMPQ  BX, $0          // if BX == 0 { return }
	JE    axpyi_end

axpyi_tail_start: // Reset loop registers
	MOVQ BX, CX // Loop counter: CX = BX

axpyi_tail: // do {
	MOVSS (SI), X2 // X2 = x[i]
	MULSS X1, X2   // X2 *= a
	ADDSS (DX), X2 // X2 += y[i]
	MOVSS X2, (DI) // dst[i] = X2
	ADDQ  R8, SI
	ADDQ  R9, DX
	ADDQ  R10, DI
/* LEAQ  (SI)(R8*1), SI  // SI = &(x[incX])
	LEAQ  (DX)(R9*1), DX  // DX = &(y[incY])
	LEAQ  (DI)(R10*1), DI // DI = &(dst[incDst]) */
	LOOP axpyi_tail // } while --CX > 0

axpyi_end:
	RET

