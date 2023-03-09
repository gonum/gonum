# Gonum cmplxs

[![go.dev reference](https://pkg.go.dev/badge/gonum.org/v1/gonum/cmplxs)](https://pkg.go.dev/gonum.org/v1/gonum/cmplxs)
[![GoDoc](https://godocs.io/gonum.org/v1/gonum/cmplxs?status.svg)](https://godocs.io/gonum.org/v1/gonum/cmplxs)

Package cmplxs provides a set of helper routines for dealing with slices of complex128.
The functions avoid allocations to allow for use within tight loops without garbage collection overhead.
