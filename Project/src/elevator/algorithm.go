//

package elevator

import (
	. "../def"
)

func OrderAtFloor(elev Elevator) bool{
	for o := 0; o < NumTypes; o++{
		if(elev.Orders[elev.State.Floor][o]){
			return true
		}
	}
	return false
}

func GetDirection(elev Elevator) MotorDirection{
	switch(elev.State.Direction) {
	case DirnUp:
		if(checkAbove(elev)) {
			return DirnUp
		} else if(checkBelow(elev)) {
			return DirnDown
		} else{
			return DirnStop
		}
	case DirnStop, DirnDown:
		if(checkBelow(elev)) {
			return DirnDown
		} else if(checkAbove(elev)) {
			return DirnUp
		} else{
			return DirnStop
		}	
	}
	return DirnStop
}

func ShouldStop(elev Elevator) bool{
	switch(elev.State.Direction) {
	case DirnDown:
		return 	elev.Orders[elev.State.Floor][OrderCallCommand] ||
				elev.Orders[elev.State.Floor][OrderCallDown] ||
				!checkBelow(elev)
	case DirnUp:
		return 	elev.Orders[elev.State.Floor][OrderCallCommand] ||
				elev.Orders[elev.State.Floor][OrderCallUp] ||
				!checkAbove(elev) 	
	}
	return true
}

func checkAbove (elev Elevator) bool{
	for f := elev.State.Floor+1; f < NumFloors; f++{
		for o := 0; o < NumTypes; o++{
			if(elev.Orders[f][o]){
				return true
			}
		}
	}
	return false
}

func checkBelow (elev Elevator) bool{
	for f := 0; f < elev.State.Floor; f++{
		for o := 0; o < NumTypes; o++{
			if(elev.Orders[f][o]){
				return true
			}
		}
	}
	return false
}