package network

import (
	//"./bcast"
	//"./localip"
	. "../def"
	"./peers"
	"./bcast"
	//"fmt"
	"time"
)

const interval = 15 * time.Millisecond

type OrderMessage struct {
	Source string
	Id string
	OrderEvent OrderEvent
}
type StateMessage struct {
	Source string
	Id string
	StateEvent StateEvent
}

type StateAck struct {
	Source string
	Id string
}
type OrderAck struct {
	Source string
	Id string
}

func Network(id string, ordersEvents OrdersNetworkEvents) {

	var elevators Elevators

	buffer := NewBuffer()
	messageId := 0
	
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	txStateMessageCh := make(chan StateMessage)
	rxStateMessageCh := make(chan StateMessage)

	txOrderMessageCh := make(chan OrderMessage)
	rxOrderMessageCh := make(chan OrderMessage)

	txStateAckCh := make(chan StateAck)
	rxStateAckCh := make(chan StateAck)

	txOrderAckCh := make(chan OrderAck)
	rxOrderAckCh := make(chan OrderAck)

	go peers.Transmitter(20004, id, peerTxEnable)
	go peers.Receiver(20004, peerUpdateCh)

	go bcast.Transmitter(26004, txStateMessageCh, txOrderMessageCh, txStateAckCh, txOrderAckCh)
	go bcast.Receiver(26004, rxStateMessageCh, rxOrderMessageCh, rxStateAckCh, rxOrderAckCh)

	for {
		select {

		// Send messages
		case stateEvent := <-ordersEvents.TxStateEvent:
			buffer.AppendStateEvent(id+string(messageId), stateEvent)
			messageId++

		case orderEvent := <-ordersEvents.TxOrderEvent:
			buffer.AppendOrderEvent(id+string(messageId), orderEvent)
			messageId++

		// Receive messages
		case stateMessage := <-rxStateMessageCh:
			if stateMessage.Source != id {
				ordersEvents.RxStateEvent <-stateMessage.StateEvent
				txStateAckCh <-StateAck{id, stateMessage.Id}
			}

		case orderMessage := <-rxOrderMessageCh:
			if orderMessage.Source != id {
				ordersEvents.RxOrderEvent <-orderMessage.OrderEvent
				txOrderAckCh <-OrderAck{id, orderMessage.Id}
			}

		// Receive acks
		case stateAck := <-rxStateAckCh:
			if stateAck.Source != id {
				buffer.RemoveStateEvent(stateAck.Id)
			}

		case orderAck := <-rxOrderAckCh:
			if orderAck.Source != id {
				buffer.RemoveOrderEvent(orderAck.Id)
			}

		// New/Lost id
		case peerUpdate := <-peerUpdateCh:
			if peerUpdate.New != "" {
				if peerUpdate.New != id {
					ordersEvents.ElevatorNew <-peerUpdate.New
					buffer.AppendStateEvent(id+string(messageId), StateEvent{id, elevators[id].State})
					messageId++
					// Merge
					for f,floorOrders := range elevators[id].Orders {
						for t,order := range floorOrders {
							if order {
								buffer.AppendOrderEvent(id+string(messageId), OrderEvent{id, Order{f,OrderType(t), true}})
								messageId++
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

		case <-time.After(interval):
			for messageId, stateEvent := range buffer.StateEvents {
				txStateMessageCh <-StateMessage{id, messageId, stateEvent}
			}
			for messageId, orderEvent := range buffer.OrderEvents {
				txOrderMessageCh <-OrderMessage{id, messageId, orderEvent}
			}
		}
	}
}