package elevator

import (
	. "../def"
)

type Actuator struct {
	DriverEvents DriverElevatorEvents
	OrdersEvents ElevatorOrdersEvents
	DoorTimerResetCh chan bool
	ErrorTimerResetCh chan bool
}
