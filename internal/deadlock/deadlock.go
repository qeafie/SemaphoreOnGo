package deadlock

import (
	"fmt"
	"math/rand"
	"parallel/internal/semaphore"
	"sync"
	"time"
)

// --------------------------
// Моделирование критической ситуации: Deadlock
// --------------------------

// simulateDeadlock демонстрирует deadlock, когда два потока захватывают семафоры в разном порядке.
func SimulateDeadlock() {
	fmt.Println("Запуск симуляции deadlock")
	sem1 := semaphore.NewSemaphore(1)
	sem2 := semaphore.NewSemaphore(1)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		fmt.Println("Горутина 1: захватывает sem1")
		sem1.Acquire()
		time.Sleep(100 * time.Millisecond) // эмуляция работы
		fmt.Println("Горутина 1: пытается захватить sem2")
		sem2.Acquire()
		fmt.Println("Горутина 1: захвачены оба семафора")
		sem2.Release()
		sem1.Release()
	}()

	go func() {
		defer wg.Done()
		fmt.Println("Горутина 2: захватывает sem2")
		sem2.Acquire()
		time.Sleep(100 * time.Millisecond) // эмуляция работы
		fmt.Println("Горутина 2: пытается захватить sem1")
		sem1.Acquire()
		fmt.Println("Горутина 2: захвачены оба семафора")
		sem1.Release()
		sem2.Release()
	}()

	wg.Wait()
	fmt.Println("Симуляция deadlock завершена (если функция не зависла, то deadlock не произошёл, но обычно ожидание остаётся бесконечным)")
}

// simulateDeadlockResolved решает проблему deadlock посредством единообразного порядка захвата семафоров.
func SimulateDeadlockResolved() {
	fmt.Println("Запуск симуляции разрешения deadlock")
	sem1 := semaphore.NewSemaphore(1)
	sem2 := semaphore.NewSemaphore(1)

	var wg sync.WaitGroup
	wg.Add(2)

	// Все горутины захватывают семафоры в одном и том же порядке.
	acquireInOrder := func(semA, semB *semaphore.SemaphoreImpl, name string) {
		fmt.Printf("%s: захватывает sem1\n", name)
		semA.Acquire()
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		fmt.Printf("%s: захватывает sem2\n", name)
		semB.Acquire()
		fmt.Printf("%s: захвачены оба семафора\n", name)
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		semB.Release()
		semA.Release()
	}

	go func() {
		defer wg.Done()
		acquireInOrder(sem1, sem2, "Горутина 1")
	}()
	go func() {
		defer wg.Done()
		acquireInOrder(sem1, sem2, "Горутина 2")
	}()
	wg.Wait()
	fmt.Println("Симуляция разрешения deadlock завершена")
}
