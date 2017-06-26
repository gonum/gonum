#!/bin/bash

go generate gonum.org/v1/gonum/blas/gonum
go generate gonum.org/v1/gonum/unit
if [ -n "$(git diff)" ]; then
	exit 1
fi
