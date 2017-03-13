package network

import (
	. "../def"
)

type Buffer struct {
	StateEvents map[string]StateEvent
	OrderEvents map[string]OrderEvent
}

func NewBuffer() Buffer {
    return Buffer{make(map[string]StateEvent), make(map[string]OrderEvent)}
}

func (b *Buffer) AppendStateEvent(id string, stateEvent StateEvent) {
	b.StateEvents[id] = stateEvent
}

func (b *Buffer) AppendOrderEvent(id string, orderEvent OrderEvent) {
	b.OrderEvents[id] = orderEvent
}

func (b *Buffer) RemoveStateEvent(id string) {
    delete(b.StateEvents, id)
}

func (b *Buffer) RemoveOrderEvent(id string) {
    delete(b.OrderEvents, id)
}