package it

import (
	"gonum.org/v1/gonum/mat"
)

// Emperical1D is an empirical estimator for a one-dimensional
// probability distribution
func Emperical1D(d []int) mat.Vector {
	max := 0
	for _, v := range d {
		if v > max {
			max = v
		}
	}

	max++

	c := make([]float64, max, max)
	l := float64(len(d))

	for _, v := range d {
		c[v] += 1.0
	}

	for i := range c {
		c[i] /= l
	}

	p := mat.NewVecDense(max, c)

	return p
}

// Emperical2D is an empirical estimator for a two-dimensional
// probability distribution
func Emperical2D(d [][]int) [][]float64 {
	max := make([]int, 2, 2)
	rows := len(d)
	for r := 0; r < rows; r++ {
		for c := 0; c < 2; c++ {
			if d[r][c] > max[c] {
				max[c] = d[r][c]
			}
		}
	}

	max[0]++
	max[1]++

	p := make([][]float64, max[0], max[0])
	for m := 0; m < max[0]; m++ {
		p[m] = make([]float64, max[1], max[1])
	}

	for r := 0; r < rows; r++ {
		p[d[r][0]][d[r][1]] += 1.0
	}

	l := float64(len(d))
	for r := 0; r < max[0]; r++ {
		for c := 0; c < max[1]; c++ {
			p[r][c] /= l
		}
	}

	return p
}

// Emperical3D is an empirical estimator for a three-dimensional
// probability distribution
func Emperical3D(d [][]int) [][][]float64 {
	max := make([]int, 3, 3)
	rows := len(d)
	for r := 0; r < rows; r++ {
		for c := 0; c < 3; c++ {
			if d[r][c] > max[c] {
				max[c] = d[r][c]
			}
		}
	}

	max[0]++
	max[1]++
	max[2]++

	p := make([][][]float64, max[0], max[0])
	for m := 0; m < max[0]; m++ {
		p[m] = make([][]float64, max[1], max[1])
		for n := 0; n < max[1]; n++ {
			p[m][n] = make([]float64, max[2], max[2])
		}
	}

	for r := 0; r < rows; r++ {
		p[d[r][0]][d[r][1]][d[r][2]] += 1.0
	}

	l := float64(len(d))
	for a := 0; a < max[0]; a++ {
		for b := 0; b < max[1]; b++ {
			for c := 0; c < max[2]; c++ {
				p[a][b][c] /= l
			}
		}
	}

	return p
}

// Emperical4D is an empirical estimator for a three-dimensional
// probability distribution
func Emperical4D(d [][]int) [][][][]float64 {
	max := make([]int, 4, 4)
	rows := len(d)
	for r := 0; r < rows; r++ {
		for c := 0; c < 4; c++ {
			if d[r][c] > max[c] {
				max[c] = d[r][c]
			}
		}
	}

	max[0]++
	max[1]++
	max[2]++
	max[3]++

	p := make([][][][]float64, max[0], max[0])
	for m := 0; m < max[0]; m++ {
		p[m] = make([][][]float64, max[1], max[1])
		for n := 0; n < max[1]; n++ {
			p[m][n] = make([][]float64, max[2], max[2])
			for k := 0; k < max[2]; k++ {
				p[m][n][k] = make([]float64, max[3], max[3])
			}
		}
	}

	for r := 0; r < rows; r++ {
		p[d[r][0]][d[r][1]][d[r][2]][d[r][3]] += 1.0
	}

	l := float64(len(d))
	for a := 0; a < max[0]; a++ {
		for b := 0; b < max[1]; b++ {
			for c := 0; c < max[2]; c++ {
				for d := 0; d < max[3]; d++ {
					p[a][b][c][d] /= l
				}
			}
		}
	}

	return p
}
