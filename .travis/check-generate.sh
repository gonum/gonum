#!/bin/bash

go generate gonum.org/v1/gonum/matrix
go generate gonum.org/v1/gonum/blas/native
go generate gonum.org/v1/gonum/lapack/native
if [ -n "$(git diff)" ]; then
	exit 1
fi
