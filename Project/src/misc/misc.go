package misc

import(
	. "../def"
)

func B2I(b bool) int {
	if b {
		return 1
	}
	return 0
}

func CopyElevators(elevators Elevators) Elevators {
	newElevators := make(Elevators)
	for k,v := range elevators {
		newElevators[k] = v
	}
	return newElevators
}

func CopyOrders(orders Orders) Orders {
	var newOrders Orders
	for f,_ := range orders {
		for t,_ := range orders[f] {
			newOrders[f][t] = orders[f][t]
		}
	}
	return newOrders
}

func Union(differentOrders []Orders) Orders {
	var newOrders Orders
	for _,orders := range differentOrders {
		for f,_ := range orders {
			for t,_ := range orders[f] {
				newOrders[f][t] = newOrders[f][t] || orders[f][t]
			}
		}
	}
	return newOrders
}

func GlobalOrders(elevators Elevators) Orders {
	var newOrders Orders
	for _,elevator := range elevators {
		for f,floorOrders := range elevator.Orders {
			for t,order := range floorOrders {
				newOrders[f][t] = newOrders[f][t] || order
			}
		}
	}
	return newOrders
}