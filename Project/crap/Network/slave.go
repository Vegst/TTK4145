package main

import (
	"./client"
	"fmt"
	"time"
)

func main() {
	fmt.Println("Start")

	serverCh := make(chan string)
	sendCh := make(chan client.Message)
	receiveCh := make(chan client.Message)
	eventCh := make(chan client.Event)

	go client.Client(20013, "test", serverCh, sendCh, receiveCh, eventCh)

	
	serverCh <- "127.0.0.1"
	
	for {
		select {
		case message := <-receiveCh:
			fmt.Println("Received: ", message)
		case event := <-eventCh:
			fmt.Println("Event: ", client.ToString(event))
		default:
			time.Sleep(50*time.Millisecond)
		}
	}
}
