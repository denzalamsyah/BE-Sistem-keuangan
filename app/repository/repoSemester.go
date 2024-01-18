package repository

import (
	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type SemesterRepository interface {
	Store(PembayaranSemester *models.PembayaranSemester) error
	Update(id int, PembayaranSemester models.PembayaranSemester) error
	Delete(id int) error
	GetByID(id int) (*models.PembayaranSemesterResponse, error)
	GetList() ([]models.PembayaranSemesterResponse, error)
}

type semesterRepository struct{
	db *gorm.DB
}

func NewSemesterRepo(db *gorm.DB) *semesterRepository{
	return &semesterRepository{db}
}

func (c *semesterRepository) Store(PembayaranSemester *models.PembayaranSemester) error {
	tx := c.db.Begin()

    if err := tx.Create(PembayaranSemester).Error; err != nil {
        tx.Rollback()
        return err
    }
	 
	var transaksi models.Transaksi
    if err := tx.First(&transaksi, PembayaranSemester.TransaksiID).Error; err != nil {
        tx.Rollback()
        return err
    }

    pemasukanSemester := models.Pemasukan{
		IDBayar: PembayaranSemester.ID,
        Nama:    transaksi.Nama,
        Tanggal: PembayaranSemester.Tanggal,
        Jumlah:  PembayaranSemester.Jumlah,
    }
    if err := tx.Create(&pemasukanSemester).Error; err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}

func (c *semesterRepository) Update(id int, PembayaranSemester models.PembayaranSemester) error {
	tx := c.db.Begin()
	var pembayaranSemester models.PembayaranSemester

	if err := tx.Where("id = ?", id).First(&pembayaranSemester).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&pembayaranSemester).Updates(&PembayaranSemester).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&models.Pemasukan{}).
	Where("id_bayar = ?", pembayaranSemester.ID).
	Updates(map[string]interface{}{
		"nama":    pembayaranSemester.Transaksi,
		"tanggal": pembayaranSemester.Tanggal,
		"jumlah":  pembayaranSemester.Jumlah,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (c *semesterRepository) Delete(id int) error {
	tx := c.db.Begin()

	var pembayaranSemester models.PembayaranSemester

	if err := tx.Where("id = ?", id).First(&pembayaranSemester).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&pembayaranSemester).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("nama = ? AND tanggal = ? AND jumlah = ? AND id_bayar = ?", pembayaranSemester.Transaksi, pembayaranSemester.Tanggal, pembayaranSemester.Jumlah, pembayaranSemester.ID).Delete(&models.Pemasukan{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}



func (c *semesterRepository) GetByID(id int) (*models.PembayaranSemesterResponse, error) {
	var PembayaranSemester models.PembayaranSemester
	err := c.db.Preload("Siswa").Preload("Penerima").Preload("Transaksi").Where("id = ?", id).First(&PembayaranSemester).Error
	if err != nil {
		return nil, err
	}

	PembayaranSemesterResponse := models.PembayaranSemesterResponse{
		ID:             PembayaranSemester.ID,
		Siswa:          PembayaranSemester.Siswa.Nama,
		Transaksi: PembayaranSemester.Transaksi.Nama,
		Semester:       PembayaranSemester.Semester,
		TahunAjar:      PembayaranSemester.TahunAjar,
		Status:         PembayaranSemester.Status,
		Tanggal:        PembayaranSemester.Tanggal,
		Jumlah:         PembayaranSemester.Jumlah,
		Penerima:       PembayaranSemester.Penerima.Nama,
	}

	return &PembayaranSemesterResponse, nil
}

func (c *semesterRepository) GetList() ([]models.PembayaranSemesterResponse, error) {
	var PembayaranSemester []models.PembayaranSemester
	err := c.db.Preload("Siswa").Preload("Penerima").Preload("Transaksi").Find(&PembayaranSemester).Error
	if err != nil {
		return nil, err
	}

	var PembayaranSemesterResponse []models.PembayaranSemesterResponse
	for _, s := range PembayaranSemester{
		PembayaranSemesterResponse = append(PembayaranSemesterResponse, models.PembayaranSemesterResponse{
			ID:             s.ID,
			Siswa:          s.Siswa.Nama,
			Transaksi: s.Transaksi.Nama,
			Semester:       s.Semester,
			TahunAjar:      s.TahunAjar,
			Status:         s.Status,
			Tanggal:        s.Tanggal,
			Jumlah:         s.Jumlah,
			Penerima:       s.Penerima.Nama,
		})
	}
	return PembayaranSemesterResponse, nil
}