package multiset

import (
	"slices"
	"testing"
)

// ---- helpers ----------------------------------------------------------------

func intLess(a, b int) bool { return a < b }

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

func buildMS(vals ...int) *Multiset[int] {
	ms := New(intLess)
	for _, v := range vals {
		ms.Insert(v)
	}
	return ms
}

// ---- New / IsEmpty / Size ---------------------------------------------------

func TestNew_Empty(t *testing.T) {
	ms := New(intLess)
	if !ms.IsEmpty() {
		t.Fatal("new multiset should be empty")
	}
	if ms.Size() != 0 {
		t.Fatalf("Size() want 0, got %d", ms.Size())
	}
}

// ---- Insert / Size / IsEmpty ------------------------------------------------

func TestInsert_UpdatesSize(t *testing.T) {
	ms := New(intLess)
	for i := 1; i <= 5; i++ {
		ms.Insert(10)
		if ms.Size() != i {
			t.Fatalf("after %d inserts: Size() want %d got %d", i, i, ms.Size())
		}
	}
	if ms.IsEmpty() {
		t.Fatal("IsEmpty() must be false after inserts")
	}
}

func TestInsert_Duplicates_SizeCountsAll(t *testing.T) {
	ms := buildMS(3, 3, 3, 1, 1)
	if ms.Size() != 5 {
		t.Fatalf("Size() want 5, got %d", ms.Size())
	}
}

// ---- ToSlice ----------------------------------------------------------------

func TestToSlice_SortedWithDuplicates(t *testing.T) {
	ms := buildMS(5, 1, 3, 1, 4, 1, 5, 9, 2, 6)
	got := ms.ToSlice()
	want := slices.Sorted(slices.Values([]int{5, 1, 3, 1, 4, 1, 5, 9, 2, 6}))
	mustSliceEqual(t, got, want)
}

func TestToSlice_Empty(t *testing.T) {
	ms := New(intLess)
	if s := ms.ToSlice(); len(s) != 0 {
		t.Fatalf("expected empty slice, got %v", s)
	}
}

// ---- Count / Contains -------------------------------------------------------

func TestCount_And_Contains(t *testing.T) {
	ms := buildMS(2, 2, 2, 5)
	if ms.Count(2) != 3 {
		t.Fatalf("Count(2) want 3, got %d", ms.Count(2))
	}
	if ms.Count(5) != 1 {
		t.Fatalf("Count(5) want 1, got %d", ms.Count(5))
	}
	if ms.Count(99) != 0 {
		t.Fatalf("Count(99) want 0, got %d", ms.Count(99))
	}
	if !ms.Contains(2) {
		t.Fatal("Contains(2) should be true")
	}
	if ms.Contains(99) {
		t.Fatal("Contains(99) should be false")
	}
}

// ---- Remove (one occurrence) ------------------------------------------------

func TestRemove_OneOccurrence(t *testing.T) {
	ms := buildMS(4, 4, 4)
	if !ms.Remove(4) {
		t.Fatal("Remove(4) should return true")
	}
	if ms.Count(4) != 2 {
		t.Fatalf("Count(4) after one Remove: want 2, got %d", ms.Count(4))
	}
	if ms.Size() != 2 {
		t.Fatalf("Size() after one Remove: want 2, got %d", ms.Size())
	}
}

func TestRemove_LastOccurrence_ElementGone(t *testing.T) {
	ms := buildMS(7)
	ms.Remove(7)
	if ms.Contains(7) {
		t.Fatal("element should be gone after removing its last occurrence")
	}
	if ms.Size() != 0 {
		t.Fatalf("Size() want 0, got %d", ms.Size())
	}
}

func TestRemove_NotPresent_ReturnsFalse(t *testing.T) {
	ms := buildMS(1, 2, 3)
	if ms.Remove(99) {
		t.Fatal("Remove of absent element should return false")
	}
	if ms.Size() != 3 {
		t.Fatal("Size should be unchanged after failed Remove")
	}
}

// ---- RemoveAll --------------------------------------------------------------

func TestRemoveAll_RemovesEveryOccurrence(t *testing.T) {
	ms := buildMS(6, 6, 6, 7)
	if !ms.RemoveAll(6) {
		t.Fatal("RemoveAll(6) should return true")
	}
	if ms.Count(6) != 0 {
		t.Fatalf("Count(6) should be 0 after RemoveAll, got %d", ms.Count(6))
	}
	if ms.Size() != 1 {
		t.Fatalf("Size() want 1, got %d", ms.Size())
	}
}

func TestRemoveAll_NotPresent_ReturnsFalse(t *testing.T) {
	ms := buildMS(1, 2)
	if ms.RemoveAll(99) {
		t.Fatal("RemoveAll of absent element should return false")
	}
}

// ---- Pop --------------------------------------------------------------------

func TestPop_RemovesMinimum(t *testing.T) {
	ms := buildMS(5, 1, 3)
	v, ok := ms.Pop()
	if !ok || v != 1 {
		t.Fatalf("Pop() want (1, true), got (%d, %v)", v, ok)
	}
	if ms.Size() != 2 {
		t.Fatalf("Size() after Pop: want 2, got %d", ms.Size())
	}
}

func TestPop_DrainsSortedOrder(t *testing.T) {
	input := []int{5, 1, 4, 1, 5, 9, 2, 6, 3}
	ms := buildMS(input...)
	want := slices.Sorted(slices.Values(input))
	for i, w := range want {
		v, ok := ms.Pop()
		if !ok || v != w {
			t.Fatalf("Pop #%d: want (%d, true), got (%d, %v)", i, w, v, ok)
		}
	}
	if !ms.IsEmpty() {
		t.Fatal("multiset should be empty after draining")
	}
}

func TestPop_OnEmpty_ReturnsFalse(t *testing.T) {
	ms := New(intLess)
	_, ok := ms.Pop()
	if ok {
		t.Fatal("Pop on empty multiset should return false")
	}
}

func TestPop_DuplicateMin(t *testing.T) {
	ms := buildMS(1, 1, 2)
	ms.Pop()
	if ms.Count(1) != 1 {
		t.Fatalf("after popping one duplicate min, Count(1) want 1, got %d", ms.Count(1))
	}
}

// ---- Min / Max --------------------------------------------------------------

func TestMin_Max(t *testing.T) {
	ms := buildMS(5, 1, 9, 3)
	min, okMin := ms.Min()
	max, okMax := ms.Max()
	if !okMin || min != 1 {
		t.Fatalf("Min() want (1, true), got (%d, %v)", min, okMin)
	}
	if !okMax || max != 9 {
		t.Fatalf("Max() want (9, true), got (%d, %v)", max, okMax)
	}
}

func TestMin_Max_OnEmpty(t *testing.T) {
	ms := New(intLess)
	if _, ok := ms.Min(); ok {
		t.Fatal("Min() on empty should return false")
	}
	if _, ok := ms.Max(); ok {
		t.Fatal("Max() on empty should return false")
	}
}

func TestMin_Max_SingleElement(t *testing.T) {
	ms := buildMS(42)
	if v, ok := ms.Min(); !ok || v != 42 {
		t.Fatalf("Min() want (42, true), got (%d, %v)", v, ok)
	}
	if v, ok := ms.Max(); !ok || v != 42 {
		t.Fatalf("Max() want (42, true), got (%d, %v)", v, ok)
	}
}

// ---- Floor / Ceiling --------------------------------------------------------

func TestFloor(t *testing.T) {
	ms := buildMS(10, 20, 30, 40, 50)

	cases := []struct{ query, want int }{
		{25, 20}, // between 20 and 30
		{20, 20}, // exact hit
		{50, 50}, // exact hit at max
		{55, 50}, // beyond max
	}
	for _, tc := range cases {
		v, ok := ms.Floor(tc.query)
		if !ok || v != tc.want {
			t.Errorf("Floor(%d) want (%d, true), got (%d, %v)", tc.query, tc.want, v, ok)
		}
	}
}

func TestFloor_NoneSmaller_ReturnsFalse(t *testing.T) {
	ms := buildMS(10, 20, 30)
	if _, ok := ms.Floor(5); ok {
		t.Fatal("Floor(5) should return false when all elements are larger")
	}
}

func TestCeiling(t *testing.T) {
	ms := buildMS(10, 20, 30, 40, 50)

	cases := []struct{ query, want int }{
		{25, 30}, // between 20 and 30
		{30, 30}, // exact hit
		{10, 10}, // exact hit at min
		{5, 10},  // below min
	}
	for _, tc := range cases {
		v, ok := ms.Ceiling(tc.query)
		if !ok || v != tc.want {
			t.Errorf("Ceiling(%d) want (%d, true), got (%d, %v)", tc.query, tc.want, v, ok)
		}
	}
}

func TestCeiling_NoneLarger_ReturnsFalse(t *testing.T) {
	ms := buildMS(10, 20, 30)
	if _, ok := ms.Ceiling(35); ok {
		t.Fatal("Ceiling(35) should return false when all elements are smaller")
	}
}

func TestFloor_Ceiling_WithDuplicates(t *testing.T) {
	ms := buildMS(5, 5, 10, 10)
	if v, ok := ms.Floor(5); !ok || v != 5 {
		t.Fatalf("Floor(5) want 5, got %d %v", v, ok)
	}
	if v, ok := ms.Ceiling(5); !ok || v != 5 {
		t.Fatalf("Ceiling(5) want 5, got %d %v", v, ok)
	}
}

// ---- Rank -------------------------------------------------------------------

func TestRank(t *testing.T) {
	// elements: 10, 20, 20, 30
	ms := buildMS(20, 10, 30, 20)

	cases := []struct {
		v, want int
	}{
		{10, 0},  // nothing is less than 10
		{20, 1},  // only 10 is less than 20
		{30, 3},  // 10, 20, 20 are less than 30
		{5, 0},   // nothing less than 5
		{100, 4}, // all 4 elements are less than 100
	}
	for _, tc := range cases {
		if got := ms.Rank(tc.v); got != tc.want {
			t.Errorf("Rank(%d) want %d, got %d", tc.v, tc.want, got)
		}
	}
}

// ---- Kth --------------------------------------------------------------------

func TestKth(t *testing.T) {
	// sorted: 1, 1, 2, 3, 5
	ms := buildMS(3, 1, 5, 1, 2)

	cases := []struct{ k, want int }{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 5},
	}
	for _, tc := range cases {
		v, ok := ms.Kth(tc.k)
		if !ok || v != tc.want {
			t.Errorf("Kth(%d) want (%d, true), got (%d, %v)", tc.k, tc.want, v, ok)
		}
	}
}

func TestKth_OutOfRange(t *testing.T) {
	ms := buildMS(1, 2, 3)
	if _, ok := ms.Kth(-1); ok {
		t.Fatal("Kth(-1) should return false")
	}
	if _, ok := ms.Kth(3); ok {
		t.Fatal("Kth(3) on 3-element set should return false")
	}
}

// ---- String type ------------------------------------------------------------

func TestStringMultiset(t *testing.T) {
	ms := New(func(a, b string) bool { return a < b })
	for _, s := range []string{"banana", "apple", "apple", "cherry"} {
		ms.Insert(s)
	}
	if ms.Count("apple") != 2 {
		t.Fatalf("Count(\"apple\") want 2, got %d", ms.Count("apple"))
	}
	v, _ := ms.Min()
	if v != "apple" {
		t.Fatalf("Min() want \"apple\", got %q", v)
	}
	v, _ = ms.Pop()
	if v != "apple" {
		t.Fatalf("Pop() want \"apple\", got %q", v)
	}
	if ms.Count("apple") != 1 {
		t.Fatalf("Count(\"apple\") after one Pop: want 1, got %d", ms.Count("apple"))
	}
}

// ---- Custom struct ----------------------------------------------------------

type point struct{ x, y int }

func TestStructMultiset(t *testing.T) {
	// Order by x, then y
	less := func(a, b point) bool {
		if a.x != b.x {
			return a.x < b.x
		}
		return a.y < b.y
	}
	ms := New(less)
	ms.Insert(point{3, 0})
	ms.Insert(point{1, 5})
	ms.Insert(point{1, 2})
	ms.Insert(point{1, 2}) // duplicate

	if ms.Count(point{1, 2}) != 2 {
		t.Fatalf("Count want 2, got %d", ms.Count(point{1, 2}))
	}

	got := ms.ToSlice()
	want := []point{{1, 2}, {1, 2}, {1, 5}, {3, 0}}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("index %d: want %v, got %v", i, w, got[i])
		}
	}
}

// ---- Interleaved insert / remove --------------------------------------------

func TestInterleaved_InsertRemove(t *testing.T) {
	ms := New(intLess)
	ms.Insert(5)
	ms.Insert(3)
	ms.Insert(8)
	ms.Remove(3)
	ms.Insert(1)
	ms.Insert(3) // re-add 3
	// present: 1, 3, 5, 8
	mustSliceEqual(t, ms.ToSlice(), []int{1, 3, 5, 8})
	if ms.Size() != 4 {
		t.Fatalf("Size() want 4, got %d", ms.Size())
	}
}

// ---- Monotone Pop with duplicates -------------------------------------------

func TestPop_MonotoneWithDuplicates(t *testing.T) {
	ms := buildMS(2, 1, 2, 1, 3)
	want := []int{1, 1, 2, 2, 3}
	for i, w := range want {
		v, ok := ms.Pop()
		if !ok || v != w {
			t.Fatalf("Pop #%d: want (%d, true), got (%d, %v)", i, w, v, ok)
		}
	}
}
