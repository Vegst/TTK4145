package elevator

import (
	. "../def"
)

type StateMachine struct {
	DriverEvents DriverElevatorEvents
	OrdersEvents ElevatorOrdersEvents
	DoorTimerResetCh chan bool
	ErrorTimerResetCh chan bool
	Elevator Elevator
}

func NewStateMachine(driverEvents DriverElevatorEvents, ordersEvents ElevatorOrdersEvents, doorTimerResetCh chan bool, errorTimerResetCh chan bool) *StateMachine {
	sm := new(StateMachine)
	sm.DriverEvents = driverEvents
	sm.OrdersEvents = ordersEvents
	sm.DoorTimerResetCh = doorTimerResetCh
	sm.ErrorTimerResetCh = errorTimerResetCh
	sm.Elevator.State = ElevatorState{Active: false, Floor: -1, Direction: DirnStop, Behaviour: ElevatorBehaviourIdle}
	sm.Elevator.Orders = Orders{{}}
	return sm
}

func (sm *StateMachine) OnInit() {
	sm.Elevator.State.Direction = DirnUp
	sm.Elevator.State.Behaviour = ElevatorBehaviourMoving
	sm.DriverEvents.MotorDirection <- DirnUp
	sm.OrdersEvents.State <- sm.Elevator.State
	sm.ErrorTimerResetCh <- true
}

func (sm *StateMachine) OnButtonPressed(button Button) {
	sm.OrdersEvents.Order <- Order{button.Floor, OrderType(button.Type), true}
}

func (sm *StateMachine) OnButtonReleased(button Button) {}

func (sm *StateMachine) OnStopPressed() {
	if sm.Elevator.State.Active {
		sm.Elevator.State.Active = false
		sm.DriverEvents.MotorDirection <- DirnStop
		sm.DriverEvents.Light <- LightEvent{LightTypeStop, 0, true}
		sm.OrdersEvents.State <- sm.Elevator.State
	}
}

func (sm *StateMachine) OnStopReleased() {}

func (sm *StateMachine) OnFloorReached(floor int) {
	sm.Elevator.State.Floor = floor
	sm.OrdersEvents.State <- sm.Elevator.State
	sm.DriverEvents.FloorIndicator <- sm.Elevator.State.Floor
	sm.ErrorTimerResetCh <- true

	switch sm.Elevator.State.Behaviour {
	case ElevatorBehaviourMoving:
		if !sm.Elevator.State.Active {
			sm.Elevator.State.Active = true
			sm.OrdersEvents.State <- sm.Elevator.State
		}
		if ShouldStop(sm.Elevator) {
			if IsOrderAtFloor(sm.Elevator) {
				if sm.Elevator.State.Direction == DirnUp {
					sm.OrdersEvents.Order <- Order{sm.Elevator.State.Floor, OrderCallUp, false}
				} else if sm.Elevator.State.Direction == DirnDown {
					sm.OrdersEvents.Order <- Order{sm.Elevator.State.Floor, OrderCallDown, false}
				}
				sm.OrdersEvents.Order <- Order{sm.Elevator.State.Floor, OrderCallCommand, false}

				sm.DriverEvents.DoorOpen <- true
				sm.DoorTimerResetCh <- true

				sm.Elevator.State.Behaviour = ElevatorBehaviourDoorOpen
			} else {
				sm.Elevator.State.Behaviour = ElevatorBehaviourIdle
				sm.Elevator.State.Direction = DirnStop
			}
			sm.DriverEvents.MotorDirection <- DirnStop
			sm.OrdersEvents.State <- sm.Elevator.State
		}
	}
}

func (sm *StateMachine) OnLocalOrdersUpdated(localOrders Orders) {
	sm.Elevator.Orders = localOrders
	for f := 0; f < NumFloors; f++ {
		sm.DriverEvents.Light <- LightEvent{LightType(OrderCallCommand), f, sm.Elevator.Orders[f][OrderCallCommand]}
	}
	if sm.Elevator.State.Active {
		switch sm.Elevator.State.Behaviour {
		case ElevatorBehaviourDoorOpen:
			if IsOrderAtFloor(sm.Elevator) {
				sm.OrdersEvents.Order <- Order{sm.Elevator.State.Floor, OrderCallUp, false}
				sm.OrdersEvents.Order <- Order{sm.Elevator.State.Floor, OrderCallDown, false}
				sm.OrdersEvents.Order <- Order{sm.Elevator.State.Floor, OrderCallCommand, false}
				sm.DoorTimerResetCh <- true
			}
		case ElevatorBehaviourIdle:
			if IsOrderAtFloor(sm.Elevator) {
				sm.OrdersEvents.Order <- Order{sm.Elevator.State.Floor, OrderCallUp, false}
				sm.OrdersEvents.Order <- Order{sm.Elevator.State.Floor, OrderCallDown, false}
				sm.OrdersEvents.Order <- Order{sm.Elevator.State.Floor, OrderCallCommand, false}
				sm.DoorTimerResetCh <- true
				sm.DriverEvents.DoorOpen <- true
				sm.Elevator.State.Behaviour = ElevatorBehaviourDoorOpen
				sm.OrdersEvents.State <- sm.Elevator.State

			} else {
				sm.Elevator.State.Direction = GetDirection(sm.Elevator)
				if sm.Elevator.State.Direction == DirnStop {
					sm.Elevator.State.Behaviour = ElevatorBehaviourIdle
				} else {
					sm.Elevator.State.Behaviour = ElevatorBehaviourMoving
				}
				sm.OrdersEvents.State <- sm.Elevator.State
				sm.DriverEvents.MotorDirection <- sm.Elevator.State.Direction
				sm.ErrorTimerResetCh <- true
			}
		}
	}
}

func (sm *StateMachine) OnGlobalOrdersUpdated(globalOrders Orders) {
	for f := 0; f < NumFloors; f++ {
		sm.DriverEvents.Light <- LightEvent{LightType(OrderCallDown), f, globalOrders[f][OrderCallDown]}
		sm.DriverEvents.Light <- LightEvent{LightType(OrderCallUp), f, globalOrders[f][OrderCallUp]}
	}
}

func (sm *StateMachine) OnDoorTimerTimeout() {
	switch sm.Elevator.State.Behaviour {
	case ElevatorBehaviourDoorOpen:
		if sm.Elevator.State.Active {
			sm.Elevator.State.Direction = GetDirection(sm.Elevator)
			if sm.Elevator.State.Direction == DirnStop {
				sm.Elevator.State.Behaviour = ElevatorBehaviourIdle
			} else {
				sm.Elevator.State.Behaviour = ElevatorBehaviourMoving
			}
			sm.OrdersEvents.State <- sm.Elevator.State
			sm.DriverEvents.MotorDirection <- sm.Elevator.State.Direction
			sm.DriverEvents.DoorOpen <- false
			sm.ErrorTimerResetCh <- true
		}
	}
}

func (sm *StateMachine) OnErrorTimerTimeout() {
	switch sm.Elevator.State.Behaviour {
	case ElevatorBehaviourMoving:
		if sm.Elevator.State.Active {
			sm.Elevator.State.Active = false
			sm.DriverEvents.MotorDirection <- DirnStop
			sm.DriverEvents.Light <- LightEvent{LightTypeStop, 0, true}
			sm.OrdersEvents.State <- sm.Elevator.State
		}
	}
}