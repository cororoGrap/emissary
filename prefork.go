package emissary

import "time"

type PreforkPool struct {
	defaultDispatcher
	max     int
	timeout time.Duration
	workerQ chan chan task
}

func (p *PreforkPool) init(buffer, max int, timeout time.Duration) {
	p.max = max
	p.timeout = timeout
	p.tasks = make(chan task, buffer)
	p.workerQ = make(chan chan task, max)
	require := p.max - len(p.workerQ)
	for i := 0; i < require; i++ {
		p.fork()
	}
}

func (p *PreforkPool) start() {
	tt := p.timeout / time.Duration(p.max)
	timer := time.Tick(tt)
	for {
		select {
		case ch := <-p.workerQ:
			task := <-p.tasks
			ch <- task
		case <-timer:
			if len(p.workerQ) < p.max {
				p.fork()
			}
		}
	}
}

func (p *PreforkPool) dispatch(t task) {
	p.tasks <- t
}

func (p *PreforkPool) fork() {
	ch := make(chan task)
	go worker(ch)
	p.workerQ <- ch
}

func worker(ch <-chan task) {
	task := <-ch
	task.Do()
	task.done()
}
