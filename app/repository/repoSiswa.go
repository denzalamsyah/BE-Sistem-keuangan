package repository

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type SiswaRepository interface {
	Store(Siswa *models.Siswa) error
	Update(nisn string, Siswa models.Siswa) error
	Delete(nisn string) error
	GetByID(nisn string) (*models.SiswaResponse, error)
	GetList(page, pageSize int) ([]models.SiswaResponse, int, error)
	HistoryPembayaranSiswa(siswaID string, page, pageSize int) ([]models.HistoryPembayaran, int, error)
	GetTotalGenderCount() (int, int, error)
	Search(name, nisn, kelas, jurusan, angkatan string) ([]models.SiswaResponse, error)
	SearchByKodeKelas(name, nisn, kodeKelas string) ([]models.SiswaResponse, error)
	GetUserNisn(nisn string) (models.Siswa, error)
	
}

type siswaRepository struct{
	db *gorm.DB
}

func NewSiswaRepo(db *gorm.DB) *siswaRepository {
	return &siswaRepository{db}
}

func (c *siswaRepository) Store(Siswa *models.Siswa) error {
	Siswa.CreatedAt = time.Now().Format("02 January 2006 15:04:05")
	if err := c.db.Create(Siswa).Error; err != nil {
		return fmt.Errorf("failed to store new siswa: %v", err)
	}

	return nil
}

func (c *siswaRepository) Update(nisn string, Siswa models.Siswa) error {
	err := c.db.Model(&models.Siswa{}).Where("nisn = ?", nisn).Updates(&Siswa).Error
	if err != nil {
		return err
	}

	err = c.db.Model(&models.Siswa{}).Where("nisn = ?", nisn).Update("updated_at", time.Now().Format("02 January 2006 15:04:05")).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *siswaRepository) Delete(nisn string) error {
	var Siswa models.Siswa
	err := c.db.Where("nisn = ?", nisn).Delete(&Siswa).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *siswaRepository) GetByID(nisn string) (*models.SiswaResponse, error) {
	var siswa models.Siswa
	err := c.db.Preload("Kelas").Preload("Jurusan").Preload("Agama").Preload("Gender").Where("nisn = ?", nisn).First(&siswa).Error
	if err != nil {
		return nil, err
	}

	siswaResponse := models.SiswaResponse{
		NISN:         siswa.Nisn,
		Nama:         siswa.Nama,
		Kelas:   siswa.Kelas.Kelas,
		Jurusan:      siswa.Jurusan.Jurusan,
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
		CreatedAt: siswa.CreatedAt,
		UpdatedAt: siswa.UpdatedAt,
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
			NISN:         s.Nisn,
			Nama:         s.Nama,
			Kelas:       s.Kelas.Kelas,
			Jurusan:      s.Jurusan.Jurusan,
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
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		})
	}
	totalPage := int(math.Ceil(float64(totalData) / float64(pageSize)))
    return SiswaResponse, totalPage, nil
}

func (c *siswaRepository) HistoryPembayaranSiswa(siswaID string, page, pageSize int) ([]models.HistoryPembayaran,int, error) {
    var historyPembayaran []models.HistoryPembayaran

    // Mengambil data pembayaran dari PembayaranSemester
    var pembayaranSemester []models.PembayaranSemester
    if err := c.db.Preload("Siswa").Preload("Transaksi").Where("siswa_id = ?", siswaID).
	Offset((page - 1) * pageSize).Limit(pageSize).Find(&pembayaranSemester).Error; err != nil {
        return nil,0, err
    }

	var totalData int64
    if err := c.db.Model(&models.PembayaranSemester{}).Count(&totalData).Error; err != nil {
        return nil, 0, err
    }

    // Mengonversi data PembayaranSemester ke HistoryPembayaran
    for _, p := range pembayaranSemester {
        historyPembayaran = append(historyPembayaran, models.HistoryPembayaran{
			ID: p.ID,
            Siswa:          p.Siswa.Nama,
			NISN: p.Siswa.Nisn,
            Nama_transaksi: p.Transaksi.Nama,
            Biaya:          p.Jumlah,
            Tanggal:        p.Tanggal,
			Status: p.Status,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
        })
    }
	totalPage := int(math.Ceil(float64(totalData) / float64(pageSize)))

    return historyPembayaran, totalPage, nil
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

func (c *siswaRepository) Search(name, nisn, kelas, jurusan, angkatan string) ([]models.SiswaResponse, error) {
    name = strings.ToLower(name)
    kelas = strings.ToLower(kelas)
    jurusan = strings.ToLower(jurusan)
	angkatan = strings.ToLower(angkatan)
	nisn = strings.ToLower(nisn)

    var siswaList []models.SiswaResponse

    // Query dengan menggunakan Select untuk menentukan kolom yang akan diambil
    query := c.db.Table("siswas").
        Select("siswas.nama, siswas.nisn, kelas.kelas as kelas, jurusans.jurusan as jurusan, agamas.nama as agama, siswas.tempat_lahir, siswas.tanggal_lahir, genders.nama as gender, siswas.nama_ayah, siswas.nama_ibu, siswas.nomor_telepon, siswas.angkatan, siswas.email, siswas.alamat, siswas.gambar, siswas.created_at, siswas.updated_at").
        Joins("JOIN kelas ON siswas.kelas_id = kelas.kode_kelas").
        Joins("JOIN jurusans ON siswas.jurusan_id = jurusans.kode_jurusan").
		Joins("JOIN agamas ON siswas.agama_id = agamas.id_agama").
		Joins("JOIN genders ON siswas.gender_id = genders.id_gender").
        Where("LOWER(siswas.nama) LIKE ? AND LOWER(siswas.nisn) LIKE ? AND LOWER(kelas.kelas) LIKE ? AND LOWER(jurusans.jurusan) LIKE ? AND LOWER(siswas.angkatan) LIKE ?", "%"+name+"%", "%"+nisn+"%", "%"+kelas+"%", "%"+jurusan+"%", "%"+angkatan+"%")
    if err := query.Find(&siswaList).Error; err != nil {
        return nil, err
    }

    return siswaList, nil
}

func (c *siswaRepository) SearchByKodeKelas(name, nisn, kodeKelas string) ([]models.SiswaResponse, error) {
	var siswaList []models.SiswaResponse
	name = strings.ToLower(name)
	kodeKelas = strings.ToLower(kodeKelas)
	nisn = strings.ToLower(nisn)
	
    // Query dengan menggunakan Select untuk menentukan kolom yang akan diambil
    query := c.db.Table("siswas").
		Select("siswas.nama, siswas.nisn, kelas.kelas as kelas, jurusans.jurusan as jurusan, agamas.nama as agama, siswas.tempat_lahir, siswas.tanggal_lahir, genders.nama as gender, siswas.nama_ayah, siswas.nama_ibu, siswas.nomor_telepon, siswas.angkatan, siswas.email, siswas.alamat, siswas.gambar, siswas.created_at, siswas.updated_at").
        Joins("JOIN kelas ON siswas.kelas_id = kelas.kode_kelas").
        Joins("JOIN jurusans ON siswas.jurusan_id = jurusans.kode_jurusan").
        Joins("JOIN agamas ON siswas.agama_id = agamas.id_agama").
        Joins("JOIN genders ON siswas.gender_id = genders.id_gender").
        Where("LOWER(siswas.nama) LIKE ? AND LOWER(siswas.nisn) LIKE ? AND LOWER(kelas.kode_kelas) LIKE ?", "%"+name+"%", "%"+nisn+"%", "%"+kodeKelas+"%")
    if err := query.Find(&siswaList).Error; err != nil {
        return nil, err
    }

    return siswaList, nil
}

func (c *siswaRepository) GetUserNisn(nisn string) (models.Siswa, error){
	var Siswa models.Siswa

	result :=c.db.Where("nisn = ?", nisn).First(&Siswa)
	if result.Error != nil{
		if result.Error == gorm.ErrRecordNotFound{
			return Siswa, nil
		}
		return Siswa, result.Error
	}
	return Siswa, nil
}



