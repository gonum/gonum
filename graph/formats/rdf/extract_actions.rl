// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

%%{
	machine extract;

	action StartIRI {
		iri = p
	}

	action EndIRI {
		if iri < 0 {
			panic("unexpected parser state: iri start not set")
		}
		iriText = unEscape(data[iri:p])
		if kind == Invalid {
			kind = IRI
		}
	}

	action StartBlank {
		blank = p
	}

	action EndBlank {
		if blank < 0 {
			panic("unexpected parser state: blank start not set")
		}
		blankText = string(data[blank:p])
		kind = Blank
	}

	action StartLiteral {
		literal = p
	}

	action EndLiteral {
		if literal < 0 {
			panic("unexpected parser state: literal start not set")
		}
		literalText = unEscape(data[literal:p])
		kind = Literal
	}

	action StartLang {
		lang = p
	}

	action EndLang {
		if lang < 0 {
			panic("unexpected parser state: lang start not set")
		}
		langText = string(data[lang:p])
	}

	action Return {
		switch kind {
		case IRI:
			return iriText, "", kind, nil
		case Blank:
			return blankText, "", kind, nil
		case Literal:
			qual = iriText
			if qual == "" {
				qual = langText
			}
			return literalText, qual, kind, nil
		default:
			return "", "", kind, ErrInvalidTerm
		}
	}

	action Error {
		if p < len(data) {
			if r := data[p]; r < unicode.MaxASCII {
				return "", "", Invalid, fmt.Errorf("%w: unexpected rune %q at %d", ErrInvalidTerm, data[p], p)
			} else {
				return "", "", Invalid, fmt.Errorf("%w: unexpected rune %q (\\u%04[2]x) at %d", ErrInvalidTerm, data[p], p)
			}
		}
		return "", "", Invalid, ErrIncompleteTerm
	}
}%%
