package hw04lrucache

import (
	"sync"
)

type Key string

// LRU - когда вытесняется элемент, к которому дольше всего не было обращений.
type CacheMethods interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type LruCache struct {
	sync.Mutex

	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *LruCache) Set(key Key, value interface{}) {

	c.Lock()

	defer c.Unlock()

	c.items[key] = &ListItem{
		Value: value,
		Next: nil,
		Prev: nil,
	}

	c.queue.PushFront(value)

}

func NewCache(capacity int) LruCache {
	return LruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
