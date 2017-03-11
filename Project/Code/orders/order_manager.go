package orders

import (
	. "../def"
	"../elevator"
	"fmt"
	"time"
)

func OrderManager(id string, orderEventCh <-chan OrderEvent, assignedOrderCh <-chan OrderEvent, stateCh <-chan ElevatorState, updateElevatorCh <-chan Elevator, localOrdersCh chan Orders, globalOrdersCh chan Orders, elevatorCh chan Elevator, elevatorsCh chan Elevators, netOrderCh chan NetOrder) {
	var elevators Elevators
	var orders Orders

	elevators = make(Elevators)
	localElevator := Elevator{ElevatorState{0, DirnStop, ElevatorBehaviourIdle}, orders, id}
	elevatorCh <- localElevator
	for {
		select {
		case elev := <-updateElevatorCh:
			elevators[id] = elev
			elevatorsCh <- elevators
		case orderEvent := <-orderEventCh:
			elev := localElevator
			//Assigning an order to an elevator
			AssignmentID := elevator.OrderAssigner(orderEvent, elevators)
			netOrderCh <- NetOrder{orderEvent, AssignmentID}

			elev.Orders[orderEvent.Floor][orderEvent.Type] = orderEvent.Flag
			elevators[id] = elev
			elevatorCh <- elev

		case <-stateCh:

		case assignedOrder := <-assignedOrderCh:
			fmt.Println("New order for elevator: %s", id)
			elev := localElevator
			elev.Orders[assignedOrder.Floor][assignedOrder.Type] = assignedOrder.Flag
			localElevator = elev
			localOrdersCh <- elev.Orders
			globalOrdersCh <- elev.Orders
			elevatorsCh <- elevators
		case <-time.After(50 * time.Millisecond):
		}
	}
}
