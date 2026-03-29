package rh

func DereferenceSlice[T any](input []*T) []T {
	result := make([]T, len(input))
	for i, v := range input {
		if v != nil {
			result[i] = *v
		}
	}
	return result
}
