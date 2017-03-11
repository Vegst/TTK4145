package orders

import (
	"time"
	. "../def"
)


func OrderManager(orderEventCh <-chan OrderEvent, stateCh <-chan Elevator, localOrdersCh chan Orders, globalOrdersCh chan Orders) {
	var orders Orders
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