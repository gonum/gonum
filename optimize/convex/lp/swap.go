package lp

import (
    "errors"
    "math"

    "gonum.org/v1/gonum/mat"
    "gonum.org/v1/gonum/floats"
    "gonum.org/v1/gonum/lapack"
)

var (
    ErrSwapSingular = errors.New("swap: swap matrix is singular")
)

// Swap is a type that represents a chain of column swaps performed on
// a matrix A implicitly as the product of rank-one updates to the identity matrix
//  E_i = I + (y - e_k) * e_k^T
// such that
//  A * E_0 * ... * E_[i - 1] * y = v_i,
// where v_i is ith column exchange to be inserted at position k.
type Swap struct {
    Dim int
    idx []int
    swapCache []float64
    cond float64
}

// Append adds the vector y to the receiver and updates the bound on the
// condition number.
func (s *Swap) Append(y []float64, k int) {
    if len(y) != s.Dim {
        panic(mat.ErrShape)
    }

    // update bound on condition number
    if c := cond(y, k, mat.CondNorm); s.Len() == 0 {
        s.cond = c
    } else {
        s.cond *= c
    }

    // append vector
    s.idx = append(s.idx, k)
    s.swapCache = append(s.swapCache, y...)
}

// Reset returns the length of the receiver to zero without decreasing the
// capacity.
func (s *Swap) Reset() {
    s.idx = s.idx[:0]
    s.swapCache = s.swapCache[:0]
}

// SolveVec solves a system of linear equations defined by the receiver.
// It computes x such that
//  E_0 * ... E_i * x = b if trans == false
//  E_i^T * ... * E_0^T * x = b if trans == true
// The systems are solved sequentually using the Sherman–Morrison formula,
// such that the ouput of one system becomes the input for the next.
//
// If A is singular exactly singular, an error is returned.
func (s *Swap) SolveVec(v *mat.VecDense, trans bool, b *mat.VecDense) error {
    n := s.Dim
    if b.Len() != n {
        panic(mat.ErrShape)
    }
    if v != nil && v.Len() != n {
        panic(mat.ErrShape)
    }

    // TODO: Add overlap check
    if v == nil {
        v = b
    } else if v != b {
        v.CopyVec(b)
    }

    m := s.Len()
    if !trans {
        var vkyk float64
        for i := 0; i < m; i++ {
            k := s.idx[i]
            y := s.swapCache[i * n : (i + 1) * n]
            yVec := mat.NewVecDense(n, y)
            if a := y[k]; a != 0 {
                vkyk = v.At(k, 0) / a
            } else {
                return ErrSwapSingular
            }
            v.AddScaledVec(v, -vkyk, yVec)
            v.SetVec(k, vkyk)
        }
    } else {
        for i := m - 1; i >= 0; i-- {
            k := s.idx[i]
            y := s.swapCache[i * n : (i + 1) * n]
            yVec := mat.NewVecDense(n, y)
            if a := y[k]; a != 0 {
                vk := v.At(k, 0)
                v.SetVec(k, vk - (mat.Dot(yVec, v) - vk) / a)
            } else {
                return ErrSwapSingular
            }
        }
    }
    return nil
}

// Len returns the number of vectors stored in the receiver.
func (s *Swap) Len() int {
    return len(s.swapCache) / s.Dim
}

// Cap returns the maximum number of vectors that can currently be stored in the receiver.
func (s *Swap) Cap() int {
    return cap(s.swapCache) / s.Dim
}

// Cond returns (an upper bound on) the condition number of the swap structure.
// Cond will panic if the receiver does not currently contain vectors.
func (s *Swap) Cond() float64 {
    if s.swapCache == nil || s.Len() == 0 {
        panic("swap: swap matrix is empty")
    }
    return s.cond
}

// condition number helper functions

func exclusiveAbsMax(y []float64, k int) float64{
    n := len(y)
    L := math.Inf(1)

    if k < 0 || k >= n {
        panic("swap: index out of bounds")
    } else if k > 0 && k < n {
        return math.Max(floats.Norm(y[:k], L), floats.Norm(y[k + 1:], L))
    } else if k == 0 {
        return floats.Norm(y[1:], L)
    }
    return floats.Norm(y[:n - 1], L)
}

// cond calculates the condition number of the matrix E = I + (y - e_k) * e_k^T
// for a given norm.
func cond(y []float64, k int, norm lapack.MatrixNorm) float64 {
    yk := math.Abs(y[k])
    if yk == 0 {
        return math.Inf(1)
    }
    beta := 1 / yk
    var normA, normAInv float64

    switch norm {
    case 'M':
        ymax := exclusiveAbsMax(y, k)
        normA = math.Max(1, math.Max(ymax, yk))
        normAInv = math.Max(1, beta * math.Max(ymax, 1))
    case 'O':
        y1norm := floats.Norm(y, 1)
        normA = math.Max(1, y1norm)
        normAInv = math.Max(1, beta * (y1norm + 1) - 1)
    case 'I':
        ymax := exclusiveAbsMax(y, k)
        normA = math.Max(1 + ymax, yk)
        normAInv = math.Max(1 + beta * ymax, beta)
    case 'F':
        n := float64(len(y))
        ydot := floats.Dot(y, y)
        normA = math.Sqrt(ydot + n - 1)
        normAInv = math.Sqrt(beta * beta * (ydot + 1) - 1 + n - 1)
    }
    return normA * normAInv
}
