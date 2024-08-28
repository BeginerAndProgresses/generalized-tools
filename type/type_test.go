package gttype

import "testing"

func TestNewHeap(t *testing.T) {
	heap := NewHeap[any]()
	heap.Insert(1)
	heap.Insert(2)
	heap.PrintHeap()
}
