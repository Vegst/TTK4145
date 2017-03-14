package orders

import (
	. "../def"
	"time"
	"math"
	"../elevator"
)

func numOrdersBelowToFloor(e Elevator, floor int) int{
	numOrders := 0
	for f := floor; f < e.State.Floor; f++{
		for o := 0; 0 < NumTypes; o++{
			if e.Orders[f][o] {
				numOrders += 1
			}
			break
		}
	}
	return numOrders
}

func numOrdersAboveToFloor(e Elevator, floor int) int{
	numOrders := 0
	for f := e.State.Floor; f < floor; f++{
		for o := 0; 0 < NumTypes; o++{
			if e.Orders[f][o] {
				numOrders += 1
			}
			break
		}
	}
	return numOrders
}

func CalculateCost(o Order, e Elevator) time.Duration {

	e.Orders[o.Floor][o.Type] = o.Flag

    dur := 0*time.Millisecond
    
    switch e.State.Behaviour {
    case ElevatorBehaviourIdle:
        e.State.Direction = elevator.GetDirection(e)
        if e.State.Direction == DirnStop {
        	return dur
        }
    case ElevatorBehaviourMoving:
        e.State.Floor = e.State.Floor + int(e.State.Direction)
        dur += TravelTime/2
    case ElevatorBehaviourDoorOpen:
        dur -= DoorOpenTime/2
    }
    
    for {
        if elevator.ShouldStop(e) {
            for b := 0; b < NumTypes; b++{
            	e.Orders[e.State.Floor][b] = false
            }
            dur += DoorOpenTime
            e.State.Direction = elevator.GetDirection(e)
            if e.State.Direction == DirnStop{
                return dur
            }
        }
        e.State.Floor = e.State.Floor + int(e.State.Direction)
        dur += TravelTime
    }
}

func OrderAssigner(id string, o Order, elevs Elevators) string {
	if o.Type == OrderCallCommand {
		return id
	}
	var assignedId string = id
	eCost := math.Inf(1)
	for k := range elevs {
		iCost := float64(CalculateCost(o, elevs[k]))
		if iCost < eCost {
			assignedId = k
			eCost = iCost
		}
	}
	return assignedId
}
