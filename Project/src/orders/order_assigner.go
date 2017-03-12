package orders

import (
	. "../def"
	"math"
	"math/rand"
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

func OrderAssigner(o OrderEvent, elevs Elevators) string {
	var id string = ""
	eCost := math.Inf(1)
	r := rand.Intn(2) + 1
	for k := range elevs {
		iCost := float64(CalculateCost(o, elevs[k]))
		if iCost < eCost {
			id = k
			eCost = iCost
		}
	}
	if(r == 1){
		id = "Heis1"
	} else if(r == 2){
		id = "Heis2"
	}
	return id
}
