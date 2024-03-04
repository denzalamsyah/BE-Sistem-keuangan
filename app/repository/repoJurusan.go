package repository

import (
	"math"

	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type JurusanRepository interface {
	Store(Jurusan *models.Jurusan) error
	Update(id int, Jurusan models.Jurusan) error
	Delete(id int) error
	GetList(page, pageSize int) ([]models.Jurusan,int, error)
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

func (c *jurusanRepository) GetList(page, pageSize int) ([]models.Jurusan, int, error) {
	var Jurusan []models.Jurusan

	// Menghitung offset berdasarkan halaman dan jumlah data per halaman
	offset := (page - 1) * pageSize

	// Mengambil data Jurusan dengan limit dan offset
	err := c.db.Limit(pageSize).Offset(offset).Find(&Jurusan).Error
	if err != nil {
		return nil, 0, err
	}

	// Menghitung total data Jurusan
	var totalData int64
	if err := c.db.Model(&models.Jurusan{}).Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	// Menghitung total halaman
	totalPage := int(math.Ceil(float64(totalData) / float64(pageSize)))

	return Jurusan, totalPage, nil
}
