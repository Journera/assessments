package common

import (
	"errors"
	"fmt"
)

var (
	ErrOutOfBounds = errors.New("index out of bounds")
)

// A LinkedList has a size, a pointer to the first node, and
// a pointer to the last node in the list.
//
//	first -> 1
//	         2
//	         3
//	 last -> 4
type LinkedList[V any] struct {
	size    int
	first   *node[V]
	last    *node[V]
	compare func(V, V) int
}

// The LinkedList's chain is made up of nodes with an element,
// a pointer to the previous node, and a pointer to the next node.
//
//	1<->2<->3<->4
type node[V any] struct {
	Value V
	next  *node[V]
	prev  *node[V]
}

func NewLinkedList[V any](values ...V) *LinkedList[V] {
	l := &LinkedList[V]{}
	for _, v := range values {
		l.AddLast(v)
	}
	return l
}

func (l *LinkedList[V]) Size() int {
	return l.size
}

func (l *LinkedList[V]) Empty() bool {
	return l.Size() == 0
}

func (l *LinkedList[V]) AddFirst(value V) {
	n := &node[V]{Value: value, next: nil, prev: nil}

	if l.size == 0 {
		l.last = n
	} else {
		n.next = l.first
		l.first.prev = n
	}

	l.first = n
	l.size += 1
}

func (l *LinkedList[V]) AddLast(value V) {
	n := &node[V]{Value: value, next: nil, prev: l.last}

	if l.size == 0 {
		l.first = n
	} else {
		l.last.next = n
	}

	l.last = n
	l.size += 1
}

func (l *LinkedList[V]) Get(i int) V {
	node, err := l.getNode(i)
	if err != nil {
		var result V
		return result
	}
	return node.Value
}

func (l *LinkedList[V]) Set(i int, value V) error {
	node, err := l.getNode(i)
	if err == nil {
		node.Value = value
	}
	return err
}

func (l *LinkedList[V]) Insert(i int, value V) error {
	if i == 0 {
		l.AddFirst(value)
		return nil
	}
	if i == l.size {
		l.AddLast(value)
		return nil
	}

	cur, err := l.getNode(i)
	if err != nil {
		return err
	}
	l.insert(cur, value)
	return nil
}

func (l *LinkedList[V]) insert(cur *node[V], value V) error {
	prev := cur.prev
	n := &node[V]{Value: value, next: cur, prev: prev}
	cur.prev = n
	if prev == nil {
		l.first = n
	} else {
		prev.next = n
	}
	l.size += 1
	return nil
}

func (l *LinkedList[V]) First() V {
	if l.size == 0 {
		var result V
		return result
	}
	return l.first.Value
}

func (l *LinkedList[V]) Last() V {
	if l.size == 0 {
		var result V
		return result
	}
	return l.last.Value
}

func (l *LinkedList[V]) RemoveFirst() V {
	if l.size == 0 {
		var result V
		return result
	}
	result := l.first
	l.removeNode(l.first)
	return result.Value
}

func (l *LinkedList[V]) RemoveLast() V {
	if l.size == 0 {
		var result V
		return result
	}
	result := l.last
	l.removeNode(l.last)
	return result.Value
}

func (l *LinkedList[V]) Clear() {
	l.size = 0
	l.first = nil
	l.last = nil
}

func (l *LinkedList[V]) Iter() chan V {
	return l.iter()
}

func (l *LinkedList[V]) ToSlice() []V {
	res := make([]V, l.size)
	i := 0
	for x := range l.iter() {
		res[i] = x
		i++
	}
	return res
}

func (l *LinkedList[V]) Append(other *LinkedList[V]) {
	if l.size == 0 {
		l.first = other.first
		l.last = other.last
		l.size = other.size
		return
	}
	if other.size == 0 {
		return
	}

	l.last.next = other.first
	other.first.prev = l.last
	l.last = other.last
	l.size += other.size
}

// String will return a string via fmt.Sprintf() renders a slice.
func (l *LinkedList[V]) String() string {
	return fmt.Sprintf("%v", l.ToSlice())
}

func (l *LinkedList[V]) Sort(comparator func(i, j V) int) {
	if l.size == 0 {
		return
	}
	current := l.first
	for current != nil {
		index := current.next
		for index != nil {
			if comparator(current.Value, index.Value) > 0 {
				current.Value, index.Value = index.Value, current.Value
			}
			index = index.next
		}
		current = current.next
	}
}

// iter is used internally and is not locked.
func (l *LinkedList[V]) iter() chan V {
	ch := make(chan V, l.size)
	for n := l.first; n != nil; n = n.next {
		ch <- n.Value
	}
	close(ch)
	return ch
}

// getNode retrieves a node given an index.
func (l *LinkedList[V]) getNode(i int) (*node[V], error) {
	if l.size == 0 || i > l.size-1 {
		return nil, ErrOutOfBounds
	}

	var n *node[V]

	if i <= l.size/2 {
		n = l.first
		for p := 0; p != i; p++ {
			n = n.next
		}
	} else {
		n = l.last
		for p := l.size - 1; p != i; p-- {
			n = n.prev
		}
	}

	return n, nil
}

// removeNode deletes the node from the given list.
// The function is considered to be used internally.
func (l *LinkedList[V]) removeNode(node *node[V]) {
	if l.size == 1 { // Only node
		l.first = nil
		l.last = nil
		l.size--
		return
	}

	if node.prev == nil { // First node
		node.next.prev = nil
		l.first = node.next
		l.size--
		return
	}

	if node.next == nil { // Last node
		node.prev.next = nil
		l.last = node.prev
		l.size--
		return
	}

	// Node in middle of chain
	node.next.prev = node.prev
	node.prev.next = node.next
	l.size--
	return
}
