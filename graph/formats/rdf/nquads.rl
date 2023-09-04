// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Ragel grammar definition derived from http://www.w3.org/TR/n-quads/#sec-grammar.

%%{
	machine nquads;

	alphtype rune;

	PN_CHARS_BASE        = [A-Za-z]
	                     | 0x00c0 .. 0x00d6
	                     | 0x00d8 .. 0x00f6
	                     | 0x00f8 .. 0x02ff
	                     | 0x0370 .. 0x037d
	                     | 0x037f .. 0x1fff
	                     | 0x200c .. 0x200d
	                     | 0x2070 .. 0x218f
	                     | 0x2c00 .. 0x2fef
	                     | 0x3001 .. 0xd7ff
	                     | 0xf900 .. 0xfdcf
	                     | 0xfdf0 .. 0xfffd
	                     | 0x10000 .. 0xeffff
	                     ;

	PN_CHARS_U           = PN_CHARS_BASE | '_' | ':' ;

	PN_CHARS             = PN_CHARS_U
	                     | '-'
	                     | [0-9]
	                     | 0xb7
	                     | 0x0300 .. 0x036f
	                     | 0x203f .. 0x2040
	                     ;

	BLANK_NODE_LABEL     = (PN_CHARS_U | [0-9]) ((PN_CHARS | '.')* PN_CHARS)? ;

	BLANK_NODE           = '_:' BLANK_NODE_LABEL ;

	ECHAR                = ('\\' [tbnrf"'\\]) ;

	UCHAR                = ('\\u' xdigit {4}
	                     | '\\U' xdigit {8})
	                     ;

	STRING_LITERAL       = (
	                       0x00 .. 0x09
	                     | 0x0b .. 0x0c
	                     | 0x0e .. '!'
	                     | '#' .. '['
	                     | ']' .. 0x10ffff
	                     | ECHAR
	                     | UCHAR)*
	                     ;

	STRING_LITERAL_QUOTE = '"' STRING_LITERAL '"' ;

	IRI                  = (
	                       '!' .. ';'
	                     | '='
	                     | '?' .. '['
	                     | ']'
	                     | '_'
	                     | 'a' .. 'z'
	                     | '~'
	                     | 0x80 .. 0x10ffff
	                     | UCHAR)*
	                     ;

	IRIREF               = '<' IRI >StartIRI %EndIRI '>' ;

	LANGTAG              = '@' [a-zA-Z]+ ('-' [a-zA-Z0-9]+)* ;

	whitespace           = [ \t] ;

	literal              = STRING_LITERAL_QUOTE ('^^' IRIREF | LANGTAG)? ;

	subject              = IRIREF | BLANK_NODE ;
	predicate            = IRIREF ;
	object               = IRIREF | BLANK_NODE | literal ;
	graphLabel           = IRIREF | BLANK_NODE ;
}%%
