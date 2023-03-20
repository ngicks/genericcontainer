package heapimpl

import (
	"testing"

	"github.com/ngicks/generic/slice"
)

var _ Lesser[*item] = ((*item)(nil))
var _ Swapper[*item] = ((*item)(nil))
var _ Popper[*item] = ((*item)(nil))
var _ Pusher[*item] = ((*item)(nil))

type item struct {
	value string
	index int
}

func (item *item) Less(i, j *item) bool {
	return i.value < j.value
}

func (it *item) Swap(slice *slice.Stack[*item], i, j int) {
	(*slice)[i], (*slice)[j] = (*slice)[j], (*slice)[i]
	(*slice)[i].index = i
	(*slice)[j].index = j
}

func (it *item) Push(slice *slice.Stack[*item], v *item) {
	v.index = slice.Len()
	slice.Push(v)
}

func (it *item) Pop(slice *slice.Stack[*item]) *item {
	popped, ok := slice.Pop()
	if !ok {
		_ = (*(*[]*item)(slice))[slice.Len()-1]
	}
	popped.index = -1
	return popped
}

func TestFilterableHeap_swap_is_used_if_implemented(t *testing.T) {
	h := NewFilterableHeap[*item]()

	h.Push(&item{value: "foo"})
	h.Push(&item{value: "bar"})
	h.Push(&item{value: "baz"})
	h.Push(&item{value: "qux"})
	h.Push(&item{value: "quux"})

	slice := ([]*item)(h.internal.Inner)

	if ele := slice[0]; ele.index != 0 || ele.value != "bar" {
		t.Fatalf("incorrect: %+v", *ele)
	}
	if ele := slice[1]; ele.index != 1 || ele.value != "foo" {
		t.Fatalf("incorrect: %+v", *ele)
	}
	if ele := slice[2]; ele.index != 2 || ele.value != "baz" {
		t.Fatalf("incorrect: %+v", *ele)
	}
	if ele := slice[3]; ele.index != 3 || ele.value != "qux" {
		t.Fatalf("incorrect: %+v", *ele)
	}
	if ele := slice[4]; ele.index != 4 || ele.value != "quux" {
		t.Fatalf("incorrect: %+v", *ele)
	}

	for i := 0; i < 5; i++ {
		if item := h.Pop(); item.index != -1 {
			t.Fatalf("incorrect: %+v", *item)
		}
	}
}
