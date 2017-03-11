package def

// EventManager

type ButtonEvent struct {
	Floor  int
	Button ButtonType
	State  bool
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
type OrderEvent struct {
	Floor int
	Type  OrderType
	Flag  bool
}

type NetOrder struct {
	OrderEvent OrderEvent
	ID         string
}

type Orders [NumFloors][NumTypes]bool

const (
	NumFloors = 4
	NumTypes  = 3
)

const NumElevators = 3

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
	ID     string
}
type ElevatorState struct {
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
