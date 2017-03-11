package elevator

import (
	. "../def"
	"fmt"
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
	fmt.Println("Cost of new order for Elevator: ", cost)
	return cost

}

func OrderAssigner(o OrderEvent, elevs Elevators) string {
	var id string = ""
	eCost := math.Inf(1)
	for k := range elevs {
		iCost := float64(CalculateCost(o, elevs[k]))
		if iCost < eCost {
			id = k
			eCost = iCost
		}
	}
	return id
}
