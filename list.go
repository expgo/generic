package generic

import "errors"

type List[T comparable] struct {
	items []T
}

func (l *List[T]) At(idx int) (t T, err error) {
	if idx >= l.Size() {
		return t, errors.New("index out of range")
	}

	return l.items[idx], nil
}

func (l *List[T]) Add(e T) bool {
	l.items = append(l.items, e)
	return true
}

func (l *List[T]) Remove(e T) bool {
	for i, item := range l.items {
		if e == item {
			l.items = append(l.items[:i], l.items[i+1:]...)
			return true
		}
	}
	return false
}

func (l *List[T]) Contains(e T) bool {
	for _, item := range l.items {
		if e == item {
			return true
		}
	}
	return false
}

func (l *List[T]) Clear() {
	l.items = nil
}

func (l *List[T]) RemoveAt(idx int) bool {
	if idx < 0 || idx >= l.Size() {
		return false
	}

	l.items = append(l.items[:idx], l.items[idx+1:]...)
	return true
}

func (l *List[T]) Size() int {
	return len(l.items)
}
