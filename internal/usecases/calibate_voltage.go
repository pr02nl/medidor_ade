package usecases

import (
	"time"

	"github.com/pr02nl/medidor_ade/internal/entity"
	"github.com/pr02nl/medidor_ade/pkg/ade9000"
)

type CalibrateVoltageUseCase struct {
	ade9000           ade9000.ADE9000Interface
	medidorRepository entity.MedidorRepositoryInterface
}

func NewCalibrateVoltageUseCase(medidorRepository entity.MedidorRepositoryInterface, ade9000 ade9000.ADE9000Interface) *CalibrateVoltageUseCase {
	return &CalibrateVoltageUseCase{medidorRepository: medidorRepository, ade9000: ade9000}
}

func (u *CalibrateVoltageUseCase) Execute() error {
	println("Calibrating Voltage...")
	calibration := ade9000.NewCalibration(u.ade9000)
	if err := calibration.GetPGA_gain(); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)
	if err := calibration.VGain_calibrate(); err != nil {
		return err
	}
	return nil
}
