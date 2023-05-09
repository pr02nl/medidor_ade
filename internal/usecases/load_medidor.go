package usecases

import "github.com/pr02nl/medidor_ade/internal/entity"

type LoadMedidorUseCase struct {
	MedidorRepository entity.MedidorRepositoryInterface
}

func NewLoadMedidorUseCase(medidorRepository entity.MedidorRepositoryInterface) *LoadMedidorUseCase {
	return &LoadMedidorUseCase{MedidorRepository: medidorRepository}
}

func (u *LoadMedidorUseCase) Execute() (*entity.Medidor, error) {
	println("Loading Medidor...")
	medidor, err := u.MedidorRepository.Load()
	if err != nil {
		return nil, err
	}
	return medidor, nil
}
