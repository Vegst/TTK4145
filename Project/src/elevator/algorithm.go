package elevator

import (
	. "../def"
)

func OrderAtFloor(orders Orders, floor int) bool{
	for o := 0; o < NumTypes; o++{
		if(orders[floor][o]){
			return true
		}
	}
	return false
}

func GetDirection(orders Orders, floor int, direction MotorDirection) MotorDirection{
	switch(direction) {
	case DirnUp:
		if(checkAbove(orders, floor)) {
			return DirnUp
		} else if(checkBelow(orders, floor)) {
			return DirnDown
		} else{
			return DirnStop
		}
	case DirnStop, DirnDown:
		if(checkBelow(orders, floor)) {
			return DirnDown
		} else if(checkAbove(orders, floor)) {
			return DirnUp
		} else{
			return DirnStop
		}	
	}
	return DirnStop
}

func ShouldStop(orders Orders, floor int, direction MotorDirection) bool{
	switch(direction) {
	case DirnDown:
		return 	orders[floor][OrderCallCommand] ||
				orders[floor][OrderCallDown] ||
				!checkBelow(orders, floor)
	case DirnUp:
		return 	orders[floor][OrderCallCommand] ||
				orders[floor][OrderCallUp] ||
				!checkAbove(orders, floor) 	
	}
	return true
}

func checkAbove(orders Orders, floor int) bool{
	for f := floor+1; f < NumFloors; f++{
		for o := 0; o < NumTypes; o++{
			if(orders[f][o]){
				return true
			}
		}
	}
	return false
}

func checkBelow(orders Orders, floor int) bool{
	for f := 0; f < floor; f++{
		for o := 0; o < NumTypes; o++{
			if(orders[f][o]){
				return true
			}
		}
	}
	return false
}