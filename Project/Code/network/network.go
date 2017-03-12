package network

import (
	//"./bcast"
	//"./localip"
	. "../def"
	//"./peers"
	//"fmt"
	//"time"
)

type Message struct {
	orders Orders
}

func broadcaster() {

}

func Network(ID string, ordersEvents OrdersNetworkEvents) {
	//var elevator Elevator
	/*
	for {
		select {
		//Change state
		case assignment := <- ordersEvents.TxAssignedOrder:
			fmt.Println("Sent assignment from ", ID, " to ", assignment.ID)
			//txAssignmentCh <- assignment
		case assignment := <- rxAssignmentCh:
			if(assignment.ID == ID){
				fmt.Println("Received assignment")
				assignedOrderCh <- assignment.OrderEvent
			}

		case peerUpdate := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", peerUpdate.Peers)
			fmt.Printf("  New:      %q\n", peerUpdate.New)
			fmt.Printf("  Lost:     %q\n", peerUpdate.Lost)
		}
	}
	*/
}
