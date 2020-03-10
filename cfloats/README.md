# Gonum cfloats [![GoDoc](https://godoc.org/gonum.org/v1/gonum/cfloats?status.svg)](https://godoc.org/gonum.org/v1/gonum/cfloats)

Package cfloats provides a set of helper routines for dealing with slices of complex128.
The functions avoid allocations to allow for use within tight loops without garbage collection overhead.
