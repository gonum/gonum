#!/usr/bin/env perl
# Copyright ©2014 The gonum Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

use strict;
use warnings;

my $clapackHeader = "lapacke.h";
my $LIB = $ARGV[0];

my $excludeComplex = 0;

my @lapack_extendedprecision_objs = (
                "zposvxx", "clagge", "clatms", "chesvxx", "cposvxx", "cgesvxx", "ssyrfssx", "csyrfsx",
                "dlagsy", "dsysvxx", "sporfsx", "slatms", "zlatms", "zherfsx", "csysvxx", "dlatms",
                );
my %xobjs;
foreach my $obj (@lapack_extendedprecision_objs) {
	$xobjs{$obj} = 1;
}

open(my $clapack, "<", $clapackHeader) or die;
open(my $golapack, ">", "lapack.go") or die;

my %done;

printf $golapack <<"EOH";
// Do not manually edit this file. It was created by the genLapack.pl script from ${clapackHeader}.

// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package clapack

/*
#cgo CFLAGS: -g -O2
#cgo LDFLAGS: -L${LIB} -lopenblas
#include "${clapackHeader}"
*/
import "C"

import (
	"github.com/gonum/blas"
	"github.com/dane-unltd/lapack"
	"unsafe"
)

type La struct{}
 
func init() {
         _ = lapack.Complex128(La{})
         _ = lapack.Float64(La{})
}

EOH

$/ = undef;
my $header = <$clapack>;

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

our %typeConv = (
	"lapack_logical" => "int32",
	"lapack_logical*" => "[]int32",
	"lapack_int*" => "[]int32",
	"lapack_int" => "int",
	"float*" => "[]float32",
	"double*" => "[]float64 ",
	"float" => "float32",
	"lapack_complex_float" => "complex64",
	"lapack_complex_float*" => "[]complex64",
	"lapack_complex_double" => "complex128",
	"lapack_complex_double*" => "[]complex128",
	"double" => "float64",
	"char" => "byte",
	"char*" => "[]byte",
	"LAPACK_S_SELECT2" => "Select2Float32",
	"LAPACK_S_SELECT3" => "Select3Float32",
	"LAPACK_D_SELECT2" => "Select2Float64",
	"LAPACK_D_SELECT3" => "Select3Float64",
	"LAPACK_C_SELECT1" => "Select1Complex64",
	"LAPACK_C_SELECT2" => "Select2Complex64",
	"LAPACK_Z_SELECT1" => "Select1Complex128",
	"LAPACK_Z_SELECT2" => "Select2Complex128",
	"void" => ""
);

foreach my $line (@lines) {
	process($line);	
}

close($golapack);
`go fmt .`;

sub process {
	my $line = shift;
	chomp $line;
	processProto($line);
}

sub processProto {
	my $proto = shift;
	if(not($proto =~ /LAPACKE/)) {
		return
	}
	my ($func, $paramList) = split /[()]/, $proto;

	(my $ret, $func) = split ' ', $func;
	(my $pack, $func, my $tail) = split '_', $func;

	if ($done{$func} or $xobjs{$func}){
		return
	}

	if (substr($func,-2) eq "xx") {
		return
	}
	if (substr($func,-3) eq "fsx") {
		return
	}
	if (substr($func,1,3) eq "lag") {
		return
	}
	if ($func eq "ilaver") {
		return
	}
	$done{$func} = 1;


	my $gofunc;
	if ($tail) {
		$gofunc = ucfirst $func . ucfirst $tail;
	}else{
		$gofunc = ucfirst $func;
	}

	my $GoRet = $typeConv{$ret};
	my $complexType = $func;
	$complexType =~ s/.*_[isd]?([zc]).*/$1/;
	my ($params,$bp) = processParamToGo($func, $paramList, $complexType);
	if ($params eq "") {
		return
	}
	print $golapack "func (La) ".$gofunc."(".$params.") ".$GoRet."{\n";
	print $golapack "\t";
	if ($ret ne 'void') {
		print $golapack "\n".$bp."\n"."return ".$GoRet."(";
	}
	print $golapack "C.LAPACKE_$func(".processParamToC($func, $paramList).")";
	if ($ret ne 'void') {
		print $golapack ")";
	}
	print $golapack "\n}\n";
}

sub Gofunc {
	my $fnName = shift;
	my ($pack, $func, $tail) = split '_', $fnName;
	return ucfirst $func . ucfirst $tail if $tail;
	return ucfirst $func;
}

sub processParamToGo {
	my $func = shift;
	my $paramList = shift;
	my $complexType = shift;
	my @processed;
	my @boilerplate;
	my @params = split ',', $paramList;
	foreach my $param (@params) {
		$param =~ s/const //g;
		my ($type,$var) = split ' ', $param;
		$var eq "matrix_order" && do {
			$var = "o";
			push @processed, $var." blas.Order"; next;
		};
		$var =~ /trans/ && do {
			my $bp = << "EOH";
var $var C.char
if go$var == blas.NoTrans{ $var = 'n' }
if go$var == blas.Trans{ $var = 't' }
if go$var == blas.ConjTrans{ $var = 'c' }
EOH
			push @boilerplate, $bp;
			push @processed, "go".$var." blas.Transpose"; next;
		};
		$var eq "uplo" && do {
			$var = "ul";
			my $bp = << "EOH";
var $var C.char
if go$var == blas.Upper{ $var = 'u' }
if go$var == blas.Lower{ $var = 'l' }
EOH
			push @boilerplate, $bp;
			push @processed, "go".$var." blas.Uplo"; next;
		};
		$var eq "diag" && do {
			$var = "d";
			my $bp = << "EOH";
var $var C.char
if go$var == blas.Unit{ $var = 'u' }
if go$var == blas.NonUnit{ $var = 'n' }
EOH
			push @boilerplate, $bp;
			push @processed, "go".$var." blas.Diag"; next;
		};
		$var eq "side" && do {
			$var = "s";
			my $bp = << "EOH";
var $var C.char
if go$var == blas.Left{ $var = 'l' }
if go$var == blas.Right{ $var = 'r' }
EOH
			push @boilerplate, $bp;
			push @processed, "go".$var." blas.Side"; next;
		};
		$var eq "compq" && do {
			push @processed, $var." lapack.CompSV"; next;
		};
		$var =~ /job+/ && do {
			push @processed, $var." lapack.Job"; next;
		};
		$var eq "select" && do {
			$var = "sel";
			return ""
		};
		$var eq "selctg" && do {
			return ""
		};
		$var eq "range" && do {
			$var = "rng";
		};

		my $goType = $typeConv{$type};

		if (substr($type,-1) eq "*") {
			my $base = $typeConv{$type};
			$base =~ s/\[\]//;
			my $bp = << "EOH";
var $var *$base
if len(go$var) > 0 {
	$var = &go$var [0]
}
EOH
			push @boilerplate, $bp;
			$var = "go".$var;
		}

		if (not $goType) {
			die "missed Go parameters from '$func', '$type', '$param'";
		}
		push @processed, $var." ".$goType; next;
	}
	return ((join ", ", @processed), (join "\n", @boilerplate));
}

sub processParamToC {
	my $func = shift;
	my $paramList = shift;
	my @processed;
	my @boilerplate;
	my @params = split ',', $paramList;
	foreach my $param (@params) {
		$param =~ s/const //g;
		my ($type,$var) = split ' ', $param;

		$var eq "matrix_order" && do {
			$var = "o";
		};
		$var =~ /trans/ && do {
		};
		$var eq "uplo" && do {
			$var = "ul";
		};
		$var eq "diag" && do {
			$var = "d";
		};
		$var eq "side" && do {
			$var = "s";
		};
		$var eq "select" && do {
			$var = "sel";
		};
		$var eq "range" && do {
			$var = "rng";
		};

		if (substr($type,-1) eq "*") {
			chop $type;

			if ($type eq "char") {
				push @processed, "(*C.".$type.")(unsafe.Pointer(".$var."))"; next;
			} else {
				push @processed, "(*C.".$type.")(".$var.")"; next;
			}
		}else{
			push @processed, "(C.".$type.")(".$var.")"; next;
		}
	}
	die "missed C parameters from '$func', '$paramList'" if scalar @processed != scalar @params;
	return join ", ", @processed;
}
