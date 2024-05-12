package hw04lrucache

import (
	"sync"
)

type Key string

type elem struct {
	Key   Key
	Value interface{}
}

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

	element := elem{Key: key, Value: value} // Готовим элемент к вставке в очередь
	if _, found := c.items[key]; found {    // если элемент найден в словаре
		c.items[key].Value = element      // присвоили элементу новое значение
		c.queue.MoveToFront(c.items[key]) // подвинули найденный элемент в начало списка
		return true                       // элемент был
	}
	// если элемент не найден в словаре - тогда это новый элемент который надо вставить
	// если емкость достигла предела - то перед вставкой нового элемента - удалить самый старый
	if c.queue.Len() >= c.capacity {
		// первым делом получить элемент (для последующего удаления из мапы) а только потом удалять его из очереди
		// тут нужно подставить ключ последнего элемента. Ключ нужно взять из элемента очереди.
		castedElement := c.queue.Back().Value.(elem) // берем интерфейс из очереди - его надо кастить чтобы добраться до ключа
		c.queue.Remove(c.queue.Back())               // получаем последний элемент из списка и удаляем его из списка
		delete(c.items, castedElement.Key)           // после удаления из списка - удалить из мапы
	}
	item := c.queue.PushFront(element) // создать новый элемент в списке (и получить созданный элемент)
	c.items[key] = item
	return false // элемента не было
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.Lock()

	defer c.Unlock()

	if _, found := c.items[key]; found { // если элемент найден в словаре
		c.queue.MoveToFront(c.items[key]) // подвинули найденный элемент в начало списка
		// вернули значение элемента и статус что он найден
		// В мапе находим по ключу элемент, но это интерфейс - кастим в структуру elem и только после этого берем Value
		return c.items[key].Value.(elem).Value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.Lock()

	defer c.Unlock()

	// Если вы хотите очистить карту от всех значений, вы можете задать
	// ее равной пустой карте того же типа. При этом будет создана новая
	// пустая карта, а сборщик мусора очистит память от старой карты.
	c.items = map[Key]*ListItem{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
