package services

import (
	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/repository"
)

type KelasService interface {
	Store(Kelas *models.Kelas) error
	Update(id int, Kelas models.Kelas) error
	Delete(id int) error
	GetList(page, pageSize int) ([]models.Kelas, int, error)
}

type kelasServices struct {
	kelasRepo repository.KelasRepository
}

func NewKelasService(kelasRepo repository.KelasRepository) KelasService {
	return &kelasServices{kelasRepo}
}

func (s *kelasServices) Store(Kelas *models.Kelas) error {
	err := s.kelasRepo.Store(Kelas)
	if err != nil {
		return err
	}
	return nil
}

func (s *kelasServices) Update(id int, Kelas models.Kelas) error {
	err := s.kelasRepo.Update(id, Kelas)
	if err != nil {
		return err
	}
	return nil
}

func (s *kelasServices) Delete(id int) error {
	err := s.kelasRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *kelasServices) GetList(page, pageSize int) ([]models.Kelas, int, error) {
	Kelas, totalPage, err := s.kelasRepo.GetList(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return Kelas, totalPage, nil
}