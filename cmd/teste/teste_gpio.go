package main

import (
	"fmt"
	"log"

	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

func main() {
	println("Starting...")
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	println("Host initialized")

	all := gpioreg.All()
	for _, pin := range all {
		fmt.Printf("Pin: %v\n", pin)
	}
}
