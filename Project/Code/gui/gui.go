package gui

import (
	. "../def"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func clearTerminal() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func print(elevators Elevators) {
	for id, elevator := range elevators {
		fmt.Println("ELEVATOR", id)
		fmt.Println("State		Up	Down	Command")
		for f := len(elevator.Orders)-1; f >= 0; f-- {
			if f == elevator.State.Floor {
				switch elevator.State.Direction {
				case DirnUp:
					fmt.Print("ˆ ")
				case DirnDown:
					fmt.Print("v ")
				case DirnStop:
					fmt.Print("  ")
				}
				switch elevator.State.Behaviour {
				case ElevatorBehaviourIdle:
					fmt.Print(" []")
				case ElevatorBehaviourMoving:
					fmt.Print(" []*")
				case ElevatorBehaviourDoorOpen:
					fmt.Print("[  ]")
				}
			}
			fmt.Print("	")
			for _, order := range elevator.Orders[f] {
				if order {
					fmt.Print("	*")
				} else {
					fmt.Print("	-")
				}
			}
			fmt.Println()
		}
	}
	fmt.Println()
}

func ElevatorVisualizer(ordersGuiEvents OrdersGuiEvents) {
	elevators := make(Elevators)
	clearTerminal()
	print(elevators)

	for {
		select {
		case elevators = <-ordersGuiEvents.Elevators:
			clearTerminal()
			print(elevators)
		case <-time.After(100 * time.Millisecond):
		}
	}
}
