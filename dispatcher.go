package emissary

import (
	"time"
)

type dispatcher interface {
	dispatch(t task)
	start()
	DoAndWait(fn func())
}

// defaultDispatcher is a serial dispatcher
type defaultDispatcher struct {
	tasks chan task
}

func (d *defaultDispatcher) dispatch(t task) {
	d.tasks <- t
}

func (d *defaultDispatcher) init(buffer int) {
	d.tasks = make(chan task, buffer)
}

func (d *defaultDispatcher) start() {
	for {
		task := <-d.tasks
		task.Do()
		task.done()
	}
}

func (d *defaultDispatcher) DoAndWait(fn func()) {
	w := make(chan int)
	task := &simpleTask{fn: fn, re: w}
	d.dispatch(task)
	task.waitForResult()
}

type timeLimitDispatcher struct {
	defaultDispatcher
	limit int
}

func (d *timeLimitDispatcher) init(buffer int, limit int) {
	d.tasks = make(chan task, buffer)
	d.limit = limit
}

func (d *timeLimitDispatcher) start() {
	var startTime int
	var endTime int
	for {
		task := <-d.tasks
		startTime = int(time.Now().UnixNano() / int64(time.Millisecond))
		task.Do()
		endTime = int(time.Now().UnixNano() / int64(time.Millisecond))
		task.done()
		if v := d.limit - (endTime - startTime); v > 0 {
			b := v * int(time.Millisecond)
			time.Sleep(time.Duration(b))
		}
	}
}

type concurrentThrottleDispatcher struct {
	defaultDispatcher
	L      int
	lambda int
	W      int
	active int
}

func (d *concurrentThrottleDispatcher) init() {

}

func (d *concurrentThrottleDispatcher) start() {
	for {
		task := <-d.tasks
		d.doTask(task)
	}
}

func (d *concurrentThrottleDispatcher) doTask(t task) {

}
