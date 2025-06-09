package vk_trainee_task

import (
	"fmt"
	"sync"
)

// Worker структура воркера
type Worker struct {
	ID int
}

// NewWorker создает воркера с ID
func NewWorker(id int) *Worker {
	return &Worker{
		ID: id,
	}
}

// LaunchWorker запускает воркера
func (w *Worker) LaunchWorker(in <-chan string, stopCh chan struct{}, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stopCh:
				fmt.Printf("worker №%d stopped\n", w.ID)
				return
			case str, ok := <-in:
				if !ok {
					fmt.Printf("worker №%d stopped\n", w.ID)
					return
				}
				w.Process(str)
			}
		}
	}()
}

// Process обрабатывает job
func (w *Worker) Process(str string) {
	fmt.Printf("string %s was processed by worker №%d\n", str, w.ID)
}
