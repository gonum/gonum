#!/bin/bash

go generate github.com/gonum/internal/asm
if [ -n "$(git diff)" ]; then
	exit 1
fi
