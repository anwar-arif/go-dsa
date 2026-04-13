// Package stack provides a generic LIFO stack backed by a slice.
//
// All operations are O(1) amortized (Push/Pop/Peek/Size/IsEmpty).
//
// Basic usage:
//
//	s := stack.New[int]()
//	s.Push(1)
//	s.Push(2)
//	v, _ := s.Pop()  // v == 2
//	v, _ = s.Peek()  // v == 1, stack unchanged
package stack

// Stack is a generic LIFO stack.
// The zero value is ready to use.
type Stack[T any] struct {
	items []T
}

// New returns an empty Stack.
func New[T any]() *Stack[T] {
	return &Stack[T]{}
}

// NewFromSlice returns a Stack pre-populated with values.
// The first element of the slice becomes the bottom of the stack;
// the last element becomes the top.
// The slice is copied; the original is not modified.
func NewFromSlice[T any](values []T) *Stack[T] {
	items := make([]T, len(values))
	copy(items, values)
	return &Stack[T]{items: items}
}

// Push adds v to the top of the stack. Amortized O(1).
func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

// Pop removes and returns the top element.
// Returns the zero value and false if the stack is empty. O(1).
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	n := len(s.items)
	v := s.items[n-1]
	var zero T
	s.items[n-1] = zero // prevent memory leak for pointer/interface types
	s.items = s.items[:n-1]
	return v, true
}

// Peek returns the top element without removing it.
// Returns the zero value and false if the stack is empty. O(1).
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// Size returns the number of elements in the stack. O(1).
func (s *Stack[T]) Size() int {
	return len(s.items)
}

// IsEmpty reports whether the stack has no elements. O(1).
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Clear removes all elements from the stack.
func (s *Stack[T]) Clear() {
	var zero T
	for i := range s.items {
		s.items[i] = zero
	}
	s.items = s.items[:0]
}

// ToSlice returns a copy of the stack contents ordered from bottom to top.
// The last element of the returned slice is the current top.
func (s *Stack[T]) ToSlice() []T {
	out := make([]T, len(s.items))
	copy(out, s.items)
	return out
}
