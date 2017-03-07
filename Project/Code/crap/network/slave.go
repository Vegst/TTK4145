package main

import (
	"./peer"
	"./localid"
	"fmt"
	"time"
	"flag"
)

func main() {
	// Get ID automatically or as passed argument from the command line 
	// using 'go run main.go -id=<any id>'
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()
	if id == "" {
		id = localid.ID()
	}

	broadcastCh := make(chan string)
	receiveCh := make(chan string)

	go peer.Peer(20013, broadcastCh, receiveCh)
	
	for {
		select {
		case data := <-receiveCh:
			fmt.Println("Received: ", data)
		case <-time.After(200*time.Millisecond):
			broadcastCh <- fmt.Sprintf("hei %s", id)
		}
	}
}
