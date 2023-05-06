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
