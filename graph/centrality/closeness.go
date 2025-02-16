package centrality

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/traverse"
	"math"
)

// ClosenessCentrality computes closeness centrality for all nodes in an unweighted weightedUndirectedGraph.
func ClosenessCentrality(graph graph.Graph) map[int64]float64 {
	nodes := graph.Nodes()
	numberOfNodes := nodes.Len()
	if numberOfNodes <= 1 {
		return nil
	}

	centrality := make(map[int64]float64)

	for nodes.Next() {
		currentNode := nodes.Node()
		shortestPath := shortestPathsBFS(graph, currentNode)

		totalDistance := 0.0
		for _, distance := range shortestPath {
			totalDistance += float64(distance)
		}

		// Compute closeness centrality
		if totalDistance > 0 {
			centrality[currentNode.ID()] = float64(numberOfNodes-1) / totalDistance
		} else {
			centrality[currentNode.ID()] = 0
		}
	}

	return centrality
}

func shortestPathsBFS(g graph.Graph, startNode graph.Node) map[int64]int {
	distance := make(map[int64]int)
	var t traverse.BreadthFirst
	t.Walk(g, startNode, func(n graph.Node, depth int) bool {
		distance[n.ID()] = depth
		return false
	})
	return distance
}

func ClosenessCentralityWeighted(g graph.Weighted) map[int64]float64 {
	nodes := g.Nodes()
	numberOfNodes := nodes.Len()

	if numberOfNodes <= 1 {
		return nil
	}
	centrality := make(map[int64]float64)

	for nodes.Next() {
		currentNode := nodes.Node()
		shortestPaths := path.DijkstraFrom(currentNode, g)

		// Sum shortest path distances
		totalDistance := 0.0
		for nodes.Reset(); nodes.Next(); {
			target := nodes.Node()
			if target.ID() == currentNode.ID() {
				continue
			}
			distance := shortestPaths.WeightTo(target.ID())
			if distance != math.Inf(1) {
				totalDistance += distance
			}
		}

		// Compute closeness centrality
		if totalDistance > 0 {
			centrality[currentNode.ID()] = float64(numberOfNodes-1) / totalDistance
		} else {
			centrality[currentNode.ID()] = 0
		}
	}

	return centrality
}
