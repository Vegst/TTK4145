package driver

import (
	"C"
	"fmt"
)

const (
	MOTOR_SPEED = 2800
	N_FLOORS = 4
	N_BUTTONS = 3
)

type elev_motor_direction_t int
const (
	DIRN_DOWN elev_motor_direction_t 	= -1
	DIRN_STOP elev_motor_direction_t 	= 0
    DIRN_UP elev_motor_direction_t 		= 1
)

type elev_button_type_t int
const (
	BUTTON_CALL_UP elev_button_type_t 	= 0
	BUTTON_CALL_DOWN elev_button_type_t = 1
    BUTTON_COMMAND elev_button_type_t 	= 2
)

var lamp_channel_matrix = [N_FLOORS][N_BUTTONS] int{
    {LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
    {LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
    {LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
    {LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var button_channel_matrix = [N_FLOORS][N_BUTTONS] int{
    {BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
    {BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
    {BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
    {BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

func Elev_init() int{
	if(!IO_init()){
		fmt.Println("Unable to initialize elevator hardware!")
		return -1
	}

	for f:= 0; f<N_FLOORS; f++{
		if f != 0 {
			Elev_set_button_lamp(BUTTON_CALL_DOWN, f, 0)
		}
		if f != N_FLOORS-1 {
			Elev_set_button_lamp(BUTTON_CALL_UP, f, 0)
		}
		Elev_set_button_lamp(BUTTON_COMMAND, f, 0)
	}
	Elev_set_stop_lamp(0)
    Elev_set_door_open_lamp(0)
    Elev_set_floor_indicator(0)

    return 0
}

func Elev_set_motor_direction(dirn elev_motor_direction_t){
	if(dirn==DIRN_STOP){
		IO_write_analog(MOTOR, 0)
	} else if(dirn==DIRN_UP){
		IO_clear_bit(MOTORDIR)
		IO_write_analog(MOTOR, MOTOR_SPEED)
	} else if(dirn==DIRN_DOWN){
		IO_set_bit(MOTORDIR)
		IO_write_analog(MOTOR, MOTOR_SPEED)
	}
}

func Elev_set_button_lamp(button elev_button_type_t, floor int, value int){
	if value == 1 {
		IO_set_bit(lamp_channel_matrix[floor][int(button)])
	} else{
		IO_clear_bit(lamp_channel_matrix[floor][int(button)])
	}
}

func Elev_set_floor_indicator(floor int){
	switch floor{
	case 0:
		IO_clear_bit(LIGHT_FLOOR_IND1)
		IO_clear_bit(LIGHT_FLOOR_IND2)
	case 1:
		IO_clear_bit(LIGHT_FLOOR_IND1)
		IO_set_bit(LIGHT_FLOOR_IND2)
	case 2:
		IO_set_bit(LIGHT_FLOOR_IND1)
		IO_clear_bit(LIGHT_FLOOR_IND2)
	case 3:
		IO_set_bit(LIGHT_FLOOR_IND1)
		IO_set_bit(LIGHT_FLOOR_IND2)
	}
}

func Elev_set_stop_lamp(value int){
	if value == 1{
		IO_set_bit(LIGHT_STOP)
	} else {
		IO_clear_bit(LIGHT_STOP)
	}
}

func Elev_set_door_open_lamp(value int) {
    if value == 1 {
        IO_set_bit(LIGHT_DOOR_OPEN);
    } else {
        IO_clear_bit(LIGHT_DOOR_OPEN);
    }
}

func Elev_get_button_signal(button elev_button_type_t, floor int) bool{
	return IO_read_bit(button_channel_matrix[floor][int(button)])
}

func Elev_get_floor_sensor_signal() int{
	 if (IO_read_bit(SENSOR_FLOOR1)) {
        return 0
    } else if (IO_read_bit(SENSOR_FLOOR2)) {
        return 1
    } else if (IO_read_bit(SENSOR_FLOOR3)) {
        return 2
    } else if (IO_read_bit(SENSOR_FLOOR4)) {
        return 3
    } else {
        return -1
    }
}

func Elev_get_stop_signal() bool{
    return IO_read_bit(STOP)
}

func Elev_get_obstruction_signal() bool{
    return IO_read_bit(OBSTRUCTION)
}
