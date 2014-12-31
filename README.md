# Gonum BLAS [![Build Status](https://travis-ci.org/gonum/blas.svg)](https://travis-ci.org/gonum/blas)  [![Coverage Status](https://img.shields.io/coveralls/gonum/blas.svg)](https://coveralls.io/r/gonum/blas)

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
  cd $GOPATH/src/github.com/gonum/blas/cblas
  CGO_LDFLAGS="-L/path/to/OpenBLAS -lopenblas" go install 
```

For Windows you can download binary packages for OpenBLAS at
[SourceForge](http://sourceforge.net/projects/openblas/files/).

If you want to use a different BLAS package such as the Intel MKL you can
adjust the `CGO_LDFLAGS` variable:
```sh
  cd $GOPATH/src/github.com/gonum/blas/cblas
  CGO_LDFLAGS="-lmkl_rt" go install
```

On OS X the easiest solution is to use the libraries provided by the system:
```sh
  cd $GOPATH/src/github.com/gonum/blas/cblas
  CGO_LDFLAGS="-framework Accelerate" go install
```

## Packages

### blas

Defines [BLAS API](www.netlib.org/blas/blast-forum/cinterface.pdf) split in several interfaces

### blas/goblas

Go implementation of the BLAS API (incomplete, implements most of the float64 API)

### blas/cblas

Binding to a C implementation of the cblas interface (e.g. ATLAS, OpenBLAS, Intel MKL)

The recommended (free) option for good performance on both Linux and Darwin is OpenBLAS.

### blas/dbw

Wrapper for an implementation of the double precision real (i.e., `float64`) part
of the blas API

You have to register an implementation before you can use the BLAS functions:

```Go
package main

import (
	"fmt"

	"github.com/gonum/blas/cblas"
	"github.com/gonum/blas/dbw"
)

func init() {
	dbw.Register(cblas.Blas{})
}

func main() {
	v := dbw.NewVector([]float64{1, 1, 1})
	fmt.Println("v has length:", dbw.Nrm2(v))
}
```

### blas/zbw

Wrapper for an implementation of the double precision complex (i.e., `complex128`)
part of the blas API

## Issues

If you find any bugs, feel free to file an issue on the github issue tracker.
Discussions on API changes, added features, code review, or similar requests
are preferred on the [gonum-dev Google Group](https://groups.google.com/forum/#!forum/gonum-dev).

## License

Please see [github.com/gonum/license](https://github.com/gonum/license) for general
license information, contributors, authors, etc on the Gonum suite of packages.
