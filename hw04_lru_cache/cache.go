package hw04lrucache

import (
	"sync"
)

type Key string

// LRU - когда вытесняется элемент, к которому дольше всего не было обращений.
type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.Lock()

	defer c.Unlock()

	if _, found := c.items[key]; found { // если элемент найден в словаре
		c.items[key].Value = value        // присвоили элементу новое значение
		c.queue.MoveToFront(c.items[key]) // подвинули найденный элемент в начало списка
		return true                       // элемент был
	}
	// если элемент не найден в словаре - тогда это новый элемент который надо вставить
	// если емкость достигла предела - то перед вставкой нового элемента - удалить самый старый
	if c.queue.Len() >= c.capacity {
		c.queue.Remove(c.queue.Back()) // получаем последний элемент из списка и удаляем его из списка
		delete(c.items, key)           // после удаления из списка - удалить из мапы
	}
	item := c.queue.PushFront(value) // создать новый элемент в списке (и получить созданный элемент)
	c.items[key] = item
	return false // элемента не было
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.Lock()

	defer c.Unlock()

	if _, found := c.items[key]; found { // если элемент найден в словаре
		c.queue.MoveToFront(c.items[key]) // подвинули найденный элемент в начало списка
		return c.items[key].Value, true   // вернули значение элемента и статус что он найден
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.Lock()

	defer c.Unlock()

	for k := range c.items {
		delete(c.items, k)
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
