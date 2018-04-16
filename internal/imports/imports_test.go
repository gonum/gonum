// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package imports

import (
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var (
	blacklist = []string{
		"github.com/gonum/", // prefer gonum.org/v1/gonum
		"math/rand",         // prefer golang.org/x/exp/rand
	}
)

func TestCheck(t *testing.T) {
	fset := token.NewFileSet()
	for _, tc := range []struct {
		pkg string
		err error
	}{
		{
			pkg: "github.com/gonum/",
			err: Error{
				File:    "file.go",
				Imports: []string{"github.com/gonum/"},
			},
		},
		{
			pkg: "github.com/gonum/floats",
			err: Error{
				File:    "file.go",
				Imports: []string{"github.com/gonum/floats"},
			},
		},
		{
			pkg: "github.com/gonum/plot",
			err: Error{
				File:    "file.go",
				Imports: []string{"github.com/gonum/plot"},
			},
		},
		{
			pkg: "gonum.org/v1/gonum/floats",
			err: nil,
		},
		{
			pkg: "gonum.org/v1/plot",
			err: nil,
		},
		{
			pkg: "github.com/gonumnum/floats",
			err: nil,
		},
	} {
		t.Run("", func(t *testing.T) {
			src := fmt.Sprintf("package foo\nimport _ %q\n", tc.pkg)
			err := checkImports(fset, []byte(src), "file.go", blacklist)
			if !reflect.DeepEqual(err, tc.err) {
				t.Fatalf("error\ngot= %v\nwant=%v", err, tc.err)
			}
		})
	}
}

func TestBlacklistedImports(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not retrieve current working directory: %v", err)
	}
	dir := filepath.Join(pwd, "../..") // Gonum top directory
	err = CheckBlacklisted(dir, blacklist)
	if err != nil {
		t.Fatal(err)
	}
}
