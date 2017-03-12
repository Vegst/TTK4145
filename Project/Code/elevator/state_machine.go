package elevator

import (
	"time"
	"./timer"
	. "../def"
	//"fmt"
)


func StateMachine(buttonEventCh chan ButtonEvent, lightEventCh chan LightEvent, stopCh chan bool, motorStateCh chan MotorDirection, floorCh chan int, doorOpenCh chan bool, floorIndicatorCh chan int, orderEventCh chan OrderEvent, stateCh chan ElevatorState, localOrdersCh chan Orders, globalOrdersCh chan Orders) {


	timerResetCh := make(chan time.Duration)
	timerTimeoutCh := make(chan bool)
	go timer.Timer(timerResetCh, timerTimeoutCh)
	
	var elev Elevator
	var globalElev Elevator

	// Initial state change
	motorStateCh <- DirnUp
	elev.State.Direction = DirnUp
	elev.State.Behaviour = ElevatorBehaviourMoving

	for {
		select {
		// Event: Button pressed
		case buttonEvent := <-buttonEventCh:
			if buttonEvent.State {
				orderEventCh <- OrderEvent{buttonEvent.Floor, OrderType(buttonEvent.Button), true}
			}
		// Event: Stop command
		case <- stopCh:
		// Event: Floor reached
		case elev.State.Floor = <-floorCh:
			if elev.State.Floor >= 0 && elev.State.Floor < NumFloors {
				floorIndicatorCh <- elev.State.Floor
				
				switch(elev.State.Behaviour) {
				case ElevatorBehaviourMoving:
					if(ShouldStop(elev)){
						if OrderAtFloor(elev) {
							// Clear elev.Orders at current floor
							if elev.State.Direction == DirnUp {
								orderEventCh <- OrderEvent{elev.State.Floor, OrderCallUp, false}
							} else if elev.State.Direction == DirnDown {
								orderEventCh <- OrderEvent{elev.State.Floor, OrderCallDown, false}
							}
							orderEventCh <- OrderEvent{elev.State.Floor, OrderCallCommand, false}

							doorOpenCh <- true
							timerResetCh <- time.Second * 3
							motorStateCh <- DirnStop

							elev.State.Behaviour = ElevatorBehaviourDoorOpen
						} else {
							motorStateCh <- DirnStop
							elev.State.Behaviour = ElevatorBehaviourIdle
							elev.State.Direction = DirnStop
						}
					}
				}
			}
			
		case elev.Orders = <- localOrdersCh:
			switch(elev.State.Behaviour) {
			case ElevatorBehaviourDoorOpen:
				if OrderAtFloor(elev) {
					orderEventCh <- OrderEvent{elev.State.Floor, OrderCallUp, false}
					orderEventCh <- OrderEvent{elev.State.Floor, OrderCallDown, false}
					orderEventCh <- OrderEvent{elev.State.Floor, OrderCallCommand, false}
					timerResetCh <- time.Second * 3
				}
			case ElevatorBehaviourIdle:
				if OrderAtFloor(elev) {
					orderEventCh <- OrderEvent{elev.State.Floor, OrderCallUp, false}
					orderEventCh <- OrderEvent{elev.State.Floor, OrderCallDown, false}
					orderEventCh <- OrderEvent{elev.State.Floor, OrderCallCommand, false}
					timerResetCh <- time.Second * 3
					doorOpenCh <- true
					elev.State.Behaviour = ElevatorBehaviourDoorOpen

				} else {
					elev.State.Direction = GetDirection(elev)
					if elev.State.Direction == DirnStop {
						elev.State.Behaviour = ElevatorBehaviourIdle
					} else {
						elev.State.Behaviour = ElevatorBehaviourMoving
					}
					motorStateCh <- elev.State.Direction
				}
			}
		case globalElev.Orders = <-globalOrdersCh:
			for f := 0; f < NumFloors; f++ {
				for b := 0; b < NumTypes; b++ {
					lightEventCh <- LightEvent{LightType(b), f, globalElev.Orders[f][b]}
				}
			}
		case <-timerTimeoutCh:
			switch(elev.State.Behaviour) {
			case ElevatorBehaviourDoorOpen:
				elev.State.Direction = GetDirection(elev)
				if elev.State.Direction == DirnStop {
					elev.State.Behaviour = ElevatorBehaviourIdle
				} else {
					elev.State.Behaviour = ElevatorBehaviourMoving
				}
				motorStateCh <- elev.State.Direction
				doorOpenCh <- false
			}

		case <-time.After(10 * time.Millisecond):
		}
	}
}
