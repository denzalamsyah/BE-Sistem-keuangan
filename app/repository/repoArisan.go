package repository

import (
	"fmt"
	"math"
	"strings"

	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type ArisanRepository interface {
	Store(Arisan *models.Arisan) error
	Update(id int, Arisan models.Arisan) error
	GetList(page, pageSize int) ([]models.Arisan,int, error)
	Delete(id int) error
	Search(nama, tanggal string) ([]models.Arisan, error)
}

type arisanRepository struct{
	db *gorm.DB
}

func NewArisanRepo(db *gorm.DB) *arisanRepository{
	return &arisanRepository{db}
}

func (c *arisanRepository) Store(Arisan *models.Arisan) error{
	if err := c.db.Create(Arisan).Error; err != nil {
		return fmt.Errorf("failed to store new siswa: %v", err)
	}
	return nil
}

func (c *arisanRepository) Update(id int, Arisan models.Arisan) error{
	err := c.db.Model(&models.Arisan{}).Where("id = ?", id).Updates(&Arisan).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *arisanRepository) Delete(id int) error{
	var Arisan models.Arisan
	err := c.db.Where("id = ?", id).Delete(&Arisan).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *arisanRepository) GetList(page, pageSize int) ([]models.Arisan,int, error){
	var Arisan []models.Arisan

	offset := (page - 1) * pageSize

	err := c.db.Limit(pageSize).Offset(offset).Find(&Arisan).Error
	if err != nil {
		return nil, 0,err
	}

	var totalData int64
	if err := c.db.Model(&models.Arisan{}).Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	totalPage := int(math.Ceil(float64(totalData) / float64(pageSize)))
	return Arisan, totalPage, nil
}

func ( c *arisanRepository) Search(nama, tanggal string) ([]models.Arisan, error){
	nama = strings.ToLower(nama)
	tanggal = strings.ToLower(tanggal)


	var Arisan []models.Arisan

	query := c.db.Table("arisans").
	Select("arisans.id, arisans.nama_arisan").
	Where("LOWER(arisans.nama_arisans) LIKE ? AND LOWER(arisans.tanggal_mulai) LIKE ?", "%" +nama+ "%", "%"+tanggal+"%")

	if err := query.Find(&Arisan).Error; err != nil {
        return nil, err
    }

	return Arisan, nil
}