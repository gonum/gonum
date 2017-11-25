// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"math"
	"sort"
	"sync"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distmv"
)

// TODO(btracey): If we ever implement the traditional CMA-ES algorithm, provide
// the base explanation there, and modify this description to just
// describe the differences.

// CmaEsChol implements the covariance matrix adaptation evolution strategy (CMA-ES)
// based on the Cholesky decomposition. The full algorithm is described in
//  Krause, Oswin, Dídac Rodríguez Arbonès, and Christian Igel. "CMA-ES with
//  optimal covariance update and storage complexity." Advances in Neural
//  Information Processing Systems. 2016.
//  https://papers.nips.cc/paper/6457-cma-es-with-optimal-covariance-update-and-storage-complexity.pdf
// CMA-ES is a global optimization method that progressively adapts a population
// of samples. CMA-ES combines techniques from local optimization with global
// optimization. Specifically, the CMA-ES algorithm uses an initial multivariate
// normal distribution to generate a population of input locations. The input locations
// with the lowest function values are used to update the parameters of the normal
// distribution, a new set of input locations are generated, and this procedure
// is iterated until convergence.
//
// As the normal distribution is progressively updated according to the best samples,
// it can be that the mean of the distribution is updated in a gradient-descent
// like fashion, followed by a shrinking covariance.
// It is recommended that the algorithm be run multiple times (with different
// InitMean) to have a better chance of finding the global minimum.
//
// The CMA-ES-Chol algorithm differs from the standard CMA-ES algorithm in that
// it directly updates the Cholesky decomposition of the normal distribution.
// This changes the runtime from O(dimension^3) to O(dimension^2*population)
// The evolution of the multi-variate normal will be similar to the baseline
// CMA-ES algorithm, but the covariance update equation is not identical.
//
// For more information about the CMA-ES algorithm, see
//  https://en.wikipedia.org/wiki/CMA-ES
//  https://arxiv.org/pdf/1604.00772.pdf
type CmaEsChol struct {
	// InitStepSize sets the initial size of the covariance matrix adaptation.
	// If InitStepSize is 0, a default value of 0.5 is used. InitStepSize cannot
	// be negative, or CmaEsChol will panic.
	InitStepSize float64
	// Population sets the population size for the algorithm. If Population is
	// 0, a default value of 4 + math.Floor(3*math.Log(float64(dim))) is used.
	// Population cannot be negative or CmaEsChol will panic.
	Population int
	// InitMean is the initial mean of the multivariate normal for sampling
	// input locations. If InitMean is nil, the zero vector is used. If InitMean
	// is not nil, it must have length equal to the problem dimension.
	InitMean []float64
	// InitCholesky specifies the Cholesky decomposition of the covariance
	// matrix for the initial sampling distribution. If InitCholesky is nil,
	// a default value of I is used. If it is non-nil, then it must have
	// InitCholesky.Size() be equal to the problem dimension.
	InitCholesky *mat.Cholesky
	// StopLogDet sets the threshold for stopping the optimization if the
	// distribution becomes too peaked. The log determinant is a measure of the
	// (log) "volume" of the normal distribution, and when it is too small
	// the samples are almost the same. If the log determinant of the covariance
	// matrix becomes less than StopLogDet, the optimization run is concluded.
	// If StopLogDet is 0, a default value of dim*log(1e-16) is used.
	// If StopLogDet is NaN, the stopping criterion is not used, though
	// this can cause numeric instabilities in the algorithm.
	StopLogDet float64
	// ForgetBest, when true, does not track the best overall function value found,
	// instead returning the new best sample in each iteration. If ForgetBest
	// is false, then the minimum value returned will be the lowest across all
	// iterations, regardless of when that sample was generated.
	ForgetBest bool
	// Src allows a random number generator to be supplied for generating samples.
	// If Src is nil the generator in golang.org/x/math/rand is used.
	Src *rand.Rand

	// Fixed algorithm parameters.
	dim                 int
	pop                 int
	weights             []float64
	muEff               float64
	cc, cs, c1, cmu, ds float64
	eChi                float64

	// Function data.
	xs *mat.Dense
	fs []float64

	// Adaptive algorithm parameters.
	invSigma float64 // inverse of the sigma parameter
	pc, ps   []float64
	mean     []float64
	chol     mat.Cholesky

	// Parallel fields.
	mux      sync.Mutex     // protect access to evals.
	wg       sync.WaitGroup // wait for simulations to finish before iterating.
	taskIdxs []int          // Stores which simulation the task ran.
	evals    []int          // remaining evaluations in this iteration.

	// Overall best.
	bestX []float64
	bestF float64
}

var (
	_ Statuser     = (*CmaEsChol)(nil)
	_ GlobalMethod = (*CmaEsChol)(nil)
)

func (cma *CmaEsChol) Needs() struct{ Gradient, Hessian bool } {
	return struct{ Gradient, Hessian bool }{false, false}
}

func (cma *CmaEsChol) Done() {}

// Status returns the status of the method.
func (cma *CmaEsChol) Status() (Status, error) {
	sd := cma.StopLogDet
	switch {
	case math.IsNaN(sd):
		return NotTerminated, nil
	case sd == 0:
		sd = float64(cma.dim) * -36.8413614879 // ln(1e-16)
	}
	if cma.chol.LogDet() < sd {
		return MethodConverge, nil
	}
	return NotTerminated, nil
}

func (cma *CmaEsChol) InitGlobal(dim, tasks int) int {
	if dim <= 0 {
		panic(nonpositiveDimension)
	}
	if tasks < 0 {
		panic(negativeTasks)
	}
	// Initialize the parameters
	if cma.InitMean != nil && len(cma.InitMean) != dim {
		panic("cma-es-chol: initial mean must be nil or have length equal to dimension")
	}

	// Set fixed algorithm parameters.
	// Parameter values are from https://arxiv.org/pdf/1604.00772.pdf .
	cma.dim = dim
	cma.pop = cma.Population
	n := float64(dim)
	if cma.pop == 0 {
		cma.pop = 4 + int(3*math.Log(n)) // Note the implicit floor.
	} else if cma.pop < 0 {
		panic("cma-es-chol: negative population size")
	}
	mu := cma.pop / 2
	cma.weights = resize(cma.weights, mu)
	for i := range cma.weights {
		v := math.Log(float64(mu)+0.5) - math.Log(float64(i)+1)
		cma.weights[i] = v
	}
	floats.Scale(1/floats.Sum(cma.weights), cma.weights)
	cma.muEff = 0
	for _, v := range cma.weights {
		cma.muEff += v * v
	}
	cma.muEff = 1 / cma.muEff

	cma.cc = (4 + cma.muEff/n) / (n + 4 + 2*cma.muEff/n)
	cma.cs = (cma.muEff + 2) / (n + cma.muEff + 5)
	cma.c1 = 2 / ((n+1.3)*(n+1.3) + cma.muEff)
	cma.cmu = math.Min(1-cma.c1, 2*(cma.muEff-2+1/cma.muEff)/((n+2)*(n+2)+cma.muEff))
	cma.ds = 1 + 2*math.Max(0, math.Sqrt((cma.muEff-1)/(n+1))-1) + cma.cs
	// E[chi] is taken from https://en.wikipedia.org/wiki/CMA-ES (there
	// listed as E[||N(0,1)||]).
	cma.eChi = math.Sqrt(n) * (1 - 1.0/(4*n) + 1/(21*n*n))

	// Allocate memory for function data.
	cma.xs = mat.NewDense(cma.pop, dim, nil)
	cma.fs = resize(cma.fs, cma.pop)

	// Allocate and initialize adaptive parameters.
	cma.invSigma = 1 / cma.InitStepSize
	if cma.InitStepSize == 0 {
		cma.invSigma = 10.0 / 3
	} else if cma.InitStepSize < 0 {
		panic("cma-es-chol: negative initial step size")
	}
	cma.pc = resize(cma.pc, dim)
	for i := range cma.pc {
		cma.pc[i] = 0
	}
	cma.ps = resize(cma.ps, dim)
	for i := range cma.ps {
		cma.ps[i] = 0
	}
	cma.mean = resize(cma.mean, dim)
	if cma.InitMean != nil {
		copy(cma.mean, cma.InitMean)
	}
	if cma.InitCholesky != nil {
		if cma.InitCholesky.Size() != dim {
			panic("cma-es-chol: incorrect InitCholesky size")
		}
		cma.chol.Clone(cma.InitCholesky)
	} else {
		// Set the initial Cholesky to I.
		b := mat.NewDiagonal(dim, nil)
		for i := 0; i < dim; i++ {
			b.SetSymBand(i, i, 1)
		}
		var chol mat.Cholesky
		ok := chol.Factorize(b)
		if !ok {
			panic("cma-es-chol: bad cholesky. shouldn't happen")
		}
		cma.chol = chol
	}

	cma.evals = make([]int, cma.pop)
	for i := range cma.evals {
		cma.evals[i] = i
	}

	cma.bestX = resize(cma.bestX, dim)
	cma.bestF = math.Inf(1)

	t := min(tasks, cma.pop)
	cma.taskIdxs = make([]int, t)
	for i := 0; i < t; i++ {
		cma.taskIdxs[i] = -1
	}
	// Get a new mutex and waitgroup so that if the structure is reused there
	// aren't residual interactions with the previous optimization.
	cma.mux = sync.Mutex{}
	cma.wg = sync.WaitGroup{}
	return t
}

func (cma *CmaEsChol) IterateGlobal(task int, loc *Location) (Operation, error) {
	// Check the status of the incoming task. If it is a number, it means
	// that task contains a valid location.
	idx := cma.taskIdxs[task]
	if idx != -1 {
		cma.fs[idx] = loc.F
		cma.wg.Done()
	}

	// Get the next task and send it to be run if there is a next task to be run.
	// If all of the tasks have been run, perform an update step. Note that the
	// use of this mutex means that only one task can proceed, all of the
	// other tasks should get stuck and then get a new location.
	cma.mux.Lock()
	if len(cma.evals) != 0 {
		// There are still tasks to evaluate. Grab one and remove it from the list.
		newIdx := cma.evals[len(cma.evals)-1]
		cma.evals = cma.evals[:len(cma.evals)-1]
		cma.wg.Add(1)
		cma.mux.Unlock()

		// Sample x and send it to be evaluated.
		distmv.NormalRand(cma.xs.RawRowView(newIdx), cma.mean, &cma.chol, cma.Src)
		copy(loc.X, cma.xs.RawRowView(newIdx))
		cma.taskIdxs[task] = newIdx
		return FuncEvaluation, nil
	}
	// There are no more tasks to evaluate. This means the iteration is over.
	// Find the best current f, update the parameters, and re-establish
	// the evaluations to run.

	// Wait for all of the outstanding tasks to finish, so the full set of functions
	// has been evaluated.
	cma.wg.Wait()

	// Find the best f out of all the tasks.
	best := floats.MinIdx(cma.fs)
	bestF := cma.fs[best]
	bestX := cma.xs.RawRowView(best)
	if cma.ForgetBest {
		loc.F = bestF
		copy(loc.X, bestX)
	} else {
		if bestF < cma.bestF {
			cma.bestF = bestF
			copy(cma.bestX, bestX)
		}
		loc.F = cma.bestF
		copy(loc.X, cma.bestX)
	}

	cma.taskIdxs[task] = -1

	// Update the parameters of the distribution
	err := cma.update()

	// Reset the tasks
	cma.evals = cma.evals[:cma.pop]

	cma.mux.Unlock()
	return MajorIteration, err
}

// update computes the new parameters (mean, cholesky, etc.)
func (cma *CmaEsChol) update() error {
	// Sort the function values to find the elite samples.
	ftmp := make([]float64, cma.pop)
	copy(ftmp, cma.fs)
	indexes := make([]int, cma.pop)
	for i := range indexes {
		indexes[i] = i
	}
	sort.Sort(bestSorter{F: ftmp, Idx: indexes})

	meanOld := make([]float64, len(cma.mean))
	copy(meanOld, cma.mean)

	// m_{t+1} = \sum_{i=1}^mu w_i x_i
	for i := range cma.mean {
		cma.mean[i] = 0
	}
	for i, w := range cma.weights {
		idx := indexes[i] // index of teh 1337 sample.
		floats.AddScaled(cma.mean, w, cma.xs.RawRowView(idx))
	}
	meanDiff := make([]float64, len(cma.mean))
	floats.SubTo(meanDiff, cma.mean, meanOld)

	// p_{c,t+1} = (1-c_c) p_{c,t} + \sqrt(c_c*(2-c_c)*mueff) (m_{t+1}-m_t)/sigma_t
	floats.Scale(1-cma.cc, cma.pc)
	scaleC := math.Sqrt(cma.cc*(2-cma.cc)*cma.muEff) * cma.invSigma
	floats.AddScaled(cma.pc, scaleC, meanDiff)

	// p_{sigma, t+1} = (1-c_sigma) p_{sigma,t} + \sqrt(c_s*(2-c_s)*mueff) A_t^-1 (m_{t+1}-m_t)/sigma_t
	floats.Scale(1-cma.cs, cma.ps)
	// First compute A_t^-1 (m_{t+1}-m_t), then add the scaled vector.
	tmp := make([]float64, cma.dim)
	tmpVec := mat.NewVecDense(cma.dim, tmp)
	diffVec := mat.NewVecDense(cma.dim, meanDiff)
	err := tmpVec.SolveVec(cma.chol.RawU().T(), diffVec)
	if err != nil {
		return err
	}
	scaleS := math.Sqrt(cma.cs*(2-cma.cs)*cma.muEff) * cma.invSigma
	floats.AddScaled(cma.ps, scaleS, tmp)

	// Compute the update to A.
	scaleChol := 1 - cma.c1 - cma.cmu
	if scaleChol == 0 {
		scaleChol = math.SmallestNonzeroFloat64 // enough to kill the old data, but still non-zero.
	}
	cma.chol.Scale(scaleChol, &cma.chol)
	cma.chol.SymRankOne(&cma.chol, cma.c1, mat.NewVecDense(cma.dim, cma.pc))
	for i, w := range cma.weights {
		idx := indexes[i]
		floats.SubTo(tmp, cma.xs.RawRowView(idx), meanOld)
		cma.chol.SymRankOne(&cma.chol, cma.cmu*w*cma.invSigma, tmpVec)
	}

	// sigma_{t+1} = sigma_t exp(c_sigma/d_sigma * norm(p_{sigma,t+1}/ E[chi] -1)
	normPs := floats.Norm(cma.ps, 2)
	cma.invSigma /= math.Exp(cma.cs / cma.ds * (normPs/cma.eChi - 1))
	return nil
}

type bestSorter struct {
	F   []float64
	Idx []int
}

func (b bestSorter) Len() int {
	return len(b.F)
}
func (b bestSorter) Less(i, j int) bool {
	return b.F[i] < b.F[j]
}
func (b bestSorter) Swap(i, j int) {
	b.F[i], b.F[j] = b.F[j], b.F[i]
	b.Idx[i], b.Idx[j] = b.Idx[j], b.Idx[i]
}
