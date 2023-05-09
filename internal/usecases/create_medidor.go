package usecases

import (
	"github.com/google/uuid"
	"github.com/pr02nl/medidor_ade/internal/entity"
)

type CreateMedidorUseCase struct {
	MedidorRepository entity.MedidorRepositoryInterface
}

func NewCreateMedidorUseCase(medidorRepository entity.MedidorRepositoryInterface) *CreateMedidorUseCase {
	return &CreateMedidorUseCase{MedidorRepository: medidorRepository}
}

func (u *CreateMedidorUseCase) Execute(medidor *entity.Medidor) (*entity.Medidor, error) {
	println("Creating Medidor...")
	if medidor.ID == "" {
		medidor.ID = uuid.New().String()
	}
	err := u.MedidorRepository.Save(medidor)
	if err != nil {
		return nil, err
	}
	return medidor, nil
}
