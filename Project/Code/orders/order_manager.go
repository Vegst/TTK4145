package orders

import (
	. "../def"
	"../elevator"
	//"fmt"
	"time"
)

func OrderManager(id string, orderEventCh <-chan OrderEvent, assignedOrderCh <-chan OrderEvent, assignmentCh chan Assignment, stateCh <-chan ElevatorState, updateElevatorCh <-chan Elevator, localOrdersCh chan Orders, globalOrdersCh chan Orders, elevatorCh chan Elevator, elevatorsCh chan Elevators) {

	var elevators Elevators
	var orders Orders

	elevators = make(Elevators)
	elevators[id] = Elevator{ElevatorState{0, DirnStop, ElevatorBehaviourIdle}, orders, id}
	elevatorsCh <- elevators

	for {
		select {

		case elev := <-updateElevatorCh:
			elevators[id] = elev
			elevatorsCh <- elevators

		case orderEvent := <-orderEventCh:
			elev := elevators[id]

			//Assigning an order to an elevator
			AssignmentID := elevator.OrderAssigner(orderEvent, elevators)
			assignmentCh <- Assignment{orderEvent, AssignmentID}

			elev.Orders[orderEvent.Floor][orderEvent.Type] = orderEvent.Flag
			elevators[id] = elev
			elevatorCh <- elev

		case <-stateCh:
		
		case assignedOrder := <-assignedOrderCh:
			elev := elevators[id]
			elev.Orders[assignedOrder.Floor][assignedOrder.Type] = assignedOrder.Flag
			elevators[id] = elev
			localOrdersCh <- elev.Orders
			globalOrdersCh <- elev.Orders
			elevatorsCh <- elevators


		case <-time.After(50 * time.Millisecond):
		}
	}
}
