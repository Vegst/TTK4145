package main

import (
	"fmt"
	"time"
	"./elevator"
)

func main() {

	fmt.Println("Start")
	go elevator.StateMachine()
	for {
		select {
		case <-time.After(10*time.Millisecond):
		}
	}
}
