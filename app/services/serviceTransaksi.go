package services

import (
	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/repository"
)

type TransaksiService interface {
	Store(Transaksi *models.Transaksi) error
	Update(id int, Transaksi models.Transaksi) error
	Delete(id int) error
	GetList(page, pageSize int) ([]models.Transaksi, int, error)
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

func (c *transaksiService) GetList(page, pageSize int) ([]models.Transaksi, int, error) {

	Transaksi, totalPage, err := c.transaksiRepo.GetList(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return Transaksi,totalPage, nil
}