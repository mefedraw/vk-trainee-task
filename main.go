package main

import (
	"github.com/brianvoe/gofakeit"
	"runtime"
	"sync"
	"vk-trainee-task/workerpool"
)

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	wg := &sync.WaitGroup{}

	maxWorkers := numCPU * 2
	workers := 3
	jobsCnt := 3000

	wp := workerpool.NewPool(maxWorkers, wg)
	wp.Run(workers)

	for i := 0; i < jobsCnt; i++ {
		wp.AddJob(gofakeit.Username())
	}

	wp.Stop()
}
