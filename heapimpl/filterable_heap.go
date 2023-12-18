package heapimpl

import "github.com/ngicks/generic/slice"

type Lesser[T any] interface {
	Less(i, j T) bool
}

type Swapper[T any] interface {
	Swap(slice *slice.Stack[T], i, j int)
}

type Pusher[T any] interface {
	Push(slice *slice.Stack[T], v T)
}

type Popper[T any] interface {
	Pop(slice *slice.Stack[T]) T
}

type FilterableHeap[T any] struct {
	*HeapWrapper[T]
	internal *SliceInterface[T]
}

// NewFilterableHeap returns newly created FilterableHeap.
// T must implement Lesser[T], otherwise it panics.
// T can optionally implement one or all of Swapper, Pusher and Popper.
// If T implements those interfaces, methods will be used in corresponding heap functions
// instead of default implementations.
func NewFilterableHeap[T any]() *FilterableHeap[T] {
	var less func(i, j T) bool
	var hooks HeapMethods[T]

	var zero T
	var asAny any = zero

	if lesser, ok := asAny.(Lesser[T]); ok {
		less = lesser.Less
	} else {
		panic("T must implements Lesser[T]")
	}

	if swapper, ok := asAny.(Swapper[T]); ok {
		hooks.Swap = swapper.Swap
	}
	if pusher, ok := asAny.(Pusher[T]); ok {
		hooks.Push = pusher.Push
	}
	if popper, ok := asAny.(Popper[T]); ok {
		hooks.Pop = popper.Pop
	}

	heapInternal, interfaceInternal := MakeHeap(less, hooks)
	return &FilterableHeap[T]{
		HeapWrapper: heapInternal,
		internal:    interfaceInternal,
	}
}

func NewFilterableHeapHooks[T any](less func(i, j T) bool, hooks HeapMethods[T]) *FilterableHeap[T] {
	heapInternal, interfaceInternal := MakeHeap(less, hooks)
	return &FilterableHeap[T]{
		HeapWrapper: heapInternal,
		internal:    interfaceInternal,
	}
}

// Peek returns most prioritized value in heap without removing it.
// If this heap contains 0 element, returned p is zero value for type T.
//
// The complexity is O(1).
func (h *FilterableHeap[T]) Peek() (p T) {
	if len(h.internal.Inner) == 0 {
		return
	}
	return h.internal.Inner[0]
}

func (h *FilterableHeap[T]) Len() int {
	return h.internal.Len()
}

// Clone clones internal slice and creates new FilterableHeap with it.
// This is done by simple slice copy, without succeeding Init call,
// which means it also clones broken invariants if any.
//
// If type T or one of its internal value is pointer type,
// mutation of T propagates cloned to original, and vice versa.
func (h *FilterableHeap[T]) Clone() *FilterableHeap[T] {
	cloned := make([]T, len(h.internal.Inner))
	copy(cloned, h.internal.Inner)

	n := NewFilterableHeapHooks[T](h.internal.less, h.internal.methods)
	n.internal.Inner = cloned
	return n
}

// Filter calls filterFuncs one by one, in given order, with following heap.Init().
// Each filterFunc receives innerSlice so that it can be mutated in that func.
//
// Filter calls heap.Init() at the end of the method.
// So the complexity is at least O(n) where n is h.Len().
func (h *FilterableHeap[T]) Filter(filterFuncs ...func(innerSlice *[]T)) {
	for _, v := range filterFuncs {
		if v != nil {
			v((*[]T)(&h.internal.Inner))
		}
	}

	h.Init()
}

func BuildExcludeFilter[T any](start, end int, filterPredicate func(T) bool) func(sl *[]T) {
	return func(sl *[]T) {
		if start < 0 {
			start = 0
		} else if start >= len(*sl) {
			return
		}
		if end > len(*sl) {
			end = len(*sl)
		}

		if start > end {
			return
		}

		for i := start; i < end; i++ {
			if filterPredicate((*sl)[i]) {
				*sl = append((*sl)[:i], (*sl)[i+1:]...)
				end--
				i--
			}
		}
	}
}
