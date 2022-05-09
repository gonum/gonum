// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !debug
// +build !debug

package rdf

type debugger bool

const debug debugger = false

func (d debugger) log(depth int, args ...interface{})                      {}
func (d debugger) logf(depth int, format string, args ...interface{})      {}
func (d debugger) logHashes(depth int, hashes map[string][]byte, size int) {}
func (d debugger) logParts(depth int, parts byLengthHash)                  {}
