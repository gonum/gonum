package discrete

import (
	"container/heap"
)

// Returns an ordered list consisting of the nodes between start and goal. The path will be the shortest path assuming the function heuristicCost is admissible.
// The second return value is the cost, and the third is the number of nodes expanded while searching (useful info for tuning heuristics). Negative Costs will cause
// bad things to happen, as well as negative heuristic estimates.
//
// A heuristic is admissible if, for any node in the graph, the heuristic estimate of the cost between the node and the goal is less than or equal to the true cost.
//
// Performance may be improved by providing a consistent heuristic (though one is not needed to find the optimal path), a heuristic is consistent if its value for a given node is less than (or equal to) the
// actual cost of reaching its neighbors + the heuristic estimate for the neighbor itself. You can force consistency by making your HeuristicCost function
// return max(NonConsistentHeuristicCost(neighbor,goal), NonConsistentHeuristicCost(self,goal) - Cost(self,neighbor)). If there are multiple neighbors, take the max of all of them.
//
// Cost and HeuristicCost take precedence for evaluating cost/heuristic distance. If one is not present (i.e. nil) the function will check the graph's interface for the respective interface:
// Coster for Cost and HeuristicCoster for HeuristicCost. If the correct one is present, it will use the graph's function for evaluation.
//
// Finally, if neither the argument nor the interface is present, the function will assume discrete.UniformCost for Cost and discrete.NullHeuristic for HeuristicCost
//
// To run Uniform Cost Search, run A* with the NullHeuristic
//
// To run Breadth First Search, run A* with both the NullHeuristic and UniformCost (or any cost function that returns a uniform positive value)
func AStar(start, goal int, graph Graph, Cost, HeuristicCost func(int, int) float64) (path []int, cost float64, nodesExpanded int) {
	if Cost == nil {
		if cgraph, ok := graph.(Coster); ok {
			Cost = cgraph.Cost
		} else {
			Cost = UniformCost
		}
	}
	if HeuristicCost == nil {
		if hgraph, ok := graph.(HeuristicCoster); ok {
			HeuristicCost = hgraph.HeuristicCost
		} else {
			HeuristicCost = NullHeuristic
		}
	}

	closedSet := make(map[int]internalNode)
	openSet := &aStarPriorityQueue{}
	heap.Init(openSet)
	node := internalNode{start, 0, HeuristicCost(start, goal)}
	heap.Push(openSet, node)
	predecessor := make(map[int]int)

	for openSet.Len() != 0 {
		curr := heap.Pop(openSet).(internalNode)

		// This isn't in most implementations of A*, it's a restructuring of the step "if node not in openSet, add it"
		// Instead of searching to check, we see if we already evaluated it. If we have we can ignore it
		if _, ok := closedSet[curr.int]; ok {
			continue
		}

		nodesExpanded += 1

		if curr.int == goal {
			return rebuildPath(predecessor, goal), curr.gscore, nodesExpanded
		}

		closedSet[curr.int] = curr

		for _, neighbor := range graph.Successors(curr.int) {
			g := curr.gscore + Cost(curr.int, neighbor)
			if _, ok := closedSet[neighbor]; ok && g >= closedSet[neighbor].gscore {
				continue
			}

			if _, ok := closedSet[neighbor]; !ok || g < closedSet[neighbor].gscore {
				node = internalNode{neighbor, g, g + HeuristicCost(neighbor, goal)}
				predecessor[node.int] = curr.int
				heap.Push(openSet, node)
			}
		}
	}

	return nil, 0.0, nodesExpanded
}

// Dijkstra's Algorithm is essentially a goalless Uniform Cost Search. That is, its results are roughly equivalent to
// running A* with the Null Heuristic from a single node to every other node in the graph -- though it's a fair bit faster
// because running A* in that way will recompute things it's already computed every call. Note that you won't necessarily get the same path
// you would get for A*, but the cost is guaranteed to be the same (that is, if multiple shortest paths exist, you may get a different shortest path).
//
// Like A*, Dijkstra's Algorithm likely won't run correctly with negative edge weights -- use Bellman-Ford for that instead
//
// Dijkstra's algorithm usually only returns a cost map, however, since the data is available this version will also reconstruct the path to every node
func Dijkstra(source int, graph Graph, Cost func(int, int) float64) (paths map[int][]int, costs map[int]float64) {
	if Cost == nil {
		if cgraph, ok := graph.(Coster); ok {
			Cost = cgraph.Cost
		} else {
			Cost = UniformCost
		}
	}
	nodes := graph.NodeList()
	openSet := &aStarPriorityQueue{}
	closedSet := NewSet()                     // This is to make use of that same
	costs = make(map[int]float64, len(nodes)) // May overallocate, will change if it becomes a problem
	predecessor := make(map[int]int, len(nodes))
	heap.Init(openSet)

	// I don't think we actually need the init step since I use a map check rather than inf to check if we're done
	/*for _, node := range nodes {
		if node == source {
			heap.Push(openSet, internalNode{node, 0, 0})
			costs[node] = 0
		} else {
			heap.Push(openSet, internalNode{node, math.MaxFloat64, math.MaxFloat64})
			predecessor[node] = -1
		}
	}*/

	costs[source] = 0
	heap.Push(openSet, internalNode{source, 0, 0})

	for openSet.Len() != 0 {
		node := heap.Pop(openSet).(internalNode)
		/* if _, ok := costs[node.int]; !ok {
			 break
		 } */

		if closedSet.Contains(node.int) { // As in A*, prevents us from having to slowly search and reorder the queue
			continue
		}

		closedSet.Add(node.int)

		for _, neighbor := range graph.Successors(node.int) {
			tmpCost := costs[node.int] + Cost(node.int, neighbor)
			if cost, ok := costs[neighbor]; !ok || tmpCost < cost {
				costs[neighbor] = cost
				predecessor[neighbor] = node.int
				heap.Push(openSet, internalNode{neighbor, cost, cost})
			}
		}
	}

	paths = make(map[int][]int, len(costs))
	for node, _ := range costs { // Only reconstruct the path if one exists
		paths[node] = rebuildPath(predecessor, node)
	}
	return paths, costs
}

// The Bellman-Ford Algorithm is the same as Dijkstra's Algorithm with a key difference. They both take a single source and find the shortest path to every other
// (reachable) node in the graph. Bellman-Ford, however, will detect negative edge loops and abort if one is present. A negative edge loop occurs when there is a cycle in the graph
// such that it can take an edge with a negative cost over and over. A -(-2)> B -(2)> C isn't a loop because A->B can only be taken once, but A<-(-2)->B-(2)>C is one because
// A and B have a bi-directional edge, and algorithms like Dijkstra's will infinitely flail between them getting progressively lower costs.
//
// That said, if you do not have a negative edge weight, use Dijkstra's Algorithm instead, because it's faster.
//
// Like Dijkstra's, along with the costs this implementation will also construct all the paths for you. In addition, it has a third return value which will be true if the algorithm was aborted
// due to the presence of a negative edge weight.
func BellmanFord(source int, graph Graph, Cost func(int, int) float64) (paths map[int][]int, costs map[int]float64, aborted bool) {
	if Cost == nil {
		if cgraph, ok := graph.(Coster); ok {
			Cost = cgraph.Cost
		} else {
			Cost = UniformCost
		}
	}

	predecessor := make(map[int]int)
	costs = make(map[int]float64)
	costs[source] = 0
	nodes := graph.NodeList()
	edges := graph.EdgeList()

	for i := 1; i < len(nodes)-1; i++ {
		for _, edge := range edges {
			weight := Cost(edge[0], edge[1])
			if dist := costs[edge[0]] + weight; dist < costs[edge[1]] {
				costs[edge[1]] = dist
				predecessor[edge[1]] = edge[0]
			}
		}
	}

	for _, edge := range edges {
		weight := Cost(edge[0], edge[1])
		if costs[edge[0]]+weight < costs[edge[1]] {
			return nil, nil, true // Abandoned because a cycle is detected
		}
	}

	paths = make(map[int][]int, len(costs))
	for node, _ := range costs {
		paths[node] = rebuildPath(predecessor, node)
	}
	return paths, costs, false
}

// Johnson's Algorithm generates the lowest cost path between every pair of nodes in the graph.
//
// It makes use of Bellman-Ford and a dummy graph. It creates a dummy node containing edges with a cost of zero to every other node. Then it runs Bellman-Ford with this
// dummy node as the source.
//
// Finally, it removes the dummy node and run Dijkstra's starting at every node.
//
// This algorithm is fairly slow. Its purpose is to remove negative edge weights to allow Dijkstra's to function properly. It's probably not worth it to run this algorithm if you have
// all non-negative edge weights. Also note that this implementation copies your whole graph into a GonumGraph (so it can add/remove the dummy node and edges and reweight the graph).
//
// Its return values are, in order: a map from the source node, to the destination node, to the path between them; a map from the source node, to the destination node, to the cost of the path between them;
// and a bool that is true if Bellman-Ford detected a negative edge weight cycle -- thus causing it (and this algorithm) to abort (if aborted is true, both maps will be nil).
func Johnson(graph Graph, Cost func(int, int) float64) (nodePaths map[int]map[int][]int, nodeCosts map[int]map[int]float64, aborted bool) {
	if Cost == nil {
		if cgraph, ok := graph.(Coster); ok {
			Cost = cgraph.Cost
		} else {
			Cost = UniformCost
		}
	}
	/* Copy graph into a mutable one since it has to be altered for this algorithm */
	dummyGraph := NewGonumGraph(true)
	for _, node := range graph.NodeList() {
		neighbors := graph.Successors(node)
		if !dummyGraph.NodeExists(node) {
			dummyGraph.AddNode(node, neighbors)
			for _, neighbor := range neighbors {
				dummyGraph.SetEdgeCost(node, neighbor, Cost(node, neighbor))
			}
		} else {
			for _, neighbor := range neighbors {
				dummyGraph.AddEdge(node, neighbor)
				dummyGraph.SetEdgeCost(node, neighbor, Cost(node, neighbor))
			}
		}
	}

	/* Step 1: Dummy node with 0 cost edge weights to every other node*/
	dummyNode := dummyGraph.NewNode(graph.NodeList())
	for _, node := range graph.NodeList() {
		dummyGraph.SetEdgeCost(dummyNode, node, 0)
	}

	/* Step 2: Run Bellman-Ford starting at the dummy node, abort if it detects a cycle */
	_, costs, aborted := BellmanFord(dummyNode, dummyGraph, nil)
	if aborted {
		return nil, nil, true
	}

	/* Step 3: reweight the graph and remove the dummy node */
	for _, edge := range graph.EdgeList() {
		dummyGraph.SetEdgeCost(edge[0], edge[1], Cost(edge[0], edge[1])+costs[edge[0]]-costs[edge[1]])
	}

	dummyGraph.RemoveNode(dummyNode)

	/* Step 4: Run Dijkstra's starting at every node */
	nodePaths = make(map[int]map[int][]int, len(graph.NodeList()))
	nodeCosts = make(map[int]map[int]float64)

	for _, node := range graph.NodeList() {
		nodePaths[node], nodeCosts[node] = Dijkstra(node, dummyGraph, nil)
	}

	return nodePaths, nodeCosts, false
}

// Expands the first node it sees trying to find the destination. Depth First Search is *not* guaranteed to find the shortest path,
// however, if a path exists DFS is guaranteed to find it (provided you don't find a way to implement a Graph with an infinite depth)
func DepthFirstSearch(start, goal int, graph Graph) []int {
	closedSet := NewSet()
	openSet := Stack([]interface{}{start})
	predecessor := make(map[int]int)

	for !openSet.IsEmpty() {
		c, err := openSet.Pop()
		if err != nil {
			return nil
		}

		curr := c.(int)

		if closedSet.Contains(curr) {
			continue
		}

		if curr == goal {
			return rebuildPath(predecessor, goal)
		}

		closedSet.Add(curr)

		for _, neighbor := range graph.Successors(curr) {
			if closedSet.Contains(neighbor) {
				continue
			}

			predecessor[neighbor] = curr
			openSet.Push(neighbor)
		}
	}

	return nil
}

// An admissible, consistent heuristic that won't speed up computation time at all.
func NullHeuristic(a, b int) float64 {
	return 0.0
}

// Assumes all edges in the graph have the same weight (including edges that don't exist!)
func UniformCost(a, b int) float64 {
	return 1.0
}

/** Keeps track of a node's scores so they can be used in a priority queue for A* **/

type internalNode struct {
	int
	gscore, fscore float64
}

type aStarPriorityQueue []internalNode

func (pq *aStarPriorityQueue) Less(i, j int) bool {
	return (*pq)[i].fscore < (*pq)[j].fscore // As the heap documentation says, a priority queue is listed if the actual values are treated as if they were negative
}

func (pq *aStarPriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
}

func (pq *aStarPriorityQueue) Len() int {
	return len(*pq)
}

func (pq *aStarPriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(internalNode))
}

func (pq *aStarPriorityQueue) Pop() interface{} {
	x := (*pq)[len(*pq)-1]
	(*pq) = (*pq)[:len(*pq)-1]

	return x
}

// Rebuilds a path backwards from the goal.
func rebuildPath(predecessors map[int]int, goal int) []int {
	path := []int{goal}
	curr := goal
	for prev, ok := predecessors[curr]; ok; prev, ok = predecessors[curr] {
		path = append([]int{prev}, path...) // Maybe do something better than prepending?
		curr = prev
	}

	return path
}
