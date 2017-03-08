package main

import (
	"./elevator"
	"fmt"
	"time"
)

func main() {

	fmt.Println("Start")

	go elevator.StateMachine()
	for {
		select {
		case <-time.After(10 * time.Millisecond):
		}
	}
}
