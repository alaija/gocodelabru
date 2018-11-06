package lru

import (
	"container/list"
	"errors"
)

type (
	// LRU структура представляющая собой LRU-cache
	LRU struct {
		size      int
		evictList *list.List
		items     map[interface{}]*list.Element
	}
	// entry используется для хранения значений в `evictList`
	entry struct {
		key   interface{}
		value interface{}
	}
)

// New создает новый `LRU` с заданным размером
func New(size int) (*LRU, error) {
	if size <= 0 {
		return nil, errors.New("Size must be greater than 0")
	}
	c := &LRU{
		size:      size,
		evictList: list.New(),
		items:     make(map[interface{}]*list.Element),
	}
	return c, nil
}

// Add добавлдяет значение в LRU.
// Возвращает `true` если произошло вытеснение елемента
func (l *LRU) Add(key, value interface{}) bool {
	if ent, ok := l.items[key]; ok {
		l.evictList.MoveToFront(ent)
		ent.Value.(*entry).value = value
		return false
	}
	ent := &entry{key, value}
	entry := l.evictList.PushFront(ent)
	l.items[key] = entry
	evict := l.Len() > l.size
	if evict {
		l.RemoveOldest()
	}
	return evict
}

// Len возвращает количество элементов в кэше
func (l *LRU) Len() int {
	return l.evictList.Len()
}

// Purge очищает кэш
func (l *LRU) Purge() {
	for k := range l.items {
		delete(l.items, k)
	}
	l.evictList.Init()
}

// Get возвращет элемент кэша по ключу
func (l *LRU) Get(key interface{}) (value interface{}, ok bool) {
	if ent, ok := l.items[key]; ok {
		l.evictList.MoveToFront(ent)
		return ent.Value.(*entry).value, true
	}
	return
}

// Remove удаляет элемент по ключу.
// Возвращает `true` если элемент был в кэше
func (l *LRU) Remove(key interface{}) bool {
	if ent, ok := l.items[key]; ok {
		l.removeElement(ent)
		return true
	}
	return false
}

// Contains проверяет наличе элемента
// не обновляя состояние кэша.
func (l *LRU) Contains(key interface{}) (ok bool) {
	_, ok = l.items[key]
	return ok
}

// RemoveOldest удаляет старейший элемент из кэша
func (l *LRU) RemoveOldest() (interface{}, interface{}, bool) {
	ent := l.evictList.Back()
	if ent != nil {
		l.removeElement(ent)
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return nil, nil, false
}

// GetOldest получает старейший элемент из кэша
func (l *LRU) GetOldest() (interface{}, interface{}, bool) {
	ent := l.evictList.Back()
	if ent != nil {
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return nil, nil, false
}

// Keys возвращает ключи элементов в кэше
func (l *LRU) Keys() []interface{} {
	keys := make([]interface{}, len(l.items))
	i := 0
	for ent := l.evictList.Back(); ent != nil; ent = ent.Prev() {
		keys[i] = ent.Value.(*entry).key
		i++
	}
	return keys
}

func (l *LRU) removeElement(e *list.Element) {
	l.evictList.Remove(e)
	kv := e.Value.(*entry)
	delete(l.items, kv.key)
}
