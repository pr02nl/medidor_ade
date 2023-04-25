package ade9000

import (
	"errors"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
)

type ADE9000Api struct {
	chipSelect_Pin gpio.PinIO
	spiConn        spi.Conn
}

func NewADE9000Api() *ADE9000Api {
	return &ADE9000Api{}
}

func (ade *ADE9000Api) SetupADE9000() error {
	if err := ade.SPI_Write_16bit(ADDR_PGA_GAIN, ADE9000_PGA_GAIN); err != nil {
		return err
	}
	err := ade.SPI_Write_32bit(ADDR_CONFIG0, ADE9000_CONFIG0)
	if err != nil {
		return err
	}
	err = ade.SPI_Write_16bit(ADDR_CONFIG1, ADE9000_CONFIG1)
	if err != nil {
		return err
	}
	err = ade.SPI_Write_16bit(ADDR_CONFIG2, ADE9000_CONFIG2)
	if err != nil {
		return err
	}
	err = ade.SPI_Write_16bit(ADDR_CONFIG3, ADE9000_CONFIG3)
	if err != nil {
		return err
	}
	err = ade.SPI_Write_16bit(ADDR_ACCMODE, ADE9000_ACCMODE)
	if err != nil {
		return err
	}
	err = ade.SPI_Write_16bit(ADDR_TEMP_CFG, ADE9000_TEMP_CFG)
	if err != nil {
		return err
	}
	err = ade.SPI_Write_16bit(ADDR_ZX_LP_SEL, ADE9000_ZX_LP_SEL)
	if err != nil {
		return err
	}
	err = ade.SPI_Write_32bit(ADDR_MASK0, ADE9000_MASK0)
	if err != nil {
		return err
	}
	err = ade.SPI_Write_32bit(ADDR_MASK1, ADE9000_MASK1)
	if err != nil {
		return err
	}
	err = ade.SPI_Write_32bit(ADDR_EVENT_MASK, ADE9000_EVENT_MASK)
	if err != nil {
		return err
	}
	err = ade.SPI_Write_16bit(ADDR_WFB_CFG, ADE9000_WFB_CFG)
	if err != nil {
		return err
	}
	err = ade.SPI_Write_32bit(ADDR_VLEVEL, ADE9000_VLEVEL)

	if err != nil {
		return err
	}
	err = ade.SPI_Write_32bit(ADDR_DICOEFF, ADE9000_DICOEFF)
	if err != nil {
		return err
	}
	err = ade.SPI_Write_16bit(ADDR_EGY_TIME, ADE9000_EGY_TIME)
	if err != nil {
		return err
	}
	err = ade.SPI_Write_16bit(ADDR_EP_CFG, ADE9000_EP_CFG) //Energy accumulation ON
	if err != nil {
		return err
	}
	err = ade.SPI_Write_16bit(ADDR_RUN, ADE9000_RUN_ON) //DSP ON
	if err != nil {
		return err
	}
	return nil
}

func (ade *ADE9000Api) SPI_Write_16bit(address uint16, data uint16) error {
	temp_address := ((address << 4) & 0xFFF0)
	var err error
	if err = ade.chipSelect_Pin.Out(gpio.Low); err != nil {
		return err
	}
	if err = ade.spiConn.Tx([]byte{byte(temp_address >> 8), byte(temp_address), byte(data >> 8), byte(data)}, nil); err != nil {
		return err
	}
	if err = ade.chipSelect_Pin.Out(gpio.High); err != nil {
		return err
	}
	return nil
}
func (ade *ADE9000Api) SPI_Write_32bit(address uint16, data uint32) error {
	temp_address := ((address << 4) & 0xFFF0)
	var err error
	if err = ade.chipSelect_Pin.Out(gpio.Low); err != nil {
		return err
	}
	if err = ade.spiConn.Tx([]byte{byte(temp_address >> 8), byte(temp_address), byte(data >> 24), byte(data >> 16), byte(data >> 8), byte(data)}, nil); err != nil {
		return err
	}
	if err = ade.chipSelect_Pin.Out(gpio.High); err != nil {
		return err
	}
	return nil
}

func (ade *ADE9000Api) SPI_Read_16bit(address uint16) (uint16, error) {
	temp_address := (((address << 4) & 0xFFF0) + 8)
	var err error
	read := make([]byte, 2)
	if err = ade.chipSelect_Pin.Out(gpio.Low); err != nil {
		return 0, err
	}
	if err = ade.spiConn.Tx([]byte{byte(temp_address >> 8), byte(temp_address)}, read); err != nil {
		return 0, err
	}
	if err = ade.chipSelect_Pin.Out(gpio.High); err != nil {
		return 0, err
	}
	return uint16(read[0])<<8 + uint16(read[1]), nil
}

func (ade *ADE9000Api) SPI_Read_32bit(address uint16) (uint32, error) {
	temp_address := (((address << 4) & 0xFFF0) + 8)
	var err error
	read := make([]byte, 4)
	if err = ade.chipSelect_Pin.Out(gpio.Low); err != nil {
		return 0, err
	}
	if err = ade.spiConn.Tx([]byte{byte(temp_address >> 8), byte(temp_address)}, nil); err != nil {
		return 0, err
	}
	if err = ade.spiConn.Tx([]byte{0x00, 0x00, 0x00, 0x00}, read); err != nil {
		return 0, err
	}
	if err = ade.chipSelect_Pin.Out(gpio.High); err != nil {
		return 0, err
	}
	return uint32(read[0])<<24 + uint32(read[1])<<16 + uint32(read[2])<<8 + uint32(read[3]), nil
}

func (ade *ADE9000Api) SPI_Init(SPI_speed uint32, chipSelect_Pin string) (spi.PortCloser, error) {
	spinSelect := gpioreg.ByName(chipSelect_Pin)
	if spinSelect == nil {
		return nil, errors.New("failed to find chip select pin" + chipSelect_Pin)
	}
	ade.chipSelect_Pin = spinSelect
	// Use spireg SPI port registry to find the first available SPI bus.
	p, err := spireg.Open("")
	if err != nil {
		return nil, err
	}
	// Convert the spi.Port into a spi.Conn so it can be used for communication.
	c, err := p.Connect(physic.MegaHertz, spi.Mode0, 8)
	if err != nil {
		return nil, err
	}
	ade.spiConn = c
	err = ade.chipSelect_Pin.Out(gpio.High)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (ade *ADE9000Api) SPI_Burst_Read_Resampled_Wfb(address uint16, read_Element_Length uint16) (*ResampledWfbData, error) {
	temp_address := (((address << 4) & 0xFFF0) + 8)
	var err error
	read := make([]byte, 2)
	if err = ade.chipSelect_Pin.Out(gpio.Low); err != nil {
		return nil, err
	}
	if err = ade.spiConn.Tx([]byte{byte(temp_address >> 8), byte(temp_address)}, nil); err != nil {
		return nil, err
	}
	data := &ResampledWfbData{}
	for i := uint16(0); i < read_Element_Length; i++ {
		if err = ade.spiConn.Tx([]byte{byte(0x00), byte(0x00)}, read); err != nil {
			return nil, err
		}
		data.IA_Resampled[i] = int16(read[0])<<8 + int16(read[1])
		if err = ade.spiConn.Tx([]byte{byte(0x00), byte(0x00)}, read); err != nil {
			return nil, err
		}
		data.VA_Resampled[i] = int16(read[0])<<8 + int16(read[1])
		if err = ade.spiConn.Tx([]byte{byte(0x00), byte(0x00)}, read); err != nil {
			return nil, err
		}
		data.IB_Resampled[i] = int16(read[0])<<8 + int16(read[1])
		if err = ade.spiConn.Tx([]byte{byte(0x00), byte(0x00)}, read); err != nil {
			return nil, err
		}
		data.VB_Resampled[i] = int16(read[0])<<8 + int16(read[1])
		if err = ade.spiConn.Tx([]byte{byte(0x00), byte(0x00)}, read); err != nil {
			return nil, err
		}
		data.IC_Resampled[i] = int16(read[0])<<8 + int16(read[1])
		if err = ade.spiConn.Tx([]byte{byte(0x00), byte(0x00)}, read); err != nil {
			return nil, err
		}
		data.VC_Resampled[i] = int16(read[0])<<8 + int16(read[1])
		if err = ade.spiConn.Tx([]byte{byte(0x00), byte(0x00)}, read); err != nil {
			return nil, err
		}
		data.IN_Resampled[i] = int16(read[0])<<8 + int16(read[1])
	}
	if err = ade.chipSelect_Pin.Out(gpio.High); err != nil {
		return nil, err
	}
	return data, nil
}

func (ade *ADE9000Api) ReadActivePowerRegs(data *ActivePowerRegs) error {
	data1, err := ade.SPI_Read_32bit(ADDR_AWATT)
	if err != nil {
		return err
	}
	data.ActivePowerReg_A = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_BWATT)
	if err != nil {
		return err
	}
	data.ActivePowerReg_B = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_CWATT)
	if err != nil {
		return err
	}
	data.ActivePowerReg_C = int32(data1)
	return nil
}

func (ade *ADE9000Api) ReadReactivePowerRegs(data *ReactivePowerRegs) error {
	data1, err := ade.SPI_Read_32bit(ADDR_AVAR)
	if err != nil {
		return err
	}
	data.ReactivePowerReg_A = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_BVAR)
	if err != nil {
		return err
	}
	data.ReactivePowerReg_B = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_CVAR)
	if err != nil {
		return err
	}
	data.ReactivePowerReg_C = int32(data1)
	return nil
}

func (ade *ADE9000Api) ReadApparentPowerRegs(data *ApparentPowerRegs) error {
	data1, err := ade.SPI_Read_32bit(ADDR_AVA)
	if err != nil {
		return err
	}
	data.ApparentPowerReg_A = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_BVA)
	if err != nil {
		return err
	}
	data.ApparentPowerReg_B = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_CVA)
	if err != nil {
		return err
	}
	data.ApparentPowerReg_C = int32(data1)
	return nil
}

func (ade *ADE9000Api) ReadVoltageRMSRegs(data *VoltageRMSRegs) error {
	data1, err := ade.SPI_Read_32bit(ADDR_AVRMS)
	if err != nil {
		return err
	}
	data.VoltageRMSReg_A = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_BVRMS)
	if err != nil {
		return err
	}
	data.VoltageRMSReg_B = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_CVRMS)
	if err != nil {
		return err
	}
	data.VoltageRMSReg_C = int32(data1)
	return nil
}

func (ade *ADE9000Api) ReadCurrentRMSRegs(data *CurrentRMSRegs) error {
	data1, err := ade.SPI_Read_32bit(ADDR_AIRMS)
	if err != nil {
		return err
	}
	data.CurrentRMSReg_A = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_BIRMS)
	if err != nil {
		return err
	}
	data.CurrentRMSReg_B = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_CIRMS)
	if err != nil {
		return err
	}
	data.CurrentRMSReg_C = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_NIRMS)
	if err != nil {
		return err
	}
	data.CurrentRMSReg_N = int32(data1)
	return nil
}

func (ade *ADE9000Api) ReadFundActivePowerRegs(data *FundActivePowerRegs) error {
	data1, err := ade.SPI_Read_32bit(ADDR_AFWATT)
	if err != nil {
		return err
	}
	data.FundActivePowerReg_A = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_BFWATT)
	if err != nil {
		return err
	}
	data.FundActivePowerReg_B = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_CFWATT)
	if err != nil {
		return err
	}
	data.FundActivePowerReg_C = int32(data1)
	return nil
}

func (ade *ADE9000Api) ReadFundReactivePowerRegs(data *FundReactivePowerRegs) error {
	data1, err := ade.SPI_Read_32bit(ADDR_AFVAR)
	if err != nil {
		return err
	}
	data.FundReactivePowerReg_A = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_BFVAR)
	if err != nil {
		return err
	}
	data.FundReactivePowerReg_B = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_CFVAR)
	if err != nil {
		return err
	}
	data.FundReactivePowerReg_C = int32(data1)
	return nil
}

func (ade *ADE9000Api) ReadFundApparentPowerRegs(data *FundApparentPowerRegs) error {
	data1, err := ade.SPI_Read_32bit(ADDR_AFVA)
	if err != nil {
		return err
	}
	data.FundApparentPowerReg_A = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_BFVA)
	if err != nil {
		return err
	}
	data.FundApparentPowerReg_B = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_CFVA)
	if err != nil {
		return err
	}
	data.FundApparentPowerReg_C = int32(data1)
	return nil
}

func (ade *ADE9000Api) ReadFundVoltageRMSRegs(data *FundVoltageRMSRegs) error {
	data1, err := ade.SPI_Read_32bit(ADDR_AVFRMS)
	if err != nil {
		return err
	}
	data.FundVoltageRMSReg_A = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_BVFRMS)
	if err != nil {
		return err
	}
	data.FundVoltageRMSReg_B = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_CVFRMS)
	if err != nil {
		return err
	}
	data.FundVoltageRMSReg_C = int32(data1)
	return nil
}

func (ade *ADE9000Api) ReadFundCurrentRMSRegs(data *FundCurrentRMSRegs) error {
	data1, err := ade.SPI_Read_32bit(ADDR_AIFRMS)
	if err != nil {
		return err
	}
	data.FundCurrentRMSReg_A = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_BIFRMS)
	if err != nil {
		return err
	}
	data.FundCurrentRMSReg_B = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_CIFRMS)
	if err != nil {
		return err
	}
	data.FundCurrentRMSReg_C = int32(data1)
	return nil
}

func (ade *ADE9000Api) ReadHalfVoltageRMSRegs(data *HalfVoltageRMSRegs) error {
	data1, err := ade.SPI_Read_32bit(ADDR_AVRMSONE)
	if err != nil {
		return err
	}
	data.HalfVoltageRMSReg_A = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_BVRMSONE)
	if err != nil {
		return err
	}
	data.HalfVoltageRMSReg_B = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_CVRMSONE)
	if err != nil {
		return err
	}
	data.HalfVoltageRMSReg_C = int32(data1)
	return nil
}

func (ade *ADE9000Api) ReadHalfCurrentRMSRegs(data *HalfCurrentRMSRegs) error {
	data1, err := ade.SPI_Read_32bit(ADDR_AIRMSONE)
	if err != nil {
		return err
	}
	data.HalfCurrentRMSReg_A = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_BIRMSONE)
	if err != nil {
		return err
	}
	data.HalfCurrentRMSReg_B = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_CIRMSONE)
	if err != nil {
		return err
	}
	data.HalfCurrentRMSReg_C = int32(data1)
	return nil
}

func (ade *ADE9000Api) ReadTen12VoltageRMSRegs(data *Ten12VoltageRMSRegs) error {
	data1, err := ade.SPI_Read_32bit(ADDR_AVRMS1012)
	if err != nil {
		return err
	}
	data.Ten12VoltageRMSReg_A = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_BVRMS1012)
	if err != nil {
		return err
	}
	data.Ten12VoltageRMSReg_B = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_CVRMS1012)
	if err != nil {
		return err
	}
	data.Ten12VoltageRMSReg_C = int32(data1)
	return nil
}

func (ade *ADE9000Api) ReadTen12CurrentRMSRegs(data *Ten12CurrentRMSRegs) error {
	data1, err := ade.SPI_Read_32bit(ADDR_AIRMS1012)
	if err != nil {
		return err
	}
	data.Ten12CurrentRMSReg_A = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_BIRMS1012)
	if err != nil {
		return err
	}
	data.Ten12CurrentRMSReg_B = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_CIRMS1012)
	if err != nil {
		return err
	}
	data.Ten12CurrentRMSReg_C = int32(data1)
	data1, err = ade.SPI_Read_32bit(ADDR_NIRMS1012)
	if err != nil {
		return err
	}
	data.Ten12CurrentRMSReg_N = int32(data1)
	return nil
}

func (ade *ADE9000Api) ReadVoltageTHDRegsnValues(data *VoltageTHDRegs) error {
	data1, err := ade.SPI_Read_32bit(ADDR_AVTHD)
	if err != nil {
		return err
	}
	data.VoltageTHDReg_A = int32(data1)
	data.VoltageTHDValue_A = float32(data.VoltageTHDReg_A*100) / 134217728
	data1, err = ade.SPI_Read_32bit(ADDR_BVTHD)
	if err != nil {
		return err
	}
	data.VoltageTHDReg_B = int32(data1)
	data.VoltageTHDValue_B = float32(data.VoltageTHDReg_B*100) / 134217728
	data1, err = ade.SPI_Read_32bit(ADDR_CVTHD)
	if err != nil {
		return err
	}
	data.VoltageTHDReg_C = int32(data1)
	data.VoltageTHDValue_C = float32(data.VoltageTHDReg_C*100) / 134217728
	return nil
}

func (ade *ADE9000Api) ReadCurrentTHDRegsnValues(data *CurrentTHDRegs) error {
	data1, err := ade.SPI_Read_32bit(ADDR_AITHD)
	if err != nil {
		return err
	}
	data.CurrentTHDReg_A = int32(data1)
	data.CurrentTHDValue_A = float32(data.CurrentTHDReg_A*100) / 134217728
	data1, err = ade.SPI_Read_32bit(ADDR_BITHD)
	if err != nil {
		return err
	}
	data.CurrentTHDReg_B = int32(data1)
	data.CurrentTHDValue_B = float32(data.CurrentTHDReg_B*100) / 134217728
	data1, err = ade.SPI_Read_32bit(ADDR_CITHD)
	if err != nil {
		return err
	}
	data.CurrentTHDReg_C = int32(data1)
	data.CurrentTHDValue_C = float32(data.CurrentTHDReg_C*100) / 134217728
	return nil
}

func (ade *ADE9000Api) ReadPowerFactorRegsnValues(data *PowerFactorRegs) error {
	data1, err := ade.SPI_Read_32bit(ADDR_APF)
	if err != nil {
		return err
	}
	data.PowerFactorReg_A = int32(data1)
	data.PowerFactorValue_A = float32(data.PowerFactorReg_A) / 134217728
	data1, err = ade.SPI_Read_32bit(ADDR_BPF)
	if err != nil {
		return err
	}
	data.PowerFactorReg_B = int32(data1)
	data.PowerFactorValue_B = float32(data.PowerFactorReg_B) / 134217728
	data1, err = ade.SPI_Read_32bit(ADDR_CPF)
	if err != nil {
		return err
	}
	data.PowerFactorReg_C = int32(data1)
	data.PowerFactorValue_C = float32(data.PowerFactorReg_C) / 134217728
	return nil
}

func (ade *ADE9000Api) ReadPeriodRegsnValues(data *PeriodRegs) error {
	data1, err := ade.SPI_Read_32bit(ADDR_APERIOD)
	if err != nil {
		return err
	}
	data.PeriodReg_A = int32(data1)
	data.FrequencyValue_A = (8000 * 65536) / float32(data.PeriodReg_A+1)
	data1, err = ade.SPI_Read_32bit(ADDR_BPERIOD)
	if err != nil {
		return err
	}
	data.PeriodReg_B = int32(data1)
	data.FrequencyValue_B = (8000 * 65536) / float32(data.PeriodReg_B+1)
	data1, err = ade.SPI_Read_32bit(ADDR_CPERIOD)
	if err != nil {
		return err
	}
	data.PeriodReg_C = int32(data1)
	data.FrequencyValue_C = (8000 * 65536) / float32(data.PeriodReg_C+1)
	return nil
}

func (ade *ADE9000Api) ReadAngleRegsnValues(data *AngleRegs) error {
	var mulConstant float32
	temp, err := ade.SPI_Read_16bit(ADDR_ACCMODE)
	if err != nil {
		return err
	}
	if (temp & 0x0100) >= 0 {
		mulConstant = 0.02109375
	} else {
		mulConstant = 0.017578125
	}
	tempReg, err := ade.SPI_Read_32bit(ADDR_ANGL_VA_VB)
	if err != nil {
		return err
	}
	data.AngleReg_VA_VB = int16(tempReg)
	data.AngleValue_VA_VB = float32(data.AngleReg_VA_VB) * mulConstant
	tempReg, err = ade.SPI_Read_32bit(ADDR_ANGL_VB_VC)
	if err != nil {
		return err
	}
	data.AngleReg_VB_VC = int16(tempReg)
	data.AngleValue_VB_VC = float32(data.AngleReg_VB_VC) * mulConstant
	tempReg, err = ade.SPI_Read_32bit(ADDR_ANGL_VA_VC)
	if err != nil {
		return err
	}
	data.AngleReg_VA_VC = int16(tempReg)
	data.AngleValue_VA_VC = float32(data.AngleReg_VA_VC) * mulConstant
	tempReg, err = ade.SPI_Read_32bit(ADDR_ANGL_VA_IA)
	if err != nil {
		return err
	}
	data.AngleReg_VA_IA = int16(tempReg)
	data.AngleValue_VA_IA = float32(data.AngleReg_VA_IA) * mulConstant
	tempReg, err = ade.SPI_Read_32bit(ADDR_ANGL_VB_IB)
	if err != nil {
		return err
	}
	data.AngleReg_VB_IB = int16(tempReg)
	data.AngleValue_VB_IB = float32(data.AngleReg_VB_IB) * mulConstant
	tempReg, err = ade.SPI_Read_32bit(ADDR_ANGL_VC_IC)
	if err != nil {
		return err
	}
	data.AngleReg_VC_IC = int16(tempReg)
	data.AngleValue_VC_IC = float32(data.AngleReg_VC_IC) * mulConstant
	tempReg, err = ade.SPI_Read_32bit(ADDR_ANGL_IA_IB)
	if err != nil {
		return err
	}
	data.AngleReg_IA_IB = int16(tempReg)
	data.AngleValue_IA_IB = float32(data.AngleReg_IA_IB) * mulConstant
	tempReg, err = ade.SPI_Read_32bit(ADDR_ANGL_IB_IC)
	if err != nil {
		return err
	}
	data.AngleReg_IB_IC = int16(tempReg)
	data.AngleValue_IB_IC = float32(data.AngleReg_IB_IC) * mulConstant
	tempReg, err = ade.SPI_Read_32bit(ADDR_ANGL_IA_IC)
	if err != nil {
		return err
	}
	data.AngleReg_IA_IC = int16(tempReg)
	data.AngleValue_IA_IC = float32(data.AngleReg_IA_IC) * mulConstant
	return nil
}

func (ade *ADE9000Api) ReadTempRegnValue(data *TemperatureRegnValue) error {
	err := ade.SPI_Write_16bit(ADDR_TEMP_CFG, ADE9000_TEMP_CFG) //Start temperature acquisition cycle with settings in defined in ADE9000_TEMP_CFG
	if err != nil {
		return err
	}
	time.Sleep(2 * time.Millisecond) //Wait for temperature acquisition cycle to complete
	trim, err := ade.SPI_Read_32bit(ADDR_TEMP_TRIM)
	if err != nil {
		return err
	}
	gain := uint16(trim & 0xFFFF)                      //Extract 16 LSB
	offset := uint16((trim >> 16) & 0xFFFF)            //Extract 16 MSB
	tempReg, err := ade.SPI_Read_16bit(ADDR_TEMP_RSLT) //Read Temperature result register
	if err != nil {
		return err
	}
	tempValue := float32(offset>>5) - (float32(tempReg*gain) / float32(65536))
	data.Temperature_Reg = int16(tempReg)
	data.Temperature_Val = tempValue
	return nil
}

func (ade *ADE9000Api) WriteByteToEeprom(dataAddress uint16, data uint8) error {
	return nil
}

func (ade *ADE9000Api) ReadByteFromEeprom(dataAddress uint16) (uint8, error) {
	return 0, nil
}

func (ade *ADE9000Api) WriteWordToEeprom(address uint16, data uint32) error {
	return nil
}

func (ade *ADE9000Api) ReadWordFromEeprom(address uint16) (uint32, error) {
	return 0, nil
}
