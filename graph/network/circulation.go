// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package network

import (
	"errors"
	"fmt"
	"math/big"

	"gonum.org/v1/gonum/graph"
)

// FlowLine identifies a line independently of graph iteration order.
type FlowLine struct {
	From, To, ID int64
}

// FlowBounds specifies inclusive bounds and a per-unit cost.
// Upper == -1 denotes unlimited capacity.
type FlowBounds struct {
	Lower, Upper, Cost int64
}

// Circulation is a primal solution and an exact dual certificate.
// Potentials have no prescribed additive normalization.
type Circulation struct {
	Flows      map[FlowLine]*big.Int
	Cost       *big.Int
	Potentials map[int64]*big.Int
}

// InfeasibleCirculationError certifies a violated cut inequality.
// Excess is the sum of balances in Cut minus its outgoing upper bounds
// plus its incoming lower bounds, and is strictly positive.
type InfeasibleCirculationError struct {
	Cut    []int64
	Excess *big.Int
}

func (e *InfeasibleCirculationError) Error() string {
	return "network: infeasible circulation"
}

// UnboundedCirculationError contains a feasible assignment and an improving
// nonnegative circulation Ray supported only on unlimited-capacity lines.
// Feasible has no dual potentials because no finite optimum exists.
type UnboundedCirculationError struct {
	Feasible *Circulation
	Ray      map[FlowLine]*big.Int
}

func (e *UnboundedCirculationError) Error() string {
	return "network: unbounded circulation cost"
}

type circulationLine struct {
	key      FlowLine
	from, to int
	bounds   FlowBounds
}

type circulationProblem struct {
	nodes    []int64
	index    map[int64]int
	lines    []circulationLine
	balances []big.Int
}

// MinCostCirculation finds a globally minimum-cost integral circulation.
// At every node, outgoing flow minus incoming flow must equal its balance.
// Missing balances are zero. Lower must be nonnegative; Upper must be -1 or at least Lower.
// The graph and the caller's data are not modified.
//
// Invalid inputs return an error. Valid infeasible inputs return a nil
// result and an *InfeasibleCirculationError. A feasible problem with no finite
// optimum returns an *UnboundedCirculationError. Feasibility takes precedence
// over unboundedness. An empty graph is feasible.
// All arithmetic, including intermediate sums and costs, is exact.
//
// The algorithm finds a feasible flow and cancels negative residual cycles.
// Each cancellation strictly lowers integral cost. Its running time is
// pseudo-polynomial in capacities and costs.
func MinCostCirculation(g graph.DirectedMultigraph, bounds func(graph.Line) FlowBounds, balances map[int64]int64) (*Circulation, error) {
	p, err := newCirculationProblem(g, bounds, balances)
	if err != nil {
		return nil, err
	}
	r, forward := p.residual()
	if cut := r.satisfy(p); cut != nil {
		return nil, cut
	}
	// Establish boundedness before cancelling finite-capacity cycles. This
	// avoids an arbitrary cancellation order indefinitely increasing flow on
	// unlimited arcs while alternating finite bottlenecks.
	cycle, _ := r.negativeCycle(len(p.nodes), true)
	var potentials []big.Int
	if cycle == nil {
		potentials, _ = r.minimize(len(p.nodes))
	}
	result := p.assignment(r, forward)
	if cycle != nil {
		ray := make(map[FlowLine]*big.Int)
		byArc := make(map[int]FlowLine, len(forward))
		for i, arc := range forward {
			byArc[arc] = p.lines[i].key
		}
		for _, arc := range cycle {
			ray[byArc[arc]] = big.NewInt(1)
		}
		return nil, &UnboundedCirculationError{Feasible: result, Ray: ray}
	}
	result.Potentials = make(map[int64]*big.Int, len(p.nodes))
	for i, id := range p.nodes {
		result.Potentials[id] = new(big.Int).Set(&potentials[i])
	}
	return result, nil
}

func (p *circulationProblem) assignment(r *circulationResidual, forward []int) *Circulation {
	result := &Circulation{
		Flows: make(map[FlowLine]*big.Int, len(p.lines)),
		Cost:  new(big.Int),
	}
	for i, line := range p.lines {
		used := new(big.Int).Add(big.NewInt(line.bounds.Lower), &r.arcs[forward[i]^1].capacity)
		result.Flows[line.key] = used
		var term big.Int
		term.Mul(used, big.NewInt(line.bounds.Cost))
		result.Cost.Add(result.Cost, &term)
	}
	return result
}

func newCirculationProblem(g graph.DirectedMultigraph, bounds func(graph.Line) FlowBounds, balances map[int64]int64) (*circulationProblem, error) {
	if g == nil || bounds == nil {
		return nil, errors.New("network: nil circulation input")
	}
	p := &circulationProblem{index: make(map[int64]int)}
	nodes := g.Nodes()
	for nodes.Next() {
		id := nodes.Node().ID()
		p.index[id] = len(p.nodes)
		p.nodes = append(p.nodes, id)
	}
	p.balances = make([]big.Int, len(p.nodes))
	var total big.Int
	for id, balance := range balances {
		i, ok := p.index[id]
		if !ok {
			return nil, fmt.Errorf("network: balance for unknown node %d", id)
		}
		p.balances[i].SetInt64(balance)
		total.Add(&total, &p.balances[i])
	}
	if total.Sign() != 0 {
		return nil, errors.New("network: balances do not sum to zero")
	}
	for i, id := range p.nodes {
		neighbors := g.From(id)
		for neighbors.Next() {
			to := neighbors.Node().ID()
			lines := g.Lines(id, to)
			for lines.Next() {
				line := lines.Line()
				b := bounds(line)
				if b.Lower < 0 || (b.Upper != -1 && b.Upper < b.Lower) {
					return nil, errors.New("network: invalid circulation bounds")
				}
				p.lines = append(p.lines, circulationLine{
					key:  FlowLine{From: id, To: to, ID: line.ID()},
					from: i, to: p.index[to], bounds: b,
				})
			}
		}
	}
	return p, nil
}

type circulationArc struct {
	from, to int
	capacity big.Int
	infinite bool
	cost     big.Int
}

type circulationResidual struct {
	arcs []circulationArc
	adj  [][]int
}

func (r *circulationResidual) add(from, to int, capacity, cost *big.Int) int {
	i := len(r.arcs)
	r.arcs = append(r.arcs, circulationArc{from: from, to: to}, circulationArc{from: to, to: from})
	if capacity == nil {
		r.arcs[i].infinite = true
	} else {
		r.arcs[i].capacity.Set(capacity)
	}
	r.arcs[i].cost.Set(cost)
	r.arcs[i+1].cost.Neg(cost)
	r.adj[from] = append(r.adj[from], i)
	r.adj[to] = append(r.adj[to], i+1)
	return i
}

func (p *circulationProblem) residual() (*circulationResidual, []int) {
	r := &circulationResidual{adj: make([][]int, len(p.nodes)+2)}
	forward := make([]int, len(p.lines))
	for i, line := range p.lines {
		var capacity *big.Int
		if line.bounds.Upper != -1 {
			capacity = big.NewInt(line.bounds.Upper - line.bounds.Lower)
		}
		forward[i] = r.add(line.from, line.to, capacity, big.NewInt(line.bounds.Cost))
	}
	return r, forward
}

func (r *circulationResidual) augment(path []int, amount *big.Int) {
	for _, i := range path {
		if !r.arcs[i].infinite {
			r.arcs[i].capacity.Sub(&r.arcs[i].capacity, amount)
		}
		if !r.arcs[i^1].infinite {
			r.arcs[i^1].capacity.Add(&r.arcs[i^1].capacity, amount)
		}
	}
}

func (r *circulationResidual) bottleneck(path []int) *big.Int {
	var amount *big.Int
	for _, i := range path {
		arc := &r.arcs[i]
		if arc.infinite {
			continue
		}
		if amount == nil {
			amount = new(big.Int).Set(&arc.capacity)
		} else if arc.capacity.Cmp(amount) < 0 {
			amount.Set(&arc.capacity)
		}
	}
	return amount
}

// path returns either an augmenting path or the reachable side of a cut.
func (r *circulationResidual) path(source, sink int) ([]int, []bool) {
	seen := make([]bool, len(r.adj))
	previous := make([]int, len(r.adj))
	queue := []int{source}
	seen[source] = true
	for head := 0; head < len(queue); head++ {
		from := queue[head]
		for _, i := range r.adj[from] {
			arc := &r.arcs[i]
			if (!arc.infinite && arc.capacity.Sign() <= 0) || seen[arc.to] {
				continue
			}
			seen[arc.to] = true
			previous[arc.to] = i
			if arc.to == sink {
				var path []int
				for v := sink; v != source; {
					j := previous[v]
					path = append(path, j)
					v = r.arcs[j].from
				}
				return path, seen
			}
			queue = append(queue, arc.to)
		}
	}
	return nil, seen
}

func (r *circulationResidual) satisfy(p *circulationProblem) *InfeasibleCirculationError {
	n := len(p.nodes)
	source, sink := n, n+1
	remaining := make([]big.Int, n)
	for i := range remaining {
		remaining[i].Set(&p.balances[i])
	}
	for _, line := range p.lines {
		lower := big.NewInt(line.bounds.Lower)
		remaining[line.from].Sub(&remaining[line.from], lower)
		remaining[line.to].Add(&remaining[line.to], lower)
	}
	firstAuxiliary := len(r.arcs)
	var required, zero big.Int
	for i := range remaining {
		switch remaining[i].Sign() {
		case 1:
			r.add(source, i, &remaining[i], &zero)
			required.Add(&required, &remaining[i])
		case -1:
			var capacity big.Int
			capacity.Neg(&remaining[i])
			r.add(i, sink, &capacity, &zero)
		}
	}
	for required.Sign() > 0 {
		path, seen := r.path(source, sink)
		if len(path) == 0 {
			cut := &InfeasibleCirculationError{Excess: new(big.Int)}
			for i, id := range p.nodes {
				if seen[i] {
					cut.Cut = append(cut.Cut, id)
				}
			}
			cut.Excess.Set(p.cutExcess(seen[:n]))
			return cut
		}
		amount := r.bottleneck(path)
		r.augment(path, amount)
		required.Sub(&required, amount)
	}
	// Auxiliary arcs must not participate in cost optimization.
	for i := firstAuxiliary; i < len(r.arcs); i++ {
		r.arcs[i].capacity.SetInt64(0)
	}
	return nil
}

func (p *circulationProblem) cutExcess(inside []bool) *big.Int {
	excess := new(big.Int)
	for i := range p.nodes {
		if inside[i] {
			excess.Add(excess, &p.balances[i])
		}
	}
	for _, line := range p.lines {
		switch {
		case inside[line.from] && !inside[line.to]:
			if line.bounds.Upper == -1 {
				return nil
			}
			excess.Sub(excess, big.NewInt(line.bounds.Upper))
		case !inside[line.from] && inside[line.to]:
			excess.Add(excess, big.NewInt(line.bounds.Lower))
		}
	}
	return excess
}

// negativeCycle runs Bellman-Ford from an implicit zero-cost source joined
// to every node. Thus disconnected improving circulations are included.
// In the absence of a negative cycle its distances are dual potentials.
func (r *circulationResidual) negativeCycle(n int, unlimitedOnly bool) ([]int, []big.Int) {
	distance := make([]big.Int, n)
	previous := make([]int, n)
	last := -1
	for pass := 0; pass < n; pass++ {
		last = -1
		for i := range r.arcs {
			arc := &r.arcs[i]
			if arc.from >= n || arc.to >= n || (!arc.infinite && arc.capacity.Sign() <= 0) || (unlimitedOnly && !arc.infinite) {
				continue
			}
			var candidate big.Int
			candidate.Add(&distance[arc.from], &arc.cost)
			if candidate.Cmp(&distance[arc.to]) < 0 {
				distance[arc.to].Set(&candidate)
				previous[arc.to] = i
				last = arc.to
			}
		}
		if last == -1 {
			return nil, distance
		}
	}
	if last == -1 {
		return nil, distance
	}
	for range n {
		last = r.arcs[previous[last]].from
	}
	start := last
	var cycle []int
	for {
		i := previous[last]
		cycle = append(cycle, i)
		last = r.arcs[i].from
		if last == start {
			return cycle, nil
		}
	}
}

func (r *circulationResidual) minimize(n int) ([]big.Int, []int) {
	for {
		cycle, potentials := r.negativeCycle(n, false)
		if len(cycle) == 0 {
			return potentials, nil
		}
		amount := r.bottleneck(cycle)
		if amount == nil {
			return nil, cycle
		}
		r.augment(cycle, amount)
	}
}

// VerifyCirculation checks a primal and dual certificate without solving
// another optimization problem. All graph lines and nodes must be present
// exactly once in the respective maps, including zero-flow lines.
// It accepts any optimum and any valid additive normalization of potentials.
// It does not mutate the graph, inputs, or certificate.
func VerifyCirculation(g graph.DirectedMultigraph, bounds func(graph.Line) FlowBounds, balances map[int64]int64, result *Circulation) error {
	p, err := newCirculationProblem(g, bounds, balances)
	if err != nil {
		return err
	}
	if err := p.verifyAssignment(result); err != nil {
		return err
	}
	if len(result.Potentials) != len(p.nodes) {
		return errors.New("network: incomplete dual certificate")
	}
	for _, id := range p.nodes {
		if result.Potentials[id] == nil {
			return fmt.Errorf("network: missing potential for %d", id)
		}
	}
	for _, line := range p.lines {
		flow := result.Flows[line.key]
		var reduced big.Int
		reduced.Sub(result.Potentials[line.key.From], result.Potentials[line.key.To])
		reduced.Add(&reduced, big.NewInt(line.bounds.Cost))
		if (line.bounds.Upper == -1 || flow.Cmp(big.NewInt(line.bounds.Upper)) < 0) && reduced.Sign() < 0 {
			return errors.New("network: negative forward reduced cost")
		}
		if flow.Cmp(big.NewInt(line.bounds.Lower)) > 0 && reduced.Sign() > 0 {
			return errors.New("network: negative reverse reduced cost")
		}
	}
	return nil
}

func (p *circulationProblem) verifyAssignment(result *Circulation) error {
	if result == nil || result.Cost == nil {
		return errors.New("network: missing circulation assignment")
	}
	if len(result.Flows) != len(p.lines) {
		return errors.New("network: incomplete circulation assignment")
	}
	net := make([]big.Int, len(p.nodes))
	var cost big.Int
	for _, line := range p.lines {
		flow := result.Flows[line.key]
		if flow == nil || flow.Cmp(big.NewInt(line.bounds.Lower)) < 0 {
			return errors.New("network: flow below lower bound")
		}
		if line.bounds.Upper != -1 && flow.Cmp(big.NewInt(line.bounds.Upper)) > 0 {
			return errors.New("network: flow above upper bound")
		}
		net[line.from].Add(&net[line.from], flow)
		net[line.to].Sub(&net[line.to], flow)
		var term big.Int
		term.Mul(flow, big.NewInt(line.bounds.Cost))
		cost.Add(&cost, &term)
	}
	for i := range net {
		if net[i].Cmp(&p.balances[i]) != 0 {
			return errors.New("network: certified flow does not conserve balance")
		}
	}
	if cost.Cmp(result.Cost) != 0 {
		return errors.New("network: incorrect certified cost")
	}
	return nil
}

// VerifyInfeasibleCirculation checks a cut certificate. Cut must contain
// distinct existing nodes and its calculated positive excess must equal
// Excess. A different violated cut from the solver's is equally valid.
func VerifyInfeasibleCirculation(g graph.DirectedMultigraph, bounds func(graph.Line) FlowBounds, balances map[int64]int64, certificate *InfeasibleCirculationError) error {
	p, err := newCirculationProblem(g, bounds, balances)
	if err != nil {
		return err
	}
	if certificate == nil || certificate.Excess == nil {
		return errors.New("network: missing infeasibility certificate")
	}
	inside := make([]bool, len(p.nodes))
	for _, id := range certificate.Cut {
		i, ok := p.index[id]
		if !ok || inside[i] {
			return errors.New("network: invalid cut node")
		}
		inside[i] = true
	}
	excess := p.cutExcess(inside)
	if excess == nil || excess.Sign() <= 0 || excess.Cmp(certificate.Excess) != 0 {
		return errors.New("network: invalid cut excess")
	}
	return nil
}

// VerifyUnboundedCirculation checks feasibility and a negative-cost ray.
// Ray may omit zero entries. Its present entries must be nonnegative,
// non-nil, and identify existing lines. A positive entry is allowed only
// on an unlimited line. All arithmetic and verification are exact.
func VerifyUnboundedCirculation(g graph.DirectedMultigraph, bounds func(graph.Line) FlowBounds, balances map[int64]int64, certificate *UnboundedCirculationError) error {
	p, err := newCirculationProblem(g, bounds, balances)
	if err != nil {
		return err
	}
	if certificate == nil {
		return errors.New("network: missing unboundedness certificate")
	}
	if err := p.verifyAssignment(certificate.Feasible); err != nil {
		return err
	}
	byKey := make(map[FlowLine]circulationLine, len(p.lines))
	for _, line := range p.lines {
		byKey[line.key] = line
	}
	net := make([]big.Int, len(p.nodes))
	var cost big.Int
	for key, value := range certificate.Ray {
		line, ok := byKey[key]
		if !ok || value == nil || value.Sign() < 0 {
			return errors.New("network: invalid ray entry")
		}
		if value.Sign() > 0 && line.bounds.Upper != -1 {
			return errors.New("network: bounded line in ray")
		}
		net[line.from].Add(&net[line.from], value)
		net[line.to].Sub(&net[line.to], value)
		var term big.Int
		term.Mul(value, big.NewInt(line.bounds.Cost))
		cost.Add(&cost, &term)
	}
	for i := range net {
		if net[i].Sign() != 0 {
			return errors.New("network: ray is not a circulation")
		}
	}
	if cost.Sign() >= 0 {
		return errors.New("network: ray does not improve cost")
	}
	return nil
}
