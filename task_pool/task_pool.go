package task_pool

type Task = func()

type TaskPool interface {
	// Submit 提交任务，在任务进入pool之前是阻塞的
	Submit(Task)

	// TrySubmit 提交任务，如果阻塞则返回false
	TrySubmit(Task) bool

	// Finish 完成，执行后拒绝Submit
	Finish()

	// TaskNumber 获取当前任务执行数量
	TaskNumber() int
}
