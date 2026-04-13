// Package priorityqueue provides a generic priority queue backed by a binary heap.
//
// The priority is determined by a user-supplied Less function:
//
//	less(a, b) returns true if a has higher priority than b
//
// Min-heap (smallest value = highest priority):
//
//	pq := priorityqueue.New(func(a, b int) bool { return a < b })
//
// Max-heap (largest value = highest priority):
//
//	pq := priorityqueue.New(func(a, b int) bool { return a > b })
//
// Custom struct with priority field:
//
//	type Task struct { Name string; Priority int }
//	pq := priorityqueue.New(func(a, b Task) bool { return a.Priority > b.Priority })
package priorityqueue

import "container/heap"

// innerHeap is the internal heap that satisfies container/heap.Interface.
type innerHeap[T any] struct {
	items []T
	less  func(a, b T) bool
}

func (h *innerHeap[T]) Len() int            { return len(h.items) }
func (h *innerHeap[T]) Less(i, j int) bool  { return h.less(h.items[i], h.items[j]) }
func (h *innerHeap[T]) Swap(i, j int)       { h.items[i], h.items[j] = h.items[j], h.items[i] }
func (h *innerHeap[T]) Push(x any)          { h.items = append(h.items, x.(T)) }
func (h *innerHeap[T]) Pop() any {
	n := len(h.items)
	x := h.items[n-1]
	var zero T
	h.items[n-1] = zero // avoid memory leak for pointer/interface types
	h.items = h.items[:n-1]
	return x
}

// PriorityQueue is a generic heap-backed priority queue.
// The zero value is not usable; use New or NewFromSlice to construct one.
type PriorityQueue[T any] struct {
	h *innerHeap[T]
}

// New returns an empty PriorityQueue using the provided Less function to
// determine priority. less(a, b) must return true when a has higher priority
// than b (i.e. a should be popped before b).
func New[T any](less func(a, b T) bool) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		h: &innerHeap[T]{less: less},
	}
}

// NewFromSlice builds a PriorityQueue pre-populated with the given values.
// The slice is copied; the original is not modified.
// Time complexity: O(n).
func NewFromSlice[T any](less func(a, b T) bool, values []T) *PriorityQueue[T] {
	items := make([]T, len(values))
	copy(items, values)
	h := &innerHeap[T]{items: items, less: less}
	heap.Init(h)
	return &PriorityQueue[T]{h: h}
}

// Push adds v to the queue.
// Time complexity: O(log n).
func (pq *PriorityQueue[T]) Push(v T) {
	heap.Push(pq.h, v)
}

// Pop removes and returns the highest-priority element.
// Panics if the queue is empty.
// Time complexity: O(log n).
func (pq *PriorityQueue[T]) Pop() T {
	return heap.Pop(pq.h).(T)
}

// Peek returns the highest-priority element without removing it.
// Panics if the queue is empty.
// Time complexity: O(1).
func (pq *PriorityQueue[T]) Peek() T {
	return pq.h.items[0]
}

// Len returns the number of elements in the queue.
func (pq *PriorityQueue[T]) Len() int {
	return pq.h.Len()
}

// IsEmpty reports whether the queue has no elements.
func (pq *PriorityQueue[T]) IsEmpty() bool {
	return pq.h.Len() == 0
}
