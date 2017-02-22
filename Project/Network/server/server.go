package server

import (
	"fmt"
	"net"
	"time"
)

const (
	EVENT_TYPE_SERVER_OPENED = 1
	EVENT_TYPE_SERVER_CLOSED = 2
	EVENT_TYPE_CLIENT_CONNECTED = 3
	EVENT_TYPE_CLIENT_DISCONNECTED = 4
)

type Event struct {
	Type int
}

type Message struct {
	Client string
	Data string
}

var clients []net.Conn

func Server(port int, id string, enableCh <-chan bool, sendCh <-chan Message, receiveCh chan Message, eventCh chan Event) {

	var serverConnection net.Listener
	var err error
	listening := false
	enabled := false

	for {
		select {
		case enabled = <-enableCh:
		case message := <-sendCh:
			for _,client := range clients {
				if message.Client == "test" {
					_,err := client.Write([]byte(message.Data))
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		default:
			if enabled {
				// Listen for clients
				if listening { 
					connection, err := serverConnection.Accept()
					if err == nil {
						eventCh <- Event{EVENT_TYPE_CLIENT_CONNECTED}
						go handleConnection(connection, receiveCh, eventCh)
					} else {
						fmt.Println(err)
					}
					
				// Open connection
				} else {
					serverConnection, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
					if err == nil {
						eventCh <- Event{EVENT_TYPE_SERVER_OPENED}
						listening = true
					} else {
						fmt.Println(err)
					}
				}
			// Stop listening and disconnect
			} else { 
				serverConnection.Close()
				eventCh <- Event{EVENT_TYPE_SERVER_CLOSED}
			}
			time.Sleep(50*time.Millisecond)
		}
	}
}


func handleConnection(connection net.Conn, receiveCh chan Message, eventCh chan Event) {
	buffer := make([]byte, 1024)

	for {
		n, err := connection.Read(buffer)
		if err == nil {
			receiveCh <- Message{connection.RemoteAddr().String(), string(buffer[0:n])}
		} else {
			eventCh <- Event{EVENT_TYPE_CLIENT_DISCONNECTED}
			return
		}
	}
}
