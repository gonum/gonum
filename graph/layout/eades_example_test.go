package layout

import (
	"fmt"
	"path/filepath"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
)

func ExampleEadesR2() {

	g := simple.NewUndirectedGraph()

	s1 := makeStar(0, 4)
	s2 := makeStar(4, 4)

	graph.Copy(g, s1)
	graph.Copy(g, s2)

	g.SetEdge(g.NewEdge(simple.Node(0), simple.Node(4)))

	eades := EadesR2{Repulsion: 1, Rate: 0.1, Updates: 30, Theta: 0.1}

	optimizer := NewOptimizerR2(g, eades.Update)

	for optimizer.Update() {

	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Add(render{optimizer})
	p.HideAxes()

	path := filepath.Join("testdata", "graph.png")
	if err := p.Save(10*vg.Centimeter, 10*vg.Centimeter, path); err != nil {
		panic(err)
	}

	fmt.Println("Saved plot to testdata/graph.png")

	// Output: Saved plot to testdata/graph.png
}

func makeStar(start, length int) *simple.UndirectedGraph {
	g := simple.NewUndirectedGraph()

	for n := start + 1; n < start+length; n++ {
		g.SetEdge(g.NewEdge(simple.Node(start), simple.Node(n)))
	}

	return g
}
