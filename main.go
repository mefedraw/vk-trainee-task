package main

import (
	"github.com/brianvoe/gofakeit"
	"vk-trainee-task/workerpool"
)

func main() {
	var allStrings []string
	for i := 0; i < 10000; i++ {
		allStrings = append(allStrings, gofakeit.Username())
	}

	pool := workerpool.NewPool(5, allStrings)
	pool.Run()
}
