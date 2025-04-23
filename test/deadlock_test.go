package test

import (
	"parallel/internal/deadlock"
	"testing"
	"time"
)

func TestSimulateDeadlockResolved(t *testing.T) {
	done := make(chan struct{})
	go func() {
		deadlock.SimulateDeadlockResolved()
		close(done)
	}()
	select {
	case <-done:
		// Успех: функция завершилась без deadlock.
	case <-time.After(2 * time.Second):
		t.Error("simulateDeadlockResolved завершилась с таймаутом – возможен deadlock")
	}
}
