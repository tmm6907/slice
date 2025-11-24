package slice

import (
	"strconv"
	"testing"
)

func assertSlicesEqual[T comparable](t *testing.T, expected, actual []T, msg string) {
	if len(expected) != len(actual) {
		t.Fatalf("%s: Length mismatch. Expected %d, got %d. Expected: %v, Actual: %v", msg, len(expected), len(actual), expected, actual)
	}
	for i := range expected {
		if expected[i] != actual[i] {
			t.Fatalf("%s: Value mismatch at index %d. Expected %v, got %v. Expected: %v, Actual: %v", msg, i, expected[i], actual[i], expected, actual)
		}
	}
}

func TestMap(t *testing.T) {
	intSlice := []int{1, 2, 3}

	t.Run("IntToString", func(t *testing.T) {
		mapper := func(i int) string {
			return "Num-" + strconv.Itoa(i)
		}
		expected := []string{"Num-1", "Num-2", "Num-3"}

		iterator := Map(NewIterator(intSlice), mapper)
		actual := iterator.Collect()
		assertSlicesEqual(t, expected, actual, "Map int to string failed")
	})

	t.Run("IntToIntDoubling", func(t *testing.T) {
		mapper := func(i int) int {
			return i * 2
		}
		expected := []int{2, 4, 6}
		iterator := Map(NewIterator(intSlice), mapper)
		actual := iterator.Collect()
		assertSlicesEqual(t, expected, actual, "Map int to int failed")
	})

	t.Run("EmptySliceMap", func(t *testing.T) {
		emptySlice := []int{}
		mapper := func(i int) string { return strconv.Itoa(i) }
		expected := []string{}
		iterator := Map(NewIterator(emptySlice), mapper)
		actual := iterator.Collect()
		assertSlicesEqual(t, expected, actual, "Map empty slice failed")
	})
}

func TestFilter(t *testing.T) {
	intSlice := []int{10, 15, 22, 30, 7}
	t.Run("FilterEvenNumbers", func(t *testing.T) {
		filter := func(i int) bool {
			return i%2 == 0
		}
		expected := []int{10, 22, 30}
		iterator := Filter(NewIterator(intSlice), filter) // Call Filter directly
		actual := iterator.Collect()
		assertSlicesEqual(t, expected, actual, "Filter even numbers failed")
	})

	t.Run("FilterGreaterThan20", func(t *testing.T) {
		filter := func(i int) bool {
			return i > 20
		}
		expected := []int{22, 30}
		iterator := Filter(NewIterator(intSlice), filter)
		actual := iterator.Collect()
		assertSlicesEqual(t, expected, actual, "Filter > 20 failed")
	})
	t.Run("FilterNoMatch", func(t *testing.T) {
		filter := func(i int) bool {
			return i > 100
		}
		expected := []int{}
		iterator := Filter(NewIterator(intSlice), filter)
		actual := iterator.Collect()
		assertSlicesEqual(t, expected, actual, "Filter no match failed")
	})
}

func TestReduce(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5}

	t.Run("SumReduction", func(t *testing.T) {
		initial := 0
		reducer := func(acc int, curr int) int {
			return acc + curr
		}
		expected := 15
		actual := Reduce(NewIterator(intSlice), initial, reducer)
		if actual != expected {
			t.Fatalf("Sum reduction failed. Expected %d, got %d", expected, actual)
		}
	})

	t.Run("StringConcatenation", func(t *testing.T) {
		initial := "Start: "
		reducer := func(acc string, curr int) string {
			return acc + strconv.Itoa(curr)
		}
		expected := "Start: 12345"
		actual := Reduce(NewIterator(intSlice), initial, reducer)
		if actual != expected {
			t.Fatalf("String concatenation failed. Expected %s, got %s", expected, actual)
		}
	})

	t.Run("ProductReduction", func(t *testing.T) {
		initial := 1
		reducer := func(acc int, curr int) int {
			return acc * curr
		}
		expected := 120 // 1 * 2 * 3 * 4 * 5
		actual := Reduce(NewIterator(intSlice), initial, reducer)
		if actual != expected {
			t.Fatalf("Product reduction failed. Expected %d, got %d", expected, actual)
		}
	})

	t.Run("EmptySliceReduction", func(t *testing.T) {
		emptySlice := []int{}
		initial := 100
		reducer := func(acc int, curr int) int {
			return acc + curr
		}
		expected := 100
		actual := Reduce(NewIterator(emptySlice), initial, reducer)
		if actual != expected {
			t.Fatalf("Empty slice reduction failed. Expected %d, got %d", expected, actual)
		}
	})
}

func TestZipBasic(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := []string{"a", "b", "c"}

	it := Zip(s1, s2)
	got := it.Collect()

	want := []Tuple[int, string]{
		{1, "a"},
		{2, "b"},
		{3, "c"},
	}

	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("at index %d: got %+v want %+v", i, got[i], want[i])
		}
	}
}

func TestZipEmpty(t *testing.T) {
	s1 := []int{}
	s2 := []int{}

	it := Zip(s1, s2)
	got := it.Collect()

	if len(got) != 0 {
		t.Fatalf("expected empty result, got %v", got)
	}
}

func TestZipDifferentLengths(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := []int{4, 5}

	it := Zip(s1, s2)
	got := it.Collect()

	if len(got) != 0 {
		t.Fatalf("expected no iteration due to mismatched lengths, got %v", got)
	}
}

func TestZipEarlyStop(t *testing.T) {
	s1 := []int{1, 2, 3, 4}
	s2 := []int{10, 20, 30, 40}

	var got []Tuple[int, int]

	it := Zip(s1, s2)

	it(func(v Tuple[int, int]) bool {
		got = append(got, v)
		return false
	})

	if len(got) != 1 {
		t.Fatalf("expected early termination after 1 item, got %d", len(got))
	}

	if got[0] != (Tuple[int, int]{1, 10}) {
		t.Errorf("unexpected yielded value: %+v", got[0])
	}
}

func TestZipGenericTypes(t *testing.T) {
	type A struct{ X int }
	type B struct{ Y string }

	s1 := []A{{1}, {2}}
	s2 := []B{{"x"}, {"y"}}

	it := Zip(s1, s2)
	got := it.Collect()

	want := []Tuple[A, B]{
		{A{1}, B{"x"}},
		{A{2}, B{"y"}},
	}

	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("at index %d: got %+v want %+v", i, got[i], want[i])
		}
	}
}

func TestConcat(t *testing.T) {
	it1 := NewIterator([]int{1, 2})
	it2 := NewIterator([]int{3, 4})
	it3 := NewIterator([]int{})

	combined := Concat(it1, it2, it3)

	result := combined.Collect()

	expected := []int{1, 2, 3, 4}
	if len(result) != len(expected) {
		t.Fatalf("unexpected result length: got %d want %d", len(result), len(expected))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Fatalf("unexpected result at index %d: got %d want %d", i, result[i], v)
		}
	}
}

func TestConcatEarlyStop(t *testing.T) {
	it1 := NewIterator([]int{1, 2, 3})
	it2 := NewIterator([]int{4, 5})

	called := 0

	Concat(it1, it2)(func(v int) bool {
		called++
		return v < 2
	})

	if called != 2 {
		t.Fatalf("Concat did not stop early: expected 2 calls, got %d", called)
	}
}

func TestCount(t *testing.T) {
	it := NewIterator([]string{"a", "b", "c"})
	result := it.Count()

	if result != 3 {
		t.Fatalf("Count returned %d, want 3", result)
	}
}

func TestCountEmpty(t *testing.T) {
	it := NewIterator([]int{})
	result := it.Count()

	if result != 0 {
		t.Fatalf("Count returned %d, want 0", result)
	}
}

func TestAny(t *testing.T) {
	it := NewIterator([]int{1, 2, 3, 4})

	if !Any(it, func(v int) bool { return v == 3 }) {
		t.Fatalf("Any should have returned true for value 3")
	}
}

func TestAnyNoMatch(t *testing.T) {
	it := NewIterator([]int{1, 2, 3})

	if Any(it, func(v int) bool { return v == 5 }) {
		t.Fatalf("Any returned true, expected false")
	}
}

func TestAnyEmpty(t *testing.T) {
	it := NewIterator([]int{})

	if Any(it, func(v int) bool { return true }) {
		t.Fatalf("Any should return false on empty iterator")
	}
}

func TestAll(t *testing.T) {
	it := NewIterator([]int{2, 4, 6})

	if !All(it, func(v int) bool { return v%2 == 0 }) {
		t.Fatalf("All should return true when all elements match predicate")
	}
}

func TestAllFail(t *testing.T) {
	it := NewIterator([]int{2, 4, 5})

	if All(it, func(v int) bool { return v%2 == 0 }) {
		t.Fatalf("All returned true, but 5 does not match predicate")
	}
}

func TestAllEmpty(t *testing.T) {
	it := NewIterator([]int{})
	if !All(it, func(v int) bool { return false }) {
		t.Fatalf("All should return true for empty iterator by definition")
	}
}

/*
BENCHMARKS
*/

func BenchmarkMap_100000000(b *testing.B) {
	const size = 100_000_000
	data := make([]int, size)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()
	for b.Loop() {
		dataIter := NewIterator(data)
		Map(dataIter, func(x int) int {
			return x ^ 2
		})
	}
}

func BenchmarkFilter_100000000(b *testing.B) {
	const size = 100_000_000
	data := make([]int, size)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()
	for b.Loop() {
		dataIter := NewIterator(data)
		Filter(dataIter, func(x int) bool {
			return x%2 == 0
		})
	}
}

func BenchmarkReduce_100000000(b *testing.B) {
	const size = 100_000_000
	data := make([]int, size)
	for i := range data {
		data[i] = i
	}

	sumReduce := func(acc int, curr int) int {
		return acc + curr
	}
	b.ResetTimer()
	for b.Loop() {
		dataIter := NewIterator(data)
		Reduce(dataIter, 0, sumReduce)
	}
}

func BenchmarkZip_100000000(b *testing.B) {
	const size = 100_000_000
	data := make([]int, size)
	for i := range data {
		data[i] = i
	}
	b.ResetTimer()

	for b.Loop() {
		Zip(data, data)
	}
}

func BenchmarkConcat_100000000(b *testing.B) {
	const size = 100_000_000

	data1 := make([]int, size)
	data2 := make([]int, size)

	for i := range data1 {
		data1[i] = i
		data2[i] = i
	}

	it1 := NewIterator(data1)
	it2 := NewIterator(data2)

	b.ResetTimer()

	for b.Loop() {
		Concat(it1, it2).Collect()
	}
}
