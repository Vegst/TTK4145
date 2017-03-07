package peer

import (
	"time"
	"./udp"
	"math/rand"
)

func Peer(port int, broadcastCh <-chan string, receiveCh chan string) {
	go Transmitter(port, broadcastCh)
	go Receiver(port, receiveCh)
}

func Transmitter(port int, broadcastCh <-chan string) {
	c := udp.Open(port)
	for {
		select {
		case msg := <- broadcastCh:
			udp.Broadcast(c, port, msg)
		case <-time.After(50*time.Millisecond):
		}
	}
	udp.Close(c)
}

func Receiver(port int, receiveCh chan string) {
	c := udp.Open(port)
	for {
		message, error := udp.Receive(c)
		if error == nil {
			if (rand.Intn(100) < 50) { // Simulate lost packets (50%)
				receiveCh <- message
			}
		}
	}
	udp.Close(c)
}
