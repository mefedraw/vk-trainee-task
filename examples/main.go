package main

import (
	"github.com/brianvoe/gofakeit"
	vk_trainee_task "github.com/mefedraw/vk-trainee-task"
	"runtime"
	"sync"
)

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	wg := &sync.WaitGroup{}

	maxWorkers := numCPU * 2
	workers := 3
	jobsCnt := 10000

	wp := vk_trainee_task.NewPool(maxWorkers, wg)
	wp.Run(workers)

	for i := 0; i < jobsCnt; i++ {
		wp.AddJob(gofakeit.Username())
	}

	wp.Stop()
}
