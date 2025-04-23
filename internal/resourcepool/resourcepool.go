package resourcepool

import (
	"fmt"
	"math/rand"
	"parallel/internal/semaphore"
	"sync"
	"time"
)

// --------------------------
// Дополнительная задача 1: Пул ресурсов
// Ограничение числа горутин, одновременно имеющих доступ к ресурсу.
// --------------------------

type ResourcePool struct {
	Semaphore semaphore.Semaphore
}

func NewResourcePool(maxResources int) *ResourcePool {
	return &ResourcePool{
		Semaphore: semaphore.NewSemaphore(maxResources),
	}
}

// AccessResource моделирует доступ горутины к ресурсу с использованием семафора.
func (rp *ResourcePool) AccessResource(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Горутина %d ожидает ресурс\n", id)
	rp.Semaphore.Acquire()
	fmt.Printf("Горутина %d получила ресурс\n", id)
	// Эмуляция работы со случайной задержкой
	time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)
	fmt.Printf("Горутина %d освобождает ресурс\n", id)
	rp.Semaphore.Release()
}
