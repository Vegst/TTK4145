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
	Message string
	Iter    int
}

func Network(txCh chan Message, rxCh <-chan Message, peerTxEnable chan bool, peerUpdateCh <-chan peers.PeerUpdate, localOrdersCh <-chan Orders, orderEventCh chan OrderEvent) {
	fmt.Println("Started")
	for {
		select {
		case rxMsg := <-rxCh:
			fmt.Printf("Peer update:%q\n", rxMsg)
		case peerUpdate := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", peerUpdate.Peers)
			fmt.Printf("  New:      %q\n", peerUpdate.New)
			fmt.Printf("  Lost:     %q\n", peerUpdate.Lost)
		case orders := <-localOrdersCh:

		}
	}
}
