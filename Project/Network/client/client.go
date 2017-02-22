package slave
 
import (
    "fmt"
    "net"
    "os"
)
 
type Message struct {
	Data string
}

func Client(ip string, port int, id string, enableCh <-chan bool, sendCh <-chan Message, receiveCh chan Message, eventCh chan Event) {

	var connection net.Listener
	var err error
	connected := false
	enabled := false

	for {
		select {
		case enabled = <-enableCh:
		case message := <-sendCh:
			if connected {
				_,err := client.Write([]byte(message.Data))
				if err != nil {
					fmt.Println(err)
				}
			}
		default:
			// Connect
			if enabled && ! connected {
				connection, err = net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
				if err == nil {
					connected = true
				}
			// Disconnect
			} else { 
				connection.Close()
				connected = false
			}
			time.Sleep(50*time.Millisecond)
		}
	}
}




