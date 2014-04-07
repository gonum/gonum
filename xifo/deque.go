// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xifo

import ()

/* Simple Stack/Queue implementations, I don't think they need to be explained much */

type Stack interface {
	Push(x interface{})
	Pop() interface{}
	PeekLast() interface{}
	IsEmpty() bool
}

type Queue interface {
	Push(x interface{})
	Poll() interface{}
	PeekFirst() interface{}
	IsEmpty() bool
}

type Deque interface {
	Push(x interface{})
	Poll() interface{}
	Pop() interface{}
	PeekFirst() interface{}
	PeekLast() interface{}
	IsEmpty() bool
}

type GonumStack []interface{}

type GonumQueue []interface{}

func (s *GonumStack) Push(x interface{}) {
	*s = append(*s, x)
}

func (s *GonumStack) Pop() interface{} {
	if len(*s) == 0 {
		panic("No element to pop")
	}

	x := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	return x
}

func (s *GonumStack) PeekLast() interface{} {
	if len(*s) == 0 {
		panic("No element to peek at")
	}

	return (*s)[len(*s)-1]
}

func (s *GonumStack) IsEmpty() bool {
	return len(*s) == 0
}

func (q *GonumQueue) Push(x interface{}) {
	*q = append(*q, x)
}

func (q *GonumQueue) Poll() interface{} {
	if len(*q) == 0 {
		panic("No element to poll")
	}

	x := (*q)[0]
	*q = (*q)[1:]

	return x
}

func (q *GonumQueue) PeekFirst() interface{} {
	if len(*q) == 0 {
		panic("No element to peek at")
	}

	return (*q)[0]
}

func (q *GonumQueue) IsEmpty() bool {
	return len(*q) == 0
}

// Deque is a stack/queue hybrid (from "deck"), I'm not sure if the type conversions will hurt
// performance or not (I suspect not.)
type GonumDeque []interface{}

func (d *GonumDeque) IsEmpty() bool {
	return len(*d) == 0
}

func (d *GonumDeque) Push(x interface{}) {
	*d = append(*d, x)
}

func (d *GonumDeque) Pop() interface{} {
	a := (*GonumStack)(d)
	return a.Pop()
}

// Poll is a queue-pop.
func (d *GonumDeque) Poll() interface{} {
	a := (*GonumQueue)(d)
	return a.Poll()
}

func (d *GonumDeque) PeekLast() interface{} {
	a := (*GonumStack)(d)
	return a.PeekLast()
}

func (d *GonumDeque) PeekFirst() interface{} {
	a := (*GonumQueue)(d)
	return a.PeekFirst()
}
