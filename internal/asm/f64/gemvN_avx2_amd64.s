// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// AVX2/FMA optimized version of GemvN
// y = alpha * A * x + beta * y

//go:build !noasm && !gccgo && !safe
// +build !noasm,!gccgo,!safe

#include "textflag.h"

// func GemvNAVX2(m, n uintptr, alpha float64, a []float64, lda uintptr, x []float64, incX uintptr, beta float64, y []float64, incY uintptr)
TEXT ·GemvNAVX2(SB), NOSPLIT, $0
	MOVQ    m+0(FP), CX           // m
	MOVQ    n+8(FP), BX           // n
	MOVSD   alpha+16(FP), X15     // alpha
	MOVQ    a_base+24(FP), DI     // a
	MOVQ    lda+48(FP), R12       // lda
	MOVQ    x_base+56(FP), SI     // x
	MOVQ    incX+80(FP), R8       // incX
	MOVSD   beta+88(FP), X14      // beta
	MOVQ    y_base+96(FP), DX     // y
	MOVQ    incY+120(FP), R10     // incY
	
	TESTQ   CX, CX
	JZ      end
	TESTQ   BX, BX
	JZ      end
	
	// Adjust Y pointer for negative increment
	XORQ    R11, R11              // R11 = 0
	MOVQ    CX, R13               // R13 = m
	SUBQ    $1, R13               // R13 = m - 1
	IMULQ   R10, R13              // R13 = (m-1) * incY
	NEGQ    R13                   // R13 = -(m-1) * incY
	CMPQ    R10, $0               // if incY < 0
	CMOVQLT R13, R11              // R11 = R13
	LEAQ    (DX)(R11*8), DX       // y += R11
	
	// Adjust X pointer for negative increment
	XORQ    R11, R11              // R11 = 0
	MOVQ    BX, R13               // R13 = n
	SUBQ    $1, R13               // R13 = n - 1
	IMULQ   R8, R13               // R13 = (n-1) * incX
	NEGQ    R13                   // R13 = -(n-1) * incX
	CMPQ    R8, $0                // if incX < 0
	CMOVQLT R13, R11              // R11 = R13
	LEAQ    (SI)(R11*8), SI       // x += R11
	
	// Convert increments to bytes
	MOVQ    R8, R13
	SHLQ    $3, R8                // incX *= 8
	MOVQ    R10, R14
	SHLQ    $3, R10               // incY *= 8
	SHLQ    $3, R12               // lda *= 8
	
	// Handle beta scaling
	VXORPD  Y13, Y13, Y13         // Y13 = 0
	VUCOMISD X13, X14             // Compare beta with 0
	JE      clear_y               // If beta == 0, clear y
	
	// Scale y by beta if beta != 0
	MOVQ    DX, R9
	MOVQ    CX, R15
scale_y_loop:
	VMOVSD  (R9), X0
	VMULSD  X14, X0, X0
	VMOVSD  X0, (R9)
	ADDQ    R10, R9
	DECQ    R15
	JNZ     scale_y_loop
	JMP     check_alpha
	
clear_y:
	// Clear y when beta == 0
	MOVQ    DX, R9
	MOVQ    CX, R15
	VXORPD  X0, X0, X0
clear_y_loop:
	VMOVSD  X0, (R9)
	ADDQ    R10, R9
	DECQ    R15
	JNZ     clear_y_loop
	
check_alpha:
	// For each row of A
	// Note: We always compute even if alpha=0 to propagate NaN
row_loop:
	MOVQ    DI, AX                // Current row pointer
	MOVQ    SI, R11               // Reset x pointer
	MOVQ    BX, R9                // n counter
	
	VXORPD  Y0, Y0, Y0            // Clear accumulator
	
	// Check if we can use unit stride for x
	CMPQ    R13, $1               // Compare original incX with 1
	JNE     non_unit_stride
	
	// Unit stride - process 4 elements at a time
	MOVQ    R9, R14
	SHRQ    $2, R14               // R14 = n / 4
	JZ      unit_remainder
	
unit_loop_4:
	VMOVUPD (AX), Y1              // Load 4 elements from A
	VMOVUPD (R11), Y2             // Load 4 elements from x
	VFMADD231PD Y2, Y1, Y0        // Y0 += Y1 * Y2
	
	ADDQ    $32, AX
	ADDQ    $32, R11
	DECQ    R14
	JNZ     unit_loop_4
	
unit_remainder:
	// First reduce Y0 to a scalar
	VEXTRACTF128 $1, Y0, X1       // X1 = upper 128 bits of Y0
	VADDPD  X1, X0, X0            // X0 = X0 + X1
	VHADDPD X0, X0, X0            // X0 = horizontal sum
	
	// Now process remainder elements
	MOVQ    R9, R14
	ANDQ    $3, R14               // R14 = n % 4
	JZ      final_scale
	
unit_loop_1:
	VMOVSD  (AX), X1
	VMOVSD  (R11), X2
	VMULSD  X2, X1, X1            // X1 = X1 * X2
	VADDSD  X1, X0, X0            // X0 += X1
	
	ADDQ    $8, AX
	ADDQ    $8, R11
	DECQ    R14
	JNZ     unit_loop_1
	JMP     final_scale
	
non_unit_stride:
	// Non-unit stride - process one element at a time
non_unit_loop:
	VMOVSD  (AX), X1
	VMOVSD  (R11), X2
	VMULSD  X2, X1, X1
	VADDSD  X1, X0, X0
	
	ADDQ    $8, AX
	ADDQ    R8, R11               // Add incX (in bytes)
	DECQ    R9
	JNZ     non_unit_loop
	
reduce:
	// Horizontal sum of Y0
	VEXTRACTF128 $1, Y0, X1       // X1 = upper 128 bits of Y0
	VADDPD  X1, X0, X0            // X0 = X0 + X1
	VHADDPD X0, X0, X0            // Horizontal add to get final sum in X0[0]
	
final_scale:
	// Scale by alpha
	VMULSD  X15, X0, X0
	
	// Load y[i] and add (beta scaling already done)
	VADDSD  (DX), X0, X0
	
	// Store result
	VMOVSD  X0, (DX)
	
	// Advance to next row
	ADDQ    R12, DI               // a += lda
	ADDQ    R10, DX               // y += incY
	
	DECQ    CX
	JNZ     row_loop
	
end:
	VZEROUPPER
	RET
