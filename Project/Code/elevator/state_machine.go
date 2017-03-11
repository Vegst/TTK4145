package elevator

import (
	. "../def"
	"./timer"
	"time"
	//"fmt"
)


func StateMachine(buttonEventCh chan ButtonEvent, lightEventCh chan LightEvent, stopCh chan bool, motorStateCh chan MotorDirection, floorCh chan int, doorOpenCh chan bool, floorIndicatorCh chan int, orderEventCh chan OrderEvent, stateCh chan ElevatorState, localOrdersCh chan elev.Orders, globalOrdersCh chan elev.Orders) {

	timerResetCh := make(chan time.Duration)
	timerTimeoutCh := make(chan bool)
	go timer.Timer(timerResetCh, timerTimeoutCh)

	var elev Elevator

	// Initial state change
	motorStateCh <- DirnUp
	elev.Direction = DirnUp
	elev.Behaviour = ElevatorBehaviourMoving

	for {
		select {
		// Event: Button pressed
		case buttonEvent := <-buttonEventCh:
			if buttonEvent.State {
				oe := OrderEvent{buttonEvent.Floor, OrderType(buttonEvent.Button), true}
				orderEventCh <- oe
				CalculateCost(oe, elev)
			}
		// Event: Stop command
		case <-stopCh:
		// Event: Floor reached
		case elev.Floor = <-floorCh:
			if elev.Floor >= 0 && elev.Floor < NumFloors {
				floorIndicatorCh <- elev.Floor

				switch elev.Behaviour {
				case ElevatorBehaviourMoving:
					if(ShouldStop(elev)){
						if OrderAtFloor(elev.Orders, elev.Floor) {
							// Clear elev.Orders at current floor
							if elev.Direction == DirnUp {
								orderEventCh <- OrderEvent{elev.Floor, OrderCallUp, false}
							} else if elev.Direction == DirnDown {
								orderEventCh <- OrderEvent{elev.Floor, OrderCallDown, false}
							}
							orderEventCh <- OrderEvent{elev.Floor, OrderCallCommand, false}

							doorOpenCh <- true
							timerResetCh <- time.Second * 3
							motorStateCh <- DirnStop

							elev.Behaviour = ElevatorBehaviourDoorOpen
						} else {
							motorStateCh <- DirnStop
							elev.Behaviour = ElevatorBehaviourIdle
							elev.Direction = DirnStop
						}
					}
				}
			}

		case elev.Orders = <-localOrdersCh:
			switch elev.Behaviour {
			case ElevatorBehaviourDoorOpen:
				if OrderAtFloor(elev.Orders, elev.Floor) {
					orderEventCh <- OrderEvent{elev.Floor, OrderCallUp, false}
					orderEventCh <- OrderEvent{elev.Floor, OrderCallDown, false}
					orderEventCh <- OrderEvent{elev.Floor, OrderCallCommand, false}
					timerResetCh <- time.Second * 3
				}
			case ElevatorBehaviourIdle:
				if OrderAtFloor(elev.Orders, elev.Floor) {
					orderEventCh <- OrderEvent{elev.Floor, OrderCallUp, false}
					orderEventCh <- OrderEvent{elev.Floor, OrderCallDown, false}
					orderEventCh <- OrderEvent{elev.Floor, OrderCallCommand, false}
					timerResetCh <- time.Second * 3
					doorOpenCh <- true
					elev.Behaviour = ElevatorBehaviourDoorOpen

				} else {
					elev.Direction = GetDirection(elev.Orders, elev)
					if elev.Direction == DirnStop {
						elev.Behaviour = ElevatorBehaviourIdle
					} else {
						elev.Behaviour = ElevatorBehaviourMoving
					}
					motorStateCh <- elev.Direction
				}
			}
		case globalOrders := <-globalOrdersCh:
			for f := 0; f < NumFloors; f++ {
				for b := 0; b < NumTypes; b++ {
					lightEventCh <- LightEvent{LightType(b), f, globalelev.Orders[f][b]}
				}
			}
		case <-timerTimeoutCh:
			switch elev.Behaviour {
			case ElevatorBehaviourDoorOpen:
				elev.Direction = GetDirection(elev.Orders, elev)
				if elev.Direction == DirnStop {
					elev.Behaviour = ElevatorBehaviourIdle
				} else {
					elev.Behaviour = ElevatorBehaviourMoving
				}
				motorStateCh <- elev.Direction
				doorOpenCh <- false
			}

		case <-time.After(10 * time.Millisecond):
		}
	}
}
