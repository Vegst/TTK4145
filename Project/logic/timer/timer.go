package timer

import (
	"time"
)

func Timer(resetCh chan time.Duration, timeoutCh chan bool) {
	var d time.Duration = 0
	for {
		select {
		case d = <- resetCh:
		case <- time.After(100*time.Millisecond):
			if (d > 0) {
				if (d > 100*time.Millisecond) {
					d -= 100*time.Millisecond
				} else {
					d = 0
					timeoutCh <- true
				}
			}
		}
	}
}
