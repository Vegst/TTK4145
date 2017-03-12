package orders

import (
	. "../def"
	//"fmt"
	"time"
)

func OrderManager(id string, elevatorEvents ElevatorOrdersEvents, networkEvents OrdersNetworkEvents, guiEvents OrdersGuiEvents) {

	var elevators Elevators
	var orders Orders

	elevators = make(Elevators)
	elevators[id] = Elevator{ElevatorState{0, DirnStop, ElevatorBehaviourIdle}, orders, id}
	guiEvents.Elevators <- elevators

	for {
		select {

			/*
		case elev := <-updateElevatorCh:
			elevators[id] = elev
			elevatorsCh <- elevators
	*/
		case orderEvent := <-elevatorEvents.Order:
			elev := elevators[id]

			//Assigning an order to an elevator
			//AssignmentID := elevator.OrderAssigner(orderEvent, elevators)
			//assignmentCh <- Assignment{orderEvent, AssignmentID}

			elev.Orders[orderEvent.Floor][orderEvent.Type] = orderEvent.Flag
			elevators[id] = elev
			guiEvents.Elevators <- elevators
			elevatorEvents.LocalOrders <- elev.Orders
			elevatorEvents.GlobalOrders <- elev.Orders

		case <-elevatorEvents.State:
		
		/*
		case assignedOrder := <-assignedOrderCh:
			elev := elevators[id]
			elev.Orders[assignedOrder.Floor][assignedOrder.Type] = assignedOrder.Flag
			elevators[id] = elev
			localOrdersCh <- elev.Orders
			globalOrdersCh <- elev.Orders
			elevatorsCh <- elevators
		*/

		case <-time.After(50 * time.Millisecond):
		}
	}
}
