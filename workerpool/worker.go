package workerpool

import (
	"fmt"
	"sync"
)

type Worker struct {
	ID int
	wg *sync.WaitGroup
}

func NewWorker(id int, wg *sync.WaitGroup) *Worker {
	return &Worker{
		ID: id,
		wg: wg,
	}
}

func (w *Worker) LaunchWorker(in chan string, stopCh chan struct{}) {
	go func() {
		defer w.wg.Done()
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

func (w *Worker) Process(str string) {
	fmt.Printf("string %s was processed by worker №%d\n", str, w.ID)
}
