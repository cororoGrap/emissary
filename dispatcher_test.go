package emissary

import (
	"testing"
)

func TestSerial(t *testing.T) {
	var st string
	d := NewSerial(10)
	d.DoAndWait(func() {
		st += "this is "
	})
	d.DoAndWait(func() {
		st += "testing"
	})

	if st != "this is testing" {
		t.Error("Wrong sequence for serial dispatcher wrong")
	}
}
