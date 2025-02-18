package centrality

import (
	"errors"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/traverse"
	"math"
)

type GraphWithEdgesAndWeights interface {
	graph.Weighted
	Edges() graph.Edges
}

// ClosenessCentrality computes normalized closeness centrality for all nodes in an unweighted undirected Graph.
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

		// Compute normalized closeness centrality
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

func ValidateNonNegativeWeights(graph GraphWithEdgesAndWeights) error {
	edges := graph.Edges()
	for edges.Next() {
		edge := edges.Edge()
		weight, _ := graph.Weight(edge.From().ID(), edge.To().ID())
		if weight < 0 {
			return errors.New("graph contains negative edge weights")
		}
	}
	return nil
}

// ClosenessCentralityWeighted computes normalized closeness centrality for all nodes in a weighted undirected Graph.
func ClosenessCentralityWeighted(graph GraphWithEdgesAndWeights) (map[int64]float64, error) {
	if err := ValidateNonNegativeWeights(graph); err != nil {
		return nil, err
	}

	nodes := graph.Nodes()
	numberOfNodes := nodes.Len()

	if numberOfNodes <= 1 {
		return make(map[int64]float64), nil
	}
	centrality := make(map[int64]float64)

	for nodes.Next() {
		currentNode := nodes.Node()
		shortestPaths := path.DijkstraFrom(currentNode, graph)

		// Sum shortest path distances
		totalDistance := 0.0
		neighbors := graph.Nodes()
		for neighbors.Next() {
			target := neighbors.Node()
			if target.ID() == currentNode.ID() {
				continue
			}
			distance := shortestPaths.WeightTo(target.ID())
			if distance != math.Inf(1) {
				totalDistance += distance
			}
		}

		if totalDistance > 0 {
			centrality[currentNode.ID()] = float64(numberOfNodes-1) / totalDistance
		} else {
			centrality[currentNode.ID()] = 0
		}
	}

	return centrality, nil
}
