// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package window provides a set of functions to perform the transformation
// of sequence by different window functions.
//
// Window functions can be used to control spectral leakage parameters
// when performing a Fourier transform on a signal of limited length.
// See https://en.wikipedia.org/wiki/Window_function for more details.
//
// # Spectral leakage parameters
//
// Application of window functions to an input will result in changes
// to the frequency content of the signal in an effect called spectral
// leakage. See https://en.wikipedia.org/wiki/Spectral_leakage.
//
// The characteristic changes associated with each window function may
// be described using a set of spectral leakage parameters; β, ΔF_0, ΔF_0.5,
// K and ɣ_max.
//
// The β, attenuation, coefficient of a window is the ratio of the
// constant component of the spectrum resulting from use of the window
// compared to that produced using the rectangular window, expressed in
// a logarithmic scale.
//
//	β_w = 20 log10(A_w / A_rect) dB
//
// The ΔF_0 parameter describes the normalized width of the main lobe of
// the frequency spectrum at zero amplitude.
//
// The ΔF_0.5 parameter describes the normalized width of the main lobe of
// the frequency spectrum at -3 dB (half maximum amplitude).
//
// The K parameter describes the relative width of the main lobe of the
// frequency spectrum produced by the window compared with the rectangular
// window. The rectangular window has the lowest ΔF_0 at a value of 2.
//
//	K_w = ΔF_0_w/ΔF_0_rect.
//
// The value of K divides windows into high resolution windows (K≤3) and
// low resolution windows (K>3).
//
// The ɣ_max parameter is the maximum level of the side lobes of the
// normalized spectrum, in decibels.
package window // import "gonum.org/v1/gonum/dsp/window"

// The article at http://www.dsplib.ru/content/win/win.html is the origin
// of much of the information used in this package. It is in Russian, but
// Google translate does a pretty good job with it.
