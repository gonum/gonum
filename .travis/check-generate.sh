#!/bin/bash
set -ex

go generate github.com/gonum/blas/cgo
go generate github.com/gonum/blas/native
if [ -n "$(git diff)" ]; then
	exit 1
fi
