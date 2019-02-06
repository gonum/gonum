#!/bin/bash

go generate gonum.org/v1/gonum/blas
go generate gonum.org/v1/gonum/blas/gonum
go generate gonum.org/v1/gonum/unit
go generate gonum.org/v1/gonum/graph/formats/dot
if [ -n "$(git diff)" ]; then
	git checkout -- go.mod # Discard changes to go.mod that have been made.
	git diff
	exit 1
fi
