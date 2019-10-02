# Gonum [![Build Status](https://www.travis-ci.org/gonum/gonum.svg?branch=master)](https://www.travis-ci.org/gonum/gonum/branches) [![Build status](https://ci.appveyor.com/api/projects/status/valslkp8sr50eepn/branch/master?svg=true)](https://ci.appveyor.com/project/Gonum/gonum/branch/master) [![codecov.io](https://codecov.io/gh/gonum/gonum/branch/master/graph/badge.svg)](https://codecov.io/gh/gonum/gonum) [![coveralls.io](https://coveralls.io/repos/gonum/gonum/badge.svg?branch=master&service=github)](https://coveralls.io/github/gonum/gonum?branch=master) [![GoDoc](https://godoc.org/gonum.org/v1/gonum?status.svg)](https://godoc.org/gonum.org/v1/gonum) [![Go Report Card](https://goreportcard.com/badge/github.com/gonum/gonum)](https://goreportcard.com/report/github.com/gonum/gonum) [![stability-unstable](https://img.shields.io/badge/stability-unstable-yellow.svg)](https://github.com/emersion/stability-badges#unstable)

## Installation

The core packages of the Gonum suite are written in pure Go with some assembly.
Installation is done using `go get`.
```
go get -u gonum.org/v1/gonum/...
```

## Supported Go versions

Gonum supports and tests on the three most recent minor versions of Go on [Linux](https://github.com/gonum/gonum/blob/master/.travis.yml#L6-L11) and [Windows](https://github.com/gonum/gonum/blob/master/appveyor.yml#L13-L18).

## Release schedule

The Gonum modules are released on a quaterly-based time schedule, aligned with the Go releases.
_i.e.:_ when `Go-1.x` is released, `Gonum-v0.n.0` is released around the same time.
Three months later, `Gonum-v0.n+1.0` is released.
Three months after, `Go-1.x+1` is released, and `Gonum-v0.n+2.0` as well.

The release schedule, based on the current Go release schedule is thus:

- `Gonum-v0.n.0`: February
- `Gonum-v0.n+1.0`: May
- `Gonum-v0.n+2.0`: August
- `Gonum-v0.n+3.0`: November

## Build tags

The Gonum packages use a variety of build tags to set non-standard build conditions.
Building Gonum applications will work without knowing how to use these tags, but they can be used during testing and to control the use of assembly and CGO code.

The current list of non-internal tags is as follows:

- appengine — do not use assembly or unsafe
- safe — synonym for appengine
- bounds — use bounds checks even in internal calls
- cblas — use CGO gonum.org/v1/netlib/blas/netlib BLAS implementation in tests (only in [mat package](https://godoc.org/gonum.org/v1/gonum/mat))
- noasm — do not use assembly implementations
- tomita — use [Tomita, Tanaka, Takahashi pivot choice](https://doi.org/10.1016%2Fj.tcs.2006.06.015) for maximimal clique calculation, otherwise use random pivot (only in [topo package](https://godoc.org/gonum.org/v1/gonum/graph/topo))


## Issues

If you find any bugs, feel free to file an issue on the github issue tracker. Discussions on API changes, added features, code review, or similar requests are preferred on the gonum-dev Google Group.

https://groups.google.com/forum/#!forum/gonum-dev

## License

Original code is licensed under the Gonum License found in the LICENSE file. Portions of the code are subject to the additional licenses found in THIRD_PARTY_LICENSES. All third party code is licensed either under a BSD or MIT license.

Code in graph/formats/dot is dual licensed [Public Domain Dedication](https://creativecommons.org/publicdomain/zero/1.0/) and Gonum License, and users are free to choose the license which suits their needs for this code.
