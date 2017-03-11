package orders

import (
	. "../def"
	"time"
)

func OrderManager(orderEventCh <-chan OrderEvent, stateCh <-chan Elevator, localOrdersCh chan Orders, globalOrdersCh chan Orders, assignedOrderCh chan AssignedOrder) {
	var orders Orders
	for {
		select {
		case orderEvent := <-orderEventCh:
			assignedOrder := AssignerOrder{orderEvent, OrderAssigner{orderEvent}}
			orders[orderEvent.Floor][orderEvent.Type] = orderEvent.Flag
			localOrdersCh <- orders
			globalOrdersCh <- orders
		case <-stateCh:
		case <-time.After(50 * time.Millisecond):
		}
	}
}
