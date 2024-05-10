package services

import (
	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/repository"
)

type JurusanService interface {
	Store(Jurusan *models.Jurusan) error
	Update(kode string, Jurusan models.Jurusan) error
	Delete(kode string) error
	GetList(page, pageSize int) ([]models.Jurusan,int, error)
	GetTotalJurusanCount() (int, error)
	Search(nama string) ([]models.Jurusan, error)
	GetKode(kode string) (models.Jurusan, error)
}

type jurusanService struct {
	jurusanRepo repository.JurusanRepository
}

func NewJurusanService(jurusanRepo repository.JurusanRepository) JurusanService {
	return &jurusanService{jurusanRepo}
}

func (s *jurusanService) Store(jurusan *models.Jurusan) error {
	err := s.jurusanRepo.Store(jurusan)
	if err != nil {
		return err
	}
	return nil
}

func (s *jurusanService) Update(kode string, jurusan models.Jurusan) error {
	err := s.jurusanRepo.Update(kode, jurusan)
	if err != nil {
		return err
	}
	return nil
}

func (s *jurusanService) Delete(kode string) error {
	err := s.jurusanRepo.Delete(kode)
	if err != nil {
		return err
	}
	return nil
}

func (s *jurusanService) GetList(page, pageSize int) ([]models.Jurusan,int, error) {
	jurusans, totalPage, err := s.jurusanRepo.GetList(page, pageSize)
	if err != nil {
		return nil,0, err
	}
	return jurusans, totalPage, nil
}

func (s *jurusanService) GetTotalJurusanCount() (int, error) {
	count, err := s.jurusanRepo.GetTotalJurusanCount()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *jurusanService) Search(nama string) ([]models.Jurusan, error){
	jurusan, err := s.jurusanRepo.Search(nama)

	if err != nil {
        return nil, err
    }
	return jurusan, nil
}


func (c *jurusanService) GetKode(kode string) (models.Jurusan, error) {
	return c.jurusanRepo.GetKode(kode)
}