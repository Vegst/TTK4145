package elevator

import (
	. "../def"
)

type StateMachine struct {
	DriverEvents DriverElevatorEvents
	OrdersEvents ElevatorOrdersEvents
	DoorTimerResetCh chan bool
	ErrorTimerResetCh chan bool



	State ElevatorState
	Orders Orders
}

func NewStateMachine(driverEvents DriverElevatorEvents, ordersEvents ElevatorOrdersEvents, doorTimerResetCh chan bool, errorTimerResetCh chan bool) *StateMachine {
	sm := new(StateMachine)
	sm.DriverEvents = driverEvents
	sm.OrdersEvents = ordersEvents
	sm.DoorTimerResetCh = doorTimerResetCh
	sm.ErrorTimerResetCh = errorTimerResetCh
	sm.State = ElevatorState{Active: false, Floor: -1, Direction: DirnStop, Behaviour: ElevatorBehaviourIdle}
	sm.Orders = Orders{{}}
	return sm
}

func (sm *StateMachine) OnInit() {
	sm.State.Direction = DirnUp
	sm.State.Behaviour = ElevatorBehaviourMoving
	sm.DriverEvents.MotorDirection <- DirnUp
	sm.OrdersEvents.State <- sm.State
}

func (sm *StateMachine) OnButtonPressed(button Button) {
	sm.OrdersEvents.Order <- Order{button.Floor, OrderType(button.Type), true}
}

func (sm *StateMachine) OnButtonReleased(button Button) {}

func (sm *StateMachine) OnStopPressed() {
	if sm.State.Active {
		sm.State.Active = false
		sm.DriverEvents.MotorDirection <- DirnStop
		sm.DriverEvents.Light <- LightEvent{LightTypeStop, 0, true}
		sm.OrdersEvents.State <- sm.State
	}
}

func (sm *StateMachine) OnStopReleased() {}

func (sm *StateMachine) OnFloorReached(floor int) {
	sm.State.Floor = floor
	sm.OrdersEvents.State <- sm.State
	sm.DriverEvents.FloorIndicator <- sm.State.Floor

	switch sm.State.Behaviour {
	case ElevatorBehaviourMoving:
		if !sm.State.Active {
			sm.State.Active = true
			sm.OrdersEvents.State <- sm.State
		}
		if ShouldStop(sm.Orders, sm.State.Floor, sm.State.Direction) {
			if IsOrderAtFloor(sm.Orders, sm.State.Floor) {
				if sm.State.Direction == DirnUp {
					sm.OrdersEvents.Order <- Order{sm.State.Floor, OrderCallUp, false}
				} else if sm.State.Direction == DirnDown {
					sm.OrdersEvents.Order <- Order{sm.State.Floor, OrderCallDown, false}
				}
				sm.OrdersEvents.Order <- Order{sm.State.Floor, OrderCallCommand, false}

				sm.DriverEvents.DoorOpen <- true
				sm.DoorTimerResetCh <- true
				sm.DriverEvents.MotorDirection <- DirnStop

				sm.State.Behaviour = ElevatorBehaviourDoorOpen
			} else {
				sm.State.Behaviour = ElevatorBehaviourIdle
				sm.State.Direction = DirnStop
			}
			sm.DriverEvents.MotorDirection <- DirnStop
			sm.OrdersEvents.State <- sm.State
		}
	
	}
}

func (sm *StateMachine) OnLocalOrdersUpdated(localOrders Orders) {
	sm.Orders = localOrders
	for f := 0; f < NumFloors; f++ {
		sm.DriverEvents.Light <- LightEvent{LightType(OrderCallCommand), f, sm.Orders[f][OrderCallCommand]}
	}
	if sm.State.Active {
		switch sm.State.Behaviour {
		case ElevatorBehaviourDoorOpen:
			if IsOrderAtFloor(sm.Orders, sm.State.Floor) {
				sm.OrdersEvents.Order <- Order{sm.State.Floor, OrderCallUp, false}
				sm.OrdersEvents.Order <- Order{sm.State.Floor, OrderCallDown, false}
				sm.OrdersEvents.Order <- Order{sm.State.Floor, OrderCallCommand, false}
				sm.DoorTimerResetCh <- true
			}
		case ElevatorBehaviourIdle:
			if IsOrderAtFloor(sm.Orders, sm.State.Floor) {
				sm.OrdersEvents.Order <- Order{sm.State.Floor, OrderCallUp, false}
				sm.OrdersEvents.Order <- Order{sm.State.Floor, OrderCallDown, false}
				sm.OrdersEvents.Order <- Order{sm.State.Floor, OrderCallCommand, false}
				sm.DoorTimerResetCh <- true
				sm.DriverEvents.DoorOpen <- true
				sm.State.Behaviour = ElevatorBehaviourDoorOpen
				sm.OrdersEvents.State <- sm.State

			} else {
				sm.State.Direction = GetDirection(sm.Orders, sm.State.Floor, sm.State.Direction)
				if sm.State.Direction == DirnStop {
					sm.State.Behaviour = ElevatorBehaviourIdle
				} else {
					sm.State.Behaviour = ElevatorBehaviourMoving
				}
				sm.OrdersEvents.State <- sm.State
				sm.DriverEvents.MotorDirection <- sm.State.Direction
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
	switch sm.State.Behaviour {
	case ElevatorBehaviourDoorOpen:
		if sm.State.Active {
			sm.State.Direction = GetDirection(sm.Orders, sm.State.Floor, sm.State.Direction)
			if sm.State.Direction == DirnStop {
				sm.State.Behaviour = ElevatorBehaviourIdle
			} else {
				sm.State.Behaviour = ElevatorBehaviourMoving
			}
			sm.OrdersEvents.State <- sm.State
			sm.DriverEvents.MotorDirection <- sm.State.Direction
			sm.DriverEvents.DoorOpen <- false
		}
	}
}

func (sm *StateMachine) OnErrorTimerTimeout() {
	if sm.State.Active {
		sm.State.Active = false
		sm.DriverEvents.MotorDirection <- DirnStop
		sm.DriverEvents.Light <- LightEvent{LightTypeStop, 0, true}
		sm.OrdersEvents.State <- sm.State
	}
}