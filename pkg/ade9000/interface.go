package ade9000

import "periph.io/x/conn/v3/spi"

const (
	WFB_ELEMENT_ARRAY_SIZE           = 512
	CALIBRATION_CONSTANTS_ARRAY_SIZE = 13
	/*Full scale Codes referred from Datasheet.Respective digital codes are produced when ADC inputs are at full scale. Donot Change. */
	ADE9000_RMS_FULL_SCALE_CODES       = 52702092
	ADE9000_WATT_FULL_SCALE_CODES      = 20694066
	ADE9000_RESAMPLED_FULL_SCALE_CODES = 18196
	ADE9000_PCF_FULL_SCALE_CODES       = 74532013
	ADE90xx_FDSP                       = 8000
)

type ResampledWfbData struct {
	VA_Resampled [WFB_ELEMENT_ARRAY_SIZE]int16
	IA_Resampled [WFB_ELEMENT_ARRAY_SIZE]int16
	VB_Resampled [WFB_ELEMENT_ARRAY_SIZE]int16
	IB_Resampled [WFB_ELEMENT_ARRAY_SIZE]int16
	VC_Resampled [WFB_ELEMENT_ARRAY_SIZE]int16
	IC_Resampled [WFB_ELEMENT_ARRAY_SIZE]int16
	IN_Resampled [WFB_ELEMENT_ARRAY_SIZE]int16
}

type ActivePowerRegs struct {
	ActivePowerReg_A int32
	ActivePowerReg_B int32
	ActivePowerReg_C int32
}

type ReactivePowerRegs struct {
	ReactivePowerReg_A int32
	ReactivePowerReg_B int32
	ReactivePowerReg_C int32
}

type ApparentPowerRegs struct {
	ApparentPowerReg_A int32
	ApparentPowerReg_B int32
	ApparentPowerReg_C int32
}

type VoltageRMSRegs struct {
	VoltageRMSReg_A int32
	VoltageRMSReg_B int32
	VoltageRMSReg_C int32
}

type CurrentRMSRegs struct {
	CurrentRMSReg_A int32
	CurrentRMSReg_B int32
	CurrentRMSReg_C int32
	CurrentRMSReg_N int32
}

type FundActivePowerRegs struct {
	FundActivePowerReg_A int32
	FundActivePowerReg_B int32
	FundActivePowerReg_C int32
}

type FundReactivePowerRegs struct {
	FundReactivePowerReg_A int32
	FundReactivePowerReg_B int32
	FundReactivePowerReg_C int32
}

type FundApparentPowerRegs struct {
	FundApparentPowerReg_A int32
	FundApparentPowerReg_B int32
	FundApparentPowerReg_C int32
}

type FundVoltageRMSRegs struct {
	FundVoltageRMSReg_A int32
	FundVoltageRMSReg_B int32
	FundVoltageRMSReg_C int32
}

type FundCurrentRMSRegs struct {
	FundCurrentRMSReg_A int32
	FundCurrentRMSReg_B int32
	FundCurrentRMSReg_C int32
}

type HalfVoltageRMSRegs struct {
	HalfVoltageRMSReg_A int32
	HalfVoltageRMSReg_B int32
	HalfVoltageRMSReg_C int32
}

type HalfCurrentRMSRegs struct {
	HalfCurrentRMSReg_A int32
	HalfCurrentRMSReg_B int32
	HalfCurrentRMSReg_C int32
}

type Ten12VoltageRMSRegs struct {
	Ten12VoltageRMSReg_A int32
	Ten12VoltageRMSReg_B int32
	Ten12VoltageRMSReg_C int32
}

type Ten12CurrentRMSRegs struct {
	Ten12CurrentRMSReg_A int32
	Ten12CurrentRMSReg_B int32
	Ten12CurrentRMSReg_C int32
	Ten12CurrentRMSReg_N int32
}

type VoltageTHDRegs struct {
	VoltageTHDReg_A   int32
	VoltageTHDReg_B   int32
	VoltageTHDReg_C   int32
	VoltageTHDValue_A float32
	VoltageTHDValue_B float32
	VoltageTHDValue_C float32
}

type CurrentTHDRegs struct {
	CurrentTHDReg_A   int32
	CurrentTHDReg_B   int32
	CurrentTHDReg_C   int32
	CurrentTHDValue_A float32
	CurrentTHDValue_B float32
	CurrentTHDValue_C float32
}

type PowerFactorRegs struct {
	PowerFactorReg_A   int32
	PowerFactorReg_B   int32
	PowerFactorReg_C   int32
	PowerFactorValue_A float32
	PowerFactorValue_B float32
	PowerFactorValue_C float32
}

type PeriodRegs struct {
	PeriodReg_A      int32
	PeriodReg_B      int32
	PeriodReg_C      int32
	FrequencyValue_A float32
	FrequencyValue_B float32
	FrequencyValue_C float32
}

type AngleRegs struct {
	AngleReg_VA_VB   int16
	AngleReg_VB_VC   int16
	AngleReg_VA_VC   int16
	AngleReg_VA_IA   int16
	AngleReg_VB_IB   int16
	AngleReg_VC_IC   int16
	AngleReg_IA_IB   int16
	AngleReg_IB_IC   int16
	AngleReg_IA_IC   int16
	AngleValue_VA_VB float32
	AngleValue_VB_VC float32
	AngleValue_VA_VC float32
	AngleValue_VA_IA float32
	AngleValue_VB_IB float32
	AngleValue_VC_IC float32
	AngleValue_IA_IB float32
	AngleValue_IB_IC float32
	AngleValue_IA_IC float32
}

type TemperatureRegnValue struct {
	Temperature_Reg int16
	Temperature_Val float32
}

type ADE9000Interface interface {
	SetupADE9000() error
	/*SPI Functions*/
	SPI_Init(SPI_speed uint32, chipSelect_Pin string) (spi.PortCloser, error)
	SPI_Write_16bit(address uint16, data uint16) error
	SPI_Write_32bit(address uint16, data uint32) error
	SPI_Read_16bit(address uint16) (uint16, error)
	SPI_Read_32bit(address uint16) (uint32, error)
	SPI_Burst_Read_Resampled_Wfb(address uint16, read_Element_Length uint16) (*ResampledWfbData, error)
	/*ADE9000 Calculated Parameter Read Functions*/
	ReadActivePowerRegs(data *ActivePowerRegs) error
	ReadReactivePowerRegs(data *ReactivePowerRegs) error
	ReadApparentPowerRegs(data *ApparentPowerRegs) error
	ReadVoltageRMSRegs(data *VoltageRMSRegs) error
	ReadCurrentRMSRegs(data *CurrentRMSRegs) error
	ReadFundActivePowerRegs(data *FundActivePowerRegs) error
	ReadFundReactivePowerRegs(data *FundReactivePowerRegs) error
	ReadFundApparentPowerRegs(data *FundApparentPowerRegs) error
	ReadFundVoltageRMSRegs(data *FundVoltageRMSRegs) error
	ReadFundCurrentRMSRegs(data *FundCurrentRMSRegs) error
	ReadHalfVoltageRMSRegs(data *HalfVoltageRMSRegs) error
	ReadHalfCurrentRMSRegs(data *HalfCurrentRMSRegs) error
	ReadTen12VoltageRMSRegs(data *Ten12VoltageRMSRegs) error
	ReadTen12CurrentRMSRegs(data *Ten12CurrentRMSRegs) error
	ReadVoltageTHDRegsnValues(data *VoltageTHDRegs) error
	ReadCurrentTHDRegsnValues(data *CurrentTHDRegs) error
	ReadPowerFactorRegsnValues(data *PowerFactorRegs) error
	ReadPeriodRegsnValues(data *PeriodRegs) error
	ReadAngleRegsnValues(data *AngleRegs) error
	ReadTempRegnValue(data *TemperatureRegnValue) error
	/*EEPROM Functions*/
	WriteByteToEeprom(dataAddress uint16, data uint8) error
	ReadByteFromEeprom(dataAddress uint16) (uint8, error)
	WriteWordToEeprom(address uint16, data uint32) error
	ReadWordFromEeprom(address uint16) (uint32, error)
}
