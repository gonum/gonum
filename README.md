# Gonum BLAS [![Build Status](https://travis-ci.org/gonum/blas.svg?branch=master)](https://travis-ci.org/gonum/blas)  [![Coverage Status](https://img.shields.io/coveralls/gonum/blas.svg)](https://coveralls.io/r/gonum/blas)

A collection of packages to provide BLAS functionality for the [Go programming
language](http://golang.org)

## Installation
```sh
  go get github.com/gonum/blas
```

### BLAS C-bindings

If you want to use OpenBLAS, install it in any directory:
```sh
  git clone https://github.com/xianyi/OpenBLAS
  cd OpenBLAS
  make
```

Then install the cblas package:
```sh
  cd $GOPATH/src/github.com/gonum/blas/cgo
  CGO_LDFLAGS="-L/path/to/OpenBLAS -lopenblas" go install
```

For Windows you can download binary packages for OpenBLAS at
[SourceForge](http://sourceforge.net/projects/openblas/files/).

If you want to use a different BLAS package such as the Intel MKL you can
adjust the `CGO_LDFLAGS` variable:
```sh
  cd $GOPATH/src/github.com/gonum/blas/cgo
  CGO_LDFLAGS="-lmkl_rt" go install
```

On OS X the easiest solution is to use the libraries provided by the system:
```sh
  cd $GOPATH/src/github.com/gonum/blas/cgo
  CGO_LDFLAGS="-framework Accelerate" go install
```

## Packages

### blas

Defines [BLAS API](http://www.netlib.org/blas/blast-forum/cinterface.pdf) split in several interfaces

### blas/native

Go implementation of the BLAS API (incomplete, implements the float64 API)

### blas/cgo

Binding to a C implementation of the cblas interface (e.g. ATLAS, OpenBLAS, Intel MKL)

The recommended (free) option for good performance on both Linux and Darwin is OpenBLAS.

### blas/blas64

Wrapper for an implementation of the double precision real (i.e., `float64`) part
of the blas API

```Go
package main

import (
	"fmt"

	"github.com/gonum/blas/blas64"
)

func main() {
	v := blas64.Vector{Inc: 1, Data: []float64{1, 1, 1}}
	fmt.Println("v has length:", blas64.Nrm2(len(v.Data), v))
}
```

### blas/cblas128

Wrapper for an implementation of the double precision complex (i.e., `complex128`)
part of the blas API

Currently blas/cblas128 requires blas/cgo and must be built with the `cblas` build tag,
assuming blas/cgo is installed:

```sh
  go install -tags cblas github.com/gonum/blas/cblas128
```


## Issues

If you find any bugs, feel free to file an issue on the github issue tracker.
Discussions on API changes, added features, code review, or similar requests
are preferred on the [gonum-dev Google Group](https://groups.google.com/forum/#!forum/gonum-dev).

## License

Please see [github.com/gonum/license](https://github.com/gonum/license) for general
license information, contributors, authors, etc on the Gonum suite of packages.
