package isunippets

func GroupBy[T any](objects []*T, keyFunc func(*T) string) map[string][]T {
	grouped := make(map[string][]T)
	for _, obj := range objects {
		key := keyFunc(obj)
		grouped[key] = append(grouped[key], *obj)
	}
	return grouped
}

func GroupById[T any, U int64](objects []*T, keyFunc func(*T) U) map[U][]T {
	grouped := make(map[U][]T)
	for _, obj := range objects {
		key := keyFunc(obj)
		grouped[key] = append(grouped[key], *obj)
	}
	return grouped
}
