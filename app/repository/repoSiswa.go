package repository

import (
	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type SiswaRepository interface {
	Store(Siswa *models.Siswa) error
	Update(id int, Siswa models.Siswa) error
	Delete(id int) error
	GetByID(id int) (*models.SiswaResponse, error)
	GetList() ([]models.SiswaResponse, error)
	HistoryPembayaranSiswa(siswaID int) ([]models.HistoryPembayaran, error)
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

func (c *siswaRepository) GetByID(id int) (*models.SiswaResponse, error) {
	var siswa models.Siswa
	err := c.db.Preload("Kelas").Preload("Jurusan").Preload("Agama").Preload("Gender").Where("id = ?", id).First(&siswa).Error
	if err != nil {
		return nil, err
	}

	siswaResponse := models.SiswaResponse{
		ID:           siswa.ID,
		Nama:         siswa.Nama,
		NISN:         siswa.NISN,
		Kelas:        siswa.Kelas.Kelas,
		Jurusan:      siswa.Jurusan.Nama,
		Agama:        siswa.Agama.Nama,
		TempatLahir:  siswa.TempatLahir,
		TanggalLahir: siswa.TanggalLahir,
		Gender:       siswa.Gender.Nama,
		NamaAyah:     siswa.NamaAyah,
		NamaIbu:      siswa.NamaIbu,
		NomorTelepon: siswa.NomorTelepon,
		Email:        siswa.Email,
		Alamat:       siswa.Alamat,
	}

	return &siswaResponse, nil
}

func (c *siswaRepository) GetList() ([]models.SiswaResponse, error) {
	var Siswa []models.Siswa
	err := c.db.Preload("Kelas").Preload("Jurusan").Preload("Agama").Preload("Gender").Find(&Siswa).Error
	if err != nil {
		return nil, err
	}

	var SiswaResponse []models.SiswaResponse
	for _, s := range Siswa{
		SiswaResponse = append(SiswaResponse, models.SiswaResponse{
			ID:           s.ID,
			Nama:         s.Nama,
			NISN:         s.NISN,
			Kelas:        s.Kelas.Kelas,
			Jurusan:      s.Jurusan.Nama,
			Agama:        s.Agama.Nama,
			TempatLahir:  s.TempatLahir,
			TanggalLahir: s.TanggalLahir,
			Gender:       s.Gender.Nama,
			NamaAyah:     s.NamaAyah,
			NamaIbu:      s.NamaIbu,
			NomorTelepon: s.NomorTelepon,
			Email:        s.Email,
			Alamat:       s.Alamat,
		})
	}
	return SiswaResponse, nil
}

func (c *siswaRepository) HistoryPembayaranSiswa(siswaID int) ([]models.HistoryPembayaran, error) {
    var historyPembayaran []models.HistoryPembayaran

    // Mengambil data pembayaran dari PembayaranSPP
    var pembayaranSPP []models.PembayaranSPP
    if err := c.db.Preload("Siswa").Preload("Penerima").Where("siswa_id = ?", siswaID).Find(&pembayaranSPP).Error; err != nil {
        return nil, err
    }

    // Mengonversi data PembayaranSPP ke HistoryPembayaran
    for _, p := range pembayaranSPP {
        historyPembayaran = append(historyPembayaran, models.HistoryPembayaran{
            Siswa:          p.Siswa.Nama,
            Nama_transaksi: "Pembayaran SPP",
            Biaya:          p.Jumlah,
            Tanggal:        p.Tanggal,
            Penerima:       p.Penerima.Nama,
        })
    }

    // Mengambil data pembayaran dari PembayaranSemester
    var pembayaranSemester []models.PembayaranSemester
    if err := c.db.Preload("Siswa").Preload("Penerima").Where("siswa_id = ?", siswaID).Find(&pembayaranSemester).Error; err != nil {
        return nil, err
    }

    // Mengonversi data PembayaranSemester ke HistoryPembayaran
    for _, p := range pembayaranSemester {
        historyPembayaran = append(historyPembayaran, models.HistoryPembayaran{
            Siswa:          p.Siswa.Nama,
            Nama_transaksi: p.Transaksi.Nama,
            Biaya:          p.Jumlah,
            Tanggal:        p.Tanggal,
            Penerima:       p.Penerima.Nama,
        })
    }

    return historyPembayaran, nil
}
