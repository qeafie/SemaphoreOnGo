package semaphore

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
