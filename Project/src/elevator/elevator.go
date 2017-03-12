package elevator
/*
import (
	. "../def"
	"./timer"
	"time"
	"fmt"
)

func StateMachine(driverEvents DriverElevatorEvents, ordersEvents ElevatorOrdersEvents) {

	timerResetCh := make(chan time.Duration)
	timerTimeoutCh := make(chan bool)
	go timer.Timer(timerResetCh, timerTimeoutCh)

	// Initial state
	OnInit()
	

	for {
		select {
		case buttonEvent := <-driverEvents.Button:
			if buttonEvent.State {
				OnButtonPressed(driverEvents, ordersEvents, buttonEvent.Button)
			}
		case stop := <-driverEvents.Stop:
			if stop {
				OnStopBegin(driverEvents, ordersEvents)
			} else {
				OnStopEnd(driverEvents, ordersEvents)
			}
		case elev.State.Floor = <-driverEvents.Floor:
			OnFloorReached(driverEvents, ordersEvents, elev.State.Floor)

		case elev.Orders = <-ordersEvents.LocalOrders:
			OnLocalOrdersUpdated(driverEvents, ordersEvents, elev.Orders)
			
		case globalOrders := <-ordersEvents.GlobalOrders:
			OnGlobalOrdersUpdated(driverEvents, ordersEvents, globalOrders)
		case <-timerTimeoutCh:
			OnTim
			switch elev.State.Behaviour {
			case ElevatorBehaviourDoorOpen:
				elev.State.Direction = GetDirection(elev)
				if elev.State.Direction == DirnStop {
					elev.State.Behaviour = ElevatorBehaviourIdle
				} else {
					elev.State.Behaviour = ElevatorBehaviourMoving
				}
				ordersEvents.State <- elev.State
				driverEvents.MotorDirection <- elev.State.Direction
				driverEvents.DoorOpen <- false
			}

		case <-time.After(10 * time.Millisecond):
		}
	}
}
*/