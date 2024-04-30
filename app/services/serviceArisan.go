package services

import (
	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/repository"
)

type ArisanService interface {
	Store(Arisan *models.Arisan) error
	Update(id int, Arisan models.Arisan) error
	GetList(page, pageSize int) ([]models.Arisan,int, error)
	Delete(id int) error
	Search(nama, tanggal string) ([]models.Arisan, error)
	
}

type arisanService struct {
	arisanRepo repository.ArisanRepository
}

func NewArisanService(arisanRepo repository.ArisanRepository) ArisanService {
	return &arisanService{arisanRepo}
	
}

func (s *arisanService) Store(Arisan *models.Arisan) error{
	err := s.arisanRepo.Store(Arisan)
	if err != nil {
		return err
	}
	return nil
}

func (s *arisanService) Update(id int, Arisan models.Arisan) error{
	err := s.arisanRepo.Update(id, Arisan)
	if err != nil {
		return err
	}
	return nil
}
func (s *arisanService) GetList(page, pageSize int) ([]models.Arisan,int, error){
	arisans, totalPage, err := s.arisanRepo.GetList(page, pageSize)
	if err != nil {
		return nil,0, err
	}
	return arisans, totalPage, nil
}
func (s *arisanService) Delete(id int) error{
	err := s.arisanRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
func (s *arisanService) Search(nama, tanggal string) ([]models.Arisan, error){
	arisan, err := s.arisanRepo.Search(nama, tanggal)

	if err != nil {
        return nil, err
    }
	return arisan, nil
}