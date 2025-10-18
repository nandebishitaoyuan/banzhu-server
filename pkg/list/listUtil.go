package list

func MapSlice[T any, R any](slice []T, f func(T) R) []R {
	result := make([]R, 0, len(slice))
	for _, v := range slice {
		result = append(result, f(v))
	}
	return result
}

func GroupBy[T any, K comparable](slice []T, keyFunc func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range slice {
		key := keyFunc(item)
		result[key] = append(result[key], item)
	}
	return result
}
