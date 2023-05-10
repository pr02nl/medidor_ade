package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/pr02nl/medidor_ade/internal/entity"
	"github.com/pr02nl/medidor_ade/internal/infra/database"
	"github.com/pr02nl/medidor_ade/internal/usecases"
	"github.com/pr02nl/medidor_ade/pkg/ade9000"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"

	//sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
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
	// IRQ0              = "GPIO10"
	IRQ1             = "GPIO20"
	CF4_DREADY_EVENT = "GPIO68"
	CF3_ZX           = "GPIO71"
)

func main() {
	println("Starting...")
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	println("Host initialized")
	// println("Loading configs...")
	// configs, err := configs.LoadConfig(".")
	// if err != nil {
	// 	panic(err)
	// }
	// println("Configs loaded")
	println("Connecting to database...")
	db, err := sql.Open("sqlite3", "medidor.db")
	if err != nil {
		panic(err)
	}
	println("Database connected")
	defer db.Close()

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

	read, err := ade.SPI_Read_16bit(ade9000.ADDR_RUN)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("RUN Register: %#X\n", read)

	medidorRepository := database.NewMedidorRepository(db)
	medidorRepository.InitTable()
	loadUseCase := usecases.NewLoadMedidorUseCase(medidorRepository)
	medidor, err := loadUseCase.Execute()
	if err != nil {
		log.Print(err)
		medidor = &entity.Medidor{
			NominalVoltage:    ade9000.NOMINAL_INPUT_VOLTAGE,
			NominalCurrent:    ade9000.NOMINAL_INPUT_CURRENT,
			Frequency:         ade9000.INPUT_FREQUENCY,
			CurrentTransfer:   ade9000.CURRENT_TRANSFER_FUNCTION,
			VoltageTransfer:   ade9000.VOLTAGE_TRANSFER_FUNCTION,
			CalIrmsCC:         ade9000.CAL_IRMS_CC,
			CalVrmsCC:         ade9000.CAL_VRMS_CC,
			CalPwrCC:          ade9000.CAL_POWER_CC,
			CalEnergyCC:       ade9000.CAL_ENERGY_CC,
			CalibratedVoltage: false,
			CalibratedCurrent: false,
			CalibratedPower:   false,
			CalibratedPhase:   false,
		}
		log.Print("Creating medidor...")
		createUseCase := usecases.NewCreateMedidorUseCase(medidorRepository)
		medidor, err = createUseCase.Execute(medidor)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("Medidor: %+v\n", medidor.ID)

	if medidor.CalibratedVoltage && medidor.CalibratedCurrent && medidor.CalibratedPower && medidor.CalibratedPhase {
		log.Println("Medidor já calibrado!")
	} else {
		cal, err := calibration(ade, *medidor)
		if err != nil {
			log.Fatal(err)
		}
		if cal {
			if !medidor.CalibratedVoltage {
				fmt.Printf("Adicione uma tensão de %vV e tecle enter\n", medidor.NominalVoltage)
				fmt.Scanln()
				err = usecases.NewCalibrateVoltageUseCase(medidorRepository, ade).Execute()
				if err != nil {
					log.Fatal(err)
				}
				medidor.CalibratedVoltage = true
				_, err = usecases.NewUpdateMedidorUseCase(medidorRepository).Execute(medidor)
				if err != nil {
					log.Fatal(err)
				}
			}
			if !medidor.CalibratedCurrent {
				fmt.Printf("Adicione uma corrente de %vA e tecle enter\n", medidor.NominalCurrent)
				fmt.Scanln()
				err = usecases.NewCalibrateCurrentUseCase(medidorRepository, ade).Execute()
				if err != nil {
					log.Fatal(err)
				}
				medidor.CalibratedCurrent = true
				_, err = usecases.NewUpdateMedidorUseCase(medidorRepository).Execute(medidor)
				if err != nil {
					log.Fatal(err)
				}
			}
			if !medidor.CalibratedPower {
				fmt.Printf("Adicione a tensão nominal (%vV) e corrente nominal (%vA) e FP = 1 e tecle enter\n", medidor.NominalVoltage, medidor.NominalCurrent)
				fmt.Scanln()
				err = usecases.NewCalibratePowerUseCase(medidorRepository, ade).Execute()
				if err != nil {
					log.Fatal(err)
				}
				medidor.CalibratedPower = true
				_, err = usecases.NewUpdateMedidorUseCase(medidorRepository).Execute(medidor)
				if err != nil {
					log.Fatal(err)
				}
			}
			if !medidor.CalibratedPhase {
				fmt.Printf("Adicione a tensão nominal (%vV) e corrente nominal (%vA) e FP = 0.5 e tecle enter\n", medidor.NominalVoltage, medidor.NominalCurrent)
				fmt.Scanln()
				err = usecases.NewCalibratePhaseUseCase(medidorRepository, ade).Execute()
				if err != nil {
					log.Fatal(err)
				}
				medidor.CalibratedPhase = true
				_, err = usecases.NewUpdateMedidorUseCase(medidorRepository).Execute(medidor)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	// println("Calibrating...")
	// calibration := ade9000.NewCalibration(ade)
	// err = calibration.GetPGA_gain()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// time.Sleep(500 * time.Millisecond)
	// err = calibration.VGain_calibrate()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// loop(ade)
}

func calibration(ade ade9000.ADE9000Interface, medidor entity.Medidor) (bool, error) {
	var calibration string
	fmt.Println("Medidor ainda não calibrado, deseja iniciar a calibração agora?")
	fmt.Scanln(&calibration)
	if calibration == "s" || calibration == "S" {
		return true, nil
	} else {
		return false, errors.New("calibração cancelada: " + calibration)
	}
}

func resetADE9000(reset_pin gpio.PinIO) {
	reset_pin.Out(gpio.Low)
	time.Sleep(50 * time.Millisecond)
	reset_pin.Out(gpio.High)
	time.Sleep(100 * time.Millisecond)
	println("Reset done")
}

// func readRegisterData(ade ade9000.ADE9000Interface) {
// 	voltageRms := ade9000.VoltageRMS{}
// 	currentRms := ade9000.CurrentRMS{}
// 	activePower := ade9000.Power{}
// 	reactivePower := ade9000.Power{}
// 	aparentPower := ade9000.Power{}
// 	voltageRmsRegs := ade9000.VoltageRMSRegs{}
// 	currentRmsRegs := ade9000.CurrentRMSRegs{}
// 	activePowerRegs := ade9000.ActivePowerRegs{}
// 	aparentPowerRegs := ade9000.ApparentPowerRegs{}
// 	reactivePowerRegs := ade9000.ReactivePowerRegs{}
// 	powerFactorRegs := ade9000.PowerFactorRegs{}
// 	periodRegs := ade9000.PeriodRegs{}
// 	angleRegs := ade9000.AngleRegs{}
// 	tempRegs := ade9000.TemperatureRegnValue{}
// 	ade.ReadVoltageRMSRegs(&voltageRmsRegs)
// 	ade.ReadCurrentRMSRegs(&currentRmsRegs)
// 	ade.ReadActivePowerRegs(&activePowerRegs)
// 	ade.ReadReactivePowerRegs(&reactivePowerRegs)
// 	ade.ReadApparentPowerRegs(&aparentPowerRegs)
// 	ade.ReadPowerFactorRegsnValues(&powerFactorRegs)
// 	ade.ReadPeriodRegsnValues(&periodRegs)
// 	ade.ReadAngleRegsnValues(&angleRegs)
// 	ade.ReadTempRegnValue(&tempRegs)
// 	voltageRms.VoltageRMS_A = (float64(voltageRmsRegs.VoltageRMSReg_A) * ade9000.CAL_VRMS_CC) / math.Pow10(6)
// 	voltageRms.VoltageRMS_B = (float64(voltageRmsRegs.VoltageRMSReg_B) * ade9000.CAL_VRMS_CC) / math.Pow10(6)
// 	voltageRms.VoltageRMS_C = (float64(voltageRmsRegs.VoltageRMSReg_C) * ade9000.CAL_VRMS_CC) / math.Pow10(6)

// 	currentRms.CurrentRMS_A = (float64(currentRmsRegs.CurrentRMSReg_A) * ade9000.CAL_IRMS_CC) / math.Pow10(6)
// 	currentRms.CurrentRMS_B = (float64(currentRmsRegs.CurrentRMSReg_B) * ade9000.CAL_IRMS_CC) / math.Pow10(6)
// 	currentRms.CurrentRMS_C = (float64(currentRmsRegs.CurrentRMSReg_C) * ade9000.CAL_IRMS_CC) / math.Pow10(6)

// 	activePower.Power_A = (float64(activePowerRegs.ActivePowerReg_A) * ade9000.CAL_POWER_CC) / math.Pow10(3)
// 	activePower.Power_B = (float64(activePowerRegs.ActivePowerReg_B) * ade9000.CAL_POWER_CC) / math.Pow10(3)
// 	activePower.Power_C = (float64(activePowerRegs.ActivePowerReg_C) * ade9000.CAL_POWER_CC) / math.Pow10(3)

// 	reactivePower.Power_A = (float64(reactivePowerRegs.ReactivePowerReg_A) * ade9000.CAL_POWER_CC) / math.Pow10(3)
// 	reactivePower.Power_B = (float64(reactivePowerRegs.ReactivePowerReg_B) * ade9000.CAL_POWER_CC) / math.Pow10(3)
// 	reactivePower.Power_C = (float64(reactivePowerRegs.ReactivePowerReg_C) * ade9000.CAL_POWER_CC) / math.Pow10(3)

// 	aparentPower.Power_A = (float64(aparentPowerRegs.ApparentPowerReg_A) * ade9000.CAL_POWER_CC) / math.Pow10(3)
// 	aparentPower.Power_B = (float64(aparentPowerRegs.ApparentPowerReg_B) * ade9000.CAL_POWER_CC) / math.Pow10(3)
// 	aparentPower.Power_C = (float64(aparentPowerRegs.ApparentPowerReg_C) * ade9000.CAL_POWER_CC) / math.Pow10(3)

// 	print("AVRMS: ")
// 	fmt.Printf("%f ", voltageRms.VoltageRMS_A)
// 	print("BVRMS: ")
// 	fmt.Printf("%f ", voltageRms.VoltageRMS_B)
// 	print("CVRMS: ")
// 	fmt.Printf("%f\n", voltageRms.VoltageRMS_C)
// }

// func readResampledData(ade ade9000.ADE9000Interface) {
// 	ade.SPI_Write_16bit(ade9000.ADDR_WFB_CFG, 0x1000)
// 	ade.SPI_Write_16bit(ade9000.ADDR_WFB_CFG, 0x1010)
// 	time.Sleep(100 * time.Millisecond)
// 	resampledData, err := ade.SPI_Burst_Read_Resampled_Wfb(0x800, ade9000.WFB_ELEMENT_ARRAY_SIZE)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for i := 0; i < ade9000.WFB_ELEMENT_ARRAY_SIZE; i++ {
// 		print("VA: ")
// 		fmt.Printf("%X\n", resampledData.VA_Resampled[i])
// 		print("VB: ")
// 		fmt.Printf("%X\n", resampledData.VB_Resampled[i])
// 		print("VC: ")
// 		fmt.Printf("%X\n", resampledData.VC_Resampled[i])
// 		print("IA: ")
// 		fmt.Printf("%X\n", resampledData.IA_Resampled[i])
// 		print("IB: ")
// 		fmt.Printf("%X\n", resampledData.IB_Resampled[i])
// 		print("IC: ")
// 		fmt.Printf("%X\n", resampledData.IC_Resampled[i])
// 	}
// }

// func loop(ade ade9000.ADE9000Interface) {
// 	for {
// 		readRegisterData(ade)
// 		// readResampledData(ade)
// 		time.Sleep(500 * time.Millisecond)
// 	}
// }
