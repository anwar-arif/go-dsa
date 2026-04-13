package main

import (
	"fmt"

	"github.com/anwar-arif/go-dsa/priorityqueue"
)

func main() {
	nums := []int{3, 3, 1, 1, 2, 2}
	pq := priorityqueue.NewFromSlice(func(a, b int) bool {
		return a < b
	}, nums)

	for !pq.IsEmpty() {
		fmt.Println(pq.Pop())
	}
}
