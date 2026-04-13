package queue

import "testing"

// ---- helpers ----------------------------------------------------------------

func mustEqual[T comparable](t *testing.T, label string, got, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: got %v, want %v", label, got, want)
	}
}

func mustSliceEqual[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("slice length: got %d want %d\n  got:  %v\n  want: %v", len(got), len(want), got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d: got %v want %v\n  full got:  %v\n  full want: %v", i, got[i], want[i], got, want)
		}
	}
}

func drainInts(q *Queue[int]) []int {
	out := make([]int, 0, q.Size())
	for !q.IsEmpty() {
		v, _ := q.Pop()
		out = append(out, v)
	}
	return out
}

// ---- New / IsEmpty / Size ---------------------------------------------------

func TestNew_EmptyQueue(t *testing.T) {
	q := New[int]()
	if !q.IsEmpty() {
		t.Fatal("new queue should be empty")
	}
	mustEqual(t, "Size()", q.Size(), 0)
}

// ---- Push / Size / IsEmpty -----------------------------------------------

func TestPush_UpdatesSize(t *testing.T) {
	q := New[int]()
	for i := 1; i <= 5; i++ {
		q.Push(i * 10)
		mustEqual(t, "Size()", q.Size(), i)
		if q.IsEmpty() {
			t.Fatal("IsEmpty() must be false after Push")
		}
	}
}

// ---- Pop ----------------------------------------------------------------

func TestPop_FIFO_Order(t *testing.T) {
	q := New[int]()
	for _, v := range []int{1, 2, 3, 4, 5} {
		q.Push(v)
	}
	want := []int{1, 2, 3, 4, 5}
	for i, w := range want {
		v, ok := q.Pop()
		if !ok {
			t.Fatalf("Pop #%d returned false unexpectedly", i)
		}
		mustEqual(t, "Pop()", v, w)
	}
	if !q.IsEmpty() {
		t.Fatal("queue should be empty after draining")
	}
}

func TestPop_OnEmpty_ReturnsFalse(t *testing.T) {
	q := New[int]()
	v, ok := q.Pop()
	if ok {
		t.Fatal("Pop on empty queue should return false")
	}
	mustEqual(t, "zero value", v, 0)
}

func TestPop_UpdatesSize(t *testing.T) {
	q := New[int]()
	q.Push(1)
	q.Push(2)
	q.Pop()
	mustEqual(t, "Size() after one Pop", q.Size(), 1)
	q.Pop()
	mustEqual(t, "Size() after draining", q.Size(), 0)
	if !q.IsEmpty() {
		t.Fatal("should be empty")
	}
}

// ---- Peek -------------------------------------------------------------------

func TestPeek_DoesNotRemove(t *testing.T) {
	q := New[int]()
	q.Push(10)
	q.Push(20)

	v, ok := q.Peek()
	if !ok {
		t.Fatal("Peek returned false on non-empty queue")
	}
	mustEqual(t, "Peek()", v, 10) // front is 10, not 20
	mustEqual(t, "Size() after Peek", q.Size(), 2)
}

func TestPeek_MatchesPop(t *testing.T) {
	q := New[int]()
	q.Push(7)
	q.Push(3)

	peekVal, _ := q.Peek()
	dequeueVal, _ := q.Pop()
	mustEqual(t, "Peek == Pop", peekVal, dequeueVal)
}

func TestPeek_OnEmpty_ReturnsFalse(t *testing.T) {
	q := New[string]()
	v, ok := q.Peek()
	if ok {
		t.Fatal("Peek on empty queue should return false")
	}
	mustEqual(t, "zero value", v, "")
}

// ---- Single element ---------------------------------------------------------

func TestSingleElement(t *testing.T) {
	q := New[int]()
	q.Push(42)
	if v, ok := q.Peek(); !ok || v != 42 {
		t.Fatalf("Peek() want (42, true), got (%d, %v)", v, ok)
	}
	if v, ok := q.Pop(); !ok || v != 42 {
		t.Fatalf("Pop() want (42, true), got (%d, %v)", v, ok)
	}
	if !q.IsEmpty() {
		t.Fatal("queue should be empty after dequeuing last element")
	}
}

// ---- NewFromSlice -----------------------------------------------------------

func TestNewFromSlice_FrontIsFirstElement(t *testing.T) {
	q := NewFromSlice([]int{1, 2, 3})
	mustEqual(t, "Size()", q.Size(), 3)
	mustSliceEqual(t, drainInts(q), []int{1, 2, 3})
}

func TestNewFromSlice_DoesNotMutateOriginal(t *testing.T) {
	original := []int{10, 20, 30}
	snapshot := []int{10, 20, 30}
	q := NewFromSlice(original)
	q.Pop()
	q.Push(99)
	for i, v := range snapshot {
		if original[i] != v {
			t.Fatalf("NewFromSlice mutated original at index %d", i)
		}
	}
}

func TestNewFromSlice_Empty(t *testing.T) {
	q := NewFromSlice([]int{})
	if !q.IsEmpty() {
		t.Fatal("queue from empty slice should be empty")
	}
}

// ---- ToSlice ----------------------------------------------------------------

func TestToSlice_FrontToBack(t *testing.T) {
	q := New[int]()
	for _, v := range []int{1, 2, 3} {
		q.Push(v)
	}
	mustSliceEqual(t, q.ToSlice(), []int{1, 2, 3})
}

func TestToSlice_DoesNotMutateQueue(t *testing.T) {
	q := New[int]()
	q.Push(1)
	q.Push(2)
	out := q.ToSlice()
	out[0] = 99
	mustEqual(t, "Size() after modifying ToSlice result", q.Size(), 2)
	v, _ := q.Pop()
	mustEqual(t, "Pop() unaffected", v, 1)
}

func TestToSlice_Empty(t *testing.T) {
	q := New[int]()
	if sl := q.ToSlice(); len(sl) != 0 {
		t.Fatalf("expected empty slice, got %v", sl)
	}
}

// ---- Clear ------------------------------------------------------------------

func TestClear_EmptiesQueue(t *testing.T) {
	q := New[int]()
	q.Push(1)
	q.Push(2)
	q.Push(3)
	q.Clear()
	if !q.IsEmpty() {
		t.Fatal("IsEmpty() should be true after Clear")
	}
	mustEqual(t, "Size() after Clear", q.Size(), 0)
}

func TestClear_ThenPush(t *testing.T) {
	q := New[int]()
	q.Push(1)
	q.Clear()
	q.Push(5)
	mustEqual(t, "Size() after Clear+Push", q.Size(), 1)
	v, ok := q.Pop()
	if !ok || v != 5 {
		t.Fatalf("Pop() after Clear+Push: want (5, true), got (%d, %v)", v, ok)
	}
}

func TestClear_OnEmpty_IsNoop(t *testing.T) {
	q := New[int]()
	q.Clear()
	if !q.IsEmpty() {
		t.Fatal("Clear on empty queue should leave it empty")
	}
}

// ---- Interleaved Push / Pop ------------------------------------------

func TestInterleaved_PushPop(t *testing.T) {
	q := New[int]()
	q.Push(1)
	q.Push(2)
	v, _ := q.Pop() // dequeues 1
	mustEqual(t, "first dequeue", v, 1)
	q.Push(3)
	q.Push(4)
	// queue front→back: 2, 3, 4
	mustSliceEqual(t, q.ToSlice(), []int{2, 3, 4})
}

// ---- Ring buffer wrap-around ------------------------------------------------
// These tests exercise the circular nature of the backing buffer by dequeuing
// enough elements to push head past the midpoint, then enqueuing more so that
// the tail wraps around.

func TestRingBuffer_WrapAround(t *testing.T) {
	q := New[int]()
	// Fill to initial capacity (minCap == 4)
	for i := range 4 {
		q.Push(i)
	}
	// Pop 3 elements — head is now at index 3
	for range 3 {
		q.Pop()
	}
	// Push 3 more — tail wraps: slots 0, 1, 2 get reused
	q.Push(10)
	q.Push(11)
	q.Push(12)
	// Queue should be: 3, 10, 11, 12
	mustSliceEqual(t, drainInts(q), []int{3, 10, 11, 12})
}

func TestRingBuffer_GrowWhileWrapped(t *testing.T) {
	q := New[int]()
	// Force a wrapped state then trigger a grow
	for i := range 4 {
		q.Push(i) // 0,1,2,3  — buffer full at cap=4
	}
	for range 2 {
		q.Pop() // consume 0,1 — head=2
	}
	// Push 3 more — fills the 2 remaining slots + forces a grow
	q.Push(10)
	q.Push(11)
	q.Push(12) // triggers grow; head was mid-buffer
	// Queue should be: 2, 3, 10, 11, 12
	mustSliceEqual(t, drainInts(q), []int{2, 3, 10, 11, 12})
}

func TestRingBuffer_LargeSequence(t *testing.T) {
	// Push 1..100, dequeue all, verify FIFO and correct size at each step.
	q := New[int]()
	const n = 100
	for i := 1; i <= n; i++ {
		q.Push(i)
	}
	mustEqual(t, "Size() after 100 enqueues", q.Size(), n)
	for i := 1; i <= n; i++ {
		v, ok := q.Pop()
		if !ok || v != i {
			t.Fatalf("Pop #%d: want (%d, true), got (%d, %v)", i, i, v, ok)
		}
	}
	if !q.IsEmpty() {
		t.Fatal("should be empty after draining")
	}
}

// ---- String type ------------------------------------------------------------

func TestStringQueue(t *testing.T) {
	q := New[string]()
	q.Push("hello")
	q.Push("world")
	v, ok := q.Pop()
	if !ok || v != "hello" {
		t.Fatalf("Pop() want (\"hello\", true), got (%q, %v)", v, ok)
	}
	v, ok = q.Pop()
	if !ok || v != "world" {
		t.Fatalf("Pop() want (\"world\", true), got (%q, %v)", v, ok)
	}
	if _, ok = q.Pop(); ok {
		t.Fatal("Pop on empty string queue should return false")
	}
}

// ---- Struct type ------------------------------------------------------------

type pair struct{ x, y int }

func TestStructQueue(t *testing.T) {
	q := New[pair]()
	q.Push(pair{1, 2})
	q.Push(pair{3, 4})

	v, ok := q.Pop()
	if !ok || v != (pair{1, 2}) {
		t.Fatalf("Pop() want ({1,2}, true), got (%v, %v)", v, ok)
	}
	v, ok = q.Peek()
	if !ok || v != (pair{3, 4}) {
		t.Fatalf("Peek() want ({3,4}, true), got (%v, %v)", v, ok)
	}
}
