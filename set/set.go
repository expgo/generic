package set

import "github.com/expgo/generic/list"

func Add[E comparable](s []E, e E) ([]E, bool) {
	if list.Contains(s, e) {
		return s, false
	}

	s = append(s, e)
	return s, true
}

func AddFunc[E any](s []E, e E, matchFunc func(E) bool) ([]E, bool) {
	if list.ContainsFunc(s, e, matchFunc) {
		return s, false
	}

	s = append(s, e)
	return s, true
}
