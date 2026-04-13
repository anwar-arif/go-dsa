package collections_test

import (
	"testing"

	collections "github.com/yourname/collections" // update to your module path
)

// ── Basic int stack ───────────────────────────────────────────────────────────

func TestPushAndPop_Int(t *testing.T) {
	s := collections.New[int]()
	s.Push(10)
	s.Push(20)
	s.Push(30)

	got, err := s.Pop()
	if err != nil || got != 30 {
		t.Errorf("expected 30, got %d (err=%v)", got, err)
	}
	if s.Size() != 2 {
		t.Errorf("expected size 2, got %d", s.Size())
	}
}

func TestPeek_DoesNotRemove(t *testing.T) {
	s := collections.New[int]()
	s.Push(42)

	val, _ := s.Peek()
	if val != 42 {
		t.Errorf("peek expected 42, got %d", val)
	}
	if s.Size() != 1 {
		t.Error("peek should not remove the element")
	}
}

func TestPopEmpty_ReturnsError(t *testing.T) {
	s := collections.New[int]()
	_, err := s.Pop()
	if err == nil {
		t.Error("expected error when popping empty stack")
	}
}

func TestIsEmpty(t *testing.T) {
	s := collections.New[string]()
	if !s.IsEmpty() {
		t.Error("new stack should be empty")
	}
	s.Push("hello")
	if s.IsEmpty() {
		t.Error("stack should not be empty after push")
	}
}

// ── String stack ──────────────────────────────────────────────────────────────

func TestStringStack(t *testing.T) {
	s := collections.New[string]()
	s.Push("alpha")
	s.Push("beta")
	s.Push("gamma")

	top, _ := s.Pop()
	if top != "gamma" {
		t.Errorf("expected gamma, got %s", top)
	}
}

// ── Struct stack ──────────────────────────────────────────────────────────────

type Point struct{ X, Y int }

func TestStructStack(t *testing.T) {
	s := collections.New[Point]()
	s.Push(Point{1, 2})
	s.Push(Point{3, 4})

	top, _ := s.Peek()
	if top.X != 3 || top.Y != 4 {
		t.Errorf("unexpected top value: %+v", top)
	}
	if s.Size() != 2 {
		t.Errorf("expected size 2, got %d", s.Size())
	}
}

// ── Utility methods ───────────────────────────────────────────────────────────

func TestClear(t *testing.T) {
	s := collections.NewFromSlice([]int{1, 2, 3})
	s.Clear()
	if !s.IsEmpty() {
		t.Error("stack should be empty after Clear()")
	}
}

func TestClone(t *testing.T) {
	original := collections.NewFromSlice([]int{1, 2, 3})
	clone := original.Clone()

	clone.Push(99)
	if original.Size() == clone.Size() {
		t.Error("clone should be independent of original")
	}
}

func TestToSlice(t *testing.T) {
	s := collections.NewFromSlice([]int{10, 20, 30})
	sl := s.ToSlice()
	if sl[2] != 30 {
		t.Errorf("expected last element 30, got %d", sl[2])
	}
}

func TestContains(t *testing.T) {
	s := collections.NewFromSlice([]int{1, 2, 3, 4, 5})
	if !s.Contains(func(v int) bool { return v == 3 }) {
		t.Error("expected stack to contain 3")
	}
	if s.Contains(func(v int) bool { return v == 99 }) {
		t.Error("stack should not contain 99")
	}
}

func TestFilter(t *testing.T) {
	s := collections.NewFromSlice([]int{1, 2, 3, 4, 5, 6})
	evens := s.Filter(func(v int) bool { return v%2 == 0 })
	if evens.Size() != 3 {
		t.Errorf("expected 3 even numbers, got %d", evens.Size())
	}
}

func TestForEach(t *testing.T) {
	s := collections.NewFromSlice([]int{1, 2, 3})
	sum := 0
	s.ForEach(func(v int) { sum += v })
	if sum != 6 {
		t.Errorf("expected sum 6, got %d", sum)
	}
}

func TestString(t *testing.T) {
	s := collections.New[int]()
	if s.String() != "Stack[empty]" {
		t.Errorf("unexpected string for empty stack: %s", s.String())
	}
	s.Push(1)
	s.Push(2)
	expected := "Stack[bottom -> top]: [1, 2]"
	if s.String() != expected {
		t.Errorf("expected %q, got %q", expected, s.String())
	}
}

// ── Capacity hint ─────────────────────────────────────────────────────────────

func TestNewWithCapacity(t *testing.T) {
	s := collections.NewWithCapacity[int](100)
	for i := range 100 {
		s.Push(i)
	}
	if s.Size() != 100 {
		t.Errorf("expected 100 elements, got %d", s.Size())
	}
}
