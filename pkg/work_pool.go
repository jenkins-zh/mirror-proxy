package pkg

type Task struct {
	Data     interface{}
	TaskFunc TaskFunc
}

type TaskFunc = func(data interface{})

type WorkPool struct {
	TaskChannel chan Task
	QuitChan    chan int
}

func (w *WorkPool) InitPool(size int) {
	w.TaskChannel = make(chan Task, size)
	w.QuitChan = make(chan int)
	go func() {
	DONE:
		for {
			select {
			case task := <-w.TaskChannel:
				task.Run()
			case <-w.QuitChan:
				break DONE
			}
		}
	}()
}

func (w *WorkPool) ClosePool() {
	w.QuitChan <- 1
}

func (w *WorkPool) AddTask(task Task) {
	w.TaskChannel <- task
}

func (t *Task) Run() {
	t.TaskFunc(t.Data)
}
