package common

func ValueOrDefault[T any](v *T, defaultValue T) T {
	if v != nil {
		return *v
	}
	return defaultValue
}
