package main

import (
	"fmt"
	"time"
	"./logic"
)

func main() {

	fmt.Println("Start")
	go logic.StateMachine()
	for {
		select {
		case <-time.After(10*time.Millisecond):
		}
	}
}
