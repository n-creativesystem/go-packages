package logging

func toInterface[T any](values []T) []any {
	results := make([]any, len(values))
	for idx, value := range values {
		results[idx] = value
	}
	return results
}
