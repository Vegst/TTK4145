package server

import (
	"fmt"
	"net"
	"time"
)

type Event struct {
	Type int
}
const (
	EVENT_TYPE_SERVER_OPENED = 3
	EVENT_TYPE_SERVER_CLOSED = 4
	EVENT_TYPE_CLIENT_CONNECTED = 10
	EVENT_TYPE_CLIENT_DISCONNECTED = 11
)
func ToString(event Event) string {
	switch event.Type {
	case EVENT_TYPE_SERVER_OPENED:
		return "Server Opened"
	case EVENT_TYPE_SERVER_CLOSED:
		return "Server Closed"
	case EVENT_TYPE_CLIENT_CONNECTED:
		return "Client Connected"
	case EVENT_TYPE_CLIENT_DISCONNECTED:
		return "Client Disconnected"
	}
	return "Unknown"
}

type Message struct {
	Client *net.Conn
	Data string
}

var clients []net.Conn

func Server(port int, enableCh <-chan bool, sendCh <-chan Message, receiveCh chan Message, eventCh chan Event) {

	var serverConnection net.Listener
	var err error
	enabled := false
	opened := false

	for {
		select {
		case enabled = <-enableCh:
		case message := <-sendCh:
			for _,client := range clients {
				if message.Client == &client {
					client.Write([]byte(message.Data))
					break
				}
			}
		default:
			if enabled {
				// Listen for clients
				if opened { 
					connection, err := serverConnection.Accept()
					if err == nil {
						clients = append(clients, connection)
						go handleConnection(&clients[len(clients)-1], receiveCh, eventCh)
						eventCh <- Event{EVENT_TYPE_CLIENT_CONNECTED}
					} else {
						opened = false
						eventCh <- Event{EVENT_TYPE_SERVER_CLOSED}
					}
					
				// Open connection
				} else {
					serverConnection, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
					if err == nil {
						opened = true
						eventCh <- Event{EVENT_TYPE_SERVER_OPENED}
					}
				}
			// Close connection
			} else if opened { 
				serverConnection.Close()
				eventCh <- Event{EVENT_TYPE_SERVER_CLOSED}
				opened = false
			}
			time.Sleep(50*time.Millisecond)
		}
	}
}


func handleConnection(connection *net.Conn, receiveCh chan Message, eventCh chan Event) {
	buffer := make([]byte, 1024)

	for {
		n, err := (*connection).Read(buffer)
		if err == nil {
			receiveCh <- Message{connection, string(buffer[0:n])}
		} else {
			eventCh <- Event{EVENT_TYPE_CLIENT_DISCONNECTED}
			return
		}
	}
}
