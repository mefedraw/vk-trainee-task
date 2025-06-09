package vk_trainee_task

import (
	"errors"
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

// AddJob добавляет job в пул, возвращает ошибку если пул закрыт или спит
func (p *Pool) AddJob(str string) error {
	if p.workerNum == 0 {
		return errors.New("zero active workers exists")
	}
	if p.isClosed {
		return errors.New("pool is closed")
	}
	p.inCh <- str
	return nil
}

// AddWorker добавляет одного воркера в пул
func (p *Pool) AddWorker() error {
	p.mu.Lock()
	if p.isClosed {
		p.mu.Unlock()
		return errors.New("pool is closed")
	}
	defer p.mu.Unlock()
	p.workerNum++
	w := NewWorker(p.workerNum)
	p.wg.Add(1)
	w.LaunchWorker(p.inCh, p.stopCh, p.wg)
	return nil
}

// RemoveWorker убирает одного воркера из пула, если их уже нуль, то ничего не происходит
func (p *Pool) RemoveWorker() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.workerNum == 0 {
		return errors.New("zero workers exists")
	}
	p.workerNum--
	p.stopCh <- struct{}{}
	return nil
}

// Stop останавливает пул и дает завершить работу оставшимся воркерам
func (p *Pool) Stop() error {
	p.mu.Lock()
	if p.isClosed {
		p.mu.Unlock()
		return errors.New("pool is closed")
	}
	close(p.inCh)
	p.isClosed = true
	p.mu.Unlock()

	p.wg.Wait()
	fmt.Println("pool stopped")
	return nil
}
