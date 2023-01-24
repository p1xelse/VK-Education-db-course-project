package usecase

import (
	"github.com/p1xelse/VK_DB_course_project/app/internal/service/repository"
	"github.com/p1xelse/VK_DB_course_project/app/models"
)

type ServiceUseCaseI interface {
	ClearData() error
	GetStatus() (*models.ServiceStatus, error)
}

type serviceUsecase struct {
	serviceRepo repository.RepositoryI
}

func (s serviceUsecase) ClearData() error {
	err := s.serviceRepo.ClearData()
	return err
}

func (s serviceUsecase) GetStatus() (*models.ServiceStatus, error) {
	status, err := s.serviceRepo.GetStatus()
	if err != nil {
		return nil, err
	}

	return status, nil
}

func NewServiceUsecase(ps repository.RepositoryI) ServiceUseCaseI {
	return &serviceUsecase{
		serviceRepo: ps,
	}
}
