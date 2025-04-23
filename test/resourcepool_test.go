package test

import (
	"math/rand"
	"parallel/internal/resourcepool"
	"sync"
	"testing"
	"time"
)

func TestResourcePool(t *testing.T) {
	maxResources := 3
	pool := resourcepool.NewResourcePool(maxResources)
	totalGoroutines := 20
	var wg sync.WaitGroup
	counterCh := make(chan int, totalGoroutines)

	var currentConcurrent int
	var concurrentMutex sync.Mutex

	for i := 0; i < totalGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Захватываем ресурс
			pool.Semaphore.Acquire()
			// Увеличиваем счётчик одновременного доступа
			concurrentMutex.Lock()
			currentConcurrent++
			if currentConcurrent > maxResources {
				t.Errorf("Превышен лимит одновременных доступов: %d", currentConcurrent)
			}
			concurrentMutex.Unlock()

			time.Sleep(time.Duration(rand.Intn(100)+50) * time.Millisecond)

			concurrentMutex.Lock()
			currentConcurrent--
			concurrentMutex.Unlock()
			pool.Semaphore.Release()
			counterCh <- 1
		}(i)
	}
	wg.Wait()
	close(counterCh)
	count := 0
	for range counterCh {
		count++
	}
	if count != totalGoroutines {
		t.Errorf("Не все горутины завершились: ожидалось %d, получили %d", totalGoroutines, count)
	}
}
