package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// --------------------------
// Интерфейс и реализация Semaphore
// --------------------------

// Semaphore — интерфейс, определяющий базовые операции
type Semaphore interface {
	Acquire()
	Release()
}

// SemaphoreImpl — реализация семафора с использованием каналов
type SemaphoreImpl struct {
	ch chan struct{}
}

// NewSemaphore создаёт новый семафор с заданной вместимостью
func NewSemaphore(n int) *SemaphoreImpl {
	return &SemaphoreImpl{
		ch: make(chan struct{}, n),
	}
}

// Acquire — захватывает ресурс, блокируя горутину, если ресурс недоступен
func (s *SemaphoreImpl) Acquire() {
	s.ch <- struct{}{}
}

// Release — освобождает ресурс
func (s *SemaphoreImpl) Release() {
	select {
	case <-s.ch:
	default:
	}
}

// --------------------------
// Моделирование критической ситуации: Deadlock
// --------------------------

// simulateDeadlock демонстрирует deadlock, когда два потока захватывают семафоры в разном порядке.
func simulateDeadlock() {
	fmt.Println("Запуск симуляции deadlock")
	sem1 := NewSemaphore(1)
	sem2 := NewSemaphore(1)

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
func simulateDeadlockResolved() {
	fmt.Println("Запуск симуляции разрешения deadlock")
	sem1 := NewSemaphore(1)
	sem2 := NewSemaphore(1)

	var wg sync.WaitGroup
	wg.Add(2)

	// Все горутины захватывают семафоры в одном и том же порядке.
	acquireInOrder := func(semA, semB *SemaphoreImpl, name string) {
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

// --------------------------
// Дополнительная задача 1: Пул ресурсов
// Ограничение числа горутин, одновременно имеющих доступ к ресурсу.
// --------------------------

type ResourcePool struct {
	semaphore Semaphore
}

func NewResourcePool(maxResources int) *ResourcePool {
	return &ResourcePool{
		semaphore: NewSemaphore(maxResources),
	}
}

// AccessResource моделирует доступ горутины к ресурсу с использованием семафора.
func (rp *ResourcePool) AccessResource(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Горутина %d ожидает ресурс\n", id)
	rp.semaphore.Acquire()
	fmt.Printf("Горутина %d получила ресурс\n", id)
	// Эмуляция работы с случайной задержкой
	time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)
	fmt.Printf("Горутина %d освобождает ресурс\n", id)
	rp.semaphore.Release()
}

// --------------------------
// Дополнительная задача 2: Читатели-писатели
// Несколько читателей могут одновременно читать, но запись выполняется эксклюзивно.
// --------------------------

type RWLock struct {
	readSemaphore  *SemaphoreImpl // для упорядочивания доступа читателей
	writeSemaphore *SemaphoreImpl // обеспечивает эксклюзивный доступ писателя
	readersCount   int
	mutex          sync.Mutex
}

func NewRWLock() *RWLock {
	return &RWLock{
		readSemaphore:  NewSemaphore(1),
		writeSemaphore: NewSemaphore(1),
		readersCount:   0,
	}
}

// RLock — захватывается при чтении.
func (rw *RWLock) RLock() {
	rw.mutex.Lock()
	rw.readersCount++
	if rw.readersCount == 1 {
		// Первый читатель блокирует писателей.
		rw.writeSemaphore.Acquire()
	}
	rw.mutex.Unlock()
}

// RUnlock — освобождение ресурса после чтения.
func (rw *RWLock) RUnlock() {
	rw.mutex.Lock()
	rw.readersCount--
	if rw.readersCount == 0 {
		rw.writeSemaphore.Release()
	}
	rw.mutex.Unlock()
}

// Lock — захват для писателя.
func (rw *RWLock) Lock() {
	// Дополнительное блокирование для предотвращения входа новых читателей.
	rw.readSemaphore.Acquire()
	// Эксклюзивный доступ.
	rw.writeSemaphore.Acquire()
}

// Unlock — освобождение писательского доступа.
func (rw *RWLock) Unlock() {
	rw.writeSemaphore.Release()
	rw.readSemaphore.Release()
}

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
	simulateDeadlockResolved()

	// Демонстрация задачи 1: Пул ресурсов.
	fmt.Println("\nЗапуск симуляции пула ресурсов")
	pool := NewResourcePool(3) // максимум 3 одновременных доступа
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
	rwLock := NewRWLock()

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
