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
	GetList(page, pageSize int) ([]models.Pengeluaran,int, error)
	Search(nama, tanggal string) ([]models.Pengeluaran, int, error)
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

func (c *pengeluaranService) GetList(page, pageSize int) ([]models.Pengeluaran,int, error) {
	Pengeluarans, totalPage, err := c.pengeluaranRepository.GetList(page, pageSize)
	if err != nil {
		return nil, 0,err
	}
	return Pengeluarans, totalPage, nil
}

func (c *pengeluaranService) Search(nama, tanggal string) ([]models.Pengeluaran, int, error){
	pengeluaran,total, err := c.pengeluaranRepository.Search(nama, tanggal)
	if err != nil {
        return nil,0, err
    }
	return pengeluaran,total, nil
}