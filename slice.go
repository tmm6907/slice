package slice

import "iter"

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

// Collect consumes the SliceIterator and returns a new slice containing all
// the yielded elements.
func (i SliceIterator[T]) Collect() []T {
	var res []T
	for v := range i {
		res = append(res, v)
	}
	return res
}

// Map applies a transformation function to every element of a slice T and returns a new iterator of V.
func Map[T, V any](s []T, transform func(t T) V) SliceIterator[V] {
	return func(yield func(v V) bool) {
		for _, v := range s {
			if !yield(transform(v)) {
				return
			}
		}
	}
}

// Filter iterates over a slice T and returns a new iterator containing only the elements
// for which the provided filter function returns true.
func Filter[T any](s []T, filter func(t T) bool) SliceIterator[T] {
	return func(yield func(v T) bool) {
		for _, v := range s {
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
