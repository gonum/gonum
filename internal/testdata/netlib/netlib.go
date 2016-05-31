// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netlib

// void dlahr2_(int* n, int* k, int* nb, double* a, int* lda, double* tau, double* t, int* ldt, double* y, int* ldy);
import "C"

func Dlahr2(n, k, nb int, a []float64, lda int, tau, t []float64, ldt int, y []float64, ldy int) {
	func() {
		n := C.int(n)
		k := C.int(k)
		nb := C.int(nb)
		lda := C.int(lda)
		ldt := C.int(ldt)
		ldy := C.int(ldy)
		C.dlahr2_((*C.int)(&n), (*C.int)(&k), (*C.int)(&nb),
			(*C.double)(&a[0]), (*C.int)(&lda),
			(*C.double)(&tau[0]),
			(*C.double)(&t[0]), (*C.int)(&ldt),
			(*C.double)(&y[0]), (*C.int)(&ldy))
	}()
}
