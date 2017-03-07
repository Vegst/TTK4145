package elevator

import (
	"time"
)

type OrderEvent struct{
	Floor int
	Type OrderType
	Flag bool
}



func OrderManager(orderEventCh <-chan OrderEvent, stateCh <-chan Elevator, localOrdersCh chan [NumFloors][NumTypes] bool, globalOrdersCh chan [NumFloors][NumTypes] bool) {
	var orders [NumFloors][NumTypes] bool
	for {
		select {
		case orderEvent := <- orderEventCh:
			orders[orderEvent.Floor][orderEvent.Type] = orderEvent.Flag
			localOrdersCh <- orders
			globalOrdersCh <- orders
		case <- stateCh:
		case <-time.After(50*time.Millisecond):
		}
	}
}