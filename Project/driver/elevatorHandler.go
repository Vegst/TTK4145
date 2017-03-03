package driver

import (
	"time"
)

type ButtonEvent struct {
	Floor  int
	Button elev_button_type_t
}

type LightType int

const (
	LIGHT_TYPE_UP      = 0
	LIGHT_TYPE_DOWN    = 1
	LIGHT_TYPE_COMMAND = 2
	LIGHT_TYPE_FLOOR   = 3
	LIGHT_TYPE_STOP    = 4
	LIGHT_TYPE_DOOR    = 5
)

type LightEvent struct {
	Floor     int
	LightType LightType
}

func Elevator_event_manager(buttonEventCh chan ButtonEvent, lightEventCh chan LightEvent, stopCh chan bool, motorStateCh chan Elev_motor_direction_t, floorCh chan int) {
	Elev_init()

	var lastButtonState = [N_FLOORS][N_BUTTONS]bool{
		{false, false, false},
		{false, false, false},
		{false, false, false},
		{false, false, false},
	}
	var buttonState bool

	var lastStopState bool
	var stopState bool

	var lastFloorState int
	var floorState int

	for {
		select {
		case ms := <-motorStateCh:
			Elev_set_motor_direction(ms)
		case l := <-lightEventCh:
			switch l.LightType {
			case LIGHT_TYPE_UP:
				Elev_set_button_lamp(BUTTON_CALL_UP, l.Floor, 1)
			case LIGHT_TYPE_DOWN:
				Elev_set_button_lamp(BUTTON_CALL_DOWN, l.Floor, 1)
			case LIGHT_TYPE_COMMAND:
				Elev_set_button_lamp(BUTTON_COMMAND, l.Floor, 1)
			case LIGHT_TYPE_FLOOR:
				Elev_set_floor_indicator(l.Floor)
			case LIGHT_TYPE_STOP:
				Elev_set_stop_lamp(1)
			case LIGHT_TYPE_DOOR:
				Elev_set_door_open_lamp(1)
			}
		default:
			stopState = Elev_get_stop_signal()
			if stopState != lastStopState {
				lastStopState = stopState
				stopCh <- stopState
			}

			floorState = Elev_get_floor_sensor_signal()
			if (floorState != lastFloorState) && (floorState != -1) {
				lastFloorState = floorState
				floorCh <- floorState
			}

			for f := 0; f < 4; f++ {
				buttonState = Elev_get_button_signal(BUTTON_COMMAND, f)
				if buttonState != lastButtonState[BUTTON_COMMAND][f] {
					lastButtonState[BUTTON_COMMAND][f] = buttonState
					buttonEventCh <- ButtonEvent{f, BUTTON_COMMAND}
				}
				if f > 0 {
					buttonState = Elev_get_button_signal(BUTTON_CALL_DOWN, f)
					if buttonState != lastButtonState[BUTTON_CALL_DOWN][f] {
						lastButtonState[BUTTON_CALL_DOWN][f] = buttonState
						buttonEventCh <- ButtonEvent{f, BUTTON_CALL_DOWN}
					}
				}
				if f < N_FLOORS-1 {
					buttonState = Elev_get_button_signal(BUTTON_CALL_UP, f)
					if buttonState != lastButtonState[BUTTON_CALL_UP][f] {
						lastButtonState[BUTTON_CALL_UP][f] = buttonState
						buttonEventCh <- ButtonEvent{f, BUTTON_CALL_UP}
					}
				}
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}
