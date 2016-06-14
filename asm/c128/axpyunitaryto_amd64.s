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

// func AxpyUnitaryTo(dst []complex128, alpha complex64, x, y []complex128)
TEXT ·AxpyUnitaryTo(SB), NOSPLIT, $0
	MOVQ    dst_base+0(FP), DI
	MOVQ    x_base+40(FP), SI
	MOVQ    y_base+64(FP), DX
	MOVQ    x_len+48(FP), CX
	CMPQ    y_len+72(FP), CX
	CMOVQLE y_len+72(FP), CX
	CMPQ    dst_len+8(FP), CX
	CMOVQLE dst_len+8(FP), CX
	CMPQ    CX, $0
	JE      caxy_end
	MOVUPS  alpha+24(FP), X0   // (ar,ai)
	MOVAPS  X0, X1
	SHUFPD  $0x1, X1, X1       // (ai,ar)
	XORQ    AX, AX
	MOVAPS  X0, X10
	MOVAPS  X1, X11
	MOVQ    CX, BX
	ANDQ    $3, CX
	SHRQ    $2, BX
	JZ      caxy_tail

caxy_loop:
	MOVUPS (SI)(AX*8), X2
	MOVUPS 16(SI)(AX*8), X4
	MOVUPS 32(SI)(AX*8), X6
	MOVUPS 48(SI)(AX*8), X8
	MOVDDUP_X2_X3           // Load and duplicate imag elements (xi, xi)
	SHUFPD $0x3, X2, X2     // duplicate real elements (xr, xr)
	MOVDDUP_X4_X5
	SHUFPD $0x3, X4, X4
	MOVDDUP_X6_X7
	SHUFPD $0x3, X6, X6
	MOVDDUP_X8_X9
	SHUFPD $0x3, X8, X8
	MULPD  X1, X2           // (ai*xr, ar*xr)
	MULPD  X0, X3           // (ar*xi, ai*xi)
	MULPD  X11, X4
	MULPD  X10, X5
	MULPD  X1, X6
	MULPD  X0, X7
	MULPD  X11, X8
	MULPD  X10, X9
	ADDSUBPD_X2_X3          // Add/Sub to (ai*xr + ar*xi , ar*xr - (ai*xi))
	ADDSUBPD_X4_X5
	ADDSUBPD_X6_X7
	ADDSUBPD_X8_X9
	ADDPD  (DX)(AX*8), X3   // Add y2,y1 to a*(x2,x1)
	ADDPD  16(DX)(AX*8), X5
	ADDPD  32(DX)(AX*8), X7
	ADDPD  48(DX)(AX*8), X9
	MOVUPS X3, (DI)(AX*8)   // Write result back to y2,y1
	MOVUPS X5, 16(DI)(AX*8)
	MOVUPS X7, 32(DI)(AX*8)
	MOVUPS X9, 48(DI)(AX*8)
	ADDQ   $8, AX
	DECQ   BX
	JNZ    caxy_loop
	CMPQ   CX, $0
	JE     caxy_end

caxy_tail: // Same calculation, but read in values to avoid trampling memory
	MOVUPS (SI)(AX*8), X2
	MOVDDUP_X2_X3         // Load and duplicate imag elements (xi, xi)
	SHUFPD $0x3, X2, X2   // duplicate real elements (xr, xr)
	MULPD  X1, X2         // (ai*x2r, ar*x2r, ai*x1r, ar*x1r)
	MULPD  X0, X3         // (ar*x2i, ai*x2i, ar*x1i, ai*x1i)
	ADDSUBPD_X2_X3        // (ai*x2r+ar*x2i, ar*x2r-ai*x2i, ai*x1r+ar*x1i, ar*x1r-ai*x1i)
	ADDPD  (DX)(AX*8), X3
	MOVUPS X3, (DI)(AX*8)
	ADDQ   $2, AX
	LOOP   caxy_tail

caxy_end:
	RET
