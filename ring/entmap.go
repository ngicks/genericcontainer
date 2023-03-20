package ring

import (
	"container/ring"
)

type entMap[T any] map[*ring.Ring]*Ring[T]

// getOrCreate gets e[l] or creates *Ring[T].
//
//   - returns nil if l is nil.
//   - creates new *Ring[T] if e does not have l.
//     this happens only after calling PushBackList / PushFrontList.
func (e entMap[T]) getOrCreate(r *ring.Ring) *Ring[T] {
	if r == nil {
		return nil
	}
	ele, ok := e[r]
	if !ok {
		return e.createInsertingRing(r)
	}
	return ele
}

func (e *entMap[T]) merge(other entMap[T]) {
	for k, v := range other {
		(*e)[k] = v
	}
}

// createInsertingRing creates *Ring[T] and insert into e,
// returning that *Ring[T].
//
//   - returns nil if l is nil.
func (e entMap[T]) createInsertingRing(r *ring.Ring) *Ring[T] {
	if r == nil {
		return nil
	}
	ent := &Ring[T]{
		inner:  r,
		entMap: e,
	}
	e[r] = ent
	return ent
}
