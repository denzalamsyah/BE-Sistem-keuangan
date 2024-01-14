package repository

import (
	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type SiswaRepository interface {
	Store(Siswa *models.Siswa) error
	Update(id int, Siswa models.Siswa) error
	Delete(id int) error
	GetByID(id int) (*models.Siswa, error)
	GetList() ([]models.Siswa, error)
}

type siswaRepository struct{
	db *gorm.DB
}

func NewSiswaRepo(db *gorm.DB) *siswaRepository {
	return &siswaRepository{db}
}
func (c *siswaRepository) Store(Siswa *models.Siswa) error {
	err := c.db.Create(Siswa).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *siswaRepository) Update(id int, Siswa models.Siswa) error {
	err := c.db.Model(&Siswa).Where("id = ?", id).Updates(&Siswa).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *siswaRepository) Delete(id int) error {
	var Siswa models.Siswa
	err := c.db.Where("id = ?", id).Delete(&Siswa).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *siswaRepository) GetByID(id int) (*models.Siswa, error) {
	var Siswa models.Siswa
	err := c.db.Where("id = ?", id).Find(&Siswa).Error
	if err != nil {
		return nil, err
	}
	return &Siswa, nil
}

func (c *siswaRepository) GetList() ([]models.Siswa, error) {
	var Siswa []models.Siswa
	err := c.db.Find(&Siswa).Error
	if err != nil {
		return nil, err
	}
	return Siswa, nil
}