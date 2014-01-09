package xifo_test

import (
	"github.com/gonum/graph/xifo"
	"testing"
)

func TestStack(t *testing.T) {
	testStack(&xifo.GonumStack{}, t)
}

func testStack(stack xifo.Stack, t *testing.T) {

	stack.Push(5)
	if stack.IsEmpty() {
		t.Error("Stack is empty after push")
	}

	if i, ok := stack.PeekLast().(int); !ok {
		t.Error("Stack peek does not return element of correct type")
	} else if i != 5 {
		t.Error("Stack peek does not return element of correct value")
	}

	if stack.IsEmpty() {
		t.Fatal("Stack peek destroys element")
	}

	if i, ok := stack.Pop().(int); !ok {
		t.Error("Stack pop does not return element of correct type")
	} else if i != 5 {
		t.Error("Stack pop does not return element of correct value")
	}

	if !stack.IsEmpty() {
		t.Fatal("Stack is not empty after popping last element")
	}

	for i := 0; i <= 10; i++ {
		stack.Push(i)
	}

	for i := 10; i >= 0; i-- {
		j := stack.Pop().(int)
		if j == (10-i) && i != 5 {
			t.Errorf("Stack is FIFO; stack value: %d; expected: %d", j, i)
		} else if j != i {
			t.Errorf("Stack returned strange element on LIFO test; stack value: %d; expected: %d", j, i)
		}
	}

	if !stack.IsEmpty() {
		t.Fatal("Stack is not empty after popping last element")
	}

	defer func() {
		err := recover()
		if err == nil {
			t.Error("Stack did not properly panic when popping an empty one")
		}
	}()

	stack.Pop()
}

func TestQueue(t *testing.T) {
	testQueue(&xifo.GonumQueue{}, t)
}

func testQueue(queue xifo.Queue, t *testing.T) {

	queue.Push(5)
	if queue.IsEmpty() {
		t.Error("Queue is empty after push")
	}

	if i, ok := queue.PeekFirst().(int); !ok {
		t.Error("Queue peek does not return element of correct type")
	} else if i != 5 {
		t.Error("Queue peek does not return element of correct value")
	}

	if queue.IsEmpty() {
		t.Fatal("Queue peek destroys element")
	}

	if i, ok := queue.Poll().(int); !ok {
		t.Error("Queue poll does not return element of correct type (or peek mutates it)")
	} else if i != 5 {
		t.Error("Queue poll does not return element of correct value (or peek mutates it)")
	}

	if !queue.IsEmpty() {
		t.Fatal("Queue is not empty after polling last element")
	}

	for i := 0; i <= 10; i++ {
		queue.Push(i)
	}

	for i := 0; i <= 10; i++ {
		j := queue.Poll().(int)
		if j == (10-i) && i != 5 {
			t.Errorf("Queue is LIFO; queue value: %d; expected: %d", j, i)
		} else if j != i {
			t.Errorf("Queue returned strange element on FIFO test; queue value: %d; expected: %d", j, i)
		}
	}

	if !queue.IsEmpty() {
		t.Fatal("Queue is not empty after popping last element")
	}

	defer func() {
		err := recover()
		if err == nil {
			t.Error("Queue did not properly panic when polling an empty one")
		}
	}()

	queue.Poll()
}

func TestDeque(t *testing.T) {
	testStack(&xifo.GonumDeque{}, t)
	testQueue(&xifo.GonumDeque{}, t)
}
