// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encoding

import (
	"sort"
	"testing"
)

var setAttributesTests = []struct {
	attr   *Attributes
	opName string
	op     func(AttributeSetter) error
	want   *Attributes
}{
	{
		attr:   &Attributes{},
		opName: "noop",
		op: func(a AttributeSetter) error {
			return a.SetAttribute(Attribute{Key: "", Value: "bar"})
		},
		want: &Attributes{},
	},
	{
		attr:   &Attributes{},
		opName: "add attr to empty",
		op: func(a AttributeSetter) error {
			return a.SetAttribute(Attribute{Key: "foo", Value: "bar"})
		},
		want: &Attributes{{Key: "foo", Value: "bar"}},
	},
	{
		attr:   &Attributes{},
		opName: "remove attr from empty",
		op: func(a AttributeSetter) error {
			return a.SetAttribute(Attribute{Key: "foo", Value: ""})
		},
		want: &Attributes{},
	},
	{
		attr:   &Attributes{{Key: "foo", Value: "bar"}},
		opName: "add attr to non-empty",
		op: func(a AttributeSetter) error {
			return a.SetAttribute(Attribute{Key: "bif", Value: "fud"})
		},
		want: &Attributes{{Key: "foo", Value: "bar"}, {Key: "bif", Value: "fud"}},
	},
	{
		attr:   &Attributes{{Key: "foo", Value: "bar"}},
		opName: "remove attr from singleton",
		op: func(a AttributeSetter) error {
			return a.SetAttribute(Attribute{Key: "foo", Value: ""})
		},
		want: &Attributes{},
	},
	{
		attr:   &Attributes{{Key: "foo", Value: "bar"}, {Key: "bif", Value: "fud"}},
		opName: "remove first attr from pair",
		op: func(a AttributeSetter) error {
			return a.SetAttribute(Attribute{Key: "foo", Value: ""})
		},
		want: &Attributes{{Key: "bif", Value: "fud"}},
	},
	{
		attr:   &Attributes{{Key: "foo", Value: "bar"}, {Key: "bif", Value: "fud"}},
		opName: "remove second attr from pair",
		op: func(a AttributeSetter) error {
			return a.SetAttribute(Attribute{Key: "bif", Value: ""})
		},
		want: &Attributes{{Key: "foo", Value: "bar"}},
	},
	{
		attr:   &Attributes{{Key: "foo", Value: "bar"}, {Key: "bif", Value: "fud"}},
		opName: "replace first attr in pair",
		op: func(a AttributeSetter) error {
			return a.SetAttribute(Attribute{Key: "foo", Value: "not bar"})
		},
		want: &Attributes{{Key: "foo", Value: "not bar"}, {Key: "bif", Value: "fud"}},
	},
	{
		attr:   &Attributes{{Key: "foo", Value: "bar"}, {Key: "bif", Value: "fud"}},
		opName: "replace second attr in pair",
		op: func(a AttributeSetter) error {
			return a.SetAttribute(Attribute{Key: "bif", Value: "not fud"})
		},
		want: &Attributes{{Key: "foo", Value: "bar"}, {Key: "bif", Value: "not fud"}},
	},
}

func TestSetAttributes(t *testing.T) {
	for _, test := range setAttributesTests {
		err := test.op(test.attr)
		if err != nil {
			t.Errorf("unexpected error for %q: %v", test.opName, err)
			continue
		}
		if !sameAttributes(test.attr, test.want) {
			t.Errorf("unexpected result from %q:\ngot: %+v\nwant:%+v", test.opName, test.attr, test.want)
		}
	}
}

func sameAttributes(a, b Attributer) bool {
	aAttr := a.Attributes()
	bAttr := b.Attributes()
	if len(aAttr) != len(bAttr) {
		return false
	}
	aAttr = append(aAttr[:0:0], aAttr...)
	sort.Slice(aAttr, func(i, j int) bool { return aAttr[i].Key < aAttr[j].Key })
	bAttr = append(bAttr[:0:0], bAttr...)
	sort.Slice(bAttr, func(i, j int) bool { return bAttr[i].Key < bAttr[j].Key })
	for i, a := range aAttr {
		if bAttr[i] != a {
			return false
		}
	}
	return true
}
