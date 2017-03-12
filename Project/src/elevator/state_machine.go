/*
package elevator

import (
	. "../def"
	"./timer"
	"time"
	"fmt"
)


func (state *ElevatorState) OnInit(driverEvents DriverElevatorEvents, ordersEvents ElevatorOrdersEvents) {
	state = ElevatorState{Floor: -1, Direction: DirnUp, Behaviour: ElevatorBehaviourMoving}
	driverEvents.MotorDirection <- DirnUp
	ordersEvents.State <- elev.State
}

func (state *ElevatorState) OnButtonPressed(driverEvents DriverElevatorEvents, ordersEvents ElevatorOrdersEvents, button Button) {
	ordersEvents.Order <- OrderEvent{button.Floor, OrderType(buttonEvent.Button), true}
}

func (state *ElevatorState) OnStopBegin(driverEvents DriverElevatorEvents, ordersEvents ElevatorOrdersEvents) {

}

func (state *ElevatorState) OnStopEnd(driverEvents DriverElevatorEvents, ordersEvents ElevatorOrdersEvents) {

}

func (state *ElevatorState) OnFloorReached(driverEvents DriverElevatorEvents, ordersEvents ElevatorOrdersEvents, floor int) {
	state.Floor = <-driverEvents.Floor:
	ordersEvents.State <- elev.State
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
				ordersEvents.State <- elev.State
				ordersEvents.Order <- OrderEvent{elev.State.Floor, OrderCallCommand, false}

				driverEvents.DoorOpen <- true
				timerResetCh <- time.Second * 3
				driverEvents.MotorDirection <- DirnStop

				elev.State.Behaviour = ElevatorBehaviourDoorOpen
				ordersEvents.State <- elev.State
			} else {
				elev.State.Behaviour = ElevatorBehaviourIdle
				elev.State.Direction = DirnStop
				ordersEvents.State <- elev.State
				driverEvents.MotorDirection <- DirnStop
			}
		}
	
	}
}

func (state *ElevatorState) OnLocalOrdersUpdated(driverEvents DriverElevatorEvents, ordersEvents ElevatorOrdersEvents, orders Orders) {
	for f := 0; f < NumFloors; f++ {
		driverEvents.Light <- LightEvent{LightType(OrderCallCommand), f, orders[f][OrderCallCommand]}
	}
	elev := Elevator{state, orders}
	switch state.Behaviour {
	case ElevatorBehaviourDoorOpen:
		if OrderAtFloor(elev) {
			ordersEvents.Order <- OrderEvent{state.Floor, OrderCallUp, false}
			ordersEvents.Order <- OrderEvent{state.Floor, OrderCallDown, false}
			ordersEvents.Order <- OrderEvent{state.Floor, OrderCallCommand, false}
			timerResetCh <- time.Second * 3
		}
	case ElevatorBehaviourIdle:
		if OrderAtFloor(elev) {
			ordersEvents.Order <- OrderEvent{state.Floor, OrderCallUp, false}
			ordersEvents.Order <- OrderEvent{state.Floor, OrderCallDown, false}
			ordersEvents.Order <- OrderEvent{state.Floor, OrderCallCommand, false}
			timerResetCh <- time.Second * 3
			driverEvents.DoorOpen <- true
			state.Behaviour = ElevatorBehaviourDoorOpen
			ordersEvents.State <- state

		} else {
			state.Direction = GetDirection(elev)
			if state.Direction == DirnStop {
				state.Behaviour = ElevatorBehaviourIdle
			} else {
				state.Behaviour = ElevatorBehaviourMoving
			}
			ordersEvents.State <- state
			driverEvents.MotorDirection <- state.Direction
		}
	}
}


func (state *ElevatorState) OnGlobalOrdersUpdated(driverEvents DriverElevatorEvents, ordersEvents ElevatorOrdersEvents, orders Orders) {
	for f := 0; f < NumFloors; f++ {
		driverEvents.Light <- LightEvent{LightType(OrderCallDown), f, globalOrders[f][OrderCallDown]}
		driverEvents.Light <- LightEvent{LightType(OrderCallUp), f, globalOrders[f][OrderCallUp]}
	}
}
*/
package elevator

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
	elev := Elevator{
		State: ElevatorState{Floor: -1, Direction: DirnUp, Behaviour: ElevatorBehaviourMoving},
		Orders: Orders{{}},
	}
	driverEvents.MotorDirection <- DirnUp
	ordersEvents.State <- elev.State

	for {
		select {
		// Event: Button pressed
		case buttonEvent := <-driverEvents.Button:
			if buttonEvent.State {
				oe := OrderEvent{buttonEvent.Button.Floor, OrderType(buttonEvent.Button.Type), true}
				ordersEvents.Order <- oe
			}
		// Event: Stop command
		case <-driverEvents.Stop:
		// Event: Floor reached
		case elev.State.Floor = <-driverEvents.Floor:
			fmt.Println("Floor reached:", elev.State.Floor)
			ordersEvents.State <- elev.State
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
						ordersEvents.State <- elev.State
						ordersEvents.Order <- OrderEvent{elev.State.Floor, OrderCallCommand, false}

						driverEvents.DoorOpen <- true
						timerResetCh <- time.Second * 3
						driverEvents.MotorDirection <- DirnStop

						elev.State.Behaviour = ElevatorBehaviourDoorOpen
						ordersEvents.State <- elev.State
					} else {
						elev.State.Behaviour = ElevatorBehaviourIdle
						elev.State.Direction = DirnStop
						ordersEvents.State <- elev.State
						driverEvents.MotorDirection <- DirnStop
					}
				}
			
			}

		case elev.Orders = <-ordersEvents.LocalOrders:
			for f := 0; f < NumFloors; f++ {
				driverEvents.Light <- LightEvent{LightType(OrderCallCommand), f, elev.Orders[f][OrderCallCommand]}
			}
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
					ordersEvents.State <- elev.State

				} else {
					elev.State.Direction = GetDirection(elev)
					if elev.State.Direction == DirnStop {
						elev.State.Behaviour = ElevatorBehaviourIdle
					} else {
						elev.State.Behaviour = ElevatorBehaviourMoving
					}
					ordersEvents.State <- elev.State
					driverEvents.MotorDirection <- elev.State.Direction
				}
			}
		case globalOrders := <-ordersEvents.GlobalOrders:
			for f := 0; f < NumFloors; f++ {
				driverEvents.Light <- LightEvent{LightType(OrderCallDown), f, globalOrders[f][OrderCallDown]}
				driverEvents.Light <- LightEvent{LightType(OrderCallUp), f, globalOrders[f][OrderCallUp]}
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
				ordersEvents.State <- elev.State
				driverEvents.MotorDirection <- elev.State.Direction
				driverEvents.DoorOpen <- false
			}

		case <-time.After(10 * time.Millisecond):
		}
	}
}
