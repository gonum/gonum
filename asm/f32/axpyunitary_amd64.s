// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

#include "textflag.h"

// func AxpyUnitary(alpha float32, x, y []float32)
TEXT ·AxpyUnitary(SB), NOSPLIT, $0
	MOVQ    x_base+8(FP), SI  // Load data buffer pointers
	MOVQ    y_base+32(FP), DI
	MOVQ    x_len+16(FP), BX  // BX = min( len(x), len(y) )
	CMPQ    y_len+40(FP), BX
	CMOVQLE y_len+40(FP), BX
	CMPQ    BX, $0
	JE      axpy_end
	MOVSS   alpha+0(FP), X0
	SHUFPS  $0, X0, X0        // Load alpha into X0 4 times
	MOVUPS  X0, X2            // Copy to X2 for pipelining
	XORQ    AX, AX            // i = 0
	PXOR    X1, X1            // 2 NOP instructions (PXOR) to align
	PXOR    X3, X3            // loop to cache line
	XORQ    AX, AX            // i = 0
	MOVQ    DI, CX
	ANDQ    $0xF, CX          // Align on 16-byte boundary for ADDPS
	JZ      axpy_no_trim
	SHRQ    $2, CX

axpy_align: // Trim first value(s) in unaligned buffer
	MOVSS (SI)(AX*4), X2
	MULSS X0, X2
	ADDSS (DI)(AX*4), X2
	MOVSS X2, (DI)(AX*4)
	INCQ  AX
	DECQ  BX
	JZ    axpy_end       // Zero check for small unaligned slices
	LOOP  axpy_align

axpy_no_trim:
	MOVUPS X0, X1          // Copy to X1 for pipelining
	MOVQ   BX, CX
	ANDQ   $0xF, BX        // BX = len % 16
	SHRQ   $4, CX          // CX = int(len / 16)
	JZ     axpy_tail_start

axpy_loop:
	MOVUPS (SI)(AX*4), X2   // xmm = x[i:i+4]
	MOVUPS 16(SI)(AX*4), X3
	MOVUPS 32(SI)(AX*4), X4
	MOVUPS 48(SI)(AX*4), X5
	MULPS  X0, X2           // xmm *= a
	MULPS  X1, X3
	MULPS  X0, X4
	MULPS  X1, X5
	ADDPS  (DI)(AX*4), X2   // xmm += y[i:i+4]
	ADDPS  16(DI)(AX*4), X3
	ADDPS  32(DI)(AX*4), X4
	ADDPS  48(DI)(AX*4), X5
	MOVUPS X2, (DI)(AX*4)   // dst[i:i+4] = xmm
	MOVUPS X3, 16(DI)(AX*4)
	MOVUPS X4, 32(DI)(AX*4)
	MOVUPS X5, 48(DI)(AX*4)
	ADDQ   $16, AX          // i+=16
	LOOP   axpy_loop        // while (--CX) > 0
	CMPQ   BX, $0
	JE     axpy_end

axpy_tail_start:
	MOVQ BX, CX

axpy_tail:
	MOVSS (SI)(AX*4), X1
	MULSS X0, X1
	ADDSS (DI)(AX*4), X1
	MOVSS X1, (DI)(AX*4)
	INCQ  AX
	LOOP  axpy_tail

axpy_end:
	RET
