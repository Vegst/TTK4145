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

	//var elevators Elevators
	
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	txAssignedStateCh := make(chan AssignedState)
	rxAssignedStateCh := make(chan AssignedState)

	txAssignedOrderCh := make(chan AssignedOrder)
	rxAssignedOrderCh := make(chan AssignedOrder)

	go bcast.Transmitter(16569, txAssignedStateCh, txAssignedOrderCh)
	go bcast.Receiver(16569, rxAssignedStateCh, rxAssignedOrderCh)

	for {
		select {
		case assignedState := <-ordersEvents.TxAssignedState:
			txAssignedStateCh <-assignedState

		case assignedState := <-rxAssignedStateCh:
			if assignedState.Id != id {
				ordersEvents.RxAssignedState <-assignedState
			}

		case assignedOrder := <-ordersEvents.TxAssignedOrder:
			txAssignedOrderCh <-assignedOrder

		case assignedOrder := <-rxAssignedOrderCh:
			if assignedOrder.Id != id {
				ordersEvents.RxAssignedOrder <-assignedOrder
			}

		case peerUpdate := <-peerUpdateCh:
			if peerUpdate.New != "" {
				if peerUpdate.New != id {
					ordersEvents.ElevatorNew <-peerUpdate.New
					txAssignedStateCh <-AssignedState{id, ElevatorState{}}
				}
			}
			for _,lostElevator := range peerUpdate.Lost {
				if lostElevator != id {
					ordersEvents.ElevatorLost <-lostElevator
				}
			}
		case <-ordersEvents.Elevators:

		}
	}
}
