package entity

import "time"

type Medicao struct {
	ID              string    `json:"id"`
	MedidorID       string    `json:"medidor_id"`
	DateTime        time.Time `json:"datetime"`
	VoltageRMS_A    float64   `json:"voltage_rms_a"`
	VoltageRMS_B    float64   `json:"voltage_rms_b"`
	VoltageRMS_C    float64   `json:"voltage_rms_c"`
	CurrentRMS_A    float64   `json:"current_rms_a"`
	CurrentRMS_B    float64   `json:"current_rms_b"`
	CurrentRMS_C    float64   `json:"current_rms_c"`
	CurrentRMS_N    float64   `json:"current_rms_n"`
	ActivePower_A   float64   `json:"active_power_a"`
	ActivePower_B   float64   `json:"active_power_b"`
	ActivePower_C   float64   `json:"active_power_c"`
	ReactivePower_A float64   `json:"reactive_power_a"`
	ReactivePower_B float64   `json:"reactive_power_b"`
	ReactivePower_C float64   `json:"reactive_power_c"`
	ApparentPower_A float64   `json:"apparent_power_a"`
	ApparentPower_B float64   `json:"apparent_power_b"`
	ApparentPower_C float64   `json:"apparent_power_c"`
	PowerFactor_A   float32   `json:"power_factor_a"`
	PowerFactor_B   float32   `json:"power_factor_b"`
	PowerFactor_C   float32   `json:"power_factor_c"`
	Frequency_A     float32   `json:"frequency_a"`
	Frequency_B     float32   `json:"frequency_b"`
	Frequency_C     float32   `json:"frequency_c"`
	Angle_VA_VB     float32   `json:"angle_va_vb"`
	Angle_VB_VC     float32   `json:"angle_vb_vc"`
	Angle_VA_VC     float32   `json:"angle_va_vc"`
	Angle_VA_IA     float32   `json:"angle_va_ia"`
	Angle_VB_IB     float32   `json:"angle_vb_ib"`
	Angle_VC_IC     float32   `json:"angle_vc_ic"`
	Angle_IA_IB     float32   `json:"angle_ia_ib"`
	Angle_IB_IC     float32   `json:"angle_ib_ic"`
	Angle_IA_IC     float32   `json:"angle_ia_ic"`
	Temperature     float32   `json:"temperature"`
	Sincronizado    bool      `json:"sincronizado"`
}
