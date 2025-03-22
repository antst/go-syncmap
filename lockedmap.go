package syncmap

// to complain if a type does not implement the required methods
var _ LockedMap[any, any] = (*lockedMap[any, any])(nil)

// LockedMap is an interface that provides safe access to the map while locked.
// It defines methods for loading, storing, deleting, and ranging over map entries.
// This interface is used internally by SyncMap to provide atomic operations.
//
// The methods in this interface assume that the caller has already acquired
// the necessary lock. Therefore, these methods should only be used within
// the context of SyncMap's DoLocked and DoLockedWithResult methods.
//
// Type parameters:
//   - K: must be a comparable type (used as map keys)
//   - V: can be any type (used as map values)
type LockedMap[K comparable, V any] interface {
	// Load retrieves the value for a key.
	// It returns the value and a boolean indicating whether the key was present.
	Load(key K) (V, bool)

	// LoadOrStore returns the existing value for the key if present.
	// Otherwise, it stores and returns the given value.
	// The loaded result is true if the value was loaded, false if stored.
	LoadOrStore(key K, value V) (V, bool)

	// Store sets the value for a key.
	Store(key K, value V)

	// LoadAndDelete removes the value for a key, returning the previous value if any.
	// The loaded result reports whether the key was present.
	LoadAndDelete(key K) (V, bool)

	// Remove deletes the value associated with the given key from the map.
	// It returns true if the key was present and removed, false otherwise.
	Remove(k K) bool

	// Purge removes all key-value pairs from the map, effectively clearing its contents.
	Purge()

	// Range calls f sequentially for each key and value present in the map.
	// If f returns false, Range stops the iteration.
	Range(f func(key K, value V) bool)

	// Filter creates a new map containing key-value pairs from the map that satisfy the given predicate function.
	// It acquires a read lock to ensure thread-safe access to the underlying data.
	Filter(predicateFn func(k K, v V) bool) map[K]V

	// Map applies a given function to all key-value pairs in the map and returns a new map with the results.
	// It acquires a read lock to ensure thread-safe access to the underlying data.
	Map(mapFn func(k K, v V) V) map[K]V

	// Len returns the number of items in the map.
	Len() int

	//  method to make sure SyncMap does not fit the LockedMap interface
	syncMap() *SyncMap[K, V]
}

// unexported type to restrict access
type lockedMap[K comparable, V any] struct {
	m *SyncMap[K, V]
}

func (lm *lockedMap[K, V]) Len() int {
	return len(lm.m.Data)
}

func (lm *lockedMap[K, V]) Load(key K) (V, bool) {
	v, ok := lm.m.Data[key]
	return v, ok
}

func (lm *lockedMap[K, V]) Store(key K, value V) {
	lm.m.Data[key] = value
}

func (lm *lockedMap[K, V]) LoadAndDelete(key K) (V, bool) {
	v, ok := lm.m.Data[key]
	if ok {
		delete(lm.m.Data, key)
	}
	return v, ok
}

func (lm *lockedMap[K, V]) Range(f func(key K, value V) bool) {
	for k, v := range lm.m.Data {
		if !f(k, v) {
			break
		}
	}
}

func (lm *lockedMap[K, V]) Purge() {
	lm.m.Data = make(map[K]V)
}

func (lm *lockedMap[K, V]) Remove(k K) bool {
	if _, ok := lm.m.Data[k]; !ok {
		return false
	}
	delete(lm.m.Data, k)
	return true
}

func (lm *lockedMap[K, V]) LoadOrStore(key K, value V) (V, bool) {

	if v, ok := lm.m.Data[key]; ok {
		return v, true
	}

	lm.m.Data[key] = value
	return value, false
}

func (lm *lockedMap[K, V]) Filter(predicateFn func(k K, v V) bool) map[K]V {
	data := make(map[K]V)

	for k, v := range lm.m.Data {
		if predicateFn(k, v) {
			data[k] = v
		}
	}

	return data
}

func (lm *lockedMap[K, V]) Map(mapFn func(k K, v V) V) map[K]V {
	data := make(map[K]V, len(lm.m.Data))

	for k, v := range lm.m.Data {
		data[k] = mapFn(k, v)
	}

	return data
}

func (lm *lockedMap[K, V]) syncMap() *SyncMap[K, V] {
	return lm.m
}
