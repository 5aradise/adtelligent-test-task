package slices

// FilterFunc filters the slice in place,
// leaving only those elements for which the function f returns true.
// FilterFunc returns slice with new len.
func FilterFunc[S ~[]E, E any](s S, f func(E) bool) S {
	if len(s) < 1 {
		return s
	}

	newLen := 0
main:
	for right := len(s) - 1; newLen <= right; newLen++ {
		if f(s[newLen]) {
			continue
		}
		for !f(s[right]) {
			if newLen >= right {
				break main
			}
			right--
		}
		s[newLen] = s[right]
		right--
	}
	return s[:newLen]
}
