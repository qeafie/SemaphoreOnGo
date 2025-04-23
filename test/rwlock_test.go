package test

import (
	"math/rand"
	"parallel/internal/rwlock"
	"sync"
	"testing"
	"time"
)

func TestRWLock(t *testing.T) {
	rwLock := rwlock.NewRWLock()
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
