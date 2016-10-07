package concurrenthashmap

import (
	"fmt"
	"strings"
	"sync"
)

// Map holds the elements in go's native map and uses RWMutex to
// guard concurrent access to it
type Map struct {
	items map[interface{}]interface{}
	// extends Read Write mutex, guards concurrent access to items.
	sync.RWMutex
}

// New instantiates a concurrent hash map
func New() *Map {
	return &Map{items: make(map[interface{}]interface{})}
}

// Put inserts an entry into the map.
func (m *Map) Put(key interface{}, value interface{}) {
	m.Lock()
	m.items[key] = value
	m.Unlock()
}

// PutIfAbsent inserts an entry into the map, if the key doesn't exists
func (m *Map) PutIfAbsent(key interface{}, value interface{}) {
	m.Lock()
	_, found := m.items[key]
	if !found {
		m.items[key] = value
	}
	m.Unlock()
}

// Get searches the element in the map by key and returns its value or nil if key doesn't exists.
// Second return parameter is true if key was found, otherwise false.
func (m *Map) Get(key interface{}) (value interface{}, found bool) {
	m.RLock()
	value, found = m.items[key]
	m.RUnlock()
	return
}

// Remove removes the element from the map by key.
func (m *Map) Remove(key interface{}) {
	m.Lock()
	delete(m.items, key)
	m.Unlock()
}

// IsEmpty returns true if map does not contain any elements
func (m *Map) IsEmpty() bool {
	return m.Size() == 0
}

// Size returns number of el
// ements in the map.
func (m *Map) Size() int {
	m.RLock()
	size := len(m.items)
	m.RUnlock()
	return size
}

// Keys returns all keys of the map(random order).
func (m *Map) Keys() []interface{} {
	keys := make([]interface{}, m.Size())
	index := 0
	m.RLock()
	for key, _ := range m.items {
		keys[index] = key
		index++
	}
	m.RUnlock()
	return keys
}

// Values returns all values of the map (random order).
func (m *Map) Values() []interface{} {
	values := make([]interface{}, m.Size())
	index := 0
	m.RLock()
	for _, value := range m.items {
		values[index] = value
		index++
	}
	m.RUnlock()
	return values
}

// Returns true if the given keys are found in the map
func (m *Map) Contains(keys ...interface{}) bool {
	m.RLock()
	for _, key := range keys {
		_, found := m.items[key]
		if !found {
			m.RUnlock()
			return false
		}
	}
	m.RUnlock()
	return true
}

// Clear removes all elements from the map.
func (m *Map) Clear() {
	m.items = make(map[interface{}]interface{})
}

// String returns a string representation of container
func (m *Map) String() string {
	str := "ConcurrentHashMap\nMap["
	m.RLock()
	for key, value := range m.items {
		str += fmt.Sprintf("%v:%v ", key, value)
	}
	m.RUnlock()
	return strings.TrimRight(str, " ") + "]"
}
