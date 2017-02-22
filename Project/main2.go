package main

import(
	"./driver"
)

func main() {
	driver.Elev_init()
	driver.Elev_set_motor_direction(driver.DIRN_UP)

	for {
        if (driver.Elev_get_floor_sensor_signal() == driver.N_FLOORS - 1) {
            driver.Elev_set_motor_direction(driver.DIRN_DOWN)
        } else if (driver.Elev_get_floor_sensor_signal() == 0) {
            driver.Elev_set_motor_direction(driver.DIRN_UP)
        }
        c := make(chan driver.Order, 1)
        go driver.OrderHandler_get_orders(c)


        // Stop elevator and exit program if the stop button is pressed
        if (driver.Elev_get_stop_signal()) {
            driver.Elev_set_motor_direction(driver.DIRN_STOP)
            return
        }
    }
}