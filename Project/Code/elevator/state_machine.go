package elevator

import (
	. "../def"
	"./timer"
	"time"
	//"fmt"
)

func StateMachine(driverEvents DriverElevatorEvents, ordersEvents ElevatorOrdersEvents) {

	timerResetCh := make(chan time.Duration)
	timerTimeoutCh := make(chan bool)
	go timer.Timer(timerResetCh, timerTimeoutCh)

	var elev Elevator

	// Initial state change
	driverEvents.MotorDirection <- DirnUp
	elev.State.Direction = DirnUp
	elev.State.Behaviour = ElevatorBehaviourMoving

	for {
		select {
		// Event: Button pressed
		case buttonEvent := <-driverEvents.Button:
			if buttonEvent.State {
				oe := OrderEvent{buttonEvent.Floor, OrderType(buttonEvent.Button), true}
				ordersEvents.Order <- oe
				CalculateCost(oe, elev)
			}
		// Event: Stop command
		case <-driverEvents.Stop:
		// Event: Floor reached
		case elev.State.Floor = <-driverEvents.Floor:
			if elev.State.Floor >= 0 && elev.State.Floor < NumFloors {
				driverEvents.FloorIndicator <- elev.State.Floor

				switch elev.State.Behaviour {
				case ElevatorBehaviourMoving:
					if ShouldStop(elev) {
						if OrderAtFloor(elev) {
							// Clear elev.Orders at current floor
							if elev.State.Direction == DirnUp {
								ordersEvents.Order <- OrderEvent{elev.State.Floor, OrderCallUp, false}
							} else if elev.State.Direction == DirnDown {
								ordersEvents.Order <- OrderEvent{elev.State.Floor, OrderCallDown, false}
							}
							ordersEvents.Order <- OrderEvent{elev.State.Floor, OrderCallCommand, false}

							driverEvents.DoorOpen <- true
							timerResetCh <- time.Second * 3
							driverEvents.MotorDirection <- DirnStop

							elev.State.Behaviour = ElevatorBehaviourDoorOpen
						} else {
							driverEvents.MotorDirection <- DirnStop
							elev.State.Behaviour = ElevatorBehaviourIdle
							elev.State.Direction = DirnStop
						}
					}
				}
			}

		case elev.Orders = <-ordersEvents.LocalOrders:
			switch elev.State.Behaviour {
			case ElevatorBehaviourDoorOpen:
				if OrderAtFloor(elev) {
					ordersEvents.Order <- OrderEvent{elev.State.Floor, OrderCallUp, false}
					ordersEvents.Order <- OrderEvent{elev.State.Floor, OrderCallDown, false}
					ordersEvents.Order <- OrderEvent{elev.State.Floor, OrderCallCommand, false}
					timerResetCh <- time.Second * 3
				}
			case ElevatorBehaviourIdle:
				if OrderAtFloor(elev) {
					ordersEvents.Order <- OrderEvent{elev.State.Floor, OrderCallUp, false}
					ordersEvents.Order <- OrderEvent{elev.State.Floor, OrderCallDown, false}
					ordersEvents.Order <- OrderEvent{elev.State.Floor, OrderCallCommand, false}
					timerResetCh <- time.Second * 3
					driverEvents.DoorOpen <- true
					elev.State.Behaviour = ElevatorBehaviourDoorOpen

				} else {
					elev.State.Direction = GetDirection(elev)
					if elev.State.Direction == DirnStop {
						elev.State.Behaviour = ElevatorBehaviourIdle
					} else {
						elev.State.Behaviour = ElevatorBehaviourMoving
					}
					driverEvents.MotorDirection <- elev.State.Direction
				}
			}
		case globalOrders := <-ordersEvents.GlobalOrders:
			for f := 0; f < NumFloors; f++ {
				for b := 0; b < NumTypes; b++ {
					driverEvents.Light <- LightEvent{LightType(b), f, globalOrders[f][b]}
				}
			}
		case <-timerTimeoutCh:
			switch elev.State.Behaviour {
			case ElevatorBehaviourDoorOpen:
				elev.State.Direction = GetDirection(elev)
				if elev.State.Direction == DirnStop {
					elev.State.Behaviour = ElevatorBehaviourIdle
				} else {
					elev.State.Behaviour = ElevatorBehaviourMoving
				}
				driverEvents.MotorDirection <- elev.State.Direction
				driverEvents.DoorOpen <- false
			}

		case <-time.After(10 * time.Millisecond):
		}
	}
}
