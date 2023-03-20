// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package list

import (
	"container/list"
	"testing"
)

func checkListLen[T any](t *testing.T, l *List[T], len int) bool {
	if n := l.Len(); n != len {
		t.Errorf("l.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func TestList(t *testing.T) {
	l := New[any]()

	// Single element list
	e := l.PushFront("a")
	l.MoveToFront(e)
	l.MoveToBack(e)
	l.Remove(e)
	// Bigger list
	e2 := l.PushFront(2)
	e1 := l.PushFront(1)
	e3 := l.PushBack(3)
	e4 := l.PushBack("banana")

	l.Remove(e2)

	l.MoveToFront(e3) // move from middle

	l.MoveToFront(e1)
	l.MoveToBack(e3) // move from middle

	l.MoveToFront(e3) // move from back
	l.MoveToFront(e3) // should be no-op

	l.MoveToBack(e3) // move from front
	l.MoveToBack(e3) // should be no-op

	e2 = l.InsertBefore(2, e1) // insert before front
	l.Remove(e2)
	e2 = l.InsertBefore(2, e4) // insert before middle
	l.Remove(e2)
	e2 = l.InsertBefore(2, e3) // insert before back
	l.Remove(e2)

	e2 = l.InsertAfter(2, e1) // insert after front
	l.Remove(e2)
	e2 = l.InsertAfter(2, e4) // insert after middle
	l.Remove(e2)
	e2 = l.InsertAfter(2, e3) // insert after back
	l.Remove(e2)

	// Check standard iteration.
	sum := 0
	for e := l.Front(); e != nil; e = e.Next() {
		got, _ := e.Get()
		if i, ok := got.(int); ok {
			sum += i
		}
	}
	if sum != 4 {
		t.Errorf("sum over l = %d, want 4", sum)
	}

	// Clear all elements by iterating
	var next *Element[any]
	for e := l.Front(); e != nil; e = next {
		next = e.Next()
		l.Remove(e)
	}
}

func checkList(t *testing.T, l *List[int], es []int) {
	if !checkListLen(t, l, len(es)) {
		return
	}

	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		le, _ := e.Get()
		if le != es[i] {
			t.Errorf("elt[%d].Value = %v, want %v", i, le, es[i])
		}
		i++
	}
}

func TestExtending(t *testing.T) {
	l1 := New[int]()
	l2 := New[int]()

	l1.PushBack(1)
	l1.PushBack(2)
	l1.PushBack(3)

	l2.PushBack(4)
	l2.PushBack(5)

	l3 := New[int]()
	l3.PushBackList(l1)
	checkList(t, l3, []int{1, 2, 3})
	l3.PushBackList(l2)
	checkList(t, l3, []int{1, 2, 3, 4, 5})

	l3 = New[int]()
	l3.PushFrontList(l2)
	checkList(t, l3, []int{4, 5})
	l3.PushFrontList(l1)
	checkList(t, l3, []int{1, 2, 3, 4, 5})

	checkList(t, l1, []int{1, 2, 3})
	checkList(t, l2, []int{4, 5})

	l3 = New[int]()
	l3.PushBackList(l1)
	checkList(t, l3, []int{1, 2, 3})
	l3.PushBackList(l3)
	checkList(t, l3, []int{1, 2, 3, 1, 2, 3})

	l3 = New[int]()
	l3.PushFrontList(l1)
	checkList(t, l3, []int{1, 2, 3})
	l3.PushFrontList(l3)
	checkList(t, l3, []int{1, 2, 3, 1, 2, 3})

	l3 = New[int]()
	l1.PushBackList(l3)
	checkList(t, l1, []int{1, 2, 3})
	l1.PushFrontList(l3)
	checkList(t, l1, []int{1, 2, 3})
}

func TestIssue4103(t *testing.T) {
	l1 := New[int]()
	l1.PushBack(1)
	l1.PushBack(2)

	l2 := New[int]()
	l2.PushBack(3)
	l2.PushBack(4)

	e := l1.Front()
	l2.Remove(e) // l2 should not change because e is not an element of l2
	if n := l2.Len(); n != 2 {
		t.Errorf("l2.Len() = %d, want 2", n)
	}

	l1.InsertBefore(8, e)
	if n := l1.Len(); n != 3 {
		t.Errorf("l1.Len() = %d, want 3", n)
	}
}

func TestIssue6349(t *testing.T) {
	l := New[int]()
	l.PushBack(1)
	l.PushBack(2)

	e := l.Front()
	l.Remove(e)
	if got, _ := e.Get(); got != 1 {
		t.Errorf("e.value = %d, want 1", got)
	}
	if e.Next() != nil {
		t.Errorf("e.Next() != nil")
	}
	if e.Prev() != nil {
		t.Errorf("e.Prev() != nil")
	}
}

// Test PushFront, PushBack, PushFrontList, PushBackList with uninitialized List
func TestZeroList(t *testing.T) {
	var l1 = New[int]()
	l1.PushFront(1)
	checkList(t, l1, []int{1})

	var l2 = New[int]()
	l2.PushBack(1)
	checkList(t, l2, []int{1})

	var l3 = New[int]()
	l3.PushFrontList(l1)
	checkList(t, l3, []int{1})

	var l4 = New[int]()
	l4.PushBackList(l2)
	checkList(t, l4, []int{1})
}

// Test that a list l is not modified when calling InsertBefore with a mark that is not an element of l.
func TestInsertBeforeUnknownMark(t *testing.T) {
	var l = New[int]()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.InsertBefore(1, &Element[int]{inner: new(list.Element), entMap: make(entMap[int])})
	checkList(t, l, []int{1, 2, 3})
}

// Test that a list l is not modified when calling InsertAfter with a mark that is not an element of l.
func TestInsertAfterUnknownMark(t *testing.T) {
	var l = New[int]()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.InsertAfter(1, &Element[int]{inner: new(list.Element), entMap: make(entMap[int])})
	checkList(t, l, []int{1, 2, 3})
}

// Test that a list l is not modified when calling MoveAfter or MoveBefore with a mark that is not an element of l.
func TestMoveUnknownMark(t *testing.T) {
	var l1 = New[int]()
	e1 := l1.PushBack(1)

	var l2 = New[int]()
	e2 := l2.PushBack(2)

	l1.MoveAfter(e1, e2)
	checkList(t, l1, []int{1})
	checkList(t, l2, []int{2})

	l1.MoveBefore(e1, e2)
	checkList(t, l1, []int{1})
	checkList(t, l2, []int{2})
}
