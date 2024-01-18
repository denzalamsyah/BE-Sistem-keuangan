package services

import (
	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/repository"
)

type PemasukanService interface {
	FindAll() ([]models.PemasukanResponse, error)
	TotalKeuangan() (int, int, int, error)
	Store(Pemasukanlainnya *models.Pemasukanlainnya) error
	Update(id int, Pemasukanlainnya models.Pemasukanlainnya) error
	Delete(id int) error
	GetByID(id int) (*models.Pemasukanlainnya, error)
	GetList() ([]models.Pemasukanlainnya, error)

	
}
type pemasukanService struct {
	pemasukanRepository repository.PemasukanRepository
}

func NewPemasukanService(pemasukanRepository repository.PemasukanRepository) PemasukanService {
	return &pemasukanService{pemasukanRepository}
}

func (c *pemasukanService) FindAll() ([]models.PemasukanResponse, error) {
	pemasukan, err := c.pemasukanRepository.FindAll()
	if err != nil {
		return nil, err
	}
	return pemasukan, nil
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

func (c *pemasukanService) GetList() ([]models.Pemasukanlainnya, error) {
	Pemasukanlainnyas, err := c.pemasukanRepository.GetList()
	if err != nil {
		return nil, err
	}
	return Pemasukanlainnyas, nil
}

func (c *pemasukanService) TotalKeuangan() (int, int, int, error) {
    saldo, totalPemasukan, totalPengeluaran, err := c.pemasukanRepository.TotalKeuangan()
    if err != nil {
        return 0, 0, 0, err
    }
    return saldo, totalPemasukan, totalPengeluaran, nil
}


