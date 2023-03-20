package list

import (
	"container/list"
)

type entMap[T any] map[*list.Element]*Element[T]

// getOrCreate gets e[l] or creates *Element[T].
//
//   - returns nil if l is nil.
//   - creates new *Element[T] if e does not have l.
//     this happens only after calling PushBackList / PushFrontList.
func (e entMap[T]) getOrCreate(l *list.Element) *Element[T] {
	if l == nil {
		return nil
	}
	ele, ok := e[l]
	if !ok {
		return e.createInsertingElement(l)
	}
	return ele
}

// createInsertingElement creates *Element[T] and insert into e,
// returning that *Element[T].
//
//   - returns nil if l is nil.
func (e entMap[T]) createInsertingElement(l *list.Element) *Element[T] {
	if l == nil {
		return nil
	}
	ent := &Element[T]{
		inner:  l,
		entMap: e,
	}
	e[l] = ent
	return ent
}
