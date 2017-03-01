package main

import (
	"./driver"
	"time"
)

func main() {

	bCh := make(chan driver.ButtonEvent)
	lCh := make(chan driver.LightEvent)
	sCh := make(chan bool)
	msCh := make(chan driver.Elev_motor_direction_t)
	fCh := make(chan int)

	msCh <- driver.DIRN_UP

	go driver.Elevator_event_manager(bCh, lCh, sCh, msCh, fCh)

	for {
		select {
		case fi := <-fCh:
			if fi == driver.N_FLOORS-1 {
				msCh <- driver.DIRN_DOWN
			} else if fi == 0 {
				msCh <- driver.DIRN_UP
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
