package driver

import (
	"time"
	"fmt"
)

type Order struct{
	floor int
	dir int
	button elev_button_type_t
}

type Light struct{
	floor int
	flag bool
	button elev_button_type_t
}

func OrderHandler_handle_orders()

func OrderHandler_get_orders(orderInternal chan Order){
	time.Sleep(10*time.Millisecond)
	var newOrder Order
	var orderList [N_BUTTONS][N_FLOORS] bool
	
	for {
		for f := 0; f < 4; f++{
			if(Elev_get_button_signal(BUTTON_COMMAND, f)){
				Elev_set_button_lamp(BUTTON_COMMAND, f, 1)
				if(!orderList[BUTTON_COMMAND][f]){
					newOrder.floor = f
					newOrder.button = BUTTON_COMMAND
					orderList[BUTTON_COMMAND][f] = true
					fmt.Println("Floor: ", newOrder.floor)
					fmt.Println("Button: ", newOrder.button)
					orderInternal <- newOrder
					
				}
			}
			if(f > 0){
				if(Elev_get_button_signal(BUTTON_CALL_DOWN, f)){
					Elev_set_button_lamp(BUTTON_CALL_DOWN, f, 1)
					if(!orderList[BUTTON_CALL_DOWN][f]){
						newOrder.floor = f
						newOrder.button = BUTTON_CALL_DOWN
						orderList[BUTTON_CALL_DOWN][f] = true
						fmt.Println("Floor: ", newOrder.floor)
						fmt.Println("Button: ", newOrder.button)
						orderInternal <- newOrder
					}
				}
			}
			if(f < N_FLOORS - 1){
				if(Elev_get_button_signal(BUTTON_CALL_UP, f)){
					Elev_set_button_lamp(BUTTON_CALL_UP, f, 1)
					if(!orderList[BUTTON_CALL_UP][f]){
						newOrder.floor = f
						newOrder.button = BUTTON_CALL_UP
						orderList[BUTTON_CALL_UP][f] = true
						fmt.Println("Floor: ", newOrder.floor)
						fmt.Println("Button: ", newOrder.button)
						orderInternal <- newOrder
					}
				}
			}

		}
	}
}
