package orders

import (
	. "../def"
	"time"
	//"math"
	"math/rand"
	"../elevator"
	"../misc"
)

func CalculateCost(o Order, e Elevator) time.Duration {

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
