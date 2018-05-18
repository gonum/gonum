// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"log"
	"os"

	"gonum.org/v1/gonum/internal/imports"
)

var blacklist = []string{
	"github.com/gonum/.*", // prefer gonum.org/v1/gonum
	"math/rand",           // prefer golang.org/x/exp/rand
}

func main() {
	log.SetPrefix("check-imports: ")
	log.SetFlags(0)

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not retrieve current working directory: %v", err)
	}
	log.Printf("analyzing imports under %q...", dir)
	err = imports.CheckBlacklisted(dir, blacklist)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("analyzing imports under %q... [OK]", dir)
}
