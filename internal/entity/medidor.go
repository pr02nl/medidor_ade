package entity

type Medidor struct {
	ID                      string  `json:"id"`
	NominalVoltage          float64 `json:"nominal_voltage"`
	NominalCurrent          float64 `json:"nominal_current"`
	CurrentTransformerRatio float64 `json:"current_transformer_ratio"`
	Frequency               float64 `json:"frequency"`
	CurrentTransfer         float64 `json:"current_transfer"`
	VoltageTransfer         float64 `json:"voltage_transfer"`
	CalIrmsCC               float64 `json:"cal_irms_cc"`
	CalVrmsCC               float64 `json:"cal_vrms_cc"`
	CalPwrCC                float64 `json:"cal_pwr_cc"`
	CalEnergyCC             float64 `json:"cal_energy_ac"`
	CalibratedVoltage       bool    `json:"calibrated_voltage"`
	CalibratedCurrent       bool    `json:"calibrated_current"`
	CalibratedPower         bool    `json:"calibrated_power"`
	CalibratedPhase         bool    `json:"calibrated_phase"`
}
