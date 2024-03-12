package hallassign

import (
	network "Elevator/networkcom"
	"Elevator/utils"
)

var OneElevRequests = [utils.N_FLOORS][utils.N_BUTTONS]bool{}

func AssignHallRequest() {
	ListOfElevators := network.ListOfElevators
	AssignedHallCalls := CalculateCostFunc(ListOfElevators)
	OneElevCabCalls := GetCabCalls(utils.Elevator_glob)
	OneElevHallCalls := AssignedHallCalls[utils.Elevator_glob.ID]

	for floor := 0; floor < utils.N_FLOORS; floor++ {
		OneElevRequests[floor][0] = OneElevHallCalls[floor][0]
		OneElevRequests[floor][1] = OneElevHallCalls[floor][1]
		OneElevRequests[floor][2] = OneElevCabCalls[floor]
		// utils.Elevator_glob.Requests[floor] = OneElevRequests[floor]
	}

	for i := 0; i < len(OneElevHallCalls); i++ {
		for j := 0; j < len(OneElevHallCalls[i]); j++ {
			if OneElevRequests[i][j] {
				utils.Elevator_glob.Requests[i][j] = true
			} else {
				utils.Elevator_glob.Requests[i][j] = false
			}
		}
	}

	// utils.Elevator_glob.Requests[0] = OneElevRequests[0]
	// utils.Elevator_glob.Requests[1] = OneElevRequests[1]
	// utils.Elevator_glob.Requests[2] = OneElevRequests[2]
	// utils.Elevator_glob.Requests[3] = OneElevRequests[3]
}
