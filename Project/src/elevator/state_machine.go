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
	ordersEvents.State <- state
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
	ordersEvents.State <- state
	driverEvents.FloorIndicator <- state.Floor

	switch state.Behaviour {
	case ElevatorBehaviourMoving:
		if ShouldStop(elev) {
			if OrderAtFloor(elev) {
				// Clear orders at current floor
				if state.Direction == DirnUp {
					ordersEvents.Order <- OrderEvent{state.Floor, OrderCallUp, false}
				} else if state.Direction == DirnDown {
					ordersEvents.Order <- OrderEvent{state.Floor, OrderCallDown, false}
				}
				ordersEvents.State <- state
				ordersEvents.Order <- OrderEvent{state.Floor, OrderCallCommand, false}

				driverEvents.DoorOpen <- true
				timerResetCh <- time.Second * 3
				driverEvents.MotorDirection <- DirnStop

				state.Behaviour = ElevatorBehaviourDoorOpen
				ordersEvents.State <- state
			} else {
				state.Behaviour = ElevatorBehaviourIdle
				state.Direction = DirnStop
				ordersEvents.State <- state
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
)

func StateMachine(driverEvents DriverElevatorEvents, ordersEvents ElevatorOrdersEvents) {

	timerResetCh := make(chan time.Duration)
	timerTimeoutCh := make(chan bool)
	go timer.Timer(timerResetCh, timerTimeoutCh)

	// Initial state
	state := ElevatorState{Active: false, Floor: -1, Direction: DirnUp, Behaviour: ElevatorBehaviourMoving}
	orders := Orders{{}}
	
	driverEvents.MotorDirection <- DirnUp
	ordersEvents.State <- state

	for {
		select {
		// Event: Button pressed
		case buttonEvent := <-driverEvents.Button:
			if buttonEvent.State {
				ordersEvents.Order <- Order{buttonEvent.Button.Floor, OrderType(buttonEvent.Button.Type), true}
			}

		// Event: Stop command
		case <-driverEvents.Stop:
			driverEvents.MotorDirection <- DirnStop
			state.Active = false
			ordersEvents.State <- state


		// Event: Floor reached
		case state.Floor = <-driverEvents.Floor:
			ordersEvents.State <- state
			driverEvents.FloorIndicator <- state.Floor

			switch state.Behaviour {
			case ElevatorBehaviourMoving:
				if !state.Active {
					state.Active = true
					ordersEvents.State <- state
				}
				if ShouldStop(orders, state.Floor, state.Direction) {
					if OrderAtFloor(orders, state.Floor) {
						if state.Direction == DirnUp {
							ordersEvents.Order <- Order{state.Floor, OrderCallUp, false}
						} else if state.Direction == DirnDown {
							ordersEvents.Order <- Order{state.Floor, OrderCallDown, false}
						}
						ordersEvents.Order <- Order{state.Floor, OrderCallCommand, false}

						driverEvents.DoorOpen <- true
						timerResetCh <- time.Second * 3
						driverEvents.MotorDirection <- DirnStop

						state.Behaviour = ElevatorBehaviourDoorOpen
					} else {
						state.Behaviour = ElevatorBehaviourIdle
						state.Direction = DirnStop
					}
					driverEvents.MotorDirection <- DirnStop
					ordersEvents.State <- state
				}
			
			}

		case orders = <-ordersEvents.LocalOrders:
			for f := 0; f < NumFloors; f++ {
				driverEvents.Light <- LightEvent{LightType(OrderCallCommand), f, orders[f][OrderCallCommand]}
			}
			switch state.Behaviour {
			case ElevatorBehaviourDoorOpen:
				if OrderAtFloor(orders, state.Floor) {
					ordersEvents.Order <- Order{state.Floor, OrderCallUp, false}
					ordersEvents.Order <- Order{state.Floor, OrderCallDown, false}
					ordersEvents.Order <- Order{state.Floor, OrderCallCommand, false}
					timerResetCh <- time.Second * 3
				}
			case ElevatorBehaviourIdle:
				if OrderAtFloor(orders, state.Floor) {
					ordersEvents.Order <- Order{state.Floor, OrderCallUp, false}
					ordersEvents.Order <- Order{state.Floor, OrderCallDown, false}
					ordersEvents.Order <- Order{state.Floor, OrderCallCommand, false}
					timerResetCh <- time.Second * 3
					driverEvents.DoorOpen <- true
					state.Behaviour = ElevatorBehaviourDoorOpen
					ordersEvents.State <- state

				} else {
					state.Direction = GetDirection(orders, state.Floor, state.Direction)
					if state.Direction == DirnStop {
						state.Behaviour = ElevatorBehaviourIdle
					} else {
						state.Behaviour = ElevatorBehaviourMoving
					}
					ordersEvents.State <- state
					driverEvents.MotorDirection <- state.Direction
				}
			}

		case globalOrders := <-ordersEvents.GlobalOrders:
			for f := 0; f < NumFloors; f++ {
				driverEvents.Light <- LightEvent{LightType(OrderCallDown), f, globalOrders[f][OrderCallDown]}
				driverEvents.Light <- LightEvent{LightType(OrderCallUp), f, globalOrders[f][OrderCallUp]}
			}

		case <-timerTimeoutCh:
			switch state.Behaviour {
			case ElevatorBehaviourDoorOpen:
				state.Direction = GetDirection(orders, state.Floor, state.Direction)
				if state.Direction == DirnStop {
					state.Behaviour = ElevatorBehaviourIdle
				} else {
					state.Behaviour = ElevatorBehaviourMoving
				}
				ordersEvents.State <- state
				driverEvents.MotorDirection <- state.Direction
				driverEvents.DoorOpen <- false
			}

		case <-time.After(10 * time.Millisecond):
		}
	}
}
