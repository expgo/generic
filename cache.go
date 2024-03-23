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
