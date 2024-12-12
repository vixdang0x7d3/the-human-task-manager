package postgres

func Reduce[T any, R any](slice []T, initial R, f func(R, T) R) R {
	result := initial
	for _, item := range slice {
		result = f(result, item)
	}
	return result
}

func ForEach[T any](slice []T, f func(T)) {
	for _, item := range slice {
		f(item)
	}
}

func Map[T, U any](slice []T, f func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
}
