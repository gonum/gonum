// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cluster

import (
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
)

// K-means++ ensures that a dissimilar point to the previously selected centroid is selected as the next
func initKmeansPlusPlus[T any](k int, data []T, r *rand.Rand, dist func(T, T) float64) []T {
	centerIndices := []int{}
	if r != nil {
		centerIndices = append(centerIndices, r.IntN(len(data)))
	} else {
		centerIndices = append(centerIndices, rand.IntN(len(data)))
	}

	for range k - 1 {
		dataClosestCenterDistSquared := make([]float64, len(data))
		dataTotDistSquared := 0.0
		for i := range data {
			if slices.Contains(centerIndices, i) {
				continue
			}
			minDist := math.MaxFloat64
			for j := range centerIndices {
				if d := dist(data[j], data[i]); d < minDist {
					minDist = d
				}
			}
			dataClosestCenterDistSquared[i] = minDist * minDist
			dataTotDistSquared += minDist * minDist
		}
		var lim float64
		if r != nil {
			lim = r.Float64() * dataTotDistSquared
		} else {
			lim = rand.Float64() * dataTotDistSquared
		}
		var center, acc = 0, 0.0
		for i := range dataClosestCenterDistSquared {
			if acc >= lim {
				center = i
				break
			}
			acc += dataClosestCenterDistSquared[i]
		}
		centerIndices = append(centerIndices, center)
	}
	centersData := make([]T, k)
	for i, c := range centerIndices {
		centersData[i] = data[c]
	}
	return centersData
}

// Init centers randomly
func initRandomPoints[T any](k int, data []T, r *rand.Rand) []T {
	centroids := make([]T, k)
	for i := range k {
		if r != nil {
			centroids[i] = data[r.IntN(len(data))]
		} else {

			centroids[i] = data[rand.IntN(len(data))]
		}
	}
	return centroids
}

type kMeansInit int

const (
	InitRandomPoints kMeansInit = iota
	InitPlusPlus
)

// Generic implementation of KMeans. Use KMeans2D or KMeans3D for float64 2D or 3D float64 data
func Kmeans[T any](k int, data []T, i kMeansInit, seed *uint64, epsilon float64, iter int, centroid func([]T) T, dist func(T, T) float64) ([]T, [][]T, error) {
	if len(data) < k {
		return []T{}, [][]T{}, fmt.Errorf("%d clusters is less than %d data points", k, len(data))
	}
	centroids := make([]T, k)
	var r *rand.Rand
	if seed != nil {
		r = rand.New(rand.NewPCG(*seed, *seed))
	}
	switch i {
	case InitPlusPlus:
		centroids = initKmeansPlusPlus(k, data, r, dist)
	default:
		centroids = initRandomPoints(k, data, r)
	}

	calcClusters := func(centroids []T, data []T) [][]T {
		clusters := make([][]T, k)
		for i := range data {
			d := data[i]
			closestIndex, minDistance := 0, dist(d, centroids[0])
			for j := range k {
				dist := dist(d, centroids[j])
				if dist < minDistance {
					closestIndex, minDistance = j, dist
				}
			}
			clusters[closestIndex] = append(clusters[closestIndex], d)
		}
		return clusters
	}
	clusters := make([][]T, k)
	if iter == 0 {
		clusters = calcClusters(centroids, data)
		return centroids, clusters, nil
	}
	for range iter {
		clusters = calcClusters(centroids, data)
		allEqual := true
		newCentroids := make([]T, k)
		for i := range k {
			var newCentroid T
			if len(clusters[i]) > 0 {
				newCentroid = centroid(clusters[i])
			} else { // empty cluster, set centroid to a random point
				if r != nil {
					newCentroid = data[r.IntN(len(data))]
				} else {
					newCentroid = data[rand.IntN(len(data))]
				}
			}
			allEqual = allEqual && dist(newCentroid, centroids[i]) <= epsilon
			newCentroids[i] = newCentroid
		}
		if allEqual {
			break
		}
		centroids = newCentroids
	}
	return centroids, clusters, nil
}

type Config struct {
	Init kMeansInit
	Seed *uint64
	Iter *int
	Eps  float64
}

// KMeans for [x,y] float64 data points. Default config run up to 5 iterations.
// Tolerance for centroid equality is default 0.01, with option epsilon in config
// Centroid initialization defaults to rand.Shuffle, with seed options in config
// Setting 0 iterations returns the initialized centroids and their clusters
func Kmeans2D(k int, data [][2]float64, config Config) ([][2]float64, [][][2]float64, error) {
	dist := func(p1, p2 [2]float64) float64 {
		return math.Hypot(p1[0]-p2[0], p1[1]-p2[1])
	}
	centroid := func(cluster [][2]float64) [2]float64 {
		x, y := 0.0, 0.0
		for _, p := range cluster {
			x += p[0]
			y += p[1]
		}
		l := float64(len(cluster))
		return [2]float64{x / l, y / l}
	}
	if config.Eps <= 0 {
		config.Eps = 0.01
	}
	if config.Iter == nil {
		iter := 5
		config.Iter = &iter
	}
	return Kmeans(k, data, config.Init, config.Seed, config.Eps, *config.Iter, centroid, dist)
}

// KMeans for [x,y,z] float64 data points. Default config run up to 5 iterations.
// Tolerance for centroid equality is default 0.01, with option epsilon in config
// Centroid initialization defaults to rand.Shuffle, with seed options in config
// Setting 0 iterations returns the initialized centroids and their clusters
func Kmeans3D(k int, data [][3]float64, config Config) ([][3]float64, [][][3]float64, error) {
	dist := func(p1, p2 [3]float64) float64 {
		return math.Sqrt(math.Pow(p1[0]-p2[0], 2) + math.Pow(p1[1]-p2[1], 2) + math.Pow(p1[2]-p2[2], 2))
	}
	centroid := func(cluster [][3]float64) [3]float64 {
		x, y, z := 0.0, 0.0, 0.0
		for _, p := range cluster {
			x += p[0]
			y += p[1]
			y += p[2]
		}
		l := float64(len(cluster))
		return [3]float64{x / l, y / l, z / l}
	}
	if config.Eps <= 0 {
		config.Eps = 0.01
	}
	if config.Iter == nil {
		*config.Iter = 5
	}
	return Kmeans(k, data, config.Init, config.Seed, config.Eps, *config.Iter, centroid, dist)
}
