package driver

/*
#cgo CFLAGS: -std=gnu11
#cgo LDFLAGS: -lcomedi -lm
#include "elev.h"
*/
import (
	"C"
)

type MotorDirection int

const (
	DirnDown MotorDirection = -1
	DirnStop MotorDirection = 0
	DirnUp   MotorDirection = 1
)

type ButtonType int

const (
	ButtonCallUp      ButtonType = 0
	ButtonCallDown    ButtonType = 1
	ButtonCallCommand ButtonType = 2
)

type Type int

const (
	TypeComedi     Type = 0
	TypeSimulation Type = 1
)

const (
	NumFloors  = int(C.N_FLOORS)
	NumButtons = int(C.N_BUTTONS)
)

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

func Init(t Type) {
	C.elev_init(C.elev_type(t))
}

func SetMotorDirection(dirn MotorDirection) {
	C.elev_set_motor_direction(C.elev_motor_direction_t(dirn))
}

func SetButtonLamp(button ButtonType, floor int, value bool) {
	C.elev_set_button_lamp(C.elev_button_type_t(button), C.int(floor), C.int(b2i(value)))
}

func SetFloorIndicator(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
}

func SetDoorOpenLamp(value bool) {
	C.elev_set_door_open_lamp(C.int(b2i(value)))
}

func SetStopLamp(value bool) {
	C.elev_set_stop_lamp(C.int(b2i(value)))
}

func GetButtonSignal(button ButtonType, floor int) bool {
	return int(C.elev_get_button_signal(C.elev_button_type_t(button), C.int(floor))) != 0
}

func GetFloorSignal() int {
	return int(C.elev_get_floor_sensor_signal())
}

func GetStopSignal() bool {
	return int(C.elev_get_stop_signal()) != 0
}

func GetObstructionSignal() bool {
	return int(C.elev_get_obstruction_signal()) != 0
}
