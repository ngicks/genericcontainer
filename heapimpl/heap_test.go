package heapimpl_test

import (
	"testing"
	"time"

	"github.com/ngicks/genericcontainer/heap"
	"github.com/ngicks/genericcontainer/heapimpl"
)

var _ heap.Interface[int] = heapimpl.NewSliceInterface(nil, nil, heapimpl.HeapMethods[int]{})

func TestSimpleHeap(t *testing.T) {
	// Seeing basic delegation.
	t.Run("int heap", func(t *testing.T) {
		h, inter := heapimpl.MakeMinHeap[int]()
		ans := []int{3, 4, 4, 5, 6}
		h.Push(5)
		h.Push(4)
		h.Push(6)
		h.Push(3)
		h.Push(4)

		for _, i := range ans {
			popped := h.Pop()
			if popped != i {
				t.Errorf("pop = %v expected %v", popped, i)
			}
		}
		if inter.Len() != 0 {
			t.Errorf("expect empty but size = %v", inter.Len())
		}
	})

	t.Run("struct heap", func(t *testing.T) {
		type testStruct struct {
			t time.Time
		}
		less := func(i, j *testStruct) bool {
			return i.t.Before(j.t)
		}

		h, inter := heapimpl.MakeHeap(less, heapimpl.HeapMethods[*testStruct]{})
		ans := []*testStruct{
			{t: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
			{t: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)},
			{t: time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)},
			{t: time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC)},
			{t: time.Date(2021, 5, 1, 0, 0, 0, 0, time.UTC)},
		}
		h.Push(ans[2])
		h.Push(ans[1])
		h.Push(ans[3])
		h.Push(ans[0])
		h.Push(ans[4])

		for _, i := range ans {
			popped := h.Pop()
			if popped.t != i.t {
				t.Errorf("pop = %v expected %v", popped.t, i.t)
			}
		}
		if inter.Len() != 0 {
			t.Errorf("expect empty but size = %v", inter.Len())
		}
	})
}
