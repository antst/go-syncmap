package syncmap

import (
	"sync"
)

// SyncMap is a generic, thread-safe map implementation.
// It uses a read-write mutex to ensure safe concurrent access to the underlying map.
//
// Type parameters:
//
//	K: must be a comparable type (used as map keys)
//	V: can be any type (used as map values)
type SyncMap[K comparable, V any] struct {
	Mu   sync.RWMutex
	Data map[K]V
}

// New creates and returns a new SyncMap with the specified initial size.
// It initializes the internal map and mutex for thread-safe operations.
func New[K comparable, V any](size int) *SyncMap[K, V] {
	return &SyncMap[K, V]{
		Mu:   sync.RWMutex{},
		Data: make(map[K]V, size),
	}
}

// Store adds or updates a key-value pair in the SyncMap.
// It acquires a write lock to ensure thread-safe access to the underlying data.
func (m *SyncMap[K, V]) Store(k K, v V) {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	m.Data[k] = v
}

// Load retrieves the value associated with the given key from the SyncMap.
// It acquires a read lock to ensure thread-safe access to the underlying data.
func (m *SyncMap[K, V]) Load(k K) (V, bool) {
	m.Mu.RLock()
	defer m.Mu.RUnlock()

	v, ok := m.Data[k]
	return v, ok
}

// Remove deletes the value associated with the given key from the SyncMap.
// It acquires a write lock to ensure thread-safe access to the underlying data.
func (m *SyncMap[K, V]) Remove(k K) bool {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	if _, ok := m.Data[k]; !ok {
		return false
	}

	delete(m.Data, k)

	return true
}

// Map applies a given function to all key-value pairs in the SyncMap and returns a new map with the results.
// It acquires a read lock to ensure thread-safe access to the underlying data.
func (m *SyncMap[K, V]) Map(mapFn func(k K, v V) V) map[K]V {
	m.Mu.RLock()
	defer m.Mu.RUnlock()

	data := make(map[K]V, len(m.Data))

	for k, v := range m.Data {
		data[k] = mapFn(k, v)
	}

	return data
}

// Filter creates a new map containing key-value pairs from the SyncMap that satisfy the given predicate function.
// It acquires a read lock to ensure thread-safe access to the underlying data.
func (m *SyncMap[K, V]) Filter(predicateFn func(k K, v V) bool) map[K]V {
	data := make(map[K]V)

	m.Mu.RLock()
	defer m.Mu.RUnlock()

	for k, v := range m.Data {
		if predicateFn(k, v) {
			data[k] = v
		}
	}

	return data
}

// Purge removes all key-value pairs from the SyncMap, effectively clearing its contents.
// It acquires a write lock to ensure thread-safe access to the underlying data.
func (m *SyncMap[K, V]) Purge() {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	m.Data = make(map[K]V)
}

// Len returns the number of key-value pairs in the SyncMap.
// It acquires a read lock to ensure thread-safe access to the underlying data.
func (m *SyncMap[K, V]) Len() int {
	m.Mu.RLock()
	defer m.Mu.RUnlock()

	return len(m.Data)
}

// DoLocked executes a function with exclusive access to the SyncMap.
// It acquires a write lock before executing the function and releases it afterward.
func (m *SyncMap[K, V]) DoLocked(f func(LockedMap[K, V])) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	f(&lockedMap[K, V]{m: m})
}

// DoLockedWithResult executes a function with exclusive access to the SyncMap and returns its result.
// It acquires a write lock before executing the function and releases it afterward.
func (m *SyncMap[K, V]) DoLockedWithResult(f func(LockedMap[K, V]) any) any {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	return f(&lockedMap[K, V]{m: m})
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
// It acquires a write lock to ensure thread-safe access to the underlying data.
func (m *SyncMap[K, V]) LoadOrStore(key K, value V) (V, bool) {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	if v, ok := m.Data[key]; ok {
		return v, true
	}

	m.Data[key] = value
	return value, false
}

// LoadAndDelete removes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
// It acquires a write lock to ensure thread-safe access to the underlying data.
func (m *SyncMap[K, V]) LoadAndDelete(key K) (V, bool) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	v, ok := m.Data[key]
	if ok {
		delete(m.Data, key)
	}
	return v, ok
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
// It acquires a read lock to ensure thread-safe access to the underlying data.
func (m *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	m.Mu.RLock()
	defer m.Mu.RUnlock()
	for k, v := range m.Data {
		if !f(k, v) {
			break
		}
	}
}
