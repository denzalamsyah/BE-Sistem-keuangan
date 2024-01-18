package repository

import (
	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type KelasRepository interface {
	Store(Kelas *models.Kelas) error
	Update(id int, Kelas models.Kelas) error
	Delete(id int) error
	GetList() ([]models.Kelas, error)
}

type kelasRepository struct{
	db *gorm.DB
}

func NewKelasRepo(db *gorm.DB) *kelasRepository {
	return &kelasRepository{db}
}

func (c *kelasRepository) Store(Kelas *models.Kelas) error {
	err := c.db.Create(Kelas).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *kelasRepository) Update(id int, Kelas models.Kelas) error {
	err := c.db.Model(&Kelas).Where("id_kelas = ?", id).Updates(&Kelas).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *kelasRepository) Delete(id int) error {
	var Kelas models.Kelas
	err := c.db.Where("id_kelas = ?", id).Delete(&Kelas).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *kelasRepository) GetList() ([]models.Kelas, error) {
	var Kelas []models.Kelas
	err := c.db.Find(&Kelas).Error
	if err != nil {
		return nil, err
	}
	return Kelas, nil
}