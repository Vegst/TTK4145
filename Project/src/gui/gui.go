package gui

import (
	. "../def"
	"fmt"
	"os"
	"os/exec"
	"time"
    "sort"
)

func clearTerminal() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func sortElevators(elevators Elevators) []string {
    mk := make([]string, len(elevators))
    i := 0
    for k, _ := range elevators {
        mk[i] = k
        i++
    }
    sort.Strings(mk)
    return mk
}

func print(id string, elevators Elevators) {
	for _,e := range sortElevators(elevators) {
		if e == id {
			fmt.Println("ELEVATOR", e, "(this)")
		} else {
			fmt.Println("ELEVATOR", e)
		}
		fmt.Println("Floor	State	Up	Down	Command")
		for f := len(elevators[e].Orders)-1; f >= 0; f-- {
			fmt.Print(f, "	")
			if f == len(elevators[e].Orders)-1 && elevators[e].State.Floor < 0 {
				fmt.Print("U")
			} else if f == elevators[e].State.Floor {
				switch elevators[e].State.Direction {
				case DirnUp:
					fmt.Print("ˆ ")
				case DirnDown:
					fmt.Print("v ")
				case DirnStop:
					fmt.Print("  ")
				}
				switch elevators[e].State.Behaviour {
				case ElevatorBehaviourIdle:
					fmt.Print(" []")
				case ElevatorBehaviourMoving:
					fmt.Print(" []*")
				case ElevatorBehaviourDoorOpen:
					fmt.Print("[  ]")
				}
			}
			for t, order := range elevators[e].Orders[f] {
				fmt.Print("	")
				if f == len(elevators[e].Orders)-1 && t == int(OrderCallUp) {
					continue
				}
				if f == 0 && t == int(OrderCallDown) {
					continue
				}
				if order {
					fmt.Print("*")
				} else {
					fmt.Print("-")
				}
			}
			fmt.Println("")
		}
		fmt.Println("")
	}
}

func ElevatorVisualizer(id string, ordersEvents OrdersGuiEvents) {
	elevators := make(Elevators)
	clearTerminal()
	print(id, elevators)

	for {
		select {
		case elevators = <-ordersEvents.Elevators:
			clearTerminal()
			print(id, elevators)
		case <-time.After(100 * time.Millisecond):
		}
	}
}