package workerpool

import (
	"fmt"
	"sync"
)

type Worker struct {
	ID       int
	dataChan chan string
}

func NewWorker(ID int, ch chan string) *Worker {
	return &Worker{
		ID:       ID,
		dataChan: ch,
	}
}

func (w *Worker) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for str := range w.dataChan {
			fmt.Printf("worker id: %d; data: %s\n", w.ID, str)
		}
	}()

}
