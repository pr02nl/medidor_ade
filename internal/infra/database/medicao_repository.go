package database

import (
	"database/sql"

	"github.com/pr02nl/medidor_ade/internal/entity"
)

type MedicaoRepository struct {
	Db *sql.DB
}

func NewMedicaoRepository(db *sql.DB) *MedicaoRepository {
	return &MedicaoRepository{Db: db}
}

func (r *MedicaoRepository) Save(medicao *entity.Medicao) error {
	stmt, err := r.Db.Prepare(`
		INSERT INTO medicao (id, datetime, voltage_rms_a, voltage_rms_b, voltage_rms_c, current_rms_a, current_rms_b, current_rms_c, current_rms_n, active_power_a, active_power_b, active_power_c, reactive_power_a, reactive_power_b, reactive_power_c, apparent_power_a, apparent_power_b, apparent_power_c, power_factor_a, power_factor_b, power_factor_c, frequency_a, frequency_b, frequency_c, angle_va_vb, angle_vb_vc, angle_va_vc, angle_va_ia, angle_vb_ib, angle_vc_ic, angle_ia_ib, angle_ib_ic, angle_ia_ic, temperature)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(medicao.ID, medicao.DateTime, medicao.VoltageRMS_A, medicao.VoltageRMS_B, medicao.VoltageRMS_C, medicao.CurrentRMS_A, medicao.CurrentRMS_B, medicao.CurrentRMS_C, medicao.CurrentRMS_N, medicao.ActivePower_A, medicao.ActivePower_B, medicao.ActivePower_C, medicao.ReactivePower_A, medicao.ReactivePower_B, medicao.ReactivePower_C, medicao.ApparentPower_A, medicao.ApparentPower_B, medicao.ApparentPower_C, medicao.PowerFactor_A, medicao.PowerFactor_B, medicao.PowerFactor_C, medicao.Frequency_A, medicao.Frequency_B, medicao.Frequency_C, medicao.Angle_VA_VB, medicao.Angle_VB_VC, medicao.Angle_VA_VC, medicao.Angle_VA_IA, medicao.Angle_VB_IB, medicao.Angle_VC_IC, medicao.Angle_IA_IB, medicao.Angle_IB_IC, medicao.Angle_IA_IC, medicao.Temperature)
	if err != nil {
		return err
	}
	return nil
}

func (r *MedicaoRepository) InitTable() error {
	stmt, err := r.Db.Prepare(`
		CREATE TABLE IF NOT EXISTS medicao (
			id VARCHAR(36) PRIMARY KEY,
			datetime DATETIME,
			voltage_rms_a FLOAT,
			voltage_rms_b FLOAT,
			voltage_rms_c FLOAT,
			current_rms_a FLOAT,
			current_rms_b FLOAT,
			current_rms_c FLOAT,
			current_rms_n FLOAT,
			active_power_a FLOAT,
			active_power_b FLOAT,
			active_power_c FLOAT,
			reactive_power_a FLOAT,
			reactive_power_b FLOAT,
			reactive_power_c FLOAT,
			apparent_power_a FLOAT,
			apparent_power_b FLOAT,
			apparent_power_c FLOAT,
			power_factor_a FLOAT,
			power_factor_b FLOAT,
			power_factor_c FLOAT,
			frequency_a FLOAT,
			frequency_b FLOAT,
			frequency_c FLOAT,
			angle_va_vb FLOAT,
			angle_vb_vc FLOAT,
			angle_va_vc FLOAT,
			angle_va_ia FLOAT,
			angle_vb_ib FLOAT,
			angle_vc_ic FLOAT,
			angle_ia_ib FLOAT,
			angle_ib_ic FLOAT,
			angle_ia_ic FLOAT,
			temperature FLOAT
		)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}
