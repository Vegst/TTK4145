package network

import (
	//"./bcast"
	//"./localip"
	. "../def"
	"./peers"
	"fmt"
	//"time"
)

type Message struct {
	orders Orders
}

func broadcaster() {

}

func Network(ID string, txStateCh chan ElevatorState, txNetOrderCh chan NetOrder, rxStateCh <-chan ElevatorState, rxNetOrderCh <-chan NetOrder, peerTxEnable chan bool, peerUpdateCh <-chan peers.PeerUpdate, stateCh <-chan ElevatorState, updateElevatorCh chan Elevator, netOrderCh chan NetOrder, assignedOrderCh chan OrderEvent) {
	var elevator Elevator
	var pendingOrder NetOrder
	for {
		select {
		//Change state
		case elevatorState := <-stateCh:
			elevator.State = elevatorState
			txStateCh <- elevatorState
		//case updateState := <-rxStateCh:
		//	updateStateCh <- updateState
		case pendingOrder = <-netOrderCh:
			fmt.Printf("Pending order on :    %s\n", pendingOrder.ID)
		case rxNetOrder := <-rxNetOrderCh:
			if rxNetOrder.ID == ID {
				fmt.Printf("Received order on :    %s\n", pendingOrder.ID)
				assignedOrderCh <- rxNetOrder.OrderEvent
			}

		case peerUpdate := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", peerUpdate.Peers)
			fmt.Printf("  New:      %q\n", peerUpdate.New)
			fmt.Printf("  Lost:     %q\n", peerUpdate.Lost)
		}
	}
}
