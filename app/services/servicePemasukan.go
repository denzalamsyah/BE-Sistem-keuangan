package services

import (
	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/repository"
)

type PemasukanService interface {
	FindAll(page, pageSize int) ([]models.PemasukanResponse, int, error)
	TotalKeuangan() (int, int, int, error)
	Store(Pemasukanlainnya *models.Pemasukanlainnya) error
	Update(id int, Pemasukanlainnya models.Pemasukanlainnya) error
	Delete(id int) error
	GetByID(id int) (*models.Pemasukanlainnya, error)
	GetList(page, pageSize int) ([]models.Pemasukanlainnya, int, error)

	
}
type pemasukanService struct {
	pemasukanRepository repository.PemasukanRepository
}

func NewPemasukanService(pemasukanRepository repository.PemasukanRepository) PemasukanService {
	return &pemasukanService{pemasukanRepository}
}

func (c *pemasukanService) FindAll(page, pageSize int) ([]models.PemasukanResponse, int, error) {
	pemasukan, totalPage, err := c.pemasukanRepository.FindAll(page, pageSize)
	if err != nil {
		return nil,0, err
	}
	return pemasukan, totalPage, nil
}

func (c *pemasukanService) Store(Pemasukanlainnya *models.Pemasukanlainnya) error {
	err := c.pemasukanRepository.Store(Pemasukanlainnya)
	if err != nil {
		return err
	}
	return nil
}

func (c *pemasukanService) Update(id int, Pemasukanlainnya models.Pemasukanlainnya) error {
	err := c.pemasukanRepository.Update(id, Pemasukanlainnya)
	if err != nil {
		
		return err
	}
	return nil
}

func (c *pemasukanService) Delete(id int) error {
	err := c.pemasukanRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *pemasukanService) GetByID(id int) (*models.Pemasukanlainnya, error) {
	Pemasukanlainnya, err := c.pemasukanRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return Pemasukanlainnya, nil
}

func (c *pemasukanService) GetList(page, pageSize int) ([]models.Pemasukanlainnya, int,error) {
	Pemasukanlainnyas, totalPage, err := c.pemasukanRepository.GetList(page, pageSize)
	if err != nil {
		return nil, 0,err
	}
	return Pemasukanlainnyas, totalPage, nil
}

func (c *pemasukanService) TotalKeuangan() (int, int, int, error) {
    saldo, totalPemasukan, totalPengeluaran, err := c.pemasukanRepository.TotalKeuangan()
    if err != nil {
        return 0, 0, 0, err
    }
    return saldo, totalPemasukan, totalPengeluaran, nil
}


