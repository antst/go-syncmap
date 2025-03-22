# SyncMap

SyncMap is a high-performance, thread-safe map implementation in Go that provides a concurrent-safe wrapper around the built-in Go map type. It ensures safe access in multi-threaded environments using efficient read-write locking mechanisms.

## Features

- Thread-safe operations
- Generic implementation
- Efficient read-write locking
- Support for common map operations (Load, Store, Remove, Filter, etc.)
- Optimized for high concurrency scenarios

## Installation

To use SyncMap in your Go project, you can install it using `go get`:

```bash
go get github.com/antst/go-syncmap
```

## Usage

Import the package in your Go code:

```go
import "github.com/antst/go-syncmap"
```

Create a new SyncMap:

```go
m := syncmap.New[string, int]()
```

Perform operations:

```go
// Set a value
m.Store("key", 42)

// Get a value
value, exists := m.Load("key")
if exists {
    fmt.Printf("Value: %d\n", value)
}

// Delete a key
m.Remove("key")

// Get the number of items in the map
count := m.Len()

// Use the Map function to double all values
doubledMap := m.Map(
func(k string, v int) int {return v * 2})

// Use the Filter function to keep only values greater than or equal to 5
filteredMap := m.Filter(
func(k string, v int) bool {return v >= 5})

// Perform some complex tasks while in the locked state and return value
// LockedMap is an interface to the internal map[K]V, while in locked scope.
result := m.DoLockedWithResult(
	func(m LockedMap[string, int]) any {
		    ....
			// do a lot of stuff while in locked scope
			.....
		return some_value
	},
)
```

## Performance

SyncMap is designed to provide high performance in concurrent scenarios. It uses a read-write mutex to allow multiple simultaneous reads while ensuring exclusive access for writes.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the BSD 3-Clause License - see the [LICENSE](LICENSE) file for details.

## Author

Anton Starikov
