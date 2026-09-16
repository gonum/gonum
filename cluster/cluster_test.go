// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cluster

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
	"testing"
)

func TestKmeansSeed(t *testing.T) {
	data := [][2]float64{{rand.NormFloat64()*1 + 3, rand.NormFloat64()*1 + 3}, {rand.NormFloat64()*1 + 3, rand.NormFloat64()*1 + 3}}
	var seed uint64 = 123
	config := Config{Seed: &seed}
	cent1, _, err := Kmeans2D(1, data, config)
	if err != nil {
		t.Error(err)
	}
	cent2, _, err := Kmeans2D(1, data, config)
	if err != nil {
		t.Error(err)
	}
	for i := range cent1 {
		if cent1[i] != cent2[i] {
			t.Error("runs with same seed doesn't yield same result")
		}
	}
}

func TestKmeans2DZeroIter(t *testing.T) {
	data := [][2]float64{{rand.NormFloat64()*1 + 3, rand.NormFloat64()*1 + 3}}
	iter := 0
	configs := []Config{
		{Iter: &iter},
		{Iter: &iter, Init: InitPlusPlus},
	}
	for i := range len(configs) {
		cent, clust, err := Kmeans2D(1, data, configs[i])
		if err != nil {
			t.Error(err)
		}
		if len(cent) != 1 || len(clust) != 1 {
			t.Errorf("kmeans invalid dimensions for config %d", i)
		}
	}
}

func TestKmeans2DEmptyData(t *testing.T) {
	data := [][2]float64{}
	iter := 10
	configs := []Config{
		{Iter: &iter},
		{Iter: &iter, Init: InitPlusPlus},
	}
	for i := range len(configs) {
		_, _, err := Kmeans2D(3, data, configs[i])
		if err == nil {
			t.Error("expected error on empty data")
		}
	}
}

func TestKmeans2D(t *testing.T) {
	data := [][2]float64{}
	for range 500 {
		data = append(data,
			[2]float64{rand.NormFloat64()*1 + 3, rand.NormFloat64()*1 + 3},
			[2]float64{rand.NormFloat64()*1 + 3, rand.NormFloat64()*1 - 3},
			[2]float64{rand.NormFloat64()*1 - 3, rand.NormFloat64()*1 + 3},
		)
	}
	expectedCenters := [][2]float64{{3, 3}, {3, -3}, {-3, 3}}
	iter := 5
	var seed uint64 = 123
	cent, clust, err := Kmeans2D(len(expectedCenters), data, Config{Seed: &seed, Iter: &iter})
	if err != nil {
		t.Error(err)
	}
	if len(cent) != len(clust) {
		t.Error("kmeans invalid dimensions")
	}
	eps := 0.15
	for _, c1 := range cent {
		found := false
		for _, c2 := range expectedCenters {
			if math.Abs(c1[0]-c2[0]) <= eps && math.Abs(c1[1]-c2[1]) <= eps {
				found = true
			}
		}
		if !found {
			t.Errorf("center %+v not correct", c1)
		}
	}
}

func TestKmeans2DVisual(t *testing.T) {
	if os.Getenv("VISUAL_TEST") == "1" {
		type point struct {
			X     float64 `json:"x"`
			Y     float64 `json:"y"`
			Color string  `json:"color"`
		}
		html := `<html><body><form>
			<b>k = %d</b><input type="range" min="1" max="40" name="k" value="%d" style="width:200px"><br>
			<b>c = %d</b><input type="range" min="1" max="40" name="c" value="%d" style="width:200px"><br>
			<b>i = %d</b><input type="range" min="1" max="40" name="i" value="%d" style="width:200px"><br>
			<b>init strategy</b><select name="init"><option %s value="random">random points</option><option %s value="plusplus">kmeans++</option></select><br>
			<b>data point random seed</b><input name="dataRandom" %s type="checkbox"><br>
			<b>kmeans init random seed</b><input name="initRandom" %s type="checkbox"><br>
			<input type="submit" value="Run"></form><canvas id="canvas" width="900" height="800"></canvas>
		<script>let h = 800; let w = 900; let zoom = 20;
			let ctx = document.getElementById("canvas").getContext("2d")
			ctx.fillRect(0,h/2,w,1); ctx.fillRect(w/2,0,1,h)
			let d = %s
			for(let i in d.clusters) {
				let [r,g,b] = [Math.random()*200+20, Math.random()*200+20, Math.random()*200+20]
				ctx.setFillColor("rgb("+r+","+g+","+b+")")
				for(let p of d.clusters[i]) ctx.fillRect((zoom*p[0])+w/2,(zoom*-p[1])+h/2,3,3)
			}
			ctx.setFillColor("#333")
			for(let i = 0; i < d.centroidIterations.length; i++) {
				let centroids = d.centroidIterations[i]
				for(let j in centroids) {
					let [x,y] = centroids[j]
					if(i > 0) {
						let [x_, y_] = d.centroidIterations[i-1][j]
						ctx.beginPath(); ctx.lineWidth = 2
						ctx.moveTo((zoom*x_)+w/2, (zoom*-y_)+h/2)
						ctx.lineTo((zoom*x)+w/2, (zoom*-y)+h/2)
						ctx.strokeStyle = "#333"; ctx.stroke()
					}
					let s = i == d.centroidIterations.length-1 ? 10 : 5
					ctx.fillRect((zoom*x)+w/2, (zoom*-y)+h/2, s, s)
				}
			}</script></html></body>`

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			k, err := strconv.Atoi(q.Get("k"))
			if err != nil {
				k = 3
			}
			clusters, err := strconv.Atoi(q.Get("c"))
			if err != nil {
				clusters = 3
			}
			iter, err := strconv.Atoi(q.Get("i"))
			if err != nil {
				iter = 3
			}
			init := InitRandomPoints
			initRandomSelected := "selected"
			initPlusPlusSelected := ""
			if q.Get("init") == "plusplus" {
				initRandomSelected = ""
				initPlusPlusSelected = "selected"
				init = InitPlusPlus
			}
			data := [][2]float64{}
			var dataSeed uint64 = 123
			var rData *rand.Rand = rand.New(rand.NewPCG(dataSeed, dataSeed))
			randomizeData := false
			randomizeDataChecked := ""
			if randomizeData = q.Get("dataRandom") == "on"; randomizeData {
				seed := rand.Uint64()
				randomizeDataChecked = "checked"
				rData = rand.New(rand.NewPCG(seed, seed))
			}
			var initSeed uint64 = 123
			randomizeInit := false
			randomizeInitChecked := ""
			if randomizeInit = q.Get("initRandom") == "on"; randomizeInit {
				randomizeInitChecked = "checked"
				initSeed = rand.Uint64()
			}
			for range clusters {
				m1, m2 := -15+rData.Float64()*30, -15+rData.Float64()*30
				std1, std2 := 1+rData.Float64()*3, 1+rData.Float64()*3
				for range 1000 {
					data = append(data, [2]float64{rData.NormFloat64()*std1 + m1, rData.NormFloat64()*std2 + m2})
				}
			}
			clust := [][][2]float64{}
			centroidIterations := make([][][2]float64, iter)
			for i := range iter {
				cent, c, err := Kmeans2D(k, data, Config{Seed: &initSeed, Iter: &i, Init: init})
				if err != nil {
					t.Error(err)
				}
				centroidIterations[i] = cent
				clust = c
			}

			d, err := json.Marshal(struct {
				Clusters           [][][2]float64 `json:"clusters"`
				CentroidIterations [][][2]float64 `json:"centroidIterations"`
			}{Clusters: clust, CentroidIterations: centroidIterations})
			if err != nil {
				t.Error(err, centroidIterations)
			}

			fmt.Fprintf(w, html, k, k, clusters, clusters, iter, iter, initRandomSelected, initPlusPlusSelected, randomizeDataChecked, randomizeInitChecked, d)
		})
		http.ListenAndServe(":8080", nil)
	}
}
