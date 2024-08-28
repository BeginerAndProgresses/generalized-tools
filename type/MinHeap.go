package gttype

import "fmt"

// Comparator 是一种函数类型，用于比较 T 类型的两个元素
type Comparator[T any] func(a, b T) bool

// Heap 表示通用的 Min-Heap 结构
type Heap[T any] struct {
	arr        []T
	comparator Comparator[T]
}

// NewHeap 使用提供的 Comparator 函数创建新的泛型堆
func NewHeap[T any](comparator Comparator[T]) *Heap[T] {
	return &Heap[T]{
		arr:        []T{},
		comparator: comparator,
	}
}

// Insert 添加一个键到堆中
func (h *Heap[T]) Insert(key T) {
	h.arr = append(h.arr, key)
	h.heapifyUp(len(h.arr) - 1)
}

// ExtractMin 移除堆中的最小值并返回
func (h *Heap[T]) ExtractMin() T {
	if len(h.arr) == 0 {
		var zero T
		fmt.Println("Heap is empty")
		return zero
	}

	min := h.arr[0]
	h.arr[0] = h.arr[len(h.arr)-1]
	h.arr = h.arr[:len(h.arr)-1]

	h.heapifyDown(0)
	return min
}

// heapifyUp 插入元素時調整堆的結構
func (h *Heap[T]) heapifyUp(index int) {
	for h.comparator(h.arr[index], h.arr[parent(index)]) {
		h.swap(parent(index), index)
		index = parent(index)
	}
}

// heapifyDown 取出後來調整堆的結構
func (h *Heap[T]) heapifyDown(index int) {
	smallest := index
	left := leftChild(index)
	right := rightChild(index)

	if left < len(h.arr) && h.comparator(h.arr[left], h.arr[smallest]) {
		smallest = left
	}
	if right < len(h.arr) && h.comparator(h.arr[right], h.arr[smallest]) {
		smallest = right
	}

	if smallest != index {
		h.swap(index, smallest)
		h.heapifyDown(smallest)
	}
}

// swap 改變兩個元素的位置
func (h *Heap[T]) swap(i, j int) {
	h.arr[i], h.arr[j] = h.arr[j], h.arr[i]
}

// parent 返回父節點的索引
func parent(index int) int {
	return (index - 1) / 2
}

// leftChild 返回左子節點的索引
func leftChild(index int) int {
	return 2*index + 1
}

// rightChild 返回右子節點的索引
func rightChild(index int) int {
	return 2*index + 2
}

// PrintHeap 打印堆
func (h *Heap[T]) PrintHeap() {
	fmt.Println(h.arr)
}
