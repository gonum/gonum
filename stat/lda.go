package stat

import (
  "sort"
  "fmt"
  "math"


	"gonum.org/v1/gonum/mat"
)

/**
    * The dimensionality of data, the number of classes
    * Constant term of discriminant function of each class
    * Mean vectors of each class
    * Eigen vectors of common covariance matrix, which transforms observations
        * to discriminant functions, normalized so that common covariance
        * matrix is spherical
    * Eigen values of common variance matrix
    */

type LD struct {
	n, p    int  //n= row p = col
	ct []float64
	mu [][]float64
  svd     *mat.SVD
  ok      bool
  eigen mat.Eigen
}

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
// find unique labels
  for i:=0 ; i<len(y); i++{
    found = false
    for j:=0 ; j<len(labels); j++{
      if y[i] == labels[j]{
        found = true
        break
      }
      }
      if !found {
        labels = append(labels, y[i])
      }

  }
  //create a new array with labels and go through the array of y values and if
  //it doesnt exist then add it to the new array
  sort.Ints(labels)
  fmt.Printf("Sorted list of labels: %v \n", labels)

if labels[0] != 0{
  panic("Label does not start from zero")
}
  for i := 0; i < len(labels); i++{
    if labels[i] < 0{
      panic("Negative class label") //: " + labels[i]
    }
    if i > 0 && labels[i] - labels[i-1] > 1{
      panic("Missing class") //: " + labels[i]+1
    }
  }

  var tol float64 = 1E-4
  var k int = len(labels)
  fmt.Printf("this is k and ld.n: %v, %v \n", k, ld.n)
  if k < 2 {
    panic("Only one class.")
  }
  if tol < 0.0 {
    panic("Invalid tol") //: " + tol
  }
  if ld.n <= k {
    panic("Sample size is too small")
  }

//var z int = ld.n
//var n int = ld.p
 ni := make([]int, k) //number of instances in each class

//common mean vector
var colmean []float64
for i:=0 ; i<ld.p ; i++{
 var col []float64 = mat.Col(nil, i, x)
 var sum float64 = 0
 for _, value:=range col{
   sum += value
 }
 colmean = append(colmean, sum/float64 (ld.n))
}
fmt.Printf("this is the array of means %v \n", colmean)


C := mat.NewDense(ld.p,ld.p,make([]float64, ld.p*ld.p,ld.p*ld.p ))
fmt.Printf("this is the zero matrix: %v \n",  C)


mu := mat.NewDense(k,ld.p,make([]float64, k*ld.p,k*ld.p )) //class mean vectors
//calculating class mean vectors
for i:= 0; i<ld.n; i++{
  //var c int = y[i]
  ni[y[i]]= ni[y[i]] + 1
  for j:= 0; j<ld.p; j++{
    mu.Set(y[i],j,((mu.At(y[i],j)) + (x.At(i,j))))
  }
}
for i:= 0; i<k; i++{
  for j:= 0; j<ld.p; j++{
    mu.Set(i,j,((mu.At(i,j)) / (float64)(ni[i])))
  }
}



  priori := make([]float64, k)
  for i:=0; i<k; i++{
    priori[i] = (float64)(ni[i]/ld.n)
  }


//TODO save priori/ct to struct
ct := make([]float64, k)
for i := 0; i<k; i++{
  ct[i] = math.Log(priori[i])
}

for i:=0; i<ld.n; i++{
  for j:= 0; j<ld.p; j++{
    for l:= 0; l<=j; l++{
      C.Set(j,l, (C.At(j,l) + ((x.At(i,j) - colmean[j]) * (x.At(i,l) - colmean[l]))))
    }
  }
}

tol = tol*tol

for j:=0; j<ld.p; j++{
  for l:= 0; l<=j; l++{
    C.Set(j,l, ((C.At(j,l))/ (float64)(ld.n-k)))
    C.Set(l,j,C.At(j,l))
  }
  if C.At(j,j) < tol{
    panic("Covarience matrix (variable %d) is close to singular")
  }
}

fmt.Printf("this is the code varience %v \n", C)

// Factorize returns whether the decomposition succeeded. If the decomposition
// failed, methods that require a successful factorization will panic.
ld.eigen.Factorize(C, false, true)

fmt.Printf("this is the eigen value %v \n", ld.eigen)


return true
}

func (ld *LD) Transform (x mat.Matrix) (*mat.Dense){
  _, p := ld.eigen.Vectors().Dims()
  result := mat.NewDense(ld.n,p, make([]float64, ld.n*p,ld.n*p ))
  result.Mul(x, ld.eigen.Vectors())
  return result


}
