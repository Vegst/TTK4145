package main

import (
	. "./def"
	"./driver"
	"./elevator"
	//"./gui"
	"./network"
	"./network/bcast"
	//"./network/conn"
	"./network/localip"
	"./network/peers"
	"./orders"
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
	stateCh := make(chan ElevatorState, 10)
	localOrdersCh := make(chan Orders, 10)
	globalOrdersCh := make(chan Orders, 10)

	// OrderManager <--> Network
	assignmentCh := make(chan Assignment, 10)
	assignedOrderCh := make(chan OrderEvent, 10)

	txAssignmentCh := make(chan Assignment, 10)
	rxAssignmentCh := make(chan Assignment, 10)
	txStateCh := make(chan ElevatorState, 10)
	rxStateCh := make(chan ElevatorState, 10)

	updateElevatorCh := make(chan Elevator, 10)
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	go bcast.Transmitter(16569, txAssignmentCh)
	go bcast.Receiver(16569, rxAssignmentCh)

	go network.Network(id, txAssignmentCh, rxAssignmentCh, assignmentCh, assignedOrderCh, txStateCh, rxStateCh,  peerTxEnable, peerUpdateCh, stateCh, updateElevatorCh)
	// OrderManager <--> GUI
	elevatorCh := make(chan Elevator, 10)
	elevatorsCh := make(chan Elevators, 10)

	go elevator.StateMachine(buttonEventCh, lightEventCh, stopCh, motorStateCh, floorCh, doorOpenCh, floorIndicatorCh, orderEventCh, stateCh, localOrdersCh, globalOrdersCh)
	go driver.EventManager(buttonEventCh, lightEventCh, stopCh, motorStateCh, floorCh, doorOpenCh, floorIndicatorCh)
	go orders.OrderManager(id, orderEventCh, assignedOrderCh, assignmentCh, stateCh, updateElevatorCh, localOrdersCh, globalOrdersCh, elevatorCh, elevatorsCh)

	//go gui.ElevatorVisualizer(elevatorsCh)

	for {
		select {
		case <-time.After(10 * time.Millisecond):
		}
	}
}
