# multiset

A generic sorted multiset backed by a **treap** (randomized binary search tree). Equivalent to C++'s `std::multiset`: elements are kept in sorted order, duplicates are allowed, and all structural operations run in O(log n) expected time.

## Import

```go
import "github.com/anwar-arif/go-dsa/dsa/multiset"
```

## Overview

A multiset is an ordered collection that permits duplicate values. Unlike a priority queue (which only exposes the top element), a multiset supports arbitrary lookups: floor/ceiling queries, rank, k-th element, and ordered iteration — all in O(log n).

The backing treap stores each **distinct** key once with a `count` field for duplicates, and a `size` field on every node for O(log n) order-statistic queries (`Rank`, `Kth`).

## API

| Method | Description | Complexity |
|---|---|---|
| `New[T](less)` | Create an empty multiset | O(1) |
| `Insert(v T)` | Add one occurrence of `v` | O(log n) |
| `Remove(v T) bool` | Remove one occurrence of `v`; `false` if absent | O(log n) |
| `RemoveAll(v T) bool` | Remove all occurrences of `v`; `false` if absent | O(log n) |
| `Pop() (T, bool)` | Remove and return the smallest element | O(log n) |
| `Count(v T) int` | Number of occurrences of `v` | O(log n) |
| `Contains(v T) bool` | Whether `v` is present at least once | O(log n) |
| `Size() int` | Total elements including duplicates | O(1) |
| `IsEmpty() bool` | Whether the multiset has no elements | O(1) |
| `Min() (T, bool)` | Smallest element | O(log n) |
| `Max() (T, bool)` | Largest element | O(log n) |
| `Floor(v T) (T, bool)` | Largest element `<= v` | O(log n) |
| `Ceiling(v T) (T, bool)` | Smallest element `>= v` | O(log n) |
| `Rank(v T) int` | Number of elements strictly less than `v` | O(log n) |
| `Kth(k int) (T, bool)` | k-th smallest element (0-indexed, duplicates counted) | O(log n) |
| `ToSlice() []T` | All elements in sorted order with duplicates | O(n) |

Methods returning `(T, bool)` use the bool to signal absence (empty set or value not found), avoiding panics.

## Examples

### Integer multiset (ascending order)

```go
ms := multiset.New(func(a, b int) bool { return a < b })

ms.Insert(3)
ms.Insert(1)
ms.Insert(3)
ms.Insert(2)

fmt.Println(ms.ToSlice()) // [1 2 3 3]
fmt.Println(ms.Count(3))  // 2
fmt.Println(ms.Size())    // 4
```

### Remove one vs all occurrences

```go
ms.Remove(3)              // removes one 3
fmt.Println(ms.Count(3))  // 1

ms.RemoveAll(3)            // removes the last 3
fmt.Println(ms.Contains(3)) // false
```

### Pop — draining in sorted order

```go
ms := multiset.New(func(a, b int) bool { return a < b })
for _, v := range []int{5, 1, 3, 1, 4} {
    ms.Insert(v)
}

for !ms.IsEmpty() {
    v, _ := ms.Pop()
    fmt.Print(v, " ") // 1 1 3 4 5
}
```

### Floor and Ceiling

```go
ms := multiset.New(func(a, b int) bool { return a < b })
for _, v := range []int{10, 20, 30, 40, 50} {
    ms.Insert(v)
}

f, _ := ms.Floor(25)    // 20 — largest element <= 25
c, _ := ms.Ceiling(25)  // 30 — smallest element >= 25
fmt.Println(f, c)
```

### Rank and Kth — order statistics

```go
ms := multiset.New(func(a, b int) bool { return a < b })
for _, v := range []int{10, 20, 20, 30} {
    ms.Insert(v)
}

fmt.Println(ms.Rank(20))    // 1 — one element (10) is strictly less than 20
fmt.Println(ms.Rank(30))    // 3 — three elements (10, 20, 20) are less than 30

v, _ := ms.Kth(0)  // 10 — 0th smallest
v, _ = ms.Kth(2)   // 20 — 2nd smallest (second duplicate)
fmt.Println(v)
```

### Custom struct ordered by a field

```go
type Task struct {
    Name     string
    Priority int
}

ms := multiset.New(func(a, b Task) bool {
    return a.Priority < b.Priority
})

ms.Insert(Task{"write docs",  1})
ms.Insert(Task{"fix bug",     5})
ms.Insert(Task{"fix bug",     5}) // duplicate allowed
ms.Insert(Task{"fix outage", 10})

fmt.Println(ms.Size())                 // 4
fmt.Println(ms.Count(Task{"fix bug", 5})) // 2

min, _ := ms.Min() // Task{write docs, 1}
max, _ := ms.Max() // Task{fix outage, 10}
```

### String multiset

```go
ms := multiset.New(func(a, b string) bool { return a < b })
for _, s := range []string{"banana", "apple", "apple", "cherry"} {
    ms.Insert(s)
}

fmt.Println(ms.ToSlice()) // [apple apple banana cherry]

v, _ := ms.Ceiling("b")  // banana
fmt.Println(v)
```

## Implementation notes

**Treap (tree + heap):** each node carries a BST key and a random priority. The random priorities maintain the heap invariant which, in expectation, keeps the tree height at O(log n) — with no rotations bookkeeping required beyond simple left/right rotations on insert.

**Duplicate handling:** distinct keys are stored once with a `count` field. `size` on every node tracks the total element count in the subtree (sum of all counts), enabling O(log n) `Rank` and `Kth` without walking the tree.

**Comparison to `std::multiset`:** C++'s multiset uses a red-black tree (O(log n) worst-case). The treap gives O(log n) **expected** time. The probability of degenerate behaviour is the same as quicksort hitting its worst case on every single partition — negligible in practice.
