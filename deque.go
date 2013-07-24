package discrete

import (
	"errors"
)

/* Simple Stack/Queue implementations, I don't think they need to be explained much */

type Pusher interface {
	Push(x interface{})
}

type Peeker interface {
	Peek() (interface{}, error)
}

type Popper interface {
	Pop() (interface{}, error)
}

type IsEmptyer interface {
	IsEmpty() bool
}

type XInFirstOut interface {
	Pusher
	Popper
	Peeker
	IsEmptyer
}

type Stack []interface{}

type Queue []interface{}

func (s *Stack) Push(x interface{}) {
	*s = append(*s, x)
}

func (s *Stack) Pop() (interface{}, error) {
	if len(*s) == 0 {
		return 0, errors.New("No element to pop")
	}

	x := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	return x, nil
}

func (s *Stack) Peek() (interface{}, error) {
	if len(*s) == 0 {
		return 0, errors.New("No element to peek at")
	}

	return (*s)[len(*s)-1], nil
}

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (q *Queue) Push(x interface{}) {
	*q = append(*q, x)
}

func (q *Queue) Pop() (interface{}, error) {
	if len(*q) == 0 {
		return 0, errors.New("No element to pop")
	}

	x := (*q)[0]
	*q = (*q)[1:]

	return x, nil
}

func (q *Queue) Peek() (interface{}, error) {
	if len(*q) == 0 {
		return 0, errors.New("No element to peek at")
	}

	return (*q)[0], nil
}

func (q *Queue) IsEmpty() bool {
	return len(*q) == 0
}

// Deque is a stack/queue hybrid (from "deck"), I'm not sure if the type conversions will hurt performance or not (I suspect not)
type Deque []interface{}

func (d *Deque) IsEmpty() bool {
	return len(*d) == 0
}

func (d *Deque) Push(x interface{}) {
	*d = append(*d, x)
}

func (d *Deque) Pop() (interface{}, error) {
	a := Stack(*d)
	return (&a).Pop()
}

// Poll is a queue-pop
func (d *Deque) Poll() (interface{}, error) {
	a := Queue(*d)
	return (&a).Pop()
}

func (d *Deque) PeekLast() (interface{}, error) {
	a := Stack(*d)
	return (&a).Peek()
}

func (d *Deque) PeekFirst() (interface{}, error) {
	a := Queue(*d)
	return (&a).Peek()
}
