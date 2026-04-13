package priorityqueue

import (
	"testing"
)

// ---- helpers ----------------------------------------------------------------

func intMinLess(a, b int) bool { return a < b }
func intMaxLess(a, b int) bool { return a > b }

func drainInts(pq *PriorityQueue[int]) []int {
	out := make([]int, 0, pq.Len())
	for !pq.IsEmpty() {
		out = append(out, pq.Pop())
	}
	return out
}

func isSorted[T any](s []T, less func(a, b T) bool) bool {
	for i := 1; i < len(s); i++ {
		if less(s[i], s[i-1]) {
			return false
		}
	}
	return true
}

// ---- New / IsEmpty / Len ----------------------------------------------------

func TestNew_EmptyQueue(t *testing.T) {
	pq := New(intMinLess)
	if !pq.IsEmpty() {
		t.Fatal("expected IsEmpty() == true for a new queue")
	}
	if pq.Len() != 0 {
		t.Fatalf("expected Len() == 0, got %d", pq.Len())
	}
}

// ---- Push / Len / IsEmpty ---------------------------------------------------

func TestPush_UpdatesLen(t *testing.T) {
	pq := New(intMinLess)
	for i := 1; i <= 5; i++ {
		pq.Push(i)
		if pq.Len() != i {
			t.Fatalf("after %d pushes expected Len() == %d, got %d", i, i, pq.Len())
		}
		if pq.IsEmpty() {
			t.Fatal("IsEmpty() must be false after Push")
		}
	}
}

// ---- Min-heap pop order -----------------------------------------------------

func TestMinHeap_PopOrder(t *testing.T) {
	pq := New(intMinLess)
	for _, v := range []int{5, 1, 4, 2, 3} {
		pq.Push(v)
	}
	got := drainInts(pq)
	if !isSorted(got, intMinLess) {
		t.Fatalf("min-heap pop order wrong: %v", got)
	}
	want := []int{1, 2, 3, 4, 5}
	for i, v := range want {
		if got[i] != v {
			t.Fatalf("index %d: want %d got %d", i, v, got[i])
		}
	}
}

// ---- Max-heap pop order -----------------------------------------------------

func TestMaxHeap_PopOrder(t *testing.T) {
	pq := New(intMaxLess)
	for _, v := range []int{3, 1, 4, 1, 5, 9, 2, 6} {
		pq.Push(v)
	}
	got := drainInts(pq)
	if !isSorted(got, intMaxLess) {
		t.Fatalf("max-heap pop order wrong: %v", got)
	}
}

// ---- Peek -------------------------------------------------------------------

func TestPeek_DoesNotRemove(t *testing.T) {
	pq := New(intMinLess)
	pq.Push(10)
	pq.Push(2)
	pq.Push(7)

	before := pq.Len()
	top := pq.Peek()
	if top != 2 {
		t.Fatalf("Peek() want 2, got %d", top)
	}
	if pq.Len() != before {
		t.Fatal("Peek() must not change Len()")
	}
}

func TestPeek_MatchesPop(t *testing.T) {
	pq := New(intMinLess)
	for _, v := range []int{8, 3, 6, 1, 9} {
		pq.Push(v)
	}
	if pq.Peek() != pq.Pop() {
		t.Fatal("Peek() and Pop() must return the same element")
	}
}

// ---- Pop reduces Len --------------------------------------------------------

func TestPop_UpdatesLen(t *testing.T) {
	pq := New(intMinLess)
	pq.Push(1)
	pq.Push(2)
	pq.Pop()
	if pq.Len() != 1 {
		t.Fatalf("expected Len() == 1 after one Pop, got %d", pq.Len())
	}
	pq.Pop()
	if !pq.IsEmpty() {
		t.Fatal("expected IsEmpty() after draining queue")
	}
}

// ---- Single element ---------------------------------------------------------

func TestSingleElement(t *testing.T) {
	pq := New(intMinLess)
	pq.Push(42)
	if pq.Peek() != 42 {
		t.Fatalf("Peek() want 42, got %d", pq.Peek())
	}
	if v := pq.Pop(); v != 42 {
		t.Fatalf("Pop() want 42, got %d", v)
	}
	if !pq.IsEmpty() {
		t.Fatal("expected IsEmpty() after popping last element")
	}
}

// ---- Duplicate values -------------------------------------------------------

func TestDuplicateValues(t *testing.T) {
	pq := New(intMinLess)
	for _, v := range []int{3, 3, 1, 1, 2, 2} {
		pq.Push(v)
	}
	if pq.Len() != 6 {
		t.Fatalf("expected Len() == 6, got %d", pq.Len())
	}
	got := drainInts(pq)
	if !isSorted(got, intMinLess) {
		t.Fatalf("duplicate values: wrong pop order %v", got)
	}
}

// ---- NewFromSlice -----------------------------------------------------------

func TestNewFromSlice_CorrectOrder(t *testing.T) {
	input := []int{9, 4, 7, 1, 5}
	pq := NewFromSlice(intMinLess, input)
	if pq.Len() != len(input) {
		t.Fatalf("expected Len() == %d, got %d", len(input), pq.Len())
	}
	got := drainInts(pq)
	if !isSorted(got, intMinLess) {
		t.Fatalf("NewFromSlice: wrong pop order %v", got)
	}
}

func TestNewFromSlice_DoesNotMutateOriginal(t *testing.T) {
	input := []int{5, 3, 8, 1}
	snapshot := make([]int, len(input))
	copy(snapshot, input)

	NewFromSlice(intMinLess, input)

	for i, v := range snapshot {
		if input[i] != v {
			t.Fatalf("NewFromSlice mutated original slice at index %d: want %d got %d", i, v, input[i])
		}
	}
}

func TestNewFromSlice_Empty(t *testing.T) {
	pq := NewFromSlice(intMinLess, []int{})
	if !pq.IsEmpty() {
		t.Fatal("expected IsEmpty() for queue built from empty slice")
	}
}

// ---- String type ------------------------------------------------------------

func TestStringMinHeap(t *testing.T) {
	pq := New(func(a, b string) bool { return a < b })
	for _, s := range []string{"banana", "apple", "cherry", "apricot"} {
		pq.Push(s)
	}
	if top := pq.Pop(); top != "apple" {
		t.Fatalf("want \"apple\", got %q", top)
	}
	if top := pq.Pop(); top != "apricot" {
		t.Fatalf("want \"apricot\", got %q", top)
	}
}

// ---- Custom struct ----------------------------------------------------------

type task struct {
	Name     string
	Priority int
}

func taskHigherFirst(a, b task) bool { return a.Priority > b.Priority }

func TestStructMaxPriority(t *testing.T) {
	pq := New(taskHigherFirst)
	pq.Push(task{"low", 1})
	pq.Push(task{"high", 10})
	pq.Push(task{"medium", 5})

	want := []string{"high", "medium", "low"}
	for _, name := range want {
		got := pq.Pop()
		if got.Name != name {
			t.Fatalf("want task %q, got %q", name, got.Name)
		}
	}
}

func TestStructPeek(t *testing.T) {
	pq := New(taskHigherFirst)
	pq.Push(task{"a", 3})
	pq.Push(task{"b", 7})

	if pq.Peek().Name != "b" {
		t.Fatalf("Peek() should return highest priority task")
	}
	if pq.Len() != 2 {
		t.Fatal("Peek() must not remove element")
	}
}

// ---- Panic on empty queue ---------------------------------------------------

func TestPop_PanicsOnEmpty(t *testing.T) {
	pq := New(intMinLess)
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when calling Pop() on an empty queue")
		}
	}()
	pq.Pop()
}

func TestPeek_PanicsOnEmpty(t *testing.T) {
	pq := New(intMinLess)
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when calling Peek() on an empty queue")
		}
	}()
	pq.Peek()
}

// ---- Interleaved push / pop -------------------------------------------------

func TestInterleavedPushPop(t *testing.T) {
	pq := New(intMinLess)
	pq.Push(5)
	pq.Push(1)
	if v := pq.Pop(); v != 1 {
		t.Fatalf("want 1, got %d", v)
	}
	pq.Push(3)
	pq.Push(2)
	// remaining: 2, 3, 5
	got := drainInts(pq)
	want := []int{2, 3, 5}
	for i, v := range want {
		if got[i] != v {
			t.Fatalf("index %d: want %d got %d", i, v, got[i])
		}
	}
}
