package related

import "sync"

var (
	orderCounter int
	counterMutex sync.Mutex
)

// generateUniqueID generates a unique integer ID using a counter
func GenerateUniqueID() int {
	counterMutex.Lock()
	defer counterMutex.Unlock()

	// Increment the counter
	orderCounter++

	return orderCounter
}
