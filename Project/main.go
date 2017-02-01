package main

import (
	"./Network/server"
	"./Network/localip"
	"fmt"
)

func main() {
	fmt.Println("Start")

	enableCh := make(chan bool)
	sendCh := make(chan string)
	receiveCh := make(chan string)
	eventCh := make(chan bool)

	var id string

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	go server.Server(20013, id, enableCh, sendCh, receiveCh, eventCh)

	enableCh <- true
  sendCh <- "hei"

	for {
		select {
		case message := <-receiveCh:
			fmt.Println("Received: ", message)
		case event := <-eventCh:
			fmt.Println("Event: ", event)
		}
	}
}
