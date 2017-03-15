package orders

import (
	. "../def"
	"time"
	"math"
	"../elevator"
)

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
	var assignedId string = id
	eDur := math.Inf(1)
	for k := range elevs {
		iDur := float64(CalculateCost(o, elevs[k]))
		if iDur < eDur {
			assignedId = k
			eDur = iDur
		}
	}
	return assignedId
}
