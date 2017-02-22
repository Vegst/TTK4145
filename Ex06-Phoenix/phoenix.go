package main

import (
	"fmt"
	"time"
	"./broadcaster"
	"./timer"
)

func main() {
	master := false
	i := 0
	inactiveTime := time.Second * 0

	// Make threads
	sendCh := make(chan string)
	receiveCh := make(chan string)
	go broadcaster.Broadcaster(sendCh, receiveCh)

	timeoutCh := make(chan bool)
	go timer.Timer(time.Second, timeoutCh)

	// Main loop
	for {
		select {
			case <- timeoutCh:
				if master {
					i++
					fmt.Println(i)
				}
			case msg := <-receiveCh:
				if !master {
					inactiveTime = 0
				}
				fmt.Println(msg)
			default:
				if master {
					sendCh <- "test"
				}
				time.Sleep(time.Second/10)
				if !master {
					inactiveTime += time.Second/10
					if (inactiveTime < time.Second/2) {
						master = true
						inactiveTime = 0
						fmt.Println("Now master!")
					}
				}

		}
	}
	/*
	if master {
		i := 0
		for {
			i++
			fmt.Println(i)
			for c := 1; c <= 10; c++ {
				//udp.Broadcast(conn, "test")
				time.Sleep(time.Second/10)
			}
		}
	} else {
		for {
			fmt.Println("test")
			time.Sleep(time.Second)
		}
	}
	*/
	
}
