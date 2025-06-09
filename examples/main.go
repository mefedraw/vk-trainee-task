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

	workers := 5
	jobsCnt := 1000

	wp := vk_trainee_task.NewPool(wg)
	wp.Run(workers)

	for i := 0; i < jobsCnt; i++ {
		wp.AddJob(gofakeit.Username())
	}

	wp.Stop()
	wp.Stop()
}
