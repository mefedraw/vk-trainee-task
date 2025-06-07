package workerpool

import "sync"

type Pool struct {
	Strings   []string
	workerNum int
	collector chan string
	wg        *sync.WaitGroup
}

func NewPool(workerNum int, strings []string) *Pool {
	return &Pool{
		Strings:   strings,
		workerNum: workerNum,
		collector: make(chan string, 1000),
		wg:        &sync.WaitGroup{},
	}
}

func (p *Pool) AddWorker(amount int) {
	p.workerNum += amount
}

func (p *Pool) DecrementWorker(amount int) {
	p.workerNum -= amount
}

func (p *Pool) Run() {
	for i := 1; i <= p.workerNum; i++ {
		worker := NewWorker(i, p.collector)
		worker.Start(p.wg)
	}
	for i := range p.Strings {
		p.collector <- p.Strings[i]
	}
	close(p.collector)
	p.wg.Wait()
}
