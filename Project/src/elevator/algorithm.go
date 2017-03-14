package elevator

import (
	. "../def"
)

func IsOrderAtFloor(e Elevator) bool {
	for o := 0; o < NumTypes; o++{
		if(e.Orders[e.State.Floor][o]){
			return true
		}
	}
	return false
}

func GetDirection(e Elevator) MotorDirection{
	switch(e.State.Direction) {
	case DirnUp:
		if(checkAbove(e)) {
			return DirnUp
		} else if(checkBelow(e)) {
			return DirnDown
		} else{
			return DirnStop
		}
	case DirnStop, DirnDown:
		if(checkBelow(e)) {
			return DirnDown
		} else if(checkAbove(e)) {
			return DirnUp
		} else{
			return DirnStop
		}	
	}
	return DirnStop
}

func ShouldStop(e Elevator) bool{
	switch(e.State.Direction) {
	case DirnDown:
		return 	e.Orders[e.State.Floor][OrderCallCommand] ||
				e.Orders[e.State.Floor][OrderCallDown] ||
				!checkBelow(e)
	case DirnUp:
		return 	e.Orders[e.State.Floor][OrderCallCommand] ||
				e.Orders[e.State.Floor][OrderCallUp] ||
				!checkAbove(e) 	
	}
	return true
}
/*
func GetOrdersToClear(orders Orders, floor int, direction MotorDirection) {

	if sm.State.Direction == DirnUp {
		sm.OrdersEvents.Order <- Order{sm.State.Floor, OrderCallUp, false}
	} else if sm.State.Direction == DirnDown {
		sm.OrdersEvents.Order <- Order{sm.State.Floor, OrderCallDown, false}
	}
	sm.OrdersEvents.Order <- Order{sm.State.Floor, OrderCallCommand, false}
}*/

func checkAbove(e Elevator) bool{
	for f := e.State.Floor+1; f < NumFloors; f++{
		for o := 0; o < NumTypes; o++{
			if(e.Orders[f][o]){
				return true
			}
		}
	}
	return false
}

func checkBelow(e Elevator) bool{
	for f := 0; f < e.State.Floor; f++{
		for o := 0; o < NumTypes; o++{
			if(e.Orders[f][o]){
				return true
			}
		}
	}
	return false
}