package orders

import (
	. "../def"
	"time"
	//"math"
	"math/rand"
	"../elevator"
	"../misc"
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

func CalculateCost(order Order, elev Elevator) time.Duration {
	e := misc.CopyElevator(elev)
	e.Orders[order.Floor][order.Type] = order.Flag

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
            e = ClearOrdersAtCurrentFloor(e)
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
	/*
	var assignedId string = id
	eDur := math.Inf(1)
	for k := range elevs {
		iDur := float64(CalculateCost(o, elevs[k]))
		if iDur < eDur {
			assignedId = k
			eDur = iDur
		}
	}
	*/
	ids := make([]string, 0)
	for id,_ := range elevs {
		ids = append(ids, id)
	}
	rand.Seed(time.Now().UTC().UnixNano())
	return ids[rand.Intn(len(ids))]
	//return assignedId
}
