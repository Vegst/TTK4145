package peer

import (
	"time"
	"./udp"
)

func Peer(port int, broadcastCh <-chan string, receiveCh chan string) {
	go Transmitter(port, broadcastCh)
	go Receiver(port, receiveCh)
}

func Transmitter(port int, broadcastCh <-chan string) {
	c := udp.Init(port)
	for {
		select {
			case msg := <- broadcastCh:
				udp.Broadcast(c, msg)
			default: 
				time.Sleep(50*time.Millisecond)
		}
	}
}

func Receiver(port int, receiveCh chan string) {
	c := udp.Init(port)
	for {
		receiveCh <- udp.Receive(c)
		time.Sleep(50*time.Millisecond)
	}
}
