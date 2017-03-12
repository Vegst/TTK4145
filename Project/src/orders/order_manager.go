package orders

import (
	. "../def"
	//"fmt"
	"time"
	"../misc"
)

func OrderManager(id string, elevatorEvents ElevatorOrdersEvents, networkEvents OrdersNetworkEvents, guiEvents OrdersGuiEvents) {

	var elevators Elevators

	elevators = make(Elevators)

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
			networkEvents.TxAssignedOrder <- AssignedOrder{id, orderEvent}

			elev.Orders[orderEvent.Floor][orderEvent.Type] = orderEvent.Flag
			elevators[id] = elev
			guiEvents.Elevators <- misc.Copy(elevators)
			networkEvents.Elevators <- misc.Copy(elevators)
			elevatorEvents.LocalOrders <- elev.Orders
			elevatorEvents.GlobalOrders <- elev.Orders

		case state := <-elevatorEvents.State:
			elev := elevators[id]
			elev.State = state
			elevators[id] = elev
			guiEvents.Elevators <- misc.Copy(elevators)
			networkEvents.Elevators <- misc.Copy(elevators)
			networkEvents.TxAssignedState <- AssignedState{id, state}

		case newElevator := <-networkEvents.ElevatorNew:
			elevators[newElevator] = Elevator{
				State: ElevatorState{Floor: -1, Direction: DirnUp, Behaviour: ElevatorBehaviourMoving},
				Orders: Orders{{}},
			}
			guiEvents.Elevators <-misc.Copy(elevators)
			networkEvents.Elevators <- misc.Copy(elevators)

		case lostElevator := <-networkEvents.ElevatorLost:
			delete(elevators, lostElevator)
			guiEvents.Elevators <-misc.Copy(elevators)
			networkEvents.Elevators <- misc.Copy(elevators)

		case assignedState := <-networkEvents.RxAssignedState:
			elev := elevators[assignedState.Id]
			elev.State = assignedState.State
			elevators[assignedState.Id] = elev
			guiEvents.Elevators <- misc.Copy(elevators)
			networkEvents.Elevators <- misc.Copy(elevators)

		case assignedOrder := <-networkEvents.RxAssignedOrder:
			elev := elevators[assignedOrder.Id]
			elev.Orders[assignedOrder.OrderEvent.Floor][assignedOrder.OrderEvent.Type] = assignedOrder.OrderEvent.Flag
			elevators[assignedOrder.Id] = elev
			guiEvents.Elevators <- misc.Copy(elevators)
			networkEvents.Elevators <- misc.Copy(elevators)

			if assignedOrder.Id == id {
				elevatorEvents.LocalOrders <- elev.Orders
			}
			elevatorEvents.GlobalOrders <- misc.GlobalOrders(elevators)
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
