package repository

import (
	"math"

	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type StakeholderRepository interface {
	Store(Stakeholder *models.Stakeholder) error
	Update(id int, Stakeholder models.Stakeholder) error
	Delete(id int) error
	GetByID(id int) (*models.StakeholderResponse, error)
	GetList(page, pageSize int) ([]models.StakeholderResponse, int, error)
}

type stakeholderRepository struct{
	db *gorm.DB
}


func NewStakeholderRepo(db *gorm.DB) *stakeholderRepository{
	return &stakeholderRepository{db}
}


func(c *stakeholderRepository) Store(Stakeholder *models.Stakeholder) error{
	err := c.db.Create(Stakeholder).Error
	if err != nil{
		return err
	}
	return nil
}

func(c *stakeholderRepository) Update(id int, Stakeholder models.Stakeholder) error{
	err := c.db.Model(&Stakeholder).Where("id = ?", id).Updates(&Stakeholder).Error

	if err != nil{
		return err
	}

	return nil
}

func(c *stakeholderRepository) Delete(id int) (error){
	var stake models.Stakeholder
	err := c.db.Where("id = ?", id).Delete(&stake).Error

	if err != nil{
		return err
	}
	return nil
}

func (c *stakeholderRepository) GetList(page, pageSize int) ([]models.StakeholderResponse, int, error){
	var stake []models.Stakeholder

	err := c.db.Preload("Jabatan").Preload("Agama").Preload("Gender").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&stake).Error
	if err != nil {
		return nil,0, err
	}

	var totalData int64
	if err := c.db.Model(&models.Stakeholder{}).Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	var stakeResponse []models.StakeholderResponse
	for _, s := range stake{
		stakeResponse = append(stakeResponse, models.StakeholderResponse{
			ID: s.ID,
			Nama: s.Nama,
			NIP: s.NIP,
			Agama: s.Agama.Nama,
			Gender: s.Gender.Nama,
			Jabatan: s.Jabatan.Nama,
			TempatLahir: s.TempatLahir,
			TanggalLahir: s.TanggalLahir,
			NomorTelepon: s.NomorTelepon,
			Email: s.Email,
			Alamat: s.Alamat,
			Gambar: s.Gambar,
		})
	}
	totalPage := int(math.Ceil(float64(totalData) / float64(pageSize)))
	return stakeResponse, totalPage, nil
}

func (c *stakeholderRepository) GetByID(id int) (*models.StakeholderResponse, error){
	var stake models.Stakeholder

	err := c.db.Preload("Jabatan").Preload("Agama").Preload("Gender").Where("id = ?", id).First(&stake).Error
	if err != nil {
		return nil, err
	}

	stakeResponse := models.StakeholderResponse{
		ID: stake.ID,
		Nama: stake.Nama,
		NIP: stake.NIP,
		Agama: stake.Agama.Nama,
		Gender: stake.Gender.Nama,
		Jabatan: stake.Jabatan.Nama,
		TempatLahir: stake.TempatLahir,
		TanggalLahir: stake.TanggalLahir,
		NomorTelepon: stake.NomorTelepon,
		Email: stake.Email,
		Alamat: stake.Alamat,
		Gambar: stake.Gambar,
	}
	return &stakeResponse, nil
}