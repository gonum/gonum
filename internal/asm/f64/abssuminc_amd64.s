// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !noasm,!gccgo,!safe

#include "textflag.h"

// func L1NormInc(x []float64, n, incX int) (sum float64)
TEXT ·L1NormInc(SB), NOSPLIT, $0
	MOVQ  x_base+0(FP), SI // SI = &x
	MOVQ  n+24(FP), CX     // CX = n
	MOVQ  incX+32(FP), AX  // AX =  increment * sizeof( float64 )
	// Zero stride follows the documented empty loop, without reading x.
	TESTQ AX, AX
	JE    absum_zero_stride
	SHLQ  $3, AX
	MOVQ  AX, DX           // DX = AX * 3
	IMULQ $3, DX
	PXOR  X0, X0           // p_sum_i = 0
	PXOR  X2, X2
	PXOR  X4, X4
	PXOR  X6, X6
	CMPQ  CX, $0           // if CX == 0 { return 0 }
	JE    absum_end
	// Clear the input sign before adding. max(sum+x, sum-x) can select
	// NaN from Inf-Inf when an accumulated +Inf meets another +Inf.
	PCMPEQL X12, X12
	PSRLQ   $1, X12 // two copies of 0x7fffffffffffffff
	MOVQ  CX, BX
	ANDQ  $7, BX           // BX = n % 8
	SHRQ  $3, CX           // CX = floor( n / 8 )
	JZ    absum_tail_start // if CX == 0 { goto absum_tail_start }

absum_loop: // do {
	// Preserve the Inc lane packing: (0,4), (1,5), (2,6), (3,7).
	MOVSD  (SI), X8        // X_i[0] = x[i]
	MOVSD  (SI)(AX*1), X9
	MOVSD  (SI)(AX*2), X10
	MOVSD  (SI)(DX*1), X11
	LEAQ   (SI)(AX*4), SI  // SI = SI + 4
	MOVHPD (SI), X8        // X_i[1] = x[i+4]
	MOVHPD (SI)(AX*1), X9
	MOVHPD (SI)(AX*2), X10
	MOVHPD (SI)(DX*1), X11
	ANDPD  X12, X8
	ANDPD  X12, X9
	ANDPD  X12, X10
	ANDPD  X12, X11
	ADDPD  X8, X0          // p_sum_i += abs(X_i)
	ADDPD  X9, X2
	ADDPD  X10, X4
	ADDPD  X11, X6
	LEAQ   (SI)(AX*4), SI  // SI = SI + 4
	LOOP   absum_loop      // } while --CX > 0

	// Preserve the original per-lane merge: (lane0+lane1)+(lane3+lane2).
	ADDPD X2, X0
	ADDPD X4, X6
	ADDPD X6, X0

	// p_sum_0[0] = p_sum_0[0] + p_sum_0[1]
	MOVAPS X0, X1
	SHUFPD $0x3, X0, X0 // lower( p_sum_0 ) = upper( p_sum_0 )
	ADDSD  X1, X0
	CMPQ   BX, $0
	JE     absum_end    // if BX == 0 { goto absum_end }

absum_tail_start: // Reset loop registers
	MOVQ  BX, CX // Loop counter:  CX = BX

absum_tail: // do {
	// Add one absolute value in the original scalar-tail order.
	MOVSD (SI), X8
	ANDPD X12, X8
	ADDSD X8, X0
	ADDQ  AX, SI     // i++
	LOOP  absum_tail // } while --CX > 0

absum_end: // return p_sum_0
	MOVSD X0, sum+40(FP)
	RET

absum_zero_stride:
	PXOR  X0, X0
	MOVSD X0, sum+40(FP)
	RET
