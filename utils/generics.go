package utils

func Filter[T any](slice []T, f func(T) bool) []T {
	var n []T
	for _, e := range slice {
		if f(e) {
			n = append(n, e)
		}
	}
	return n
}

func In[T comparable](slice []T, x T) bool {
	for _, e := range slice {
		if e == x {
			return true
		}
	}

	return false
}

func Ptr[T any](v T) *T {
	return &v
}
