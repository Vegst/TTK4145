package def


// EventManager

type ButtonEvent struct {
	Button Button
	State  bool
}


type Button struct {
	Floor int
	Type ButtonType
}

type LightType int

const (
	LightTypeUp      = 0
	LightTypeDown    = 1
	LightTypeCommand = 2
	LightTypeStop    = 3
)

type LightEvent struct {
	LightType LightType
	Floor     int
	Value     bool
}

// Orders
type Order struct {
	Floor int
	Type  OrderType
	Flag  bool
}

type OrderEvent struct {
	Target string
	Order Order
}
type StateEvent struct {
	Target string
	State ElevatorState
}

type Orders [NumFloors][NumTypes]bool

const (
	NumFloors = 4
	NumTypes  = 3
)

//const NumElevators = 3

type OrderType int

const (
	OrderCallUp      OrderType = 0
	OrderCallDown    OrderType = 1
	OrderCallCommand OrderType = 2
)

// Elevator
type ElevatorBehaviour int

const (
	ElevatorBehaviourIdle     = 0
	ElevatorBehaviourMoving   = 1
	ElevatorBehaviourDoorOpen = 2
)

type Elevators map[string]Elevator

type Elevator struct {
	State  ElevatorState
	Orders Orders
}
type ElevatorState struct {
	Active	  bool
	Floor     int
	Direction MotorDirection
	Behaviour ElevatorBehaviour
}

// Driver
type MotorDirection int

const (
	DirnDown MotorDirection = -1
	DirnStop MotorDirection = 0
	DirnUp   MotorDirection = 1
)

type ButtonType int

const (
	ButtonCallUp      ButtonType = 0
	ButtonCallDown    ButtonType = 1
	ButtonCallCommand ButtonType = 2
)


type DriverElevatorEvents struct {
	Button chan ButtonEvent
	Light chan LightEvent
	Stop chan bool
	MotorDirection chan MotorDirection
	Floor chan int
	DoorOpen chan bool
	FloorIndicator chan int
}

type ElevatorOrdersEvents struct {
	Order chan Order
	State chan ElevatorState
	LocalOrders chan Orders
	GlobalOrders chan Orders
}

type OrdersNetworkEvents struct {
	TxOrderEvent chan OrderEvent
	RxOrderEvent chan OrderEvent
	TxStateEvent chan StateEvent
	RxStateEvent chan StateEvent
	ElevatorNew chan string
	ElevatorLost chan string
	Elevators chan Elevators
}

type OrdersGuiEvents struct {
	Elevators chan Elevators
}


type MessageElevator struct {
	Id string
	Elevator Elevator
}
