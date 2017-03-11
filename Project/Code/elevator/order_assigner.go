package elevator

/*
import (
	"../driver"
	"fmt"
	"math"
	. "../def"
)

func CalculateCost(o OrderEvent, e Elevator) float64 {
	cost := float64(0)
	cost += math.Abs(float64(o.Floor - e.Floor))

	//Adds cost for each f
	for i := 0; i < NumFloors; i++ {
		for j := 0; j < NumTypes; j++ {
			if e.Orders[i][j] {
				cost += 0.3
			}
			break
		}
	}

	if o.Floor < e.Floor && e.Direction == DirnUp ||
		o.Floor > e.Floor && e.Direction == DirnDown {
		cost += 2
	}
	fmt.Println("Cost of new order for Elevator: ", cost)
	return cost

}

func OrderAssigner(o OrderEvent, elevs [NumElevators]Elevator) int {
	e := -1
	eCost := math.Inf(1)
	for i := 0; i < NumElevators; i++ {
		iCost := float64(CalculateCost(o, elevs[i]))
		if iCost < eCost {
			e = i
			eCost = iCost
		}
	}
	return e
}
*/