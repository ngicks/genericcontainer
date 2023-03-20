package ring

import (
	"container/ring"
	"reflect"
)

type Ring[T any] struct {
	entMap entMap[T]
	inner  *ring.Ring
}

func New[T any](n int) *Ring[T] {
	r := ring.New(n)
	if r == nil {
		return nil
	}
	entMap := make(entMap[T])
	for i := 0; i < n; i++ {
		entMap.getOrCreate(r)
		r = r.Next()
	}
	return entMap.getOrCreate(r)
}

func (r *Ring[T]) lazyInit() {
	if r.inner == nil {
		r.inner = ring.New(1)
		r.entMap = make(entMap[T])
		r.entMap[r.inner] = r
	}
}

// Unwrap returns internal *`container/ring`.Ring.
// Setting non-T value may cause runtime panic in succeeding Get or Do call.
func (r *Ring[T]) Unwrap() *ring.Ring {
	r.lazyInit()
	return r.inner
}

// Get returns internal Value. If internal Value is non-nil and then returns value and true.
// Otherwise returns zero of T and false.
func (r *Ring[T]) Get() (v T, ok bool) {
	r.lazyInit()
	if r.inner.Value == nil {
		return
	}
	return r.inner.Value.(T), true
}

// Set sets value to inner `ring.Value`
// and returns that Ring.
func (r *Ring[T]) Set(v T) *Ring[T] {
	r.lazyInit()
	r.inner.Value = v
	return r
}

// Do is equivalent of `container/ring`.Do but added hasValue boolean.
// hasValue is false if internal Value is nil, indicating passed value is zero of T. hasValue is true otherwise.
func (r *Ring[T]) Do(f func(v T, hasValue bool)) {
	r.lazyInit()
	r.inner.Do(func(a any) {
		if a == nil {
			var zero T
			f(zero, false)
		} else {
			f(a.(T), true)
		}
	})
}
func (r *Ring[T]) Len() int {
	r.lazyInit()
	return r.inner.Len()
}
func (r *Ring[T]) Link(s *Ring[T]) *Ring[T] {
	r.lazyInit()
	if reflect.ValueOf(r.entMap).Pointer() != reflect.ValueOf(s.entMap).Pointer() {
		r.entMap.merge(s.entMap)
		s.entMap.merge(r.entMap)
	}
	return r.entMap.getOrCreate(r.inner.Link(s.inner))
}
func (r *Ring[T]) Move(n int) *Ring[T] {
	r.lazyInit()
	return r.entMap.getOrCreate(r.inner.Move(n))
}
func (r *Ring[T]) Next() *Ring[T] {
	r.lazyInit()
	return r.entMap.getOrCreate(r.inner.Next())
}
func (r *Ring[T]) Prev() *Ring[T] {
	r.lazyInit()
	return r.entMap.getOrCreate(r.inner.Prev())
}
func (r *Ring[T]) Unlink(n int) *Ring[T] {
	r.lazyInit()
	return r.entMap.getOrCreate(r.inner.Unlink(n))
}
