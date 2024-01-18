package repository

import (
	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type JurusanRepository interface {
	Store(Jurusan *models.Jurusan) error
	Update(id int, Jurusan models.Jurusan) error
	Delete(id int) error
	GetList() ([]models.Jurusan, error)
}

type jurusanRepository struct {
	db *gorm.DB
}

func NewJurusanRepo(db *gorm.DB) *jurusanRepository {
	return &jurusanRepository{db}
}

func (c *jurusanRepository) Store(Jurusan *models.Jurusan) error {
	err := c.db.Create(Jurusan).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *jurusanRepository) Update(id int, Jurusan models.Jurusan) error {
	err := c.db.Model(&Jurusan).Where("id_jurusan = ?", id).Updates(&Jurusan).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *jurusanRepository) Delete(id int) error {
	err := c.db.Where("id_jurusan = ?", id).Delete(&models.Jurusan{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *jurusanRepository) GetList() ([]models.Jurusan, error) {
	var Jurusan []models.Jurusan
	err := c.db.Find(&Jurusan).Error
	if err != nil {
		return nil, err
	}
	return Jurusan, nil
}