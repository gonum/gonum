package optimize

/*
import "github.com/gonum/floats"

// NelderMead implements the Nelder-Mead simplex algorithm. Nelder-Mead is a
// gradient-free optimization technique which sequentially refines a simplex
// to find a local minimum of a function.
type NelderMead struct {
	vertices []Location
	centroid []float64

	// Initial sets the initial exploration vertices. If Initial is nil
	// a random initialization will be used.
	Initial []Location

	Reflection  float64 // alpha
	Expansion   float64 // gamma
	Contraction float64 // rho
	Shrink      float64 //sigma
}

func (n *NelderMead) Init(l Location, f *FunctionInfo, xNext []float64) (EvaluationType, IterationType, error) {
	dim := len(l.X)
	if cap(n.vertices) < dim {
		n.vertices = append(n.vertices, dim+1-len(n.vertices))
	} else {
		n.vertices = n.vertices[:n]
	}
	for i := range n.vertices {
		n.vertices[i] = resize(n.vertices[i], dim)
	}
	n.centroid = resize(n.centroid, dim)
	n.f = resize(n.f, dim)

	if n.Reflection == 0 {
		n.Reflection = 1
	}
	if n.Expansion == 0 {
		n.Expansion = 2
	}
	if n.Contraction == 0 {
		n.Contraction = -0.5
	}
	if n.Shrink == 0 {
		n.Shrink = 0.5
	}

	// TODO: Add something about initialization of simplex
}

func (n *NelderMead) Iterate(l Location, xNext []float64) (EvaluationType, IterationType, error) {
	// TODO: Code about if still during initialization
	// When done initialization, need to have sorted

	// Compute centroid of all but the best point
	copy(n.centroid, n.vertices[0])
	for i := 1; i < len(n.vertices)-1; i++ {
		floats.AddScaled(n.centroid, 1, n.vertices[i])
	}
	floats.Scale(n.centroid, 1/(len(n.vertices)-1))

}

/*
// Vertices returns a copy of the vertices of the simplex. The vertex locations
// will be copied into dst. If dst is nil, new slices will be allocated.
func (n *NelderMead) Vertices() []Location {

}
*/
