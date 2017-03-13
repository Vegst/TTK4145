package orders

import (
	. "../def"
	"math"
)

func CalculateCost(o OrderEvent, e Elevator) float64 {
	cost := float64(0)
	cost += math.Abs(float64((o.Floor - e.State.Floor) * 10))

	//Adds cost for each f
	for i := 0; i < NumFloors; i++ {
		for j := 0; j < NumTypes; j++ {
			if e.Orders[i][j] {
				cost += 3
			}
			break
		}
	}

	if o.Floor < e.State.Floor && e.State.Direction == DirnUp ||
		o.Floor > e.State.Floor && e.State.Direction == DirnDown {
		cost += 20
	}
	return cost

}

func OrderAssigner(id string, o OrderEvent, elevs Elevators) string {
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
