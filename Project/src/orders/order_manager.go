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
			networkEvents.TxOrderEvent <- OrderEvent{assignedId, orderEvent}
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
			networkEvents.TxStateEvent <- StateEvent{id, state}

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
						order := Order{f,OrderType(t),true}
						assignedId := OrderAssigner(id, order, elevators)
						elev := elevators[assignedId]
						networkEvents.TxOrderEvent <- OrderEvent{assignedId, order}
						elev.Orders[order.Floor][order.Type] = true
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

		case messageState := <-networkEvents.RxStateEvent:
			elev := elevators[messageState.Target]
			elev.State = messageState.State
			elevators[messageState.Target] = elev
			guiEvents.Elevators <- misc.CopyElevators(elevators)
			networkEvents.Elevators <- misc.CopyElevators(elevators)

		case orderEvent := <-networkEvents.RxOrderEvent:
			elev := elevators[orderEvent.Target]
			elev.Orders[orderEvent.Order.Floor][orderEvent.Order.Type] = orderEvent.Order.Flag
			elevators[orderEvent.Target] = elev
			guiEvents.Elevators <- misc.CopyElevators(elevators)
			networkEvents.Elevators <- misc.CopyElevators(elevators)

			if orderEvent.Target == id {
				elevatorEvents.LocalOrders <- elevators[id].Orders
			}
			elevatorEvents.GlobalOrders <- misc.GlobalOrders(elevators)

		case <-time.After(50 * time.Millisecond):
		}
	}
}