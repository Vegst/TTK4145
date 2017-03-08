//

package elevator

import (
	"../driver"
)

func checkAbove(localOrders Orders, elev Elevator) bool {
	for f := elev.Floor + 1; f < NumFloors; f++ {
		for o := 0; o < NumTypes; o++ {
			if localOrders[f][o] {
				return true
			}
		}
	}
	return false
}

func checkBelow(localOrders Orders, elev Elevator) bool {
	for f := 0; f < elev.Floor; f++ {
		for o := 0; o < NumTypes; o++ {
			if localOrders[f][o] {
				return true
			}
		}
	}
	return false
}

func OrderAtFloor(localOrders Orders, floor int) bool {
	for o := 0; o < NumTypes; o++ {
		if localOrders[floor][o] {
			return true
		}
	}
	return false
}

func GetDirection(localOrders Orders, elev Elevator) driver.MotorDirection {
	switch elev.Direction {
	case driver.DirnUp:
		if checkAbove(localOrders, elev) {
			return driver.DirnUp
		} else if checkBelow(localOrders, elev) {
			return driver.DirnDown
		} else {
			return driver.DirnStop
		}
	case driver.DirnStop, driver.DirnDown:
		if checkBelow(localOrders, elev) {
			return driver.DirnDown
		} else if checkAbove(localOrders, elev) {
			return driver.DirnUp
		} else {
			return driver.DirnStop
		}
	}
	return driver.DirnStop
}

func ShouldStop(localOrders Orders, elev Elevator) bool {
	switch elev.Direction {
	case driver.DirnDown:
		return localOrders[elev.Floor][OrderCallCommand] ||
			localOrders[elev.Floor][OrderCallDown] ||
			!checkBelow(localOrders, elev)
	case driver.DirnUp:
		return localOrders[elev.Floor][OrderCallCommand] ||
			localOrders[elev.Floor][OrderCallUp] ||
			!checkAbove(localOrders, elev)
	}
	return true
}
