// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

%%{
	machine nquads;

	action StartSubject {
		subject = p
	}

	action StartPredicate {
		predicate = p
	}

	action StartObject {
		object = p
	}

	action StartLabel {
		label = p
	}

	action StartIRI {
		iri = p
	}

	action SetSubject {
		if subject < 0 {
			panic("unexpected parser state: subject start not set")
		}
		s.Subject.Value = string(data[subject:p])
	}

	action SetPredicate {
		if predicate < 0 {
			panic("unexpected parser state: predicate start not set")
		}
		s.Predicate.Value = string(data[predicate:p])
	}

	action SetObject {
		if object < 0 {
			panic("unexpected parser state: object start not set")
		}
		s.Object.Value = string(data[object:p])
	}

	action SetLabel {
		if label < 0 {
			panic("unexpected parser state: label start not set")
		}
		s.Label.Value = string(data[label:p])
	}

	action EndIRI {
		if iri < 0 {
			panic("unexpected parser state: iri start not set")
		}
		switch u, err := url.Parse(string(data[iri:p])); {
		case err != nil:
			return s, err
		case !u.IsAbs():
			return s, fmt.Errorf("%w: relative IRI ref %q", ErrInvalid, string(data[iri:p]))
		}
	}

	action Return {
		return s, nil
	}

	action Comment {
	}

	action Error {
		if p < len(data) {
			if r := data[p]; r < unicode.MaxASCII {
				return s, fmt.Errorf("%w: unexpected rune %q at %d", ErrInvalid, data[p], p)
			} else {
				return s, fmt.Errorf("%w: unexpected rune %q (\\u%04[2]x) at %d", ErrInvalid, data[p], p)
			}
		}
		return s, ErrIncomplete
	}
}%%
