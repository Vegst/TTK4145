package logic

func checkAbove (localOrders [][]bool, state State) bool{
	for f := state.Floor; f < NumFloors; f++{
		for o := 0; o < NumOrderTypes; o++{
			if(localOrders[f][o]){
				return true
			}
		}
	}
}

func checkBelow (localOrders [][]bool, state State) bool{
	for f := 0; f < state.Floor; f++{
		for o := 0; o < NumOrderTypes; o++{
			if(localOrders[f][o]){
				return true
			}
		}
	}
}

func checkFloor (localOrders [][]bool, state State) bool{
	for o := 0; o < NumOrderTypes; o++{
			if(localOrders[f][o]){
				return true
			}
		}
}

func GetDirection (localOrders [][]bool, state State) elevator.MotorDirection{
	ordersAbove := checkAbove(localOrders, state)
	ordersBelow := checkBelow(localOrders, state)

	if(!ordersAbove && !ordersBelow){
		return elevator.DirnStop
	} else if(state.Direction == elevator.DirnUp || state.Direction == elevator.DirnStop){
		if(ordersAbove){
			return elevator.DirnUp
		}else if(ordersBelow){
			return elevator.DirnDown
		}
		else{
			return elevator.DirnStop
		}
	}
	else{
		if(ordersBelow){
			return elevator.DirnDown
		}else if(ordersAbove){
			return elevator.DirnUp
		}
		else{
			return elevator.DirnStop
		}
	}
}

func ShouldStop (localOrders [][]bool, state State) bool{

	if(localOrders [state.Floor][OrderCallCommand]){
		return true
	} else if(state.Direction == elevator.DirnUp){
		if(localOrders[state.Floor][OrderCallUp]){
			return true
		}
	} else{
		if(localOrders[state.Floor][OrderCallDown]){
			return true
		}
	}
	return false
}