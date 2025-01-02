// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmv

import (
	"math"
	"math/rand/v2"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

// Normal is a multivariate normal distribution (also known as the multivariate
// Gaussian distribution). Its pdf in k dimensions is given by
//
//	(2 π)^(-k/2) |Σ|^(-1/2) exp(-1/2 (x-μ)'Σ^-1(x-μ))
//
// where μ is the mean vector and Σ the covariance matrix. Σ must be symmetric
// and positive definite. Use NewNormal to construct.
type Normal struct {
	mu []float64

	sigma mat.SymDense

	chol       mat.Cholesky
	logSqrtDet float64
	dim        int

	// If src is altered, rnd must be updated.
	src rand.Source
	rnd *rand.Rand
}

// NewNormal creates a new Normal with the given mean and covariance matrix.
// NewNormal panics if len(mu) == 0, or if len(mu) != sigma.N. If the covariance
// matrix is not positive-definite, the returned boolean is false.
func NewNormal(mu []float64, sigma mat.Symmetric, src rand.Source) (*Normal, bool) {
	if len(mu) == 0 {
		panic(badZeroDimension)
	}
	dim := sigma.SymmetricDim()
	if dim != len(mu) {
		panic(badSizeMismatch)
	}
	n := &Normal{
		src: src,
		rnd: rand.New(src),
		dim: dim,
		mu:  make([]float64, dim),
	}
	copy(n.mu, mu)
	ok := n.chol.Factorize(sigma)
	if !ok {
		return nil, false
	}
	n.sigma = *mat.NewSymDense(dim, nil)
	n.sigma.CopySym(sigma)
	n.logSqrtDet = 0.5 * n.chol.LogDet()
	return n, true
}

// NewNormalChol creates a new Normal distribution with the given mean and
// covariance matrix represented by its Cholesky decomposition. NewNormalChol
// panics if len(mu) is not equal to chol.Size().
func NewNormalChol(mu []float64, chol *mat.Cholesky, src rand.Source) *Normal {
	dim := len(mu)
	if dim != chol.SymmetricDim() {
		panic(badSizeMismatch)
	}
	n := &Normal{
		src: src,
		rnd: rand.New(src),
		dim: dim,
		mu:  make([]float64, dim),
	}
	n.chol.Clone(chol)
	copy(n.mu, mu)
	n.logSqrtDet = 0.5 * n.chol.LogDet()
	return n
}

// NewNormalPrecision creates a new Normal distribution with the given mean and
// precision matrix (inverse of the covariance matrix). NewNormalPrecision
// panics if len(mu) is not equal to prec.SymmetricDim(). If the precision matrix
// is not positive-definite, NewNormalPrecision returns nil for norm and false
// for ok.
func NewNormalPrecision(mu []float64, prec *mat.SymDense, src rand.Source) (norm *Normal, ok bool) {
	if len(mu) == 0 {
		panic(badZeroDimension)
	}
	dim := prec.SymmetricDim()
	if dim != len(mu) {
		panic(badSizeMismatch)
	}
	// TODO(btracey): Computing a matrix inverse is generally numerically unstable.
	// This only has to compute the inverse of a positive definite matrix, which
	// is much better, but this still loses precision. It is worth considering if
	// instead the precision matrix should be stored explicitly and used instead
	// of the Cholesky decomposition of the covariance matrix where appropriate.
	var chol mat.Cholesky
	ok = chol.Factorize(prec)
	if !ok {
		return nil, false
	}
	var sigma mat.SymDense
	err := chol.InverseTo(&sigma)
	if err != nil {
		return nil, false
	}
	return NewNormal(mu, &sigma, src)
}

// ConditionNormal returns the Normal distribution that is the receiver conditioned
// on the input evidence. The returned multivariate normal has dimension
// n - len(observed), where n is the dimension of the original receiver. The updated
// mean and covariance are
//
//	mu = mu_un + sigma_{ob,un}ᵀ * sigma_{ob,ob}^-1 (v - mu_ob)
//	sigma = sigma_{un,un} - sigma_{ob,un}ᵀ * sigma_{ob,ob}^-1 * sigma_{ob,un}
//
// where mu_un and mu_ob are the original means of the unobserved and observed
// variables respectively, sigma_{un,un} is the unobserved subset of the covariance
// matrix, sigma_{ob,ob} is the observed subset of the covariance matrix, and
// sigma_{un,ob} are the cross terms. The elements of x_2 have been observed with
// values v. The dimension order is preserved during conditioning, so if the value
// of dimension 1 is observed, the returned normal represents dimensions {0, 2, ...}
// of the original Normal distribution.
//
// ConditionNormal returns {nil, false} if there is a failure during the update.
// Mathematically this is impossible, but can occur with finite precision arithmetic.
func (n *Normal) ConditionNormal(observed []int, values []float64, src rand.Source) (*Normal, bool) {
	if len(observed) == 0 {
		panic("normal: no observed value")
	}
	if len(observed) != len(values) {
		panic(badInputLength)
	}
	for _, v := range observed {
		if v < 0 || v >= n.Dim() {
			panic("normal: observed value out of bounds")
		}
	}

	_, mu1, sigma11 := studentsTConditional(observed, values, math.Inf(1), n.mu, &n.sigma)
	if mu1 == nil {
		return nil, false
	}
	return NewNormal(mu1, sigma11, src)
}

// CovarianceMatrix stores the covariance matrix of the distribution in dst.
// Upon return, the value at element {i, j} of the covariance matrix is equal
// to the covariance of the i^th and j^th variables.
//
//	covariance(i, j) = E[(x_i - E[x_i])(x_j - E[x_j])]
//
// If the dst matrix is empty it will be resized to the correct dimensions,
// otherwise dst must match the dimension of the receiver or CovarianceMatrix
// will panic.
func (n *Normal) CovarianceMatrix(dst *mat.SymDense) {
	if dst.IsEmpty() {
		*dst = *(dst.GrowSym(n.dim).(*mat.SymDense))
	} else if dst.SymmetricDim() != n.dim {
		panic("normal: input matrix size mismatch")
	}
	dst.CopySym(&n.sigma)
}

// Dim returns the dimension of the distribution.
func (n *Normal) Dim() int {
	return n.dim
}

// Entropy returns the differential entropy of the distribution.
func (n *Normal) Entropy() float64 {
	return float64(n.dim)/2*(1+logTwoPi) + n.logSqrtDet
}

// LogProb computes the log of the pdf of the point x.
func (n *Normal) LogProb(x []float64) float64 {
	dim := n.dim
	if len(x) != dim {
		panic(badSizeMismatch)
	}
	return normalLogProb(x, n.mu, &n.chol, n.logSqrtDet)
}

// NormalLogProb computes the log probability of the location x for a Normal
// distribution the given mean and Cholesky decomposition of the covariance matrix.
// NormalLogProb panics if len(x) is not equal to len(mu), or if len(mu) != chol.Size().
//
// This function saves time and memory if the Cholesky decomposition is already
// available. Otherwise, the NewNormal function should be used.
func NormalLogProb(x, mu []float64, chol *mat.Cholesky) float64 {
	dim := len(mu)
	if len(x) != dim {
		panic(badSizeMismatch)
	}
	if chol.SymmetricDim() != dim {
		panic(badSizeMismatch)
	}
	logSqrtDet := 0.5 * chol.LogDet()
	return normalLogProb(x, mu, chol, logSqrtDet)
}

// normalLogProb is the same as NormalLogProb, but does not make size checks and
// additionally requires log(|Σ|^-0.5)
func normalLogProb(x, mu []float64, chol *mat.Cholesky, logSqrtDet float64) float64 {
	dim := len(mu)
	c := -0.5*float64(dim)*logTwoPi - logSqrtDet
	dst := stat.Mahalanobis(mat.NewVecDense(dim, x), mat.NewVecDense(dim, mu), chol)
	return c - 0.5*dst*dst
}

// MarginalNormal returns the marginal distribution of the given input variables.
// That is, MarginalNormal returns
//
//	p(x_i) = \int_{x_o} p(x_i | x_o) p(x_o) dx_o
//
// where x_i are the dimensions in the input, and x_o are the remaining dimensions.
// See https://en.wikipedia.org/wiki/Marginal_distribution for more information.
//
// The input src is passed to the call to NewNormal.
func (n *Normal) MarginalNormal(vars []int, src rand.Source) (*Normal, bool) {
	newMean := make([]float64, len(vars))
	for i, v := range vars {
		newMean[i] = n.mu[v]
	}
	var s mat.SymDense
	s.SubsetSym(&n.sigma, vars)
	return NewNormal(newMean, &s, src)
}

// MarginalNormalSingle returns the marginal of the given input variable.
// That is, MarginalNormal returns
//
//	p(x_i) = \int_{x_¬i} p(x_i | x_¬i) p(x_¬i) dx_¬i
//
// where i is the input index.
// See https://en.wikipedia.org/wiki/Marginal_distribution for more information.
//
// The input src is passed to the constructed distuv.Normal.
func (n *Normal) MarginalNormalSingle(i int, src rand.Source) distuv.Normal {
	return distuv.Normal{
		Mu:    n.mu[i],
		Sigma: math.Sqrt(n.sigma.At(i, i)),
		Src:   src,
	}
}

// Mean returns the mean of the probability distribution.
//
// If dst is not nil, the mean will be stored in-place into dst and returned,
// otherwise a new slice will be allocated first. If dst is not nil, it must
// have length equal to the dimension of the distribution.
func (n *Normal) Mean(dst []float64) []float64 {
	dst = reuseAs(dst, n.dim)
	copy(dst, n.mu)
	return dst
}

// Prob computes the value of the probability density function at x.
func (n *Normal) Prob(x []float64) float64 {
	return math.Exp(n.LogProb(x))
}

// Quantile returns the value of the multi-dimensional inverse cumulative
// distribution function at p.
//
// If dst is not nil, the quantile will be stored in-place into dst and
// returned, otherwise a new slice will be allocated first. If dst is not nil,
// it must have length equal to the dimension of the distribution. Quantile will
// also panic if the length of p is not equal to the dimension of the
// distribution.
//
// All of the values of p must be between 0 and 1, inclusive, or Quantile will
// panic.
func (n *Normal) Quantile(dst, p []float64) []float64 {
	if len(p) != n.dim {
		panic(badInputLength)
	}
	dst = reuseAs(dst, n.dim)

	// Transform to a standard normal and then transform to a multivariate Gaussian.
	for i, v := range p {
		dst[i] = distuv.UnitNormal.Quantile(v)
	}
	n.TransformNormal(dst, dst)
	return dst
}

// Rand generates a random sample according to the distribution.
//
// If dst is not nil, the sample will be stored in-place into dst and returned,
// otherwise a new slice will be allocated first. If dst is not nil, it must
// have length equal to the dimension of the distribution.
func (n *Normal) Rand(dst []float64) []float64 {
	return NormalRand(dst, n.mu, &n.chol, n.src)
}

// NormalRand generates a random sample from a multivariate normal distribution
// given by the mean and the Cholesky factorization of the covariance matrix.
//
// If dst is not nil, the sample will be stored in-place into dst and returned,
// otherwise a new slice will be allocated first. If dst is not nil, it must
// have length equal to the dimension of the distribution.
//
// This function saves time and memory if the Cholesky factorization is already
// available. Otherwise, the NewNormal function should be used.
func NormalRand(dst, mean []float64, chol *mat.Cholesky, src rand.Source) []float64 {
	if len(mean) != chol.SymmetricDim() {
		panic(badInputLength)
	}
	dst = reuseAs(dst, len(mean))

	if src == nil {
		for i := range dst {
			dst[i] = rand.NormFloat64()
		}
	} else {
		rnd := rand.New(src)
		for i := range dst {
			dst[i] = rnd.NormFloat64()
		}
	}
	transformNormal(dst, dst, mean, chol)
	return dst
}

// EigenSym is an eigendecomposition of a symmetric matrix.
type EigenSym interface {
	mat.Symmetric
	// RawValues returns all eigenvalues in ascending order. The returned slice
	// must not be modified.
	RawValues() []float64
	// RawQ returns an orthogonal matrix whose columns contain the eigenvectors.
	// The returned matrix must not be modified.
	RawQ() mat.Matrix
}

// PositivePartEigenSym is an EigenSym that sets any negative eigenvalues from
// the given eigendecomposition to zero but otherwise returns the values
// unchanged.
//
// This is useful for filtering eigenvalues of positive semi-definite matrices
// that are almost zero but negative due to rounding errors.
type PositivePartEigenSym struct {
	ed   *mat.EigenSym
	vals []float64
}

var _ EigenSym = (*PositivePartEigenSym)(nil)
var _ EigenSym = (*mat.EigenSym)(nil)

// NewPositivePartEigenSym returns a new PositivePartEigenSym, wrapping the
// given eigendecomposition.
func NewPositivePartEigenSym(ed *mat.EigenSym) *PositivePartEigenSym {
	n := ed.SymmetricDim()
	vals := make([]float64, n)
	for i, lamda := range ed.RawValues() {
		if lamda > 0 {
			vals[i] = lamda
		}
	}
	return &PositivePartEigenSym{
		ed:   ed,
		vals: vals,
	}
}

// SymmetricDim returns the value from the wrapped eigendecomposition.
func (ed *PositivePartEigenSym) SymmetricDim() int { return ed.ed.SymmetricDim() }

// Dims returns the dimensions from the wrapped eigendecomposition.
func (ed *PositivePartEigenSym) Dims() (r, c int) { return ed.ed.Dims() }

// At returns the value from the wrapped eigendecomposition.
func (ed *PositivePartEigenSym) At(i, j int) float64 { return ed.ed.At(i, j) }

// T returns the transpose from the wrapped eigendecomposition.
func (ed *PositivePartEigenSym) T() mat.Matrix { return ed.ed.T() }

// RawQ returns the orthogonal matrix Q from the wrapped eigendecomposition. The
// returned matrix must not be modified.
func (ed *PositivePartEigenSym) RawQ() mat.Matrix { return ed.ed.RawQ() }

// RawValues returns the eigenvalues from the wrapped eigendecomposition in
// ascending order with any negative value replaced by zero. The returned slice
// must not be modified.
func (ed *PositivePartEigenSym) RawValues() []float64 { return ed.vals }

// NormalRandCov generates a random sample from a multivariate normal
// distribution given by the mean and the covariance matrix.
//
// If dst is not nil, the sample will be stored in-place into dst and returned,
// otherwise a new slice will be allocated first. If dst is not nil, it must
// have length equal to the dimension of the distribution.
//
// cov should be *mat.Cholesky, *mat.PivotedCholesky or EigenSym, otherwise
// NormalRandCov will be very inefficient because a pivoted Cholesky
// factorization of cov will be computed for every sample.
//
// If cov is an EigenSym, all eigenvalues returned by RawValues must be
// non-negative, otherwise NormalRandCov will panic.
func NormalRandCov(dst, mean []float64, cov mat.Symmetric, src rand.Source) []float64 {
	n := len(mean)
	if cov.SymmetricDim() != n {
		panic(badInputLength)
	}
	dst = reuseAs(dst, n)
	if src == nil {
		for i := range dst {
			dst[i] = rand.NormFloat64()
		}
	} else {
		rnd := rand.New(src)
		for i := range dst {
			dst[i] = rnd.NormFloat64()
		}
	}

	switch cov := cov.(type) {
	case *mat.Cholesky:
		dstVec := mat.NewVecDense(n, dst)
		dstVec.MulVec(cov.RawU().T(), dstVec)
	case *mat.PivotedCholesky:
		dstVec := mat.NewVecDense(n, dst)
		dstVec.MulVec(cov.RawU().T(), dstVec)
		dstVec.Permute(cov.ColumnPivots(nil), true)
	case EigenSym:
		vals := cov.RawValues()
		if vals[0] < 0 {
			panic("distmv: covariance matrix is not positive semi-definite")
		}
		for i, val := range vals {
			dst[i] *= math.Sqrt(val)
		}
		dstVec := mat.NewVecDense(n, dst)
		dstVec.MulVec(cov.RawQ(), dstVec)
	default:
		var chol mat.PivotedCholesky
		chol.Factorize(cov, -1)
		dstVec := mat.NewVecDense(n, dst)
		dstVec.MulVec(chol.RawU().T(), dstVec)
		dstVec.Permute(chol.ColumnPivots(nil), true)
	}
	floats.Add(dst, mean)

	return dst
}

// ScoreInput returns the gradient of the log-probability with respect to the
// input x. That is, ScoreInput computes
//
//	∇_x log(p(x))
//
// If dst is not nil, the score will be stored in-place into dst and returned,
// otherwise a new slice will be allocated first. If dst is not nil, it must
// have length equal to the dimension of the distribution.
func (n *Normal) ScoreInput(dst, x []float64) []float64 {
	// Normal log probability is
	//  c - 0.5*(x-μ)' Σ^-1 (x-μ).
	// So the derivative is just
	//  -Σ^-1 (x-μ).
	if len(x) != n.Dim() {
		panic(badInputLength)
	}
	dst = reuseAs(dst, n.dim)

	floats.SubTo(dst, x, n.mu)
	dstVec := mat.NewVecDense(len(dst), dst)
	err := n.chol.SolveVecTo(dstVec, dstVec)
	if err != nil {
		panic(err)
	}
	floats.Scale(-1, dst)
	return dst
}

// SetMean changes the mean of the normal distribution. SetMean panics if len(mu)
// does not equal the dimension of the normal distribution.
func (n *Normal) SetMean(mu []float64) {
	if len(mu) != n.Dim() {
		panic(badSizeMismatch)
	}
	copy(n.mu, mu)
}

// TransformNormal transforms x generated from a standard multivariate normal
// into a vector that has been generated under the normal distribution of the
// receiver.
//
// If dst is not nil, the result will be stored in-place into dst and returned,
// otherwise a new slice will be allocated first. If dst is not nil, it must
// have length equal to the dimension of the distribution. TransformNormal will
// also panic if the length of x is not equal to the dimension of the receiver.
func (n *Normal) TransformNormal(dst, x []float64) []float64 {
	if len(x) != n.dim {
		panic(badInputLength)
	}
	dst = reuseAs(dst, n.dim)
	transformNormal(dst, x, n.mu, &n.chol)
	return dst
}

// transformNormal performs the same operation as Normal.TransformNormal except
// no safety checks are performed and all memory must be provided.
func transformNormal(dst, normal, mu []float64, chol *mat.Cholesky) []float64 {
	dim := len(mu)
	dstVec := mat.NewVecDense(dim, dst)
	srcVec := mat.NewVecDense(dim, normal)
	// If dst and normal are the same slice, make them the same Vector otherwise
	// mat complains about being tricky.
	if &normal[0] == &dst[0] {
		srcVec = dstVec
	}
	dstVec.MulVec(chol.RawU().T(), srcVec)
	floats.Add(dst, mu)
	return dst
}
