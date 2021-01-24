package input

import (
	"time"
)

func Debounce(interval time.Duration, input chan string, cb func(arg string)) {
	var curr string
	var previtem string
	timer := time.NewTimer(interval)
	for {
		select {
		case curr = <-input:
			timer.Reset(interval)
		case <-timer.C:
			if curr != "" && previtem != curr {
				cb(curr)
				previtem = curr
			}
		}
	}
}
