package syncmap

import (
	"fmt"
)

func Example() {
	// Create a new SyncMap with an initial size of 10
	sm := New[string, int](10)

	// Store key-value pairs
	sm.Store("apple", 5)
	sm.Store("banana", 3)
	sm.Store("cherry", 8)

	// Load a value
	if value, ok := sm.Load("banana"); ok {
		fmt.Printf("Value of 'banana': %d\n", value)
	}

	// Use LoadOrStore to add a new key-value pair
	if value, loaded := sm.LoadOrStore("date", 6); !loaded {
		fmt.Printf("Stored new value for 'date': %d\n", value)
	}

	// Use LoadOrStore with an existing key
	if value, loaded := sm.LoadOrStore("apple", 10); loaded {
		fmt.Printf("Loaded existing value for 'apple': %d\n", value)
	}

	// Remove a key-value pair
	sm.Remove("cherry")

	// Get the length of the map
	fmt.Printf("Number of items: %d\n", sm.Len())

	// Use the Map function to double all values
	doubledMap := sm.Map(
		func(k string, v int) int {
			return v * 2
		},
	)
	fmt.Printf("Doubled map: %v\n", doubledMap)

	// Use the Filter function to keep only values greater than or equal to 5
	filteredMap := sm.Filter(
		func(k string, v int) bool {
			return v >= 5
		},
	)
	fmt.Printf("Filtered map: %v\n", filteredMap)

	// Use LoadAndDelete
	if value, ok := sm.LoadAndDelete("apple"); ok {
		fmt.Printf("Loaded and deleted 'apple': %d\n", value)
	}

	// Use Range to iterate over remaining items
	sm.Range(
		func(key string, value int) bool {
			fmt.Printf("Key: %s, Value: %d\n", key, value)
			return true
		},
	)

	// Use DoLocked to perform multiple operations atomically
	sm.DoLocked(
		func(m LockedMap[string, int]) {
			m.Store("grape", 7)
			m.Store("kiwi", 4)
		},
	)

	// Use DoLockedWithResult to perform operations and return a result
	result := sm.DoLockedWithResult(
		func(m LockedMap[string, int]) any {
			return m.Len()
		},
	)
	fmt.Printf("Number of items after DoLocked: %d\n", result)

	// Purge the map
	sm.Purge()
	fmt.Printf("After purge, number of items: %d\n", sm.Len())

	// Output:
	// Value of 'banana': 3
	// Stored new value for 'date': 6
	// Loaded existing value for 'apple': 5
	// Number of items: 3
	// Doubled map: map[apple:10 banana:6 date:12]
	// Filtered map: map[apple:5 date:6]
	// Loaded and deleted 'apple': 5
	// Key: banana, Value: 3
	// Key: date, Value: 6
	// Number of items after DoLocked: 4
	// After purge, number of items: 0
}
