package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (lru *lruCache) Clear() {
	for i := lru.queue.Front(); i != nil; i = i.Next {
		i.Prev = nil
		i.Next = nil
	}

	lru.queue = NewList()
	lru.items = make(map[Key]*ListItem, lru.capacity)
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	if li, ok := lru.items[key]; ok {
		lru.queue.MoveToFront(li)

		return li.Value.(cacheItem).value, ok
	}

	return nil, false
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	li, ok := lru.items[key]

	if ok {
		li.Value = cacheItem{
			key:   key,
			value: value,
		}

		lru.queue.MoveToFront(li)
	} else {
		lru.items[key] = lru.queue.PushFront(
			cacheItem{
				key:   key,
				value: value,
			})

		if lru.queue.Len() > lru.capacity {
			back := lru.queue.Back()
			lru.queue.Remove(back)

			delete(lru.items, back.Value.(cacheItem).key)
		}
	}

	return ok
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
