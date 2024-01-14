package repository

import (
	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type SemesterRepository interface {
	Store(PembayaranSemester *models.PembayaranSemester) error
	Update(id int, PembayaranSemester models.PembayaranSemester) error
	Delete(id int) error
	GetByID(id int) (*models.PembayaranSemester, error)
	GetList() ([]models.PembayaranSemester, error)
}

type semesterRepository struct{
	db *gorm.DB
}

func NewSemesterRepo(db *gorm.DB) *semesterRepository{
	return &semesterRepository{db}
}

func (c *semesterRepository) Store(PembayaranSemester *models.PembayaranSemester) error {
	err := c.db.Create(PembayaranSemester).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *semesterRepository) Update(id int, PembayaranSemester models.PembayaranSemester) error {
	err := c.db.Model(&PembayaranSemester).Where("id = ?", id).Updates(&PembayaranSemester).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *semesterRepository) Delete(id int) error {
	var PembayaranSemester models.PembayaranSemester
	err := c.db.Where("id = ?", id).Delete(&PembayaranSemester).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *semesterRepository) GetByID(id int) (*models.PembayaranSemester, error) {
	var PembayaranSemester models.PembayaranSemester
	err := c.db.Where("id = ?", id).Find(&PembayaranSemester).Error
	if err != nil {
		return nil, err
	}
	return &PembayaranSemester, nil
}

func (c *semesterRepository) GetList() ([]models.PembayaranSemester, error) {
	var PembayaranSemester []models.PembayaranSemester
	err := c.db.Find(&PembayaranSemester).Error
	if err != nil {
		return nil, err
	}
	return PembayaranSemester, nil
}