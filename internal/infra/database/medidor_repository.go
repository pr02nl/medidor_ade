package database

import (
	"database/sql"

	"github.com/pr02nl/medidor_ade/internal/entity"
)

type MedidorRepository struct {
	Db *sql.DB
}

func NewMedidorRepository(db *sql.DB) *MedidorRepository {
	return &MedidorRepository{Db: db}
}

func (r *MedidorRepository) Save(medidor *entity.Medidor) error {
	stmt, err := r.Db.Prepare(`
		INSERT INTO medidor (id, nominalVoltage, nominalCurrent, currentTransformerRatio, frequency, currentTransfer, 
			voltageTransfer, calIrmsCC, calVrmsCC, calPwrCC, calEnergyCC, calibratedVoltage, calibratedCurrent, calibratedPower, calibratedPhase)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,	?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(medidor.ID, medidor.NominalVoltage, medidor.NominalCurrent, medidor.CurrentTransformerRatio, medidor.Frequency, medidor.CurrentTransfer, medidor.VoltageTransfer, medidor.CalIrmsCC, medidor.CalVrmsCC, medidor.CalPwrCC, medidor.CalEnergyCC, medidor.CalibratedVoltage, medidor.CalibratedCurrent, medidor.CalibratedPower, medidor.CalibratedPhase)
	if err != nil {
		return err
	}
	return nil
}

func (r *MedidorRepository) InitTable() error {
	stmt, err := r.Db.Prepare(`
		CREATE TABLE IF NOT EXISTS medidor (
			id VARCHAR(36) PRIMARY KEY,
			nominalVoltage FLOAT,
			nominalCurrent FLOAT,
			currentTransformerRatio FLOAT,
			frequency FLOAT,
			currentTransfer FLOAT,
			voltageTransfer FLOAT,
			calIrmsCC FLOAT,
			calVrmsCC FLOAT,
			calPwrCC FLOAT,
			calEnergyCC FLOAT,
			calibratedVoltage BOOLEAN,
			calibratedCurrent BOOLEAN,
			calibratedPower BOOLEAN,
			calibratedPhase BOOLEAN
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

func (r *MedidorRepository) Load() (*entity.Medidor, error) {
	stmt, err := r.Db.Prepare(`
		SELECT id, nominalVoltage, nominalCurrent, currentTransformerRatio, frequency, 
				currentTransfer, voltageTransfer, calIrmsCC, calVrmsCC, calPwrCC, calEnergyCC, 
				calibratedVoltage, calibratedCurrent, calibratedPower, calibratedPhase FROM medidor limit 1
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow()
	var medidor entity.Medidor
	err = row.Scan(&medidor.ID, &medidor.NominalVoltage, &medidor.NominalCurrent, &medidor.CurrentTransformerRatio,
		&medidor.Frequency, &medidor.CurrentTransfer, &medidor.VoltageTransfer, &medidor.CalIrmsCC,
		&medidor.CalVrmsCC, &medidor.CalPwrCC, &medidor.CalEnergyCC, &medidor.CalibratedVoltage,
		&medidor.CalibratedCurrent, &medidor.CalibratedPower, &medidor.CalibratedPhase)
	if err != nil {
		return nil, err
	}
	return &medidor, nil
}

func (r *MedidorRepository) Update(medidor *entity.Medidor) error {
	stmt, err := r.Db.Prepare(`
		UPDATE medidor SET nominalVoltage=?, nominalCurrent=?, currentTransformerRatio=?, frequency=?, 
			currentTransfer=?, voltageTransfer=?, calIrmsCC=?, calVrmsCC=?, calPwrCC=?, calEnergyCC=?, 
			calibratedVoltage=?, calibratedCurrent=?, calibratedPower=?, calibratedPhase=? WHERE id=?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(medidor.NominalVoltage, medidor.NominalCurrent, medidor.CurrentTransformerRatio, medidor.Frequency, medidor.CurrentTransfer, medidor.VoltageTransfer, medidor.CalIrmsCC, medidor.CalVrmsCC, medidor.CalPwrCC, medidor.CalEnergyCC, medidor.CalibratedVoltage, medidor.CalibratedCurrent, medidor.CalibratedPower, medidor.CalibratedPhase, medidor.ID)
	if err != nil {
		return err
	}
	return nil
}
