package logic

import(
	"time"
	"../elevator"
	"../orders"
)

const NumOrderTypes = 3
const NumFloors = 4


type State struct{
	Floor int
	Direction elevator.MotorDirection
}

var localOrders orders.Orders
var state State

func StateMachine(){
	buttonEventCh := make(chan elevator.ButtonEvent)
	lightEventCh := make(chan elevator.LightEvent)
	stopCh := make(chan bool)
	motorStateCh := make(chan elevator.MotorDirection)
	floorCh := make(chan int)

	orderCh := make(chan OrderEvent)
	orderEventCh := make(chan OrderEvent)

	go elevator.EventManager(buttonEventCh, lightEventCh, stopCh, motorStateCh, floorCh)

	for {
		select {
		case be := <-buttonEventCh:
			orderEventch <- OrderEvent{be.Floor, be.Button, true}
		case <- stopCh:

		case floor := <-floorCh:
			state.Floor = floor
			lightEventCh <- elevator.LightType{LIGHT_TYPE_FLOOR, state.floor, true}

			stop := ShouldStop(localOrder, state)
			if(stop){
				if state.Direction == elevator.DirnUp {
					orderEventCh <- OrderEvent{floor, orders.OrderCallUp, false}
				}
				else if state.Direction == elevator.DirnDown {
					orderEventCh <- OrderEvent{floor, orders.OrderCallDown, false}
				}
				orderEventCh <- OrderEvent{floor, orders.OrderCommand, false}
				motorStateCh <- stop

			}
			motorStateCh <- GetDirection(localOrders, state)
			
		case LocalOrders = <- orderCh:
		case <-time.After(10*time.Millisecond):
		}
	}
}
