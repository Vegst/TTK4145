package main

import (
	. "./def"
	"./driver"
	"./elevator"
	"./network"
	"./orders"
	"./gui"
	"./network/localip"
	. "./def"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {

	var id string
	var simulator string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.StringVar(&simulator, "sim", "", "simulator config file")
	flag.Parse()

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}
	if simulator == "" {
		simulator = "simulator1.con"
	}


	// Initialize system
	driver.Init(driver.TypeSimulation, simulator)

	// See documentation for full communication structure between main goroutines

	// Driver <--> StateMachine
	buttonEventCh := make(chan ButtonEvent, 10)
	lightEventCh := make(chan LightEvent, 10)
	stopCh := make(chan bool, 10)
	motorStateCh := make(chan MotorDirection, 10)
	floorCh := make(chan int, 10)
	doorOpenCh := make(chan bool, 10)
	floorIndicatorCh := make(chan int, 10)

	// StateMachine <--> OrderManager
	orderEventCh := make(chan OrderEvent, 10)
	stateCh := make(chan Elevator, 10)
	localOrdersCh := make(chan Orders, 10)
	globalOrdersCh := make(chan Orders, 10)

	// OrderManager <--> StateMachine

	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	helloTx := make(chan HelloMsg)
	helloRx := make(chan HelloMsg)

	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	go bcast.Transmitter(16569, helloTx)
	go bcast.Receiver(16569, helloRx)

	go network.Network(helloTx, helloRx, peerTxEnable, peerUpdateCh, localOrdersCh, orderEventCh)

	// OrderManager <--> GUI
	elevatorsCh := make(chan Elevators, 10)

	go elevator.StateMachine(buttonEventCh, lightEventCh, stopCh, motorStateCh, floorCh, doorOpenCh, floorIndicatorCh, orderEventCh, stateCh, localOrdersCh, globalOrdersCh)
	go driver.EventManager(buttonEventCh, lightEventCh, stopCh, motorStateCh, floorCh, doorOpenCh, floorIndicatorCh)
	go orders.OrderManager(id, orderEventCh, stateCh, localOrdersCh, globalOrdersCh, elevatorsCh)
	go gui.ElevatorVisualizer(elevatorsCh)

	for {
		select {
		case <-time.After(10 * time.Millisecond):
		}
	}
}
