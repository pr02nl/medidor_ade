package ade9000

import (
	"fmt"
	"math"
	"time"
)

const (
	IGAIN_CAL_REG_SIZE = 4
	VGAIN_CAL_REG_SIZE = 3
	PHCAL_CAL_REG_SIZE = 3
	PGAIN_CAL_REG_SIZE = 3
	EGY_REG_SIZE       = 3
	ACCUMULATION_TIME  = 5
)

type Calibration struct {
	ADE                                 ADE9000Interface
	XIgain_registers                    [IGAIN_CAL_REG_SIZE]int32 //order [AIGAIN, BIGAIN, CIGAIN, NIGAIN]
	XIgain_register_address             [IGAIN_CAL_REG_SIZE]int32
	XIrms_registers                     [IGAIN_CAL_REG_SIZE]int32
	XIrms_registers_address             [IGAIN_CAL_REG_SIZE]int32
	XVgain_registers                    [VGAIN_CAL_REG_SIZE]int32
	XVgain_register_address             [VGAIN_CAL_REG_SIZE]int32
	XVrms_registers                     [VGAIN_CAL_REG_SIZE]int32
	XVrms_registers_address             [VGAIN_CAL_REG_SIZE]int32
	XPhcal_registers                    [PHCAL_CAL_REG_SIZE]int32 //order [APHCAL, BPHCAL, CPHCAL]
	XPhcal_register_address             [PHCAL_CAL_REG_SIZE]int32
	XWATTHRHI_registers                 [PHCAL_CAL_REG_SIZE]int32
	XWATTHRHI_registers_address         [PHCAL_CAL_REG_SIZE]int32
	XVARHRHI_registers                  [PHCAL_CAL_REG_SIZE]int32
	XVARHRHI_registers_address          [PHCAL_CAL_REG_SIZE]int32
	XPgain_registers                    [PGAIN_CAL_REG_SIZE]int32 //order [APGAIN, BPGAIN, CPGAIN]
	XPgain_register_address             [PGAIN_CAL_REG_SIZE]int32
	AccumulatedActiveEnergy_registers   [EGY_REG_SIZE]int32
	AccumulatedReactiveEnergy_registers [EGY_REG_SIZE]int32
	CalCurrentPGA_gain                  int8
	CalVoltagePGA_gain                  int8
	CalibrationDataToEEPROM             [CALIBRATION_CONSTANTS_ARRAY_SIZE]uint32
}

func NewCalibration(ade9000Api ADE9000Interface) *Calibration {
	return &Calibration{
		ADE:                         ade9000Api,
		XIgain_register_address:     [IGAIN_CAL_REG_SIZE]int32{ADDR_AIGAIN, ADDR_BIGAIN, ADDR_CIGAIN, ADDR_NIGAIN},
		XIrms_registers_address:     [IGAIN_CAL_REG_SIZE]int32{ADDR_AIRMS, ADDR_BIRMS, ADDR_CIRMS, ADDR_NIRMS},
		XVgain_register_address:     [VGAIN_CAL_REG_SIZE]int32{ADDR_AVGAIN, ADDR_BVGAIN, ADDR_CVGAIN},
		XVrms_registers_address:     [VGAIN_CAL_REG_SIZE]int32{ADDR_AVRMS, ADDR_BVRMS, ADDR_CVRMS},
		XPhcal_register_address:     [PHCAL_CAL_REG_SIZE]int32{ADDR_APHCAL0, ADDR_BPHCAL0, ADDR_CPHCAL0},
		XWATTHRHI_registers_address: [PHCAL_CAL_REG_SIZE]int32{ADDR_AWATTHR_HI, ADDR_BWATTHR_HI, ADDR_CWATTHR_HI},
		XVARHRHI_registers_address:  [PHCAL_CAL_REG_SIZE]int32{ADDR_AVARHR_HI, ADDR_BVARHR_HI, ADDR_CVARHR_HI},
		XPgain_register_address:     [PGAIN_CAL_REG_SIZE]int32{ADDR_APGAIN, ADDR_BPGAIN, ADDR_CPGAIN},
		CalCurrentPGA_gain:          0,
		CalVoltagePGA_gain:          0,
	}
}

func (calibration *Calibration) GetPGA_gain() error {
	var pgaGainRegister uint16
	var temp uint16
	pgaGainRegister, err := calibration.ADE.SPI_Read_16bit(ADDR_PGA_GAIN)
	if err != nil {
		return err
	}
	fmt.Printf("PGA Gain Register: %#X\n", pgaGainRegister)
	temp = pgaGainRegister & (0x0003) //extract gain of current channel
	// 00-->Gain 1: 01-->Gain 2: 10/11-->Gain 4
	if temp == 0 {
		calibration.CalCurrentPGA_gain = 1
	} else if temp == 1 {
		calibration.CalCurrentPGA_gain = 2
	} else {
		calibration.CalCurrentPGA_gain = 4
	}
	temp = (pgaGainRegister >> 8) & (0x0003) //extract gain of voltage channel
	// 00-->Gain 1: 01-->Gain 2: 10/11-->Gain 4
	if temp == 0 {
		calibration.CalVoltagePGA_gain = 1
	} else if temp == 1 {
		calibration.CalVoltagePGA_gain = 2
	} else {
		calibration.CalVoltagePGA_gain = 4
	}
	return nil
}

func (calibration *Calibration) IGain_calibrate() error {
	temp := ADE9000_RMS_FULL_SCALE_CODES * CURRENT_TRANSFER_FUNCTION * float32(calibration.CalCurrentPGA_gain) * NOMINAL_INPUT_CURRENT * math.Sqrt2
	expectedCodes := int32(temp) //Round off
	for i := 0; i < IGAIN_CAL_REG_SIZE; i++ {
		actualCodes, err := calibration.ADE.SPI_Read_32bit(uint16(calibration.XIrms_registers_address[i]))
		if err != nil {
			return err
		}
		temp = ((float32(expectedCodes) / float32(actualCodes)) - 1) * 134217728 //calculate the gain.
		calibration.XIgain_registers[i] = int32(temp)
	}
	return nil
}

func (calibration *Calibration) VGain_calibrate() error {
	temp := ADE9000_RMS_FULL_SCALE_CODES * VOLTAGE_TRANSFER_FUNCTION * float32(calibration.CalVoltagePGA_gain) * NOMINAL_INPUT_VOLTAGE * math.Sqrt2
	expectedCodes := int32(temp) //Round off
	fmt.Printf("Expected VRMS Code: %#X\n", expectedCodes)
	for i := 0; i < VGAIN_CAL_REG_SIZE; i++ {
		actualCodes, err := calibration.ADE.SPI_Read_32bit(uint16(calibration.XVrms_registers_address[i]))
		if err != nil {
			return err
		}
		fmt.Printf("CH %d, Actual VRMS Code: %#X\n", i, actualCodes)
		temp = ((float32(expectedCodes) / float32(actualCodes)) - 1) * 134217728 //calculate the gain.
		calibration.XVgain_registers[i] = int32(temp)
		fmt.Printf("CH %d, Voltage Gain Register: %#X\n", i, calibration.XVgain_registers[i])
		err = calibration.ADE.SPI_Write_32bit(uint16(calibration.XVgain_register_address[i]), uint32(calibration.XVgain_registers[i]))
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Calibration) DegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func (calibration *Calibration) Phase_calibrate() error {
	time.Sleep((ACCUMULATION_TIME + 1) * time.Millisecond)
	omega := 2 * math.Pi * INPUT_FREQUENCY / ADE90xx_FDSP
	for i := 0; i < PHCAL_CAL_REG_SIZE; i++ {
		actualActiveEnergyCode := calibration.AccumulatedActiveEnergy_registers[i]
		actualReactiveEnergyCode := calibration.AccumulatedReactiveEnergy_registers[i]
		errorAngle := -1 * math.Atan((float64(actualActiveEnergyCode)*math.Sin(calibration.DegreesToRadians(CALIBRATION_ANGLE_DEGREES))-float64(actualReactiveEnergyCode)*math.Cos(calibration.DegreesToRadians(CALIBRATION_ANGLE_DEGREES)))/(float64(actualActiveEnergyCode)*math.Cos(calibration.DegreesToRadians(CALIBRATION_ANGLE_DEGREES))+float64(actualReactiveEnergyCode)*math.Sin(calibration.DegreesToRadians(CALIBRATION_ANGLE_DEGREES))))
		temp := ((math.Sin(errorAngle-omega) + math.Sin(omega)) / (math.Sin(2*omega - errorAngle))) * 134217728
		calibration.XPhcal_registers[i] = int32(temp)
	}
	return nil
}

func (calibration *Calibration) PGain_calibrate(pGaincalPF float32) error {
	time.Sleep((ACCUMULATION_TIME + 1) * time.Millisecond)
	temp := (ADE90xx_FDSP * NOMINAL_INPUT_VOLTAGE * NOMINAL_INPUT_CURRENT * CALIBRATION_ACC_TIME * CURRENT_TRANSFER_FUNCTION * float32(calibration.CalCurrentPGA_gain) * VOLTAGE_TRANSFER_FUNCTION * float32(calibration.CalVoltagePGA_gain) * ADE9000_WATT_FULL_SCALE_CODES * 2 * (pGaincalPF)) / (8192)
	expectedActiveEnergyCode := int32(temp)
	for i := 0; i < PGAIN_CAL_REG_SIZE; i++ {
		actualActiveEnergyCode := calibration.AccumulatedActiveEnergy_registers[i]
		temp = ((float32(expectedActiveEnergyCode) / float32(actualActiveEnergyCode)) - 1) * 134217728 //calculate the gain.
		calibration.XPgain_registers[i] = int32(temp)
	}
	return nil
}
