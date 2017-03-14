package network

import (
	. "../def"
)

type StateMachine struct {
	id string 
	OrdersEvents OrdersNetworkEvents
	TxStateMessageCh chan StateMessage
	TxOrderMessageCh chan OrderMessage
	TxStateAckCh chan StateAck
	TxOrderAckCh chan OrderAck

	Buffer Buffer
	messageId int
	elevators Elevators
}


func NewStateMachine(id string, ordersEvents OrdersNetworkEvents, txStateMessageCh chan StateMessage, txOrderMessageCh chan OrderMessage, txStateAckCh chan StateAck, txOrderAckCh chan OrderAck) *StateMachine {
	sm := new(StateMachine)
	sm.id = id
	sm.OrdersEvents = ordersEvents
	sm.TxStateMessageCh = txStateMessageCh
	sm.TxOrderMessageCh = txOrderMessageCh
	sm.TxStateAckCh = txStateAckCh
	sm.TxOrderAckCh = txOrderAckCh

	sm.Buffer = NewBuffer()
	sm.messageId = 0
	sm.elevators = make(Elevators)
	return sm
}

func (sm *StateMachine) OnStateEventTransmit(stateEvent StateEvent) {
	sm.Buffer.EnqueueStateMessage(StateMessage{sm.id, sm.id+string(sm.messageId), stateEvent})
	sm.messageId++
}

func (sm *StateMachine) OnOrderEventTransmit(orderEvent OrderEvent) {
	sm.Buffer.EnqueueOrderMessage(OrderMessage{sm.id, sm.id+string(sm.messageId), orderEvent})
	sm.messageId++
}

func (sm *StateMachine) OnStateMessageReceived(stateMessage StateMessage) {
	sm.OrdersEvents.RxStateEvent <-stateMessage.StateEvent
	sm.TxStateAckCh <-StateAck{sm.id, stateMessage.Id}
}

func (sm *StateMachine) OnOrderMessageReceived(orderMessage OrderMessage) {
	sm.OrdersEvents.RxOrderEvent <-orderMessage.OrderEvent
	sm.TxOrderAckCh <-OrderAck{sm.id, orderMessage.Id}
}

func (sm *StateMachine) OnStateAckReceived(stateAck StateAck) {
	if stateAck.Id == sm.Buffer.TopStateMessage().Id {
		sm.Buffer.DequeueStateMessage()
	}
}

func (sm *StateMachine) OnOrderAckReceived(orderAck OrderAck) {
	if orderAck.Id == sm.Buffer.TopOrderMessage().Id {
		sm.Buffer.DequeueOrderMessage()
	}
}

func (sm *StateMachine) OnPeerNew(peer string) {
	sm.OrdersEvents.ElevatorNew <-peer
	sm.Buffer.EnqueueStateMessage(StateMessage{sm.id, sm.id+string(sm.messageId), StateEvent{sm.id, sm.elevators[sm.id].State}})
	sm.messageId++
	// Merge
	for f,_ := range sm.elevators[sm.id].Orders {
		for t,_ := range sm.elevators[sm.id].Orders[f] {
			if sm.elevators[sm.id].Orders[f][t] {
				sm.Buffer.EnqueueOrderMessage(OrderMessage{sm.id, sm.id+string(sm.messageId), OrderEvent{sm.id, Order{f,OrderType(t), true}}})
				sm.messageId++
			}
		}
	}
}

func (sm *StateMachine) OnPeerLost(peer string) {
	sm.OrdersEvents.ElevatorLost <-peer
}

func (sm *StateMachine) OnElevatorsUpdated(elevators Elevators) {
	sm.elevators = elevators
}

func (sm *StateMachine) OnInterval() {
	if sm.Buffer.HasStateMessage() {
		sm.TxStateMessageCh <- sm.Buffer.TopStateMessage()
	}
	if sm.Buffer.HasOrderMessage() {
		sm.TxOrderMessageCh <- sm.Buffer.TopOrderMessage()
	}
}