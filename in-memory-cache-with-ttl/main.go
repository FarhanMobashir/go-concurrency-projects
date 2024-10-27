package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type CacheItem struct {
	Value      interface{}
	Expiration time.Time
}

type Cache struct {
	items           map[string]CacheItem
	mu              sync.RWMutex
	cleanupInterval time.Duration
	stopCleanup     chan struct{}
	hitCount        int64 // atomic counter for hits
	missCount       int64 // atomic counter for misses
}

func (c *Cache) GetHitCount() int64 {
	return atomic.LoadInt64(&c.hitCount)
}

func (c *Cache) GetMissCount() int64 {
	return atomic.LoadInt64(&c.missCount)
}

func NewCache(cleanupInterval time.Duration) *Cache {
	cache := &Cache{
		items:           make(map[string]CacheItem),
		cleanupInterval: cleanupInterval,
		stopCleanup:     make(chan struct{}),
	}

	go cache.startCleanup()
	return cache
}

// start cleanup runs in the backgroud to periodically remove the expired items
func (c *Cache) startCleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.deleteExpiredItems()
		case <-c.stopCleanup:
			return
		}
	}
}

// deleteExpiredItems removes expired items from the cache
func (c *Cache) deleteExpiredItems() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, item := range c.items {
		if time.Now().After(item.Expiration) {
			delete(c.items, key)
		}
	}
}

// Close stops the background cleanup process.
func (c *Cache) Close() {
	close(c.stopCleanup)
}

// set method on the Cache
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	expiration := time.Now().Add(ttl)

	// lock the cache for writing
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Value:      value,
		Expiration: expiration,
	}
}

// get method on cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	if !found {
		atomic.AddInt64(&c.missCount, 1) // Increment hit count
		return nil, false
	}

	// check for expiry
	if time.Now().After(item.Expiration) {
		// if expired remove from the cache
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		return nil, false
	}

	atomic.AddInt64(&c.hitCount, 1) // Increment hit count
	return item.Value, true

}

func main() {
	// Initialize the cache with a cleanup interval of 2 seconds
	cache := NewCache(2 * time.Second)

	// Define the number of workers and items to set/get
	const numWorkers = 100
	const itemsPerWorker = 5
	var wg sync.WaitGroup

	// Start goroutines to set values in the cache
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < itemsPerWorker; j++ {
				key := fmt.Sprintf("worker%d_item%d", workerID, j)
				value := fmt.Sprintf("value%d", workerID*j)
				cache.Set(key, value, 3*time.Second) // Each item expires in 3 seconds
				fmt.Printf("Worker %d set %s = %s\n", workerID, key, value)
				time.Sleep(100 * time.Millisecond) // Sleep to simulate work
			}
		}(i)
	}

	// Start goroutines to get values from the cache
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < itemsPerWorker; j++ {
				key := fmt.Sprintf("worker%d_item%d", workerID, j)
				if value, found := cache.Get(key); found {
					fmt.Printf("Worker %d got %s = %s\n", workerID, key, value)
				} else {
					fmt.Printf("Worker %d did not find %s\n", workerID, key)
				}
				time.Sleep(150 * time.Millisecond) // Sleep to simulate work
			}
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Print hit and miss counts
	fmt.Printf("Total hits: %d\n", cache.GetHitCount())
	fmt.Printf("Total misses: %d\n", cache.GetMissCount())

	// Cleanup
	cache.Close()
	fmt.Println("Cache cleanup stopped.")
}
