// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

%%{
	machine check;

	action Return {
		return nil
	}

	action Error {
		if p < len(data) {
			if r := data[p]; r < unicode.MaxASCII {
				return fmt.Errorf("%w: unexpected rune %q at %d", ErrInvalidTerm, data[p], p)
			} else {
				return fmt.Errorf("%w: unexpected rune %q (\\u%04[2]x) at %d", ErrInvalidTerm, data[p], p)
			}
		}
		return ErrIncompleteTerm
	}
}%%
