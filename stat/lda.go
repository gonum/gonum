package stat

import (
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/gonum/mat"
)

type LD struct {
	n, p  int         //n = row, p = col
	ct    []float64   //Constant term of discriminant function of each class
	mu    [][]float64 //Mean vectors of each class
	svd   *mat.SVD
	ok    bool
	eigen mat.Eigen //Eigen values of common variance matrix
}

/**
 * @param x is the training samples
 * @param y is the training labels in [0,k)
 * where k is the number of classes
 */
func (ld *LD) LinearDiscriminant(x mat.Matrix, y []int) (ok bool) {
	ld.n, ld.p = x.Dims()
	fmt.Printf("This is the matrix: %v \n", x)
	fmt.Printf("This is the array: %v \n", y)
	fmt.Printf("x dims: %v, %v \n", ld.n, ld.p)
	if y != nil && len(y) != ld.n {
		panic("The sizes of X and Y don't match")
	}
	var labels []int
	var found bool
	//Find unique labels
	for i := 0; i < len(y); i++ {
		found = false
		for j := 0; j < len(labels); j++ {
			if y[i] == labels[j] {
				found = true
				break
			}
		}
		if !found {
			labels = append(labels, y[i])
		}

	}
	//Create a new array with labels and go through the array of y values and if
	//it doesnt exist then add it to the new array
	sort.Ints(labels)
	fmt.Printf("Sorted list of labels: %v \n", labels)

	if labels[0] != 0 {
		panic("Label does not start from zero")
	}
	for i := 0; i < len(labels); i++ {
		if labels[i] < 0 {
			panic("Negative class label")
		}
		if i > 0 && labels[i]-labels[i-1] > 1 {
			panic("Missing class")
		}
	}
	//Tol is a tolerence to decide if a covariance matrix is singular
	//Tol will reject variables whose variance is less than tol
	var tol float64 = 1E-4
	//k is the number of classes
	var k int = len(labels)
	fmt.Printf("this is k and ld.n: %v, %v \n", k, ld.n)
	if k < 2 {
		panic("Only one class.")
	}
	if tol < 0.0 {
		panic("Invalid tol")
	}
	if ld.n <= k {
		panic("Sample size is too small")
	}

	//Number of instances in each class
	ni := make([]int, k)

	//Common mean vector
	var colmean []float64
	for i := 0; i < ld.p; i++ {
		var col []float64 = mat.Col(nil, i, x)
		var sum float64 = 0
		for _, value := range col {
			sum += value
		}
		colmean = append(colmean, sum/float64(ld.n))
	}
	fmt.Printf("this is the array of means %v \n", colmean)

	//C is a matrix of zeros with dimensions: ld.p x ld.p
	C := mat.NewDense(ld.p, ld.p, make([]float64, ld.p*ld.p, ld.p*ld.p))
	fmt.Printf("this is the zero matrix: %v \n", C)

	//Class mean vectors
	//mu is a matrix with dimensions: k x ld.p
	mu := mat.NewDense(k, ld.p, make([]float64, k*ld.p, k*ld.p))
	for i := 0; i < ld.n; i++ {
		ni[y[i]] = ni[y[i]] + 1
		for j := 0; j < ld.p; j++ {
			mu.Set(y[i], j, ((mu.At(y[i], j)) + (x.At(i, j))))
		}
	}
	for i := 0; i < k; i++ {
		for j := 0; j < ld.p; j++ {
			mu.Set(i, j, ((mu.At(i, j)) / (float64)(ni[i])))
		}
	}

	//priori is the priori probability of each class
	priori := make([]float64, k)
	for i := 0; i < k; i++ {
		priori[i] = (float64)(ni[i] / ld.n)
	}

	//ct is the constant term of discriminant function of each class
	ct := make([]float64, k)
	for i := 0; i < k; i++ {
		ct[i] = math.Log(priori[i])
	}

	for i := 0; i < ld.n; i++ {
		for j := 0; j < ld.p; j++ {
			for l := 0; l <= j; l++ {
				C.Set(j, l, (C.At(j, l) + ((x.At(i, j) - colmean[j]) * (x.At(i, l) - colmean[l]))))
			}
		}
	}

	tol = tol * tol

	for j := 0; j < ld.p; j++ {
		for l := 0; l <= j; l++ {
			C.Set(j, l, ((C.At(j, l)) / (float64)(ld.n-k)))
			C.Set(l, j, C.At(j, l))
		}
		if C.At(j, j) < tol {
			panic("Covarience matrix (variable %d) is close to singular")
		}
	}

	fmt.Printf("this is the code varience %v \n", C)

	//Factorize returns whether the decomposition succeeded
	//If the decomposition failed, methods that require a successful factorization will panic
	ld.eigen.Factorize(C, false, true)
	fmt.Printf("this is the eigen value %v \n", ld.eigen)
	return true
}

func (ld *LD) Transform(x mat.Matrix) *mat.Dense {
	_, p := ld.eigen.Vectors().Dims()
	result := mat.NewDense(ld.n, p, make([]float64, ld.n*p, ld.n*p))
	result.Mul(x, ld.eigen.Vectors())
	return result
}
