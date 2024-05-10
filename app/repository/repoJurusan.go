package repository

import (
	"math"
	"strings"
	"time"

	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type JurusanRepository interface {
	Store(Jurusan *models.Jurusan) error
	Update(kode string, Jurusan models.Jurusan) error
	Delete(kode string) error
	GetList(page, pageSize int) ([]models.Jurusan,int, error)
	GetTotalJurusanCount() (int, error)
	Search(nama string) ([]models.Jurusan, error)
	GetKode(kode string) (models.Jurusan, error)

}

type jurusanRepository struct {
	db *gorm.DB
}

func NewJurusanRepo(db *gorm.DB) *jurusanRepository {
	return &jurusanRepository{db}
}

func (c *jurusanRepository) Store(Jurusan *models.Jurusan) error {
	Jurusan.CreatedAt = time.Now()
	err := c.db.Create(Jurusan).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *jurusanRepository) Update(kode string, Jurusan models.Jurusan) error {
	err := c.db.Model(&models.Jurusan{}).Where("kode_jurusan = ?", kode).Updates(&Jurusan).Error
	if err != nil {
		return err
	}
	// Atur nilai UpdatedAt dengan waktu sekarang
	err = c.db.Model(&models.Jurusan{}).Where("kode_jurusan = ?", kode).Update("updated_at", time.Now()).Error
	if err != nil {
			return err
		}
	
	return nil
}

func (c *jurusanRepository) Delete(kode string) error {
	var Jurusan models.Jurusan
	err := c.db.Where("kode_jurusan = ?", kode).Delete(&Jurusan).Error
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

func (c *jurusanRepository) GetTotalJurusanCount() (int, error) {
    var count int64
    if err := c.db.Model(&models.Jurusan{}).Count(&count).Error; err != nil {
        return 0, err
    }
    return int(count), nil
}

func (c *jurusanRepository) Search(nama string) ([]models.Jurusan, error){
	nama = strings.ToLower(nama)

	var jurusan []models.Jurusan

	query := c.db.Table("jurusans").
	Select("jurusans.kode_jurusan, jurusans.jurusan, jurusans.created_at, jurusans.updated_at").
	Where("LOWER(jurusans.jurusan) LIKE ?", "%" +nama+ "%")

	if err := query.Find(&jurusan).Error; err != nil {
        return nil, err
    }

	return jurusan, nil
}

func (c *jurusanRepository) GetKode(kode string) (models.Jurusan, error){
	var Jurusan models.Jurusan

	result :=c.db.Where("kode_jurusan = ?", kode).First(&Jurusan)
	if result.Error != nil{
		if result.Error == gorm.ErrRecordNotFound{
			return Jurusan, nil
		}
		return Jurusan, result.Error
	}
	return Jurusan, nil
}