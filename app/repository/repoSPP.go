package repository

import (
	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type SPPRepository interface {
	Store(PembayaranSPP *models.PembayaranSPP) error
	Update(id int, PembayaranSPP models.PembayaranSPP) error
	Delete(id int) error
	GetByID(id int) (*models.PembayaranSPP, error)
	GetList() ([]models.PembayaranSPP, error)
}
type sppRepository struct{
	db *gorm.DB
}

func NewSPPRepo(db *gorm.DB) *sppRepository{
	return &sppRepository{db}
}

func (c *sppRepository) Store(PembayaranSPP *models.PembayaranSPP) error {
	err := c.db.Create(PembayaranSPP).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *sppRepository) Update(id int, PembayaranSPP models.PembayaranSPP) error {
	err := c.db.Model(&PembayaranSPP).Where("id = ?", id).Updates(&PembayaranSPP).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *sppRepository) Delete(id int) error {
	var PembayaranSPP models.PembayaranSPP
	err := c.db.Where("id = ?", id).Delete(&PembayaranSPP).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *sppRepository) GetByID(id int) (*models.PembayaranSPP, error) {
	var PembayaranSPP models.PembayaranSPP
	err := c.db.Where("id = ?", id).Find(&PembayaranSPP).Error
	if err != nil {
		return nil, err
	}
	return &PembayaranSPP, nil
}

func (c *sppRepository) GetList() ([]models.PembayaranSPP, error) {
	var PembayaranSPP []models.PembayaranSPP
	err := c.db.Find(&PembayaranSPP).Error
	if err != nil {
		return nil, err
	}
	return PembayaranSPP, nil
}
