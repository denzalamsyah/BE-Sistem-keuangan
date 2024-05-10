package repository

import (
	"errors"
	"math"
	"strings"
	"time"

	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type TransaksiRepository interface {
	Store(Transaksi *models.Transaksi) error
	Update(id int, Transaksi models.Transaksi) error
	Delete(id int) error
	GetList(page, paedSize int) ([]models.Transaksi, int, error)	
	Search(nama, kategori  string) ([]models.Transaksi, error)

}

type transaksiRepository struct {
	db *gorm.DB
}

func NewTransaksiRepo(db *gorm.DB) *transaksiRepository {
	return &transaksiRepository{db}
}

func (c *transaksiRepository) Store(Transaksi *models.Transaksi) error {
	var count int64
    if err := c.db.Model(&models.Transaksi{}).
        Where("nama = ?", Transaksi.Nama).
        Where("kategori = ?", Transaksi.Kategori).
        Count(&count).Error; err != nil {
    
        return err
    }

    if count > 0 {
        return errors.New("Transaksi sudah ada")
    }

	Transaksi.CreatedAt = time.Now()
	err := c.db.Create(Transaksi).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *transaksiRepository) Update(id int, Transaksi models.Transaksi) error {
	err := c.db.Model(&Transaksi).Where("id = ?", id).Updates(&Transaksi).Error
	if err != nil {
		return err
	}
	err = c.db.Model(&models.Transaksi{}).Where("id = ?", id).Update("updated_at", time.Now()).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *transaksiRepository) Delete(id int) error {
	var Transaksi models.Transaksi
	err := c.db.Where("id = ?", id).Delete(&Transaksi).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *transaksiRepository) GetList(page, pageSize int) ([]models.Transaksi, int, error) {
	var Transaksi []models.Transaksi
	offset := (page - 1) * pageSize

	err := c.db.Limit(pageSize).Offset(offset).Find(&Transaksi).Error
	if err != nil {
		return nil, 0, err
	}

	var totalData int64
	if err := c.db.Model(&models.Transaksi{}).Count(&totalData).Error; err != nil {
		return nil, 0, err
	}
	totalPage := int(math.Ceil(float64(totalData) / float64(pageSize)))

	return Transaksi, totalPage, nil
}

func (c *transaksiRepository) Search(nama, kategori string) ([]models.Transaksi, error){
	nama = strings.ToLower(nama)
	kategori = strings.ToLower(kategori)

	var Transaksi []models.Transaksi

	query := c.db.Table("transaksis").
	Select("transaksis.id, transaksis.nama, transaksis.jumlah, transaksis.kategori, transaksis.created_at, transaksis.updated_at").
	Where("LOWER(transaksis.nama) LIKE ? AND LOWER(transaksis.kategori) LIKE ?", "%" +nama+ "%", "%" +kategori+ "%" )

	if err := query.Find(&Transaksi).Error; err != nil {
        return nil, err
    }

	return Transaksi, nil
}