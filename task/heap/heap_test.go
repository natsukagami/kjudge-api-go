package heap

import (
	"math/rand"
	"sort"
	"testing"
)

type TestingTask struct {
	T int
	P int
}

type ByPriority []*TestingTask

func (b ByPriority) Len() int {
	return len(b)
}
func (b ByPriority) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b ByPriority) Less(i, j int) bool {
	return b[i].P > b[j].P
}

func TestHeap(t *testing.T) {
	const n = 10000
	slice := make([]*TestingTask, 0)
	perm := rand.Perm(n)
	H := MakeHeap()
	for i := 0; i < n; i++ {
		slice = append(slice, &TestingTask{T: i, P: perm[i]})
		H.Append(slice[i].T, slice[i].P)
	}
	sort.Sort(ByPriority(slice))
	for i := 0; i < n; i++ {
		item, err := H.Pop()
		if err != nil {
			t.Error("Unexpected error: " + err.Error())
			continue
		}
		t.Logf("Item popped: %d", item)
		if item != slice[i].T {
			t.Errorf("Wrong item popped, expected %d (P = %d) got %d", slice[i].T, slice[i].P, item)
		}
	}
	item, err := H.Pop()
	if _, ok := err.(EmptyError); item != nil || !ok {
		t.Error("Does not emit error when EmptyError is expected")
	}
	if e := err.(EmptyError); e.Error() != "Popping a heap when it is empty" {
		t.Error("Wrong error message?")
	}
}

func BenchmarkHeapPush(b *testing.B) {
	var n = b.N
	perm := rand.Perm(n)
	H := MakeHeap()
	for i := 0; i < n; i++ {
		H.Append(i, perm[i])
	}
}

func BenchmarkHeapPop(b *testing.B) {
	var n = b.N
	perm := rand.Perm(n)
	H := MakeHeap()
	for i := 0; i < n; i++ {
		H.Append(i, perm[i])
	}
	b.ResetTimer()
	for i := 0; i < n; i++ {
		H.Pop()
	}
}
