package usecases

import "github.com/pr02nl/medidor_ade/internal/entity"

type UpdateMedidorUseCase struct {
	MedidorRepository entity.MedidorRepositoryInterface
}

func NewUpdateMedidorUseCase(medidorRepository entity.MedidorRepositoryInterface) *UpdateMedidorUseCase {
	return &UpdateMedidorUseCase{MedidorRepository: medidorRepository}
}

func (u *UpdateMedidorUseCase) Execute(medidor *entity.Medidor) (*entity.Medidor, error) {
	println("Updating Medidor...")
	err := u.MedidorRepository.Update(medidor)
	if err != nil {
		return nil, err
	}
	return medidor, nil
}
