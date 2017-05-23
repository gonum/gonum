package cgo

import (
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/testblas"
)

const (
	Sm  = testblas.SmallMat
	Med = testblas.MediumMat
	Lg  = testblas.LargeMat
	Hg  = testblas.HugeMat
)

const (
	T  = blas.Trans
	NT = blas.NoTrans
)
