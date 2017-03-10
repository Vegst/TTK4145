//

package elevator

import (
	. "../def"
)

func OrderAtFloor(localOrders Orders, floor int) bool{
	for o := 0; o < NumTypes; o++{
		if(localOrders[floor][o]){
			return true
		}
	}
	return false
}

func GetDirection(localOrders Orders, elev Elevator) MotorDirection{
	switch(elev.Direction) {
	case DirnUp:
		if(checkAbove(localOrders, elev)) {
			return DirnUp
		} else if(checkBelow(localOrders, elev)) {
			return DirnDown
		} else{
			return DirnStop
		}
	case DirnStop, DirnDown:
		if(checkBelow(localOrders, elev)) {
			return DirnDown
		} else if(checkAbove(localOrders, elev)) {
			return DirnUp
		} else{
			return DirnStop
		}	
	}
	return DirnStop
}

func ShouldStop(localOrders Orders, elev Elevator) bool{
	switch(elev.Direction) {
	case DirnDown:
		return 	localOrders[elev.Floor][OrderCallCommand] ||
				localOrders[elev.Floor][OrderCallDown] ||
				!checkBelow(localOrders, elev)
	case DirnUp:
		return 	localOrders[elev.Floor][OrderCallCommand] ||
				localOrders[elev.Floor][OrderCallUp] ||
				!checkAbove(localOrders, elev) 	
	}
	return true
}

func checkAbove (localOrders Orders, elev Elevator) bool{
	for f := elev.Floor+1; f < NumFloors; f++{
		for o := 0; o < NumTypes; o++{
			if(localOrders[f][o]){
				return true
			}
		}
	}
	return false
}

func checkBelow (localOrders Orders, elev Elevator) bool{
	for f := 0; f < elev.Floor; f++{
		for o := 0; o < NumTypes; o++{
			if(localOrders[f][o]){
				return true
			}
		}
	}
	return false
}