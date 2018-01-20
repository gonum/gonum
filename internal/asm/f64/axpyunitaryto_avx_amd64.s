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
#define Y_PTR DX
#define DST_PTR DI
#define IDX AX
#define LEN CX
#define TAIL BX
#define ALPHA Y0
#define ALPHA_X X0

// func AxpyUnitaryToAVX(dst []float64, alpha float64, x, y []float64)
TEXT ·AxpyUnitaryToAVX(SB), NOSPLIT, $0
	MOVQ    dst_base+0(FP), DST_PTR // DST_PTR := &dst
	MOVQ    x_base+32(FP), X_PTR    // X_PTR := &x
	MOVQ    y_base+56(FP), Y_PTR    // Y_PTR := &y
	MOVQ    x_len+40(FP), LEN       // LEN = min( len(x), len(y), len(dst) )
	CMPQ    y_len+64(FP), LEN
	CMOVQLE y_len+64(FP), LEN
	CMPQ    dst_len+8(FP), LEN
	CMOVQLE dst_len+8(FP), LEN

	CMPQ LEN, $0
	JE   end     // if LEN == 0 { return }

	XORQ         IDX, IDX            // IDX = 0
	VBROADCASTSD alpha+24(FP), ALPHA // ALPHA := { alpha, alpha }

no_trim:
	MOVQ LEN, TAIL
	ANDQ $15, TAIL  // TAIL := n % 8
	SHRQ $4, LEN    // LEN = floor( n / 8 )
	JZ   axpy_tail8 // if LEN == 0 { goto axpy_tail8 }

loop:  // do {
	// y[i] += alpha * x[i] unrolled 8x.
	VMOVUPS (X_PTR)(IDX*8), Y2   // X_i = x[i]
	VMOVUPS 32(X_PTR)(IDX*8), Y3
	VMOVUPS 64(X_PTR)(IDX*8), Y4
	VMOVUPS 96(X_PTR)(IDX*8), Y5

	VFMADD213PD (Y_PTR)(IDX*8), ALPHA, Y2   // X_i = X_i * a + y[i]
	VFMADD213PD 32(Y_PTR)(IDX*8), ALPHA, Y3
	VFMADD213PD 64(Y_PTR)(IDX*8), ALPHA, Y4
	VFMADD213PD 96(Y_PTR)(IDX*8), ALPHA, Y5

	VMOVUPS Y2, (DST_PTR)(IDX*8)   // y[i] = X_i
	VMOVUPS Y3, 32(DST_PTR)(IDX*8)
	VMOVUPS Y4, 64(DST_PTR)(IDX*8)
	VMOVUPS Y5, 96(DST_PTR)(IDX*8)

	ADDQ $16, IDX // i += 8
	DECQ LEN
	JNZ  loop     // } while --LEN > 0
	CMPQ TAIL, $0 // if TAIL == 0 { return }
	JE   end

axpy_tail8:
	TESTQ $8, TAIL   // if TAIL & 8 == 0 { goto axpy_tail4}
	JZ    axpy_tail4

	VMOVUPS (X_PTR)(IDX*8), Y2   // X_i = x[i]
	VMOVUPS 32(X_PTR)(IDX*8), Y3

	VFMADD213PD (Y_PTR)(IDX*8), ALPHA, Y2   // X_i = X_i * a + y[i]
	VFMADD213PD 32(Y_PTR)(IDX*8), ALPHA, Y3

	VMOVUPS Y2, (DST_PTR)(IDX*8)   // y[i] = X_i
	VMOVUPS Y3, 32(DST_PTR)(IDX*8)

	ADDQ $8, IDX // i += 8

axpy_tail4:
	TESTQ $4, TAIL   // if TAIL & 4 == 0 { goto axpy_tail2}
	JZ    axpy_tail2

	VMOVUPS     (X_PTR)(IDX*8), Y2        // X_i = x[i]
	VFMADD213PD (Y_PTR)(IDX*8), ALPHA, Y2 // X_i = X_i * a + y[i]
	VMOVUPS     Y2, (DST_PTR)(IDX*8)      // y[i] = X_i

	ADDQ $4, IDX // i += 8

axpy_tail2:
	TESTQ $2, TAIL   // if TAIL & 2 == 0 { goto axpy_tail1}
	JZ    axpy_tail1

	VMOVUPS     (X_PTR)(IDX*8), X2          // X_i = x[i]
	VFMADD213PD (Y_PTR)(IDX*8), ALPHA_X, X2 // X_i = X_i * a + y[i]
	VMOVUPS     X2, (DST_PTR)(IDX*8)        // y[i] = X_i

	ADDQ $2, IDX // i += 8

axpy_tail1:
	TESTQ $1, TAIL
	JZ    end

	VMOVSD      (X_PTR)(IDX*8), X2          // X2 = x[i]
	VFMADD213SD (Y_PTR)(IDX*8), ALPHA_X, X2 // X2 = X2 * a + y[i]
	VMOVSD      X2, (DST_PTR)(IDX*8)        // y[i] = X2

end:
	VZEROALL
	RET
