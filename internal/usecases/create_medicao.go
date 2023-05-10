package usecases

import (
	"log"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/pr02nl/medidor_ade/internal/entity"
	"github.com/pr02nl/medidor_ade/pkg/ade9000"
)

type CreateMedicaoUseCase struct {
	Medidor           *entity.Medidor
	MedicaoRepository entity.MedicaoRepositoryInterface
	ade9000           ade9000.ADE9000Interface
}

func NewCreateMedicaoUseCase(medidor *entity.Medidor, medicaoRepository entity.MedicaoRepositoryInterface, ade9000 ade9000.ADE9000Interface) *CreateMedicaoUseCase {
	return &CreateMedicaoUseCase{MedicaoRepository: medicaoRepository, ade9000: ade9000, Medidor: medidor}
}

func (u *CreateMedicaoUseCase) Execute() error {
	log.Println("Creating Medicao...")
	voltageRms := ade9000.VoltageRMS{}
	currentRms := ade9000.CurrentRMS{}
	activePower := ade9000.Power{}
	reactivePower := ade9000.Power{}
	aparentPower := ade9000.Power{}
	voltageRmsRegs := ade9000.VoltageRMSRegs{}
	currentRmsRegs := ade9000.CurrentRMSRegs{}
	activePowerRegs := ade9000.ActivePowerRegs{}
	aparentPowerRegs := ade9000.ApparentPowerRegs{}
	reactivePowerRegs := ade9000.ReactivePowerRegs{}
	powerFactorRegs := ade9000.PowerFactorRegs{}
	periodRegs := ade9000.PeriodRegs{}
	angleRegs := ade9000.AngleRegs{}
	tempRegs := ade9000.TemperatureRegnValue{}
	if err := u.ade9000.ReadVoltageRMSRegs(&voltageRmsRegs); err != nil {
		return err
	}
	if err := u.ade9000.ReadCurrentRMSRegs(&currentRmsRegs); err != nil {
		return err
	}
	if err := u.ade9000.ReadActivePowerRegs(&activePowerRegs); err != nil {
		return err
	}
	if err := u.ade9000.ReadReactivePowerRegs(&reactivePowerRegs); err != nil {
		return err
	}
	if err := u.ade9000.ReadApparentPowerRegs(&aparentPowerRegs); err != nil {
		return err
	}
	if err := u.ade9000.ReadPowerFactorRegsnValues(&powerFactorRegs); err != nil {
		return err
	}
	if err := u.ade9000.ReadPeriodRegsnValues(&periodRegs); err != nil {
		return err
	}
	if err := u.ade9000.ReadAngleRegsnValues(&angleRegs); err != nil {
		return err
	}
	if err := u.ade9000.ReadTempRegnValue(&tempRegs); err != nil {
		return err
	}
	voltageRms.VoltageRMS_A = (float64(voltageRmsRegs.VoltageRMSReg_A) * ade9000.CAL_VRMS_CC) / math.Pow10(6)
	voltageRms.VoltageRMS_B = (float64(voltageRmsRegs.VoltageRMSReg_B) * ade9000.CAL_VRMS_CC) / math.Pow10(6)
	voltageRms.VoltageRMS_C = (float64(voltageRmsRegs.VoltageRMSReg_C) * ade9000.CAL_VRMS_CC) / math.Pow10(6)

	currentRms.CurrentRMS_A = (float64(currentRmsRegs.CurrentRMSReg_A) * ade9000.CAL_IRMS_CC) / math.Pow10(6)
	currentRms.CurrentRMS_B = (float64(currentRmsRegs.CurrentRMSReg_B) * ade9000.CAL_IRMS_CC) / math.Pow10(6)
	currentRms.CurrentRMS_C = (float64(currentRmsRegs.CurrentRMSReg_C) * ade9000.CAL_IRMS_CC) / math.Pow10(6)

	activePower.Power_A = (float64(activePowerRegs.ActivePowerReg_A) * ade9000.CAL_POWER_CC) / math.Pow10(3)
	activePower.Power_B = (float64(activePowerRegs.ActivePowerReg_B) * ade9000.CAL_POWER_CC) / math.Pow10(3)
	activePower.Power_C = (float64(activePowerRegs.ActivePowerReg_C) * ade9000.CAL_POWER_CC) / math.Pow10(3)

	reactivePower.Power_A = (float64(reactivePowerRegs.ReactivePowerReg_A) * ade9000.CAL_POWER_CC) / math.Pow10(3)
	reactivePower.Power_B = (float64(reactivePowerRegs.ReactivePowerReg_B) * ade9000.CAL_POWER_CC) / math.Pow10(3)
	reactivePower.Power_C = (float64(reactivePowerRegs.ReactivePowerReg_C) * ade9000.CAL_POWER_CC) / math.Pow10(3)

	aparentPower.Power_A = (float64(aparentPowerRegs.ApparentPowerReg_A) * ade9000.CAL_POWER_CC) / math.Pow10(3)
	aparentPower.Power_B = (float64(aparentPowerRegs.ApparentPowerReg_B) * ade9000.CAL_POWER_CC) / math.Pow10(3)
	aparentPower.Power_C = (float64(aparentPowerRegs.ApparentPowerReg_C) * ade9000.CAL_POWER_CC) / math.Pow10(3)

	medicao := entity.Medicao{
		ID:              uuid.New().String(),
		DateTime:        time.Now(),
		MedidorID:       u.Medidor.ID,
		VoltageRMS_A:    voltageRms.VoltageRMS_A,
		VoltageRMS_B:    voltageRms.VoltageRMS_B,
		VoltageRMS_C:    voltageRms.VoltageRMS_C,
		CurrentRMS_A:    currentRms.CurrentRMS_A,
		CurrentRMS_B:    currentRms.CurrentRMS_B,
		CurrentRMS_C:    currentRms.CurrentRMS_C,
		CurrentRMS_N:    currentRms.CurrentRMS_N,
		ActivePower_A:   activePower.Power_A,
		ActivePower_B:   activePower.Power_B,
		ActivePower_C:   activePower.Power_C,
		ReactivePower_A: reactivePower.Power_A,
		ReactivePower_B: reactivePower.Power_B,
		ReactivePower_C: reactivePower.Power_C,
		ApparentPower_A: aparentPower.Power_A,
		ApparentPower_B: aparentPower.Power_B,
		ApparentPower_C: aparentPower.Power_C,
		PowerFactor_A:   powerFactorRegs.PowerFactorValue_A,
		PowerFactor_B:   powerFactorRegs.PowerFactorValue_B,
		PowerFactor_C:   powerFactorRegs.PowerFactorValue_C,
		Frequency_A:     periodRegs.FrequencyValue_A,
		Frequency_B:     periodRegs.FrequencyValue_B,
		Frequency_C:     periodRegs.FrequencyValue_C,
		Angle_VA_VB:     angleRegs.AngleValue_VA_VB,
		Angle_VB_VC:     angleRegs.AngleValue_VB_VC,
		Angle_VA_VC:     angleRegs.AngleValue_VA_VC,
		Angle_VA_IA:     angleRegs.AngleValue_VA_IA,
		Angle_VB_IB:     angleRegs.AngleValue_VB_IB,
		Angle_VC_IC:     angleRegs.AngleValue_VC_IC,
		Angle_IA_IB:     angleRegs.AngleValue_IA_IB,
		Angle_IB_IC:     angleRegs.AngleValue_IB_IC,
		Angle_IA_IC:     angleRegs.AngleValue_IA_IC,
		Temperature:     tempRegs.Temperature_Val,
		Sincronizado:    false,
	}
	err := u.MedicaoRepository.Save(&medicao)
	if err != nil {
		return err
	}
	return nil
}
