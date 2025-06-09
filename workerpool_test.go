package vk_trainee_task

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
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

func TestPool_AddJob_ProcessesEnqueuedJobs(t *testing.T) {
	wp := initPool()
	wp.Run(3)

	output := captureStdout(func() {
		for i := 1; i <= 5; i++ {
			err := wp.AddJob(fmt.Sprintf("string#%d", i))
			require.NoError(t, err)
		}
		err := wp.Stop()
		require.NoError(t, err)
	})

	require.NotEmpty(t, output)
	require.Contains(t, output, "string string#1 was processed by worker №")
	require.Contains(t, output, "string string#2 was processed by worker №")
	require.Contains(t, output, "string string#3 was processed by worker №")
	require.Contains(t, output, "string string#4 was processed by worker №")
	require.Contains(t, output, "string string#5 was processed by worker №")
}

func TestPool_AddJob_LimitsToSubmittedJobs(t *testing.T) {
	wp := initPool()
	wp.Run(3)

	output := captureStdout(func() {
		for i := 1; i <= 2; i++ {
			err := wp.AddJob(fmt.Sprintf("string#%d", i))
			require.NoError(t, err)
		}
		err := wp.Stop()
		require.NoError(t, err)
	})

	require.NotContains(t, output, "string string#3 was processed by worker №")
}

func TestPool_AddJob_ReturnsErrorAfterStop(t *testing.T) {
	wp := initPool()
	wp.Run(3)
	err := wp.Stop()
	require.NoError(t, err)
	err = wp.AddJob("string#1")
	require.Error(t, err)
	require.ErrorContains(t, err, "pool is closed")
}

func TestPool_Run_InitializesWorkersAndAllowsJobs(t *testing.T) {
	wp := initPool()
	wp.Run(3)
	err := wp.AddJob("string#1")
	require.NoError(t, err)
	err = wp.Stop()
	require.NoError(t, err)
	require.Equal(t, wp.workerNum, 3)
}

func TestPool_AddJob_ErrorWhenNotStarted(t *testing.T) {
	wp := initPool()
	err := wp.AddJob("string#1")
	require.Error(t, err)
	require.Equal(t, wp.workerNum, 0)
	require.ErrorContains(t, err, "zero active workers exists")
}

func TestPool_AddWorker_IncrementsWorkerCount(t *testing.T) {
	wp := initPool()
	wp.Run(3)
	err := wp.AddWorker()
	require.NoError(t, err)
	require.Equal(t, wp.workerNum, 4)
	err = wp.AddWorker()
	require.NoError(t, err)
	require.Equal(t, wp.workerNum, 5)
	err = wp.AddWorker()
	require.NoError(t, err)
	err = wp.Stop()

	require.NoError(t, err)
	require.Equal(t, wp.workerNum, 6)
}

func TestPool_RemoveWorker_DecrementsWorkerCount(t *testing.T) {
	wp := initPool()
	wp.Run(5)
	err := wp.RemoveWorker()
	require.Equal(t, wp.workerNum, 4)
	require.NoError(t, err)
	err = wp.RemoveWorker()
	require.NoError(t, err)
	require.Equal(t, wp.workerNum, 3)
	err = wp.Stop()
	require.NoError(t, err)
}

func TestPool_RemoveWorker_ErrorWhenNoWorkers(t *testing.T) {
	wp := initPool()
	err := wp.RemoveWorker()
	require.Error(t, err)
	require.ErrorContains(t, err, "zero workers exists")
	require.Equal(t, wp.workerNum, 0)
}

func TestPool_Stop_StopsWorkersAndPreventsNewJobs(t *testing.T) {
	wp := initPool()
	wp.Run(3)

	output := captureStdout(func() {
		err := wp.Stop()
		require.NoError(t, err)
		err = wp.AddJob("string#1")
		require.Error(t, err)
		require.ErrorContains(t, err, "pool is closed")
	})

	require.NotContains(t, output, "string string#1 was processed by worker")
	require.Contains(t, output, "pool stopped")
}

func TestPool_Stop_ErrorOnSecondCall(t *testing.T) {
	wp := initPool()
	wp.Run(2)
	err := wp.Stop()
	require.NoError(t, err)
	err = wp.Stop()
	require.Error(t, err)
	require.ErrorContains(t, err, "pool is closed")
}
