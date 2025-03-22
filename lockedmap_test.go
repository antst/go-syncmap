package syncmap

import (
	"sort"
	"testing"
)

func TestLockedMap(t *testing.T) {
	sm := New[string, int](10)
	var lm LockedMap[string, int]

	sm.DoLocked(
		func(m LockedMap[string, int]) {
			lm = m
			lm.Store("key1", 1)
			lm.Store("key2", 2)
			lm.Store("key3", 3)
		},
	)

	t.Run(
		"Load", func(t *testing.T) {
			sm.DoLocked(
				func(m LockedMap[string, int]) {
					if v, ok := m.Load("key1"); !ok || v != 1 {
						t.Errorf("Expected 1, got %v", v)
					}
					if _, ok := m.Load("non-existent"); ok {
						t.Error("Load should return false for non-existent key")
					}
				},
			)
		},
	)

	t.Run(
		"Store", func(t *testing.T) {
			sm.DoLocked(
				func(m LockedMap[string, int]) {
					m.Store("key4", 4)
					if v, ok := m.Load("key4"); !ok || v != 4 {
						t.Errorf("Expected 4, got %v", v)
					}
				},
			)
		},
	)

	t.Run(
		"LoadAndDelete", func(t *testing.T) {
			sm.DoLocked(
				func(m LockedMap[string, int]) {
					if v, ok := m.LoadAndDelete("key2"); !ok || v != 2 {
						t.Errorf("Expected 2, got %v", v)
					}
					if _, ok := m.Load("key2"); ok {
						t.Error("Key should not exist after LoadAndDelete")
					}
				},
			)
		},
	)

	t.Run(
		"Range", func(t *testing.T) {
			keys := make([]string, 0)
			values := make([]int, 0)

			sm.DoLocked(
				func(m LockedMap[string, int]) {
					m.Range(
						func(key string, value int) bool {
							keys = append(keys, key)
							values = append(values, value)
							return true
						},
					)
				},
			)

			expectedKeys := []string{"key1", "key3", "key4"}
			expectedValues := []int{1, 3, 4}

			sort.Strings(keys)
			sort.Ints(values)

			if !slicesEqual(keys, expectedKeys) || !slicesEqual(values, expectedValues) {
				t.Errorf("Range didn't return expected results. Got keys: %v, values: %v", keys, values)
			}
		},
	)

	t.Run(
		"Len", func(t *testing.T) {
			sm.DoLocked(
				func(m LockedMap[string, int]) {
					if m.Len() != 3 {
						t.Errorf("Expected length 3, got %d", m.Len())
					}
				},
			)
		},
	)

	t.Run(
		"Purge", func(t *testing.T) {
			sm.DoLocked(
				func(m LockedMap[string, int]) {
					m.Purge()
					if m.Len() != 0 {
						t.Errorf("Expected length 0 after Purge, got %d", m.Len())
					}
				},
			)
		},
	)

	t.Run(
		"Remove", func(t *testing.T) {
			sm.DoLocked(
				func(m LockedMap[string, int]) {
					m.Store("key1", 1)
					m.Store("key2", 2)

					if !m.Remove("key1") {
						t.Error("Remove should return true for existing key")
					}
					if m.Remove("non-existent") {
						t.Error("Remove should return false for non-existent key")
					}
					if _, ok := m.Load("key1"); ok {
						t.Error("Key should not exist after Remove")
					}
					if m.Len() != 1 {
						t.Errorf("Expected length 1 after Remove, got %d", m.Len())
					}
				},
			)
		},
	)

	t.Run(
		"LoadOrStore", func(t *testing.T) {
			sm.DoLocked(
				func(m LockedMap[string, int]) {
					// Store a new value
					value, loaded := m.LoadOrStore("key3", 3)
					if loaded {
						t.Error("LoadOrStore should not have loaded for a new key")
					}
					if value != 3 {
						t.Errorf("Expected 3, got %v", value)
					}

					// Try to store an existing value
					value, loaded = m.LoadOrStore("key3", 4)
					if !loaded {
						t.Error("LoadOrStore should have loaded for an existing key")
					}
					if value != 3 {
						t.Errorf("Expected 3, got %v", value)
					}

					// Verify the value wasn't changed
					if v, ok := m.Load("key3"); !ok || v != 3 {
						t.Errorf("Expected 3, got %v", v)
					}
				},
			)
		},
	)
	t.Run(
		"Filter", func(t *testing.T) {
			sm.DoLocked(
				func(m LockedMap[string, int]) {
					m.Purge()
					m.Store("key1", 1)
					m.Store("key2", 2)
					m.Store("key3", 3)
					m.Store("key4", 4)

					filtered := m.Filter(
						func(k string, v int) bool {
							return v%2 == 0
						},
					)

					if len(filtered) != 2 {
						t.Errorf("Expected 2 items after filtering, got %d", len(filtered))
					}

					expectedFiltered := map[string]int{"key2": 2, "key4": 4}
					for k, v := range expectedFiltered {
						if filteredV, ok := filtered[k]; !ok || filteredV != v {
							t.Errorf("Expected filtered map to contain %s: %d, but got %v", k, v, filteredV)
						}
					}
				},
			)
		},
	)

	t.Run(
		"Map", func(t *testing.T) {
			sm.DoLocked(
				func(m LockedMap[string, int]) {
					m.Purge()
					m.Store("key1", 1)
					m.Store("key2", 2)
					m.Store("key3", 3)

					mapped := m.Map(
						func(k string, v int) int {
							return v * 2
						},
					)

					if len(mapped) != 3 {
						t.Errorf("Expected 3 items after mapping, got %d", len(mapped))
					}

					expectedMapped := map[string]int{"key1": 2, "key2": 4, "key3": 6}
					for k, v := range expectedMapped {
						if mappedV, ok := mapped[k]; !ok || mappedV != v {
							t.Errorf("Expected mapped map to contain %s: %d, but got %v", k, v, mappedV)
						}
					}
				},
			)
		},
	)

}
