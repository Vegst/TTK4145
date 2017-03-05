package logic

import(
	"time"
	"./timer"
	//"fmt"
	"./elevator"
)


type State struct {
	Floor int
	Direction elevator.MotorDirection
	DoorOpen bool
}

var localOrders [NumFloors][NumTypes] bool
var state State

func StateMachine() {
	
	buttonEventCh := make(chan elevator.ButtonEvent, 10)
	lightEventCh := make(chan elevator.LightEvent, 10)
	stopCh := make(chan bool, 10)
	motorStateCh := make(chan elevator.MotorDirection, 10)
	floorCh := make(chan int, 10)

	orderEventCh := make(chan OrderEvent, 10)
	stateCh := make(chan State, 10)
	localOrdersCh := make(chan [NumFloors][NumTypes] bool, 10)
	globalOrdersCh := make(chan [NumFloors][NumTypes] bool, 10)

	timerResetCh := make(chan time.Duration)
	timerTimeoutCh := make(chan bool)

	go timer.Timer(timerResetCh, timerTimeoutCh)
	go OrderManager(orderEventCh, stateCh, localOrdersCh, globalOrdersCh)
	go elevator.EventManager(buttonEventCh, lightEventCh, stopCh, motorStateCh, floorCh)
	
	motorStateCh <- elevator.DirnUp
	state.Direction = elevator.DirnUp


	for {
		select {
		case be := <-buttonEventCh:
			orderEventCh <- OrderEvent{be.Floor, OrderType(be.Button), true}
		case <- stopCh:

		case state.Floor = <-floorCh:
			if state.Floor >= 0 && state.Floor < NumFloors {
				lightEventCh <- elevator.LightEvent{elevator.LIGHT_TYPE_FLOOR, state.Floor, true}


				stop := ShouldStop(localOrders, state)

				if(stop){
					
					if state.Direction == elevator.DirnUp {
						orderEventCh <- OrderEvent{state.Floor, OrderCallUp, false}
					} else if state.Direction == elevator.DirnDown {
						orderEventCh <- OrderEvent{state.Floor, OrderCallDown, false}
					}
					orderEventCh <- OrderEvent{state.Floor, OrderCallCommand, false}

					state.Direction = elevator.DirnStop
					lightEventCh <- elevator.LightEvent{elevator.LIGHT_TYPE_DOOR, state.Floor, true}
					motorStateCh <- state.Direction
					state.DoorOpen = true
					timerResetCh <- time.Second * 3
				} else {
					state.Direction = GetDirection(localOrders, state)
					motorStateCh <- state.Direction
				}
				
			}
			
		case localOrders = <- localOrdersCh:
			if state.Direction == elevator.DirnStop && !state.DoorOpen {
				state.Direction = GetDirection(localOrders, state)
				motorStateCh <- state.Direction
			}
		case globalOrders := <- globalOrdersCh:
			for f := 0; f < NumFloors; f++ {
				for b := 0; b < NumTypes; b++ {
					lightEventCh <- elevator.LightEvent{elevator.LightType(b), f, globalOrders[f][b]}
				}
			}
		case <-timerTimeoutCh:
			state.Direction = GetDirection(localOrders, state)
			state.DoorOpen = false
			motorStateCh <- state.Direction
			lightEventCh <- elevator.LightEvent{elevator.LIGHT_TYPE_DOOR, state.Floor, false}

		case <- time.After(10*time.Millisecond):
		}
	}
}
