package cache

import (
	"sync"
	"time"
)

type Cache[K comparable, V any] struct {
	store      sync.Map
	TTL        time.Duration
	ReloadFunc func(key K) (V, error)
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	val, ok := c.store.Load(key)
	if ok {
		return val.(V), ok
	} else if c.ReloadFunc != nil {
		return c.reloadAndSetTTL(key, c.TTL)
	}
	var defaultValue V
	return defaultValue, false
}

func (c *Cache[K, V]) setAndSetTTL(key K, val V, ttl time.Duration) (V, bool) {
	c.store.Store(key, val)
	time.AfterFunc(ttl, func() {
		c.store.Delete(key)
		c.reloadAndSetTTL(key, c.TTL)
	})
	return val, true
}

func (c *Cache[K, V]) reloadAndSetTTL(key K, ttl time.Duration) (V, bool) {
	newVal, err := c.ReloadFunc(key)
	if err == nil {
		return c.setAndSetTTL(key, newVal, c.TTL)
	}
	var defaultValue V
	return defaultValue, false
}
