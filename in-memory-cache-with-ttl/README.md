# Concurrent In-Memory Cache with TTL

This project implements a concurrent in-memory cache in Go (Golang) that supports key-value storage with a Time-To-Live (TTL) feature. It serves as a practical learning resource to explore various concepts related to concurrency, data structures, and memory management in Go.

## Overview

The cache is designed to store items as key-value pairs with an expiration mechanism. Items are automatically removed after a specified TTL, ensuring that the cache remains efficient and relevant.

## Key Concepts

1. **Key-Value Store**:

   - The cache stores data as key-value pairs, allowing for quick data retrieval based on unique keys.
   - This structure facilitates efficient lookup, insertion, and deletion operations.

2. **TTL (Time-To-Live)**:

   - Each cache item is associated with an expiration time, defined by a TTL.
   - Once the TTL expires, the item is considered stale and is automatically removed from the cache during periodic cleanup.

3. **Concurrency**:

   - The cache is designed to handle concurrent access from multiple goroutines, using:
     - **Mutexes**:
       - `sync.RWMutex` is implemented to allow multiple concurrent reads while ensuring exclusive writes. This enhances performance in read-heavy workloads.
     - **Atomic Operations**:
       - `sync/atomic` is used for counters that track cache hits and misses, enabling thread-safe updates without the overhead of locking.

4. **Background Cleanup**:

   - A dedicated goroutine runs periodically to remove expired items from the cache.
   - This cleanup mechanism helps maintain cache efficiency by ensuring that memory is freed up and stale data is not retained.

5. **Hit and Miss Counters**:
   - The cache maintains atomic counters for tracking successful data retrievals (hits) and unsuccessful lookups (misses).
   - This metric is crucial for understanding cache performance and effectiveness.

## Concurrency Design

The design supports concurrent operations through the following mechanisms:

- **Read-Write Locks**:

  - Using `sync.RWMutex`, multiple goroutines can read from the cache simultaneously while writes are serialized. This approach balances concurrency and data integrity.

- **Atomic Counters**:
  - The use of `sync/atomic` ensures that updates to hit and miss counts are performed atomically, preventing race conditions and ensuring accurate statistics.

### Example Usage

Multiple goroutines can interact with the cache simultaneously, as shown in the following pattern:

```go
for i := 0; i < numWorkers; i++ {
    go func(workerID int) {
        // Set cache items
        // Retrieve cache items
    }(i)
}
```

This demonstrates how the cache can handle concurrent operations, making it suitable for high-performance applications.

## Performance Considerations

The performance of the cache is influenced by:

- **Hit-to-Miss Ratio**:

  - A higher ratio of hits to misses indicates an effective cache, reducing the need for expensive data retrieval operations.

- **Cleanup Interval**:
  - The frequency of the background cleanup process can impact overall performance. Adjusting the cleanup interval and TTL values based on application requirements can help optimize resource usage.

## Conclusion

This concurrent in-memory cache implementation in Go provides a solid foundation for understanding key concepts in concurrency, data management, and performance optimization. It is a valuable resource for learning how to build efficient, thread-safe data structures in Go.
