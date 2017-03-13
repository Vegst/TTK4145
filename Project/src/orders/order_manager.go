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
		case orderEvent := <-elevatorEvents.Order:
			assignedId := OrderAssigner(id, orderEvent, elevators)
			elev := elevators[assignedId]
			networkEvents.TxMessageOrder <- MessageOrder{id, assignedId, orderEvent}
			elev.Orders[orderEvent.Floor][orderEvent.Type] = orderEvent.Flag
			elevators[assignedId] = elev

			guiEvents.Elevators <- misc.CopyElevators(elevators)
			networkEvents.Elevators <- misc.CopyElevators(elevators)
			if assignedId == id {
				elevatorEvents.LocalOrders <- elevators[id].Orders
			}
			elevatorEvents.GlobalOrders <- misc.GlobalOrders(elevators)

		case state := <-elevatorEvents.State:
			elev := elevators[id]
			elev.State = state
			elevators[id] = elev
			guiEvents.Elevators <- misc.CopyElevators(elevators)
			networkEvents.Elevators <- misc.CopyElevators(elevators)
			networkEvents.TxMessageState <- MessageState{id, id, state}

		case newElevator := <-networkEvents.ElevatorNew:
			elevators[newElevator] = Elevator{
				State: ElevatorState{Floor: -1, Direction: DirnUp, Behaviour: ElevatorBehaviourMoving},
				Orders: Orders{{}},
			}
			guiEvents.Elevators <-misc.CopyElevators(elevators)
			networkEvents.Elevators <- misc.CopyElevators(elevators)

		case lostElevator := <-networkEvents.ElevatorLost:
			lostOrders := misc.CopyOrders(elevators[lostElevator].Orders)
			delete(elevators, lostElevator)
			for f,floorOrders := range lostOrders {
				for t,order := range floorOrders {
					if order {
						orderEvent := OrderEvent{f,OrderType(t),true}
						assignedId := OrderAssigner(id, orderEvent, elevators)
						elev := elevators[assignedId]
						networkEvents.TxMessageOrder <- MessageOrder{id, assignedId, orderEvent}
						elev.Orders[orderEvent.Floor][orderEvent.Type] = true
						elevators[assignedId] = elev
						if assignedId == id {
							elevatorEvents.LocalOrders <- elevators[id].Orders
						}
						elevatorEvents.GlobalOrders <- misc.GlobalOrders(elevators)
					}
				}
			}

			guiEvents.Elevators <-misc.CopyElevators(elevators)
			networkEvents.Elevators <- misc.CopyElevators(elevators)

		case assignedState := <-networkEvents.RxMessageState:
			elev := elevators[assignedState.Id]
			elev.State = assignedState.State
			elevators[assignedState.Id] = elev
			guiEvents.Elevators <- misc.CopyElevators(elevators)
			networkEvents.Elevators <- misc.CopyElevators(elevators)

		case assignedOrder := <-networkEvents.RxMessageOrder:
			elev := elevators[assignedOrder.Id]
			elev.Orders[assignedOrder.OrderEvent.Floor][assignedOrder.OrderEvent.Type] = assignedOrder.OrderEvent.Flag
			elevators[assignedOrder.Id] = elev
			guiEvents.Elevators <- misc.CopyElevators(elevators)
			networkEvents.Elevators <- misc.CopyElevators(elevators)

			if assignedOrder.Id == id {
				elevatorEvents.LocalOrders <- elevators[id].Orders
			}
			elevatorEvents.GlobalOrders <- misc.GlobalOrders(elevators)

		case <-time.After(50 * time.Millisecond):
		}
	}
}