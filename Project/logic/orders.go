package logic

import (
	"./elevator"
)


const (
	NumFloors = elevator.NumFloors
	NumTypes  = elevator.NumButtons
)

type OrderType int

const (
	OrderCallUp      OrderType = 0
	OrderCallDown    OrderType = 1
	OrderCallCommand OrderType = 2
)