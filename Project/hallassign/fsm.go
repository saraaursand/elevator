package hallassign

import (
	"Elevator/driver-go-master/elevio"
	network "Elevator/networkcom"
	"Elevator/utils"
	"fmt"
	"time"
)

func FSM(helloTx chan network.HelloMsg, helloRx chan network.HelloMsg, drv_buttons chan elevio.ButtonEvent, drv_floors chan int, drv_obstr chan bool, drv_stop chan bool) {
	var d elevio.MotorDirection = elevio.MD_Up
	for {
		select {
		case E := <-drv_buttons:
			println("Got button press")
			utils.Elevator_glob.Requests[E.Floor][E.Button] = true
			fmt.Println("Updatet: ", utils.Elevator_glob.Requests)

			e := utils.Elevator_glob
			helloMsg := network.HelloMsg{e, 0}
			helloMsg.Iter++
			helloMsg.Elevator = utils.Elevator_glob
			helloTx <- helloMsg

		case F := <-drv_floors:
			println("Arrived at floor", F)
			utils.FsmOnFloorArrival(F)
			// OneElevRequests[F] = [utils.N_BUTTONS]bool{false, false, false}

			fmt.Println("OneElevRequests: ", OneElevRequests[F][0], OneElevRequests[F][1], OneElevRequests[F][2])
		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(d)
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < utils.N_FLOORS; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
		case <-time.After(time.Millisecond * time.Duration(utils.DoorOpenDuration*1000)):
			utils.FsmOnDoorTimeout()
		case elev := <-helloRx:
			flag := 0
			for i, element := range network.ListOfElevators {
				if element.ID == elev.Elevator.ID {
					network.ListOfElevators[i] = elev.Elevator
					flag = 1
				}
			}
			if flag == 0 {
				network.ListOfElevators = append(network.ListOfElevators, elev.Elevator)

			}
			AssignHallRequest()
		}
		// AssignHallRequest()
		// utils.Elevator_glob.Behaviour = utils.EB_Idle
		for floor_num, floor := range utils.Elevator_glob.Requests {
			for btn_num, _ := range floor {
				if utils.Elevator_glob.Requests[floor_num][btn_num] {
					utils.FsmOnRequestButtonPress(floor_num, utils.Button(btn_num))
				}
			}
		}
	}
}
