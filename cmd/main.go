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

	print("RUN Register: ")
	read, err := ade.SPI_Read_16bit(ade9000.ADDR_RUN)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%X\n", read)

	println("Calibrating...")
	calibration := ade9000.NewCalibration(ade)
	err = calibration.GetPGA_gain()
	if err != nil {
		log.Fatal(err)
	}
	err = calibration.VGain_calibrate()
	if err != nil {
		log.Fatal(err)
	}
	loop(ade)
}

func resetADE9000(reset_pin gpio.PinIO) {
	reset_pin.Out(gpio.Low)
	time.Sleep(50 * time.Millisecond)
	reset_pin.Out(gpio.High)
	time.Sleep(100 * time.Millisecond)
	println("Reset done")
}

func readRegisterData(ade ade9000.ADE9000Interface) {
	// print("AIRMS: ")
	// airms, err := ade.SPI_Read_32bit(ade9000.ADDR_AIRMS)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%X\n", airms)
	voltageRms := ade9000.VoltageRMSRegs{}
	// activePower := ade9000.ActivePowerRegs{}
	ade.ReadVoltageRMSRegs(&voltageRms)
	// ade.ReadActivePowerRegs(&activePower)
	print("AVRMS: ")
	fmt.Printf("%#X ", uint32(voltageRms.VoltageRMSReg_A))
	print("BVRMS: ")
	fmt.Printf("%#X ", uint32(voltageRms.VoltageRMSReg_B))
	print("CVRMS: ")
	fmt.Printf("%#X\n", uint32(voltageRms.VoltageRMSReg_C))
	// print("AWATT: ")
	// fmt.Printf("%X\n", activePower.ActivePowerReg_A)
	// print("BWATT: ")
	// fmt.Printf("%X\n", activePower.ActivePowerReg_B)
	// print("CWATT: ")
	// fmt.Printf("%X\n", activePower.ActivePowerReg_C)
}

func readResampledData(ade ade9000.ADE9000Interface) {
	ade.SPI_Write_16bit(ade9000.ADDR_WFB_CFG, 0x1000)
	ade.SPI_Write_16bit(ade9000.ADDR_WFB_CFG, 0x1010)
	time.Sleep(100 * time.Millisecond)
	resampledData, err := ade.SPI_Burst_Read_Resampled_Wfb(0x800, ade9000.WFB_ELEMENT_ARRAY_SIZE)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < ade9000.WFB_ELEMENT_ARRAY_SIZE; i++ {
		print("VA: ")
		fmt.Printf("%X\n", resampledData.VA_Resampled[i])
		print("VB: ")
		fmt.Printf("%X\n", resampledData.VB_Resampled[i])
		print("VC: ")
		fmt.Printf("%X\n", resampledData.VC_Resampled[i])
		print("IA: ")
		fmt.Printf("%X\n", resampledData.IA_Resampled[i])
		print("IB: ")
		fmt.Printf("%X\n", resampledData.IB_Resampled[i])
		print("IC: ")
		fmt.Printf("%X\n", resampledData.IC_Resampled[i])
	}
}

func loop(ade ade9000.ADE9000Interface) {
	for {
		readRegisterData(ade)
		// readResampledData(ade)
		time.Sleep(500 * time.Millisecond)
	}
}
