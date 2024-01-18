package services

import (
	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/repository"
)

type TransaksiService interface {
	Store(Transaksi *models.Transaksi) error
	Update(id int, Transaksi models.Transaksi) error
	Delete(id int) error
	GetList() ([]models.Transaksi, error)
}
type transaksiService struct {
	transaksiRepo repository.TransaksiRepository
}

func NewTransaksiService(transaksiRepo repository.TransaksiRepository) TransaksiService {
	return &transaksiService{transaksiRepo}
}

func (c *transaksiService) Store(Transaksi *models.Transaksi) error {
	err := c.transaksiRepo.Store(Transaksi)
	if err != nil {
		return err
	}
	return nil
}

func (c *transaksiService) Update(id int, Transaksi models.Transaksi) error {
	err := c.transaksiRepo.Update(id, Transaksi)
	if err != nil {
		return err
	}
	return nil
}

func (c *transaksiService) Delete(id int) error {
	err := c.transaksiRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *transaksiService) GetList() ([]models.Transaksi, error) {

	Transaksi, err := c.transaksiRepo.GetList()
	if err != nil {
		return nil, err
	}
	return Transaksi, nil
}