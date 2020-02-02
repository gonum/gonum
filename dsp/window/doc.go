// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package window provides a set of functions to perform the transformation of sequence by different window functions.
// The main purpose of the window function is to control spectral leakage parameters when performing Fourier transform on a signal of limit length.
// 
// Indicators of quality
// 
// To ensure the selection criteria for a window function, each function is documented by a set of Indicators of quality. Here is a description.
// 
// ΔF_0 - 
// Normalized width of the main AFC lobe at zero amplitude level.
// 
// ΔF_0.5 - 
// Normalized width of the main AFC lobe at 0.5 amplitude level (-6 dB).
// 
// K - 
// For a rectangular window, ΔF_0 has the smallest value and is equal to 2. 
// K Shows how many times the normalized width ΔF_0 of the main lobe by the zero level of the specified window is wider than the rectangular window. 
// Depending on the parameter, the windows are divided into high resolution windows (K≤3) and low resolution windows (K>3).
// 
// ɣ_max - 
// The maximum level of side lobes.
// If the maximum level of side lobes divided by maximum level of main lobe is a g, the ɣ_max is a g expressed in logarithmic scale:  ɣ_max=20*log_10(g), [dB].
// 
// 
// β - 
// The reduction coefficient. 
// The meaning of the reduction coefficient is that the amplitudes of all spectral components after multiplying by the window function are reduced by b times compared to the rectangular window. The reduction coefficient β is a b expressed in logarithmic scale:  β=20*log_10(b), [dB].
package window // import "gonum.org/v1/gonum/dsp/window"
