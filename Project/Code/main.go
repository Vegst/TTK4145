package main

import (
	. "./def"
	"./driver"
	"./elevator"
	//"./gui"
	"./network"
	//"./network/conn"
	"./network/localip"
	"./orders"
	"./gui"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {

	// Handle terminal input: id and simulation config file
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

	driverElevatorEvents := DriverElevatorEvents {
		Button: make(chan ButtonEvent, 10),
		Light: make(chan LightEvent, 1),
		Stop: make(chan bool, 10),
		MotorDirection: make(chan MotorDirection, 10),
		Floor: make(chan int, 10),
		DoorOpen: make(chan bool, 10),
		FloorIndicator: make(chan int, 10),
	}

	elevatorOrdersEvents := ElevatorOrdersEvents {
		Order: make(chan OrderEvent, 10),
		State: make(chan ElevatorState, 10),
		LocalOrders: make(chan Orders, 10),
		GlobalOrders: make(chan Orders, 10),
	}

	ordersNetworkEvents := OrdersNetworkEvents {
		TxAssignedOrder: make(chan AssignedOrder, 10),
		RxAssignedOrder: make(chan AssignedOrder, 10),
		TxAssignedState: make(chan AssignedState, 10),
		RxAssignedState: make(chan AssignedState, 10),
	}

	ordersGuiEvents := OrdersGuiEvents {
		Elevators: make(chan Elevators, 10),
	}

	go network.Network(id, ordersNetworkEvents)
	go elevator.StateMachine(driverElevatorEvents, elevatorOrdersEvents)
	go driver.EventManager(driverElevatorEvents)
	go orders.OrderManager(id, elevatorOrdersEvents, ordersNetworkEvents, ordersGuiEvents)
	go gui.ElevatorVisualizer(ordersGuiEvents)

	for {
		select {
		case <-time.After(5 * time.Second):
		}
	}
}
