package common

import (
	"errors"
	"fmt"
)

type LinkedList[V any] interface {
	// Size returns the size of the list.
	//  (1,2,3).Size() => 3
	Size() int
	// Empty returns true if the list is empty.
	//  ().Empty() => true
	Empty() bool
	// AddFirst adds a node at the start of the list with the given element.
	//  (1,2,3).AddFirst(0) => (0,1,2,3)
	AddFirst(value V)
	// AddLast adds a node at the end of the list with the given element.
	//  (1,2,3).AddLast(4) => (1,2,3,4)
	AddLast(value V)
	// Get returns the node's element at the given index.
	//  (1,2,3).Get(2) => 3
	Get(i int) V
	// Set updates a node's element given its index.
	//  (1,2,3).Set(1, 8) => (1,8,3)
	Set(i int, value V) error
	// Insert add a node's element at the given index.
	//  (1,2,3).Insert(1, 8) => (1,8,2,3)
	//  (1,2,3).Insert(3, 8) => (1,2,3,8)
	Insert(i int, value V) error
	// First returns the first nodes' element.
	//  (1,2,3).First() => 1
	First() V
	// Last returns the last node's element.
	//  (1,2,3).Last() => 3
	Last() V
	// RemoveFirst deletes the first node in the list.
	//  (1,2,3).RemoveFirst() => (2,3)
	RemoveFirst() V
	// RemoveLast deletes the last node in the list.
	//  (1,2,3).RemoveLast() => (1,2)
	RemoveLast() V
	// Clear deletes all nodes and empties the list
	//  (1,2,3).Clear() => ()
	Clear()
	// Iter is an iterator to be used for iterate the linkedlist (front first) easily.
	//  for x := range list.Iter() { ... }
	Iter() chan V
	// ToSlice returns a slice representation of the LinkedList.
	ToSlice() []V
	// Append will add the other items to the current list
	//  (1,2,3).Append((4,5,6)) => (1,2,3,4,5,6)
	Append(other LinkedList[V])
}

var (
	ErrNotFound    = errors.New("item not found in list")
	ErrOutOfBounds = errors.New("index out of bounds")
)

// A linkedList has a size, a pointer to the first node, and
// a pointer to the last node in the list.
//
//	first -> 1
//	         2
//	         3
//	 last -> 4
type linkedList[V any] struct {
	size    int
	first   *node[V]
	last    *node[V]
	compare func(V, V) int
}

// The linkedList's chain is made up of nodes with an element,
// a pointer to the previous node, and a pointer to the next node.
//
//	1<->2<->3<->4
type node[V any] struct {
	Value V
	next  *node[V]
	prev  *node[V]
}

// New creates an unsynchronized linked list
//
//	mylist := list.New()
func NewLinkedList[V any](values ...V) *linkedList[V] {
	l := &linkedList[V]{}
	for _, v := range values {
		l.AddLast(v)
	}
	return l
}

func (l *linkedList[V]) Size() int {
	return l.size
}

func (l *linkedList[V]) Empty() bool {
	return l.Size() == 0
}

func (l *linkedList[V]) AddFirst(value V) {
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

func (l *linkedList[V]) AddLast(value V) {
	n := &node[V]{Value: value, next: nil, prev: l.last}

	if l.size == 0 {
		l.first = n
	} else {
		l.last.next = n
	}

	l.last = n
	l.size += 1
}

func (l *linkedList[V]) Get(i int) V {
	node, err := l.getNode(i)
	if err != nil {
		var result V
		return result
	}
	return node.Value
}

func (l *linkedList[V]) Set(i int, value V) error {
	node, err := l.getNode(i)
	if err == nil {
		node.Value = value
	}
	return err
}

func (l *linkedList[V]) Insert(i int, value V) error {
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

func (l *linkedList[V]) insert(cur *node[V], value V) error {
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

func (l *linkedList[V]) First() V {
	if l.size == 0 {
		var result V
		return result
	}
	return l.first.Value
}

func (l *linkedList[V]) Last() V {
	if l.size == 0 {
		var result V
		return result
	}
	return l.last.Value
}

func (l *linkedList[V]) RemoveFirst() V {
	if l.size == 0 {
		var result V
		return result
	}
	result := l.first
	l.removeNode(l.first)
	return result.Value
}

func (l *linkedList[V]) RemoveLast() V {
	if l.size == 0 {
		var result V
		return result
	}
	result := l.last
	l.removeNode(l.last)
	return result.Value
}

func (l *linkedList[V]) Clear() {
	l.size = 0
	l.first = nil
	l.last = nil
}

func (l *linkedList[V]) Iter() chan V {
	return l.iter()
}

func (l *linkedList[V]) ToSlice() []V {
	res := make([]V, l.size)
	i := 0
	for x := range l.iter() {
		res[i] = x
		i++
	}
	return res
}

func (l *linkedList[V]) Append(other LinkedList[V]) {
	ol := other.(*linkedList[V])
	if l.size == 0 {
		l.first = ol.first
		l.last = ol.last
		l.size = ol.size
		return
	}
	if ol.size == 0 {
		return
	}

	l.last.next = ol.first
	ol.first.prev = l.last
	l.last = ol.last
	l.size += ol.size
}

// String will return a string via fmt.Sprintf() renders a slice.
func (l *linkedList[V]) String() string {
	return fmt.Sprintf("%v", l.ToSlice())
}

// iter is used internally and is not locked.
func (l *linkedList[V]) iter() chan V {
	ch := make(chan V, l.size)
	for n := l.first; n != nil; n = n.next {
		ch <- n.Value
	}
	close(ch)
	return ch
}

// getNode retrieves a node given an index.
func (l *linkedList[V]) getNode(i int) (*node[V], error) {
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
func (l *linkedList[V]) removeNode(node *node[V]) {
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
