package generic

import (
	"github.com/expgo/generic/stream"
	"sync"
)

type Set[T comparable] struct {
	itemList List[T]
	elemMap  sync.Map
}

func (s *Set[T]) Add(e T) bool {
	_, loaded := s.elemMap.LoadOrStore(e, true)
	if !loaded {
		s.itemList.Add(e)
	}
	return !loaded
}

func (s *Set[T]) Remove(e T) {
	s.elemMap.Delete(e)
	s.itemList.Remove(e)
}

func (s *Set[T]) Clear() {
	s.itemList.Clear()
	s.elemMap = sync.Map{}
}

func (s *Set[T]) Contains(e T) bool {
	_, ok := s.elemMap.Load(e)
	return ok
}

func (s *Set[T]) Size() int {
	return s.itemList.Size()
}

func (s *Set[T]) ToStream() (result stream.Stream[T]) {
	return stream.Of(s.itemList.items)
}
