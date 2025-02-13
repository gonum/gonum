#!/bin/bash

set -e
go tool gonum.org/v1/tools/cmd/check-imports -b "math/rand,github.com/gonum/.*"
