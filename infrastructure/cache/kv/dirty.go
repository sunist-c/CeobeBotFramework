package kv

import (
	"sync"
	"time"
)

type entity struct {
	timeout time.Time
	data    interface{}
}

type dirtyMap struct {
	lock  sync.RWMutex
	dirty map[interface{}]*entity
}
