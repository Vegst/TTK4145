package main

import (
	"fmt"
	"time"
	"./broadcaster/udp"
)

/*
func main() {

	// Make threads
	sendCh := make(chan string)
	receiveCh := make(chan string)
	go broadcaster.Broadcaster(sendCh, receiveCh)

	// Main loop
	for {
		select {
			case msg := <-receiveCh:
				fmt.Println(msg)
			default:
				time.Sleep(time.Second/10)

		}
	}
}
*/

func main() {

	conn := udp.Init(20009)

	for {
		msg := udp.Receive(conn)
		fmt.Println(msg)
		time.Sleep(time.Second/10)

	}
}
