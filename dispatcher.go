package emissary

import (
	"time"
)

type TimeLimitDispatcher struct {
	jobs  chan Task
	limit int
}

func (d *TimeLimitDispatcher) init(buffer int, limit int) {
	d.jobs = make(chan Task, buffer)
	d.limit = limit
}

func (d *TimeLimitDispatcher) dispatch(t Task) {
	d.jobs <- t
}

func (d *TimeLimitDispatcher) start() {
	var startTime int
	var endTime int
	for {
		job := <-d.jobs
		startTime = int(time.Now().UnixNano() / int64(time.Millisecond))
		job.Do()
		endTime = int(time.Now().UnixNano() / int64(time.Millisecond))
		job.done()
		if v := d.limit - (endTime - startTime); v > 0 {
			b := v * int(time.Millisecond)
			time.Sleep(time.Duration(b))
		}
	}
}

func (d *TimeLimitDispatcher) DoAndWait(fn func()) {
	w := make(chan int)
	task := &SimpleTask{fn: fn, re: w}
	d.dispatch(task)
	task.waitForResult()
}

type SerialDispatcher struct {
	TimeLimitDispatcher
}

func (d *SerialDispatcher) init(buffer int) {
	d.jobs = make(chan Task, buffer)
	d.limit = 0
}
