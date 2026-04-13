// Package queue provides a generic FIFO queue backed by a ring buffer
// (circular buffer), giving amortized O(1) Enqueue and O(1) Dequeue.
//
// Using a plain slice would make Dequeue O(n) due to element shifting.
// The ring buffer avoids that by maintaining head/tail indices that wrap
// around the backing array, only reallocating when the buffer is full.
//
// Basic usage:
//
//	q := queue.New[int]()
//	q.Push(1)
//	q.Push(2)
//	v, _ := q.Pop()  // v == 1  (FIFO)
//	v, _ = q.Peek()  // v == 2, queue unchanged
package queue

const minCap = 4 // minimum backing array size

// Queue is a generic FIFO queue backed by a ring buffer.
// The zero value is ready to use.
type Queue[T any] struct {
	items []T
	head  int // index of the front element
	count int // number of live elements
}

// New returns an empty Queue.
func New[T any]() *Queue[T] {
	return &Queue[T]{}
}

// NewFromSlice returns a Queue pre-populated with values.
// The first element of the slice becomes the front of the queue.
// The slice is copied; the original is not modified.
func NewFromSlice[T any](values []T) *Queue[T] {
	q := &Queue[T]{}
	for _, v := range values {
		q.Push(v)
	}
	return q
}

// bufCap returns the current capacity of the backing buffer.
func (q *Queue[T]) bufCap() int { return len(q.items) }

// grow doubles the backing buffer, preserving element order from the front.
func (q *Queue[T]) grow() {
	newCap := max(q.bufCap()*2, minCap)
	newItems := make([]T, newCap)
	for i := 0; i < q.count; i++ {
		newItems[i] = q.items[(q.head+i)%q.bufCap()]
	}
	q.items = newItems
	q.head = 0
}

// Push adds v to the back of the queue. Amortized O(1).
func (q *Queue[T]) Push(v T) {
	if q.count == q.bufCap() {
		q.grow()
	}
	tail := (q.head + q.count) % q.bufCap()
	q.items[tail] = v
	q.count++
}

// Pop removes and returns the front element.
// Returns the zero value and false if the queue is empty. O(1).
func (q *Queue[T]) Pop() (T, bool) {
	if q.count == 0 {
		var zero T
		return zero, false
	}
	v := q.items[q.head]
	var zero T
	q.items[q.head] = zero // prevent memory leak for pointer/interface types
	q.head = (q.head + 1) % q.bufCap()
	q.count--
	return v, true
}

// Peek returns the front element without removing it.
// Returns the zero value and false if the queue is empty. O(1).
func (q *Queue[T]) Peek() (T, bool) {
	if q.count == 0 {
		var zero T
		return zero, false
	}
	return q.items[q.head], true
}

// Size returns the number of elements in the queue. O(1).
func (q *Queue[T]) Size() int { return q.count }

// IsEmpty reports whether the queue has no elements. O(1).
func (q *Queue[T]) IsEmpty() bool { return q.count == 0 }

// Clear removes all elements from the queue.
func (q *Queue[T]) Clear() {
	var zero T
	for i := 0; i < q.count; i++ {
		q.items[(q.head+i)%q.bufCap()] = zero
	}
	q.head = 0
	q.count = 0
}

// ToSlice returns a copy of the queue contents ordered front to back.
// The first element of the returned slice is the current front.
func (q *Queue[T]) ToSlice() []T {
	out := make([]T, q.count)
	for i := 0; i < q.count; i++ {
		out[i] = q.items[(q.head+i)%q.bufCap()]
	}
	return out
}
