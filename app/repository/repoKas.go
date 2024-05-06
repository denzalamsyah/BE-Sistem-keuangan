package repository

import (
	"fmt"
	"math"
	"strings"

	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type KasRepository interface {
	Store(KasGuru *models.KasGuru) error
	Update(id int, KasGuru models.KasGuru) error
	GetList(page, pageSize int) ([]models.KasGuruResponse,int, error)
	Delete(id int) error
	Search(nama, tanggal string) ([]models.KasGuruResponse, error)
	GetByID(id int) (*models.HistoryPembayaranKas, error)
	GetAmbilByID(id int) (*models.HistoryPengambilanKas, error)
}

type kasRepository struct{
	db *gorm.DB
}

func NewKasRepo(db *gorm.DB) *kasRepository{
	return &kasRepository{db}
}

func (c *kasRepository) Store(KasGuru *models.KasGuru) error{
	if err := c.db.Create(KasGuru).Error; err != nil {
		return fmt.Errorf("failed to store new siswa: %v", err)
	}
	return nil
}

func (c *kasRepository) Update(id int, KasGuru models.KasGuru) error{
	err := c.db.Model(&models.KasGuru{}).Where("id = ?", id).Updates(&KasGuru).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *kasRepository) Delete(id int) error{
	var KasGuru models.KasGuru
	err := c.db.Where("id = ?", id).Delete(&KasGuru).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *kasRepository) GetList(page, pageSize int) ([]models.KasGuruResponse,int, error){
	var KasGuru []models.KasGuru

	err := c.db.Preload("Guru").Offset((page - 1) * pageSize).Limit(pageSize).Find(&KasGuru).Error
	if err != nil {
		return nil, 0, err
	}

	var totalData int64
	if err := c.db.Model(&models.KasGuru{}).Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	var KasGuruResponse []models.KasGuruResponse
	for _, s := range KasGuru{
		KasGuruResponse = append(KasGuruResponse, models.KasGuruResponse{
			ID:           s.ID,
			NIP: s.GuruID,
			NamaGuru: s.Guru.Nama,
			Jumlah: s.Jumlah,
			Saldo: s.Saldo,
			TanggalBayar: s.TanggalBayar,
			
		})
	}
	totalPage := int(math.Ceil(float64(totalData) / float64(pageSize)))
	return KasGuruResponse, totalPage, nil
	}

func ( c *kasRepository) Search(nama, tanggal string) ([]models.KasGuruResponse, error){
	nama = strings.ToLower(nama)
	tanggal = strings.ToLower(tanggal)


	var KasGuru []models.KasGuruResponse

	query := c.db.Table("kas_gurus").
	Select("kas_gurus.id, gurus.nama as nama_guru, kas_gurus.jumlah, kas_gurus.tanggal_bayar").
	Joins("JOIN gurus ON kas_gurus.guru_id = gurus.nip").
	Where("LOWER(gurus.nama) LIKE ? AND LOWER(kas_gurus.tanggal_bayar) LIKE ?", "%" +nama+ "%", "%"+tanggal+"%")

	if err := query.Find(&KasGuru).Error; err != nil {
        return nil, err
    }

	return KasGuru, nil
}

func (c *kasRepository) GetByID(id int) (*models.HistoryPembayaranKas, error){
	var historiBayarKas models.KasGuru
	err := c.db.Preload("Guru").Where("id = ?", id).First(&historiBayarKas).Error
	if err != nil {
		return nil, err
	}

	kasResponse := models.HistoryPembayaranKas{
		ID: uint(historiBayarKas.ID),
		NIP: historiBayarKas.Guru.Nip,
		Nama: historiBayarKas.Guru.Nama,
		Jumlah_Bayar: historiBayarKas.Jumlah,
		Tanggal: historiBayarKas.TanggalBayar,

	}

	return &kasResponse, nil
}

func (c *kasRepository) GetAmbilByID(id int) (*models.HistoryPengambilanKas, error){
	var historiAmbilKas models.PengambilanKas
	err := c.db.Where("id = ?", id).First(&historiAmbilKas).Error
	if err != nil {
		return nil, err
	}

	kasResponse := models.HistoryPengambilanKas{
		ID: uint(historiAmbilKas.ID),
		NIP: historiAmbilKas.GuruID,
		Nama: historiAmbilKas.Nama,
		TanggalAmbil: historiAmbilKas.TanggalAmbil,
		JumlahAmbil: historiAmbilKas.JumlahAmbil,
	}

	return &kasResponse, nil
}
