# priorityqueue

A generic priority queue backed by a binary heap (`container/heap`). The element ordering is fully controlled by a caller-supplied `Less` function, making it suitable for min-heaps, max-heaps, and arbitrary struct orderings with no extra boilerplate.

## Import

```go
import "github.com/anwar-arif/go-dsa/dsa/priorityqueue"
```

## Overview

Internally the queue wraps Go's standard `container/heap`, which implements a binary heap on a slice. Each `Push` and `Pop` runs in O(log n). The `Less(a, b)` comparator you provide determines which element surfaces first:

- `less(a, b) = a < b` → **min-heap** (smallest element popped first)
- `less(a, b) = a > b` → **max-heap** (largest element popped first)
- Any arbitrary comparator for structs

## API

| Method | Description | Complexity |
|---|---|---|
| `New[T](less)` | Create an empty priority queue | O(1) |
| `NewFromSlice[T](less, values)` | Build from an existing slice | O(n) |
| `Push(v T)` | Add an element | O(log n) |
| `Pop() T` | Remove and return the highest-priority element | O(log n) |
| `Peek() T` | Return the highest-priority element without removing it | O(1) |
| `Len() int` | Number of elements | O(1) |
| `IsEmpty() bool` | Whether the queue has no elements | O(1) |

> `Pop` and `Peek` panic if called on an empty queue. Check `IsEmpty` or `Len` first if needed.

## Examples

### Min-heap (int)

```go
pq := priorityqueue.New(func(a, b int) bool { return a < b })

pq.Push(5)
pq.Push(1)
pq.Push(3)

fmt.Println(pq.Pop())  // 1
fmt.Println(pq.Pop())  // 3
fmt.Println(pq.Pop())  // 5
```

### Max-heap (int)

```go
pq := priorityqueue.New(func(a, b int) bool { return a > b })

pq.Push(5)
pq.Push(1)
pq.Push(3)

fmt.Println(pq.Pop())  // 5
fmt.Println(pq.Pop())  // 3
fmt.Println(pq.Pop())  // 1
```

### Custom struct — task scheduler

```go
type Task struct {
    Name     string
    Priority int
}

pq := priorityqueue.New(func(a, b Task) bool {
    return a.Priority > b.Priority // higher int = higher priority
})

pq.Push(Task{"send email",   1})
pq.Push(Task{"fix outage",  10})
pq.Push(Task{"write docs",   3})

for !pq.IsEmpty() {
    fmt.Println(pq.Pop().Name)
}
// fix outage
// write docs
// send email
```

### Pre-populate from a slice

```go
nums := []int{9, 4, 7, 1, 5}
pq := priorityqueue.NewFromSlice(
    func(a, b int) bool { return a < b },
    nums,
)
// Builds the heap in O(n) instead of O(n log n) repeated pushes.

fmt.Println(pq.Pop()) // 1
```

### Peek without consuming

```go
pq := priorityqueue.New(func(a, b int) bool { return a < b })
pq.Push(3)
pq.Push(1)
pq.Push(2)

fmt.Println(pq.Peek()) // 1  — queue still has 3 elements
fmt.Println(pq.Len())  // 3
```

## Implementation notes

- Backed by `container/heap` (standard library binary heap on a slice).
- `NewFromSlice` copies the input slice before calling `heap.Init`, so the original is never mutated.
- The internal `Pop` zeroes out the vacated slot to prevent memory leaks when `T` is a pointer or interface type.
- The comparator convention (`less(a, b)` returns `true` when `a` should be popped before `b`) is identical to `sort.Slice`, so existing comparators can be reused directly.
