
# mycacheengine
`mycacheengine` is a custom caching library implemented in Go, designed to optimize data retrieval and improve application performance by temporarily storing frequently accessed data in memory. It features thread-safe operations and employs an N-way set associative caching mechanism.

## Features
-  **In-Memory Storage**: Utilizes Go's `container/list` for efficient data storage and retrieval.
-  **Configurable Capacity**: Allows setting a maximum cache size to control memory usage.
-  **Automatic Eviction**: Implements strategies to remove the least recently used (LRU) or most recently used (MRU) items when the cache reaches its capacity. The eviction policy (LRU or MRU) is defined when the cache instance is initialized, defaulting to LRU if no specific algorithm is specified.
-  **Data type flexibility**: This implementation allows saving any data type, from primitive to more complex data types.
-  **Thread-Safe Operations**: Ensures safe concurrent access using mutex locks.


## Installation
To integrate `mycacheengine` into your Go project, ensure you have Go installed and set up. Then, run:
```bash
go  get  github.com/azlancpool/mycacheengine
```
## Usage

Here's a basic example of how to use `mycacheengine`:
```go
package main

import (
	"fmt"
	"log"

	"github.com/azlancpool/mycacheengine/cache"
)

func main() {
	// It allows you to define a specific strategy for refreshing cached data.
	// To define the MRU strategy, initialize the cache as follows:
	// 		- cache.NewCache[int, any](5, cache.MRU_ALGO)
	// To explicitly define LRU, you can use:
	// 		- cache.NewCache[int, any](5, cache.LRU_ALGO)
	cache, err := cache.NewCache[int, any](5) // Initialize the cache with a capacity of 5. By default, it will use the LRU algorithm.
	if err != nil {
		log.Fatal("couldn't initialize cache, error: ", err)
	}

	// Add items to the cache
	cache.Put(123, "value1")
	cache.Put(456, true)

	// Retrieve items from the cache
	if value, found := cache.Get(456); found {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Key not found")
	}
}

```
## Thread-Safe Functionality
`mycacheengine` ensures thread safety by utilizing mutex locks (`sync.Mutex`) to manage concurrent access to the cache. This design prevents race conditions and ensures data integrity when multiple goroutines interact with the cache simultaneously. The mutex is locked during write operations and unlocked upon completion, allowing safe concurrent reads and writes.
  

## Internal Mechanics: N-Way Set Associative Cache
`mycacheengine` employs an N-way set associative caching mechanism, which balances the simplicity of direct-mapped caches with the flexibility of fully associative caches. In this design, the cache is divided into multiple sets, each containing a fixed number of ways (slots). Once the set is determined, the data is aggregated in the next available slot. If no more spaces are available, one of the previously defined algorithms (LRU or MRU) is implemented to remove the corresponding information.

**Difference between LRU and MRU**

**LRU - Least Recently Used**
- **How Eviction Policy Works:** For explanation purposes let's assume that all this data is going to be saved in the same set.
```bash
# Initial data
set = []
set_capacity = 4
input = [1,7,9,15,9,7,45]

## 1º iteration - current input: [[1],7,9,15,9,7,45]
current_input_index = 0
current_input_value = 1
# set has enough space, so it saves the provided value in the next available space
set = [] => set = [1]

## 2º iteration - current input: [1,[7],9,15,9,7,45]
current_input_index = 1
current_input_value = 7
# set has enough space, so it saves the provided value in the next available space
set = [1] => set = [1,7]

## 3º iteration - current input: [1,7,[9],15,9,7,45]
current_input_index = 2
current_input_value = 9
# set has enough space, so it saves the provided value in the next available space
set = [1,7] => set = [1,7,9]

## 4º iteration - current input: [1,7,9,[15],9,7,45]
current_input_index = 3
current_input_value = 15
# set has enough space, so it saves the provided value in the next available space
set = [1,7,9] => set = [1,7,9,15]

## 5º iteration - current input: [1,7,9,15,[9],7,45]
current_input_index = 4
current_input_value = 9
# set has NOT enough space, check if the value is already in the set and move it to the front
set = [1,7,9,15] => set = [1,7,15,9]

## 6º iteration - current input: [1,7,9,15,9,[7],45]
current_input_index = 5
current_input_value = 7
# set has NOT enough space, check if the value is already in the set and move it to the front
set = [1,7,15,9] => set = [1,15,9,7]

## 7º iteration - current input: [1,7,9,15,9,7,[45]]
current_input_index = 6
current_input_value = 45
# set has NOT enough space, If the value is not in the set. It removes the Least Recently Used item, in this case: 1, having the next final result
set = [1,15,9,7] => set = [15,9,7,45]

#Pd. This specific scenario could be validated in unit tests, scenario: LRU - SHOW CASE
```

**MRU - Most Recently Used**
- **How Eviction Policy Works:** For explanation purposes let's assume that all this data is going to be saved in the same set. In order to see the difference between LRU and MRU we are going to use the same input and set size.
```bash
# Initial data
set = []
set_capacity = 4
input = [1,7,9,15,9,7,45]

## 1º iteration - current input: [[1],7,9,15,9,7,45]
current_input_index = 0
current_input_value = 1
# set has enough space, so it saves the provided value in the next available space
set = [] => set = [1]

## 2º iteration - current input: [1,[7],9,15,9,7,45]
current_input_index = 1
current_input_value = 7
# set has enough space, so it saves the provided value in the next available space
set = [1] => set = [1,7]

## 3º iteration - current input: [1,7,[9],15,9,7,45]
current_input_index = 2
current_input_value = 9
# set has enough space, so it saves the provided value in the next available space
set = [1,7] => set = [1,7,9]

## 4º iteration - current input: [1,7,9,[15],9,7,45]
current_input_index = 3
current_input_value = 15
# set has enough space, so it saves the provided value in the next available space
set = [1,7,9] => set = [1,7,9,15]

## 5º iteration - current input: [1,7,9,15,[9],7,45]
current_input_index = 4
current_input_value = 9
# set has NOT enough space, check if the value is already in the set and move it to the front
set = [1,7,9,15] => set = [1,7,15,9]

## 6º iteration - current input: [1,7,9,15,9,[7],45]
current_input_index = 5
current_input_value = 7
# set has NOT enough space, check if the value is already in the set and move it to the front
set = [1,7,15,9] => set = [1,15,9,7]

## 7º iteration - current input: [1,7,9,15,9,7,[45]]
current_input_index = 6
current_input_value = 45
# set has NOT enough space, If the value is not in the set. It removes the Most Recently Used item, in this case: 7, having the next final result:
set = [1,15,9,7] => set = [1,15,9,45]

#Pd. This specific scenario could be validated in unit tests, scenario: MRU - SHOW CASE
``` 

**Advantages:**
-  **Reduced Collisions:** By allowing multiple slots per set, the cache reduces the likelihood of collisions compared to direct-mapped caches.
-  **Balanced Performance:** Offers a compromise between the fast access of direct-mapped caches and the flexibility of fully associative caches.

**Considerations:**
-  **Memory Overhead:** Requires additional memory to maintain metadata (e.g., usage history for eviction policies).

## Suitable Use Cases
### LRU
- Web browsers (caching recently accessed pages).
- Database query caching.
- Memory management in OS page replacement.

### MRU
- Audio/video streaming buffers.
- Real-time applications where the last accessed data is temporary.
- Workloads where the newest data is often discarded quickly.

## Strengths
-  **Speed:** Offers rapid data retrieval due to in-memory storage and efficient data structures.
-  **Thread Safety:** Ensures safe concurrent access, making it suitable for multi-threaded applications.
-  **Data type flexibility**: This implementation allows saving any data type, from primitive to more complex data types.
  
## Limitations
-  **Volatility:** As an in-memory cache, all stored data is lost upon application termination or crash.
-  **Memory Constraints:** The cache size is limited by the available system memory; storing large amounts of data may lead to increased memory consumption.

## Contributing
Contributions are welcome! Please fork the repository and submit a pull request with your enhancements or bug fixes. I’d really appreciate it if you could take a look at breaking down my code. Do you have any edge cases in mind? Let’s give it a try!

## Changelog
See the full changelog in [CHANGELOG](CHANGELOG.md).

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for.
