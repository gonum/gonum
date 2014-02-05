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
type SinglePathFunc func(start, goal gr.Node) (path []gr.Node, cost float64, err error)

<<<<<<< HEAD
<<<<<<< HEAD
// This function returns two functions: one that will generate all shortest paths between two nodes with ids i and j, and one that will generate just one path.
=======
// This function returns two functions that will generate all shortest paths between two nodes with ids i and j.
>>>>>>> Basic skeleton implementation of FW, tests to follow
=======
// This function returns two functions: one that will generate all shortest paths between two nodes with ids i and j, and one that will generate just one path.
>>>>>>> Basic skeleton implementation of FW, tests to follow
//
// This algorithm requires the CrunchGraph interface which means it only works on nodes with dense ids since it uses an adjacency matrix.
//
// This algorithm isn't blazingly fast, but is relatively fast for the domain. It runs at O((number of vertices)^3), and successfully computes
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> Basic skeleton implementation of FW, tests to follow
// the cost between all pairs of vertices.
//
// Generating a single path should be pretty cheap after FW is done running. The AllPathFunc is likely to be considerably more expensive,
// simply because it has to effectively generate all combinations of known valid paths at each recursive step of the algorithm.
<<<<<<< HEAD
=======
// the cost between all pairs of vertices. Using just a little extra memory, we can remember all shortest paths
>>>>>>> Basic skeleton implementation of FW, tests to follow
=======
>>>>>>> Basic skeleton implementation of FW, tests to follow
func FloydWarshall(graph gr.CrunchGraph, cost func(gr.Node, gr.Node) float64) (AllPathFunc, SinglePathFunc) {
	graph.Crunch()
	_, _, _, _, _, _, cost, _ = setupFuncs(graph, cost, nil)

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

	edges := graph.EdgeList()
	for _, edge := range edges {
		u := edge.Head().ID()
		v := edge.Tail().ID()

		dist[u+v*numNodes] = cost(edge.Head(), edge.Tail())
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
<<<<<<< HEAD
<<<<<<< HEAD
					//
					// NOTE: This may be a straight else, awaiting tests.
=======
>>>>>>> Basic skeleton implementation of FW, tests to follow
=======
					//
					// NOTE: This may be a straight else, awaiting tests.
>>>>>>> Added a note
				} else if math.Abs(dist[i+k*numNodes]+dist[k+j*numNodes]-dist[i+j*numNodes]) < 0.00001 {
					next[i+j*numNodes] = append(next[i+j*numNodes], k)
				}
			}
		}
	}

	return genAllPathsFunc(dist, next, nodes), genSinglePathFunc(dist, next, nodes)
}

func genAllPathsFunc(dist []float64, next [][]int, nodes []gr.Node) func(start, goal gr.Node) ([][]gr.Node, float64, error) {
	numNodes := len(nodes)

	// A recursive function to reconstruct all possible paths.
	// It's not fast, but it's about as fast as can be reasonably expected
	var allPathFinder func(i, j int) ([][]gr.Node, error)
	allPathFinder = func(i, j int) ([][]gr.Node, error) {
		if dist[i+j*numNodes] == math.Inf(1) {
			return nil, errors.New("No path")
		}
		intermediates := next[i+j*numNodes]
		if intermediates == nil {
			return [][]gr.Node{}, nil
		}

		toReturn := make([][]gr.Node, 0, len(intermediates))

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
			return nil, math.Inf(1), nil
		}

		for i := range paths {
			// Prepend start and postpend goal. pathFinder only does the intermediate steps
			paths[i] = append(paths[i], nil)
			copy(paths[i][1:], paths[i][:len(paths[i])-1])
			paths[i][0] = start
			paths[i] = append(paths[i], goal)
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
		if intermediates == nil {
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
			return nil, math.Inf(1), nil
		}

		path = append(path, nil)
		copy(path[1:], path[:len(path)-1])
		path[0] = start
		path = append(path, goal)

		return path, dist[start.ID()+goal.ID()*numNodes], nil
	}
}
