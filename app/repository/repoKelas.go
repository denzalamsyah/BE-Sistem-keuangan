package repository

import (
	"math"
	"strings"
	"time"

	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type KelasRepository interface {
	Store(Kelas *models.Kelas) error
	Update(kode string, Kelas models.Kelas) error
	Delete(kode string) error
	GetList(page, pageSize int) ([]models.Kelas, int, error)
	GetTotalKelasCount() (int, error)
	Search(nama string) ([]models.Kelas, error)
	GetKode(kode string) (models.Kelas, error)

}

type kelasRepository struct{
	db *gorm.DB
}

func NewKelasRepo(db *gorm.DB) *kelasRepository {
	return &kelasRepository{db}
}

func (c *kelasRepository) Store(Kelas *models.Kelas) error {
	Kelas.CreatedAt = time.Now()
	err := c.db.Create(Kelas).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *kelasRepository) Update(kode string, Kelas models.Kelas) error {
	err := c.db.Model(&models.Kelas{}).Where("kode_kelas = ?", kode).Updates(&Kelas).Error
	if err != nil {
		return err
	}
	err = c.db.Model(&models.Kelas{}).Where("kode_kelas = ?", kode).Update("updated_at", time.Now()).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *kelasRepository) Delete(kode string) error {
	var Kelas models.Kelas
	err := c.db.Where("kode_kelas = ?", kode).Delete(&Kelas).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *kelasRepository) GetList(page, pageSize int) ([]models.Kelas, int, error) {
	var Kelas []models.Kelas
	offset := (page - 1) * pageSize

	err := c.db.Limit(pageSize).Offset(offset).Find(&Kelas).Error
	if err != nil {
		return nil, 0,err
	}

	var totalData int64
	if err := c.db.Model(&models.Kelas{}).Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	totalPage := int(math.Ceil(float64(totalData) / float64(pageSize)))
	return Kelas, totalPage, nil
}

func (c *kelasRepository) GetTotalKelasCount() (int, error) {
    var count int64
    if err := c.db.Model(&models.Kelas{}).Count(&count).Error; err != nil {
        return 0, err
    }
    return int(count), nil
}

func (c *kelasRepository) Search(nama string) ([]models.Kelas, error){
	nama = strings.ToLower(nama)

	var Kelas []models.Kelas

	query := c.db.Table("kelas").
	Select("kelas.kode_kelas, kelas.kelas, kelas.created_at, kelas.updated_at").
	Where("LOWER(kelas.kelas) LIKE ?", "%" +nama+ "%")

	if err := query.Find(&Kelas).Error; err != nil {
        return nil, err
    }

	return Kelas, nil
}

func (c *kelasRepository) GetKode(kode string) (models.Kelas, error){
	var Kelas models.Kelas

	result :=c.db.Where("kode_kelas = ?", kode).First(&Kelas)
	if result.Error != nil{
		if result.Error == gorm.ErrRecordNotFound{
			return Kelas, nil
		}
		return Kelas, result.Error
	}
	return Kelas, nil
}
