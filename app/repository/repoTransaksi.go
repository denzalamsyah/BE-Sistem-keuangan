package repository

import (
	"math"

	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type TransaksiRepository interface {
	Store(Transaksi *models.Transaksi) error
	Update(id int, Transaksi models.Transaksi) error
	Delete(id int) error
	GetList(page, paedSize int) ([]models.Transaksi, int, error)
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

func (c *transaksiRepository) GetList(page, pageSize int) ([]models.Transaksi, int, error) {
	var Transaksi []models.Transaksi
	offset := (page - 1) * pageSize

	err := c.db.Limit(pageSize).Offset(offset).Find(&Transaksi).Error
	if err != nil {
		return nil, 0, err
	}

	var totalData int64
	if err := c.db.Model(&models.Transaksi{}).Count(&totalData).Error; err != nil {
		return nil, 0, err
	}
	totalPage := int(math.Ceil(float64(totalData) / float64(pageSize)))

	return Transaksi, totalPage, nil
}