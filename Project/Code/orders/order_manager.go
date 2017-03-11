package orders

import (
	"time"
	. "../def"
)


func OrderManager(id string, orderEventCh <-chan OrderEvent, stateCh <-chan Elevator, localOrdersCh chan Orders, globalOrdersCh chan Orders, elevatorsCh chan Elevators) {
	var elevators Elevators
	var orders Orders
	elevators = make(Elevators)
	elevators[id] = Elevator{0,DirnStop,ElevatorBehaviourIdle,orders}
	elevatorsCh <- elevators
	for {
		select {
		case orderEvent := <- orderEventCh:
			elevator := elevators[id]
			elevator.Orders[orderEvent.Floor][orderEvent.Type] = orderEvent.Flag
			elevators[id] = elevator
			localOrdersCh <- elevator.Orders
			globalOrdersCh <- elevator.Orders
			elevatorsCh <- elevators
		case <- stateCh:
		case <-time.After(50*time.Millisecond):
		}
	}
}