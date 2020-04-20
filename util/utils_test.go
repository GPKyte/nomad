package util

import (
	"testing"
)

func TestPermute(t *testing.T) {
	permutations := permute("a", "b", "c")
	expectations := factorial(3) / factorial(3-3) /* n! / (n-r)! is #permutations */

	if int64(len(permutations)) != expectations {
		t.Errorf("Expected 6 got %d:\n%s", len(permutations), permutations)
	}
	var countDistinct map[interface{}]int

	for p := 0; p < len(permutations); p++ {
		countDistinct[permutations[p]]++
	}
	for _, val := range countDistinct {
		if val > 1 {
			t.Error(countDistinct)
		}
	}
}

func TestPickUnique(t *testing.T) {
	type testType generic
	nums := []testType{
		{2},
		{3},
		{5},
		{6},
		{9},
		{1},
		{5},
		{0},
	}

	var chosen = pickUnique(4, &nums)
	if len(chosen) != 4 {
		panic(len(chosen))
	}
}

func TestFactorial(t *testing.T) {
	tests := []struct{ A, B int64 }{
		{factorial(0), 1},
		{factorial(1), 1},
		{factorial(2), 2},
		{factorial(3), 6},
		{factorial(4), 24},
		{factorial(5), 120},
		{factorial(6), 720},
		{factorial(15), 1307674368000},
	}

	for i := 0; i < len(tests); i++ {
		got := tests[i].A
		want := tests[i].B

		if want != got {
			t.Errorf("Want %d, but got %d", want, got)
		}
	}
}
