package orders

import (
	"../backup"
	. "../def"
	"../misc"
)

type StateMachine struct {
	Id             string
	ElevatorEvents ElevatorOrdersEvents
	NetworkEvents  OrdersNetworkEvents
	GuiEvents      OrdersGuiEvents

	Elevators Elevators
}

func NewStateMachine(id string, elevatorEvents ElevatorOrdersEvents, networkEvents OrdersNetworkEvents, guiEvents OrdersGuiEvents) *StateMachine {
	sm := new(StateMachine)
	sm.Id = id
	sm.ElevatorEvents = elevatorEvents
	sm.NetworkEvents = networkEvents
	sm.GuiEvents = guiEvents
	sm.Elevators = make(Elevators)
	elevator := sm.Elevators[id]
	elevator.Orders = backup.ReadFromBackup(BackupFile)
	sm.Elevators[id] = elevator
	return sm
}

func (this *StateMachine) OnOrderReceived(order Order) {
	assignedId := this.Id
	if order.Flag {
		assignedId = OrderAssigner(this.Id, order, this.Elevators)
	}
	elev := this.Elevators[assignedId]
	elev.Orders[order.Floor][order.Type] = order.Flag
	this.Elevators[assignedId] = elev

	this.NetworkEvents.TxOrderEvent <- OrderEvent{assignedId, order}
	this.GuiEvents.Elevators <- misc.CopyElevators(this.Elevators)
	this.NetworkEvents.Elevators <- misc.CopyElevators(this.Elevators)
	if assignedId == this.Id {
		this.ElevatorEvents.LocalOrders <- this.Elevators[this.Id].Orders
	}
	this.ElevatorEvents.GlobalOrders <- misc.GlobalOrders(this.Elevators)
	backup.WriteToBackup(BackupFile, elev.Orders)
}

func (this *StateMachine) OnStateReceived(state ElevatorState) {
	elev := this.Elevators[this.Id]
	elev.State = state
	this.Elevators[this.Id] = elev
	this.GuiEvents.Elevators <- misc.CopyElevators(this.Elevators)
	this.NetworkEvents.Elevators <- misc.CopyElevators(this.Elevators)
	this.NetworkEvents.TxStateEvent <- StateEvent{this.Id, state}
}

func (this *StateMachine) OnElevatorNew(id string) {
	this.Elevators[id] = Elevator{
		State:  ElevatorState{Floor: -1, Direction: DirnUp, Behaviour: ElevatorBehaviourMoving},
		Orders: Orders{{}},
	}
	this.GuiEvents.Elevators <- misc.CopyElevators(this.Elevators)
	this.NetworkEvents.Elevators <- misc.CopyElevators(this.Elevators)
}

func (this *StateMachine) OnElevatorLost(id string) {
	lostOrders := misc.CopyOrders(this.Elevators[id].Orders)
	delete(this.Elevators, id)
	for f, _ := range lostOrders {
		for t, _ := range []OrderType{OrderCallDown, OrderCallUp} {
			if lostOrders[f][t] {
				lostOrder := Order{f, OrderType(t), true}
				assignedId := OrderAssigner(this.Id, lostOrder, this.Elevators)
				elev := this.Elevators[assignedId]
				elev.Orders[f][t] = true
				this.Elevators[assignedId] = elev

				this.NetworkEvents.TxOrderEvent <- OrderEvent{assignedId, lostOrder}
				if assignedId == this.Id {
					this.ElevatorEvents.LocalOrders <- this.Elevators[this.Id].Orders
				}
				this.ElevatorEvents.GlobalOrders <- misc.GlobalOrders(this.Elevators)
			}
		}
	}

	this.GuiEvents.Elevators <- misc.CopyElevators(this.Elevators)
	this.NetworkEvents.Elevators <- misc.CopyElevators(this.Elevators)
}

func (this *StateMachine) OnStateEventReceived(stateEvent StateEvent) {
	elev := this.Elevators[stateEvent.Target]
	elev.State = stateEvent.State
	this.Elevators[stateEvent.Target] = elev
	this.GuiEvents.Elevators <- misc.CopyElevators(this.Elevators)
	this.NetworkEvents.Elevators <- misc.CopyElevators(this.Elevators)
}

func (this *StateMachine) OnOrderEventReceived(orderEvent OrderEvent) {
	elev := this.Elevators[orderEvent.Target]
	elev.Orders[orderEvent.Order.Floor][orderEvent.Order.Type] = orderEvent.Order.Flag
	this.Elevators[orderEvent.Target] = elev
	this.GuiEvents.Elevators <- misc.CopyElevators(this.Elevators)
	this.NetworkEvents.Elevators <- misc.CopyElevators(this.Elevators)

	if orderEvent.Target == this.Id {
		this.ElevatorEvents.LocalOrders <- this.Elevators[this.Id].Orders
	}
	this.ElevatorEvents.GlobalOrders <- misc.GlobalOrders(this.Elevators)
}
