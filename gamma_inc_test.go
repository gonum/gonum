// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathext

import (
	"math"
	"testing"
)

func TestGammaInc(t *testing.T) {
	for i, test := range []struct {
		a, x, want float64
	}{
		// Results computed using scipy.special.gamminc
		{0, 0, 0},
		{0.001, 0.005, 0.99528424172333985},
		{0.01, 10, 0.99999995718295021},
		{0.1, 10, 0.99999944520142825},
		{0.5, 2, 0.95449973610364147},
		{1, 0.5, 0.39346934028736652},
		{1, 1, 0.63212055882855778},
		{1.5, 0.75, 0.31772966966378746},
		{2.5, 1, 0.15085496391539038},
		{5, 50, 1},
		{10, 0.9, 4.2519575433351128e-08},
		{10, 5, 0.031828057306204811},
	} {
		if got := GammaInc(test.a, test.x); math.Abs(got-test.want) > 1e-10 {
			t.Errorf("test %d GammaInc(%g, %g) failed: got %g want %g", i, test.a, test.x, got, test.want)
		}
	}
}

func TestGammaIncComp(t *testing.T) {
	for i, test := range []struct {
		a, x, want float64
	}{
		// Results computed using scipy.special.gammincc
		{0.001, 0.005, 0.0047157582766601536},
		{0.01, 0.9, 0.0026263432520514662},
		{0.25, 0.75, 0.10006348671550169},
		{0.5, 0.5, 0.31731050786291404},
		{0.75, 0.25, 0.65343980284081038},
		{0.9, 0.01, 0.98359881081593148},
		{1, 0, 1},
		{1, 1, 0.36787944117144233},
		{1, 10, 4.5399929762484861e-05},
		{5, 1, 0.99634015317265634},
		{5, 10, 0.029252688076961127},
		{100, 10, 1},
	} {
		if got := GammaIncComp(test.a, test.x); math.Abs(got-test.want) > 1e-10 {
			t.Errorf("test %d GammaIncComp(%g, %g) failed: got %g want %g", i, test.a, test.x, got, test.want)
		}
	}
}

func TestGammaIncCompInv(t *testing.T) {
	for i, test := range []struct {
		a, x, want float64
	}{
		// Results computed using scipy.special.gamminccinv
		{0.1, 0.5, 0.00059339110446022798},
		{0.1, 0.75, 5.7917132949696076e-07},
		{0.5, 0.1, 1.3527717270477047},
		{0.25, 0.25, 0.26062600197823282},
		{0.5, 0.5, 0.22746821155978625},
		{0.75, 0.25, 1.0340914067758025},
		{1, 0.5, 0.69314718055994529},
		{1, 0, math.MaxFloat64},
		{1, 1, 0},
		{10, 0.5, 9.6687146147141299},
		{100, 0.25, 106.5510925269767},
		{1000, 0.01, 1075.0328320864389},
	} {
		if got := GammaIncCompInv(test.a, test.x); math.Abs(got-test.want) > 1e-10 {
			t.Errorf("test %d GammaIncCompInv(%g, %g) failed: got %g want %g", i, test.a, test.x, got, test.want)
		}
	}
}
