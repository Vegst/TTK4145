package udp

import (
	"fmt"
	"net"
	"./conn"
)

type Connection struct {
	PacketConn net.PacketConn
	Addr *net.UDPAddr
}

func Init(port int) Connection {
	packetConn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))
	return Connection{packetConn, addr};
}

func Broadcast(c Connection, msg string) {
	c.PacketConn.WriteTo([]byte(msg), c.Addr)
}

func Receive(c Connection) (string, error) {
	buf := make([]byte, 1024)
	n, _, err := c.PacketConn.ReadFrom(buf[0:])
	if err != nil {
		return ""
	}

	return string(buf[0:n]), err
}

