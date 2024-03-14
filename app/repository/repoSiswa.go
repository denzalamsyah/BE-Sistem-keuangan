package repository

import (
	"fmt"
	"math"
	"strings"

	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type SiswaRepository interface {
	Store(Siswa *models.Siswa) error
	Update(id int, Siswa models.Siswa) error
	Delete(id int) error
	GetByID(id int) (*models.SiswaResponse, error)
	GetList(page, pageSize int) ([]models.SiswaResponse, int, error)
	HistoryPembayaranSiswa(siswaID int) ([]models.HistoryPembayaran, error)
	GetTotalGenderCount() (int, int, error)
	Search(name, nisn, kelas, jurusan string) ([]models.SiswaResponse, error)

}

type siswaRepository struct{
	db *gorm.DB
}

func NewSiswaRepo(db *gorm.DB) *siswaRepository {
	return &siswaRepository{db}
}

func (c *siswaRepository) Store(siswa *models.Siswa) error {
	if err := c.db.Create(siswa).Error; err != nil {
		return fmt.Errorf("failed to store new siswa: %v", err)
	}

	return nil
}

func (c *siswaRepository) Update(id int, Siswa models.Siswa) error {
	err := c.db.Model(&models.Siswa{}).Where("id = ?", id).Updates(&Siswa).Error
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
		Kelas:   siswa.Kelas.Kelas,
		Jurusan:      siswa.Jurusan.Nama,
		Agama:        siswa.Agama.Nama,
		TempatLahir:  siswa.TempatLahir,
		TanggalLahir: siswa.TanggalLahir,
		Gender:       siswa.Gender.Nama,
		NamaAyah:     siswa.NamaAyah,
		NamaIbu:      siswa.NamaIbu,
		Angkatan:     siswa.Angkatan,
		NomorTelepon: siswa.NomorTelepon,
		Email:        siswa.Email,
		Alamat:       siswa.Alamat,
		Gambar:       siswa.Gambar,
	}

	return &siswaResponse, nil
}

func (c *siswaRepository) GetList(page, pageSize int) ([]models.SiswaResponse, int, error) {
	var Siswa []models.Siswa
    err := c.db.Preload("Kelas").Preload("Jurusan").Preload("Agama").Preload("Gender").
        Offset((page - 1) * pageSize).Limit(pageSize).Find(&Siswa).Error
    if err != nil {
        return nil, 0, err
    }

    var totalData int64
    if err := c.db.Model(&models.Siswa{}).Count(&totalData).Error; err != nil {
        return nil, 0, err
    }

	var SiswaResponse []models.SiswaResponse
	for _, s := range Siswa{
		SiswaResponse = append(SiswaResponse, models.SiswaResponse{
			ID:           s.ID,
			Nama:         s.Nama,
			NISN:         s.NISN,
			Kelas:       s.Kelas.Kelas,
			Jurusan:      s.Jurusan.Nama,
			Agama:        s.Agama.Nama,
			TempatLahir:  s.TempatLahir,
			TanggalLahir: s.TanggalLahir,
			Gender:       s.Gender.Nama,
			NamaAyah:     s.NamaAyah,
			NamaIbu:      s.NamaIbu,
			Angkatan:     s.Angkatan,
			NomorTelepon: s.NomorTelepon,
			Email:        s.Email,
			Alamat:       s.Alamat,
			Gambar:       s.Gambar,
		})
	}
	totalPage := int(math.Ceil(float64(totalData) / float64(pageSize)))
    return SiswaResponse, totalPage, nil
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

func (c *siswaRepository) GetTotalGenderCount() (int, int, error) {
    var countLakiLaki, countPerempuan int64
    if err := c.db.Model(&models.Siswa{}).Where("gender_id = ?", 1).Count(&countLakiLaki).Error; err != nil {
        return 0, 0, err
    }

    if err := c.db.Model(&models.Siswa{}).Where("gender_id = ?", 2).Count(&countPerempuan).Error; err != nil {
        return 0, 0, err
    }

    return int(countLakiLaki), int(countPerempuan), nil
}

func (c *siswaRepository) Search(name, nisn, kelas, jurusan string) ([]models.SiswaResponse, error) {
    name = strings.ToLower(name)
    kelas = strings.ToLower(kelas)
    jurusan = strings.ToLower(jurusan)

    var siswaList []models.SiswaResponse

    // Query dengan menggunakan Select untuk menentukan kolom yang akan diambil
    query := c.db.Table("siswas").
        Select("siswas.id, siswas.nama, siswas.nisn, kelas.kelas as kelas, jurusans.nama as jurusan, agamas.nama as agama, siswas.tempat_lahir, siswas.tanggal_lahir, genders.nama as gender, siswas.nama_ayah, siswas.nama_ibu, siswas.nomor_telepon, siswas.angkatan, siswas.email, siswas.alamat, siswas.gambar").
        Joins("JOIN kelas ON siswas.kelas_id = kelas.id_kelas").
        Joins("JOIN jurusans ON siswas.jurusan_id = jurusans.id_jurusan").
		Joins("JOIN agamas ON siswas.agama_id = agamas.id_agama").
		Joins("JOIN genders ON siswas.gender_id = genders.id_gender").
        Where("LOWER(siswas.nama) LIKE ? AND siswas.nisn::TEXT LIKE ? AND LOWER(kelas.kelas) LIKE ? AND LOWER(jurusans.nama) LIKE ?", "%"+name+"%", "%"+nisn+"%", "%"+kelas+"%", "%"+jurusan+"%")

    if err := query.Find(&siswaList).Error; err != nil {
        return nil, err
    }

    return siswaList, nil
}



