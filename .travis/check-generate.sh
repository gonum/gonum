#!/bin/bash

go generate gonum.org/v1/gonum/blas/gonum
go generate gonum.org/v1/gonum/unit
go generate gonum.org/v1/gonum/graph/formats/dot
if [ -n "$(git diff)" ]; then
	exit 1
fi
