package cache

import (
	"sync"
)

// type Cache[K comparable, V any] struct {
// 	store      sync.Map
// 	TTL        time.Duration
// 	ReloadFunc func(key K) (V, error)
// }

// func (c *Cache[K, V]) Get(key K) (V, bool) {
// 	val, ok := c.store.Load(key)
// 	if ok {
// 		fmt.Println("có liền")
// 		return val.(V), ok
// 	} else if c.ReloadFunc != nil {
// 		return c.reloadAndSetTTL(key)
// 	}
// 	var defaultValue V
// 	return defaultValue, false
// }

// func (c *Cache[K, V]) reloadAndSetTTL(key K) (V, bool) {
// 	newVal, err := c.ReloadFunc(key)
// 	if err == nil {
// 		c.store.Store(key, newVal)
// 		if c.TTL > 0 {
// 			time.AfterFunc(c.TTL, func() {
// 				c.store.Delete(key)
// 				fmt.Println("re fetch")
// 				c.reloadAndSetTTL(key)
// 			})
// 		}
// 		return newVal, true
// 	}
// 	var defaultValue V
// 	return defaultValue, false
// }

type Cache[K comparable, V any] struct {
	store sync.Map
}

func (c *Cache[K, V]) Set(key K, val V) {
	c.store.Store(key, val)
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	val, ok := c.store.Load(key)
	if ok {
		return val.(V), ok
	}
	var defaultValue V
	return defaultValue, false
}

func (c *Cache[K, V]) Forget(key K) {
	c.store.Delete(key)
}

func (c *Cache[K, V]) Warmup(warmupFunc func() (map[K]V, error)) {
	if warmupFunc != nil {
		m, err := warmupFunc()
		if err != nil {
			panic(err)
		}
		for k, v := range m {
			c.Set(k, v)
		}
	}
}
