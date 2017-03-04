package orders

import (
	"../logic"
	"../elevator"
	"time"
)

type Orders [elevator.NumFloors][elevator.NumOrderTypes] bool

type OrderType int

const (
	OrderCallUp      OrderType = 0
	OrderCallDown    OrderType = 1
	OrderCallCommand OrderType = 2
)

type OrderEvent struct{
	Floor int
	Type OrderType
	Flag bool
}



func OrderManager(orderEventCh <-chan OrderEvent, stateCh <-chan elevator.State, localOrdersCh chan Orders, globalOrdersCh chan Orders) {
	var orders Orders
	for {
		select {
		case orderEvent := <- orderEventCh:
			orders[orderEvent.Floor][orderEvent.Type] = orderEvent.Flag
			localOrdersCh <- orders
			globalOrdersCh <- orders
		case <- stateCh:
		case <-time.After(50*time.Milliseconds):
		}
	}
}