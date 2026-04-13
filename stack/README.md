# stack

A generic LIFO (Last-In, First-Out) stack backed by a slice. All operations are O(1) amortized.

## Import

```go
import "github.com/anwar-arif/go-dsa/dsa/stack"
```

## Overview

A stack is the simplest sequential data structure: the last element pushed is the first one popped. The slice backing gives amortized O(1) `Push` (via `append`) and true O(1) `Pop` and `Peek` (index the last element).

## API

| Method | Description | Complexity |
|---|---|---|
| `New[T]()` | Create an empty stack | O(1) |
| `NewFromSlice[T](values)` | Pre-populate; last element of slice becomes the top | O(n) |
| `Push(v T)` | Add `v` to the top | Amortized O(1) |
| `Pop() (T, bool)` | Remove and return the top element; `false` if empty | O(1) |
| `Peek() (T, bool)` | Return the top element without removing it; `false` if empty | O(1) |
| `Size() int` | Number of elements | O(1) |
| `IsEmpty() bool` | Whether the stack has no elements | O(1) |
| `Clear()` | Remove all elements | O(n) |
| `ToSlice() []T` | Copy of contents ordered bottom → top | O(n) |

`Pop` and `Peek` return `(zero value, false)` on an empty stack — no panics.

## Examples

### Basic integer stack

```go
s := stack.New[int]()

s.Push(1)
s.Push(2)
s.Push(3)

fmt.Println(s.Peek())  // (3, true)  — top, not removed
fmt.Println(s.Pop())   // (3, true)
fmt.Println(s.Pop())   // (2, true)
fmt.Println(s.Pop())   // (1, true)
fmt.Println(s.Pop())   // (0, false) — empty
```

### Safe pop with ok check

```go
s := stack.New[string]()
s.Push("hello")

if v, ok := s.Pop(); ok {
    fmt.Println(v) // hello
}

if _, ok := s.Pop(); !ok {
    fmt.Println("stack is empty")
}
```

### Pre-populate from a slice

```go
s := stack.NewFromSlice([]int{1, 2, 3})
// bottom → top: 1, 2, 3

v, _ := s.Pop() // 3  (last element of slice = top of stack)
v, _ = s.Pop()  // 2
fmt.Println(v)
```

### String stack

```go
s := stack.New[string]()
s.Push("first")
s.Push("second")
s.Push("third")

for !s.IsEmpty() {
    v, _ := s.Pop()
    fmt.Println(v)
}
// third
// second
// first
```

### Custom struct

```go
type Frame struct {
    FuncName string
    Line     int
}

callStack := stack.New[Frame]()
callStack.Push(Frame{"main", 10})
callStack.Push(Frame{"parse", 42})
callStack.Push(Frame{"lex", 7})

top, _ := callStack.Peek()
fmt.Println(top.FuncName) // lex
```

### ToSlice snapshot

```go
s := stack.New[int]()
s.Push(10)
s.Push(20)
s.Push(30)

snap := s.ToSlice()
fmt.Println(snap) // [10 20 30] — bottom to top
// Modifying snap does not affect the stack.
```

### Clear and reuse

```go
s := stack.New[int]()
s.Push(1)
s.Push(2)
s.Clear()

fmt.Println(s.IsEmpty()) // true
s.Push(99)
v, _ := s.Pop()
fmt.Println(v) // 99
```

## Implementation notes

- Backed by a plain `[]T` slice; `Push` uses `append` which doubles capacity when full (amortized O(1)).
- `Pop` and `Clear` explicitly zero out vacated slots to allow the garbage collector to reclaim memory when `T` is a pointer or interface type.
- `NewFromSlice` and `ToSlice` both copy the underlying slice so external mutations cannot corrupt the stack's internal state.
