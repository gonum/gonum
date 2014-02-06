package search

import (
	"errors"
	"math"
	"sort"

	gr "github.com/gonum/graph"
)

// Finds all shortest paths between start and goal
type AllPathFunc func(start, goal gr.Node) (path [][]gr.Node, cost float64, err error)

// Finds one path between start and goal, which it finds is arbitrary
type PathFunc func(start, goal gr.Node) (path []gr.Node, cost float64, err error)

// This function returns two functions: one that will generate all shortest paths between two nodes with ids i and j, and one that will generate just one path.
//
// This algorithm requires the CrunchGraph interface which means it only works on graphs with dense node ids since it uses an adjacency matrix.
//
// This algorithm isn't blazingly fast, but is relatively fast for the domain. It runs at O((number of vertices)^3) in best, worst, and average case,
//  and successfully computes the cost between all pairs of vertices.
//
// This function operates slightly differently from the others for convenience -- rather than generating paths and returning them to you,
// it gives you the option of calling one of two functions for each start/goal pair you need info for. One will return the path, cost,
// or an error if no path exists.
//
// The other will return the cost and an error if no path exists, but it will also return ALL possible shortest paths between start and goal.
// This is not too much more expensive than generating one path, but it does obviously increase with the number of paths.
func FloydWarshall(graph gr.CrunchGraph, cost gr.CostFun) (AllPathFunc, PathFunc) {
	graph.Crunch()
	sf := setupFuncs(graph, cost, nil)
	successors, isSuccessor, cost := sf.successors, sf.isSuccessor, sf.cost

	nodes := denseNodeSorter(graph.NodeList())
	sort.Sort(nodes)
	numNodes := len(nodes)

	dist := make([]float64, numNodes*numNodes)
	next := make([][]int, numNodes*numNodes)
	for i := 0; i < numNodes; i++ {
		for j := 0; j < numNodes; j++ {
			if j != i {
				dist[i+j*numNodes] = math.Inf(1)
			}
		}
	}

	for _, node := range nodes {
		for _, succ := range successors(node) {
			dist[node.ID()+succ.ID()*numNodes] = cost(node, succ)
		}
	}

	for k := 0; k < numNodes; k++ {
		for i := 0; i < numNodes; i++ {
			for j := 0; j < numNodes; j++ {
				if dist[i+j*numNodes] > dist[i+k*numNodes]+dist[k+j*numNodes] {
					dist[i+j*numNodes] = dist[i+k*numNodes] + dist[k+j*numNodes]

					// Avoid generating too much garbage by reusing the memory in the list if we've allocated one already
					if next[i+j*numNodes] == nil {
						next[i+j*numNodes] = []int{k}
					} else {
						next[i+j*numNodes] = next[i+j*numNodes][:1]
						next[i+j*numNodes][0] = k
					}
					// If the cost between the nodes happens to be the same cost as what we know, add the approriate
					// intermediary to the list
				} else if math.Abs(dist[i+k*numNodes]+dist[k+j*numNodes]-dist[i+j*numNodes]) < 0.00001 && i != k && i != j && j != k {
					next[i+j*numNodes] = append(next[i+j*numNodes], k)
				}
			}
		}
	}

	return genAllPathsFunc(dist, next, nodes, graph, cost, isSuccessor), genSinglePathFunc(dist, next, nodes)
}

func genAllPathsFunc(dist []float64, next [][]int, nodes []gr.Node, graph gr.Graph, cost func(gr.Node, gr.Node) float64, isSuccessor func(gr.Node, gr.Node) bool) func(start, goal gr.Node) ([][]gr.Node, float64, error) {
	numNodes := len(nodes)

	// A recursive function to reconstruct all possible paths.
	// It's not fast, but it's about as fast as can be reasonably expected
	var allPathFinder func(i, j int) ([][]gr.Node, error)
	allPathFinder = func(i, j int) ([][]gr.Node, error) {
		if dist[i+j*numNodes] == math.Inf(1) {
			return nil, errors.New("No path")
		}
		intermediates := next[i+j*numNodes]
		if intermediates == nil || len(intermediates) == 0 {
			// There is exactly one path
			return [][]gr.Node{[]gr.Node{}}, nil
		}

		toReturn := make([][]gr.Node, 0, len(intermediates))
		// Special case: if intermediates exist we need to explicitly check to see if i and j is also an optimal path
		if isSuccessor(nodes[i], nodes[j]) && math.Abs(dist[i+j*numNodes]-cost(nodes[i], nodes[j])) < .000001 {
			toReturn = append(toReturn, []gr.Node{})
		}

		// This step is a tad convoluted: we have some list of intermediates.
		// We can think of each intermediate as a path junction
		//
		// At this junction, we can find all the shortest paths back to i,
		// and all the shortest paths down to j. Since this is a junction,
		// any predecessor path that runs through this intermediate may
		// freely choose any successor path to get to j. They'll all be
		// of equivalent length.
		//
		// Thus, for each intermediate, we run through and join each predecessor
		// path with each successor path via its junction.
		for _, intermediate := range intermediates {

			// Find predecessors
			preds, err := allPathFinder(i, intermediate)
			if err != nil {
				return nil, err
			}

			// Join each predecessor with its junction
			for a := range preds {
				preds[a] = append(preds[a], nodes[intermediate])
			}

			// Find successors
			succs, err := allPathFinder(intermediate, j)
			if err != nil {
				return nil, err
			}

			// Join each successor with its predecessor at the junction.
			// (the copying stuff is because slices are reference types)
			for a := range succs {
				for b := range preds {
					path := make([]gr.Node, len(succs[a]), len(succs[a])+len(preds[b]))
					copy(path, succs[a])
					path = append(path, preds[b]...)
					toReturn = append(toReturn, path)
				}
			}

		}

		return toReturn, nil
	}

	return func(start, goal gr.Node) ([][]gr.Node, float64, error) {
		paths, err := allPathFinder(start.ID(), goal.ID())
		if err != nil {
			return nil, math.Inf(1), err
		}

		for i := range paths {
			// Prepend start and postpend goal, but don't repeat start/goal

			if len(paths[i]) != 0 {
				if paths[i][0].ID() != start.ID() {
					paths[i] = append(paths[i], nil)
					copy(paths[i][1:], paths[i][:len(paths[i])-1])
					paths[i][0] = start
				}

				if paths[i][len(paths[i])-1].ID() != goal.ID() {
					paths[i] = append(paths[i], goal)
				}
			} else {
				paths[i] = append(paths[i], start, goal)
			}
		}

		return paths, dist[start.ID()+goal.ID()*numNodes], nil
	}
}

func genSinglePathFunc(dist []float64, next [][]int, nodes []gr.Node) func(start, goal gr.Node) ([]gr.Node, float64, error) {
	numNodes := len(nodes)

	var singlePathFinder func(i, j int) ([]gr.Node, error)
	singlePathFinder = func(i, j int) ([]gr.Node, error) {
		if dist[i+j*numNodes] == math.Inf(1) {
			return nil, errors.New("No path")
		}

		intermediates := next[i+j*numNodes]
		if intermediates == nil || len(intermediates) == 0 {
			return []gr.Node{}, nil
		}

		intermediate := intermediates[0]
		path, err := singlePathFinder(i, intermediate)
		if err != nil {
			return nil, err
		}

		path = append(path, nodes[intermediate])
		p, err := singlePathFinder(intermediate, j)
		if err != nil {
			return nil, err
		}
		path = append(path, p...)

		return path, nil
	}

	return func(start, goal gr.Node) ([]gr.Node, float64, error) {
		path, err := singlePathFinder(start.ID(), goal.ID())
		if err != nil {
			return nil, math.Inf(1), err
		}

		path = append(path, nil)
		copy(path[1:], path[:len(path)-1])
		path[0] = start
		path = append(path, goal)

		return path, dist[start.ID()+goal.ID()*numNodes], nil
	}
}
