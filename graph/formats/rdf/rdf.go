// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate ragel -Z -G2 parse.rl
//go:generate ragel -Z -G2 extract.rl
//go:generate ragel -Z -G2 check.rl
//go:generate stringer -type=Kind

package rdf

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"gonum.org/v1/gonum/graph"
)

var (
	_ graph.Node = Term{}
	_ graph.Edge = (*Statement)(nil)
	_ graph.Line = (*Statement)(nil)
)

var (
	ErrInvalid        = errors.New("invalid N-Quad")
	ErrIncomplete     = errors.New("incomplete N-Quad")
	ErrInvalidTerm    = errors.New("invalid term")
	ErrIncompleteTerm = errors.New("incomplete term")
)

// Kind represents the kind of an RDF term.
type Kind int

const (
	// Invalid is an invalid RDF term.
	Invalid Kind = iota

	// IRI is the kind of an IRI term.
	// https://www.w3.org/TR/n-quads/#sec-iri
	IRI

	// Literal is the kind of an RDF literal.
	// https://www.w3.org/TR/n-quads/#sec-literals
	Literal

	// Blank is the kind of an RDF blank node term.
	// https://www.w3.org/TR/n-quads/#BNodes
	Blank
)

// Term is an RDF term. It implements the graph.Node interface.
type Term struct {
	// Value is the text value of term.
	Value string

	// UID is the unique ID for the term
	// in a collection of RDF terms.
	UID int64
}

// NewBlankTerm returns a Term based on the provided RDF blank node
// label. The label should not include the "_:" prefix. The returned
// Term will not have the UID set.
func NewBlankTerm(label string) (Term, error) {
	err := checkLabelText([]rune(label))
	if err != nil {
		return Term{}, err
	}
	return Term{Value: blankPrefix + label}, nil
}

const blankPrefix = "_:"

func isBlank(s string) bool {
	return strings.HasPrefix(s, blankPrefix)
}

// NewIRITerm returns a Term based on the provided IRI which must
// be valid and include a scheme. The returned Term will not have
// the UID set.
func NewIRITerm(iri string) (Term, error) {
	err := checkIRIText(iri)
	if err != nil {
		return Term{}, err
	}
	return Term{Value: escape("<", iri, ">")}, nil
}

func isIRI(s string) bool {
	return strings.HasPrefix(s, "<") && strings.HasSuffix(s, ">")
}

// NewLiteralTerm returns a Term based on the literal text and an
// optional qualifier which may either be a "@"-prefixed language
// tag or a valid IRI. The text will be escaped if necessary and quoted,
// and if an IRI is given it will be escaped if necessary. The returned
// Term will not have the UID set.
func NewLiteralTerm(text, qual string) (Term, error) {
	text = escape(`"`, text, `"`)
	if qual == "" {
		return Term{Value: text}, nil
	}
	if strings.HasPrefix(qual, "@") {
		err := checkLangText([]byte(qual))
		if err != nil {
			return Term{}, err
		}
		return Term{Value: text + qual}, nil
	}
	err := checkIRIText(qual)
	if err != nil {
		return Term{}, err
	}
	return Term{Value: text + escape("^^<", qual, ">")}, nil
}

func checkIRIText(iri string) error {
	switch u, err := url.Parse(iri); {
	case err != nil:
		return err
	case u.Scheme == "":
		return fmt.Errorf("rdf: %w: relative IRI ref %q", ErrInvalidTerm, iri)
	default:
		return nil
	}
}

func isLiteral(s string) bool {
	return strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`)
}

// Parts returns the parts of the term and the kind of the term.
// IRI node text is returned as a valid IRI with the quoting angle
// brackets removed and escape sequences interpreted, and blank
// nodes are stripped of the "_:" prefix.
// When the term is a literal, qual will either be empty, an unescaped
// IRI, or an RDF language tag prefixed with an @ symbol. The literal
// text is returned unquoted and unescaped.
func (t Term) Parts() (text, qual string, kind Kind, err error) {
	return extract([]rune(t.Value))
}

// ID returns the value of the Term's UID field.
func (t Term) ID() int64 { return t.UID }

// Statement is an RDF statement. It implements the graph.Edge and graph.Line
// interfaces.
type Statement struct {
	Subject   Term
	Predicate Term
	Object    Term
	Label     Term
}

// String returns the RDF 1.1 N-Quad formatted statement.
func (s *Statement) String() string {
	if s.Label.Value == "" {
		return fmt.Sprintf("%s %s %s .", s.Subject.Value, s.Predicate.Value, s.Object.Value)
	}
	return fmt.Sprintf("%s %s %s %s .", s.Subject.Value, s.Predicate.Value, s.Object.Value, s.Label.Value)
}

// From returns the subject of the statement.
func (s *Statement) From() graph.Node { return s.Subject }

// To returns the object of the statement.
func (s *Statement) To() graph.Node { return s.Object }

// ID returns the UID of the Predicate field.
func (s *Statement) ID() int64 { return s.Predicate.UID }

// ReversedEdge returns the receiver unaltered. If there is a semantically
// valid edge reversal operation for the data, the user should implement
// this by wrapping Statement in a type performing that operation.
// See the ReversedLine example for details.
func (s *Statement) ReversedEdge() graph.Edge { return s }

// ReversedLine returns the receiver unaltered. If there is a semantically
// valid line reversal operation for the data, the user should implement
// this by wrapping Statement in a type performing that operation.
func (s *Statement) ReversedLine() graph.Line { return s }

// ParseNQuad parses the statement and returns the corresponding Statement.
// All Term UID fields are zero on return.
func ParseNQuad(statement string) (*Statement, error) {
	s, err := parse([]rune(statement))
	if err != nil {
		return nil, err
	}
	return &s, err
}

// Decoder is an RDF stream decoder. Statements returned by calls to the
// Unmarshal method have their Terms' UID fields set so that unique terms
// will have unique IDs and so can be used directly in a graph.Multi, or
// in a graph.Graph if all predicate terms are identical. IDs created by
// the decoder all exist within a single namespace and so Terms can be
// uniquely identified by their UID. Term UIDs are based from 1 to allow
// RDF-aware client graphs to assign ID if no ID has been assigned.
type Decoder struct {
	scanner *bufio.Scanner

	strings store
	ids     map[string]int64
}

// NewDecoder returns a new Decoder that takes input from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		scanner: bufio.NewScanner(r),
		strings: make(store),
		ids:     make(map[string]int64),
	}
}

// Reset resets the decoder to use the provided io.Reader, retaining
// the existing Term ID mapping.
func (dec *Decoder) Reset(r io.Reader) {
	dec.scanner = bufio.NewScanner(r)
	dec.strings = make(store)
	if dec.ids == nil {
		dec.ids = make(map[string]int64)
	}
}

// Unmarshal returns the next statement from the input stream.
func (dec *Decoder) Unmarshal() (*Statement, error) {
	for dec.scanner.Scan() {
		data := bytes.TrimSpace(dec.scanner.Bytes())
		if len(data) == 0 || data[0] == '#' {
			continue
		}

		s, err := ParseNQuad(string(data))
		if err != nil {
			return nil, fmt.Errorf("rdf: failed to parse %q: %w", data, err)
		}
		if s == nil {
			continue
		}

		s.Subject.Value = dec.strings.intern(s.Subject.Value)
		s.Predicate.Value = dec.strings.intern(s.Predicate.Value)
		s.Object.Value = dec.strings.intern(s.Object.Value)
		s.Subject.UID = dec.idFor(s.Subject.Value)
		s.Object.UID = dec.idFor(s.Object.Value)
		s.Predicate.UID = dec.idFor(s.Predicate.Value)
		if s.Label.Value != "" {
			s.Label.Value = dec.strings.intern(s.Label.Value)
			s.Label.UID = dec.idFor(s.Label.Value)
		}
		return s, nil
	}
	dec.strings = nil
	err := dec.scanner.Err()
	if err != nil {
		return nil, err
	}
	return nil, io.EOF
}

func (dec *Decoder) idFor(s string) int64 {
	id, ok := dec.ids[s]
	if ok {
		return id
	}
	id = int64(len(dec.ids)) + 1
	dec.ids[s] = id
	return id
}

// Terms returns the mapping between terms and graph node IDs constructed
// during decoding the RDF statement stream.
func (dec *Decoder) Terms() map[string]int64 {
	return dec.ids
}

// store is a string internment implementation.
type store map[string]string

// intern returns an interned version of the parameter.
func (is store) intern(s string) string {
	if s == "" {
		return ""
	}

	if len(s) < 2 || len(s) > 512 {
		// Not enough benefit on average with real data.
		return s
	}

	t, ok := is[s]
	if ok {
		return t
	}
	is[s] = s
	return s
}

func escape(lq, s, rq string) string {
	var buf strings.Builder
	if lq != "" {
		buf.WriteString(lq)
	}
	for _, r := range s {
		var c byte
		switch r {
		case '\n':
			c = 'n'
		case '\r':
			c = 'r'
		case '"', '\\':
			c = byte(r)
		default:
			const hex = "0123456789abcdef"
			switch {
			case r <= unicode.MaxASCII || strconv.IsPrint(r):
				buf.WriteRune(r)
			case r > utf8.MaxRune:
				r = 0xFFFD
				fallthrough
			case r < 0x10000:
				buf.WriteString("\\u")
				for s := 12; s >= 0; s -= 4 {
					buf.WriteByte(hex[r>>uint(s)&0xf])
				}
			default:
				buf.WriteString("\\U")
				for s := 28; s >= 0; s -= 4 {
					buf.WriteByte(hex[r>>uint(s)&0xf])
				}
			}
			continue
		}
		buf.Write([]byte{'\\', c})
	}
	if rq != "" {
		buf.WriteString(rq)
	}
	return buf.String()
}

func unEscape(r []rune) string {
	var buf strings.Builder
	for i := 0; i < len(r); {
		switch r[i] {
		case '\\':
			i++
			var c byte
			switch r[i] {
			case 't':
				c = '\t'
			case 'b':
				c = '\b'
			case 'n':
				c = '\n'
			case 'r':
				c = '\r'
			case 'f':
				c = '\f'
			case '"':
				c = '"'
			case '\\':
				c = '\\'
			case '\'':
				c = '\''
			case 'u':
				rc, err := strconv.ParseInt(string(r[i+1:i+5]), 16, 32)
				if err != nil {
					panic(fmt.Errorf("internal parser error: %w", err))
				}
				buf.WriteRune(rune(rc))
				i += 5
				continue
			case 'U':
				rc, err := strconv.ParseInt(string(r[i+1:i+9]), 16, 32)
				if err != nil {
					panic(fmt.Errorf("internal parser error: %w", err))
				}
				buf.WriteRune(rune(rc))
				i += 9
				continue
			}
			buf.WriteByte(c)
		default:
			buf.WriteRune(r[i])
		}
		i++
	}

	return buf.String()
}
