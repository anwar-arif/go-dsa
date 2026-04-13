// Package multiset provides a generic sorted multiset backed by a treap
// (randomized binary search tree), giving O(log n) expected time for all
// structural operations — equivalent to C++'s std::multiset.
//
// Priority is determined by a user-supplied Less function:
//
//	less(a, b) returns true if a should be ordered before b
//
// Ascending order (default for numbers):
//
//	ms := multiset.New(func(a, b int) bool { return a < b })
//
// Descending order:
//
//	ms := multiset.New(func(a, b int) bool { return a > b })
//
// Custom struct ordered by a field:
//
//	type Task struct { Name string; Priority int }
//	ms := multiset.New(func(a, b Task) bool { return a.Priority < b.Priority })
package multiset

import "math/rand/v2"

// node is a treap node. Each distinct key is stored once; `count` tracks
// duplicate occurrences. `size` is the total element count in the subtree
// (i.e. sum of all counts), which enables O(log n) Rank and Kth queries.
type node[T any] struct {
	key      T
	priority uint32
	count    int // occurrences of this exact key
	size     int // total elements in this subtree (counts included)
	left     *node[T]
	right    *node[T]
}

func newNode[T any](key T) *node[T] {
	return &node[T]{key: key, priority: rand.Uint32(), count: 1, size: 1}
}

func sz[T any](n *node[T]) int {
	if n == nil {
		return 0
	}
	return n.size
}

// pull recomputes n.size from its children and its own count.
func pull[T any](n *node[T]) {
	if n != nil {
		n.size = n.count + sz(n.left) + sz(n.right)
	}
}

func rotateRight[T any](n *node[T]) *node[T] {
	l := n.left
	n.left = l.right
	l.right = n
	pull(n)
	pull(l)
	return l
}

func rotateLeft[T any](n *node[T]) *node[T] {
	r := n.right
	n.right = r.left
	r.left = n
	pull(n)
	pull(r)
	return r
}

// insert adds one occurrence of key, maintaining the treap invariant.
func insert[T any](n *node[T], key T, less func(a, b T) bool) *node[T] {
	if n == nil {
		return newNode[T](key)
	}
	if less(key, n.key) {
		n.left = insert(n.left, key, less)
		if n.left.priority > n.priority {
			n = rotateRight(n)
		}
	} else if less(n.key, key) {
		n.right = insert(n.right, key, less)
		if n.right.priority > n.priority {
			n = rotateLeft(n)
		}
	} else {
		// key == n.key: just bump the count, no structural change needed.
		n.count++
	}
	pull(n)
	return n
}

// mergeNodes merges two treap subtrees where every key in l is less than
// every key in r.
func mergeNodes[T any](l, r *node[T]) *node[T] {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}
	if l.priority > r.priority {
		l.right = mergeNodes(l.right, r)
		pull(l)
		return l
	}
	r.left = mergeNodes(l, r.left)
	pull(r)
	return r
}

// removeOne removes exactly one occurrence of key.
// Returns the updated root and whether the key was found.
func removeOne[T any](n *node[T], key T, less func(a, b T) bool) (*node[T], bool) {
	if n == nil {
		return nil, false
	}
	var ok bool
	if less(key, n.key) {
		n.left, ok = removeOne(n.left, key, less)
	} else if less(n.key, key) {
		n.right, ok = removeOne(n.right, key, less)
	} else {
		ok = true
		if n.count > 1 {
			n.count--
		} else {
			return mergeNodes(n.left, n.right), true
		}
	}
	pull(n)
	return n, ok
}

// removeAll removes every occurrence of key.
func removeAll[T any](n *node[T], key T, less func(a, b T) bool) (*node[T], bool) {
	if n == nil {
		return nil, false
	}
	var ok bool
	if less(key, n.key) {
		n.left, ok = removeAll(n.left, key, less)
	} else if less(n.key, key) {
		n.right, ok = removeAll(n.right, key, less)
	} else {
		return mergeNodes(n.left, n.right), true
	}
	pull(n)
	return n, ok
}

func countKey[T any](n *node[T], key T, less func(a, b T) bool) int {
	if n == nil {
		return 0
	}
	if less(key, n.key) {
		return countKey(n.left, key, less)
	}
	if less(n.key, key) {
		return countKey(n.right, key, less)
	}
	return n.count
}

func minNode[T any](n *node[T]) *node[T] {
	for n.left != nil {
		n = n.left
	}
	return n
}

func maxNode[T any](n *node[T]) *node[T] {
	for n.right != nil {
		n = n.right
	}
	return n
}

// floor returns the node with the largest key <= v.
func floor[T any](n *node[T], v T, less func(a, b T) bool) *node[T] {
	if n == nil {
		return nil
	}
	// v == n.key
	if !less(n.key, v) && !less(v, n.key) {
		return n
	}
	// v < n.key: answer is in left subtree
	if less(v, n.key) {
		return floor(n.left, v, less)
	}
	// v > n.key: n is a candidate; a better answer may exist in right subtree
	if best := floor(n.right, v, less); best != nil {
		return best
	}
	return n
}

// ceiling returns the node with the smallest key >= v.
func ceiling[T any](n *node[T], v T, less func(a, b T) bool) *node[T] {
	if n == nil {
		return nil
	}
	// v == n.key
	if !less(n.key, v) && !less(v, n.key) {
		return n
	}
	// v > n.key: answer is in right subtree
	if less(n.key, v) {
		return ceiling(n.right, v, less)
	}
	// v < n.key: n is a candidate; a better answer may exist in left subtree
	if best := ceiling(n.left, v, less); best != nil {
		return best
	}
	return n
}

// rank returns the number of elements strictly less than key.
func rank[T any](n *node[T], key T, less func(a, b T) bool) int {
	if n == nil {
		return 0
	}
	if !less(n.key, key) && !less(key, n.key) {
		// key == n.key: everything in left subtree is smaller
		return sz(n.left)
	}
	if less(key, n.key) {
		return rank(n.left, key, less)
	}
	// key > n.key: left subtree + this node's count + recurse right
	return sz(n.left) + n.count + rank(n.right, key, less)
}

// kth returns the node containing the k-th element (0-indexed, duplicates counted).
func kth[T any](n *node[T], k int) *node[T] {
	if n == nil {
		return nil
	}
	if k < sz(n.left) {
		return kth(n.left, k)
	}
	k -= sz(n.left)
	if k < n.count {
		return n
	}
	return kth(n.right, k-n.count)
}

func toSlice[T any](n *node[T], out *[]T) {
	if n == nil {
		return
	}
	toSlice(n.left, out)
	for i := 0; i < n.count; i++ {
		*out = append(*out, n.key)
	}
	toSlice(n.right, out)
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

// Multiset is a generic sorted multiset. Duplicate elements are allowed.
// The zero value is not usable; construct with New.
type Multiset[T any] struct {
	root *node[T]
	less func(a, b T) bool
}

// New returns an empty Multiset ordered by less.
// less(a, b) must return true when a is ordered strictly before b.
func New[T any](less func(a, b T) bool) *Multiset[T] {
	return &Multiset[T]{less: less}
}

// Insert adds one occurrence of v. O(log n) expected.
func (ms *Multiset[T]) Insert(v T) {
	ms.root = insert(ms.root, v, ms.less)
}

// Remove removes one occurrence of v.
// Returns false if v is not present. O(log n) expected.
func (ms *Multiset[T]) Remove(v T) bool {
	var ok bool
	ms.root, ok = removeOne(ms.root, v, ms.less)
	return ok
}

// RemoveAll removes every occurrence of v.
// Returns false if v is not present. O(log n) expected.
func (ms *Multiset[T]) RemoveAll(v T) bool {
	var ok bool
	ms.root, ok = removeAll(ms.root, v, ms.less)
	return ok
}

// Pop removes and returns the smallest (first-ordered) element.
// Returns the zero value and false if the multiset is empty. O(log n) expected.
func (ms *Multiset[T]) Pop() (T, bool) {
	if ms.root == nil {
		var zero T
		return zero, false
	}
	v := minNode(ms.root).key
	ms.root, _ = removeOne(ms.root, v, ms.less)
	return v, true
}

// Count returns the number of occurrences of v. O(log n) expected.
func (ms *Multiset[T]) Count(v T) int {
	return countKey(ms.root, v, ms.less)
}

// Contains reports whether v is present at least once. O(log n) expected.
func (ms *Multiset[T]) Contains(v T) bool {
	return countKey(ms.root, v, ms.less) > 0
}

// Size returns the total number of elements including duplicates. O(1).
func (ms *Multiset[T]) Size() int {
	return sz(ms.root)
}

// IsEmpty reports whether the multiset has no elements. O(1).
func (ms *Multiset[T]) IsEmpty() bool {
	return ms.root == nil
}

// Min returns the smallest element.
// Returns the zero value and false if the multiset is empty. O(log n).
func (ms *Multiset[T]) Min() (T, bool) {
	if ms.root == nil {
		var zero T
		return zero, false
	}
	return minNode(ms.root).key, true
}

// Max returns the largest element.
// Returns the zero value and false if the multiset is empty. O(log n).
func (ms *Multiset[T]) Max() (T, bool) {
	if ms.root == nil {
		var zero T
		return zero, false
	}
	return maxNode(ms.root).key, true
}

// Floor returns the largest element that is <= v.
// Returns the zero value and false if no such element exists. O(log n) expected.
func (ms *Multiset[T]) Floor(v T) (T, bool) {
	n := floor(ms.root, v, ms.less)
	if n == nil {
		var zero T
		return zero, false
	}
	return n.key, true
}

// Ceiling returns the smallest element that is >= v.
// Returns the zero value and false if no such element exists. O(log n) expected.
func (ms *Multiset[T]) Ceiling(v T) (T, bool) {
	n := ceiling(ms.root, v, ms.less)
	if n == nil {
		var zero T
		return zero, false
	}
	return n.key, true
}

// Rank returns the number of elements strictly less than v. O(log n) expected.
func (ms *Multiset[T]) Rank(v T) int {
	return rank(ms.root, v, ms.less)
}

// Kth returns the k-th smallest element (0-indexed, duplicates counted separately).
// Returns the zero value and false if k is out of range. O(log n) expected.
func (ms *Multiset[T]) Kth(k int) (T, bool) {
	if k < 0 || k >= sz(ms.root) {
		var zero T
		return zero, false
	}
	n := kth(ms.root, k)
	if n == nil {
		var zero T
		return zero, false
	}
	return n.key, true
}

// ToSlice returns all elements in sorted order with duplicates. O(n).
func (ms *Multiset[T]) ToSlice() []T {
	out := make([]T, 0, sz(ms.root))
	toSlice(ms.root, &out)
	return out
}
