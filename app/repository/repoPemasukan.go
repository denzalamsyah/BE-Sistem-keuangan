package repository

import (
	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type PemasukanRepository interface {
	FindAll() ([]models.PemasukanResponse, error)
    TotalKeuangan() (int, int, int, error)
	Store(pemasukan *models.Pemasukanlainnya) error
    Update(id int, pemasukan models.Pemasukanlainnya) error
	Delete(id int) error
	GetByID(id int) (*models.Pemasukanlainnya, error)
    GetList() ([]models.Pemasukanlainnya, error)
}

type pemasukanRepository struct {
	db *gorm.DB
}

func NewPemasukanRepo(db *gorm.DB) *pemasukanRepository {
	return &pemasukanRepository{db}
}

func (c *pemasukanRepository) Store(pemasukan *models.Pemasukanlainnya) error {
    tx := c.db.Begin()

    // Create PembayaranSPP
    if err := tx.Create(pemasukan).Error; err != nil {
        tx.Rollback()
        return err
    }

    pemasukanLainnya := models.Pemasukan{
		IDBayar: pemasukan.ID,
        Nama:    pemasukan.Nama,
        Tanggal: pemasukan.Tanggal,
        Jumlah:  pemasukan.Jumlah,
    }
    if err := tx.Create(&pemasukanLainnya).Error; err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}

func (c *pemasukanRepository) Update(id int, pemasukan models.Pemasukanlainnya) error{
    tx := c.db.Begin()

    var pemasukanlainnya models.Pemasukanlainnya

    if err := tx.Where("id = ?", id).First(&pemasukanlainnya).Error; err != nil {
        tx.Rollback()
        return err
    }

    if err := tx.Model(&pemasukanlainnya).Updates(&pemasukan).Error; err != nil {
        tx.Rollback()
        return err
    }

    
    if err := tx.Model(&models.Pemasukan{}).
        Where("id_bayar = ?", pemasukanlainnya.ID).
        Updates(map[string]interface{}{
            "nama" : pemasukanlainnya.Nama,
            "tanggal": pemasukanlainnya.Tanggal,
            "jumlah":  pemasukanlainnya.Jumlah,
        }).Error; err != nil {
        tx.Rollback()
        return err
    }


    return tx.Commit().Error
}

func (c *pemasukanRepository) Delete(id int) error {
	tx := c.db.Begin()

	var pemasukanLainnya models.Pemasukanlainnya

	if err := tx.Where("id = ?", id).First(&pemasukanLainnya).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&pemasukanLainnya).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("nama = ? AND tanggal = ? AND jumlah = ? AND id_bayar = ?", pemasukanLainnya.Nama, pemasukanLainnya.Tanggal, pemasukanLainnya.Jumlah, pemasukanLainnya.ID).Delete(&models.Pemasukan{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (c *pemasukanRepository) GetByID(id int) (*models.Pemasukanlainnya, error){
    var pemasukan models.Pemasukanlainnya

    err := c.db.Where("id = ?", id).First(&pemasukan).Error

    if err != nil {
		return nil, err
	}

    return &pemasukan, nil
}

func (c *pemasukanRepository) GetList() ([]models.Pemasukanlainnya, error){
    var pemasukan []models.Pemasukanlainnya

    err := c.db.Find(&pemasukan).Error
    if err != nil {
		return nil, err
	}
    return pemasukan, nil
}

func (c *pemasukanRepository) FindAll() ([]models.PemasukanResponse, error) {
	var pemasukan []models.Pemasukan
	err := c.db.Find(&pemasukan).Error
	if err != nil {
		return nil, err
	}
    var pemasukanResponse []models.PemasukanResponse
    for _,s := range pemasukan{
        pemasukanResponse = append(pemasukanResponse, models.PemasukanResponse{
            ID: s.ID,
            Nama: s.Nama,
            Tanggal: s.Tanggal,
            Jumlah: s.Jumlah,
        })
    }

	return pemasukanResponse, nil
}
func (c *pemasukanRepository) TotalKeuangan() (int, int, int, error) {
    var totalPemasukan int
    var totalPengeluaran int

    // Menghitung total pemasukan
    if err := c.db.Model(&models.Pemasukan{}).Select("COALESCE(SUM(jumlah), 0) as total_pemasukan").Scan(&totalPemasukan).Error; err != nil {
        return 0, 0, 0, err
    }

    // Menghitung total pengeluaran
    if err := c.db.Model(&models.Pengeluaran{}).Select("COALESCE(SUM(jumlah), 0) as total_pengeluaran").Scan(&totalPengeluaran).Error; err != nil {
        return 0, 0, 0, err
    }

    saldo := totalPemasukan - totalPengeluaran
    return saldo, totalPengeluaran, totalPemasukan, nil
}

