package stack

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

// ---- New / IsEmpty / Size ---------------------------------------------------

func TestNew_EmptyStack(t *testing.T) {
	s := New[int]()
	if !s.IsEmpty() {
		t.Fatal("new stack should be empty")
	}
	mustEqual(t, "Size()", s.Size(), 0)
}

// ---- Push / Size / IsEmpty --------------------------------------------------

func TestPush_UpdatesSize(t *testing.T) {
	s := New[int]()
	for i := 1; i <= 5; i++ {
		s.Push(i * 10)
		mustEqual(t, "Size()", s.Size(), i)
		if s.IsEmpty() {
			t.Fatal("IsEmpty() must be false after Push")
		}
	}
}

// ---- Pop --------------------------------------------------------------------

func TestPop_LIFO_Order(t *testing.T) {
	s := New[int]()
	for _, v := range []int{1, 2, 3, 4, 5} {
		s.Push(v)
	}
	want := []int{5, 4, 3, 2, 1}
	for i, w := range want {
		v, ok := s.Pop()
		if !ok {
			t.Fatalf("Pop #%d returned false unexpectedly", i)
		}
		mustEqual(t, "Pop()", v, w)
	}
	if !s.IsEmpty() {
		t.Fatal("stack should be empty after draining")
	}
}

func TestPop_OnEmpty_ReturnsFalse(t *testing.T) {
	s := New[int]()
	v, ok := s.Pop()
	if ok {
		t.Fatal("Pop on empty stack should return false")
	}
	mustEqual(t, "zero value", v, 0)
}

func TestPop_UpdatesSize(t *testing.T) {
	s := New[int]()
	s.Push(1)
	s.Push(2)
	s.Pop()
	mustEqual(t, "Size() after one Pop", s.Size(), 1)
	s.Pop()
	mustEqual(t, "Size() after draining", s.Size(), 0)
	if !s.IsEmpty() {
		t.Fatal("should be empty")
	}
}

// ---- Peek -------------------------------------------------------------------

func TestPeek_DoesNotRemove(t *testing.T) {
	s := New[int]()
	s.Push(10)
	s.Push(20)

	v, ok := s.Peek()
	if !ok {
		t.Fatal("Peek returned false on non-empty stack")
	}
	mustEqual(t, "Peek()", v, 20)
	mustEqual(t, "Size() after Peek", s.Size(), 2)
}

func TestPeek_MatchesPop(t *testing.T) {
	s := New[int]()
	s.Push(7)
	s.Push(3)

	peekVal, _ := s.Peek()
	popVal, _ := s.Pop()
	mustEqual(t, "Peek == Pop", peekVal, popVal)
}

func TestPeek_OnEmpty_ReturnsFalse(t *testing.T) {
	s := New[string]()
	v, ok := s.Peek()
	if ok {
		t.Fatal("Peek on empty stack should return false")
	}
	mustEqual(t, "zero value", v, "")
}

// ---- Single element ---------------------------------------------------------

func TestSingleElement(t *testing.T) {
	s := New[int]()
	s.Push(42)
	if v, ok := s.Peek(); !ok || v != 42 {
		t.Fatalf("Peek() want (42, true), got (%d, %v)", v, ok)
	}
	if v, ok := s.Pop(); !ok || v != 42 {
		t.Fatalf("Pop() want (42, true), got (%d, %v)", v, ok)
	}
	if !s.IsEmpty() {
		t.Fatal("stack should be empty after popping last element")
	}
}

// ---- NewFromSlice -----------------------------------------------------------

func TestNewFromSlice_TopIsLastElement(t *testing.T) {
	s := NewFromSlice([]int{1, 2, 3})
	mustEqual(t, "Size()", s.Size(), 3)

	// pop order should be 3, 2, 1
	for _, want := range []int{3, 2, 1} {
		v, ok := s.Pop()
		if !ok || v != want {
			t.Fatalf("Pop() want (%d, true), got (%d, %v)", want, v, ok)
		}
	}
}

func TestNewFromSlice_DoesNotMutateOriginal(t *testing.T) {
	original := []int{10, 20, 30}
	snapshot := []int{10, 20, 30}
	s := NewFromSlice(original)
	s.Pop()
	s.Push(99)
	for i, v := range snapshot {
		if original[i] != v {
			t.Fatalf("NewFromSlice mutated original at index %d", i)
		}
	}
}

func TestNewFromSlice_Empty(t *testing.T) {
	s := NewFromSlice([]int{})
	if !s.IsEmpty() {
		t.Fatal("stack from empty slice should be empty")
	}
}

// ---- ToSlice ----------------------------------------------------------------

func TestToSlice_BottomToTop(t *testing.T) {
	s := New[int]()
	for _, v := range []int{1, 2, 3} {
		s.Push(v)
	}
	mustSliceEqual(t, s.ToSlice(), []int{1, 2, 3})
}

func TestToSlice_DoesNotMutateStack(t *testing.T) {
	s := New[int]()
	s.Push(1)
	s.Push(2)
	out := s.ToSlice()
	out[0] = 99
	// stack should be unaffected
	mustEqual(t, "Size() after modifying ToSlice result", s.Size(), 2)
	v, _ := s.Pop()
	mustEqual(t, "Pop() unaffected", v, 2)
}

func TestToSlice_Empty(t *testing.T) {
	s := New[int]()
	if sl := s.ToSlice(); len(sl) != 0 {
		t.Fatalf("expected empty slice, got %v", sl)
	}
}

// ---- Clear ------------------------------------------------------------------

func TestClear_EmptiesStack(t *testing.T) {
	s := New[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Clear()
	if !s.IsEmpty() {
		t.Fatal("IsEmpty() should be true after Clear")
	}
	mustEqual(t, "Size() after Clear", s.Size(), 0)
}

func TestClear_ThenPush(t *testing.T) {
	s := New[int]()
	s.Push(1)
	s.Clear()
	s.Push(5)
	mustEqual(t, "Size() after Clear+Push", s.Size(), 1)
	v, ok := s.Pop()
	if !ok || v != 5 {
		t.Fatalf("Pop() after Clear+Push: want (5, true), got (%d, %v)", v, ok)
	}
}

func TestClear_OnEmpty_IsNoop(t *testing.T) {
	s := New[int]()
	s.Clear()
	if !s.IsEmpty() {
		t.Fatal("Clear on empty stack should leave it empty")
	}
}

// ---- Interleaved Push / Pop -------------------------------------------------

func TestInterleaved_PushPop(t *testing.T) {
	s := New[int]()
	s.Push(1)
	s.Push(2)
	v, _ := s.Pop() // pops 2
	mustEqual(t, "first pop", v, 2)
	s.Push(3)
	s.Push(4)
	// stack bottom→top: 1, 3, 4
	mustSliceEqual(t, s.ToSlice(), []int{1, 3, 4})
}

// ---- String type ------------------------------------------------------------

func TestStringStack(t *testing.T) {
	s := New[string]()
	s.Push("hello")
	s.Push("world")
	v, ok := s.Pop()
	if !ok || v != "world" {
		t.Fatalf("Pop() want (\"world\", true), got (%q, %v)", v, ok)
	}
	v, ok = s.Pop()
	if !ok || v != "hello" {
		t.Fatalf("Pop() want (\"hello\", true), got (%q, %v)", v, ok)
	}
	if _, ok = s.Pop(); ok {
		t.Fatal("Pop on empty string stack should return false")
	}
}

// ---- Struct type ------------------------------------------------------------

type pair struct{ x, y int }

func TestStructStack(t *testing.T) {
	s := New[pair]()
	s.Push(pair{1, 2})
	s.Push(pair{3, 4})

	v, ok := s.Pop()
	if !ok || v != (pair{3, 4}) {
		t.Fatalf("Pop() want ({3,4}, true), got (%v, %v)", v, ok)
	}
	v, ok = s.Peek()
	if !ok || v != (pair{1, 2}) {
		t.Fatalf("Peek() want ({1,2}, true), got (%v, %v)", v, ok)
	}
}
