package udp

import (
	"fmt"
	"net"
	"./conn"
)


func Open(port int) net.PacketConn {
	return conn.DialBroadcastUDP(port)
}

func Close(c net.PacketConn) {
	c.Close()
}

func Broadcast(c net.PacketConn, port int, message string) {
	address, error := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))
	if error == nil {
		c.WriteTo([]byte(message), address)
	}
}

func Receive(c net.PacketConn) (string, error) {
	buffer := make([]byte, 1024)
	n, _, error := c.ReadFrom(buffer[0:])
	if error != nil {
		return "", error
	}
	return string(buffer[0:n]), error
}