// Package heap provides an implementation of a Task heap
package heap

import "container/heap"

// EmptyError is an error that is thrown when an attempt to pop an empty heap is made.
type EmptyError struct {
}

// Error implements the error interface
func (h EmptyError) Error() string {
	return "Popping a heap when it is empty"
}

// Item is a member of Heap
type item struct {
	T interface{}
	P int
}

// heapSlice provides an interface to heap library.
type heapSlice []*item

// Len returns the length of the Heap. (heap interface)
func (h heapSlice) Len() int { return len(h) }

// Swap swaps two items in the heap. (heap interface)
func (h heapSlice) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// Less compares two items in the heap. (heap interface)
// It returns true if the former has higher priority than the latter.
func (h heapSlice) Less(i, j int) bool { return h[i].P > h[j].P }

// Push adds a new item into the end of a heap. (heap interface)
func (h *heapSlice) Push(i interface{}) {
	*h = append(*h, i.(*item))
}

// Pop removes an item from the end of the heap. (heap interface)
func (h *heapSlice) Pop() interface{} {
	old := *h
	i := old[len(old)-1]
	*h = old[:len(old)-1]
	return i
}

// Heap is a data structure containing Tasks with priority.
// It supports queuing and querying the largest priority task in logarithmic time.
type Heap struct {
	slice heapSlice
}

// Len returns the length of the Heap
func (h Heap) Len() int {
	return h.slice.Len()
}

// Append adds a new item into the heap
func (h *Heap) Append(T interface{}, P int) {
	heap.Push(&h.slice, &item{T, P})
}

// Pop removes an item from the heap and returns it
func (h *Heap) Pop() (T interface{}, err error) {
	if h.Len() == 0 {
		err = EmptyError{}
		return
	}
	t := heap.Pop(&h.slice).(*item)
	T = t.T
	return
}

// MakeHeap creates a heap
func MakeHeap() Heap {
	h := Heap{make(heapSlice, 0)}
	heap.Init(&h.slice)
	return h
}
