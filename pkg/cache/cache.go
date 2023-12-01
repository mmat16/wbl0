package cache

import (
	"errors"
	"time"
)

func New(defaultExpiry, cleanupInterval time.Duration) *Cache {
	items := make(map[string]Item)

	cache := Cache{
		items:           items,
		defaultExpiry:   defaultExpiry,
		cleanupInterval: cleanupInterval,
	}

	if cleanupInterval > 0 {
		cache.StartGC()
	}

	return &cache
}

func (c *Cache) Set(key string, value any, duration time.Duration) {
	var expiry int64

	if duration == 0 {
		duration = c.defaultExpiry
	}

	if duration > 0 {
		expiry = time.Now().Add(duration).UnixNano()
	}

	c.Lock()
	defer c.Unlock()

	c.items[key] = Item{
		Value:   value,
		Expiry:  expiry,
		Created: time.Now(),
	}
}

func (c *Cache) Get(key string) (any, bool) {
	c.RLock()
	defer c.RUnlock()

	item, found := c.items[key]

	if !found {
		return nil, false
	}

	if item.Expiry > 0 {
		if time.Now().UnixNano() > item.Expiry {
			return nil, false
		}
	}

	return item.Value, true
}

func (c *Cache) Delete(key string) error {
	c.Lock()
	defer c.Unlock()

	_, found := c.items[key]
	if !found {
		return errors.New("item not found in cache")
	}

	delete(c.items, key)
	return nil
}

func (c *Cache) StartGC() {
	go c.GC()
}

func (c *Cache) GC() {
	for {
		<-time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		keys := c.expiredKeys()
		if len(keys) > 0 {
			c.clearItems(keys)
		}
	}
}

func (c *Cache) expiredKeys() []string {
	c.RLock()
	defer c.RUnlock()

	var keys []string
	for key, item := range c.items {
		if time.Now().UnixNano() > item.Expiry && item.Expiry > 0 {
			keys = append(keys, key)
		}
	}

	return keys
}

func (c *Cache) clearItems(keys []string) {
	c.Lock()
	defer c.Unlock()

	for _, key := range keys {
		delete(c.items, key)
	}
}

func (c *Cache) GetAllItems() []any {
	res := make([]any, len(c.items))
	for _, value := range c.items {
		res = append(res, value)
	}
	return res
}