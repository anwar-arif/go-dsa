// Package collections provides C++-like generic data structures for Go.
// Requires Go 1.21+
package collections

import (
	"errors"
	"fmt"
	"strings"
)

// ErrStackEmpty is returned when an operation requires a non-empty stack.
var ErrStackEmpty = errors.New("stack is empty")

// Stack is a generic LIFO (Last-In, First-Out) data structure.
// It supports any type T, including int, string, and custom structs.
//
// Example:
//
//	s := collections.New[int]()
//	s.Push(1)
//	s.Push(2)
//	val, _ := s.Pop() // returns 2
type Stack[T any] struct {
	items []T
}

// New creates and returns a new empty Stack for type T.
//
//	intStack := collections.New[int]()
//	strStack := collections.New[string]()
//	myStack := collections.New[MyStruct]()
func New[T any]() *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0),
	}
}

// NewWithCapacity creates a Stack with a pre-allocated capacity hint.
// Useful when you know the approximate maximum size upfront.
func NewWithCapacity[T any](capacity int) *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0, capacity),
	}
}

// NewFromSlice creates a Stack pre-populated from a slice.
// The last element of the slice becomes the top of the stack.
//
//	s := collections.NewFromSlice([]int{1, 2, 3})
//	s.Peek() // returns 3
func NewFromSlice[T any](items []T) *Stack[T] {
	copied := make([]T, len(items))
	copy(copied, items)
	return &Stack[T]{items: copied}
}

// Push adds an element onto the top of the stack.
// Time complexity: O(1) amortized.
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top element of the stack.
// Returns ErrStackEmpty if the stack is empty.
// Time complexity: O(1).
func (s *Stack[T]) Pop() (T, error) {
	if s.IsEmpty() {
		var zero T
		return zero, ErrStackEmpty
	}
	top := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return top, nil
}

// MustPop removes and returns the top element, panicking if the stack is empty.
// Use when you are certain the stack is non-empty.
func (s *Stack[T]) MustPop() T {
	val, err := s.Pop()
	if err != nil {
		panic(err)
	}
	return val
}

// Peek returns the top element without removing it.
// Returns ErrStackEmpty if the stack is empty.
// Time complexity: O(1).
func (s *Stack[T]) Peek() (T, error) {
	if s.IsEmpty() {
		var zero T
		return zero, ErrStackEmpty
	}
	return s.items[len(s.items)-1], nil
}

// MustPeek returns the top element without removing it, panicking if empty.
func (s *Stack[T]) MustPeek() T {
	val, err := s.Peek()
	if err != nil {
		panic(err)
	}
	return val
}

// IsEmpty returns true if the stack has no elements.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Size returns the number of elements in the stack.
func (s *Stack[T]) Size() int {
	return len(s.items)
}

// Clear removes all elements from the stack.
func (s *Stack[T]) Clear() {
	s.items = s.items[:0]
}

// ToSlice returns a copy of the stack's elements as a slice,
// ordered from bottom to top (index 0 = bottom, last = top).
func (s *Stack[T]) ToSlice() []T {
	result := make([]T, len(s.items))
	copy(result, s.items)
	return result
}

// Clone returns a deep copy of the stack.
func (s *Stack[T]) Clone() *Stack[T] {
	return NewFromSlice(s.items)
}

// Contains checks if the stack holds the given value using the provided equality function.
//
//	s.Contains(func(v int) bool { return v == 42 })
func (s *Stack[T]) Contains(predicate func(T) bool) bool {
	for _, item := range s.items {
		if predicate(item) {
			return true
		}
	}
	return false
}

// Filter returns a new stack with only elements satisfying the predicate.
// Order (bottom-to-top) is preserved.
func (s *Stack[T]) Filter(predicate func(T) bool) *Stack[T] {
	result := New[T]()
	for _, item := range s.items {
		if predicate(item) {
			result.Push(item)
		}
	}
	return result
}

// ForEach applies fn to each element, from bottom to top.
func (s *Stack[T]) ForEach(fn func(T)) {
	for _, item := range s.items {
		fn(item)
	}
}

// String returns a human-readable representation of the stack.
// Top of the stack is shown on the right side.
//
//	Stack[bottom -> top]: [1, 2, 3]
func (s *Stack[T]) String() string {
	if s.IsEmpty() {
		return "Stack[empty]"
	}
	parts := make([]string, len(s.items))
	for i, item := range s.items {
		parts[i] = fmt.Sprintf("%v", item)
	}
	return fmt.Sprintf("Stack[bottom -> top]: [%s]", strings.Join(parts, ", "))
}
