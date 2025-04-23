package rwlock

import (
	"parallel/internal/semaphore"
	"sync"
)

// --------------------------
// Дополнительная задача 2: Читатели-писатели
// Несколько читателей могут одновременно читать, но запись выполняется эксклюзивно.
// --------------------------

type RWLock struct {
	readSemaphore  *semaphore.SemaphoreImpl // для упорядочивания доступа читателей
	writeSemaphore *semaphore.SemaphoreImpl // обеспечивает эксклюзивный доступ писателя
	readersCount   int
	mutex          sync.Mutex
}

func NewRWLock() *RWLock {
	return &RWLock{
		readSemaphore:  semaphore.NewSemaphore(1),
		writeSemaphore: semaphore.NewSemaphore(1),
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
