package main

import (
	"time"
	"./elevator"
	"./driver"
	"./orders"
	. "./def"
)

func main() {
	// Initialize system
	driver.Init(driver.TypeSimulation)


	// See documentation for full communication structure between main goroutines

	// Driver <--> StateMachine
	buttonEventCh := make(chan ButtonEvent, 10)
	lightEventCh := make(chan LightEvent, 10)
	stopCh := make(chan bool, 10)
	motorStateCh := make(chan MotorDirection, 10)
	floorCh := make(chan int, 10)
	doorOpenCh := make(chan bool, 10)
	floorIndicatorCh := make(chan int, 10)

	// StateMachine <--> OrderManager
	orderEventCh := make(chan OrderEvent, 10)
	stateCh := make(chan Elevator, 10)
	localOrdersCh := make(chan Orders, 10)
	globalOrdersCh := make(chan Orders, 10)

	go elevator.StateMachine(buttonEventCh, lightEventCh, stopCh, motorStateCh, floorCh, doorOpenCh, floorIndicatorCh, orderEventCh, stateCh, localOrdersCh, globalOrdersCh)
	go driver.EventManager(buttonEventCh, lightEventCh, stopCh, motorStateCh, floorCh, doorOpenCh, floorIndicatorCh)
	go orders.OrderManager(orderEventCh, stateCh, localOrdersCh, globalOrdersCh)

	for {
		select {
		case <-time.After(10*time.Millisecond):
		}
	}
}
