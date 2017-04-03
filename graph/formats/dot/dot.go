// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dot implements a parser for Graphviz DOT files.
package dot

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/gonum/gonum/graph/formats/dot/ast"
	"github.com/gonum/gonum/graph/formats/dot/internal/lexer"
	"github.com/gonum/gonum/graph/formats/dot/internal/parser"
)

// ParseFile parses the given Graphviz DOT file into an AST.
func ParseFile(path string) (*ast.File, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseBytes(buf)
}

// Parse parses the given Graphviz DOT file into an AST, reading from r.
func Parse(r io.Reader) (*ast.File, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return ParseBytes(buf)
}

// ParseBytes parses the given Graphviz DOT file into an AST, reading from b.
func ParseBytes(b []byte) (*ast.File, error) {
	l := lexer.NewLexer(b)
	p := parser.NewParser()
	file, err := p.Parse(l)
	if err != nil {
		return nil, err
	}
	f, ok := file.(*ast.File)
	if !ok {
		return nil, fmt.Errorf("invalid file type; expected *ast.File, got %T", file)
	}
	if err := check(f); err != nil {
		return nil, err
	}
	return f, nil
}

// ParseString parses the given Graphviz DOT file into an AST, reading from s.
func ParseString(s string) (*ast.File, error) {
	return ParseBytes([]byte(s))
}
