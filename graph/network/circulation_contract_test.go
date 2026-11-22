//go:build challenge

// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package network_test

import (
	"errors"
	"math"
	"math/big"
	"reflect"
	"testing"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/multi"
	"gonum.org/v1/gonum/graph/network"
)

type circulationFixture struct {
	*multi.DirectedGraph
	ids      []int64
	limits   map[network.FlowLine]network.FlowBounds
	balances map[int64]int64
}

// Iteration order is a public graph-interface choice, fixed for reproducibility.
func (f *circulationFixture) Nodes() graph.Nodes {
	nodes := make([]graph.Node, len(f.ids))
	for i, id := range f.ids {
		nodes[i] = multi.Node(id)
	}
	return iterator.NewOrderedNodes(nodes)
}

func (f *circulationFixture) From(id int64) graph.Nodes {
	var nodes []graph.Node
	for _, to := range f.ids {
		if f.HasEdgeFromTo(id, to) {
			nodes = append(nodes, multi.Node(to))
		}
	}
	return iterator.NewOrderedNodes(nodes)
}

func fixture(ids ...int64) *circulationFixture {
	f := &circulationFixture{DirectedGraph: multi.NewDirectedGraph(), ids: ids,
		limits: make(map[network.FlowLine]network.FlowBounds), balances: make(map[int64]int64)}
	for _, id := range ids {
		f.AddNode(multi.Node(id))
	}
	return f
}

func (f *circulationFixture) edge(from, to, id, lower, upper, cost int64) network.FlowLine {
	key := network.FlowLine{From: from, To: to, ID: id}
	f.SetLine(multi.Line{F: multi.Node(from), T: multi.Node(to), UID: id})
	f.limits[key] = network.FlowBounds{Lower: lower, Upper: upper, Cost: cost}
	return key
}

func (f *circulationFixture) bounds(line graph.Line) network.FlowBounds {
	return f.limits[network.FlowLine{From: line.From().ID(), To: line.To().ID(), ID: line.ID()}]
}

func integer(s string) *big.Int {
	n, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic(s)
	}
	return n
}

// This oracle checks public mathematical conditions independently of the
// implementation's verifier. It deliberately permits multiple optima.
func optimum(t *testing.T, f *circulationFixture, want string) *network.Circulation {
	t.Helper()
	r, err := network.MinCostCirculation(f, f.bounds, f.balances)
	if err != nil || r == nil {
		t.Fatalf("expected optimum, result=%v error=%v", r, err)
	}
	if len(r.Flows) != len(f.limits) || len(r.Potentials) != len(f.ids) {
		t.Fatal("incomplete result maps")
	}
	net := make(map[int64]*big.Int)
	for _, id := range f.ids {
		net[id] = new(big.Int)
		if r.Potentials[id] == nil {
			t.Fatal("missing potential")
		}
	}
	var total big.Int
	for key, b := range f.limits {
		v := r.Flows[key]
		if v == nil || v.Cmp(big.NewInt(b.Lower)) < 0 || (b.Upper != -1 && v.Cmp(big.NewInt(b.Upper)) > 0) {
			t.Fatalf("invalid flow %v: %v", key, v)
		}
		net[key.From].Add(net[key.From], v)
		net[key.To].Sub(net[key.To], v)
		total.Add(&total, new(big.Int).Mul(v, big.NewInt(b.Cost)))
		reduced := new(big.Int).Sub(r.Potentials[key.From], r.Potentials[key.To])
		reduced.Add(reduced, big.NewInt(b.Cost))
		if (b.Upper == -1 || v.Cmp(big.NewInt(b.Upper)) < 0) && reduced.Sign() < 0 {
			t.Fatal("negative forward residual cost")
		}
		if v.Cmp(big.NewInt(b.Lower)) > 0 && reduced.Sign() > 0 {
			t.Fatal("negative reverse residual cost")
		}
	}
	for id, n := range net {
		if n.Cmp(big.NewInt(f.balances[id])) != 0 {
			t.Fatalf("unbalanced node %d", id)
		}
	}
	if r.Cost == nil || total.Cmp(r.Cost) != 0 || total.Cmp(integer(want)) != 0 {
		t.Fatalf("cost %v, reported %v, want %s", &total, r.Cost, want)
	}
	if err := network.VerifyCirculation(f, f.bounds, f.balances, r); err != nil {
		t.Fatalf("valid optimum rejected: %v", err)
	}
	return r
}

func infeasible(t *testing.T, f *circulationFixture) *network.InfeasibleCirculationError {
	t.Helper()
	r, err := network.MinCostCirculation(f, f.bounds, f.balances)
	var cut *network.InfeasibleCirculationError
	if r != nil || !errors.As(err, &cut) {
		t.Fatalf("expected infeasible certificate, got %v, %v", r, err)
	}
	inside := map[int64]bool{}
	var excess big.Int
	for _, id := range cut.Cut {
		if f.Node(id) == nil || inside[id] {
			t.Fatal("invalid cut membership")
		}
		inside[id] = true
		excess.Add(&excess, big.NewInt(f.balances[id]))
	}
	for key, b := range f.limits {
		if inside[key.From] && !inside[key.To] {
			if b.Upper == -1 {
				t.Fatal("unlimited edge crosses violated cut")
			}
			excess.Sub(&excess, big.NewInt(b.Upper))
		}
		if !inside[key.From] && inside[key.To] {
			excess.Add(&excess, big.NewInt(b.Lower))
		}
	}
	if excess.Sign() <= 0 || cut.Excess == nil || cut.Excess.Cmp(&excess) != 0 {
		t.Fatal("invalid cut excess")
	}
	if err := network.VerifyInfeasibleCirculation(f, f.bounds, f.balances, cut); err != nil {
		t.Fatal(err)
	}
	return cut
}

func TestCirculationBalances(t *testing.T) {
	f := fixture(1, 2, 3)
	f.balances[1] = 4
	f.balances[3] = -4
	f.edge(1, 2, 0, 0, 9, 2)
	f.edge(2, 3, 0, 0, 9, 3)
	optimum(t, f, "20")
	optimum(t, fixture(), "0")
	optimum(t, fixture(-99), "0")
}

func TestCirculationLowerBounds(t *testing.T) {
	f := fixture(1, 2, 3)
	f.edge(1, 2, 0, 3, 3, 5)
	f.edge(2, 3, 0, 0, 5, 2)
	f.edge(3, 1, 0, 0, 5, 1)
	optimum(t, f, "24")
	g := fixture(1, 2)
	g.edge(1, 2, 0, 1, 2, 0)
	infeasible(t, g)
}

func TestCirculationCapacityBoundaries(t *testing.T) {
	for _, amount := range []int64{0, 3, 4} {
		f := fixture(1, 2)
		f.edge(1, 2, 0, 0, 3, 1)
		f.balances[1] = amount
		f.balances[2] = -amount
		if amount == 4 {
			infeasible(t, f)
		} else {
			optimum(t, f, big.NewInt(amount).String())
		}
	}
	f := fixture(1, 2)
	f.edge(1, 2, 0, 0, 0, -10)
	optimum(t, f, "0")
}

func TestCirculationInputValidation(t *testing.T) {
	f := fixture(1, 2)
	key := f.edge(1, 2, 0, 0, 3, 1)
	for _, b := range []network.FlowBounds{{Lower: -1, Upper: 1}, {Lower: 2, Upper: 1}, {Upper: -2}} {
		f.limits[key] = b
		if r, e := network.MinCostCirculation(f, f.bounds, nil); e == nil || r != nil {
			t.Fatal("invalid bounds accepted")
		}
	}
	f.limits[key] = network.FlowBounds{Upper: 3}
	for _, balances := range []map[int64]int64{{7: 0}, {1: 1}, {1: math.MaxInt64, 2: math.MaxInt64}} {
		if r, e := network.MinCostCirculation(f, f.bounds, balances); e == nil || r != nil {
			t.Fatal("invalid balance accepted")
		}
	}
	if r, e := network.MinCostCirculation(nil, f.bounds, nil); e == nil || r != nil {
		t.Fatal("nil graph accepted")
	}
	if r, e := network.MinCostCirculation(f, nil, nil); e == nil || r != nil {
		t.Fatal("nil bounds accepted")
	}
	g := fixture(1, 2, 3)
	invalid := map[int64]int64{1: math.MaxInt64, 2: math.MaxInt64, 3: 2}
	cut := &network.InfeasibleCirculationError{Cut: []int64{1, 2, 3}, Excess: integer("18446744073709551616")}
	if network.VerifyInfeasibleCirculation(g, g.bounds, invalid, cut) == nil {
		t.Fatal("overflowing nonzero balance sum accepted")
	}
}

func TestCirculationParallelIdentity(t *testing.T) {
	f := fixture(-12, 99)
	f.balances[-12] = 4
	f.balances[99] = -4
	f.edge(-12, 99, -8, 0, 1, 2)
	f.edge(-12, 99, 40, 0, 2, 3)
	f.edge(-12, 99, 100, 0, 5, 7)
	f.edge(99, -12, 40, 0, 1, 10)
	optimum(t, f, "15")
}

func TestCirculationDisconnectedCycles(t *testing.T) {
	f := fixture(0, 1, 80, 90)
	f.balances[0] = 1
	f.balances[1] = -1
	f.edge(0, 1, 0, 0, 1, 2)
	f.edge(80, 90, 0, 0, 4, -8)
	f.edge(90, 80, 0, 0, 3, 1)
	optimum(t, f, "-19")
}

func TestCirculationSelfLoops(t *testing.T) {
	f := fixture(12)
	f.edge(12, 12, -4, 2, 6, -3)
	f.edge(12, 12, 9, 2, 8, 5)
	optimum(t, f, "-8")
}

func TestCirculationRerouting(t *testing.T) {
	f := fixture(0, 1, 2, 3)
	f.balances[0] = 2
	f.balances[3] = -2
	f.edge(0, 1, 0, 0, 1, 0)
	f.edge(0, 2, 0, 0, 1, 0)
	f.edge(1, 3, 0, 0, 1, 9)
	f.edge(2, 3, 0, 0, 1, 1)
	f.edge(1, 2, 0, 0, 1, -5)
	f.edge(2, 1, 0, 0, 1, -2)
	optimum(t, f, "3")
}

func TestCirculationExactArithmetic(t *testing.T) {
	f := fixture(1, 2, 3)
	f.balances[1] = math.MaxInt64
	f.balances[3] = -math.MaxInt64
	f.edge(1, 2, 0, 0, -1, math.MinInt64)
	f.edge(2, 3, 0, 0, -1, math.MinInt64)
	want := new(big.Int).Mul(big.NewInt(math.MaxInt64), big.NewInt(math.MinInt64))
	want.Mul(want, big.NewInt(2))
	optimum(t, f, want.String())
	g := fixture(1, 2, 3, 4)
	g.balances[1] = math.MaxInt64
	g.balances[2] = math.MaxInt64
	g.balances[3] = -math.MaxInt64
	g.balances[4] = -math.MaxInt64
	g.edge(1, 2, 0, 0, -1, 0)
	bridge := g.edge(2, 3, 0, 0, -1, 1)
	g.edge(3, 4, 0, 0, -1, 0)
	r := optimum(t, g, "18446744073709551614")
	if r.Flows[bridge].String() != "18446744073709551614" {
		t.Fatal("aggregate flow truncated")
	}
	h := fixture(1, 2)
	h.balances[1] = 1
	h.balances[2] = -1
	h.edge(1, 2, 0, 0, -1, math.MinInt64)
	h.edge(1, 2, 1, 0, 1, math.MaxInt64)
	optimum(t, h, big.NewInt(math.MinInt64).String())
}

func TestCirculationUnlimitedFinite(t *testing.T) {
	f := fixture(1, 2)
	f.edge(1, 2, 0, 0, -1, -5)
	f.edge(2, 1, 0, 0, 3, 2)
	optimum(t, f, "-9")
	g := fixture(1)
	g.edge(1, 1, 0, 2, -1, 0)
	optimum(t, g, "0")
}

func TestCirculationUnbounded(t *testing.T) {
	for _, self := range []bool{false, true} {
		f := fixture(1, 2, 3, 4)
		f.balances[1] = 1
		f.balances[2] = -1
		f.edge(1, 2, 0, 0, 1, 5)
		if self {
			f.edge(3, 3, 0, 0, -1, -1)
		} else {
			f.edge(3, 4, 0, 0, -1, math.MinInt64)
			f.edge(4, 3, 0, 0, -1, -1)
		}
		r, err := network.MinCostCirculation(f, f.bounds, f.balances)
		var u *network.UnboundedCirculationError
		if r != nil || !errors.As(err, &u) {
			t.Fatalf("expected unbounded: %v", err)
		}
		if err := network.VerifyUnboundedCirculation(f, f.bounds, f.balances, u); err != nil {
			t.Fatal(err)
		}
		if len(u.Feasible.Flows) != len(f.limits) {
			t.Fatal("incomplete feasible witness")
		}
		assignmentNet := map[int64]*big.Int{1: new(big.Int), 2: new(big.Int), 3: new(big.Int), 4: new(big.Int)}
		var assignmentCost big.Int
		for k, b := range f.limits {
			v := u.Feasible.Flows[k]
			if v == nil || v.Cmp(big.NewInt(b.Lower)) < 0 || (b.Upper != -1 && v.Cmp(big.NewInt(b.Upper)) > 0) {
				t.Fatal("invalid feasible witness flow")
			}
			assignmentNet[k.From].Add(assignmentNet[k.From], v)
			assignmentNet[k.To].Sub(assignmentNet[k.To], v)
			assignmentCost.Add(&assignmentCost, new(big.Int).Mul(v, big.NewInt(b.Cost)))
		}
		for id, v := range assignmentNet {
			if v.Cmp(big.NewInt(f.balances[id])) != 0 {
				t.Fatal("infeasible witness balance")
			}
		}
		if u.Feasible.Cost == nil || u.Feasible.Cost.Cmp(&assignmentCost) != 0 {
			t.Fatal("incorrect witness cost")
		}
		var cost big.Int
		net := map[int64]*big.Int{1: new(big.Int), 2: new(big.Int), 3: new(big.Int), 4: new(big.Int)}
		for k, v := range u.Ray {
			b, ok := f.limits[k]
			if !ok || v == nil || v.Sign() < 0 || (v.Sign() > 0 && b.Upper != -1) {
				t.Fatal("invalid ray entry")
			}
			net[k.From].Add(net[k.From], v)
			net[k.To].Sub(net[k.To], v)
			cost.Add(&cost, new(big.Int).Mul(v, big.NewInt(b.Cost)))
		}
		for _, v := range net {
			if v.Sign() != 0 {
				t.Fatal("unbalanced ray")
			}
		}
		if cost.Sign() >= 0 {
			t.Fatal("non-improving ray")
		}
	}
}

func TestCirculationInfeasiblePrecedence(t *testing.T) {
	f := fixture(1, 2, 3)
	f.balances[1] = 1
	f.balances[2] = -1
	f.edge(3, 3, 0, 0, -1, -1)
	infeasible(t, f)
}

func TestCirculationCutCertificate(t *testing.T) {
	f := fixture(1, 2, 3)
	f.edge(1, 2, 0, 5, 6, 0)
	f.edge(2, 1, 0, 0, 2, 0)
	f.edge(2, 3, 0, 0, -1, 0)
	infeasible(t, f)
	valid := &network.InfeasibleCirculationError{Cut: []int64{2, 3}, Excess: big.NewInt(3)}
	if err := network.VerifyInfeasibleCirculation(f, f.bounds, nil, valid); err != nil {
		t.Fatal(err)
	}
	for _, c := range []*network.InfeasibleCirculationError{nil, {}, {Cut: []int64{2, 3}, Excess: big.NewInt(4)}, {Cut: []int64{2}, Excess: big.NewInt(3)}, {Cut: []int64{99}, Excess: big.NewInt(3)}, {Cut: []int64{2, 2, 3}, Excess: big.NewInt(3)}, {Cut: []int64{1}, Excess: big.NewInt(1)}} {
		if network.VerifyInfeasibleCirculation(f, f.bounds, nil, c) == nil {
			t.Fatal("invalid cut accepted")
		}
	}
	for _, c := range []*network.InfeasibleCirculationError{{Excess: big.NewInt(0)}, {Cut: []int64{1, 2, 3}, Excess: big.NewInt(0)}, {Cut: []int64{1}, Excess: big.NewInt(-6)}} {
		if network.VerifyInfeasibleCirculation(f, f.bounds, nil, c) == nil {
			t.Fatal("nonpositive cut excess accepted")
		}
	}
	g := fixture(1, 2, 3)
	g.balances[1] = 1
	g.balances[2] = -1
	alternative := &network.InfeasibleCirculationError{Cut: []int64{3, 1}, Excess: big.NewInt(1)}
	if e := network.VerifyInfeasibleCirculation(g, g.bounds, g.balances, alternative); e != nil {
		t.Fatal("valid alternative cut rejected", e)
	}
	if !reflect.DeepEqual(alternative.Cut, []int64{3, 1}) || alternative.Excess.Cmp(big.NewInt(1)) != 0 {
		t.Fatal("cut certificate modified")
	}
}

func cloneAssignment(r *network.Circulation) *network.Circulation {
	c := &network.Circulation{Flows: map[network.FlowLine]*big.Int{}, Cost: new(big.Int).Set(r.Cost), Potentials: map[int64]*big.Int{}}
	for k, v := range r.Flows {
		c.Flows[k] = new(big.Int).Set(v)
	}
	for k, v := range r.Potentials {
		c.Potentials[k] = new(big.Int).Set(v)
	}
	return c
}

func TestCirculationVerifyPrimal(t *testing.T) {
	f := fixture(1, 2)
	key := f.edge(1, 2, 0, 1, 3, 2)
	f.balances[1] = 2
	f.balances[2] = -2
	r := optimum(t, f, "4")
	for _, change := range []func(*network.Circulation){
		func(c *network.Circulation) { c.Cost = nil }, func(c *network.Circulation) { c.Cost.Add(c.Cost, big.NewInt(1)) },
		func(c *network.Circulation) { delete(c.Flows, key) }, func(c *network.Circulation) { c.Flows[key] = nil },
		func(c *network.Circulation) { c.Flows[network.FlowLine{ID: 99}] = big.NewInt(0) },
		func(c *network.Circulation) { delete(c.Flows, key); c.Flows[network.FlowLine{ID: 99}] = big.NewInt(2) },
		func(c *network.Circulation) { c.Flows[key].SetInt64(0); c.Cost.SetInt64(0) },
		func(c *network.Circulation) { c.Flows[key].SetInt64(4); c.Cost.SetInt64(8) },
		func(c *network.Circulation) { c.Flows[key].SetInt64(1); c.Cost.SetInt64(2) },
	} {
		c := cloneAssignment(r)
		change(c)
		if network.VerifyCirculation(f, f.bounds, f.balances, c) == nil {
			t.Fatal("invalid primal accepted")
		}
	}
	if network.VerifyCirculation(f, f.bounds, f.balances, nil) == nil {
		t.Fatal("nil certificate accepted")
	}
}

func TestCirculationVerifyDual(t *testing.T) {
	f := fixture(1, 2)
	f.edge(1, 2, 0, 0, 3, 2)
	f.balances[1] = 1
	f.balances[2] = -1
	r := optimum(t, f, "2")
	for _, change := range []func(*network.Circulation){
		func(c *network.Circulation) { delete(c.Potentials, 1) }, func(c *network.Circulation) { c.Potentials[1] = nil },
		func(c *network.Circulation) { c.Potentials[9] = big.NewInt(0) },
		func(c *network.Circulation) { delete(c.Potentials, 1); c.Potentials[9] = big.NewInt(0) },
		func(c *network.Circulation) { c.Potentials[1].Add(c.Potentials[1], big.NewInt(1)) },
		func(c *network.Circulation) { c.Potentials[1].Sub(c.Potentials[1], big.NewInt(1)) },
	} {
		c := cloneAssignment(r)
		change(c)
		if network.VerifyCirculation(f, f.bounds, f.balances, c) == nil {
			t.Fatal("invalid dual accepted")
		}
	}
	for _, v := range r.Potentials {
		v.Add(v, integer("999999999999999999999999"))
	}
	if err := network.VerifyCirculation(f, f.bounds, f.balances, r); err != nil {
		t.Fatal("normalization was imposed", err)
	}
}

func TestCirculationAlternateOptima(t *testing.T) {
	f := fixture(1, 2)
	a := f.edge(1, 2, 0, 0, 3, 0)
	b := f.edge(1, 2, 1, 0, 3, 0)
	f.balances[1] = 2
	f.balances[2] = -2
	for _, split := range []int64{0, 1, 2} {
		r := &network.Circulation{Flows: map[network.FlowLine]*big.Int{a: big.NewInt(split), b: big.NewInt(2 - split)}, Cost: big.NewInt(0), Potentials: map[int64]*big.Int{1: big.NewInt(7), 2: big.NewInt(7)}}
		if err := network.VerifyCirculation(f, f.bounds, f.balances, r); err != nil {
			t.Fatal("valid alternative rejected", err)
		}
	}
}

func TestCirculationVerifyRay(t *testing.T) {
	f := fixture(1, 2)
	a := f.edge(1, 2, 0, 0, -1, -2)
	b := f.edge(2, 1, 0, 0, -1, 1)
	c := f.edge(1, 1, 0, 0, 4, -1)
	feasible := &network.Circulation{Flows: map[network.FlowLine]*big.Int{a: big.NewInt(0), b: big.NewInt(0), c: big.NewInt(0)}, Cost: big.NewInt(0)}
	valid := &network.UnboundedCirculationError{Feasible: feasible, Ray: map[network.FlowLine]*big.Int{a: big.NewInt(7), b: big.NewInt(7), c: big.NewInt(0)}}
	if err := network.VerifyUnboundedCirculation(f, f.bounds, nil, valid); err != nil {
		t.Fatal(err)
	}
	for _, ray := range []map[network.FlowLine]*big.Int{nil, {}, {a: big.NewInt(1)}, {a: big.NewInt(-1), b: big.NewInt(-1)}, {c: big.NewInt(1)}, {a: nil}, {network.FlowLine{ID: 77}: big.NewInt(0)}} {
		bad := &network.UnboundedCirculationError{Feasible: feasible, Ray: ray}
		if network.VerifyUnboundedCirculation(f, f.bounds, nil, bad) == nil {
			t.Fatal("invalid ray accepted")
		}
	}
	if network.VerifyUnboundedCirculation(f, f.bounds, nil, nil) == nil {
		t.Fatal("nil unbounded certificate accepted")
	}
	valid.Feasible = cloneAssignment(feasible)
	valid.Feasible.Cost.SetInt64(1)
	if network.VerifyUnboundedCirculation(f, f.bounds, nil, valid) == nil {
		t.Fatal("invalid feasible witness accepted")
	}
	g := fixture(1)
	k := g.edge(1, 1, 0, 0, -1, 0)
	zero := &network.UnboundedCirculationError{Feasible: &network.Circulation{Flows: map[network.FlowLine]*big.Int{k: big.NewInt(0)}, Cost: big.NewInt(0)}, Ray: map[network.FlowLine]*big.Int{k: big.NewInt(1)}}
	if network.VerifyUnboundedCirculation(g, g.bounds, nil, zero) == nil {
		t.Fatal("zero cost ray accepted")
	}
}

func TestCirculationReplayOwnership(t *testing.T) {
	f := fixture(1, 2)
	f.edge(1, 2, 0, 0, 3, 2)
	f.balances[1] = 1
	f.balances[2] = -1
	r := optimum(t, f, "2")
	before := cloneAssignment(r)
	if network.VerifyCirculation(f, f.bounds, f.balances, r) != nil || !reflect.DeepEqual(r, before) {
		t.Fatal("verification changed certificate")
	}
	for _, v := range r.Flows {
		v.SetInt64(999)
	}
	r.Cost.SetInt64(999)
	for _, v := range r.Potentials {
		v.SetInt64(999)
	}
	optimum(t, f, "2")
	if f.balances[1] != 1 || f.balances[2] != -1 || len(f.limits) != 1 || f.Nodes().Len() != 2 {
		t.Fatal("input modified")
	}
}

func TestCirculationExhaustiveSmall(t *testing.T) {
	for code := 0; code < 81; code++ {
		f := fixture(1, 2, 3)
		f.balances[1] = 1
		f.balances[3] = -1
		ends := [][2]int64{{1, 2}, {2, 3}, {1, 3}, {3, 1}}
		n := code
		costs := make([]int64, 4)
		for i, e := range ends {
			costs[i] = int64(n%3 - 1)
			n /= 3
			f.edge(e[0], e[1], 0, 0, 2, costs[i])
		}
		best := int64(math.MaxInt64)
		for a := int64(0); a <= 2; a++ {
			for b := int64(0); b <= 2; b++ {
				for c := int64(0); c <= 2; c++ {
					for d := int64(0); d <= 2; d++ {
						if a != b || a+c-d != 1 {
							continue
						}
						v := a*costs[0] + b*costs[1] + c*costs[2] + d*costs[3]
						if v < best {
							best = v
						}
					}
				}
			}
		}
		optimum(t, f, big.NewInt(best).String())
	}
}
