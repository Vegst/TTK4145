package main

import (
	"fmt"
	"time"
	"./broadcaster/udp"
)

func main() {

	// Make threads
	//sendCh := make(chan string)
	//receiveCh := make(chan string)
	//go broadcaster.Broadcaster(sendCh, receiveCh)

	//conn := udp.Init(20009)
	// Main loop
	//for {
	/*
		select {
			case msg := <-receiveCh:
				fmt.Println(msg)
			default:
			*/
				udp.Broadcast(20009, "Test")
				//sendCh <- "test"
				fmt.Println("sent")
				time.Sleep(time.Second/10)

		//}
	//}	
}
