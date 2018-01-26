// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Some of the loop unrolling code is copied from:
// http://golang.org/src/math/big/arith_amd64.s
// which is distributed under these terms:
//
// Copyright (c) 2012 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

//+build !noasm,!appengine

#include "textflag.h"

#define X_PTR SI
#define Y_PTR DI
#define DST_PTR DI
#define IDX AX
#define LEN CX
#define TAIL BX
#define INC_X R8
#define INCx3_X R11
#define INC_Y R9
#define INCx3_Y R12
#define INC_DST R9
#define INCx3_DST R12
#define ALPHA X0
#define ALPHA_Y Y0
#define G_IDX Y1

// func AxpyIncAVX(alpha float64, x, y []float64, n, incX, incY, ix, iy uintptr)
TEXT ·AxpyIncAVX(SB), NOSPLIT, $0
	MOVQ x_base+8(FP), X_PTR  // X_PTR = &x
	MOVQ y_base+32(FP), Y_PTR // Y_PTR = &y
	MOVQ n+56(FP), LEN        // LEN = n
	CMPQ LEN, $0              // if LEN == 0 { return }
	JE   end

	MOVQ ix+80(FP), INC_X
	MOVQ iy+88(FP), INC_Y
	LEAQ (X_PTR)(INC_X*8), X_PTR // X_PTR = &(x[ix])
	LEAQ (Y_PTR)(INC_Y*8), Y_PTR // Y_PTR = &(y[iy])
	MOVQ Y_PTR, DST_PTR          // DST_PTR = Y_PTR  // Write pointer

	MOVQ incX+64(FP), INC_X // INC_X = incX * sizeof(float64)
	SHLQ $3, INC_X
	MOVQ incY+72(FP), INC_Y // INC_Y = incY * sizeof(float64)
	SHLQ $3, INC_Y

	VMOVSD alpha+0(FP), ALPHA // ALPHA = alpha
	MOVQ   LEN, TAIL          // TAIL = LEN
	SHRQ   $2, LEN            // LEN = floor( n / 4 )
	JZ     tail_start         // if LEN == 0 { goto tail_start }

	LEAQ (INC_X)(INC_X*2), INCx3_X // INCx3_X = INC_X * 3
	LEAQ (INC_Y)(INC_Y*2), INCx3_Y // INCx3_Y = INC_Y * 3
	CMPQ INC_Y, $1
	JNE  loop
	CMPQ LEN, $4
	JLE  loop
	JMP  axpy_gather

// TODO: Branch on incY==1 to use VGATHER instructions.
loop:  // do {  // y[i] += alpha * x[i] unrolled 4x.
	VMOVSD (X_PTR), X2            // X_i = x[i]
	VMOVSD (X_PTR)(INC_X*1), X3
	VMOVSD (X_PTR)(INC_X*2), X4
	VMOVSD (X_PTR)(INCx3_X*1), X5

	VFMADD213SD (Y_PTR), ALPHA, X2            // X_i = X_i * a + y[i]
	VFMADD213SD (Y_PTR)(INC_Y*1), ALPHA, X3
	VFMADD213SD (Y_PTR)(INC_Y*2), ALPHA, X4
	VFMADD213SD (Y_PTR)(INCx3_Y*1), ALPHA, X5

	VMOVSD X2, (DST_PTR)              // y[i] = X_i
	VMOVSD X3, (DST_PTR)(INC_DST*1)
	VMOVSD X4, (DST_PTR)(INC_DST*2)
	VMOVSD X5, (DST_PTR)(INCx3_DST*1)

	LEAQ (X_PTR)(INC_X*4), X_PTR // X_PTR = &(X_PTR[incX*4])
	LEAQ (Y_PTR)(INC_Y*4), Y_PTR // Y_PTR = &(Y_PTR[incY*4])
	DECQ LEN
	JNZ  loop                    // } while --LEN > 0
	CMPQ TAIL, $0                // if TAIL == 0 { return }
	JE   end

tail_start: // Reset Loop registers
	TESTQ $2, TAIL // if TAIL & 2 == 0
	JZ    tail_one // { goto tail_one }

tail_two:
	VMOVSD      (X_PTR), X2                 // X_i = x[i]
	VMOVSD      (X_PTR)(INC_X*1), X3
	VFMADD213SD (Y_PTR), ALPHA, X2          // X_i = X_i * a + y[i]
	VFMADD213SD (Y_PTR)(INC_Y*1), ALPHA, X3
	VMOVSD      X2, (DST_PTR)               // y[i] = X_i
	VMOVSD      X3, (DST_PTR)(INC_DST*1)

	LEAQ (X_PTR)(INC_X*2), X_PTR // X_PTR = &(X_PTR[incX*2])
	LEAQ (Y_PTR)(INC_Y*2), Y_PTR // Y_PTR = &(Y_PTR[incY*2])

	TESTQ $1, TAIL
	JZ    end      // if TAIL == 0 { goto end }

tail_one:
	// y[i] += alpha * x[i] for the last n % 4 iterations.
	VMOVSD      (X_PTR), X2        // X2 = x[i]
	VFMADD213SD (Y_PTR), ALPHA, X2 // X_i = X_i * a + y[i]
	VMOVSD      X2, (DST_PTR)      // y[i] = X2

end:
	RET

axpy_gather:
	XORQ         IDX, IDX             // IDX = 0
	VPXOR        X1, X1, X1           // X1 = { 0, 0 }
	VMOVQ        INC_X, X1            // X1 = { 0, INC_X }
	VMOVQ        INCx3_X, X3          // X3 = { 0, INC_X }
	VPADDQ       X1, X1, X2           // X2 = 2 * INC_X
	VPADDQ       X3, X1, X4           // X4 = 4 * INC_X
	VSHUFPD      $1, X1, X1, X1       // X1 = { 0, INC_X }
	VSHUFPD      $0, X3, X2, X3       // X3 = { 2 * INC_X, 3 * INC_X }
	VINSERTI128  $1, X3, G_IDX, G_IDX // G_IDX = { 0, INC_X, 2 * INC_X, 3 * INC_X }
	VPCMPEQD     Y10, Y10, Y10        // set mask register to all 1's
	VBROADCASTSD ALPHA, ALPHA_Y       // ALPHA_Y = { alpha, alpha, alpha, alpha }

g_loop:
	VMOVUPS    Y10, Y9
	VGATHERQPD Y9, (X_PTR)(G_IDX * 1), Y2 // Y_i = X[IDX:IDX+3]

	VFMADD213PD (Y_PTR)(IDX*8), ALPHA_Y, Y2 // Y_i = Y_i * a + y[i]
	VMOVUPS     Y2, (DST_PTR)(IDX*8)        // y[i] = Y_i

	LEAQ (X_PTR)(INC_X*4), X_PTR // X_PTR = &(X_PTR[incX*4])
	ADDQ $4, IDX                 // i += 4
	DECQ LEN
	JNZ  g_loop

	CMPQ TAIL, $0 // if TAIL == 0 { return }
	JE   g_end

g_tail2:
	TESTQ $2, TAIL
	JZ    g_tail1

	VMOVUPS    Y10, Y9
	VGATHERQPD Y9, (X_PTR)(G_IDX * 1), Y2 // Y_i = X[IDX:IDX+3]

	VFMADD213PD (Y_PTR)(IDX*8), ALPHA_Y, Y2 // Y_i = Y_i * a + y[i]
	VMOVUPS     Y2, (DST_PTR)(IDX*8)        // y[i] = Y_i

	LEAQ (X_PTR)(INC_X*2), X_PTR // X_PTR = &(X_PTR[incX*4])
	ADDQ $2, IDX                 // i += 4

g_tail1:
	TESTQ $1, TAIL
	JZ    g_end

	VMOVUPS     X10, X9
	VGATHERQPD  X9, (X_PTR)(X1 * 1), X2
	VFMADD213PD (Y_PTR)(IDX*8), ALPHA, X2
	VMOVUPS     X2, (DST_PTR)(IDX*8)

g_end:
	RET
