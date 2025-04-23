package main

import (
	"fmt"
	"math/rand"
	"parallel/internal/deadlock"
	"parallel/internal/resourcepool"
	"parallel/internal/rwlock"
	"sync"
	"time"
)

// --------------------------
// Основная функция: демонстрация работы всех компонентов
// --------------------------

func main() {
	rand.Seed(time.Now().UnixNano())

	// Для критичной ситуации:
	// ВНИМАНИЕ: функция simulateDeadlock может привести к зависанию из-за deadlock.
	// Раскомментируйте, чтобы увидеть поведение.
	// simulateDeadlock()

	// Решённый вариант без deadlock:
	deadlock.SimulateDeadlockResolved()

	// Демонстрация задачи 1: Пул ресурсов.
	fmt.Println("\nЗапуск симуляции пула ресурсов")
	pool := resourcepool.NewResourcePool(3) // максимум 3 одновременных доступа
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go pool.AccessResource(i, &wg)
	}
	wg.Wait()
	fmt.Println("Симуляция пула ресурсов завершена")

	// Демонстрация задачи 2: Читатели-писатели.
	fmt.Println("\nЗапуск симуляции схемы читатели-писатели")
	var rwWg sync.WaitGroup
	rwLock := rwlock.NewRWLock()

	// Моделирование нескольких читателей.
	for i := 0; i < 5; i++ {
		rwWg.Add(1)
		go func(reader int) {
			defer rwWg.Done()
			rwLock.RLock()
			fmt.Printf("Читатель %d читает данные\n", reader)
			time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)
			fmt.Printf("Читатель %d завершил чтение\n", reader)
			rwLock.RUnlock()
		}(i)
	}

	// Моделирование писателя.
	rwWg.Add(1)
	go func() {
		defer rwWg.Done()
		time.Sleep(50 * time.Millisecond) // даём читателям время начать
		rwLock.Lock()
		fmt.Println("Писатель записывает данные")
		time.Sleep(300 * time.Millisecond)
		fmt.Println("Писатель завершил запись")
		rwLock.Unlock()
	}()
	rwWg.Wait()
	fmt.Println("Симуляция схемы читатели-писатели завершена")
}
