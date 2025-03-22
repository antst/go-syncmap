package syncmap

import (
	"fmt"
	"sort"
	"sync"
	"testing"
)

func TestSyncMap(t *testing.T) {
	sm := New[string, int](10)

	t.Run(
		"Store and Load", func(t *testing.T) {
			sm.Store("key1", 1)
			sm.Store("key2", 2)

			if v, ok := sm.Load("key1"); !ok || v != 1 {
				t.Errorf("Expected 1, got %v", v)
			}
			if v, ok := sm.Load("key2"); !ok || v != 2 {
				t.Errorf("Expected 2, got %v", v)
			}
		},
	)

	t.Run(
		"LoadOrStore", func(t *testing.T) {
			sm := New[string, int](10)

			// Store a new value
			value, loaded := sm.LoadOrStore("key1", 1)
			if loaded {
				t.Error("LoadOrStore should not have loaded for a new key")
			}
			if value != 1 {
				t.Errorf("Expected 1, got %v", value)
			}

			// Try to store an existing value
			value, loaded = sm.LoadOrStore("key1", 2)
			if !loaded {
				t.Error("LoadOrStore should have loaded for an existing key")
			}
			if value != 1 {
				t.Errorf("Expected 1, got %v", value)
			}

			// Verify the value wasn't changed
			if v, ok := sm.Load("key1"); !ok || v != 1 {
				t.Errorf("Expected 1, got %v", v)
			}
		},
	)

	t.Run(
		"Remove", func(t *testing.T) {
			if !sm.Remove("key1") {
				t.Error("Remove should return true for existing key")
			}
			if sm.Remove("non-existent") {
				t.Error("Remove should return false for non-existent key")
			}
			if _, ok := sm.Load("key1"); ok {
				t.Error("Key should not exist after removal")
			}
		},
	)

	t.Run(
		"Map", func(t *testing.T) {
			sm.Store("key3", 3)
			sm.Store("key4", 4)

			doubled := sm.Map(
				func(k string, v int) int {
					return v * 2
				},
			)

			expected := map[string]int{"key2": 4, "key3": 6, "key4": 8}
			if !mapsEqual(doubled, expected) {
				t.Errorf("Expected %v, got %v", expected, doubled)
			}
		},
	)

	t.Run(
		"Filter", func(t *testing.T) {
			filtered := sm.Filter(
				func(k string, v int) bool {
					return v > 2
				},
			)

			expected := map[string]int{"key3": 3, "key4": 4}
			if !mapsEqual(filtered, expected) {
				t.Errorf("Expected %v, got %v", expected, filtered)
			}
		},
	)

	t.Run(
		"Len", func(t *testing.T) {
			if sm.Len() != 3 {
				t.Errorf("Expected length 3, got %d", sm.Len())
			}
		},
	)

	t.Run(
		"LoadAndDelete", func(t *testing.T) {
			v, ok := sm.LoadAndDelete("key3")
			if !ok || v != 3 {
				t.Errorf("Expected 3, got %v", v)
			}
			if _, ok := sm.Load("key3"); ok {
				t.Error("Key should not exist after LoadAndDelete")
			}
		},
	)

	t.Run(
		"Range", func(t *testing.T) {
			keys := make([]string, 0)
			values := make([]int, 0)

			sm.Range(
				func(key string, value int) bool {
					keys = append(keys, key)
					values = append(values, value)
					return true
				},
			)

			expectedKeys := []string{"key2", "key4"}
			expectedValues := []int{2, 4}

			sort.Strings(keys)
			sort.Ints(values)

			if !slicesEqual(keys, expectedKeys) || !slicesEqual(values, expectedValues) {
				t.Errorf("Range didn't return expected results")
			}
		},
	)

	t.Run(
		"DoLocked", func(t *testing.T) {
			sm.DoLocked(
				func(m LockedMap[string, int]) {
					m.Store("key5", 5)
					m.Store("key6", 6)
				},
			)

			if v, ok := sm.Load("key5"); !ok || v != 5 {
				t.Errorf("Expected 5, got %v", v)
			}
			if v, ok := sm.Load("key6"); !ok || v != 6 {
				t.Errorf("Expected 6, got %v", v)
			}
		},
	)

	t.Run(
		"DoLockedWithResult", func(t *testing.T) {
			result := sm.DoLockedWithResult(
				func(m LockedMap[string, int]) any {
					return m.Len()
				},
			)

			if result != 4 {
				t.Errorf("Expected length 4, got %v", result)
			}
		},
	)

	t.Run(
		"Purge", func(t *testing.T) {
			sm.Purge()
			if sm.Len() != 0 {
				t.Errorf("Expected length 0 after purge, got %d", sm.Len())
			}
		},
	)

	t.Run(
		"Concurrent access", func(t *testing.T) {
			const goroutines = 100
			const iterations = 1000

			var wg sync.WaitGroup
			wg.Add(goroutines)

			for i := 0; i < goroutines; i++ {
				go func(id int) {
					defer wg.Done()
					for j := 0; j < iterations; j++ {
						key := fmt.Sprintf("key%d-%d", id, j)
						sm.Store(key, j)
						if _, ok := sm.Load(key); !ok {
							t.Errorf("Failed to load stored value")
						}
						sm.Remove(key)
					}
				}(i)
			}

			wg.Wait()

			if sm.Len() != 0 {
				t.Errorf("Expected length 0 after concurrent operations, got %d", sm.Len())
			}
		},
	)
}

func mapsEqual[K comparable, V comparable](m1, m2 map[K]V) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || v1 != v2 {
			return false
		}
	}
	return true
}

func slicesEqual[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
