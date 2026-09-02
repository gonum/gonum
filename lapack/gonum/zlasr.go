// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/lapack"
)

// Zlasr applies a sequence of real plane rotations to the complex m×n matrix A.
// The series of plane rotations is implicitly represented by a matrix P, which
// is multiplied by A on the left or right depending on side:
//
//	A = P * A  if side == blas.Left
//	A = A * P^T  if side == blas.Right
//
// The layout of P is determined by pivot as in Dlasr. If direct == lapack.Forward
// the rotations are applied in order P = P(z-1) * ... * P(2) * P(1); if
// direct == lapack.Backward the order is reversed.
//
// c and s are real (cosine and sine of the plane rotations). They must each
// have length m - 1 if side == blas.Left and n - 1 if side == blas.Right.
//
// Zlasr is an internal routine. It is exported for testing purposes.
func (impl Implementation) Zlasr(side blas.Side, pivot lapack.Pivot, direct lapack.Direct, m, n int, c, s []float64, a []complex128, lda int) {
	switch {
	case side != blas.Left && side != blas.Right:
		panic(badSide)
	case pivot != lapack.Variable && pivot != lapack.Top && pivot != lapack.Bottom:
		panic(badPivot)
	case direct != lapack.Forward && direct != lapack.Backward:
		panic(badDirect)
	case m < 0:
		panic(mLT0)
	case n < 0:
		panic(nLT0)
	case lda < max(1, n):
		panic(badLdA)
	}

	if m == 0 || n == 0 {
		return
	}

	if side == blas.Left {
		if len(c) < m-1 {
			panic(shortC)
		}
		if len(s) < m-1 {
			panic(shortS)
		}
	} else {
		if len(c) < n-1 {
			panic(shortC)
		}
		if len(s) < n-1 {
			panic(shortS)
		}
	}
	if len(a) < (m-1)*lda+n {
		panic(shortA)
	}

	if side == blas.Left {
		if pivot == lapack.Variable {
			if direct == lapack.Forward {
				for j := 0; j < m-1; j++ {
					ctmp := complex(c[j], 0)
					stmp := complex(s[j], 0)
					if c[j] != 1 || s[j] != 0 {
						for i := 0; i < n; i++ {
							tmp2 := a[j*lda+i]
							tmp := a[(j+1)*lda+i]
							a[(j+1)*lda+i] = ctmp*tmp - stmp*tmp2
							a[j*lda+i] = stmp*tmp + ctmp*tmp2
						}
					}
				}
				return
			}
			for j := m - 2; j >= 0; j-- {
				ctmp := complex(c[j], 0)
				stmp := complex(s[j], 0)
				if c[j] != 1 || s[j] != 0 {
					for i := 0; i < n; i++ {
						tmp2 := a[j*lda+i]
						tmp := a[(j+1)*lda+i]
						a[(j+1)*lda+i] = ctmp*tmp - stmp*tmp2
						a[j*lda+i] = stmp*tmp + ctmp*tmp2
					}
				}
			}
			return
		} else if pivot == lapack.Top {
			if direct == lapack.Forward {
				for j := 1; j < m; j++ {
					ctmp := complex(c[j-1], 0)
					stmp := complex(s[j-1], 0)
					if c[j-1] != 1 || s[j-1] != 0 {
						for i := 0; i < n; i++ {
							tmp := a[j*lda+i]
							tmp2 := a[i]
							a[j*lda+i] = ctmp*tmp - stmp*tmp2
							a[i] = stmp*tmp + ctmp*tmp2
						}
					}
				}
				return
			}
			for j := m - 1; j >= 1; j-- {
				ctmp := complex(c[j-1], 0)
				stmp := complex(s[j-1], 0)
				if c[j-1] != 1 || s[j-1] != 0 {
					for i := 0; i < n; i++ {
						tmp := a[j*lda+i]
						tmp2 := a[i]
						a[j*lda+i] = ctmp*tmp - stmp*tmp2
						a[i] = stmp*tmp + ctmp*tmp2
					}
				}
			}
			return
		}
		// pivot == lapack.Bottom.
		if direct == lapack.Forward {
			for j := 0; j < m-1; j++ {
				ctmp := complex(c[j], 0)
				stmp := complex(s[j], 0)
				if c[j] != 1 || s[j] != 0 {
					for i := 0; i < n; i++ {
						tmp := a[j*lda+i]
						tmp2 := a[(m-1)*lda+i]
						a[j*lda+i] = stmp*tmp2 + ctmp*tmp
						a[(m-1)*lda+i] = ctmp*tmp2 - stmp*tmp
					}
				}
			}
			return
		}
		for j := m - 2; j >= 0; j-- {
			ctmp := complex(c[j], 0)
			stmp := complex(s[j], 0)
			if c[j] != 1 || s[j] != 0 {
				for i := 0; i < n; i++ {
					tmp := a[j*lda+i]
					tmp2 := a[(m-1)*lda+i]
					a[j*lda+i] = stmp*tmp2 + ctmp*tmp
					a[(m-1)*lda+i] = ctmp*tmp2 - stmp*tmp
				}
			}
		}
		return
	}
	// side == blas.Right.
	if pivot == lapack.Variable {
		if direct == lapack.Forward {
			for j := 0; j < n-1; j++ {
				ctmp := complex(c[j], 0)
				stmp := complex(s[j], 0)
				if c[j] != 1 || s[j] != 0 {
					for i := 0; i < m; i++ {
						tmp := a[i*lda+j+1]
						tmp2 := a[i*lda+j]
						a[i*lda+j+1] = ctmp*tmp - stmp*tmp2
						a[i*lda+j] = stmp*tmp + ctmp*tmp2
					}
				}
			}
			return
		}
		for j := n - 2; j >= 0; j-- {
			ctmp := complex(c[j], 0)
			stmp := complex(s[j], 0)
			if c[j] != 1 || s[j] != 0 {
				for i := 0; i < m; i++ {
					tmp := a[i*lda+j+1]
					tmp2 := a[i*lda+j]
					a[i*lda+j+1] = ctmp*tmp - stmp*tmp2
					a[i*lda+j] = stmp*tmp + ctmp*tmp2
				}
			}
		}
		return
	} else if pivot == lapack.Top {
		if direct == lapack.Forward {
			for j := 1; j < n; j++ {
				ctmp := complex(c[j-1], 0)
				stmp := complex(s[j-1], 0)
				if c[j-1] != 1 || s[j-1] != 0 {
					for i := 0; i < m; i++ {
						tmp := a[i*lda+j]
						tmp2 := a[i*lda]
						a[i*lda+j] = ctmp*tmp - stmp*tmp2
						a[i*lda] = stmp*tmp + ctmp*tmp2
					}
				}
			}
			return
		}
		for j := n - 1; j >= 1; j-- {
			ctmp := complex(c[j-1], 0)
			stmp := complex(s[j-1], 0)
			if c[j-1] != 1 || s[j-1] != 0 {
				for i := 0; i < m; i++ {
					tmp := a[i*lda+j]
					tmp2 := a[i*lda]
					a[i*lda+j] = ctmp*tmp - stmp*tmp2
					a[i*lda] = stmp*tmp + ctmp*tmp2
				}
			}
		}
		return
	}
	// pivot == lapack.Bottom.
	if direct == lapack.Forward {
		for j := 0; j < n-1; j++ {
			ctmp := complex(c[j], 0)
			stmp := complex(s[j], 0)
			if c[j] != 1 || s[j] != 0 {
				for i := 0; i < m; i++ {
					tmp := a[i*lda+j]
					tmp2 := a[i*lda+n-1]
					a[i*lda+j] = stmp*tmp2 + ctmp*tmp
					a[i*lda+n-1] = ctmp*tmp2 - stmp*tmp
				}
			}
		}
		return
	}
	for j := n - 2; j >= 0; j-- {
		ctmp := complex(c[j], 0)
		stmp := complex(s[j], 0)
		if c[j] != 1 || s[j] != 0 {
			for i := 0; i < m; i++ {
				tmp := a[i*lda+j]
				tmp2 := a[i*lda+n-1]
				a[i*lda+j] = stmp*tmp2 + ctmp*tmp
				a[i*lda+n-1] = ctmp*tmp2 - stmp*tmp
			}
		}
	}
}
