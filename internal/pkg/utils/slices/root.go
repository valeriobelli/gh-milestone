package slices

func Contains[T comparable](values []T, element T) bool {
	for _, currentValue := range values {
		if currentValue == element {
			return true
		}
	}

	return false
}

func Map[T comparable](values []T, valueMapper func(T) T) []T {
	mappedValues := make([]T, len(values))

	for index, currentValue := range values {
		mappedValues[index] = valueMapper(currentValue)
	}

	return mappedValues
}
