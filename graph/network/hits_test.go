// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package network

import (
	"cmp"
	"fmt"
	"math"
	"slices"
	"testing"

	"gonum.org/v1/gonum/floats/scalar"
	"gonum.org/v1/gonum/graph/simple"
)

var hitsTests = []struct {
	g   []set
	tol float64

	wantTol float64
	want    map[int64]HubAuthority
}{
	{
		// Example graph from http://www.cis.hut.fi/Opinnot/T-61.6020/2008/pagerank_hits.pdf page 8.
		g: []set{
			A: linksTo(B, C, D),
			B: linksTo(C, D),
			C: linksTo(B),
			D: nil,
		},
		tol: 1e-4,

		wantTol: 1e-4,
		want: map[int64]HubAuthority{
			A: {Hub: 0.7887, Authority: 0},
			B: {Hub: 0.5774, Authority: 0.4597},
			C: {Hub: 0.2113, Authority: 0.6280},
			D: {Hub: 0, Authority: 0.6280},
		},
	},
}

func TestHITS(t *testing.T) {
	for i, test := range hitsTests {
		g := simple.NewDirectedGraph()
		for u, e := range test.g {
			// Add nodes that are not defined by an edge.
			if g.Node(int64(u)) == nil {
				g.AddNode(simple.Node(u))
			}
			for v := range e {
				g.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
			}
		}
		got := HITS(g, test.tol)
		prec := 1 - int(math.Log10(test.wantTol))
		for n := range test.g {
			if !scalar.EqualWithinAbsOrRel(got[int64(n)].Hub, test.want[int64(n)].Hub, test.wantTol, test.wantTol) {
				t.Errorf("unexpected HITS result for test %d:\ngot: %v\nwant:%v",
					i, orderedHubAuth(got, prec), orderedHubAuth(test.want, prec))
				break
			}
			if !scalar.EqualWithinAbsOrRel(got[int64(n)].Authority, test.want[int64(n)].Authority, test.wantTol, test.wantTol) {
				t.Errorf("unexpected HITS result for test %d:\ngot: %v\nwant:%v",
					i, orderedHubAuth(got, prec), orderedHubAuth(test.want, prec))
				break
			}
		}
	}
}

func orderedHubAuth(w map[int64]HubAuthority, prec int) []keyHubAuthVal {
	o := make([]keyHubAuthVal, 0, len(w))
	for k, v := range w {
		o = append(o, keyHubAuthVal{prec: prec, key: k, val: v})
	}
	slices.SortFunc(o, func(a, b keyHubAuthVal) int { return cmp.Compare(a.key, b.key) })
	return o
}

type keyHubAuthVal struct {
	prec int
	key  int64
	val  HubAuthority
}

func (kv keyHubAuthVal) String() string {
	return fmt.Sprintf("%d:{H:%.*f, A:%.*f}",
		kv.key, kv.prec, kv.val.Hub, kv.prec, kv.val.Authority,
	)
}
