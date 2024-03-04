package repository

import (
	"math"

	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type PengeluaranRepository interface {
	Store(pengeluaran *models.Pengeluaran) error
	Update(id int, pengeluaran models.Pengeluaran) error
	Delete(id int) error
	GetByID(id int) (*models.Pengeluaran, error)
	GetList(page, pageSize int) ([]models.Pengeluaran, int, error)
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

func (c *pengeluaranRepository) GetList(page, pageSize int) ([]models.Pengeluaran,int, error){
    var pengeluaran []models.Pengeluaran
	offset := (page - 1) * pageSize

	err := c.db.Limit(pageSize).Offset(offset).Find(&pengeluaran).Error
    if err != nil {
		return nil, 0,err
	}

	var totalData int64
	if err := c.db.Model(&models.Pengeluaran{}).Count(&totalData).Error; err != nil {
		return nil, 0, err
	}


	totalPage := int(math.Ceil(float64(totalData) / float64(pageSize)))
    return pengeluaran, totalPage, nil
}