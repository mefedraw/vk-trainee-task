package vk_trainee_task

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"runtime"
	"sync"
	"testing"
)

func initPool() *Pool {
	wg := &sync.WaitGroup{}
	wp := NewPool(wg)
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	return wp
}

func captureStdout(f func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestPool_addJob(t *testing.T) {
	wp := initPool()
	wp.Run(3)

	output := captureStdout(func() {
		for i := 1; i <= 5; i++ {
			wp.AddJob(fmt.Sprintf("string#%d", i))
		}
		wp.Stop()
	})

	assert.NotEmpty(t, output)
	assert.Contains(t, output, "string string#1 was processed by worker №")
	assert.Contains(t, output, "string string#2 was processed by worker №")
	assert.Contains(t, output, "string string#3 was processed by worker №")
	assert.Contains(t, output, "string string#4 was processed by worker №")
	assert.Contains(t, output, "string string#5 was processed by worker №")
}

func TestPool_AddJob_Fail(t *testing.T) {
	wp := initPool()
	wp.Run(3)

	output := captureStdout(func() {
		for i := 1; i <= 2; i++ {
			wp.AddJob(fmt.Sprintf("string#%d", i))
		}
		wp.Stop()
	})

	assert.NotContains(t, output, "string string#3 was processed by worker №")
}

func TestPool_Run(t *testing.T) {
	wp := initPool()
	wp.Run(3)
	wp.AddJob("string#1")
	wp.Stop()

	assert.Equal(t, wp.workerNum, 3)
}

func TestPool_Run_PanicsWhenNotStarted(t *testing.T) {
	wp := initPool()

	defer func() {
		if r := recover(); r == nil {
			assert.FailNow(t, "pool did not panic")
		}
	}()
	wp.AddJob("string#1")
}

func TestPool_AddWorker(t *testing.T) {
	wp := initPool()
	wp.Run(3)
	wp.AddWorker()
	assert.Equal(t, wp.workerNum, 4)
	wp.AddWorker()
	assert.Equal(t, wp.workerNum, 5)
	wp.AddWorker()
	wp.Stop()

	assert.Equal(t, wp.workerNum, 6)
}

func TestPool_RemoveWorker(t *testing.T) {
	wp := initPool()
	wp.Run(5)
	wp.RemoveWorker()
	assert.Equal(t, wp.workerNum, 4)
	wp.RemoveWorker()
	assert.Equal(t, wp.workerNum, 3)
	wp.Stop()
}

func TestPool_Stop(t *testing.T) {
	wp := initPool()
	wp.Run(3)

	defer func() {
		if r := recover(); r == nil {
			assert.FailNow(t, "pool did not panic")
		}
	}()
	
	output := captureStdout(func() {
		wp.Stop()
		wp.AddJob("string#1")
	})

	assert.NotContains(t, output, "string string#1 was processed by worker")
	assert.Contains(t, output, "pool stopped")
}
