package driver

import (
	"time"
	"../def"
)

func EventManager(buttonEventCh chan def.ButtonEvent, lightEventCh chan def.LightEvent, stopCh chan bool, motorStateCh chan def.MotorDirection, floorCh chan int, doorOpenCh chan bool, floorIndicatorCh chan int) {


	// Storage of last states to detect change of state
	var lastButtonState [NumFloors][NumButtons]bool
	for f := 0; f < NumFloors; f++ {
		for b := 0; b < NumButtons; b++ {
			lastButtonState[f][b] = GetButtonSignal(def.ButtonType(b), f)
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
			SetMotorDirection(def.MotorDirection(ms))
		case l := <-lightEventCh:
			switch l.LightType {
			case def.LightTypeUp:
				SetButtonLamp(def.ButtonCallUp, l.Floor, l.Value)
			case def.LightTypeDown:
				SetButtonLamp(def.ButtonCallDown, l.Floor, l.Value)
			case def.LightTypeCommand:
				SetButtonLamp(def.ButtonCallCommand, l.Floor, l.Value)
			case def.LightTypeStop:
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
					if def.ButtonType(b) == def.ButtonCallUp && f == NumFloors-1 {
						continue
					}
					if def.ButtonType(b) == def.ButtonCallDown && f == 0 {
						continue
					}
					buttonState = GetButtonSignal(def.ButtonType(b), f)
					if buttonState != lastButtonState[f][b] {
						lastButtonState[f][b] = buttonState
						buttonEventCh <- def.ButtonEvent{f, def.ButtonType(b), buttonState}
					}
				}
			}
		}
	}
}
