package services

import (
	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/repository"
)

type PengeluaranService interface {
	Store(Pengeluaran *models.Pengeluaran) error
	Update(id int, Pengeluaran models.Pengeluaran) error
	Delete(id int) error
	GetByID(id int) (*models.Pengeluaran, error)
	GetList() ([]models.Pengeluaran, error)
}

type pengeluaranService struct {
	pengeluaranRepository repository.PengeluaranRepository
}

func NewPengeluaranService(pengeluaranRepository repository.PengeluaranRepository) PengeluaranService {
	return &pengeluaranService{pengeluaranRepository}
}

func (c *pengeluaranService) Store(Pengeluaran *models.Pengeluaran) error {
	err := c.pengeluaranRepository.Store(Pengeluaran)
	if err != nil {
		return err
	}
	return nil
}

func (c *pengeluaranService) Update(id int, Pengeluaran models.Pengeluaran) error {
	err := c.pengeluaranRepository.Update(id, Pengeluaran)
	if err != nil {
		
		return err
	}
	return nil
}

func (c *pengeluaranService) Delete(id int) error {
	err := c.pengeluaranRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *pengeluaranService) GetByID(id int) (*models.Pengeluaran, error) {
	Pengeluaran, err := c.pengeluaranRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return Pengeluaran, nil
}

func (c *pengeluaranService) GetList() ([]models.Pengeluaran, error) {
	Pengeluarans, err := c.pengeluaranRepository.GetList()
	if err != nil {
		return nil, err
	}
	return Pengeluarans, nil
}