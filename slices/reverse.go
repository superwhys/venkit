package slices

func Reverse[T any](slice []T) []T {
	ret := make([]T, len(slice))
	ret = append(slice[:0:0], slice...)
	for i := len(ret)/2 - 1; i >= 0; i-- {
		opp := len(ret) - 1 - i
		ret[i], ret[opp] = ret[opp], ret[i]
	}
	return ret
}
