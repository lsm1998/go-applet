package task_pool

import (
	"context"
	"golang.org/x/sync/semaphore"
	"sync/atomic"
)

func init() {
	semaphore.NewWeighted(10)
}

type semaphoreTaskPool struct {
	weighted *semaphore.Weighted
	number   int32
	finish   bool
}

func NewSemaphoreTaskPool(cap int32) TaskPool {
	return &semaphoreTaskPool{
		weighted: semaphore.NewWeighted(int64(cap)),
	}
}

func (s *semaphoreTaskPool) Submit(task Task) {
	if s.finish {
		return
	}
	_ = s.weighted.Acquire(context.Background(), 1)
	go func() {
		atomic.AddInt32(&s.number, 1)
		defer s.release()
		task()
	}()
}

func (s *semaphoreTaskPool) TrySubmit(task Task) bool {
	if !s.finish && s.weighted.TryAcquire(1) {
		go func() {
			atomic.AddInt32(&s.number, 1)
			defer s.release()
			task()
		}()
		return true
	} else {
		return false
	}
}

func (s *semaphoreTaskPool) release() {
	atomic.AddInt32(&s.number, -1)
	s.weighted.Release(1)
}

func (s *semaphoreTaskPool) Finish() {
	s.finish = true
}

func (s *semaphoreTaskPool) TaskNumber() int {
	return int(atomic.LoadInt32(&s.number))
}
