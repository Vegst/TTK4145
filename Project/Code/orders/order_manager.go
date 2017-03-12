package orders

import (
	. "../def"
	//"../elevator"
	"fmt"
	"time"
)

func OrderManager(id string, orderEventCh <-chan OrderEvent, assignedOrderCh <-chan OrderEvent, assignmentCh chan Assignment, stateCh <-chan ElevatorState, updateElevatorCh <-chan Elevator, localOrdersCh chan Orders, globalOrdersCh chan Orders, elevatorCh chan Elevator, elevatorsCh chan Elevators) {
	/*
	var elevators Elevators
	var orders Orders

	elevators = make(Elevators)
	localElevator := Elevator{ElevatorState{0, DirnStop, ElevatorBehaviourIdle}, orders, id}
	elevatorCh <- localElevator
	*/
	for {
		select {
		case <-assignedOrderCh:
			fmt.Println("OrderManager for ", id, " received an order.")
		case <-time.After(1 * time.Second):
			if(id == "Heis1"){
				assignmentCh <- Assignment{OrderEvent{2, OrderCallUp, true}, "Heis2"}
			}
		/*
		case elev := <-updateElevatorCh:
			elevators[id] = elev
			elevatorsCh <- elevators

		case orderEvent := <-orderEventCh:
			elev := localElevator

			//Assigning an order to an elevator
			AssignmentID := elevator.OrderAssigner(orderEvent, elevators)
			assignmentCh <- Assignment{orderEvent, AssignmentID}

			elev.Orders[orderEvent.Floor][orderEvent.Type] = orderEvent.Flag
			elevators[id] = elev
			elevatorCh <- elev

		case <-stateCh:

		
			
			elev := localElevator
			elev.Orders[assignedOrder.Floor][assignedOrder.Type] = assignedOrder.Flag
			localElevator = elev
			localOrdersCh <- elev.Orders
			globalOrdersCh <- elev.Orders
			elevatorsCh <- elevators
			
		case <-time.After(5 * time.Second):
			for k := range elevators{
				fmt.Println("Elevator: %s", k)
			}
		*/
		}
	}
}
