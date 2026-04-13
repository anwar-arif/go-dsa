# queue

A generic FIFO (First-In, First-Out) queue backed by a **ring buffer** (circular buffer). Both `Push` and `Pop` run in amortized O(1).

## Import

```go
import "github.com/anwar-arif/go-dsa/dsa/queue"
```

## Overview

A plain-slice queue makes `Pop` O(n) because every dequeue shifts all remaining elements left. The ring buffer avoids this by maintaining `head` and `tail` indices that advance modulo the buffer capacity — no data is ever moved unless the buffer needs to grow.

```
Capacity = 4, after Push(1,2,3,4) then Pop() twice:

index:  0    1   [2]  [3]
                  ↑
                head = 2   (slots 0,1 are free and will be reused)
```

When the buffer is full, it doubles in size and the elements are copied into a new contiguous layout with `head` reset to 0.

## API

| Method | Description | Complexity |
|---|---|---|
| `New[T]()` | Create an empty queue | O(1) |
| `NewFromSlice[T](values)` | Pre-populate; first element becomes the front | O(n) |
| `Push(v T)` | Add `v` to the back | Amortized O(1) |
| `Pop() (T, bool)` | Remove and return the front element; `false` if empty | O(1) |
| `Peek() (T, bool)` | Return the front element without removing it; `false` if empty | O(1) |
| `Size() int` | Number of elements | O(1) |
| `IsEmpty() bool` | Whether the queue has no elements | O(1) |
| `Clear()` | Remove all elements | O(n) |
| `ToSlice() []T` | Copy of contents ordered front → back | O(n) |

`Pop` and `Peek` return `(zero value, false)` on an empty queue — no panics.

## Examples

### Basic integer queue

```go
q := queue.New[int]()

q.Push(1)
q.Push(2)
q.Push(3)

fmt.Println(q.Peek()) // (1, true)  — front, not removed
fmt.Println(q.Pop())  // (1, true)
fmt.Println(q.Pop())  // (2, true)
fmt.Println(q.Pop())  // (3, true)
fmt.Println(q.Pop())  // (0, false) — empty
```

### Safe pop with ok check

```go
q := queue.New[string]()
q.Push("hello")

if v, ok := q.Pop(); ok {
    fmt.Println(v) // hello
}

if _, ok := q.Pop(); !ok {
    fmt.Println("queue is empty")
}
```

### Pre-populate from a slice

```go
q := queue.NewFromSlice([]int{10, 20, 30})
// front → back: 10, 20, 30

v, _ := q.Pop() // 10 (first element of slice = front of queue)
v, _ = q.Pop()  // 20
fmt.Println(v)
```

### BFS with a queue

```go
type Node struct {
    Val   int
    Left  *Node
    Right *Node
}

func bfs(root *Node) []int {
    if root == nil {
        return nil
    }
    q := queue.New[*Node]()
    q.Push(root)
    var result []int

    for !q.IsEmpty() {
        node, _ := q.Pop()
        result = append(result, node.Val)
        if node.Left != nil {
            q.Push(node.Left)
        }
        if node.Right != nil {
            q.Push(node.Right)
        }
    }
    return result
}
```

### String queue

```go
q := queue.New[string]()
q.Push("first")
q.Push("second")
q.Push("third")

for !q.IsEmpty() {
    v, _ := q.Pop()
    fmt.Println(v)
}
// first
// second
// third
```

### ToSlice snapshot

```go
q := queue.New[int]()
q.Push(10)
q.Push(20)
q.Push(30)

snap := q.ToSlice()
fmt.Println(snap) // [10 20 30] — front to back
// Modifying snap does not affect the queue.
```

### Clear and reuse

```go
q := queue.New[int]()
q.Push(1)
q.Push(2)
q.Clear()

fmt.Println(q.IsEmpty()) // true
q.Push(99)
v, _ := q.Pop()
fmt.Println(v) // 99
```

## Implementation notes

**Ring buffer mechanics:** the backing slice is treated as circular. `head` points to the front element; `tail` is computed as `(head + count) % capacity`. Both advance by 1 (mod capacity) on each `Pop`/`Push` without touching any other element.

**Growth:** when `count == capacity` a new slice of double the capacity is allocated. Elements are copied in order from `head` so that the new layout always starts at index 0, simplifying subsequent wrap-around arithmetic.

**Memory safety:** `Pop` and `Clear` explicitly zero out vacated slots to allow the GC to reclaim memory when `T` is a pointer or interface type.

**vs plain slice:** a naive `[]T` queue uses `append` for push (fast) but requires either `slice[1:]` for pop (O(1) but leaks the underlying array) or copy-shifting (O(n)). The ring buffer gives true O(1) pop with no leaks.
