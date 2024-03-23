package generic

import (
	"github.com/expgo/generic/stream"
	"sync"
)

type Map[K comparable, V any] struct {
	innerMap sync.Map
}

type CachePair[K comparable, V any] struct {
	K K
	V V
}

func FromMap[K comparable, V any](m map[K]V) *Map[K, V] {
	result := &Map[K, V]{}
	for k, v := range m {
		result.Store(k, v)
	}
	return result
}

func (m *Map[K, V]) Load(k K) (V, bool) {
	item, ok := m.innerMap.Load(k)
	return item.(V), ok
}

func (m *Map[K, V]) Store(k K, v V) {
	m.innerMap.Store(k, v)
}

func (m *Map[K, V]) LoadOrStore(k K, v V) (V, bool) {
	item, ok := m.innerMap.LoadOrStore(k, v)
	return item.(V), ok
}

func (m *Map[K, V]) LoadAndDelete(k K) (V, bool) {
	item, loaded := m.innerMap.LoadAndDelete(k)
	return item.(V), loaded
}

func (m *Map[K, V]) Delete(k K) {
	m.innerMap.Delete(k)
}

func (m *Map[K, V]) Swap(k K, v V) (V, bool) {
	item, loaded := m.innerMap.Swap(k, v)
	return item.(V), loaded
}

func (m *Map[K, V]) CompareAndSwap(k K, old, new V) bool {
	return m.innerMap.CompareAndSwap(k, old, new)
}

func (m *Map[K, V]) CompareAndDelete(k K, old V) (deleted bool) {
	return m.innerMap.CompareAndDelete(k, old)
}

func (m *Map[K, V]) Range(rangeFunc func(k K, v V) bool) {
	m.innerMap.Range(func(key, value any) bool {
		return rangeFunc(key.(K), value.(V))
	})
}

func (m *Map[K, V]) Filter(filterFunc func(k K, v V) bool) *Map[K, V] {
	filteredCache := &Map[K, V]{}
	m.innerMap.Range(func(key, value any) bool {
		k := key.(K)
		v := value.(V)
		if filterFunc(k, v) {
			filteredCache.Store(k, v)
		}
		return true
	})
	return filteredCache
}

func (m *Map[K, V]) FilterToStream(filterFunc func(k K, v V) bool) stream.Stream[*CachePair[K, V]] {
	result := stream.Stream[*CachePair[K, V]]{}

	m.innerMap.Range(func(key, value any) bool {
		k := key.(K)
		v := value.(V)
		if filterFunc(k, v) {
			result = result.Append(&CachePair[K, V]{K: k, V: v})
		}
		return true
	})

	return result
}

func (m *Map[K, V]) ToStream() stream.Stream[*CachePair[K, V]] {
	result := stream.Stream[*CachePair[K, V]]{}

	m.innerMap.Range(func(key, value any) bool {
		result = result.Append(&CachePair[K, V]{K: key.(K), V: value.(V)})
		return true
	})

	return result
}

func (m *Map[K, V]) ToMap() map[K]V {
	result := make(map[K]V)
	m.innerMap.Range(func(key, value any) bool {
		result[key.(K)] = value.(V)
		return true
	})
	return result
}

func (m *Map[K, V]) Size() int {
	size := 0
	m.innerMap.Range(func(key, value any) bool {
		size += 1
		return true
	})

	return size
}
