// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright © 2008 World Wide Web Consortium, (MIT, ERCIM, Keio, Beihang)
// and others. All Rights Reserved.
// http://www.w3.org/Consortium/Legal/2008/04-testsuite-copyright.html
// Used under https://www.w3.org/Consortium/Legal/2008/03-bsd-license.

package rdf

type statement struct {
	input string

	subject, predicate, object, label term
}

type term struct {
	text string
	qual string
	kind Kind
}

// Test suite values were extracted from the test case archives in this directory.
// The archives were obtained from https://w3c.github.io/rdf-tests/ntriples/ and
// https://w3c.github.io/rdf-tests/nquads/.
var testSuite = map[string][]statement{
	"comment_following_triple.nq": {
		{
			input:     "<http://example/s> <http://example/p> <http://example/o> . # comment",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
		},
		{
			input:     "<http://example/s> <http://example/p> _:o . # comment",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "o", kind: Blank},
		},
		{
			input:     "<http://example/s> <http://example/p> \"o\" . # comment",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "o", kind: Literal},
		},
		{
			input:     "<http://example/s> <http://example/p> \"o\"^^<http://example/dt> . # comment",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "o", qual: "http://example/dt", kind: Literal},
		},
		{
			input:     "<http://example/s> <http://example/p> \"o\"@en . # comment",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "o", qual: "@en", kind: Literal},
		},
	},
	"comment_following_triple.nt": {
		{
			input:     "<http://example/s> <http://example/p> <http://example/o> . # comment",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
		},
		{
			input:     "<http://example/s> <http://example/p> _:o . # comment",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "o", kind: Blank},
		},
		{
			input:     "<http://example/s> <http://example/p> \"o\" . # comment",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "o", kind: Literal},
		},
		{
			input:     "<http://example/s> <http://example/p> \"o\"^^<http://example/dt> . # comment",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "o", qual: "http://example/dt", kind: Literal},
		},
		{
			input:     "<http://example/s> <http://example/p> \"o\"@en . # comment",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "o", qual: "@en", kind: Literal},
		},
	},
	"langtagged_string.nq": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"chat\"@en .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "chat", qual: "@en", kind: Literal},
		},
	},
	"langtagged_string.nt": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"chat\"@en .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "chat", qual: "@en", kind: Literal},
		},
	},
	"lantag_with_subtag.nq": {
		{
			input:     "<http://example.org/ex#a> <http://example.org/ex#b> \"Cheers\"@en-UK .",
			subject:   term{text: "http://example.org/ex#a", kind: IRI},
			predicate: term{text: "http://example.org/ex#b", kind: IRI},
			object:    term{text: "Cheers", qual: "@en-UK", kind: Literal},
		},
	},
	"lantag_with_subtag.nt": {
		{
			input:     "<http://example.org/ex#a> <http://example.org/ex#b> \"Cheers\"@en-UK .",
			subject:   term{text: "http://example.org/ex#a", kind: IRI},
			predicate: term{text: "http://example.org/ex#b", kind: IRI},
			object:    term{text: "Cheers", qual: "@en-UK", kind: Literal},
		},
	},
	"literal.nq": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"x\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "x", kind: Literal},
		},
	},
	"literal.nt": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"x\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "x", kind: Literal},
		},
	},
	"literal_all_controls.nq": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\\u0000\\u0001\\u0002\\u0003\\u0004\\u0005\\u0006\\u0007\\u0008\\t\\u000B\\u000C\\u000E\\u000F\\u0010\\u0011\\u0012\\u0013\\u0014\\u0015\\u0016\\u0017\\u0018\\u0019\\u001A\\u001B\\u001C\\u001D\\u001E\\u001F\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "\x00\x01\x02\x03\x04\x05\x06\a\b\t\v\f\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f", kind: Literal},
		},
	},
	"literal_all_controls.nt": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\\u0000\\u0001\\u0002\\u0003\\u0004\\u0005\\u0006\\u0007\\u0008\\t\\u000B\\u000C\\u000E\\u000F\\u0010\\u0011\\u0012\\u0013\\u0014\\u0015\\u0016\\u0017\\u0018\\u0019\\u001A\\u001B\\u001C\\u001D\\u001E\\u001F\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "\x00\x01\x02\x03\x04\x05\x06\a\b\t\v\f\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f", kind: Literal},
		},
	},
	"literal_all_punctuation.nq": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \" !\\\"#$%&():;<=>?@[]^_`{|}~\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: " !\"#$%&():;<=>?@[]^_`{|}~", kind: Literal},
		},
	},
	"literal_all_punctuation.nt": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \" !\\\"#$%&():;<=>?@[]^_`{|}~\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: " !\"#$%&():;<=>?@[]^_`{|}~", kind: Literal},
		},
	},
	"literal_ascii_boundaries.nq": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\x00\t\v\f\x0e&([]\u007f\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "\x00\t\v\f\x0e&([]\u007f", kind: Literal},
		},
	},
	"literal_ascii_boundaries.nt": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\x00\t\v\f\x0e&([]\u007f\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "\x00\t\v\f\x0e&([]\u007f", kind: Literal},
		},
	},
	"literal_false.nq": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"false\"^^<http://www.w3.org/2001/XMLSchema#boolean> .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "false", qual: "http://www.w3.org/2001/XMLSchema#boolean", kind: Literal},
		},
	},
	"literal_false.nt": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"false\"^^<http://www.w3.org/2001/XMLSchema#boolean> .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "false", qual: "http://www.w3.org/2001/XMLSchema#boolean", kind: Literal},
		},
	},
	"literal_true.nq": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"true\"^^<http://www.w3.org/2001/XMLSchema#boolean> .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "true", qual: "http://www.w3.org/2001/XMLSchema#boolean", kind: Literal},
		},
	},
	"literal_true.nt": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"true\"^^<http://www.w3.org/2001/XMLSchema#boolean> .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "true", qual: "http://www.w3.org/2001/XMLSchema#boolean", kind: Literal},
		},
	},
	"literal_with_2_dquotes.nq": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"x\\\"\\\"y\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "x\"\"y", kind: Literal},
		},
	},
	"literal_with_2_dquotes.nt": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"x\\\"\\\"y\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "x\"\"y", kind: Literal},
		},
	},
	"literal_with_2_squotes.nt": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"x''y\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "x''y", kind: Literal},
		},
	},
	"literal_with_2_squotes.nq": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"x''y\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "x''y", kind: Literal},
		},
	},
	"literal_with_BACKSPACE.nq": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\\b\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "\b", kind: Literal},
		},
	},
	"literal_with_BACKSPACE.nt": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\\b\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "\b", kind: Literal},
		},
	},
	"literal_with_CARRIAGE_RETURN.nq": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\\r\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "\r", kind: Literal},
		},
	},
	"literal_with_CARRIAGE_RETURN.nt": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\\r\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "\r", kind: Literal},
		},
	},
	"literal_with_CHARACTER_TABULATION.nq": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\\t\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "\t", kind: Literal},
		},
	},
	"literal_with_CHARACTER_TABULATION.nt": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\\t\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "\t", kind: Literal},
		},
	},
	"literal_with_FORM_FEED.nq": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\\f\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "\f", kind: Literal},
		},
	},
	"literal_with_FORM_FEED.nt": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\\f\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "\f", kind: Literal},
		},
	},
	"literal_with_LINE_FEED.nq": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\\n\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "\n", kind: Literal},
		},
	},
	"literal_with_LINE_FEED.nt": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\\n\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "\n", kind: Literal},
		},
	},
	"literal_with_REVERSE_SOLIDUS.nq": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\\\\\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "\\", kind: Literal},
		},
	},
	"literal_with_REVERSE_SOLIDUS.nt": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\\\\\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "\\", kind: Literal},
		},
	},
	"literal_with_REVERSE_SOLIDUS2.nq": {
		{
			input:     "<http://example.org/ns#s> <http://example.org/ns#p1> \"test-\\\\\" .",
			subject:   term{text: "http://example.org/ns#s", kind: IRI},
			predicate: term{text: "http://example.org/ns#p1", kind: IRI},
			object:    term{text: "test-\\", kind: Literal},
		},
	},
	"literal_with_REVERSE_SOLIDUS2.nt": {
		{
			input:     "<http://example.org/ns#s> <http://example.org/ns#p1> \"test-\\\\\" .",
			subject:   term{text: "http://example.org/ns#s", kind: IRI},
			predicate: term{text: "http://example.org/ns#p1", kind: IRI},
			object:    term{text: "test-\\", kind: Literal},
		},
	},
	"literal_with_UTF8_boundaries.nq": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\u0080߿ࠀ\u0fffက쿿퀀\ud7ff\ue000�𐀀\U0003fffd\U00040000\U000ffffd\U00100000\U0010fffd\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "\u0080߿ࠀ\u0fffက쿿퀀\ud7ff\ue000�𐀀\U0003fffd\U00040000\U000ffffd\U00100000\U0010fffd", kind: Literal},
		},
	},
	"literal_with_UTF8_boundaries.nt": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\u0080߿ࠀ\u0fffက쿿퀀\ud7ff\ue000�𐀀\U0003fffd\U00040000\U000ffffd\U00100000\U0010fffd\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "\u0080߿ࠀ\u0fffက쿿퀀\ud7ff\ue000�𐀀\U0003fffd\U00040000\U000ffffd\U00100000\U0010fffd", kind: Literal},
		},
	},
	"literal_with_dquote.nq": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"x\\\"y\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "x\"y", kind: Literal},
		},
	},
	"literal_with_dquote.nt": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"x\\\"y\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "x\"y", kind: Literal},
		},
	},
	"literal_with_numeric_escape4.nq": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\\u006F\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "o", kind: Literal},
		},
	},
	"literal_with_numeric_escape4.nt": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\\u006F\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "o", kind: Literal},
		},
	},
	"literal_with_numeric_escape8.nq": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\\U0000006F\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "o", kind: Literal},
		},
	},
	"literal_with_numeric_escape8.nt": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"\\U0000006F\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "o", kind: Literal},
		},
	},
	"literal_with_squote.nq": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"x'y\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "x'y", kind: Literal},
		},
	},
	"literal_with_squote.nt": {
		{
			input:     "<http://a.example/s> <http://a.example/p> \"x'y\" .",
			subject:   term{text: "http://a.example/s", kind: IRI},
			predicate: term{text: "http://a.example/p", kind: IRI},
			object:    term{text: "x'y", kind: Literal},
		},
	},
	"minimal_whitespace.nq": {
		{
			input:     "<http://example/s><http://example/p><http://example/o>.",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
		},
		{
			input:     "<http://example/s><http://example/p>\"Alice\".",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "Alice", kind: Literal},
		},
		{
			input:     "<http://example/s><http://example/p>_:o.",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "o", kind: Blank},
		},
		{
			input:     "_:s<http://example/p><http://example/o>.",
			subject:   term{text: "s", kind: Blank},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
		},
		{
			input:     "_:s<http://example/p>\"Alice\".",
			subject:   term{text: "s", kind: Blank},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "Alice", kind: Literal},
		},
		{
			input:     "_:s<http://example/p>_:bnode1.",
			subject:   term{text: "s", kind: Blank},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "bnode1", kind: Blank},
		},
	},
	"minimal_whitespace.nt": {
		{
			input:     "<http://example/s><http://example/p><http://example/o>.",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
		},
		{
			input:     "<http://example/s><http://example/p>\"Alice\".",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "Alice", kind: Literal},
		},
		{
			input:     "<http://example/s><http://example/p>_:o.",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "o", kind: Blank},
		},
		{
			input:     "_:s<http://example/p><http://example/o>.",
			subject:   term{text: "s", kind: Blank},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
		},
		{
			input:     "_:s<http://example/p>\"Alice\".",
			subject:   term{text: "s", kind: Blank},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "Alice", kind: Literal},
		},
		{
			input:     "_:s<http://example/p>_:bnode1.",
			subject:   term{text: "s", kind: Blank},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "bnode1", kind: Blank},
		},
	},
	"nq-syntax-bnode-01.nq": {
		{
			input:     "<http://example/s> <http://example/p> <http://example/o> _:g .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
			label:     term{text: "g", kind: Blank},
		},
	},
	"nq-syntax-bnode-02.nq": {
		{
			input:     "_:s <http://example/p> <http://example/o> _:g .",
			subject:   term{text: "s", kind: Blank},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
			label:     term{text: "g", kind: Blank},
		},
	},
	"nq-syntax-bnode-03.nq": {
		{
			input:     "<http://example/s> <http://example/p> _:o _:g .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "o", kind: Blank},
			label:     term{text: "g", kind: Blank},
		},
	},
	"nq-syntax-bnode-04.nq": {
		{
			input:     "<http://example/s> <http://example/p> \"o\" _:g .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "o", kind: Literal},
			label:     term{text: "g", kind: Blank},
		},
	},
	"nq-syntax-bnode-05.nq": {
		{
			input:     "<http://example/s> <http://example/p> \"o\"@en _:g .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "o", qual: "@en", kind: Literal},
			label:     term{text: "g", kind: Blank},
		},
	},
	"nq-syntax-bnode-06.nq": {
		{
			input:     "<http://example/s> <http://example/p> \"o\"^^<http://www.w3.org/2001/XMLSchema#string> _:g .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "o", qual: "http://www.w3.org/2001/XMLSchema#string", kind: Literal},
			label:     term{text: "g", kind: Blank},
		},
	},
	"nq-syntax-uri-01.nq": {
		{
			input:     "<http://example/s> <http://example/p> <http://example/o> <http://example/g> .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
			label:     term{text: "http://example/g", kind: IRI},
		},
	},
	"nq-syntax-uri-02.nq": {
		{
			input:     "_:s <http://example/p> <http://example/o> <http://example/g> .",
			subject:   term{text: "s", kind: Blank},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
			label:     term{text: "http://example/g", kind: IRI},
		},
	},
	"nq-syntax-uri-03.nq": {
		{
			input:     "<http://example/s> <http://example/p> _:o <http://example/g> .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "o", kind: Blank},
			label:     term{text: "http://example/g", kind: IRI},
		},
	},
	"nq-syntax-uri-04.nq": {
		{
			input:     "<http://example/s> <http://example/p> \"o\" <http://example/g> .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "o", kind: Literal},
			label:     term{text: "http://example/g", kind: IRI},
		},
	},
	"nq-syntax-uri-05.nq": {
		{
			input:     "<http://example/s> <http://example/p> \"o\"@en <http://example/g> .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "o", qual: "@en", kind: Literal},
			label:     term{text: "http://example/g", kind: IRI},
		},
	},
	"nq-syntax-uri-06.nq": {
		{
			input:     "<http://example/s> <http://example/p> \"o\"^^<http://www.w3.org/2001/XMLSchema#string> <http://example/g> .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "o", qual: "http://www.w3.org/2001/XMLSchema#string", kind: Literal},
			label:     term{text: "http://example/g", kind: IRI},
		},
	},
	"nt-syntax-bnode-01.nq": {
		{
			input:     "_:a  <http://example/p> <http://example/o> .",
			subject:   term{text: "a", kind: Blank},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
		},
	},
	"nt-syntax-bnode-01.nt": {
		{
			input:     "_:a  <http://example/p> <http://example/o> .",
			subject:   term{text: "a", kind: Blank},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
		},
	},
	"nt-syntax-bnode-02.nq": {
		{
			input:     "<http://example/s> <http://example/p> _:a .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "a", kind: Blank},
		},
		{
			input:     "_:a  <http://example/p> <http://example/o> .",
			subject:   term{text: "a", kind: Blank},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
		},
	},
	"nt-syntax-bnode-02.nt": {
		{
			input:     "<http://example/s> <http://example/p> _:a .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "a", kind: Blank},
		},
		{
			input:     "_:a  <http://example/p> <http://example/o> .",
			subject:   term{text: "a", kind: Blank},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
		},
	},
	"nt-syntax-bnode-03.nq": {
		{
			input:     "<http://example/s> <http://example/p> _:1a .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "1a", kind: Blank},
		},
		{
			input:     "_:1a  <http://example/p> <http://example/o> .",
			subject:   term{text: "1a", kind: Blank},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
		},
	},
	"nt-syntax-bnode-03.nt": {
		{
			input:     "<http://example/s> <http://example/p> _:1a .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "1a", kind: Blank},
		},
		{
			input:     "_:1a  <http://example/p> <http://example/o> .",
			subject:   term{text: "1a", kind: Blank},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
		},
	},
	"nt-syntax-datatypes-01.nq": {
		{
			input:     "<http://example/s> <http://example/p> \"123\"^^<http://www.w3.org/2001/XMLSchema#byte> .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "123", qual: "http://www.w3.org/2001/XMLSchema#byte", kind: Literal},
		},
	},
	"nt-syntax-datatypes-01.nt": {
		{
			input:     "<http://example/s> <http://example/p> \"123\"^^<http://www.w3.org/2001/XMLSchema#byte> .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "123", qual: "http://www.w3.org/2001/XMLSchema#byte", kind: Literal},
		},
	},
	"nt-syntax-datatypes-02.nq": {
		{
			input:     "<http://example/s> <http://example/p> \"123\"^^<http://www.w3.org/2001/XMLSchema#string> .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "123", qual: "http://www.w3.org/2001/XMLSchema#string", kind: Literal},
		},
	},
	"nt-syntax-datatypes-02.nt": {
		{
			input:     "<http://example/s> <http://example/p> \"123\"^^<http://www.w3.org/2001/XMLSchema#string> .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "123", qual: "http://www.w3.org/2001/XMLSchema#string", kind: Literal},
		},
	},
	"nt-syntax-str-esc-01.nq": {
		{
			input:     "<http://example/s> <http://example/p> \"a\\n\" .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "a\n", kind: Literal},
		},
	},
	"nt-syntax-str-esc-01.nt": {
		{
			input:     "<http://example/s> <http://example/p> \"a\\n\" .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "a\n", kind: Literal},
		},
	},
	"nt-syntax-str-esc-02.nq": {
		{
			input:     "<http://example/s> <http://example/p> \"a\\u0020b\" .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "a b", kind: Literal},
		},
	},
	"nt-syntax-str-esc-02.nt": {
		{
			input:     "<http://example/s> <http://example/p> \"a\\u0020b\" .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "a b", kind: Literal},
		},
	},
	"nt-syntax-str-esc-03.nq": {
		{
			input:     "<http://example/s> <http://example/p> \"a\\U00000020b\" .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "a b", kind: Literal},
		},
	},
	"nt-syntax-str-esc-03.nt": {
		{
			input:     "<http://example/s> <http://example/p> \"a\\U00000020b\" .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "a b", kind: Literal},
		},
	},
	"nt-syntax-string-01.nq": {
		{
			input:     "<http://example/s> <http://example/p> \"string\" .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "string", kind: Literal},
		},
	},
	"nt-syntax-string-01.nt": {
		{
			input:     "<http://example/s> <http://example/p> \"string\" .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "string", kind: Literal},
		},
	},
	"nt-syntax-string-02.nq": {
		{
			input:     "<http://example/s> <http://example/p> \"string\"@en .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "string", qual: "@en", kind: Literal},
		},
	},
	"nt-syntax-string-02.nt": {
		{
			input:     "<http://example/s> <http://example/p> \"string\"@en .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "string", qual: "@en", kind: Literal},
		},
	},
	"nt-syntax-string-03.nq": {
		{
			input:     "<http://example/s> <http://example/p> \"string\"@en-uk .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "string", qual: "@en-uk", kind: Literal},
		},
	},
	"nt-syntax-string-03.nt": {
		{
			input:     "<http://example/s> <http://example/p> \"string\"@en-uk .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "string", qual: "@en-uk", kind: Literal},
		},
	},
	"nt-syntax-subm-01.nq": {
		{
			input:     "<http://example.org/resource1> <http://example.org/property> <http://example.org/resource2> .",
			subject:   term{text: "http://example.org/resource1", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "http://example.org/resource2", kind: IRI},
		},
		{
			input:     "_:anon <http://example.org/property> <http://example.org/resource2> .",
			subject:   term{text: "anon", kind: Blank},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "http://example.org/resource2", kind: IRI},
		},
		{
			input:     "<http://example.org/resource2> <http://example.org/property> _:anon .",
			subject:   term{text: "http://example.org/resource2", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "anon", kind: Blank},
		},
		{
			input:     "<http://example.org/resource3> \t <http://example.org/property>\t <http://example.org/resource2> \t.",
			subject:   term{text: "http://example.org/resource3", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "http://example.org/resource2", kind: IRI},
		},
		{
			input:     "<http://example.org/resource4> <http://example.org/property> <http://example.org/resource2> .",
			subject:   term{text: "http://example.org/resource4", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "http://example.org/resource2", kind: IRI},
		},
		{
			input:     "<http://example.org/resource5> <http://example.org/property> <http://example.org/resource2> .",
			subject:   term{text: "http://example.org/resource5", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "http://example.org/resource2", kind: IRI},
		},
		{
			input:     "<http://example.org/resource6> <http://example.org/property> <http://example.org/resource2> .",
			subject:   term{text: "http://example.org/resource6", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "http://example.org/resource2", kind: IRI},
		},
		{
			input:     "<http://example.org/resource7> <http://example.org/property> \"simple literal\" .",
			subject:   term{text: "http://example.org/resource7", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "simple literal", kind: Literal},
		},
		{
			input:     "<http://example.org/resource8> <http://example.org/property> \"backslash:\\\\\" .",
			subject:   term{text: "http://example.org/resource8", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "backslash:\\", kind: Literal},
		},
		{
			input:     "<http://example.org/resource9> <http://example.org/property> \"dquote:\\\"\" .",
			subject:   term{text: "http://example.org/resource9", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "dquote:\"", kind: Literal},
		},
		{
			input:     "<http://example.org/resource10> <http://example.org/property> \"newline:\\n\" .",
			subject:   term{text: "http://example.org/resource10", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "newline:\n", kind: Literal},
		},
		{
			input:     "<http://example.org/resource11> <http://example.org/property> \"return\\r\" .",
			subject:   term{text: "http://example.org/resource11", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "return\r", kind: Literal},
		},
		{
			input:     "<http://example.org/resource12> <http://example.org/property> \"tab:\\t\" .",
			subject:   term{text: "http://example.org/resource12", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "tab:\t", kind: Literal},
		},
		{
			input:     "<http://example.org/resource13> <http://example.org/property> <http://example.org/resource2>.",
			subject:   term{text: "http://example.org/resource13", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "http://example.org/resource2", kind: IRI},
		},
		{
			input:     "<http://example.org/resource14> <http://example.org/property> \"x\".",
			subject:   term{text: "http://example.org/resource14", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "x", kind: Literal},
		},
		{
			input:     "<http://example.org/resource15> <http://example.org/property> _:anon.",
			subject:   term{text: "http://example.org/resource15", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "anon", kind: Blank},
		},
		{
			input:     "<http://example.org/resource16> <http://example.org/property> \"\\u00E9\" .",
			subject:   term{text: "http://example.org/resource16", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "é", kind: Literal},
		},
		{
			input:     "<http://example.org/resource17> <http://example.org/property> \"\\u20AC\" .",
			subject:   term{text: "http://example.org/resource17", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "€", kind: Literal},
		},
		{
			input:     "<http://example.org/resource21> <http://example.org/property> \"\"^^<http://www.w3.org/2000/01/rdf-schema#XMLLiteral> .",
			subject:   term{text: "http://example.org/resource21", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "", qual: "http://www.w3.org/2000/01/rdf-schema#XMLLiteral", kind: Literal},
		},
		{
			input:     "<http://example.org/resource22> <http://example.org/property> \" \"^^<http://www.w3.org/2000/01/rdf-schema#XMLLiteral> .",
			subject:   term{text: "http://example.org/resource22", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: " ", qual: "http://www.w3.org/2000/01/rdf-schema#XMLLiteral", kind: Literal},
		},
		{
			input:     "<http://example.org/resource23> <http://example.org/property> \"x\"^^<http://www.w3.org/2000/01/rdf-schema#XMLLiteral> .",
			subject:   term{text: "http://example.org/resource23", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "x", qual: "http://www.w3.org/2000/01/rdf-schema#XMLLiteral", kind: Literal},
		},
		{
			input:     "<http://example.org/resource23> <http://example.org/property> \"\\\"\"^^<http://www.w3.org/2000/01/rdf-schema#XMLLiteral> .",
			subject:   term{text: "http://example.org/resource23", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "\"", qual: "http://www.w3.org/2000/01/rdf-schema#XMLLiteral", kind: Literal},
		},
		{
			input:     "<http://example.org/resource24> <http://example.org/property> \"<a></a>\"^^<http://www.w3.org/2000/01/rdf-schema#XMLLiteral> .",
			subject:   term{text: "http://example.org/resource24", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "<a></a>", qual: "http://www.w3.org/2000/01/rdf-schema#XMLLiteral", kind: Literal},
		},
		{
			input:     "<http://example.org/resource25> <http://example.org/property> \"a <b></b>\"^^<http://www.w3.org/2000/01/rdf-schema#XMLLiteral> .",
			subject:   term{text: "http://example.org/resource25", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "a <b></b>", qual: "http://www.w3.org/2000/01/rdf-schema#XMLLiteral", kind: Literal},
		},
		{
			input:     "<http://example.org/resource26> <http://example.org/property> \"a <b></b> c\"^^<http://www.w3.org/2000/01/rdf-schema#XMLLiteral> .",
			subject:   term{text: "http://example.org/resource26", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "a <b></b> c", qual: "http://www.w3.org/2000/01/rdf-schema#XMLLiteral", kind: Literal},
		},
		{
			input:     "<http://example.org/resource26> <http://example.org/property> \"a\\n<b></b>\\nc\"^^<http://www.w3.org/2000/01/rdf-schema#XMLLiteral> .",
			subject:   term{text: "http://example.org/resource26", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "a\n<b></b>\nc", qual: "http://www.w3.org/2000/01/rdf-schema#XMLLiteral", kind: Literal},
		},
		{
			input:     "<http://example.org/resource27> <http://example.org/property> \"chat\"^^<http://www.w3.org/2000/01/rdf-schema#XMLLiteral> .",
			subject:   term{text: "http://example.org/resource27", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "chat", qual: "http://www.w3.org/2000/01/rdf-schema#XMLLiteral", kind: Literal},
		},
		{
			input:     "<http://example.org/resource30> <http://example.org/property> \"chat\"@fr .",
			subject:   term{text: "http://example.org/resource30", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "chat", qual: "@fr", kind: Literal},
		},
		{
			input:     "<http://example.org/resource31> <http://example.org/property> \"chat\"@en .",
			subject:   term{text: "http://example.org/resource31", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "chat", qual: "@en", kind: Literal},
		},
		{
			input:     "<http://example.org/resource32> <http://example.org/property> \"abc\"^^<http://example.org/datatype1> .",
			subject:   term{text: "http://example.org/resource32", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "abc", qual: "http://example.org/datatype1", kind: Literal},
		},
	},
	"nt-syntax-subm-01.nt": {
		{
			input:     "<http://example.org/resource1> <http://example.org/property> <http://example.org/resource2> .",
			subject:   term{text: "http://example.org/resource1", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "http://example.org/resource2", kind: IRI},
		},
		{
			input:     "_:anon <http://example.org/property> <http://example.org/resource2> .",
			subject:   term{text: "anon", kind: Blank},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "http://example.org/resource2", kind: IRI},
		},
		{
			input:     "<http://example.org/resource2> <http://example.org/property> _:anon .",
			subject:   term{text: "http://example.org/resource2", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "anon", kind: Blank},
		},
		{
			input:     "<http://example.org/resource3> \t <http://example.org/property>\t <http://example.org/resource2> \t.",
			subject:   term{text: "http://example.org/resource3", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "http://example.org/resource2", kind: IRI},
		},
		{
			input:     "<http://example.org/resource4> <http://example.org/property> <http://example.org/resource2> .",
			subject:   term{text: "http://example.org/resource4", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "http://example.org/resource2", kind: IRI},
		},
		{
			input:     "<http://example.org/resource5> <http://example.org/property> <http://example.org/resource2> .",
			subject:   term{text: "http://example.org/resource5", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "http://example.org/resource2", kind: IRI},
		},
		{
			input:     "<http://example.org/resource6> <http://example.org/property> <http://example.org/resource2> .",
			subject:   term{text: "http://example.org/resource6", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "http://example.org/resource2", kind: IRI},
		},
		{
			input:     "<http://example.org/resource7> <http://example.org/property> \"simple literal\" .",
			subject:   term{text: "http://example.org/resource7", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "simple literal", kind: Literal},
		},
		{
			input:     "<http://example.org/resource8> <http://example.org/property> \"backslash:\\\\\" .",
			subject:   term{text: "http://example.org/resource8", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "backslash:\\", kind: Literal},
		},
		{
			input:     "<http://example.org/resource9> <http://example.org/property> \"dquote:\\\"\" .",
			subject:   term{text: "http://example.org/resource9", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "dquote:\"", kind: Literal},
		},
		{
			input:     "<http://example.org/resource10> <http://example.org/property> \"newline:\\n\" .",
			subject:   term{text: "http://example.org/resource10", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "newline:\n", kind: Literal},
		},
		{
			input:     "<http://example.org/resource11> <http://example.org/property> \"return\\r\" .",
			subject:   term{text: "http://example.org/resource11", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "return\r", kind: Literal},
		},
		{
			input:     "<http://example.org/resource12> <http://example.org/property> \"tab:\\t\" .",
			subject:   term{text: "http://example.org/resource12", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "tab:\t", kind: Literal},
		},
		{
			input:     "<http://example.org/resource13> <http://example.org/property> <http://example.org/resource2>.",
			subject:   term{text: "http://example.org/resource13", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "http://example.org/resource2", kind: IRI},
		},
		{
			input:     "<http://example.org/resource14> <http://example.org/property> \"x\".",
			subject:   term{text: "http://example.org/resource14", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "x", kind: Literal},
		},
		{
			input:     "<http://example.org/resource15> <http://example.org/property> _:anon.",
			subject:   term{text: "http://example.org/resource15", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "anon", kind: Blank},
		},
		{
			input:     "<http://example.org/resource16> <http://example.org/property> \"\\u00E9\" .",
			subject:   term{text: "http://example.org/resource16", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "é", kind: Literal},
		},
		{
			input:     "<http://example.org/resource17> <http://example.org/property> \"\\u20AC\" .",
			subject:   term{text: "http://example.org/resource17", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "€", kind: Literal},
		},
		{
			input:     "<http://example.org/resource21> <http://example.org/property> \"\"^^<http://www.w3.org/2000/01/rdf-schema#XMLLiteral> .",
			subject:   term{text: "http://example.org/resource21", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "", qual: "http://www.w3.org/2000/01/rdf-schema#XMLLiteral", kind: Literal},
		},
		{
			input:     "<http://example.org/resource22> <http://example.org/property> \" \"^^<http://www.w3.org/2000/01/rdf-schema#XMLLiteral> .",
			subject:   term{text: "http://example.org/resource22", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: " ", qual: "http://www.w3.org/2000/01/rdf-schema#XMLLiteral", kind: Literal},
		},
		{
			input:     "<http://example.org/resource23> <http://example.org/property> \"x\"^^<http://www.w3.org/2000/01/rdf-schema#XMLLiteral> .",
			subject:   term{text: "http://example.org/resource23", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "x", qual: "http://www.w3.org/2000/01/rdf-schema#XMLLiteral", kind: Literal},
		},
		{
			input:     "<http://example.org/resource23> <http://example.org/property> \"\\\"\"^^<http://www.w3.org/2000/01/rdf-schema#XMLLiteral> .",
			subject:   term{text: "http://example.org/resource23", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "\"", qual: "http://www.w3.org/2000/01/rdf-schema#XMLLiteral", kind: Literal},
		},
		{
			input:     "<http://example.org/resource24> <http://example.org/property> \"<a></a>\"^^<http://www.w3.org/2000/01/rdf-schema#XMLLiteral> .",
			subject:   term{text: "http://example.org/resource24", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "<a></a>", qual: "http://www.w3.org/2000/01/rdf-schema#XMLLiteral", kind: Literal},
		},
		{
			input:     "<http://example.org/resource25> <http://example.org/property> \"a <b></b>\"^^<http://www.w3.org/2000/01/rdf-schema#XMLLiteral> .",
			subject:   term{text: "http://example.org/resource25", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "a <b></b>", qual: "http://www.w3.org/2000/01/rdf-schema#XMLLiteral", kind: Literal},
		},
		{
			input:     "<http://example.org/resource26> <http://example.org/property> \"a <b></b> c\"^^<http://www.w3.org/2000/01/rdf-schema#XMLLiteral> .",
			subject:   term{text: "http://example.org/resource26", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "a <b></b> c", qual: "http://www.w3.org/2000/01/rdf-schema#XMLLiteral", kind: Literal},
		},
		{
			input:     "<http://example.org/resource26> <http://example.org/property> \"a\\n<b></b>\\nc\"^^<http://www.w3.org/2000/01/rdf-schema#XMLLiteral> .",
			subject:   term{text: "http://example.org/resource26", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "a\n<b></b>\nc", qual: "http://www.w3.org/2000/01/rdf-schema#XMLLiteral", kind: Literal},
		},
		{
			input:     "<http://example.org/resource27> <http://example.org/property> \"chat\"^^<http://www.w3.org/2000/01/rdf-schema#XMLLiteral> .",
			subject:   term{text: "http://example.org/resource27", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "chat", qual: "http://www.w3.org/2000/01/rdf-schema#XMLLiteral", kind: Literal},
		},
		{
			input:     "<http://example.org/resource30> <http://example.org/property> \"chat\"@fr .",
			subject:   term{text: "http://example.org/resource30", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "chat", qual: "@fr", kind: Literal},
		},
		{
			input:     "<http://example.org/resource31> <http://example.org/property> \"chat\"@en .",
			subject:   term{text: "http://example.org/resource31", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "chat", qual: "@en", kind: Literal},
		},
		{
			input:     "<http://example.org/resource32> <http://example.org/property> \"abc\"^^<http://example.org/datatype1> .",
			subject:   term{text: "http://example.org/resource32", kind: IRI},
			predicate: term{text: "http://example.org/property", kind: IRI},
			object:    term{text: "abc", qual: "http://example.org/datatype1", kind: Literal},
		},
	},
	"nt-syntax-uri-01.nq": {
		{
			input:     "<http://example/s> <http://example/p> <http://example/o> .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
		},
	},
	"nt-syntax-uri-01.nt": {
		{
			input:     "<http://example/s> <http://example/p> <http://example/o> .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
		},
	},
	"nt-syntax-uri-02.nq": {
		{
			input:     "<http://example/\\u0053> <http://example/p> <http://example/o> .",
			subject:   term{text: "http://example/S", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
		},
	},
	"nt-syntax-uri-02.nt": {
		{
			input:     "<http://example/\\u0053> <http://example/p> <http://example/o> .",
			subject:   term{text: "http://example/S", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
		},
	},
	"nt-syntax-uri-03.nq": {
		{
			input:     "<http://example/\\U00000053> <http://example/p> <http://example/o> .",
			subject:   term{text: "http://example/S", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
		},
	},
	"nt-syntax-uri-03.nt": {
		{
			input:     "<http://example/\\U00000053> <http://example/p> <http://example/o> .",
			subject:   term{text: "http://example/S", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "http://example/o", kind: IRI},
		},
	},
	"nt-syntax-uri-04.nq": {
		{
			input:     "<http://example/s> <http://example/p> <scheme:!$%25&'()*+,-./0123456789:/@ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz~?#> .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "scheme:!$%25&'()*+,-./0123456789:/@ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz~?#", kind: IRI},
		},
	},
	"nt-syntax-uri-04.nt": {
		{
			input:     "<http://example/s> <http://example/p> <scheme:!$%25&'()*+,-./0123456789:/@ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz~?#> .",
			subject:   term{text: "http://example/s", kind: IRI},
			predicate: term{text: "http://example/p", kind: IRI},
			object:    term{text: "scheme:!$%25&'()*+,-./0123456789:/@ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz~?#", kind: IRI},
		},
	},

	// Empty valid syntax.
	"nt-syntax-file-01.nq": nil,
	"nt-syntax-file-01.nt": nil,
	"nt-syntax-file-02.nq": nil,
	"nt-syntax-file-02.nt": nil,
	"nt-syntax-file-03.nq": nil,
	"nt-syntax-file-03.nt": nil,

	// Invalid syntax.
	"nq-syntax-bad-literal-01.nq": nil,
	"nq-syntax-bad-literal-02.nq": nil,
	"nq-syntax-bad-literal-03.nq": nil,
	"nq-syntax-bad-quint-01.nq":   nil,
	"nq-syntax-bad-uri-01.nq":     nil,
	"nt-syntax-bad-base-01.nq":    nil,
	"nt-syntax-bad-base-01.nt":    nil,
	"nt-syntax-bad-esc-01.nq":     nil,
	"nt-syntax-bad-esc-01.nt":     nil,
	"nt-syntax-bad-esc-02.nq":     nil,
	"nt-syntax-bad-esc-02.nt":     nil,
	"nt-syntax-bad-esc-03.nq":     nil,
	"nt-syntax-bad-esc-03.nt":     nil,
	"nt-syntax-bad-lang-01.nq":    nil,
	"nt-syntax-bad-lang-01.nt":    nil,
	"nt-syntax-bad-num-01.nq":     nil,
	"nt-syntax-bad-num-01.nt":     nil,
	"nt-syntax-bad-num-02.nq":     nil,
	"nt-syntax-bad-num-02.nt":     nil,
	"nt-syntax-bad-num-03.nq":     nil,
	"nt-syntax-bad-num-03.nt":     nil,
	"nt-syntax-bad-prefix-01.nq":  nil,
	"nt-syntax-bad-prefix-01.nt":  nil,
	"nt-syntax-bad-string-01.nq":  nil,
	"nt-syntax-bad-string-01.nt":  nil,
	"nt-syntax-bad-string-02.nq":  nil,
	"nt-syntax-bad-string-02.nt":  nil,
	"nt-syntax-bad-string-03.nq":  nil,
	"nt-syntax-bad-string-03.nt":  nil,
	"nt-syntax-bad-string-04.nq":  nil,
	"nt-syntax-bad-string-04.nt":  nil,
	"nt-syntax-bad-string-05.nq":  nil,
	"nt-syntax-bad-string-05.nt":  nil,
	"nt-syntax-bad-string-06.nq":  nil,
	"nt-syntax-bad-string-06.nt":  nil,
	"nt-syntax-bad-string-07.nq":  nil,
	"nt-syntax-bad-string-07.nt":  nil,
	"nt-syntax-bad-struct-01.nq":  nil,
	"nt-syntax-bad-struct-01.nt":  nil,
	"nt-syntax-bad-struct-02.nq":  nil,
	"nt-syntax-bad-struct-02.nt":  nil,
	"nt-syntax-bad-uri-01.nq":     nil,
	"nt-syntax-bad-uri-01.nt":     nil,
	"nt-syntax-bad-uri-02.nq":     nil,
	"nt-syntax-bad-uri-02.nt":     nil,
	"nt-syntax-bad-uri-03.nq":     nil,
	"nt-syntax-bad-uri-03.nt":     nil,
	"nt-syntax-bad-uri-04.nq":     nil,
	"nt-syntax-bad-uri-04.nt":     nil,
	"nt-syntax-bad-uri-05.nq":     nil,
	"nt-syntax-bad-uri-05.nt":     nil,
	"nt-syntax-bad-uri-06.nq":     nil,
	"nt-syntax-bad-uri-06.nt":     nil,
	"nt-syntax-bad-uri-07.nq":     nil,
	"nt-syntax-bad-uri-07.nt":     nil,
	"nt-syntax-bad-uri-08.nq":     nil,
	"nt-syntax-bad-uri-08.nt":     nil,
	"nt-syntax-bad-uri-09.nq":     nil,
	"nt-syntax-bad-uri-09.nt":     nil,
}
