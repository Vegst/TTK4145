package elevator

import (
	"../driver"
)

type Orders [NumFloors][NumTypes] bool

const (
	NumFloors = driver.NumFloors
	NumTypes  = driver.NumButtons
)

type OrderType int

const (
	OrderCallUp      OrderType = 0
	OrderCallDown    OrderType = 1
	OrderCallCommand OrderType = 2
)