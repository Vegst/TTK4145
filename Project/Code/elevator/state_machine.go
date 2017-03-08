package elevator

import (
	"./timer"
	"time"
	//"fmt"
	"../driver"
)

const NumElevators = 3

type Behaviour int

const (
	BehaviourIdle     = 0
	BehaviourMoving   = 1
	BehaviourDoorOpen = 2
)

type Elevator struct {
	Floor     int
	Direction driver.MotorDirection
	Behaviour Behaviour
	Orders    Orders
}

/*
   States:
	* Idle
	* Moving
	* DoorOpen
   Events:
	* Button
	* Stop
	* Floor
	* Local Orders
	* Global Orders
*/

func StateMachine() {

	// Init channels
	buttonEventCh := make(chan driver.ButtonEvent, 10)
	lightEventCh := make(chan driver.LightEvent, 10)
	stopCh := make(chan bool, 10)
	motorStateCh := make(chan driver.MotorDirection, 10)
	floorCh := make(chan int, 10)
	doorOpenCh := make(chan bool, 10)
	floorIndicatorCh := make(chan int, 10)

	orderEventCh := make(chan OrderEvent, 10)
	stateCh := make(chan Elevator, 10)
	localOrdersCh := make(chan [NumFloors][NumTypes]bool, 10)
	globalOrdersCh := make(chan [NumFloors][NumTypes]bool, 10)

	timerResetCh := make(chan time.Duration)
	timerTimeoutCh := make(chan bool)

	var elev Elevator

	go timer.Timer(timerResetCh, timerTimeoutCh)
	go OrderManager(orderEventCh, stateCh, localOrdersCh, globalOrdersCh, elev)
	go driver.EventManager(buttonEventCh, lightEventCh, stopCh, motorStateCh, floorCh, doorOpenCh, floorIndicatorCh)

	// Initial state change
	motorStateCh <- driver.DirnUp
	elev.Direction = driver.DirnUp
	elev.Behaviour = BehaviourMoving

	for {
		select {
		case buttonEvent := <-buttonEventCh:
			if buttonEvent.State {
				orderEventCh <- OrderEvent{buttonEvent.Floor, OrderType(buttonEvent.Button), true}
			}
		case <-stopCh:

		case elev.Floor = <-floorCh:
			if elev.Floor >= 0 && elev.Floor < NumFloors {
				floorIndicatorCh <- elev.Floor

				switch elev.Behaviour {
				case BehaviourMoving:
					if ShouldStop(elev.Orders, elev) {
						if OrderAtFloor(elev.Orders, elev.Floor) {
							// Clear orders at current floor
							if elev.Direction == driver.DirnUp {
								orderEventCh <- OrderEvent{elev.Floor, OrderCallUp, false}
							} else if elev.Direction == driver.DirnDown {
								orderEventCh <- OrderEvent{elev.Floor, OrderCallDown, false}
							}
							orderEventCh <- OrderEvent{elev.Floor, OrderCallCommand, false}

							doorOpenCh <- true
							timerResetCh <- time.Second * 3
							motorStateCh <- driver.DirnStop

							elev.Behaviour = BehaviourDoorOpen
						} else {
							motorStateCh <- driver.DirnStop
							elev.Behaviour = BehaviourIdle
							elev.Direction = driver.DirnStop
						}

					}
				}

			}

		case elev.Orders = <-localOrdersCh:
			switch elev.Behaviour {
			case BehaviourDoorOpen:
				if OrderAtFloor(elev.Orders, elev.Floor) {
					orderEventCh <- OrderEvent{elev.Floor, OrderCallUp, false}
					orderEventCh <- OrderEvent{elev.Floor, OrderCallDown, false}
					orderEventCh <- OrderEvent{elev.Floor, OrderCallCommand, false}
					timerResetCh <- time.Second * 3
				}
			case BehaviourIdle:
				if OrderAtFloor(elev.Orders, elev.Floor) {
					orderEventCh <- OrderEvent{elev.Floor, OrderCallUp, false}
					orderEventCh <- OrderEvent{elev.Floor, OrderCallDown, false}
					orderEventCh <- OrderEvent{elev.Floor, OrderCallCommand, false}
					timerResetCh <- time.Second * 3
					doorOpenCh <- true
					elev.Behaviour = BehaviourDoorOpen

				} else {
					elev.Direction = GetDirection(elev.Orders, elev)
					if elev.Direction == driver.DirnStop {
						elev.Behaviour = BehaviourIdle
					} else {
						elev.Behaviour = BehaviourMoving
					}
					motorStateCh <- elev.Direction
				}
			}
		case globalOrders := <-globalOrdersCh:
			for f := 0; f < NumFloors; f++ {
				for b := 0; b < NumTypes; b++ {
					lightEventCh <- driver.LightEvent{driver.LightType(b), f, globalOrders[f][b]}
				}
			}
		case <-timerTimeoutCh:
			switch elev.Behaviour {
			case BehaviourDoorOpen:
				elev.Direction = GetDirection(elev.Orders, elev)
				if elev.Direction == driver.DirnStop {
					elev.Behaviour = BehaviourIdle
				} else {
					elev.Behaviour = BehaviourMoving
				}
				motorStateCh <- elev.Direction
				doorOpenCh <- false
			}

		case <-time.After(10 * time.Millisecond):
		}
	}
}
