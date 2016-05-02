#!/bin/bash

go generate github.com/gonum/blas/cgo
go generate github.com/gonum/blas/native
if [ -n "$(git diff -- . ':!.travis')" ]; then
	exit 1
fi
