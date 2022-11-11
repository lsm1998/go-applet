package task_pool

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	taskPool := New(2)
	go func() {
		for i := 0; i < 10; i++ {
			taskPool.Submit(func() {
				time.Sleep(time.Second)
				fmt.Println("完成一次")
			})
		}
	}()
	time.Sleep(time.Second)
	taskPool.Finish()
}
