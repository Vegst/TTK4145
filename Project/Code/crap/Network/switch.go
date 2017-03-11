package main

import (
	"./server"
	"./client"
	//"./Network/localip"
	"fmt"
	"time"
)

func main() {
	port := 20013

	serverEnableCh := make(chan bool)
	serverSendCh := make(chan server.Message)
	serverReceiveCh := make(chan server.Message)
	serverEventCh := make(chan server.Event)

	clientServerCh := make(chan string)
	clientSendCh := make(chan client.Message)
	clientReceiveCh := make(chan client.Message)
	clientEventCh := make(chan client.Event)


	go server.Server(port, serverEnableCh, serverSendCh, serverReceiveCh, serverEventCh)
	go client.Client(port, clientServerCh, clientSendCh, clientReceiveCh, clientEventCh)

	
	//serverEnableCh <- true
	clientServerCh <- "127.0.0.1"
	
	for {
		select {
		case message := <-serverReceiveCh:
			fmt.Println("Received: ", message)
		case event := <-serverEventCh:
			fmt.Println("Event: ", server.ToString(event))
		case message := <-clientReceiveCh:
			fmt.Println("Received: ", message)
		case event := <-clientEventCh:
			fmt.Println("Event: ", client.ToString(event))
		default:
			time.Sleep(50*time.Millisecond)
		}
	}
}
