package vk_trainee_task

import "testing"

func benchMarkWorkerPool(t *testing.B) {
	wp := NewPool(5)
}
