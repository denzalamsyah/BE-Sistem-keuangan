package repository

import (
	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type PengeluaranRepository interface {
	Store(pengeluaran *models.Pengeluaran) error
	Update(id int, pengeluaran models.Pengeluaran) error
	Delete(id int) error
	GetByID(id int) (*models.Pengeluaran, error)
	GetList() ([]models.Pengeluaran, error)
}

type pengeluaranRepository struct {
	db *gorm.DB
}

func NewPengeluaranRepo(db *gorm.DB) *pengeluaranRepository {
	return &pengeluaranRepository{db}
}

func (c *pengeluaranRepository) Store(pengeluaran *models.Pengeluaran) error{
	err := c.db.Create(pengeluaran).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *pengeluaranRepository) Update(id int, pengeluaran models.Pengeluaran) error{
	err := c.db.Model(&pengeluaran).Where("id = ?", id).Updates(&pengeluaran).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *pengeluaranRepository) Delete(id int) error {
	var Pengeluaran models.Pengeluaran
	err := c.db.Where("id = ?", id).Delete(&Pengeluaran).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *pengeluaranRepository) GetByID(id int) (*models.Pengeluaran, error){
    var pengeluaran models.Pengeluaran

    err := c.db.Where("id = ?", id).First(&pengeluaran).Error

    if err != nil {
		return nil, err
	}

    return &pengeluaran, nil
}

func (c *pengeluaranRepository) GetList() ([]models.Pengeluaran, error){
    var pengeluaran []models.Pengeluaran

    err := c.db.Find(&pengeluaran).Error
    if err != nil {
		return nil, err
	}
    return pengeluaran, nil
}