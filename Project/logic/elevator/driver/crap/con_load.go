package driver // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.h and driver.go
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "con_load.h"
*/
import "C"
import (
	"fmt"
)

func ConLoad(file string) {
	C.con_load(C.string(file), )
}

func IO_set_bit(channel int) {
	C.io_set_bit(C.int(channel))
}

func IO_clear_bit(channel int) {
	C.io_clear_bit(C.int(channel))
}

func IO_read_bit(channel int) bool {
	return bool(int(C.io_read_bit(C.int(channel))) != 0)
}

func IO_write_analog(channel int, value int) {
	C.io_write_analog(C.int(channel), C.int(value))
}

func IO_read_analog(channel int) int {
	return int(C.io_read_analog(C.int(channel)))
}
