package slice

import (
	"iter"
)

// SliceIterator is a type definition of an iter.Seq for value of type T
type SliceIterator[T any] iter.Seq[T]

// NewIterator creates a SliceIterator that iterates over the elements of the
// provided slice s.
func NewIterator[T any](s []T) SliceIterator[T] {
	return func(yield func(v T) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

type Enumerated[U any] struct {
	Index int
	Value U
}
type Enumerator[V any] iter.Seq[Enumerated[V]]

func (s SliceIterator[T]) Enumerate() Enumerator[T] {
	return func(yield func(v Enumerated[T]) bool) {
		i := 0
		for v := range s {
			if !yield(Enumerated[T]{i, v}) {
				return
			}
			i++
		}
	}
}

// Returns count of elements in iterator
func (s SliceIterator[T]) Count() int {
	cnt := 0
	for range s {
		cnt++
	}
	return cnt
}

// Collect consumes the SliceIterator and returns a new slice containing all
// the yielded elements.
func (i SliceIterator[T]) Collect() []T {
	out := make([]T, i.Count())
	for e := range i.Enumerate() {
		out[e.Index] = e.Value
	}
	return out
}

// Map applies a transformation function to every element of a slice T and returns a new iterator of V.
func Map[T, V any](s SliceIterator[T], transform func(t T) V) SliceIterator[V] {
	return func(yield func(v V) bool) {
		for v := range s {
			if !yield(transform(v)) {
				return
			}
		}
	}
}

// Filter iterates over a slice T and returns a new iterator containing only the elements
// for which the provided filter function returns true.
func Filter[T any](s SliceIterator[T], filter func(t T) bool) SliceIterator[T] {
	return func(yield func(v T) bool) {
		for v := range s {
			if filter(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Reduce combines all elements of a iterator of type T into a single accumulated value V,
// starting with an initial value
func Reduce[T, V any](i SliceIterator[T], initial V, reduce func(acc V, curr T) V) V {
	res := initial
	for v := range i {
		res = reduce(res, v)
	}
	return res
}

type Tuple[T, V any] struct {
	Left  T
	Right V
}

// Combines two slices and returns an iterator of tuples
func Zip[T, V any](s1 []T, s2 []V) SliceIterator[Tuple[T, V]] {
	return func(yield func(Tuple[T, V]) bool) {
		if len(s1) != len(s2) {
			return
		}
		for i, v1 := range s1 {
			if !yield(Tuple[T, V]{v1, s2[i]}) {
				return
			}
		}
	}
}

// Concatenates multiple iterators into a single iterator
func Concat[T any](iters ...SliceIterator[T]) SliceIterator[T] {
	return func(yield func(T) bool) {
		for _, it := range iters {
			for v := range it {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Returns true if any elements meet predicate
func Any[T any](s SliceIterator[T], predicate func(v T) bool) bool {
	for el := range s {
		if predicate(el) {
			return true
		}
	}
	return false
}

// Returns true if all elements meet predicate
func All[T any](s SliceIterator[T], predicate func(v T) bool) bool {
	for el := range s {
		if !predicate(el) {
			return false
		}
	}
	return true
}
