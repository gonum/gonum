// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file

package cscalar

import (
	"math"
	"math/cmplx"
	"testing"
)

var parseTests = []struct {
	s       string
	want    complex128
	wantErr error
}{
	// Simple error states:
	{s: "", wantErr: parseError{state: -1}},
	{s: "()", wantErr: parseError{string: "()", state: -1}},
	{s: "(1", wantErr: parseError{string: "(1", state: -1}},
	{s: "1)", wantErr: parseError{string: "1)", state: -1}},

	// Ambiguous parse error states:
	{s: "1+2i+3i", wantErr: parseError{string: "1+2i+3i", state: -1}},
	{s: "1e-4i+", wantErr: parseError{string: "1e-4i+", state: -1}},
	{s: "1e-4i-", wantErr: parseError{string: "1e-4i-", state: -1}},

	// Valid input:
	{s: "1+4i", want: 1 + 4i},
	{s: "4i+1", want: 1 + 4i},
	{s: "+1+4i", want: 1 + 4i},
	{s: "+4i+1", want: 1 + 4i},
	{s: ".1+.4i", want: 0.1 + 0.4i},
	{s: ".4i+.1", want: 0.1 + 0.4i},
	{s: "+.1+.4i", want: 0.1 + 0.4i},
	{s: "+.4i+.1", want: 0.1 + 0.4i},
	{s: "1.+4.i", want: 1 + 4i},
	{s: "4.i+1.", want: 1 + 4i},
	{s: "+1.+4.i", want: 1 + 4i},
	{s: "+4.i+1.", want: 1 + 4i},
	{s: "1.0+4.0i", want: 1 + 4i},
	{s: "4.0i+1.0", want: 1 + 4i},
	{s: "+1.0+4.0i", want: 1 + 4i},
	{s: "+4.0i+1.0", want: 1 + 4i},
	{s: "1.0e-4+1i", want: 1e-4 + 1i},
	{s: "1.0e-4+i", want: 1e-4 + 1i},
	{s: "1.0e-4-i", want: 1e-4 - 1i},
	{s: "1.0e-4i-1", want: -1 + 1e-4i},
	{s: "1.0e-4i+1", want: 1 + 1e-4i},
	{s: "1e-4+1i", want: 1e-4 + 1i},
	{s: "1e-4+i", want: 1e-4 + 1i},
	{s: "1e-4-i", want: 1e-4 - 1i},
	{s: "1e-4i-1", want: -1 + 1e-4i},
	{s: "1e-4i+1", want: 1 + 1e-4i},
	{s: "(1+4i)", want: 1 + 4i},
	{s: "(4i+1)", want: 1 + 4i},
	{s: "(+1+4i)", want: 1 + 4i},
	{s: "(+4i+1)", want: 1 + 4i},
	{s: "(1e-4+1i)", want: 1e-4 + 1i},
	{s: "(1e-4+i)", want: 1e-4 + 1i},
	{s: "(1e-4-i)", want: 1e-4 - 1i},
	{s: "(1e-4i-1)", want: -1 + 1e-4i},
	{s: "(1e-4i+1)", want: 1 + 1e-4i},
	{s: "NaN", want: cmplx.NaN()},
	{s: "nan", want: cmplx.NaN()},
	{s: "Inf", want: cmplx.Inf()},
	{s: "inf", want: cmplx.Inf()},
	{s: "(Inf+Infi)", want: complex(math.Inf(1), math.Inf(1))},
	{s: "(-Inf+Infi)", want: complex(math.Inf(-1), math.Inf(1))},
	{s: "(+Inf-Infi)", want: complex(math.Inf(1), math.Inf(-1))},
	{s: "(inf+infi)", want: complex(math.Inf(1), math.Inf(1))},
	{s: "(-inf+infi)", want: complex(math.Inf(-1), math.Inf(1))},
	{s: "(+inf-infi)", want: complex(math.Inf(1), math.Inf(-1))},
	{s: "(nan+nani)", want: complex(math.NaN(), math.NaN())},
	{s: "(nan-nani)", want: complex(math.NaN(), math.NaN())},
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		got, err := parse(test.s)
		if err != test.wantErr {
			t.Errorf("unexpected error for Parse(%q): got:%#v, want:%#v", test.s, err, test.wantErr)
		}
		if err != nil {
			continue
		}
		if !Same(got, test.want) {
			t.Errorf("unexpected result for Parse(%q): got:%v, want:%v", test.s, got, test.want)
		}
	}
}
