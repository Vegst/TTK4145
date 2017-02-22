package broadcaster

import (
	"net"
	"time"
	"./udp"
)

func Broadcaster(sendCh <-chan string, receiveCh chan string) {
	conn := udp.Init()
	go listener(conn, receiveCh)
	for {
		select {
			case msg := <- sendCh:
				udp.Broadcast(conn, msg)
			default: 
				time.Sleep(50*time.Millisecond)
		}
	}
}

func listener(conn *net.UDPConn, receiveCh chan string) {
	receiveCh <- udp.Receive(conn)
}