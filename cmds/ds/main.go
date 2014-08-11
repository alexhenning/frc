package main

import (
	"fmt"
	"log"
	"strings"

	"go.thoriumrobotics.com/frc.v0/driverstation"
)

func main() {
	fmt.Println("Started")

	ds := driverstation.New(190).
		SetAlliance(driverstation.Blue).
		SetStation(driverstation.Station2)

	fmt.Println("Connecting")
	err := ds.Connect()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Running")
	go ds.Run()
	for {
		fmt.Print("> ")
		var input string
		fmt.Scanln(&input)

		if strings.Contains(input, "e") {
			ds.SetEnabled(true)
		} else if strings.Contains(input, "d") {
			ds.SetEnabled(false)
		}

		if strings.Contains(input, "t") {
			ds.SetState(driverstation.Teleop)
		} else if strings.Contains(input, "a") {
			ds.SetState(driverstation.Auto)
		} else if strings.Contains(input, "l") {
			ds.SetState(driverstation.Test)
		}

		if strings.Contains(input, "q") {
			fmt.Println("Quiting")
			return
		}
	}
}
