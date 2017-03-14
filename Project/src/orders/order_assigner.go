package orders

import (
	. "../def"
	"math"
	"fmt"
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

func CalculateCost(o Order, e Elevator) float64 {
	const TimeBetweenFloors = 2

	cost := 0

	if(e.State.Direction == DirnDown && o.Floor >= e.State.Floor){
		cost += 3 * numOrdersBelowToFloor(e, 0)
		cost += ((e.State.Floor - 0) * TimeBetweenFloors)*2
		cost += 3 * numOrdersAboveToFloor(e, o.Floor)
		cost += (o.Floor - e.State.Floor) * TimeBetweenFloors
	} else if(e.State.Direction == DirnUp && o.Floor < e.State.Floor){
		cost += 3 * numOrdersAboveToFloor(e, NumFloors)
		cost += ((NumFloors - e.State.Floor) * TimeBetweenFloors)*2
		cost += 3 * numOrdersBelowToFloor(e, o.Floor)
		cost += (e.State.Floor - o.Floor) * TimeBetweenFloors
	} else if(e.State.Direction == DirnDown && o.Floor < e.State.Floor){
		cost += 3 * numOrdersBelowToFloor(e, o.Floor)
		cost += (e.State.Floor - o.Floor) * TimeBetweenFloors
	} else if(e.State.Direction == DirnUp && o.Floor > e.State.Floor){ 
		cost += 3 * numOrdersAboveToFloor(e, o.Floor)
		cost += (o.Floor - e.State.Floor ) * TimeBetweenFloors
	}
	return float64(cost)
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
	fmt.Printf("It will take %g seconds to reach floor %d",eCost, o.Floor)
	return assignedId
}
