#!/bin/bash

go generate github.com/gonum/blas/cgo
go generate github.com/gonum/blas/native
echo "$(git diff -- . ':!.travis')"
if [ -n "$(git diff -- . ':!.travis')" ]; then
	exit 1
fi
