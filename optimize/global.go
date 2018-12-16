// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"math"
)

// DefaultSettingsGlobal returns the default settings for a global optimization.
func DefaultSettingsGlobal() *Settings {
	return &Settings{
		FunctionThreshold: math.Inf(-1),
		Converger: &FunctionConverge{
			Absolute:   1e-10,
			Iterations: 100,
		},
	}
}
