package usecases

import (
	"time"

	"github.com/pr02nl/medidor_ade/internal/entity"
	"github.com/pr02nl/medidor_ade/pkg/ade9000"
)

type CalibratePowerUseCase struct {
	ade9000           ade9000.ADE9000Interface
	medidorRepository entity.MedidorRepositoryInterface
}

func NewCalibratePowerUseCase(medidorRepository entity.MedidorRepositoryInterface, ade9000 ade9000.ADE9000Interface) *CalibratePowerUseCase {
	return &CalibratePowerUseCase{medidorRepository: medidorRepository, ade9000: ade9000}
}

func (u *CalibratePowerUseCase) Execute() error {
	println("Calibrating Power...")
	calibration := ade9000.NewCalibration(u.ade9000)
	time.Sleep(500 * time.Millisecond)
	if err := calibration.PGain_calibrate(1); err != nil {
		return err
	}
	return nil
}
