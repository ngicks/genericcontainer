// package list wraps doubly-linked list
// implemented in go programming language standard library,
// to accept type parameter.
package list

import "container/list"

type Element[T any] struct {
	inner  *list.Element
	entMap entMap[T]
}

func (e *Element[T]) ensureValid() {
	if e.entMap == nil {
		panic("not initialized: zero *Element[T] is invalid. e.entMap is nil")
	}
	if e.inner == nil {
		panic("not initialized: zero *Element[T] is invalid. e.inner is nil")
	}
}

// Unwrap returns internal *`container/list`.Element.
// Setting non-T value may cause runtime panic in succeeding Get call.
func (e *Element[T]) Unwrap() *list.Element {
	e.ensureValid()
	return e.inner
}

// Get returns internal Value. If internal Value is non-nil and then returns value and true.
// Otherwise returns zero of T and false.
func (e *Element[T]) Get() (v T, ok bool) {
	e.ensureValid()

	if e.inner.Value == nil {
		return
	} else {
		return e.inner.Value.(T), true
	}
}

// Set is equivalent to `element.Value = v`
func (e *Element[T]) Set(v T) {
	e.ensureValid()
	e.inner.Value = v
}

func (e *Element[T]) Next() *Element[T] {
	e.ensureValid()
	return e.entMap.getOrCreate(e.inner.Next())
}
func (e *Element[T]) Prev() *Element[T] {
	e.ensureValid()
	return e.entMap.getOrCreate(e.inner.Prev())
}

// List[T] is `container/list` wrapper that is safe to use with type T.
//
// It holds *`container/list`.List and one extra map
// (map[*list.Element]*Element[T]) to maintain pointer consistency.
// The map remembers pointer values it previously returned so a returned pointer is always same
// pointer for same *`container/list`.Element.
//
// List[T] tries best to be consistent with `container/list`.
// The zero value of List[T] is a valid empty list (which will be lazily initialized).
type List[T any] struct {
	inner  *list.List
	entMap entMap[T]
}

func New[T any]() *List[T] {
	return &List[T]{
		inner:  list.New(),
		entMap: make(entMap[T]),
	}
}

func (l *List[T]) lazyInit() {
	if l.inner == nil {
		l.inner = list.New()
	}
	if l.entMap == nil {
		l.entMap = make(entMap[T])
	}
}

func (l *List[T]) Back() *Element[T] {
	l.lazyInit()
	return l.entMap.getOrCreate(l.inner.Back())
}
func (l *List[T]) Front() *Element[T] {
	l.lazyInit()
	return l.entMap.getOrCreate(l.inner.Front())
}
func (l *List[T]) Init() *List[T] {
	l.lazyInit()
	l.inner.Init()
	l.entMap = make(entMap[T])
	return l
}
func (l *List[T]) InsertAfter(v T, mark *Element[T]) *Element[T] {
	l.lazyInit()
	return l.entMap.createInsertingElement(l.inner.InsertAfter(v, mark.inner))
}
func (l *List[T]) InsertBefore(v T, mark *Element[T]) *Element[T] {
	l.lazyInit()
	return l.entMap.createInsertingElement(l.inner.InsertBefore(v, mark.inner))
}
func (l *List[T]) Len() int {
	l.lazyInit()
	return l.inner.Len()
}
func (l *List[T]) MoveAfter(e, mark *Element[T]) {
	l.lazyInit()
	l.inner.MoveAfter(e.inner, mark.inner)
}
func (l *List[T]) MoveBefore(e, mark *Element[T]) {
	l.lazyInit()
	l.inner.MoveBefore(e.inner, mark.inner)
}
func (l *List[T]) MoveToBack(e *Element[T]) {
	l.lazyInit()
	l.inner.MoveToBack(e.inner)
}
func (l *List[T]) MoveToFront(e *Element[T]) {
	l.lazyInit()
	l.inner.MoveToFront(e.inner)
}

// PushBack inserts v at end of the list
// and returns inserted element.
func (l *List[T]) PushBack(v T) *Element[T] {
	l.lazyInit()
	return l.entMap.createInsertingElement(l.inner.PushBack(v))
}

// PushBackList copies values of other and inserts them into back of l.
//
// Both l and other must not be nil.
func (l *List[T]) PushBackList(other *List[T]) {
	l.lazyInit()
	other.lazyInit()
	l.inner.PushBackList(other.inner)
}

// PushFront inserts v at front of the list
// and returns inserted element.
func (l *List[T]) PushFront(v T) *Element[T] {
	l.lazyInit()
	return l.entMap.createInsertingElement(l.inner.PushFront(v))
}

// PushFrontList copies values of other and inserts them into front of l.
//
// Both l and other must not be nil.
func (l *List[T]) PushFrontList(other *List[T]) {
	l.lazyInit()
	other.lazyInit()
	l.inner.PushFrontList(other.inner)
}

// Remove removes element from list l
// and returns value of e.
// hadValue is false when e.Value is nil.
//
// e must not be nil.
func (l *List[T]) Remove(e *Element[T]) (v T, hadValue bool) {
	l.lazyInit()
	delete(l.entMap, e.inner)
	if ele := l.inner.Remove(e.inner); ele != nil {
		return ele.(T), true
	} else {
		return
	}
}
