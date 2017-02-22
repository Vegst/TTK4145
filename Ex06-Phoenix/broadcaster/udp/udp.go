package udp

import (
	"fmt"
	"net"
	"./conn"
)

func checkError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
    }
}

func Init(port int) net.PacketConn {
/*
	serverAddr,err := net.ResolveUDPAddr("udp","127.0.0.1:20009")
	checkError(err)

	localAddr,err := net.ResolveUDPAddr("udp","127.0.0.1:0")
	checkError(err)
	 */
	c := conn.DialBroadcastUDP(port) //"udp", localAddr, serverAddr
	return c
}

func Broadcast(port int, msg string) {
	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))

	conn.WriteTo([]byte(msg), addr)
}

func Receive(c net.PacketConn) string {
	buf := make([]byte, 1024)
	n, _, err := c.ReadFrom(buf[0:])
	checkError(err)

	return string(buf[0:n])
}

