# Gonum

[![Build status](https://github.com/gonum/gonum/workflows/CI/badge.svg)](https://github.com/gonum/gonum/actions)
[![Build status](https://ci.appveyor.com/api/projects/status/valslkp8sr50eepn/branch/master?svg=true)](https://ci.appveyor.com/project/Gonum/gonum/branch/master)
[![codecov.io](https://codecov.io/gh/gonum/gonum/branch/master/graph/badge.svg)](https://codecov.io/gh/gonum/gonum)
[![go.dev reference](https://pkg.go.dev/badge/gonum.org/v1/gonum)](https://pkg.go.dev/gonum.org/v1/gonum)
[![GoDoc](https://godoc.org/gonum.org/v1/gonum?status.svg)](https://godoc.org/gonum.org/v1/gonum)
[![Go Report Card](https://goreportcard.com/badge/github.com/gonum/gonum)](https://goreportcard.com/report/github.com/gonum/gonum)
[![stability-unstable](https://img.shields.io/badge/stability-unstable-yellow.svg)](https://github.com/emersion/stability-badges#unstable)

## Installation

The core packages of the Gonum suite are written in pure Go with some assembly.
Installation is done using `go get`.
```
go get -u gonum.org/v1/gonum/...
```

## Supported Go versions

Gonum supports and tests using the gc compiler on the [two most recent Go releases](https://github.com/gonum/gonum/blob/master/.github/workflows/ci.yml#L14-L15) on Linux (386, amd64 and arm64), macOS and Windows (both on amd64).

Note that floating point behavior may differ between compiler versions and between architectures due to differences in floating point operation implementations.

## Release schedule

The Gonum modules are released on a six-month release schedule, aligned with the Go releases.
_i.e.:_ when `Go-1.x` is released, `Gonum-v0.n.0` is released around the same time.
Six months after, `Go-1.x+1` is released, and `Gonum-v0.n+1.0` as well.

The release schedule, based on the current Go release schedule is thus:

- `Gonum-v0.n.0`: February
- `Gonum-v0.n+1.0`: August

## Build tags

The Gonum packages use a variety of build tags to set non-standard build conditions.
Building Gonum applications will work without knowing how to use these tags, but they can be used during testing and to control the use of assembly and CGO code.

The current list of non-internal tags is as follows:

- safe — do not use assembly or unsafe
- bounds — use bounds checks even in internal calls
- noasm — do not use assembly implementations
- tomita — use [Tomita, Tanaka, Takahashi pivot choice](https://doi.org/10.1016%2Fj.tcs.2006.06.015) for maximimal clique calculation, otherwise use random pivot (only in [topo package](https://godoc.org/gonum.org/v1/gonum/graph/topo))


## Issues [![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/gonum/gonum)](https://www.tickgit.com/browse?repo=github.com/gonum/gonum)

If you find any bugs, feel free to file an issue on the github issue tracker. Discussions on API changes, added features, code review, or similar requests are preferred on the gonum-dev Google Group.

https://groups.google.com/forum/#!forum/gonum-dev

## License

Original code is licensed under the Gonum License found in the LICENSE file. Portions of the code are subject to the additional licenses found in THIRD_PARTY_LICENSES. All third party code is licensed either under a BSD or MIT license.

Code in graph/formats/dot is dual licensed [Public Domain Dedication](https://creativecommons.org/publicdomain/zero/1.0/) and Gonum License, and users are free to choose the license which suits their needs for this code.

The W3C test suites in graph/formats/rdf are distributed under both the [W3C Test Suite License](http://www.w3.org/Consortium/Legal/2008/04-testsuite-license) and the [W3C 3-clause BSD License](http://www.w3.org/Consortium/Legal/2008/03-bsd-license).
