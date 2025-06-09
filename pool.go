package vk_trainee_task

import (
	"fmt"
	"sync"
)

type WorkerLauncher interface {
	LaunchWorker(in chan string, stopCh chan struct{})
}

// Pool - структура пула
type Pool struct {
	inCh      chan string
	stopCh    chan struct{}
	wg        *sync.WaitGroup
	mu        sync.Mutex
	workerNum int
	isClosed  bool
}

// NewPool создает новый объект пула, но не запускает его
func NewPool(wg *sync.WaitGroup) *Pool {
	return &Pool{
		inCh:      make(chan string),
		stopCh:    make(chan struct{}),
		workerNum: 0,
		mu:        sync.Mutex{},
		wg:        wg,
	}
}

// Run принимает начальное кол-во воркеров, которых в последствии запускает
func (p *Pool) Run(n int) {
	for i := 1; i <= n; i++ {
		p.AddWorker()
	}
}

// AddJob добавляет job в пул, паникует если пул закрыт или спит
func (p *Pool) AddJob(str string) {
	if p.workerNum == 0 {
		panic("pool not started")
	}
	if p.isClosed {
		panic("pool is closed")
	}
	p.inCh <- str
}

// AddWorker добавляет одного воркера в пул
func (p *Pool) AddWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.workerNum++
	w := NewWorker(p.workerNum)
	p.wg.Add(1)
	w.LaunchWorker(p.inCh, p.stopCh, p.wg)
}

// RemoveWorker убирает одного воркера из пула, если их уже нуль, то ничего не происходит
func (p *Pool) RemoveWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.workerNum == 0 {
		return
	}
	p.workerNum--
	p.stopCh <- struct{}{}
}

// Stop останавливает пул и дает завершить работу оставшимся воркерам. Если пул уже закрыт, то паникует
func (p *Pool) Stop() {
	p.mu.Lock()
	if p.isClosed {
		p.mu.Unlock()
		panic("Pool already stopped")
	}
	close(p.inCh)
	p.isClosed = true
	p.mu.Unlock()

	p.wg.Wait()
	fmt.Println("pool stopped")
}
