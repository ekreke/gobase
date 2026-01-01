package arrayutil

import "testing"

func TestDifferenceLogical(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []int
		want []int
	}{
		{name: "empty", a: nil, b: nil, want: nil},
		{name: "basic", a: []int{1, 1, 2}, b: []int{1, 3}, want: []int{2}},
		{name: "allRemoved", a: []int{1}, b: []int{1, 1, 1}, want: nil},
		{name: "dedupeOutput", a: []int{2, 2, 2}, b: []int{1}, want: []int{2}},
		{name: "preserveOrderUnique", a: []int{3, 2, 3, 1}, b: []int{1}, want: []int{3, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DiffLogical(tt.a, tt.b)
			if !slicesEqual(got, tt.want) {
				t.Fatalf("DifferenceLogical(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDifferenceCount(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []int
		want []int
	}{
		{name: "empty", a: nil, b: nil, want: nil},
		{name: "respectCounts", a: []int{1, 1, 2}, b: []int{1, 3}, want: []int{1, 2}},
		{name: "removeAllCounts", a: []int{1, 1}, b: []int{1, 1, 2}, want: nil},
		{name: "bHasMoreCounts", a: []int{1, 2}, b: []int{1, 1}, want: []int{2}},
		{name: "preserveOrder", a: []int{2, 1, 2, 3}, b: []int{2}, want: []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DiffCount(tt.a, tt.b)
			if !slicesEqual(got, tt.want) {
				t.Fatalf("DifferenceCount(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSymmetricDifferenceLogical(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []int
		want []int
	}{
		{name: "empty", a: nil, b: nil, want: nil},
		{name: "basic", a: []int{1, 1, 2}, b: []int{1, 3}, want: []int{2, 3}},
		{name: "disjoint", a: []int{1, 2}, b: []int{3, 4}, want: []int{1, 2, 3, 4}},
		{name: "same", a: []int{1, 2}, b: []int{2, 1}, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SymmetricDiffLogical(tt.a, tt.b)
			if !slicesEqual(got, tt.want) {
				t.Fatalf("SymmetricDifferenceLogical(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSymmetricDifferenceCount(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []int
		want []int
	}{
		{name: "empty", a: nil, b: nil, want: nil},
		{name: "basic", a: []int{1, 1, 2}, b: []int{1, 3}, want: []int{1, 2, 3}},
		{name: "countCancels", a: []int{1, 1}, b: []int{1}, want: []int{1}},
		{name: "countCancelsOtherSide", a: []int{1}, b: []int{1, 1}, want: []int{1}},
		{name: "disjoint", a: []int{1, 2}, b: []int{3, 4}, want: []int{1, 2, 3, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SymmetricDiffCount(tt.a, tt.b)
			if !slicesEqual(got, tt.want) {
				t.Fatalf("SymmetricDifferenceCount(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestIsSubsetLogical(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []int
		want bool
	}{
		{name: "emptyInEmpty", a: nil, b: nil, want: true},
		{name: "emptyInNonEmpty", a: []int{}, b: []int{1}, want: true},
		{name: "nonEmptyInEmpty", a: []int{1}, b: nil, want: false},
		{name: "ignoreDuplicatesTrue", a: []int{1, 1}, b: []int{1, 2}, want: true},
		{name: "missingElementFalse", a: []int{1, 3}, b: []int{1, 2}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSubsetLogical(tt.a, tt.b); got != tt.want {
				t.Fatalf("IsSubsetLogical(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestIsSubsetCount(t *testing.T) {
	tests := []struct {
		name string
		a    []int
		b    []int
		want bool
	}{
		{name: "emptyInEmpty", a: nil, b: nil, want: true},
		{name: "emptyInNonEmpty", a: []int{}, b: []int{1}, want: true},
		{name: "nonEmptyInEmpty", a: []int{1}, b: nil, want: false},
		{name: "countMattersFalse", a: []int{1, 1}, b: []int{1, 2}, want: false},
		{name: "countEnoughTrue", a: []int{1, 1}, b: []int{1, 1, 2}, want: true},
		{name: "orderIrrelevantTrue", a: []int{2, 1, 2}, b: []int{1, 2, 2, 3}, want: true},
		{name: "countNotEnoughFalse", a: []int{2, 2, 2}, b: []int{2, 2, 3}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSubsetCount(tt.a, tt.b); got != tt.want {
				t.Fatalf("IsSubsetCount(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestIsSubsetAlias(t *testing.T) {
	a := []int{1, 1}
	b := []int{1, 2}
	if got, want := IsSubset(a, b), IsSubsetLogical(a, b); got != want {
		t.Fatalf("IsSubset should behave like IsSubsetLogical: got %v, want %v", got, want)
	}
}

func slicesEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
