package emissary

type task interface {
	Do()
	done()
	waitForResult()
}

type simpleTask struct {
	fn func()
	re chan int
}

func (t *simpleTask) Do() {
	t.fn()
}

func (t *simpleTask) done() {
	t.re <- 1
}

func (t *simpleTask) waitForResult() {
	<-t.re
}
