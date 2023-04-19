package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pr02nl/medidor_ade/pkg/ade9000"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

/*
// PA8 -> CS
// PA7 -> PM_1
// PA9 -> ADE9000_RESET_PIN
// PA10 -> IRQ0
// PA20 -> IRQ1
// PC4 -> CF4/DREADY/EVENT
// PC7 -> CF3/ZX
*/

const (
	CS                = "GPIO8"
	PM_1              = "GPIO7"
	ADE9000_RESET_PIN = "GPIO9"
	IRQ0              = "GPIO10"
	IRQ1              = "GPIO20"
	CF4_DREADY_EVENT  = "GPIO68"
	CF3_ZX            = "GPIO71"
)

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// all := gpioreg.All()
	// for _, pin := range all {
	// 	fmt.Printf("Pin: %v\n", pin)
	// }

	pm_1 := gpioreg.ByName(PM_1)
	if pm_1 == nil {
		log.Fatal("Failed to find " + PM_1)
	}
	pm_1.Out(gpio.Low)
	reset_pin := gpioreg.ByName(ADE9000_RESET_PIN)
	if reset_pin == nil {
		log.Fatal("Failed to find " + ADE9000_RESET_PIN)
	}
	reset_pin.Out(gpio.High)
	resetADE9000(reset_pin)
	time.Sleep(100 * time.Millisecond)

	ade := ade9000.NewADE9000Api()
	spi, err := ade.SPI_Init(1, CS)
	if err != nil {
		log.Fatal(err)
	}
	defer spi.Close()

	err = ade.SetupADE9000()
	if err != nil {
		log.Fatal(err)
	}

	print("Reading ADE9000 registers\n")
	read, err := ade.SPI_Read_16bit(ade9000.ADDR_RUN)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", read)
}

func resetADE9000(reset_pin gpio.PinIO) {
	reset_pin.Out(gpio.Low)
	time.Sleep(50 * time.Millisecond)
	reset_pin.Out(gpio.High)
	time.Sleep(100 * time.Millisecond)
	println("Reset done")
}
