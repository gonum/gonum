// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

#include "textflag.h"

// func AxpyUnitaryTo(dst []float32, alpha float32, x, y []float32)
TEXT ·AxpyUnitaryTo(SB), NOSPLIT, $0
	MOVQ    dst_base+0(FP), DI // Load data buffer pointers
	MOVQ    x_base+32(FP), SI
	MOVQ    y_base+56(FP), DX
	MOVQ    x_len+40(FP), BX   // BX = min( len(x), len(y), len(dst) )
	CMPQ    y_len+64(FP), BX
	CMOVQLE y_len+64(FP), BX
	CMPQ    dst_len+8(FP), BX
	CMOVQLE dst_len+8(FP), BX
	CMPQ    BX, $0             // Empty return
	JE      axpy_end
	MOVSS   alpha+24(FP), X0
	SHUFPS  $0, X0, X0         // Load alpha 4 times
	XORQ    AX, AX             // i = 0
	MOVQ    DX, CX
	ANDQ    $0xF, CX           // Align on 16-byte boundary for ADDPS
	JZ      axpy_no_trim

	XORQ $0xF, CX
	INCQ CX
	SHRQ $2, CX

axpy_align: // Trim first value(s) in unaligned buffer
	MOVSS (SI)(AX*4), X2
	MULSS X0, X2
	ADDSS (DX)(AX*4), X2
	MOVSS X2, (DI)(AX*4)
	INCQ  AX
	DECQ  BX
	JZ    axpy_end       // Zero check for small unaligned slices
	LOOP  axpy_align

axpy_no_trim:
	MOVUPS X0, X1           // Copy to X1 for pipelining
	MOVQ   BX, CX
	ANDQ   $0xF, BX         // BX = len % 16
	SHRQ   $4, CX           // CX = int(len / 16)
	JZ     axpy_tail4_start

axpy_loop: // Loop unrolled 16x
	MOVUPS (SI)(AX*4), X2   // xmm = x[i:i+4]
	MOVUPS 16(SI)(AX*4), X3
	MOVUPS 32(SI)(AX*4), X4
	MOVUPS 48(SI)(AX*4), X5
	MULPS  X0, X2           // xmm *= a
	MULPS  X1, X3
	MULPS  X0, X4
	MULPS  X1, X5
	ADDPS  (DX)(AX*4), X2   // xmm += y[i:i+4]
	ADDPS  16(DX)(AX*4), X3
	ADDPS  32(DX)(AX*4), X4
	ADDPS  48(DX)(AX*4), X5
	MOVUPS X2, (DI)(AX*4)   // dst[i:i+4] = xmm
	MOVUPS X3, 16(DI)(AX*4)
	MOVUPS X4, 32(DI)(AX*4)
	MOVUPS X5, 48(DI)(AX*4)
	ADDQ   $16, AX          // i+=16
	LOOP   axpy_loop        // while (--CX) > 0
	CMPQ   BX, $0
	JE     axpy_end

axpy_tail4_start:
	MOVQ BX, CX
	SHRQ $2, CX
	JZ   axpy_tail_start

axpy_tail4:
	MOVUPS (SI)(AX*4), X2
	MULPS  X0, X2
	ADDPS  (DX)(AX*4), X2
	MOVUPS X2, (DI)(AX*4)
	ADDQ   $4, AX
	LOOP   axpy_tail4

axpy_tail_start:
	MOVQ BX, CX
	ANDQ $3, CX
	JZ   axpy_end

axpy_tail:
	MOVSS (SI)(AX*4), X1
	MULSS X0, X1
	ADDSS (DX)(AX*4), X1
	MOVSS X1, (DI)(AX*4)
	INCQ  AX
	LOOP  axpy_tail

axpy_end:
	RET
