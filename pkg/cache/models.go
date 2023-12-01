package cache

import (
	"sync"
	"time"
)

const (
	NoExpire  = -1
	NoCleanup = -1
)

type Cache struct {
	sync.RWMutex
	defaultExpiry   time.Duration
	cleanupInterval time.Duration
	items           map[string]Item
}

type Item struct {
	Value   any
	Created time.Time
	Expiry  int64
}
