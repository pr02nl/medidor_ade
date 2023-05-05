package usecases

import (
	"time"

	"github.com/pr02nl/medidor_ade/internal/entity"
	"github.com/pr02nl/medidor_ade/pkg/ade9000"
)

type CalibratePhaseUseCase struct {
	ade9000           ade9000.ADE9000Interface
	medidorRepository entity.MedidorRepositoryInterface
}

func NewCalibratePhaseUseCase(medidorRepository entity.MedidorRepositoryInterface, ade9000 ade9000.ADE9000Interface) *CalibratePhaseUseCase {
	return &CalibratePhaseUseCase{medidorRepository: medidorRepository, ade9000: ade9000}
}

func (u *CalibratePhaseUseCase) calibrationEnergyRegisterSetup() error {
	if err := u.ade9000.SPI_Write_32bit(ade9000.ADDR_MASK0, ade9000.EGY_INTERRUPT_MASK0); err != nil { //Enable EGYRDY interrupt
		return err
	}
	if err := u.ade9000.SPI_Write_16bit(ade9000.ADDR_EGY_TIME, ade9000.EGYACCTIME); err != nil { //accumulate EGY_TIME+1 samples (8000 = 1sec)
		return err
	}
	epcfgRegister, err := u.ade9000.SPI_Read_16bit(ade9000.ADDR_EP_CFG) //Read EP_CFG register
	if err != nil {
		return err
	}
	epcfgRegister |= ade9000.CALIBRATION_EGY_CFG //Write the settings and enable accumulation
	if err = u.ade9000.SPI_Write_16bit(ade9000.ADDR_EP_CFG, epcfgRegister); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	u.ade9000.SPI_Write_32bit(ade9000.ADDR_STATUS0, 0xFFFFFFFF) //Clear all interrupts
	return nil
}

func (u *CalibratePhaseUseCase) Execute() error {
	println("Calibrating Phase...")
	calibration := ade9000.NewCalibration(u.ade9000)
	if err := calibration.GetPGA_gain(); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)
	if err := calibration.Phase_calibrate(); err != nil {
		return err
	}
	return nil
}
