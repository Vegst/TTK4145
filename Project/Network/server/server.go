package server

import (
	"fmt"
	"net"
)

func Server(port int, id string, enableCh <-chan bool, sendCh <-chan string, receiveCh chan string, eventCh chan bool) {

	var serverConnection net.Listener
	var err error
	listening := false
	enabled := false

	for {
		select {
		case enabled = <-enableCh:
		default:
			if enabled && !listening { 		// Start listening
				serverConnection, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
				if err == nil {
					listening = true
					defer serverConnection.Close()
				} else {
					fmt.Println(err)
				}
			} else if !enabled { // Stop listening and disconnect
				serverConnection.Close()
			} else if enabled && listening { // Open connection to client
				connection, err := serverConnection.Accept()
				if err == nil {
					go handleConnection(connection, eventCh)
					defer connection.Close()
				}
			}
		}
	}
}

func handleConnection(connection net.Conn, eventCh chan bool) {
	buffer := make([]byte, 1024)

	for {
		_, error := connection.Read(buffer)
		if error != nil {
			return
		}
	}
}
