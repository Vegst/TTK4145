package client
 
import (
    "fmt"
    "net"
    "time"
)
 
type Event int
const (
	EVENT_CONNECTED = 1
	EVENT_DISCONNECTED = 2
)
func ToString(event Event) string {
	switch event {
	case EVENT_CONNECTED:
		return "Connected"
	case EVENT_DISCONNECTED:
		return "Disconnected"
	}
	return "Unknown"
}

type Message struct {
	Data string
}

func Client(port int, serverCh <-chan string, sendCh <-chan Message, receiveCh chan Message, eventCh chan Event) {

	var connection net.Conn
	var err error
	connected := false
	server := ""

	for {
		select {
		case server = <-serverCh:
		case message := <-sendCh:
			if connected {
				_,err := connection.Write([]byte(message.Data))
				if err != nil {
					fmt.Println(err)
				}
			}
		default:
			// Connect
			if server != "" {
				if !connected {
					connection, err = net.Dial("tcp", fmt.Sprintf("%s:%d", server, port))
					if err == nil {
						connected = true
						eventCh <- EVENT_CONNECTED
					}
				}
			// Disconnect
			} else if connected { 
				connection.Close()
				connected = false
				eventCh <- EVENT_DISCONNECTED
			}
			// Alive
			if connected {
				buffer := make([]byte, 1024)

				n, err := connection.Read(buffer)
				if err == nil {
					receiveCh <- Message{string(buffer[0:n])}
				} else {
					connection.Close()
					connected = false
					eventCh <- EVENT_DISCONNECTED
				}
			}
			time.Sleep(50*time.Millisecond)
		}
	}
}
