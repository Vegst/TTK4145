package logic

import(
	"time"
	"../elevator"
)

const NumOrderTypes = 3
const NumFloors = 4


type State struct{
	Floor int
	Direction elevator.MotorDirection
}

buttonEventCh := make(chan elevator.ButtonEvent)
lightEventCh := make(chan elevator.LightEvent)
stopCh := make(chan bool)
motorStateCh := make(chan elevator.MotorDirection)
floorCh := make(chan int)

orderCh := make(chan OrderEvent)
orderEventCh := make(chan OrderEvent)

var localOrders[NumFloors][NumOrderTypes] bool
var state State


go elevator.EventManager(buttonEventCh, lightEventCh, stopCh, motorStateCh, floorCh)

for {
	select {
	case be := <-buttonEventCh:
		orderEventch <- OrderEvent{be.Floor, be.Button, true}
	case <- lightEventCh:
	case <- stopCh:

	case <- motorStateCh:
	case floor := <-floorCh:

		for b := 0; b < NumOrderTypes; b++{
			orderEventch <- OrderEvent{floor, b, false}
		}


		
	case LocalOrders = <- orderCh:
	case <-time.After(10*time.Millisecond):
	}
}