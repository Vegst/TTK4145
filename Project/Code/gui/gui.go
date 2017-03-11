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
	for id,elevator := range elevators {
		fmt.Println(id)
		fmt.Println("	Direction: ", elevator.Direction)
		fmt.Println("	Behaviour: ", elevator.Behaviour)
		for f, orders := range elevator.Orders {
			if f == elevator.Floor {
				switch(elevator.Direction) {
				case DirnUp:
					fmt.Print("^")

				}
			}
		    for _, order := range orders {
		        if order {
					fmt.Print("	*")
				} else {
					fmt.Print("	-")
				}
		    }
		    fmt.Println()
		}
	}
}

func ElevatorVisualizer(elevatorsCh chan Elevators) {
	elevators := make(Elevators)
	clearTerminal()
	print(elevators)

	for {
		select {
		case elevators = <-elevatorsCh:
			clearTerminal()
			print(elevators)
		case <-time.After(100 * time.Millisecond):
		}
	}
}