package cache

import (
	"sync"
	"time"
)

const (
	NoExpire       = -1
	NoCleanup      = -1
	DefaultExpiry  = time.Minute * 5
	DefaultCleanup = time.Minute * 5
)

// Cacher instance contains unimported fields such as
// sync.RWMutex, default expiration and cleanup time intervals and
// a map to store the items in the in-memory cache
type Cache struct {
	sync.RWMutex
	defaultExpiry   time.Duration
	cleanupInterval time.Duration
	items           map[string]Item
}

// Instance to store in the cache
// Contains value to store, its expiration interval and "created" timestamp
type Item struct {
	Value   any
	Created time.Time
	Expiry  int64
}
