package list

func Map[T any, R any](slice []T, f func(T) R) []R {
	result := make([]R, 0, len(slice))
	ForEach(slice, func(item T) {
		result = append(result, f(item))
	})
	return result
}

func GroupBy[T any, K comparable](slice []T, keyFunc func(T) K) map[K][]T {
	result := make(map[K][]T)
	ForEach(slice, func(item T) {
		key := keyFunc(item)
		result[key] = append(result[key], item)
	})
	return result
}

func ForEach[T any](slice []T, f func(T)) {
	for i := range slice {
		f(slice[i])
	}
}

func ForEachIndex[T any](slice []T, f func(int, T)) {
	for i := range slice {
		f(i, slice[i])
	}
}
