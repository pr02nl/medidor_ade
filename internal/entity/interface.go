package entity

type MedicaoRepositoryInterface interface {
	Save(medicao *Medicao) error
	InitTable() error
}

type MedidorRepositoryInterface interface {
	Save(medidor *Medidor) error
	Update(medidor *Medidor) error
	InitTable() error
	Load() (*Medidor, error)
}
