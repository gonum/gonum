// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// AVX2/FMA3 optimized version of AxpyUnitary
// y[i] += alpha * x[i]

// +build !noasm,!gccgo,!safe

#include "textflag.h"

#define X_PTR SI
#define Y_PTR DI
#define DST_PTR DI
#define IDX AX
#define LEN CX
#define TAIL BX
#define ALPHA Y0
#define ALPHA_2 Y1

// func AxpyUnitaryFMA(alpha float64, x, y []float64)
TEXT ·AxpyUnitaryFMA(SB), NOSPLIT, $0
	MOVQ    x_base+8(FP), X_PTR  // X_PTR := &x
	MOVQ    y_base+32(FP), Y_PTR // Y_PTR := &y
	MOVQ    x_len+16(FP), LEN    // LEN = min( len(x), len(y) )
	CMPQ    y_len+40(FP), LEN
	CMOVQLE y_len+40(FP), LEN
	CMPQ    LEN, $0              // if LEN == 0 { return }
	JE      end
	XORQ    IDX, IDX
	
	// Broadcast alpha to all lanes of YMM register
	VBROADCASTSD alpha+0(FP), ALPHA   // ALPHA := { alpha, alpha, alpha, alpha }
	VMOVUPD      ALPHA, ALPHA_2       // ALPHA_2 := ALPHA for pipelining
	
	MOVQ    Y_PTR, TAIL          // Check memory alignment
	ANDQ    $31, TAIL            // TAIL = &y % 32 (for 32-byte AVX2 alignment)
	JZ      no_trim              // if TAIL == 0 { goto no_trim }

	// Align on 32-byte boundary using scalar operations
align_loop:
	MOVSD (X_PTR)(IDX*8), X2     // X2 := x[i]
	MULSD alpha+0(FP), X2        // X2 *= alpha (scalar)
	ADDSD (Y_PTR)(IDX*8), X2     // X2 += y[i]
	MOVSD X2, (DST_PTR)(IDX*8)   // y[i] = X2
	INCQ  IDX                    // i++
	DECQ  LEN                    // LEN--
	JZ    end                    // if LEN == 0 { return }
	MOVQ  Y_PTR, TAIL
	LEAQ  (TAIL)(IDX*8), TAIL
	ANDQ  $31, TAIL
	JNZ   align_loop

no_trim:
	MOVQ LEN, TAIL
	ANDQ $15, TAIL   // TAIL := n % 16 (16 doubles with AVX2)
	SHRQ $4, LEN     // LEN = floor( n / 16 )
	JZ   tail_start  // if LEN == 0 { goto tail_start }

loop:  // Main loop: process 16 doubles per iteration using AVX2/FMA
	// Load x[i:i+16] into 4 YMM registers
	VMOVUPD (X_PTR)(IDX*8), Y2        // Y2 = x[i:i+4]
	VMOVUPD 32(X_PTR)(IDX*8), Y3      // Y3 = x[i+4:i+8]
	VMOVUPD 64(X_PTR)(IDX*8), Y4      // Y4 = x[i+8:i+12]
	VMOVUPD 96(X_PTR)(IDX*8), Y5      // Y5 = x[i+12:i+16]
	
	// Use FMA to compute y[i] = y[i] + alpha * x[i] in one instruction
	// VFMADD213PD: Y2 = Y2 * ALPHA + (Y_PTR)
	VMOVUPD (Y_PTR)(IDX*8), Y6         // Y6 = y[i:i+4]
	VMOVUPD 32(Y_PTR)(IDX*8), Y7       // Y7 = y[i+4:i+8]
	VMOVUPD 64(Y_PTR)(IDX*8), Y8       // Y8 = y[i+8:i+12]
	VMOVUPD 96(Y_PTR)(IDX*8), Y9       // Y9 = y[i+12:i+16]
	
	VFMADD213PD Y6, ALPHA, Y2          // Y2 = x[i:i+4] * alpha + y[i:i+4]
	VFMADD213PD Y7, ALPHA_2, Y3        // Y3 = x[i+4:i+8] * alpha + y[i+4:i+8]
	VFMADD213PD Y8, ALPHA, Y4          // Y4 = x[i+8:i+12] * alpha + y[i+8:i+12]
	VFMADD213PD Y9, ALPHA_2, Y5        // Y5 = x[i+12:i+16] * alpha + y[i+12:i+16]
	
	// Store results back to y
	VMOVUPD Y2, (DST_PTR)(IDX*8)       // y[i:i+4] = Y2
	VMOVUPD Y3, 32(DST_PTR)(IDX*8)    // y[i+4:i+8] = Y3
	VMOVUPD Y4, 64(DST_PTR)(IDX*8)    // y[i+8:i+12] = Y4
	VMOVUPD Y5, 96(DST_PTR)(IDX*8)    // y[i+12:i+16] = Y5
	
	ADDQ $16, IDX  // i += 16
	DECQ LEN
	JNZ  loop      // } while --LEN > 0
	CMPQ TAIL, $0  // if TAIL == 0 { return }
	JE   end

tail_start: // Handle remaining elements
	MOVQ TAIL, LEN // Loop counter: LEN = TAIL
	SHRQ $2, LEN   // LEN = floor( TAIL / 4 ) for AVX2
	JZ   tail_sse  // if LEN == 0 { goto tail_sse }

tail_avx2: // Process 4 doubles at a time with FMA
	VMOVUPD (X_PTR)(IDX*8), Y2         // Y2 = x[i:i+4]
	VMOVUPD (Y_PTR)(IDX*8), Y3         // Y3 = y[i:i+4]
	VFMADD213PD Y3, ALPHA, Y2          // Y2 = x * alpha + y
	VMOVUPD Y2, (DST_PTR)(IDX*8)       // y[i:i+4] = Y2
	ADDQ   $4, IDX                     // i += 4
	DECQ   LEN
	JNZ    tail_avx2                   // } while --LEN > 0

tail_sse:
	MOVQ TAIL, LEN
	ANDQ $3, LEN    // LEN = TAIL % 4
	SHRQ $1, LEN    // LEN = LEN / 2
	JZ   tail_one   // if TAIL == 0 { goto tail_one }

	// Process 2 doubles using SSE
	MOVUPD (X_PTR)(IDX*8), X2          // X2 = x[i:i+2]
	MOVSD  alpha+0(FP), X3             // X3 = alpha
	SHUFPD $0, X3, X3                  // X3 = {alpha, alpha}
	MULPD  X3, X2                      // X2 *= alpha
	ADDPD  (Y_PTR)(IDX*8), X2          // X2 += y[i:i+2]
	MOVUPD X2, (DST_PTR)(IDX*8)        // y[i:i+2] = X2
	ADDQ   $2, IDX                     // i += 2

tail_one:
	ANDQ   $1, TAIL
	JZ     end      // if no more elements { goto end }
	
	// Process last element
	MOVSD (X_PTR)(IDX*8), X2          // X2 = x[i]
	MULSD alpha+0(FP), X2              // X2 *= alpha
	ADDSD (Y_PTR)(IDX*8), X2           // X2 += y[i]
	MOVSD X2, (DST_PTR)(IDX*8)         // y[i] = X2

end:
	VZEROUPPER  // Clear upper YMM state
	RET
