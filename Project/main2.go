package main

import (
	"./elevator"
	"time"
	"fmt"
)

func main() {

	fmt.Println("Start")

	buttonEventCh := make(chan elevator.ButtonEvent)
	lightEventCh := make(chan elevator.LightEvent)
	stopCh := make(chan bool)
	motorStateCh := make(chan elevator.MotorDirection)
	floorCh := make(chan int)


	go elevator.EventManager(buttonEventCh, lightEventCh, stopCh, motorStateCh, floorCh)

	motorStateCh <- elevator.DirnUp

	for {
		select {
		case <- buttonEventCh:
		case <- lightEventCh:
		case <- stopCh:
		case <- motorStateCh:
		case floor := <-floorCh:
			if floor == elevator.NumFloors-1 {
				motorStateCh <- elevator.DirnDown
			} else if floor == 0 {
				motorStateCh <- elevator.DirnUp
			}
		case <-time.After(10*time.Millisecond):
		}
	}
}
