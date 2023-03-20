// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ring

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func genRangeSlice(start, end int) []int {
	out := make([]int, end-start)
	for i := 0; i < (end - start); i++ {
		out[i] = start + i
	}
	return out
}

func reverseSlice[T any](sl []T) []T {
	out := make([]T, len(sl))
	for idx, v := range sl {
		out[len(sl)-1-idx] = v
	}
	return out
}

func TestExampleDo(t *testing.T) {
	// Create a new ring of size 5
	r := New[int](5)

	// Get the length of the ring
	n := r.Len()

	// Initialize the ring with some integer values
	for i := 0; i < n; i++ {
		r.Set(i)
		r = r.Next()
	}

	output := []int{}
	// Iterate through the ring and print its contents
	r.Do(func(v int, hasValue bool) {
		output = append(output, v)
	})

	expected := genRangeSlice(0, 5)

	require.Condition(
		t,
		func() (success bool) { return cmp.Equal(expected, ([]int)(output)) },
		cmp.Diff(expected, ([]int)(output)),
	)
}

func TestExampleLen(t *testing.T) {
	// Create a new ring of size 4
	r := New[int](4)

	if l := r.Len(); l != 4 {
		t.Fatalf("not equal. expected = %+v, actual = %+v", 4, l)
	}
}

func TestExampleLink(t *testing.T) {
	// Create two rings, r and s, of size 2
	r := New[int](2)
	s := New[int](2)

	// Get the length of the ring
	lr := r.Len()
	ls := s.Len()

	// Initialize r with 0s
	for i := 0; i < lr; i++ {
		r.Set(0)
		r = r.Next()
	}

	// Initialize s with 1s
	for j := 0; j < ls; j++ {
		s.Set(1)
		s = s.Next()
	}

	// Link ring r and ring s
	rs := r.Link(s)

	output := []int{}
	// Iterate through the combined ring and print its contents
	rs.Do(func(v int, hasValue bool) {
		output = append(output, v)
	})

	expected := []int{0, 0, 1, 1}
	require.Condition(
		t,
		func() (success bool) { return cmp.Equal(expected, ([]int)(output)) },
		cmp.Diff(expected, ([]int)(output)),
	)
}

func TestExampleMove(t *testing.T) {
	// Create a new ring of size 5
	r := New[int](5)

	// Get the length of the ring
	n := r.Len()

	// Initialize the ring with some integer values
	for i := 0; i < n; i++ {
		r.Set(i)
		r = r.Next()
	}

	// Move the pointer forward by three steps
	r = r.Move(3)

	output := []int{}
	// Iterate through the combined ring and print its contents
	r.Do(func(v int, hasValue bool) {
		output = append(output, v)
	})

	expected := []int{3, 4, 0, 1, 2}
	require.Condition(
		t,
		func() (success bool) { return cmp.Equal(expected, ([]int)(output)) },
		cmp.Diff(expected, ([]int)(output)),
	)
}

func TestExampleNext(t *testing.T) {
	// Create a new ring of size 5
	r := New[int](5)

	// Get the length of the ring
	n := r.Len()

	// Initialize the ring with some integer values
	for i := 0; i < n; i++ {
		r.Set(i)
		r = r.Next()
	}

	output := []int{}
	// Iterate through the ring and print its contents
	for j := 0; j < n; j++ {
		v, _ := r.Get()
		output = append(output, v)
		r = r.Next()
	}

	expected := genRangeSlice(0, 5)
	require.Condition(
		t,
		func() (success bool) { return cmp.Equal(expected, ([]int)(output)) },
		cmp.Diff(expected, ([]int)(output)),
	)
}

func TestExamplePrev(t *testing.T) {
	// Create a new ring of size 5
	r := New[int](5)

	// Get the length of the ring
	n := r.Len()

	// Initialize the ring with some integer values
	for i := 0; i < n; i++ {
		r.Set(i)
		r = r.Next()
	}

	output := []int{}
	// Iterate through the ring backwards and print its contents
	for j := 0; j < n; j++ {
		r = r.Prev()
		v, _ := r.Get()
		output = append(output, v)
	}

	expected := reverseSlice(genRangeSlice(0, 5))
	require.Condition(
		t,
		func() (success bool) { return cmp.Equal(expected, ([]int)(output)) },
		cmp.Diff(expected, ([]int)(output)),
	)
}

func TestExampleUnlink(t *testing.T) {
	// Create a new ring of size 6
	r := New[int](6)

	// Get the length of the ring
	n := r.Len()

	// Initialize the ring with some integer values
	for i := 0; i < n; i++ {
		r.Set(i)
		r = r.Next()
	}

	// Unlink three elements from r, starting from r.Next()
	r.Unlink(3)

	output := []int{}
	// Iterate through the remaining ring and print its contents
	r.Do(func(v int, _ bool) {
		output = append(output, v)
	})

	expected := []int{0, 4, 5}
	require.Condition(
		t,
		func() (success bool) { return cmp.Equal(expected, ([]int)(output)) },
		cmp.Diff(expected, ([]int)(output)),
	)
}

func TestPointerConsistency(t *testing.T) {
	r := New[int](6)
	s := New[int](6)

	if len(r.entMap) != 6 {
		t.Fatalf("not equal. expected = %+v, actual = %+v", 6, len(r.entMap))
	}
	if len(s.entMap) != 6 {
		t.Fatalf("not equal. expected = %+v, actual = %+v", 6, len(s.entMap))
	}

	rpInit := r
	spInit := s

	// next
	rp := make([]*Ring[int], 6)
	for i := 0; i < 6; i++ {
		rp[i] = r
		r = r.Next()
	}
	sp := make([]*Ring[int], 6)
	for i := 0; i < 6; i++ {
		sp[i] = s
		s = s.Next()
	}

	if rpInit != r {
		t.Fatalf("not equal. expected = %p, actual = %p", rpInit, r)
	}
	if spInit != s {
		t.Fatalf("not equal. expected = %p, actual = %p", spInit, s)
	}

	rpAgain := make([]*Ring[int], 6)
	for i := 0; i < 6; i++ {
		rpAgain[i] = r
		r = r.Next()
	}
	spAgain := make([]*Ring[int], 6)
	for i := 0; i < 6; i++ {
		spAgain[i] = s
		s = s.Next()
	}

	// reflect.DeepEqual ignores index order. This ensures left and right are completely same slice.
	equal := func(left, right []*Ring[int]) bool {
		if len(left) != len(right) {
			return false
		}
		for i := 0; i < len(left); i++ {
			if left[i] != right[i] {
				return false
			}
		}
		return true
	}

	require.Condition(
		t,
		func() (success bool) {
			return equal(rp, rpAgain)
		},
	)
	require.Condition(
		t,
		func() (success bool) {
			return equal(sp, spAgain)
		},
	)

	// check if no ent is added.
	if len(r.entMap) != 6 {
		t.Fatalf("not equal. expected = %+v, actual = %+v", 6, len(r.entMap))
	}
	if len(s.entMap) != 6 {
		t.Fatalf("not equal. expected = %+v, actual = %+v", 6, len(s.entMap))
	}

	// prev
	rpAgain = make([]*Ring[int], 6)
	for i := 0; i < 6; i++ {
		// making sure reversed order.
		r = r.Prev()
		rpAgain[i] = r
	}
	spAgain = make([]*Ring[int], 6)
	for i := 0; i < 6; i++ {
		s = s.Prev()
		spAgain[i] = s
	}

	var rev []*Ring[int]
	rev = reverseSlice(rpAgain)
	for i := 0; i < len(rev); i++ {
		if rp[i] != rev[i] {
			t.Fatalf("not equal. expected = %p, actual = %p", rp[i], rev[i])
		}
	}
	rev = reverseSlice(spAgain)
	for i := 0; i < len(rev); i++ {
		if sp[i] != rev[i] {
			t.Fatalf("not equal. expected = %p, actual = %p", sp[i], rev[i])
		}
	}

	// Link
	// adding at tail.
	r.Prev().Link(s)

	var expected []*Ring[int]
	expected = append(expected, rp...)
	expected = append(expected, sp...)
	for i := 0; i < 12; i++ {
		if expected[i] != r {
			t.Fatalf("not equal. at index %d, expected = %p, actual = %p", i, expected[i], r)
		}
		r = r.Next()
	}

	if r != rpInit {
		t.Fatal("unknown behavior")
	}

	// Unlink
	r.Next().Next().Unlink(6)

	expected = []*Ring[int]{rp[0], rp[1], rp[2], sp[3], sp[4], sp[5]}
	if len(expected) != r.Len() {
		t.Fatalf("not equal. expected = %d, actual = %d", len(expected), r.Len())
	}
	for i := 0; i < 12-6; i++ {
		if expected[i] != r {
			t.Fatalf("not equal. at index %d, expected = %p, actual = %p", i, expected[i], r)
		}
		r = r.Next()
	}
}
