package emissary

import (
	"fmt"
	"testing"
	"time"
)

func TestPrefork(t *testing.T) {
	d := NewPrefork(100, 2, 1*time.Second)
	for i := 0; i < 5; i++ {
		d.DoAndWait(func() {
			fmt.Println("task: ", i)
		})
	}
}
