package heapimpl

import (
	"github.com/ngicks/generic/slice"
	"github.com/ngicks/genericcontainer/heap"
	"golang.org/x/exp/constraints"
)

// HeapMethods is replacement methods for heap.
// Non nil fields will be used in heap functions of corresponding name.
type HeapMethods[T any] struct {
	Swap func(slice *slice.Stack[T], i, j int)
	Push func(slice *slice.Stack[T], v T)
	Pop  func(slice *slice.Stack[T]) T
}

// MakeHeap makes a heap for the type T using a less[T] function.
//
// To add your own heap methods, embed *HeapWrapper[T] to your own struct type
// to expose basic heap methods,
// and manipulate SliceInterface[T].Inner slice in its own way.
// Mutating inner slice may need succeeding Init() or Fix() call.
func MakeHeap[T any](less func(i, j T) bool, methods HeapMethods[T]) (*HeapWrapper[T], *SliceInterface[T]) {
	internal := NewSliceInterface(nil, less, methods)
	return NewHeapWrapper[T](internal), internal
}

func less[T constraints.Ordered](i, j T) bool {
	return i < j
}

func more[T constraints.Ordered](i, j T) bool {
	return i > j
}

// MakeMinHeap makes a MinHeap for the type T.
//
// MakeMinHeap does what MakeHeap does but with pre-declared less function.
// T is constrained to pre-declared primitive types which are compatible with less and greater comparison operations.
func MakeMinHeap[T constraints.Ordered]() (*HeapWrapper[T], *SliceInterface[T]) {
	internal := NewSliceInterface(nil, less[T], HeapMethods[T]{})
	return NewHeapWrapper[T](internal), internal
}

// MakeMaxHeap makes a MaxHeap for the type T.
// This is same as MakeMinHeap but uses more[T] instead.
func MakeMaxHeap[T constraints.Ordered]() (*HeapWrapper[T], *SliceInterface[T]) {
	internal := NewSliceInterface(nil, more[T], HeapMethods[T]{})
	return NewHeapWrapper[T](internal), internal
}

type HeapWrapper[T any] struct {
	inter heap.Interface[T]
}

func NewHeapWrapper[T any](inter heap.Interface[T]) *HeapWrapper[T] {
	return &HeapWrapper[T]{
		inter: inter,
	}
}

func (hw *HeapWrapper[T]) Fix(i int) {
	heap.Fix(hw.inter, i)
}
func (hw *HeapWrapper[T]) Init() {
	heap.Init(hw.inter)
}
func (hw *HeapWrapper[T]) Pop() T {
	return heap.Pop(hw.inter)
}
func (hw *HeapWrapper[T]) Push(x T) {
	heap.Push(hw.inter, x)
}
func (hw *HeapWrapper[T]) Remove(i int) T {
	return heap.Remove(hw.inter, i)
}

type SliceInterface[T any] struct {
	Inner   slice.Stack[T]
	less    func(i, j T) bool
	methods HeapMethods[T]
}

// NewSliceInterface returns a newly created SliceInterface.
// less is mandatory. Any or all fields of HeapMethods can be nil.
func NewSliceInterface[T any](
	init []T,
	less func(i, j T) bool,
	methods HeapMethods[T],
) *SliceInterface[T] {
	if init == nil {
		init = make([]T, 0)
	}
	return &SliceInterface[T]{
		Inner:   init,
		less:    less,
		methods: methods,
	}
}

func (sl *SliceInterface[T]) Len() int {
	return len(sl.Inner)
}
func (sl *SliceInterface[T]) Less(i, j int) bool {
	return sl.less(sl.Inner[i], sl.Inner[j])
}
func (sl *SliceInterface[T]) Swap(i, j int) {
	if sl.methods.Swap != nil {
		sl.methods.Swap(&sl.Inner, i, j)
	} else {
		sl.Inner[i], sl.Inner[j] = sl.Inner[j], sl.Inner[i]
	}
}
func (sl *SliceInterface[T]) Push(x T) {
	if sl.methods.Push != nil {
		sl.methods.Push(&sl.Inner, x)
	} else {
		sl.Inner.Push(x)
	}
}
func (sl *SliceInterface[T]) Pop() T {
	if sl.methods.Pop != nil {
		return sl.methods.Pop(&sl.Inner)
	} else {
		v, popped := sl.Inner.Pop()
		if !popped {
			// preserving the original error message.
			_ = sl.Inner[len(sl.Inner)-1]
		}
		return v
	}
}
