package gmap

import "github.com/expgo/sync"

type Map[K comparable, V any] struct {
	items map[K]V
	lock  sync.RWMutex
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		items: map[K]V{},
		lock:  sync.NewRWMutex(),
	}
}

func Clone[K comparable, V any](originalMap map[K]V) map[K]V {
	cloned := make(map[K]V)

	for key, value := range originalMap {
		cloned[key] = value
	}

	return cloned
}

func Load[K comparable, V any](m *Map[K, V], key K) (value V, ok bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	value, ok = m.items[key]
	return
}

func LoadAndDelete[K comparable, V any](m *Map[K, V], key K) (value V, loaded bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	value, loaded = m.items[key]
	delete(m.items, key)

	return
}

func LoadOrStore[K comparable, V any](m *Map[K, V], key K, value V) (actual V, loaded bool) {
	m.lock.RLock()
	actual, loaded = m.items[key]
	m.lock.RUnlock()

	if !loaded {
		m.lock.Lock()
		m.items[key] = value
		m.lock.Unlock()

		actual = value
	}

	return
}

func Store[K comparable, V any](m *Map[K, V], key K, value V) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.items[key] = value
}

func Delete[K comparable, V any](m *Map[K, V], key K) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.items, key)
}

func Range[K comparable, V any](m *Map[K, V], f func(key K, value V) bool) {
	m.lock.RLock()
	mm := Clone(m.items)
	m.lock.RUnlock()

	for key, value := range mm {
		if !f(key, value) {
			break
		}
	}
}
