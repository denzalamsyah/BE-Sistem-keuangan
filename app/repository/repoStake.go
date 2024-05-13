package repository

import (
	"errors"
	"math"
	"strings"
	"time"

	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type GuruRepository interface {
	Store(Guru *models.Guru) error
	Update(nip string, Guru models.Guru) error
	Delete(nip string) error
	GetByID(nip string) (*models.GuruResponse, error)
	GetList(page, pageSize int) ([]models.GuruResponse, int, error)
	GetTotalGenderCount() (int, int, error)
	Search(nama, nip, jabatan string) ([]models.GuruResponse, error)
	GetUserNIP(nip string) (models.Guru, error)
	HistoryPembayaranGuru(nip string, page, pageSize int) ([]models.HistoryPembayaranKas, int, error)
	SaldoKasByNIP(nip string) (int, error)
	AmbilKasGuru( jumlah int, nip, nama, tanggal string) error
	HistoryPengambilanKas(nip string, page, pageSize int) ([]models.HistoryPengambilanKas, int, error)
	// HistoryPembayaranKas(GuruID, page, pageSize int) ([]models.HistoryPembayaranKas, int, error)

}

type guruRepository struct{
	db *gorm.DB
}


func NewGuruRepo(db *gorm.DB) *guruRepository{
	return &guruRepository{db}
}


func(c *guruRepository) Store(Guru *models.Guru) error{
	Guru.CreatedAt = time.Now().Format("02 January 2006 15:04:05")
	err := c.db.Create(Guru).Error
	if err != nil{
		return err
	}
	return nil
}

func (c *guruRepository) Update(nip string, Guru models.Guru) error {
	err := c.db.Model(&models.Guru{}).Where("nip = ?", nip).Updates(&Guru).Error
	if err != nil {
		return err
	}
	err = c.db.Model(&models.Guru{}).Where("nip = ?", nip).Update("updated_at", time.Now().Format("02 January 2006 15:04:05")).Error
	if err != nil {
		return err
	}
	return nil
}

func(c *guruRepository) Delete(nip string) (error){
	var guru models.Guru
	err := c.db.Where("nip = ?", nip).Delete(&guru).Error

	if err != nil{
		return err
	}
	return nil
}

func (c *guruRepository) GetList(page, pageSize int) ([]models.GuruResponse, int, error){
	var guru []models.Guru

	err := c.db.Preload("Jabatan").Preload("Agama").Preload("Gender").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&guru).Error
	if err != nil {
		return nil,0, err
	}

	var totalData int64
	if err := c.db.Model(&models.Guru{}).Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	var guruResponse []models.GuruResponse
	for _, s := range guru{
		guruResponse = append(guruResponse, models.GuruResponse{
			Nip: s.Nip,
			Nama: s.Nama,
			Agama: s.Agama.Nama,
			Gender: s.Gender.Nama,
			Jabatan: s.Jabatan.Nama,
			TempatLahir: s.TempatLahir,
			TanggalLahir: s.TanggalLahir,
			NomorTelepon: s.NomorTelepon,
			Email: s.Email,
			Alamat: s.Alamat,
			Gambar: s.Gambar,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		})
	}
	totalPage := int(math.Ceil(float64(totalData) / float64(pageSize)))
	return guruResponse, totalPage, nil
}

func (c *guruRepository) GetByID(nip string) (*models.GuruResponse, error){
	var guru models.Guru

	err := c.db.Preload("Jabatan").Preload("Agama").Preload("Gender").Where("nip = ?", nip).First(&guru).Error
	if err != nil {
		return nil, err
	}

	guruResponse := models.GuruResponse{
		Nama: guru.Nama,
		Nip: guru.Nip,
		Agama: guru.Agama.Nama,
		Gender: guru.Gender.Nama,
		Jabatan: guru.Jabatan.Nama,
		TempatLahir: guru.TempatLahir,
		TanggalLahir: guru.TanggalLahir,
		NomorTelepon: guru.NomorTelepon,
		Email: guru.Email,
		Alamat: guru.Alamat,
		Gambar: guru.Gambar,
		CreatedAt: guru.CreatedAt,
		UpdatedAt: guru.UpdatedAt,
	}
	return &guruResponse, nil
}

func (c *guruRepository) GetTotalGenderCount() (int, int, error) {
    var countLakiLaki, countPerempuan int64
    if err := c.db.Model(&models.Guru{}).Where("gender_id = ?", 1).Count(&countLakiLaki).Error; err != nil {
        return 0, 0, err
    }

    if err := c.db.Model(&models.Guru{}).Where("gender_id = ?", 2).Count(&countPerempuan).Error; err != nil {
        return 0, 0, err
    }

    return int(countLakiLaki), int(countPerempuan), nil
}

func (c *guruRepository) Search(nama, nip, jabatan string) ([]models.GuruResponse, error){
	nama = strings.ToLower(nama)
	jabatan = strings.ToLower(jabatan)
	nip = strings.ToLower(nip)

	var guruList []models.GuruResponse

	query := c.db.Table("gurus").
	Select("gurus.nama, gurus.nip, agamas.nama as agama, jabatans.nama as jabatan, gurus.tempat_lahir, gurus.tanggal_lahir, genders.nama as gender, gurus.nomor_telepon, gurus.email, gurus.alamat, gurus.gambar, gurus.created_at, gurus.updated_at ").
	Joins("JOIN agamas ON gurus.agama_id = agamas.id_agama").
	Joins("JOIN jabatans on gurus.jabatan_id = jabatans.id_jabatan").
	Joins("JOIN genders ON gurus.gender_id = genders.id_gender").
	Where("LOWER(gurus.nama) LIKE ? AND LOWER(gurus.nip) LIKE ? AND LOWER(jabatans.nama) LIKE ?", "%"+nama+"%", "%"+nip+"%", "%"+jabatan+"%")

	if err := query.Find(&guruList).Error; err != nil {
        return nil, err
    }
	
	return guruList, nil

}

func (c *guruRepository) GetUserNIP(nip string) (models.Guru, error){
	var Guru models.Guru

	result :=c.db.Where("nip = ?", nip).First(&Guru)
	if result.Error != nil{
		if result.Error == gorm.ErrRecordNotFound{
			return Guru, nil
		}
		return Guru, result.Error
	}
	return Guru, nil
}

func (c *guruRepository) HistoryPembayaranGuru(nip string, page, pageSize int) ([]models.HistoryPembayaranKas, int, error){
	var historiPembayaran []models.HistoryPembayaranKas
    var kasGuru []models.KasGuru

    if err := c.db.Preload("Guru").Where("guru_id = ?", nip).
        Offset((page - 1) * pageSize).Limit(pageSize).Find(&kasGuru).Error; err != nil {
        return nil, 0, err
    }

	var totalData int64
    if err := c.db.Model(&models.KasGuru{}).Where("guru_id = ?", nip).Count(&totalData).Error; err != nil {
        return nil, 0, err
    }

	  for _, p := range kasGuru {
        historiPembayaran = append(historiPembayaran, models.HistoryPembayaranKas{
			ID: uint(p.ID),
            Nama:         p.Guru.Nama,
            NIP:          p.Guru.Nip,
            Jumlah_Bayar: p.Jumlah,
            Tanggal:      p.TanggalBayar,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
        })
    }
	totalPage := int(math.Ceil(float64(totalData) / float64(pageSize)))

    return historiPembayaran, totalPage, nil
}

func ( c *guruRepository) SaldoKasByNIP(nip string) (int, error){
	var totalKas int

	if err := c.db.Model(&models.KasGuru{}).Where("guru_id = ?", nip).Select("SUM(saldo)").Scan(&totalKas).Error; err != nil {
		return 0, err
	}

	return totalKas, nil
}

func (c *guruRepository) AmbilKasGuru(jumlah int, nip, nama, tanggal string) error {
	var kasGuru models.KasGuru

	if err := c.db.Where("guru_id = ?", nip).Last(&kasGuru).Error; err != nil {
		return err
	}

	if kasGuru.Saldo < jumlah {
		return errors.New("Saldo tidak mencukupi")
	}

	kasGuru.Saldo -= jumlah

	if err := c.db.Save(&kasGuru).Error; err != nil {
		return err
	}
	// Menyimpan histori pengambilan kas
	pengambilanKas := models.PengambilanKas{
		GuruID:      nip,
		Nama: nama,
		JumlahAmbil: jumlah,
		TanggalAmbil: tanggal,
		CreatedAt:time.Now().Format("02 January 2006 15:04:05"),
		UpdatedAt:time.Now().Format("02 January 2006 15:04:05"),
	}
	if err := c.db.Create(&pengambilanKas).Error; err != nil {
		return err
	}


	return nil
}

func ( c *guruRepository) HistoryPengambilanKas(nip string, page, pageSize int) ([]models.HistoryPengambilanKas, int, error){
	var historiPengambilan []models.HistoryPengambilanKas

	var kasGuru []models.PengambilanKas

	if err := c.db.Where("guru_id = ?", nip).
	Offset((page - 1) * pageSize).Limit(pageSize).Find(&kasGuru).Error; err != nil {
        return nil,0, err
    }
	var totalData int64
    if err := c.db.Model(&models.KasGuru{}).Count(&totalData).Error; err != nil {
        return nil, 0, err
    }

	for _, p := range kasGuru{
		historiPengambilan = append(historiPengambilan, models.HistoryPengambilanKas{
			ID: p.ID,
			NIP: p.GuruID,
			Nama: p.Nama,
			JumlahAmbil: p.JumlahAmbil,
			TanggalAmbil: p.TanggalAmbil,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}
	totalPage := int(math.Ceil(float64(totalData) / float64(pageSize)))

    return historiPengambilan, totalPage, nil
}

