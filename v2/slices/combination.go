package slices

import (
	"math/rand"
)

func CombinationsSlice[T any](list []T, n int) [][]T {
	if n == 1 {
		result := make([][]T, 0)
		for i := 0; i < len(list); i++ {
			result = append(result, []T{list[i]})
		}
		return result
	}

	result := make([][]T, 0)
	for i := 0; i < len(list)-n+1; i++ {
		for _, c := range CombinationsSlice(list[i+1:], n-1) {
			result = append(result, append([]T{list[i]}, c...))
		}
	}
	return result
}

func RandomSelect[T any](source []T, count int) ([]T, error) {
	n := len(source)

	shuffled := make([]T, n)
	indexes := rand.Perm(n)

	for i, j := range indexes {
		shuffled[i] = source[j]
	}

	if count > n {
		// the count is greater than the length of the slice
		// so we need to append some random elements
		for i := 0; i < count-n; i++ {
			shuffled = append(shuffled, source[rand.Intn(n)])
		}
	}

	return shuffled[0:count], nil
}
