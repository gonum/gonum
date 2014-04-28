# Gonum BLAS [![Build Status](https://travis-ci.org/gonum/blas.png?branch=master)](https://travis-ci.org/gonum/blas)

A collection of packages to provide BLAS functionality for the Go programming
language (http://golang.org)

## Installation 

```
  go get github.com/gonum/blas
```

For the BLAS C-bindings:

If you want to use OpenBLAS install it in any directory
```
  git clone https://github.com/xianyi/OpenBLAS
  cd OpenBLAS
  make
```

Then install the cblas package
```
  cd $GOPATH/src/github.com/gonum/blas/cblas
  CGO_LDFLAGS="-L/path/to/OpenBLAS -lopenblas" go install 
```

For Windows you can download binary packages for OpenBLAS at
http://sourceforge.net/projects/openblas/files/

If you want to use a different BLAS package such as the Intel MKL you can adjust the CGO_LDFLAGS variable, e.g.
```
  cd $GOPATH/src/github.com/gonum/blas/cblas
  CGO_LDFLAGS="-lmkl_rt" go install 
```

On OS X you can also use the libraries provieded by the system:
```
  cd $GOPATH/src/github.com/gonum/blas/cblas
  CGO_LDFLAGS="LDFLAGS: -DYA_BLAS -DYA_LAPACK -DYA_BLASMULT -framework vecLib" go install 
```



## Packages

### blas

Defines BLAS API (www.netlib.org/blas/blast-forum/cinterface.pdf) split in several interfaces

### blas/goblas

Go implementation of the BLAS API (incomplete, implements most of the float64 API)

### blas/cblas

Binding to a C implementation of the cblas interface (e.g. ATLAS, OpenBLAS, intel MKL)

The linker flags (i.e. path to the BLAS library and library name) might have to be adapted.

The recommended (free) option for good performance on both linux and darwin is OpenBLAS.

### blas/dbw

Wrapper for an implementation of the double precision real (i.e. float64) part of the blas API

You have to register an implementation before you can use the BLAS functions:

```
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

Wrapper for an implementation of the double precision complex (i.e. complex128) part of the blas API
