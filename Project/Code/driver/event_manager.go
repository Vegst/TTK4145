package driver

import (
	"time"
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
	State  bool
}

type LightType int

const (
	LIGHT_TYPE_UP      = 0
	LIGHT_TYPE_DOWN    = 1
	LIGHT_TYPE_COMMAND = 2
	LIGHT_TYPE_STOP    = 4
)

type LightEvent struct {
	LightType LightType
	Floor     int
	Value     bool
}

func EventManager(buttonEventCh chan ButtonEvent, lightEventCh chan LightEvent, stopCh chan bool, motorStateCh chan MotorDirection, floorCh chan int, doorOpenCh chan bool, floorIndicatorCh chan int) {

	Init(TypeComedi)

	// Storage of last states to detect change of state
	var lastButtonState [NumFloors][NumButtons]bool
	for f := 0; f < NumFloors; f++ {
		for b := 0; b < NumButtons; b++ {
			lastButtonState[f][b] = GetButtonSignal(ButtonType(b), f)
		}
	}

	lastStopState := GetStopSignal()
	lastFloorState := GetFloorSignal()

	var buttonState bool
	var stopState bool
	var floorState int

	for {
		select {
		case ms := <-motorStateCh:
			SetMotorDirection(MotorDirection(ms))
		case l := <-lightEventCh:
			switch l.LightType {
			case LIGHT_TYPE_UP:
				SetButtonLamp(ButtonCallUp, l.Floor, l.Value)
			case LIGHT_TYPE_DOWN:
				SetButtonLamp(ButtonCallDown, l.Floor, l.Value)
			case LIGHT_TYPE_COMMAND:
				SetButtonLamp(ButtonCallCommand, l.Floor, l.Value)
			case LIGHT_TYPE_STOP:
				SetStopLamp(l.Value)
			}
		case doorOpen := <-doorOpenCh:
			SetDoorOpenLamp(doorOpen)
		case floorIndicator := <-floorIndicatorCh:
			SetFloorIndicator(floorIndicator)
		case <-time.After(10 * time.Millisecond):
			stopState = GetStopSignal()
			if stopState != lastStopState {
				lastStopState = stopState
				stopCh <- stopState
			}

			floorState = GetFloorSignal()
			if floorState != lastFloorState {
				lastFloorState = floorState
				floorCh <- floorState
			}

			for f := 0; f < NumFloors; f++ {
				for b := 0; b < NumButtons; b++ {
					if ButtonType(b) == ButtonCallUp && f == NumFloors-1 {
						continue
					}
					if ButtonType(b) == ButtonCallDown && f == 0 {
						continue
					}
					buttonState = GetButtonSignal(ButtonType(b), f)
					if buttonState != lastButtonState[f][b] {
						lastButtonState[f][b] = buttonState
						buttonEventCh <- ButtonEvent{f, ButtonType(b), buttonState}
					}
				}
			}
		}
	}
}
