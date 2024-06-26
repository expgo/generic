package generic

import (
	"errors"
	"sync"
)

type Cache[K comparable, V any] struct {
	innerMap sync.Map
}

type innerItem[V any] struct {
	value V
	err   error
	once  sync.Once
}

// GetOrLoad retrieves the value associated with the specified key from the cache.
// If the entry does not exist, it calls the provided `loadFunc` function to load the value and store it in the cache.
// The `loadFunc` function should have the signature `func(k K) (V, error)`.
func (c *Cache[K, V]) GetOrLoad(k K, loadFunc func(k K) (V, error)) (v V, err error) {
	if loadFunc == nil {
		panic(errors.New("load function must not be nil"))
	}

	item, _ := c.innerMap.LoadOrStore(k, &innerItem[V]{})
	iItem := item.(*innerItem[V])

	iItem.once.Do(func() {
		iItem.value, iItem.err = loadFunc(k)
	})

	return iItem.value, iItem.err
}

// Evict removes the entry with the specified key from the cache.
// It returns true if the entry was successfully evicted, and false otherwise.
func (c *Cache[K, V]) Evict(k K) bool {
	_, ok := c.innerMap.LoadAndDelete(k)
	return ok
}

// Clear removes all entries from the cache.
// It resets the innerMap to an empty state.
func (c *Cache[K, V]) Clear() {
	c.innerMap = sync.Map{}
}
