// Copyright ©2023 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"fmt"
	"math"
	"math/rand/v2"
	"reflect"
	"slices"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/graphs/gen"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/internal/order"
)

var shortestTests = []struct {
	n, d int
	p    float64
	seed uint64
}{
	{n: 100, d: 2, p: 0.5, seed: 1},
	{n: 200, d: 2, p: 0.5, seed: 1},
	{n: 100, d: 4, p: 0.25, seed: 1},
	{n: 200, d: 4, p: 0.25, seed: 1},
	{n: 100, d: 16, p: 0.1, seed: 1},
	{n: 200, d: 16, p: 0.1, seed: 1},
}

func TestShortestAlts(t *testing.T) {
	for _, test := range shortestTests {
		t.Run(fmt.Sprintf("AllTo_%d×%d|%v", test.n, test.d, test.p), func(t *testing.T) {
			g := simple.NewDirectedGraph()
			gen.SmallWorldsBB(g, test.n, test.d, test.p, rand.New(rand.NewPCG(test.seed, test.seed)))
			all := allShortest(DijkstraAllPaths(g))

			for uid := int64(0); uid < int64(test.n); uid++ {
				p := DijkstraAllFrom(g.Node(uid), g)
				for vid := int64(0); vid < int64(test.n); vid++ {
					got, gotW := p.AllTo(vid)
					want, wantW := all.AllBetween(uid, vid)
					if gotW != wantW {
						t.Errorf("mismatched weight: got:%f want:%f", gotW, wantW)
						continue
					}

					var gotPaths [][]int64
					if len(got) != 0 {
						gotPaths = make([][]int64, len(got))
					}
					for i, p := range got {
						for _, v := range p {
							gotPaths[i] = append(gotPaths[i], v.ID())
						}
					}
					order.BySliceValues(gotPaths)
					var wantPaths [][]int64
					if len(want) != 0 {
						wantPaths = make([][]int64, len(want))
					}
					for i, p := range want {
						for _, v := range p {
							wantPaths[i] = append(wantPaths[i], v.ID())
						}
					}
					order.BySliceValues(wantPaths)
					if !reflect.DeepEqual(gotPaths, wantPaths) {
						t.Errorf("unexpected shortest paths %d --> %d:\ngot: %v\nwant:%v",
							uid, vid, gotPaths, wantPaths)
					}
				}
			}
		})
	}
}

func TestAllShortest(t *testing.T) {
	for _, test := range shortestTests {
		t.Run(fmt.Sprintf("AllBetween_%d×%d|%v", test.n, test.d, test.p), func(t *testing.T) {
			g := simple.NewDirectedGraph()
			gen.SmallWorldsBB(g, test.n, test.d, test.p, rand.New(rand.NewPCG(test.seed, test.seed)))

			p := DijkstraAllPaths(g)
			for uid := int64(0); uid < int64(test.n); uid++ {
				for vid := int64(0); vid < int64(test.n); vid++ {
					got, gotW := p.AllBetween(uid, vid)
					want, wantW := allShortest(p).AllBetween(uid, vid) // Compare to naive.
					if gotW != wantW {
						t.Errorf("mismatched weight: got:%f want:%f", gotW, wantW)
						continue
					}

					var gotPaths [][]int64
					if len(got) != 0 {
						gotPaths = make([][]int64, len(got))
					}
					for i, p := range got {
						for _, v := range p {
							gotPaths[i] = append(gotPaths[i], v.ID())
						}
					}
					order.BySliceValues(gotPaths)
					var wantPaths [][]int64
					if len(want) != 0 {
						wantPaths = make([][]int64, len(want))
					}
					for i, p := range want {
						for _, v := range p {
							wantPaths[i] = append(wantPaths[i], v.ID())
						}
					}
					order.BySliceValues(wantPaths)
					if !reflect.DeepEqual(gotPaths, wantPaths) {
						t.Errorf("unexpected shortest paths %d --> %d:\ngot: %v\nwant:%v",
							uid, vid, gotPaths, wantPaths)
					}
				}
			}
		})
	}
}

// allShortest implements an allocation-naive AllBetween.
type allShortest AllShortest

// at returns a slice of node indexes into p.nodes for nodes that are mid points
// between nodes indexed by from and to.
func (p allShortest) at(from, to int) (mid []int) {
	return p.next[from+to*len(p.nodes)]
}

// AllBetween returns all shortest paths from u to v and the weight of the paths. Paths
// containing zero-weight cycles are not returned. If a negative cycle exists between
// u and v, paths is returned nil and weight is returned as -Inf.
func (p allShortest) AllBetween(uid, vid int64) (paths [][]graph.Node, weight float64) {
	from, fromOK := p.indexOf[uid]
	to, toOK := p.indexOf[vid]
	if !fromOK || !toOK || len(p.at(from, to)) == 0 {
		if uid == vid {
			if !fromOK {
				return [][]graph.Node{{node(uid)}}, 0
			}
			return [][]graph.Node{{p.nodes[from]}}, 0
		}
		return nil, math.Inf(1)
	}

	weight = p.dist.At(from, to)
	if math.Float64bits(weight) == defacedBits {
		return nil, math.Inf(-1)
	}

	var n graph.Node
	if p.forward {
		n = p.nodes[from]
	} else {
		n = p.nodes[to]
	}
	seen := make([]bool, len(p.nodes))
	paths = p.allBetween(from, to, seen, []graph.Node{n}, nil)

	return paths, weight
}

// allBetween recursively constructs a slice of paths extending from the node
// indexed into p.nodes by from to the node indexed by to. len(seen) must match
// the number of nodes held by the receiver. The path parameter is the current
// working path and the results are written into paths.
func (p allShortest) allBetween(from, to int, seen []bool, path []graph.Node, paths [][]graph.Node) [][]graph.Node {
	if p.forward {
		seen[from] = true
	} else {
		seen[to] = true
	}
	if from == to {
		if path == nil {
			return paths
		}
		if !p.forward {
			slices.Reverse(path)
		}
		return append(paths, path)
	}
	first := true
	for _, n := range p.at(from, to) {
		if seen[n] {
			continue
		}
		if first {
			path = append([]graph.Node(nil), path...)
			first = false
		}
		if p.forward {
			from = n
		} else {
			to = n
		}
		path = path[:len(path):len(path)]
		paths = p.allBetween(from, to, append([]bool(nil), seen...), append(path, p.nodes[n]), paths)
	}
	return paths
}

var shortestBenchmarks = []struct {
	n, d int
	p    float64
	seed uint64
}{
	{n: 100, d: 2, p: 0.5, seed: 1},
	{n: 1000, d: 2, p: 0.5, seed: 1},
	{n: 100, d: 4, p: 0.25, seed: 1},
	{n: 1000, d: 4, p: 0.25, seed: 1},
	{n: 100, d: 16, p: 0.1, seed: 1},
	{n: 1000, d: 16, p: 0.1, seed: 1},
}

func BenchmarkShortestAlts(b *testing.B) {
	for _, bench := range shortestBenchmarks {
		g := simple.NewDirectedGraph()
		gen.SmallWorldsBB(g, bench.n, bench.d, bench.p, rand.New(rand.NewPCG(bench.seed, bench.seed)))

		// Find the widest path set.
		var (
			bestP   ShortestAlts
			bestVid int64
			n       int
		)
		for uid := int64(0); uid < int64(bench.n); uid++ {
			p := DijkstraAllFrom(g.Node(uid), g)
			for vid := int64(0); vid < int64(bench.n); vid++ {
				paths, _ := p.AllTo(vid)
				if len(paths) > n {
					n = len(paths)
					bestP = p
					bestVid = vid
				}
			}
		}

		b.Run(fmt.Sprintf("AllTo_%d×%d|%v(%d)", bench.n, bench.d, bench.p, n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				paths, _ := bestP.AllTo(bestVid)
				if len(paths) != n {
					b.Errorf("unexpected number of paths: got:%d want:%d", len(paths), n)
				}
			}
		})
		b.Run(fmt.Sprintf("AllToFunc_%d×%d|%v(%d)", bench.n, bench.d, bench.p, n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var paths int
				bestP.AllToFunc(bestVid, func(_ []graph.Node) { paths++ })
				if paths != n {
					b.Errorf("unexpected number of paths: got:%d want:%d", paths, n)
				}
			}
		})
	}
}

func BenchmarkAllShortest(b *testing.B) {
	shortestPathAlgs := []struct {
		name string
		fn   func(g graph.Graph) AllShortest
	}{
		{
			name: "DijkstraAllPaths",
			fn:   DijkstraAllPaths,
		},
		{
			name: "FloydWarshall",
			fn: func(g graph.Graph) AllShortest {
				p, _ := FloydWarshall(g)
				return p
			},
		},
	}

	for _, bench := range shortestBenchmarks {
		for _, f := range shortestPathAlgs {
			g := simple.NewDirectedGraph()
			gen.SmallWorldsBB(g, bench.n, bench.d, bench.p, rand.New(rand.NewPCG(bench.seed, bench.seed)))
			p := f.fn(g)

			// Find the widest path set.
			var (
				bestUid, bestVid int64
				n                int
			)
			for uid := int64(0); uid < int64(bench.n); uid++ {
				for vid := int64(0); vid < int64(bench.n); vid++ {
					paths, _ := p.AllBetween(uid, vid)
					if len(paths) > n {
						n = len(paths)
						bestUid = uid
						bestVid = vid
					}
				}
			}

			b.Run(fmt.Sprintf("%s_Between_%d×%d|%v(%d)", f.name, bench.n, bench.d, bench.p, n), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					path, _, _ := p.Between(bestUid, bestVid)
					if len(path) == 0 {
						b.Errorf("unexpected empty path: got:%d want:%d", len(path), 0)
					}
				}
			})
			b.Run(fmt.Sprintf("%s_AllBetween_%d×%d|%v(%d)", f.name, bench.n, bench.d, bench.p, n), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					paths, _ := p.AllBetween(bestUid, bestVid)
					if len(paths) != n {
						b.Errorf("unexpected number of paths: got:%d want:%d", len(paths), n)
					}
				}
			})
			b.Run(fmt.Sprintf("%s_AllBetweenFunc_%d×%d|%v(%d)", f.name, bench.n, bench.d, bench.p, n), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					var paths int
					p.AllBetweenFunc(bestUid, bestVid, func(_ []graph.Node) { paths++ })
					if paths != n {
						b.Errorf("unexpected number of paths: got:%d want:%d", paths, n)
					}
				}
			})
		}
	}
}
