package network

import (
	"./bcast"
	//"./localip"
	. "../def"
	"./peers"
	"fmt"
	//"time"
)

type Message struct {
	orders Orders
}

func broadcaster(){

}

func Network(txStateCh chan State, txOrderCh chan State, rxStateCh <-chan Message , rxOrderCh <-chan Message, ;
	peerTxEnable chan bool, peerUpdateCh <-chan peers.PeerUpdate, localOrdersCh <-chan Orders, stateCh <-chan ElevatorState, ;
	assignedOrderCh <- chan AssignedOrder, orderEventCh chan OrderEvent) {
	var orders Orders
	for {
		select {
		case rxState := <-rxStateCh:

		case rxOrder := <-rxOrderCh:

		case peerUpdate := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", peerUpdate.Peers)
			fmt.Printf("  New:      %q\n", peerUpdate.New)
			fmt.Printf("  Lost:     %q\n", peerUpdate.Lost)
		case elevState = <- stateCh:
			txCh <- elevState
		case a := assignedOrderCh:
			pendingOrders[a.OrderEvent.Floor][a.OrderEvent.Type] == a.OrderEvent
			txCh <- assignedOrder
		case
}
