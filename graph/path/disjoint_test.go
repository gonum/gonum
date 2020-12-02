// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"testing"
)

func TestDisjointSetMakeSet(t *testing.T) {
	t.Parallel()

	ds := make(djSet)
	ds.add(3)
	if len(ds) != 1 {
		t.Error("Disjoint set master map of wrong size")
	}

	node, ok := ds[3]
	if !ok {
		t.Error("Make set did not successfully add element")
	} else {
		if node == nil {
			t.Fatal("Disjoint set node from add is nil")
		}

		if node.rank != 0 {
			t.Error("Node rank set incorrectly")
		}

		if node.parent != nil {
			t.Error("Node parent set incorrectly")
		}
	}
}

func TestDisjointSetFind(t *testing.T) {
	t.Parallel()

	ds := make(djSet)
	ds.add(3)
	ds.add(4)
	ds.add(5)
	ds.union(ds.find(3), ds.find(4))

	if ds.find(3) == ds.find(5) {
		t.Error("Disjoint sets incorrectly found to be the same")
	}
}

func TestUnion(t *testing.T) {
	t.Parallel()

	ds := make(djSet)
	ds.add(3)
	ds.add(4)
	ds.add(5)
	ds.union(ds.find(3), ds.find(4))
	ds.union(ds.find(4), ds.find(5))

	if ds.find(3) != ds.find(5) {
		t.Error("Sets found to be disjoint after union")
	}
}
