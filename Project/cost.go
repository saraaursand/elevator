package main

import (
	//"Elevator/driver-go-master/elevio"
	"Elevator/utils"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

type HRAElevState struct {
	ElevID 		int         `json:"id"`
    Behaviour   string      `json:"behaviour"`
    Floor       int         `json:"floor"` 
    Direction   string      `json:"direction"`
    CabRequests []bool      `json:"cabRequests"`
}

type HRAElevStatetemp struct {
	ElevID 		int         `json:"id"`
    Behaviour   utils.ElevatorBehaviour      `json:"behaviour"`
    Floor       int         `json:"floor"` 
    Direction   utils.Dirn      `json:"direction"`
    CabRequests []bool      `json:"cabRequests"`
}

type HRAInput struct {
    HallRequests    [][2]bool                   `json:"hallRequests"`
    States          map[string]HRAElevState     `json:"states"`
}

func GetHallCalls(elevators []utils.Elevator) [][2]bool{	
	var n_elevators int = len(elevators)
	GlobalHallCalls := [][2]bool{}
	for floor := 0; floor < utils.N_FLOORS; floor++{
		HallCalls := [2]bool{}
		up := false
		down := false
		for i := 0; i < n_elevators; i++ {
			if elevators[i].Requests[floor][0] == true{
				up = true
			} 
			if elevators[i].Requests[floor][1] {
				down = true
			}
		}
		HallCalls[0] = up
		HallCalls[1] = down
		GlobalHallCalls[floor] = HallCalls
	}
	return GlobalHallCalls
}

func GetMyStates(elevators []utils.Elevator) []HRAElevStatetemp{
	var n_elevators int = len(elevators)
	myStates := []HRAElevStatetemp{}
	for i := 0; i < n_elevators; i++ {
		CabCalls := []bool{} 
		for floor := 0; floor < utils.N_FLOORS; floor++{
			CabCalls[floor] = elevators[i].Requests[floor][2]
		}
		elevastate := HRAElevStatetemp{i, elevators[i].Behaviour, elevators[i].Floor, elevators[i].Dirn, CabCalls}
		myStates[i] = elevastate
	}
	return myStates

} 

func CalculateCostFunc(elevators []utils.Elevator) *map[string][][2]bool{

    hraExecutable := ""
    switch runtime.GOOS {
        case "linux":   hraExecutable  = "hall_request_assigner"
        case "windows": hraExecutable  = "hall_request_assigner.exe"
        default:        panic("OS not supported")
    }	
	
	input := HRAInput{
		HallRequests: GetHallCalls(elevators),
		States: make(map[string]HRAElevState),
	}

	for _, elevatorStatus := range GetMyStates(elevators) {
		input.States[strconv.Itoa(elevatorStatus.ElevID)] = HRAElevState{
			Behaviour : func() string {
				if elevatorStatus.Behaviour == 0 {
					return "idle"
				} else if elevatorStatus.Behaviour == 1 {
					return "door open"
				} else {
					return "moving"
				}
			}(),	
			Floor : elevatorStatus.Floor,
			Direction : func() string {
				if elevatorStatus.Direction == -1 {
					return "down"
				} else if elevatorStatus.Direction == 0 {
					return "stop"
				} else {
					return "up"
				}
			}(),
			CabRequests : elevatorStatus.CabRequests,	
		}
	}



	//Convert input to json format
    jsonBytes, err := json.Marshal(input)
    if err != nil {
        fmt.Println("json.Marshal error: ", err)
        //return
    }
    
	//runds the hall_request_assigner file
    ret, err := exec.Command(hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
    if err != nil {
        fmt.Println("exec.Command error: ", err)
        fmt.Println(string(ret))
        //return
    }
    
	//convert the json received from hall_request_assigner to output
    output := new(map[string][][2]bool)
    err = json.Unmarshal(ret, &output)
    if err != nil {
        fmt.Println("json.Unmarshal error: ", err)
        //return
    }
    
	
    fmt.Printf("output: \n")
    for k, v := range *output {
        fmt.Printf("%6v :  %+v\n", k, v)
    }
	
	return output
}
