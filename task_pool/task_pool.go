package task_pool

import (
	"context"
	"sync"
)

type Task = func()

type TaskPool interface {
	Submit(Task)

	TrySubmit(Task) bool

	Finish()

	TaskNumber() int
}

type channelTaskPool struct {
	taskC    chan Task
	cap      int32
	cancel   context.CancelFunc
	ctx      context.Context
	finishWG sync.WaitGroup
}

func New(cap int32) TaskPool {
	pool := &channelTaskPool{
		cap:   cap,
		taskC: make(chan Task, cap),
	}
	pool.ctx, pool.cancel = context.WithCancel(context.Background())
	pool.process()
	return pool
}

func (c *channelTaskPool) process() {
	c.finishWG.Add(int(c.cap))
	for i := 0; i < int(c.cap); i++ {
		go func() {
			defer c.finishWG.Done()
			for {
				select {
				case task := <-c.taskC:
					task()
				case <-c.ctx.Done():
					return
				}
			}
		}()
	}
}

func (c *channelTaskPool) TaskNumber() int {
	return len(c.taskC)
}

func (c *channelTaskPool) Submit(task Task) {
	c.taskC <- task
	return
}

func (c *channelTaskPool) TrySubmit(task Task) bool {
	select {
	case c.taskC <- task:
		return true
	default:
		return false
	}
}

func (c *channelTaskPool) Finish() {
	c.cancel()
	c.finishWG.Wait()
}
