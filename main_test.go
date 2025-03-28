package main

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestSimulateDeadlockResolved(t *testing.T) {
	done := make(chan struct{})
	go func() {
		simulateDeadlockResolved()
		close(done)
	}()
	select {
	case <-done:
		// Успех: функция завершилась без deadlock.
	case <-time.After(2 * time.Second):
		t.Error("simulateDeadlockResolved завершилась с таймаутом – возможен deadlock")
	}
}

func TestResourcePool(t *testing.T) {
	maxResources := 3
	pool := NewResourcePool(maxResources)
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
			pool.semaphore.Acquire()
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
			pool.semaphore.Release()
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

func TestRWLock(t *testing.T) {
	rwLock := NewRWLock()
	var wg sync.WaitGroup
	readersCount := 5
	writerStarted := make(chan int, 1)
	writerFinished := make(chan int, 1)

	// Запуск читателей.
	for i := 0; i < readersCount; i++ {
		wg.Add(1)
		go func(reader int) {
			defer wg.Done()
			rwLock.RLock()
			// Если писатель уже начал, чтение не должно происходить одновременно с записью.
			select {
			case <-writerStarted:
				t.Error("Читатель начал чтение, когда писатель активен")
			default:
			}
			time.Sleep(time.Duration(rand.Intn(100)+50) * time.Millisecond)
			rwLock.RUnlock()
		}(i)
	}

	// Запуск писателя.
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(50 * time.Millisecond) // даём читателям стартовать
		writerStarted <- 1
		rwLock.Lock()
		time.Sleep(150 * time.Millisecond)
		rwLock.Unlock()
		writerFinished <- 1
	}()

	wg.Wait()
	select {
	case <-writerFinished:
		// Писатель завершил работу – тест успешен.
	default:
		t.Error("Писатель не завершил работу как ожидалось")
	}
}
