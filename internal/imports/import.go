// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package imports provides an API to check whether Gonum code does
// not import deprecated packages.
package imports

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// CheckBlacklisted analyzes all Go files under dir for deprecated and
// blacklisted imports.
// If Check encounters multiple files importing deprecated imports, the
// first error is returned to the user.
func CheckBlacklisted(dir string, blacklist []string) error {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		switch {
		case info.IsDir():
			switch info.Name() {
			case "testdata":
				return filepath.SkipDir
			}
		default:
			if filepath.Ext(info.Name()) != ".go" {
				return nil
			}
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	for _, fname := range files {
		e := process(fname, fset, blacklist)
		if e != nil {
			if err == nil {
				err = e
			}
		}
	}
	return err
}

func process(fname string, fset *token.FileSet, blacklist []string) error {
	src, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	return checkImports(fset, src, fname, blacklist)
}

func checkImports(fset *token.FileSet, src []byte, fname string, blacklist []string) error {
	f, err := parser.ParseFile(fset, fname, src, parser.ImportsOnly)
	if err != nil {
		return err
	}

	imp := Error{File: fname}
	for _, s := range f.Imports {
		path := strings.Trim(s.Path.Value, `"`)
		if blacklisted(path, blacklist) {
			imp.Imports = append(imp.Imports, path)
		}
	}
	if len(imp.Imports) > 0 {
		return imp
	}
	return nil
}

func blacklisted(path string, blacklist []string) bool {
	for _, v := range blacklist {
		if strings.HasPrefix(path, v) {
			return true
		}
	}
	return false
}

// Error stores information about a deprecated import in a Go file.
type Error struct {
	File    string
	Imports []string
}

func (e Error) Error() string {
	return fmt.Sprintf(
		"%s: deprecated imports: %v",
		e.File,
		strings.Join(e.Imports, ", "),
	)
}
