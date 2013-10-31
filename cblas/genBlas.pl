#!/usr/bin/env perl
# Copyright ©2012 The bíogo.blas Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

use strict;
use warnings;

my $cblasHeader = "cblas.h";
my $LIB = "/usr/lib/";

my $excludeComplex = 0;
my $excludeAtlas = 1;


open(my $cblas, "<", $cblasHeader) or die;
open(my $goblas, ">", "blas.go") or die;

my %done = ("cblas_errprn"     => 1,
	        "cblas_srotg"      => 1,
	        "cblas_srotmg"     => 1,
	        "cblas_srotm"      => 1,
	        "cblas_drotg"      => 1,
	        "cblas_drotmg"     => 1,
	        "cblas_drotm"      => 1,
	        "cblas_crotg"      => 1,
	        "cblas_zrotg"      => 1,
	        "cblas_cdotu_sub"  => 1,
	        "cblas_cdotc_sub"  => 1,
	        "cblas_zdotu_sub"  => 1,
	        "cblas_zdotc_sub"  => 1,
	        );

my $atlas = "";
if ($excludeAtlas) {
	$done{'cblas_csrot'} = 1;
	$done{'cblas_zdrot'} = 1;
} else {
	$atlas = " -latlas";
}
printf $goblas <<EOH;
// Do not manually edit this file. It was created by the genBlas.pl script from ${cblasHeader}.

// Copyright ©2012 The bíogo.blas Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cblas implements the blas interfaces.
package cblas

/*
#cgo CFLAGS: -g -O2
#cgo linux LDFLAGS: -L/usr/lib/ -lcblas
#cgo darwin LDFLAGS: -DYA_BLAS -DYA_LAPACK -DYA_BLASMULT -framework vecLib
#include "${cblasHeader}"
*/
import "C"

import (
	"github.com/gonum/blas"
	"unsafe"
)

// Type check assertions:
var (
	_ blas.Float32    = Blas{}
	_ blas.Float64    = Blas{}
	_ blas.Complex64  = Blas{}
	_ blas.Complex128 = Blas{}
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type Blas struct{}

// Special cases...

func (Blas) Srotg(a float32, b float32) (c float32, s float32, r float32, z float32) {
	C.cblas_srotg((*C.float)(&a), (*C.float)(&b), (*C.float)(&c), (*C.float)(&s))
	return c, s, a, b
}
func (Blas) Srotmg(d1 float32, d2 float32, b1 float32, b2 float32) (p *blas.SrotmParams, rd1 float32, rd2 float32, rb1 float32) {
	p = &blas.SrotmParams{}
	C.cblas_srotmg((*C.float)(&d1), (*C.float)(&d2), (*C.float)(&b1), C.float(b2), (*C.float)(unsafe.Pointer(p)))
	return p, d1, d2, b1
}
func (Blas) Srotm(n int, x []float32, incX int, y []float32, incY int, p *blas.SrotmParams) {
	if n < 0 {
		panic("cblas: n < 0")
	}
	if n*incX > len(x) {
		panic("cblas: index out of range")
	}
	if n*incY > len(y) {
		panic("cblas: index out of range")
	}
	C.cblas_srotm(C.int(n), (*C.float)(&x[0]), C.int(incX), (*C.float)(&y[0]), C.int(incY), (*C.float)(unsafe.Pointer(p)))
}
func (Blas) Drotg(a float64, b float64) (c float64, s float64, r float64, z float64) {
	C.cblas_drotg((*C.double)(&a), (*C.double)(&b), (*C.double)(&c), (*C.double)(&s))
	return c, s, a, b
}
func (Blas) Drotmg(d1 float64, d2 float64, b1 float64, b2 float64) (p *blas.DrotmParams, rd1 float64, rd2 float64, rb1 float64) {
	p = &blas.DrotmParams{}
	C.cblas_drotmg((*C.double)(&d1), (*C.double)(&d2), (*C.double)(&b1), C.double(b2), (*C.double)(unsafe.Pointer(p)))
	return p, d1, d2, b1
}
func (Blas) Drotm(n int, x []float64, incX int, y []float64, incY int, p *blas.DrotmParams) {
	if n < 0 {
		panic("cblas: n < 0")
	}
	if n*incX > len(x) {
		panic("cblas: index out of range")
	}
	if n*incY > len(y) {
		panic("cblas: index out of range")
	}
	C.cblas_drotm(C.int(n), (*C.double)(&x[0]), C.int(incX), (*C.double)(&y[0]), C.int(incY), (*C.double)(unsafe.Pointer(p)))
}
func (Blas) Cdotu(n int, x []complex64, incX int, y []complex64, incY int) (dotu complex64) {
	if n < 0 {
		panic("cblas: n < 0")
	}
	if incX <= 0 || n*incX > len(x) {
		panic("cblas: index out of range")
	}
	if incY <= 0 || n*incY > len(y) {
		panic("cblas: index out of range")
	}
	C.cblas_cdotu_sub(C.int(n), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY), unsafe.Pointer(&dotu))
	return dotu
}
func (Blas) Cdotc(n int, x []complex64, incX int, y []complex64, incY int) (dotc complex64) {
	if n < 0 {
		panic("cblas: n < 0")
	}
	if incX <= 0 || n*incX > len(x) {
		panic("cblas: index out of range")
	}
	if incY <= 0 || n*incY > len(y) {
		panic("cblas: index out of range")
	}
	C.cblas_cdotc_sub(C.int(n), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY), unsafe.Pointer(&dotc))
	return dotc
}
func (Blas) Zdotu(n int, x []complex128, incX int, y []complex128, incY int) (dotu complex128) {
	if n < 0 {
		panic("cblas: n < 0")
	}
	if incX <= 0 || n*incX > len(x) {
		panic("cblas: index out of range")
	}
	if incY <= 0 || n*incY > len(y) {
		panic("cblas: index out of range")
	}
	C.cblas_zdotu_sub(C.int(n), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY), unsafe.Pointer(&dotu))
	return dotu
}
func (Blas) Zdotc(n int, x []complex128, incX int, y []complex128, incY int) (dotc complex128) {
	if n < 0 {
		panic("cblas: n < 0")
	}
	if incX <= 0 || n*incX > len(x) {
		panic("cblas: index out of range")
	}
	if incY <= 0 || n*incY > len(y) {
		panic("cblas: index out of range")
	}
	C.cblas_zdotc_sub(C.int(n), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY), unsafe.Pointer(&dotc))
	return dotc
}
EOH

printf $goblas <<EOH unless $excludeAtlas;
func (Blas) Crotg(a complex64, b complex64) (c complex64, s complex64, r complex64, z complex64) {
	C.cblas_srotg(unsafe.Pointer(&a), unsafe.Pointer(&b), unsafe.Pointer(&c), unsafe.Pointer(&s))
	return c, s, a, b
}
func (Blas) Zrotg(a complex128, b complex128) (c complex128, s complex128, r complex128, z complex128) {
	C.cblas_drotg(unsafe.Pointer(&a), unsafe.Pointer(&b), unsafe.Pointer(&c), unsafe.Pointer(&s))
	return c, s, a, b
}
EOH

print $goblas "\n";

$/ = undef;
my $header = <$cblas>;

# horrible munging of text...
$header =~ s/#[^\n\r]*//g;                 # delete cpp lines
$header =~ s/\n +([^\n\r]*)/\n$1/g;        # remove starting space
$header =~ s/(?:\n ?\n)+/\n/g;             # delete empty lines
$header =~ s! ((['"]) (?: \\. | .)*? \2) | # skip quoted strings
             /\* .*? \*/ |                 # delete C comments
             // [^\n\r]*                   # delete C++ comments just in case
             ! $1 || ' '                   # change comments to a single space
             !xseg;    	                   # ignore white space, treat as single line
                                           # evaluate result, repeat globally
$header =~ s/([^;])\n/$1/g;                # join prototypes into single lines
$header =~ s/, +/,/g;
$header =~ s/ +/ /g;
$header =~ s/ +}/}/g;
$header =~ s/\n+//;

$/ = "\n";
my @lines = split ";\n", $header;

our %retConv = (
	"int" => "int ",
	"float" => "float32 ",
	"double" => "float64 ",
	"CBLAS_INDEX" => "int ",
	"void" => ""
);

foreach my $line (@lines) {
	process($line);
}

close($goblas);
`go fmt .`;

sub process {
	my $line = shift;
	chomp $line;
	if (not $line =~ m/^enum/) {
		processProto($line);
	}
}

sub processProto {
	my $proto = shift;
	my ($func, $paramList) = split /[()]/, $proto;
	(my $ret, $func) = split ' ', $func;
	if ($done{$func} or $excludeComplex && $func =~ m/_[isd]?[zc]/ or $excludeAtlas && $func =~ m/^catlas_/) {
		return
	}
	$done{$func} = 1;
	my $GoRet = $retConv{$ret};
	my $complexType = $func;
	$complexType =~ s/.*_[isd]?([zc]).*/$1/;
	print $goblas "func (Blas) ".Gofunc($func)."(".processParamToGo($func, $paramList, $complexType).") ".$GoRet."{\n";
	print $goblas processParamToChecks($func, $paramList);
	print $goblas "\t";
	if ($ret ne 'void') {
		chop($GoRet);
		print $goblas "return ".$GoRet."(";
	}
	print $goblas "C.$func(".processParamToC($func, $paramList).")";
	if ($ret ne 'void') {
		print $goblas ")";
	}
	print $goblas "\n}\n";
}

sub Gofunc {
	my $fnName = shift;
	$fnName =~ s/_sub//;
	my ($pack, $func, $tail) = split '_', $fnName;
	if ($pack eq 'cblas') {
		$pack = "";
	} else {
		$pack = substr $pack, 1;
	}

	return ucfirst $pack . ucfirst $func . ucfirst $tail if $tail;
	return ucfirst $pack . ucfirst $func;
}

sub processParamToGo {
	my $func = shift;
	my $paramList = shift;
	my $complexType = shift;
	my @processed;
	my @params = split ',', $paramList;
	foreach my $param (@params) {
		my @parts = split /[ *]/, $param;
		my $var = lcfirst $parts[scalar @parts - 1];
		$param =~ m/^(?:const )?int/ && do {
			push @processed, $var." int"; next;
		};
		$param =~ m/^(?:const )?void/ && do {
			my $type;
			if ($var eq "alpha" || $var eq "beta") {
				$type = " ";
			} else {
				$type = " []";
			}
			if ($complexType eq 'c') {
				push @processed, $var.$type."complex64"; next;
			} elsif ($complexType eq 'z') {
				push @processed, $var.$type."complex128"; next;
			} else {
				die "unexpected complex type for '$func' - '$complexType'";
			}
		};
		$param =~ m/^(?:const )?char \*/ && do {
			push @processed, $var." *byte"; next;
		};
		$param =~ m/^(?:const )?float \*/ && do {
			push @processed, $var." []float32"; next;
		};
		$param =~ m/^(?:const )?double \*/ && do {
			push @processed, $var." []float64"; next;
		};
		$param =~ m/^(?:const )?float/ && do {
			push @processed, $var." float32"; next;
		};
		$param =~ m/^(?:const )?double/ && do {
			push @processed, $var." float64"; next;
		};
		$param =~ m/^const enum/ && do {
			$var eq "order" && do {
				$var = "o";
				push @processed, $var." blas.Order"; next;
			};
			$var =~ /trans/ && do {
				$var =~ s/trans([AB]?)/t$1/;
				push @processed, $var." blas.Transpose"; next;
			};
			$var eq "uplo" && do {
				$var = "ul";
				push @processed, $var." blas.Uplo"; next;
			};
			$var eq "diag" && do {
				$var = "d";
				push @processed, $var." blas.Diag"; next;
			};
			$var eq "side" && do {
				$var = "s";
				push @processed, $var." blas.Side"; next;
			};
		};
	}
	die "missed Go parameters from '$func', '$paramList'" if scalar @processed != scalar @params;
	return join ", ", @processed;
}

sub processParamToChecks {
	my $func = shift;
	my $paramList = shift;
	my @processed;
	my @params = split ',', $paramList;
	my %arrayArgs;
	my %scalarArgs;
	foreach my $param (@params) {
		my @parts = split /[ *]/, $param;
		my $var = lcfirst $parts[scalar @parts - 1];
		$param =~ m/^(?:const )?int \*[a-zA-Z]/ && do {
			$scalarArgs{$var} = 1; next;
		};
		$param =~ m/^(?:const )?void \*[a-zA-Z]/ && do {
			if ($var ne "alpha" && $var ne "beta") {
				$arrayArgs{$var} = 1;
			}
			next;
		};
		$param =~ m/^(?:const )?(?:float|double) \*[a-zA-Z]/ && do {
			$arrayArgs{$var} = 1; next;
		};
		$param =~ m/^(?:const )?(?:int|float|double) [a-zA-Z]/ && do {
			$scalarArgs{$var} = 1; next;
		};
		$param =~ m/^const enum [a-zA-Z]/ && do {
			$var eq "order" && do {
				$scalarArgs{'o'} = 1;
				push @processed, "if o != blas.RowMajor && o != blas.ColMajor { panic(\"cblas: illegal order\") }"; next;
			};
			$var =~ /trans/ && do {
				$var =~ s/trans([AB]?)/t$1/;
				$scalarArgs{$var} = 1; next;
				if ($func =~ m/cblas_[cz]h/) {
					push @processed, "if $var != blas.NoTrans && $var != blas.ConjTrans { panic(\"cblas: illegal transpose\") }"; next;
				} elsif ($func =~ m/cblas_[cz]s/) {
					push @processed, "if $var != blas.NoTrans && $var != blas.Trans { panic(\"cblas: illegal transpose\") }"; next;
				} else {
					push @processed, "if $var != blas.NoTrans && $var != blas.Trans && $var != blas.ConjTrans { panic(\"cblas: illegal transpose\") }"; next;
				}
			};
			$var eq "uplo" && do {
				push @processed, "if ul != blas.Upper && ul != blas.Lower { panic(\"cblas: illegal triangle\") }"; next;
			};
			$var eq "diag" && do {
				push @processed, "if d != blas.NonUnit && d != blas.Unit { panic(\"cblas: illegal diagonal\") }"; next;
			};
			$var eq "side" && do {
				$scalarArgs{'s'} = 1;
				push @processed, "if s != blas.Left && s != blas.Right { panic(\"cblas: illegal side\") }"; next;
			};
		};
	}

	# shape checks
	foreach my $ref ('m', 'n', 'k', 'kL', 'kU') {
		push @processed, "if $ref < 0 { panic(\"cblas: $ref < 0\") }" if $scalarArgs{$ref};
	}

	if ($arrayArgs{'ap'}) {
		push @processed, "if n*(n + 1)/2 > len(ap) { panic(\"cblas: index of ap out of range\") }"
	}

	if ($func =~ m/cblas_[sdcz]g[eb]mv/) {
		push @processed, "if incX <= 0 || incY <= 0 { panic(\"cblas: index increment out of range\") }";
		push @processed, "var lenX, lenY int";
		push @processed, "if tA == blas.NoTrans { lenX, lenY = n, m } else { lenX, lenY = m, n }";
		push @processed, "if (lenX-1)*incX > len(x) { panic(\"cblas: index of x out of range\") }";
		push @processed, "if (lenY-1)*incY > len(y) { panic(\"cblas: index of y out of range\") }";
	} elsif ($scalarArgs{'m'}) {
		push @processed, "if incX <= 0 || (m-1)*incX > len(x) { panic(\"cblas: index of x out of range\") }" if $scalarArgs{'incX'};
		push @processed, "if incY <= 0 || (n-1)*incY > len(y) { panic(\"cblas: index of y out of range\") }" if $scalarArgs{'incY'};
	} else {
		push @processed, "if incX <= 0 || (n-1)*incX > len(x) { panic(\"cblas: index of x out of range\") }" if $scalarArgs{'incX'};
		push @processed, "if incY <= 0 || (n-1)*incY > len(y) { panic(\"cblas: index of y out of range\") }" if $scalarArgs{'incY'};
	}

	if (not $func =~ m/(?:mm|r2?k)$/) {
		if ($arrayArgs{'a'}) {
			if (($scalarArgs{'kL'} && $scalarArgs{'kU'}) || $scalarArgs{'m'}) {
				push @processed, "if o == blas.RowMajor {";
				if ($scalarArgs{'kL'} && $scalarArgs{'kU'}) {
					push @processed, "if lda*(m-1)+n > len(a) || lda < kL+kU+1 { panic(\"cblas: index of a out of range\") }";
				} else {
					push @processed, "if lda*(m-1)+n > len(a) || lda < max(1, n) { panic(\"cblas: index of a out of range\") }";
				}
				push @processed, "} else {";
				if ($scalarArgs{'kL'} && $scalarArgs{'kU'}) {
					push @processed, "if lda*(n-1)+m > len(a) || lda < kL+kU+1 { panic(\"cblas: index of a out of range\") }";
				} else {
					push @processed, "if lda*(n-1)+m > len(a) || lda < max(1, m) { panic(\"cblas: index of a out of range\") }";
				}
				push @processed, "}";
			} else {
				if ($scalarArgs{'k'}) {
					push @processed, "if lda*(n-1)+n > len(a) || lda < k+1 { panic(\"cblas: index of a out of range\") }";
				} else {
					push @processed, "if lda*(n-1)+n > len(a) || lda < max(1, n) { panic(\"cblas: index of a out of range\") }";
				}
			}
		}
	} else {
		if ($scalarArgs{'s'}) {
			push @processed, "var k int";
			push @processed, "if s == blas.Left { k = m } else { k = n }";
			push @processed, "if lda*(k-1)+k > len(a) || lda < max(1, k) { panic(\"cblas: index of a out of range\") }";
			push @processed, "if o == blas.RowMajor {";
			push @processed, "if ldb*(m-1)+n > len(b) || ldb < max(1, n) { panic(\"cblas: index of b out of range\") }";
			if ($arrayArgs{'c'}) {
				push @processed, "if ldc*(m-1)+n > len(c) || ldc < max(1, n) { panic(\"cblas: index of c out of range\") }";
			}
			push @processed, "} else {";
			push @processed, "if ldb*(n-1)+m > len(b) || ldb < max(1, m) { panic(\"cblas: index of b out of range\") }";
			if ($arrayArgs{'c'}) {
				push @processed, "if ldc*(n-1)+m > len(c) || ldc < max(1, m) { panic(\"cblas: index of c out of range\") }";
			}
			push @processed, "}";
		}
		if ($scalarArgs{'t'}) {
			push @processed, "var row, col int";
			push @processed, "if t == blas.NoTrans { row, col = n, k } else { row, col = k, n }";
			push @processed, "if o == blas.RowMajor {";
			foreach my $ref ('a', 'b') {
				if ($arrayArgs{$ref}) {
					push @processed, "if ld${ref}*(row-1)+col > len(${ref}) || ld${ref} < max(1, col) { panic(\"cblas: index of ${ref} out of range\") }";
				}
			}
			push @processed, "} else {";
			foreach my $ref ('a', 'b') {
				if ($arrayArgs{$ref}) {
					push @processed, "if ld${ref}*(col-1)+row > len(${ref}) || ld${ref} < max(1, row) { panic(\"cblas: index of ${ref} out of range\") }";
				}
			}
			push @processed, "}";
			if ($arrayArgs{'c'}) {
				push @processed, "if ldc*(n-1)+n > len(c) || ldc < max(1, n) { panic(\"cblas: index of c out of range\") }";
			}
		}
		if ($scalarArgs{'tA'} && $scalarArgs{'tB'}) {
			push @processed, "var rowA, colA, rowB, colB int";
			push @processed, "if tA == blas.NoTrans { rowA, colA = m, k } else { rowA, colA = k, m }";
			push @processed, "if tB == blas.NoTrans { rowB, colB = k, n } else { rowB, colB = n, k }";
			push @processed, "if o == blas.RowMajor {";
			push @processed, "if lda*(rowA-1)+colA > len(a) || lda < max(1, colA) { panic(\"cblas: index of a out of range\") }";
			push @processed, "if ldb*(rowB-1)+colB > len(b) || ldb < max(1, colB) { panic(\"cblas: index of b out of range\") }";
			push @processed, "if ldc*(m-1)+n > len(c) || ldc < max(1, n) { panic(\"cblas: index of c out of range\") }";
			push @processed, "} else {";
			push @processed, "if lda*(colA-1)+rowA > len(a) || lda < max(1, rowA) { panic(\"cblas: index of a out of range\") }";
			push @processed, "if ldb*(colB-1)+rowB > len(b) || ldb < max(1, rowB) { panic(\"cblas: index of b out of range\") }";
			push @processed, "if ldc*(n-1)+m > len(c) || ldc < max(1, m) { panic(\"cblas: index of c out of range\") }";
			push @processed, "}";
		}
	}

	my $checks = join "\n", @processed;
	$checks .= "\n" if scalar @processed > 0;
	return $checks
}

sub processParamToC {
	my $func = shift;
	my $paramList = shift;
	my @processed;
	my @params = split ',', $paramList;
	foreach my $param (@params) {
		my @parts = split /[ *]/, $param;
		my $var = lcfirst $parts[scalar @parts - 1];
		$param =~ m/^(?:const )?int \*[a-zA-Z]/ && do {
			push @processed, "(*C.int)(&".$var.")"; next;
		};
		$param =~ m/^(?:const )?void \*[a-zA-Z]/ && do {
			my $type;
			if ($var eq "alpha" || $var eq "beta") {
				$type = "";
			} else {
				$type = "[0]";
			}
			push @processed, "unsafe.Pointer(&".$var.$type.")"; next;
		};
		$param =~ m/^(?:const )?char \*[a-zA-Z]/ && do {
			push @processed, "(*C.char)(&".$var.")"; next;
		};
		$param =~ m/^(?:const )?float \*[a-zA-Z]/ && do {
			push @processed, "(*C.float)(&".$var."[0])"; next;
		};
		$param =~ m/^(?:const )?double \*[a-zA-Z]/ && do {
			push @processed, "(*C.double)(&".$var."[0])"; next;
		};
		$param =~ m/^(?:const )?int [a-zA-Z]/ && do {
			push @processed, "C.int(".$var.")"; next;
		};
		$param =~ m/^(?:const )float [a-zA-Z]/ && do {
			push @processed, "C.float(".$var.")"; next;
		};
		$param =~ m/^(?:const )double [a-zA-Z]/ && do {
			push @processed, "C.double(".$var.")"; next;
		};
		$param =~ m/^const enum [a-zA-Z]/ && do {
			$var eq "order" && do {
				$var = "o";
				push @processed, "C.enum_$parts[scalar @parts - 2](".$var.")"; next;
			};
			$var =~ /trans/ && do {
				$var =~ s/trans([AB]?)/t$1/;
				push @processed, "C.enum_$parts[scalar @parts - 2](".$var.")"; next;
			};
			$var eq "uplo" && do {
				$var = "ul";
				push @processed, "C.enum_$parts[scalar @parts - 2](".$var.")"; next;
			};
			$var eq "diag" && do {
				$var = "d";
				push @processed, "C.enum_$parts[scalar @parts - 2](".$var.")"; next;
			};
			$var eq "side" && do {
				$var = "s";
				push @processed, "C.enum_$parts[scalar @parts - 2](".$var.")"; next;
			};
		};
	}
	die "missed C parameters from '$func', '$paramList'" if scalar @processed != scalar @params;
	return join ", ", @processed;
}
