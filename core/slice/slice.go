package slice

func Map[T any, U any](s []T, mapper func(T) U) []U {
	result := make([]U, len(s), len(s))
	for i, value := range s {
		result[i] = mapper(value)
	}
	return result
}
