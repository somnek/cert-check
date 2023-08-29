package main

func Find[T comparable](l []T, v T) int {
	for i, e := range l {
		if e == v {
			return i
		}
	}
	return -1
}

func Delete[T any](l []T, i int) []T {
	if len(l) < 1 {
		return l
	}

	return append(l[:i], l[i+1:]...)
}
