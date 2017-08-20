# Gonum [![Build Status](https://travis-ci.org/gonum/gonum.svg?branch=master)](https://travis-ci.org/gonum/gonum) [![Coverage Status](https://coveralls.io/repos/gonum/gonum/badge.svg?branch=master&service=github)](https://coveralls.io/github/gonum/gonum?branch=master) [![GoDoc](https://godoc.org/gonum.org/v1/gonum?status.svg)](https://godoc.org/gonum.org/v1/gonum) [![Go Report Card](https://goreportcard.com/badge/github.com/gonum/gonum)](https://goreportcard.com/report/github.com/gonum/gonum)

## Installation

The core packages of the gonum suite are written in pure Go with some assembly.
Installation is done using `go get`.
```
go get -u gonum.org/v1/gonum/...
```

## Build tags

The gonum packages use a variety of build tags to set non-standard build conditions.
Building gonum applications will work without knowing how to use these tags, but they can be used during testing and to control the use of assembly and CGO code.

The current list of non-internal tags is as follows:

- appengine — do not use assembly or unsafe
- bounds — use bounds checks even in internal calls
- cblas — use CGO gonum.org/v1/netlib/blas/netlib BLAS implementation in tests (only in [mat package](https://godoc.org/gonum.org/v1/gonum/mat))
- go1.7 — use go1.7 style sub tests and benchmarks where implemented
- noasm — do not use assembly implementations
- tomita — use [Tomita, Tanaka, Takahashi pivot choice](https://doi.org/10.1016%2Fj.tcs.2006.06.015) for maximimal clique calculation, otherwise use random pivot (only in [topo package](https://godoc.org/gonum.org/v1/gonum/graph/topo))


## Issues

If you find any bugs, feel free to file an issue on the github issue tracker. Discussions on API changes, added features, code review, or similar requests are preferred on the gonum-dev Google Group.

https://groups.google.com/forum/#!forum/gonum-dev

