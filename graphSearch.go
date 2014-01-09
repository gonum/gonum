package graph

import (
	"container/heap"
	"github.com/gonum/graph/set"
	"github.com/gonum/graph/xifo"
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
func AStar(start, goal Node, graph Graph, Cost, HeuristicCost func(Node, Node) float64) (path []Node, cost float64, nodesExpanded int) {
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
	predecessor := make(map[int]Node)

	for openSet.Len() != 0 {
		curr := heap.Pop(openSet).(internalNode)

		// This isn't in most implementations of A*, it's a restructuring of the step "if node not in openSet, add it"
		// Instead of searching to check, we see if we already evaluated it. If we have we can ignore it
		if _, ok := closedSet[curr.ID()]; ok {
			continue
		}

		nodesExpanded += 1

		if curr.ID() == goal.ID() {
			return rebuildPath(predecessor, goal), curr.gscore, nodesExpanded
		}

		closedSet[curr.ID()] = curr

		for _, neighbor := range graph.Successors(curr.Node) {
			g := curr.gscore + Cost(curr.Node, neighbor)
			if _, ok := closedSet[neighbor.ID()]; ok && g >= closedSet[neighbor.ID()].gscore {
				continue
			}

			if _, ok := closedSet[neighbor.ID()]; !ok || g < closedSet[neighbor.ID()].gscore {
				node = internalNode{neighbor, g, g + HeuristicCost(neighbor, goal)}
				predecessor[node.ID()] = curr
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
func Dijkstra(source Node, graph Graph, Cost func(Node, Node) float64) (paths map[int][]Node, costs map[int]float64) {
	if Cost == nil {
		if cgraph, ok := graph.(Coster); ok {
			Cost = cgraph.Cost
		} else {
			Cost = UniformCost
		}
	}
	nodes := graph.NodeList()
	openSet := &aStarPriorityQueue{}
	closedSet := set.NewSet()                 // This is to make use of that same
	costs = make(map[int]float64, len(nodes)) // May overallocate, will change if it becomes a problem
	predecessor := make(map[int]Node, len(nodes))
	nodeIDMap := make(map[int]Node, len(nodes))
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

	costs[source.ID()] = 0
	heap.Push(openSet, internalNode{source, 0, 0})

	for openSet.Len() != 0 {
		node := heap.Pop(openSet).(internalNode)
		/* if _, ok := costs[node.int]; !ok {
			 break
		 } */

		if closedSet.Contains(node.ID()) { // As in A*, prevents us from having to slowly search and reorder the queue
			continue
		}

		nodeIDMap[node.ID()] = node

		closedSet.Add(node.ID())

		for _, neighbor := range graph.Successors(node) {
			tmpCost := costs[node.ID()] + Cost(node, neighbor)
			if cost, ok := costs[neighbor.ID()]; !ok || tmpCost < cost {
				costs[neighbor.ID()] = cost
				predecessor[neighbor.ID()] = node
				heap.Push(openSet, internalNode{neighbor, cost, cost})
			}
		}
	}

	paths = make(map[int][]Node, len(costs))
	for node, _ := range costs { // Only reconstruct the path if one exists
		paths[node] = rebuildPath(predecessor, nodeIDMap[node])
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
// due to the presence of a negative edge weight cycle.
func BellmanFord(source Node, graph Graph, Cost func(Node, Node) float64) (paths map[int][]Node, costs map[int]float64, aborted bool) {
	if Cost == nil {
		if cgraph, ok := graph.(Coster); ok {
			Cost = cgraph.Cost
		} else {
			Cost = UniformCost
		}
	}

	predecessor := make(map[int]Node)
	costs = make(map[int]float64)
	nodeIDMap := make(map[int]Node)
	nodeIDMap[source.ID()] = source
	costs[source.ID()] = 0
	nodes := graph.NodeList()
	edges := graph.EdgeList()

	for i := 1; i < len(nodes)-1; i++ {
		for _, edge := range edges {
			weight := Cost(edge.Head(), edge.Tail())
			nodeIDMap[edge.Head().ID()] = edge.Head()
			nodeIDMap[edge.Tail().ID()] = edge.Tail()
			if dist := costs[edge.Head().ID()] + weight; dist < costs[edge.Tail().ID()] {
				costs[edge.Tail().ID()] = dist
				predecessor[edge.Tail().ID()] = edge.Head()
			}
		}
	}

	for _, edge := range edges {
		weight := Cost(edge.Head(), edge.Tail())
		if costs[edge.Head().ID()]+weight < costs[edge.Tail().ID()] {
			return nil, nil, true // Abandoned because a cycle is detected
		}
	}

	paths = make(map[int][]Node, len(costs))
	for node, _ := range costs {
		paths[node] = rebuildPath(predecessor, nodeIDMap[node])
	}
	return paths, costs, false
}

// Johnson's Algorithm generates the lowest cost path between every pair of nodes in the graph.
//
// It makes use of Bellman-Ford and a dummy graph. It creates a dummy node containing edges with a cost of zero to every other node. Then it runs Bellman-Ford with this
// dummy node as the source. It then modifies the all the nodes' edge weights (which gets rid of all negative weights).
//
// Finally, it removes the dummy node and runs Dijkstra's starting at every node.
//
// This algorithm is fairly slow. Its purpose is to remove negative edge weights to allow Dijkstra's to function properly. It's probably not worth it to run this algorithm if you have
// all non-negative edge weights. Also note that this implementation copies your whole graph into a GonumGraph (so it can add/remove the dummy node and edges and reweight the graph).
//
// Its return values are, in order: a map from the source node, to the destination node, to the path between them; a map from the source node, to the destination node, to the cost of the path between them;
// and a bool that is true if Bellman-Ford detected a negative edge weight cycle -- thus causing it (and this algorithm) to abort (if aborted is true, both maps will be nil).
func Johnson(graph Graph, Cost func(Node, Node) float64) (nodePaths map[int]map[int][]Node, nodeCosts map[int]map[int]float64, aborted bool) {
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
				dummyGraph.SetEdgeCost(GonumEdge{node, neighbor}, Cost(node, neighbor))
			}
		} else {
			for _, neighbor := range neighbors {
				dummyGraph.AddEdge(GonumEdge{node, neighbor})
				dummyGraph.SetEdgeCost(GonumEdge{node, neighbor}, Cost(node, neighbor))
			}
		}
	}

	/* Step 1: Dummy node with 0 cost edge weights to every other node*/
	dummyNode := dummyGraph.NewNode(graph.NodeList())
	for _, node := range graph.NodeList() {
		dummyGraph.SetEdgeCost(GonumEdge{dummyNode, node}, 0)
	}

	/* Step 2: Run Bellman-Ford starting at the dummy node, abort if it detects a cycle */
	_, costs, aborted := BellmanFord(dummyNode, dummyGraph, nil)
	if aborted {
		return nil, nil, true
	}

	/* Step 3: reweight the graph and remove the dummy node */
	for _, edge := range graph.EdgeList() {
		dummyGraph.SetEdgeCost(edge, Cost(edge.Head(), edge.Tail())+costs[edge.Head().ID()]-costs[edge.Tail().ID()])
	}

	dummyGraph.RemoveNode(dummyNode)

	/* Step 4: Run Dijkstra's starting at every node */
	nodePaths = make(map[int]map[int][]Node, len(graph.NodeList()))
	nodeCosts = make(map[int]map[int]float64)

	for _, node := range graph.NodeList() {
		nodePaths[node.ID()], nodeCosts[node.ID()] = Dijkstra(node, dummyGraph, nil)
	}

	return nodePaths, nodeCosts, false
}

// Expands the first node it sees trying to find the destination. Depth First Search is *not* guaranteed to find the shortest path,
// however, if a path exists DFS is guaranteed to find it (provided you don't find a way to implement a Graph with an infinite depth)
func DepthFirstSearch(start, goal Node, graph Graph) []Node {
	closedSet := set.NewSet()
	openSet := xifo.GonumStack([]interface{}{start})
	predecessor := make(map[int]Node)

	for !openSet.IsEmpty() {
		c := openSet.Pop()

		curr := c.(Node)

		if closedSet.Contains(curr.ID()) {
			continue
		}

		if curr == goal {
			return rebuildPath(predecessor, goal)
		}

		closedSet.Add(curr.ID())

		for _, neighbor := range graph.Successors(curr) {
			if closedSet.Contains(neighbor.ID()) {
				continue
			}

			predecessor[neighbor.ID()] = curr
			openSet.Push(neighbor)
		}
	}

	return nil
}

// An admissible, consistent heuristic that won't speed up computation time at all.
func NullHeuristic(a, b Node) float64 {
	return 0.0
}

// Assumes all edges in the graph have the same weight (including edges that don't exist!)
func UniformCost(a, b Node) float64 {
	return 1.0
}

/** Keeps track of a node's scores so they can be used in a priority queue for A* **/

type internalNode struct {
	Node
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
func rebuildPath(predecessors map[int]Node, goal Node) []Node {
	path := []Node{goal}
	curr := goal
	for prev, ok := predecessors[curr.ID()]; ok; prev, ok = predecessors[curr.ID()] {
		path = append(path, prev)
		curr = prev
	}

	// Reverse the path since it was built backwards
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}
