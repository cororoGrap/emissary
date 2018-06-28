package emissary

import (
	"time"
)

func NewSerial(buffer int) dispatcher {
	d := &defaultDispatcher{}
	d.init(buffer)
	go d.start()
	return d
}

func NewPrefork(buffer, max int, timeout time.Duration) dispatcher {
	d := &PreforkPool{}
	d.init(buffer, max, timeout)
	go d.start()
	return d
}
