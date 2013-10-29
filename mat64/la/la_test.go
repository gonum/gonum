// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"github.com/gonum/blas/cblas"
	"github.com/gonum/matrix/mat64"

	check "launchpad.net/gocheck"
	"testing"
)

// Helpers

func mustDense(m *mat64.Dense, err error) *mat64.Dense {
	if err != nil {
		panic(err)
	}
	return m
}

// Tests

func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

func (s *S) SetUpSuite(c *check.C) { mat64.Register(cblas.Blas{}) }

var _ = check.Suite(&S{})
