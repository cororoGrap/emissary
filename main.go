package emissary

type Dispatcher interface {
	dispatch(t Task)
	start()
	DoAndWait(fn func())
}

type Task interface {
	Do()
	done()
	waitForResult()
}

type SimpleTask struct {
	fn func()
	re chan int
}

func (t *SimpleTask) Do() {
	t.fn()
}

func (t *SimpleTask) done() {
	t.re <- 1
}

func (t *SimpleTask) waitForResult() {
	<-t.re
}

func NewSerial(buffer int) Dispatcher {
	d := &SerialDispatcher{}
	d.init(buffer)
	go d.start()
	return d
}
