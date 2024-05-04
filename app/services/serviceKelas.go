package services

import (
	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/repository"
)

type KelasService interface {
	Store(Kelas *models.Kelas) error
	Update(kode int, Kelas models.Kelas) error
	Delete(kode int) error
	GetList(page, pageSize int) ([]models.Kelas, int, error)
	GetTotalKelasCount() (int, error)
	Search(nama string) ([]models.Kelas, error)
	GetKode(kode int) (models.Kelas, error)

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

func (s *kelasServices) Update(kode int, Kelas models.Kelas) error {
	err := s.kelasRepo.Update(kode, Kelas)
	if err != nil {
		return err
	}
	return nil
}

func (s *kelasServices) Delete(kode int) error {
	err := s.kelasRepo.Delete(kode)
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

func (s *kelasServices) GetTotalKelasCount() (int, error) {
	count, err := s.kelasRepo.GetTotalKelasCount()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func(s *kelasServices) Search(nama string) ([]models.Kelas, error){
	kelas, err := s.kelasRepo.Search(nama)

	if err != nil {
        return nil, err
    }
	return kelas, nil
}
func (s *kelasServices) GetKode(kode int) (models.Kelas, error) {
	return s.kelasRepo.GetKode(kode)
}
