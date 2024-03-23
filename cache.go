package generic

import "sync"

type Cache[K comparable, V any] struct {
	innerMap Map[K, *innerItem[V]]
}

type innerItem[V any] struct {
	value V
	err   error
	once  sync.Once
}

// GetOrNew retrieves the value associated with the key from the cache,
// or creates a new value using the provided new function if the key does not exist.
//
// It returns the retrieved value and any error that occurred during creation.
//
// If the new function is nil, a panic will occur.
//
// The logic of this method is as follows:
// - Load the value associated with the key from the cache. If the key does not exist, store an empty inner item.
// - Perform the new function inside a sync.Once to ensure it is only executed once.
// - Return the value and error of the inner item.
//
// Example usage:
//
//	cache := &Cache{}
//	key := "example"
//	newFunc := func() (string, error) {
//	    // Logic to create a new value
//	    return "new value", nil
//	}
//	value, err := cache.GetOrNew(key, newFunc)
func (c *Cache[K, V]) GetOrNew(k K, newFunc func() (V, error)) (v V, err error) {
	if newFunc == nil {
		panic("new function must not be null")
	}

	item, _ := c.innerMap.LoadOrStore(k, &innerItem[V]{})

	item.once.Do(func() {
		item.value, item.err = newFunc()
	})

	return item.value, item.err
}
