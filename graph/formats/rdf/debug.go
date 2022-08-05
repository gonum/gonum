// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build debug
// +build debug

package rdf

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

type debugger bool

const debug debugger = true

func (d debugger) log(depth int, args ...interface{}) {
	if !d {
		return
	}
	fmt.Fprint(os.Stderr, strings.Repeat("\t", depth))
	fmt.Fprintln(os.Stderr, args...)
}

func (d debugger) logf(depth int, format string, args ...interface{}) {
	if !d {
		return
	}
	fmt.Fprint(os.Stderr, strings.Repeat("\t", depth))
	fmt.Fprintf(os.Stderr, format, args...)
}

func (d debugger) logHashes(depth int, hashes map[string][]byte, size int) {
	if !d {
		return
	}
	prefix := strings.Repeat("\t", depth)
	if len(hashes) != 0 {
		keys := make([]string, len(hashes))
		i := 0
		for k := range hashes {
			keys[i] = k
			i++
		}
		sort.Strings(keys)
		w := tabwriter.NewWriter(os.Stderr, 0, 4, 8, ' ', 0)
		for _, k := range keys {
			fmt.Fprintf(w, prefix+"%s\t%0*x\n", k, 2*size, hashes[k])
		}
		w.Flush()
	} else {
		fmt.Fprintln(os.Stderr, prefix+"none")
	}
	fmt.Fprintln(os.Stderr)
}

func (d debugger) logParts(depth int, parts byLengthHash) {
	if !d {
		return
	}
	prefix := strings.Repeat("\t", depth)
	if parts.Len() != 0 {
		w := tabwriter.NewWriter(os.Stderr, 0, 4, 8, ' ', 0)
		for i, p := range parts.nodes {
			fmt.Fprintf(w, prefix+"%v\t%x\n", p, parts.hashes[i])
		}
		w.Flush()
	} else {
		fmt.Fprintln(os.Stderr, prefix+"none")
	}
	fmt.Fprintln(os.Stderr)
}
