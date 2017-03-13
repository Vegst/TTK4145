package network

import (
	//"./bcast"
	//"./localip"
	. "../def"
	"./peers"
	"./bcast"
	//"fmt"
	//"time"
)

func Network(id string, ordersEvents OrdersNetworkEvents) {

	var elevators Elevators
	
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	go peers.Transmitter(20004, id, peerTxEnable)
	go peers.Receiver(20004, peerUpdateCh)

	txMessageStateCh := make(chan MessageState)
	rxMessageStateCh := make(chan MessageState)

	txMessageOrderCh := make(chan MessageOrder)
	rxMessageOrderCh := make(chan MessageOrder)

	go bcast.Transmitter(26004, txMessageStateCh, txMessageOrderCh)
	go bcast.Receiver(26004, rxMessageStateCh, rxMessageOrderCh)

	for {
		select {
		case messageState := <-ordersEvents.TxMessageState:
			txMessageStateCh <-messageState

		case messageOrder := <-ordersEvents.TxMessageOrder:
			txMessageOrderCh <-messageOrder

		case assignedState := <-rxMessageStateCh:
			if assignedState.Source != id {
				ordersEvents.RxMessageState <-assignedState
			}
		case assignedOrder := <-rxMessageOrderCh:
			if assignedOrder.Source != id {
				ordersEvents.RxMessageOrder <-assignedOrder
			}

		case peerUpdate := <-peerUpdateCh:
			if peerUpdate.New != "" {
				if peerUpdate.New != id {
					ordersEvents.ElevatorNew <-peerUpdate.New
					txMessageStateCh <-MessageState{id, id, elevators[id].State}
					// Merge
					for f,floorOrders := range elevators[id].Orders {
						for t,order := range floorOrders {
							if order {
								txMessageOrderCh <-MessageOrder{id, id, OrderEvent{f,OrderType(t), true}}
							}
						}
					}
				}
			}
			for _,lostElevator := range peerUpdate.Lost {
				if lostElevator != id {
					ordersEvents.ElevatorLost <-lostElevator
				}
			}
		case elevators = <-ordersEvents.Elevators:

		}
	}
}
