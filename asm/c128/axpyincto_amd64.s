// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

#include "textflag.h"

// MOVDDUP X2, X3
#define MOVDDUP_X2_X3 BYTE $0xF2; BYTE $0x0F; BYTE $0x12; BYTE $0xDA
// MOVDDUP X4, X5
#define MOVDDUP_X4_X5 BYTE $0xF2; BYTE $0x0F; BYTE $0x12; BYTE $0xEC
// MOVDDUP X6, X7
#define MOVDDUP_X6_X7 BYTE $0xF2; BYTE $0x0F; BYTE $0x12; BYTE $0xFE
// MOVDDUP X8, X9
#define MOVDDUP_X8_X9 BYTE $0xF2; BYTE $0x45; BYTE $0x0F; BYTE $0x12; BYTE $0xC8

// ADDSUBPD X2, X3
#define ADDSUBPD_X2_X3 BYTE $0x66; BYTE $0x0F; BYTE $0xD0; BYTE $0xDA
// ADDSUBPD X4, X5
#define ADDSUBPD_X4_X5 BYTE $0x66; BYTE $0x0F; BYTE $0xD0; BYTE $0xEC
// ADDSUBPD X6, X7
#define ADDSUBPD_X6_X7 BYTE $0x66; BYTE $0x0F; BYTE $0xD0; BYTE $0xFE
// ADDSUBPD X8, X9
#define ADDSUBPD_X8_X9 BYTE $0x66; BYTE $0x45; BYTE $0x0F; BYTE $0xD0; BYTE $0xC8

// func AxpyIncTo(dst []complex128, incDst, idst uintptr, alpha complex128, x, y []complex128, n, incX, incY, ix, iy uintptr)
TEXT ·AxpyIncTo(SB), NOSPLIT, $0
	MOVQ   dst_base+0(FP), DI
	MOVQ   x_base+56(FP), SI
	MOVQ   y_base+80(FP), DX
	MOVQ   n+104(FP), CX
	CMPQ   CX, $0             // if n==0, return
	JE     axpyi_end
	MOVQ   ix+128(FP), R8     // Load the first indicies
	SHLQ   $1, R8             // Double to adjust for 16-byte size
	MOVQ   iy+136(FP), R9
	SHLQ   $1, R9
	MOVQ   idst+32(FP), R10
	SHLQ   $1, R10
	LEAQ   (SI)(R8*8), SI     // Calculate addrress of first indicies
	LEAQ   (DX)(R9*8), DX
	LEAQ   (DI)(R10*8), DI
	MOVQ   incX+112(FP), R8   // Incrementors*8 for easy iteration (ADDQ)
	SHLQ   $4, R8
	MOVQ   incY+120(FP), R9
	SHLQ   $4, R9
	MOVQ   incDst+24(FP), R10
	SHLQ   $4, R10
	MOVUPS alpha+40(FP), X0   // (ar,ai)
	MOVAPS X0, X1
	SHUFPD $0x1, X1, X1       // (ai,ar)
	MOVAPS X0, X10
	MOVAPS X1, X11
	MOVQ   CX, BX
	ANDQ   $3, CX
	SHRQ   $2, BX
	JZ     axpyi_tail

axpyi_loop:
	MOVUPS (SI), X2
	MOVUPS (SI)(R8*1), X4
	LEAQ   (SI)(R8*2), SI
	MOVUPS (SI), X6
	MOVUPS (SI)(R8*1), X8
	MOVDDUP_X2_X3          // Load and duplicate imag elements (xi, xi)
	SHUFPD $0x3, X2, X2    // duplicate real elements (xr, xr)
	MOVDDUP_X4_X5
	SHUFPD $0x3, X4, X4
	MOVDDUP_X6_X7
	SHUFPD $0x3, X6, X6
	MOVDDUP_X8_X9
	SHUFPD $0x3, X8, X8
	MULPD  X1, X2          // (ai*xr, ar*xr)
	MULPD  X0, X3          // (ar*xi, ai*xi)
	MULPD  X11, X4
	MULPD  X10, X5
	MULPD  X1, X6
	MULPD  X0, X7
	MULPD  X11, X8
	MULPD  X10, X9
	ADDSUBPD_X2_X3         // Add/Sub to (ai*xr + ar*xi , ar*xr - (ai*xi))
	ADDSUBPD_X4_X5
	ADDSUBPD_X6_X7
	ADDSUBPD_X8_X9
	ADDPD  (DX), X3
	ADDPD  (DX)(R9*1), X5
	LEAQ   (DX)(R9*2), DX
	ADDPD  (DX), X7
	ADDPD  (DX)(R9*1), X9
	MOVUPS X3, (DI)        // Write result back to dst
	MOVUPS X5, (DI)(R10*1)
	LEAQ   (DI)(R10*2), DI
	MOVUPS X7, (DI)
	MOVUPS X9, (DI)(R10*1)
	LEAQ   (SI)(R8*2), SI  // Increment addresses
	LEAQ   (DX)(R9*2), DX
	LEAQ   (DI)(R10*2), DI
	DECQ   BX
	JNZ    axpyi_loop
	CMPQ   CX, $0
	JE     axpyi_end

axpyi_tail:
	MOVUPS (SI), X2
	MOVDDUP_X2_X3       // Load and duplicate imag elements (xi, xi)
	SHUFPD $0x3, X2, X2 // duplicate real elements (xr, xr)
	MULPD  X1, X2       // (ai*x2r, ar*x2r, ai*x1r, ar*x1r)
	MULPD  X0, X3       // (ar*x2i, ai*x2i, ar*x1i, ai*x1i)
	ADDSUBPD_X2_X3      // (ai*x2r+ar*x2i, ar*x2r-ai*x2i, ai*x1r+ar*x1i, ar*x1r-ai*x1i)
	ADDPD  (DX), X3
	MOVUPS X3, (DI)
	ADDQ   R8, SI       // Increment addresses
	ADDQ   R9, DX
	ADDQ   R10, DI
	LOOP   axpyi_tail

axpyi_end:
	RET
