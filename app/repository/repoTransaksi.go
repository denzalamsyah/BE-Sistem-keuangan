package repository

import (
	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type TransaksiRepository interface {
	Store(Transaksi *models.Transaksi) error
	Update(id int, Transaksi models.Transaksi) error
	Delete(id int) error
	GetList() ([]models.Transaksi, error)
}

type transaksiRepository struct {
	db *gorm.DB
}

func NewTransaksiRepo(db *gorm.DB) *transaksiRepository {
	return &transaksiRepository{db}
}

func (c *transaksiRepository) Store(Transaksi *models.Transaksi) error {
	err := c.db.Create(Transaksi).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *transaksiRepository) Update(id int, Transaksi models.Transaksi) error {
	err := c.db.Model(&Transaksi).Where("id = ?", id).Updates(&Transaksi).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *transaksiRepository) Delete(id int) error {
	var Transaksi models.Transaksi
	err := c.db.Where("id = ?", id).Delete(&Transaksi).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *transaksiRepository) GetList() ([]models.Transaksi, error) {
	var Transaksi []models.Transaksi
	err := c.db.Find(&Transaksi).Error
	if err != nil {
		return nil, err
	}
	return Transaksi, nil
}