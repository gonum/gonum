#!/bin/bash

go generate gonum.org/v1/gonum/blas/native
if [ -n "$(git diff)" ]; then
	exit 1
fi
