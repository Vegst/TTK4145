package main

import (
	"./Network/server"
	//"./Network/localip"
	"fmt"
	"time"
)

func main() {
	fmt.Println("Start")

	enableCh := make(chan bool)
	sendCh := make(chan server.Message)
	receiveCh := make(chan server.Message)
	eventCh := make(chan server.Event)


	var id string
	id = "test";
	/*
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}
	*/

	go server.Server(20013, id, enableCh, sendCh, receiveCh, eventCh)

	
	enableCh <- true
	
	for {
		select {
		case message := <-receiveCh:
			fmt.Println("Received: ", message)
		case event := <-eventCh:
			fmt.Println("Event: ", event.Type)
		default:
			//sendCh <- "hei"
			time.Sleep(50*time.Millisecond)
		}
	}
}
