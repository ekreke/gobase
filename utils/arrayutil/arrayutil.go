package arrayutil

import (
	"math"
	"slices"
)

// InArray check if value is in array
func InArray[T comparable](val T, arr []T) bool {
	return slices.Contains(arr, val)
}

// IsDuplicate check if array has duplicate element
func IsDuplicate[T comparable](arr []T) bool {
	keys := make(map[T]bool)
	for _, entry := range arr {
		if _, value := keys[entry]; value {
			return true
		}
		keys[entry] = true
	}
	return false
}

// RemoveDuplicate remove duplicate element in array
func RemoveDuplicate[T comparable](arr []T) []T {
	keys := make(map[T]bool)
	list := []T{}
	for _, entry := range arr {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// Intersect return intersection of two array
func Intersect[T comparable](a, b []T) []T {
	m := make(map[T]struct{})
	result := make([]T, 0)

	for _, v := range a {
		m[v] = struct{}{}
	}

	for _, v := range b {
		if _, ok := m[v]; ok {
			result = append(result, v)
			delete(m, v)
		}
	}
	return result
}

// Diff return set difference of two arrays (ignore duplicate counts).
// exp:
//
//	a: [1, 1, 2], b: [1, 3] => [2]
//	a: ["p"],     b: ["p", "p"] => []
func Diff[T comparable](a, b []T) []T {
	return DiffLogical(a, b)
}

// DiffLogical return set difference of two arrays (ignore duplicate counts).
// exp:
//
//	a: [1, 1, 2], b: [1, 3] => [2]
//	a: ["p"],     b: ["p", "p"] => []
func DiffLogical[T comparable](a, b []T) []T {
	inB := make(map[T]struct{}, len(b))
	for _, v := range b {
		inB[v] = struct{}{}
	}

	seen := make(map[T]struct{}, len(a))
	result := make([]T, 0)
	for _, v := range a {
		if _, ok := inB[v]; ok {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}

// DiffCount return multiset difference of two arrays (respect duplicate counts).
// exp:
//
//	a: [1, 1, 2], b: [1, 3] => [1, 2]
//	a: [1, 2],    b: [1, 1] => [2]
func DiffCount[T comparable](a, b []T) []T {
	counter := make(map[T]int, len(b))
	for _, v := range b {
		counter[v]++
	}

	result := make([]T, 0)
	for _, v := range a {
		if counter[v] > 0 {
			counter[v]--
			continue
		}
		result = append(result, v)
	}
	return result
}

// SymmetricDiff return symmetric set difference of two arrays (ignore duplicate counts).
// exp:
//
//	a: [1, 1, 2], b: [1, 3] => [2, 3]
func SymmetricDiff[T comparable](a, b []T) []T {
	return SymmetricDiffLogical(a, b)
}

// SymmetricDiffLogical return symmetric set difference of two arrays (ignore duplicate counts).
// exp:
//
//	a: [1, 1, 2], b: [1, 3] => [2, 3]
func SymmetricDiffLogical[T comparable](a, b []T) []T {
	result := DiffLogical(a, b)
	result = append(result, DiffLogical(b, a)...)
	return result
}

// SymmetricDiffCount return symmetric multiset difference of two arrays (respect duplicate counts).
// exp:
//
//	a: [1, 1, 2], b: [1, 3] => [1, 2, 3]
func SymmetricDiffCount[T comparable](a, b []T) []T {
	result := DiffCount(a, b)
	result = append(result, DiffCount(b, a)...)
	return result
}

// IsSubset check if a is subset of b
func IsSubset[T comparable](a, b []T) bool {
	return IsSubsetLogical(a, b)
}

// IsSubsetLogical check if a is subset of b (ignore duplicate counts)
// exp:
//
//	a: [1, 1, 2], b: [1, 2, 3] => true
//	a: [1, 4],    b: [1, 2, 3] => false
//	a: ["p"],     b: ["p", "p"] => true
func IsSubsetLogical[T comparable](a, b []T) bool {
	m := make(map[T]struct{})
	for _, v := range b {
		m[v] = struct{}{}
	}

	for _, v := range a {
		if _, ok := m[v]; !ok {
			return false
		}
	}
	return true
}

// IsSubsetCount check if a is subset of b (respect duplicate counts)
// exp:
//
//	a: [1, 1], b: [1, 2, 3] => false
//	a: [1, 1], b: [1, 1, 2] => true
//	a: [1, 2], b: [1, 1, 2] => true
func IsSubsetCount[T comparable](a, b []T) bool {
	if len(a) > len(b) {
		return false
	}

	counter := make(map[T]int, len(b))
	for _, v := range b {
		counter[v]++
	}

	for _, v := range a {
		if counter[v] <= 0 {
			return false
		}
		counter[v]--
	}
	return true
}

// Union return union of two arrays (ignore duplicate counts).
func Union[T comparable](a, b []T) []T {
	m := make(map[T]struct{})
	result := make([]T, 0)
	for _, v := range a {
		m[v] = struct{}{}
	}
	for _, v := range b {
		m[v] = struct{}{}
	}
	for v := range m {
		result = append(result, v)
	}
	return result
}

// Chunk split array into chunks of size n
func Chunk[T any](arr []T, size int) [][]T {
	if size <= 0 {
		return [][]T{}
	}
	length := len(arr)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	result := make([][]T, 0, chunks)
	for i := 0; i < chunks; i += size {
		end := i + size
		if end > length {
			end = length
		}
		result = append(result, arr[i:end])
	}
	return result
}
