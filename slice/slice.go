package slice

// Map applies a transformation function to every element of a slice T and returns a new slice of V.
func Map[T, V any](s []T, transform func(t T) V) []V {
	res := make([]V, len(s))
	for i, v := range s {
		res[i] = transform(v)
	}
	return res
}

// Filter iterates over a slice T and returns a new slice containing only the elements
// for which the provided filter function returns true.
func Filter[T any](s []T, filter func(t T) bool) []T {
	res := make([]T, 0, len(s))
	for _, v := range s {
		if filter(v) {
			res = append(res, v)
		}
	}
	return res
}

// Reduce combines all elements of a slice T into a single accumulated value V,
// starting with an initial value.
func Reduce[T, V any](s []T, initial V, reduce func(acc V, curr T) V) V {
	res := initial
	for _, v := range s {
		res = reduce(res, v)
	}
	return res
}
