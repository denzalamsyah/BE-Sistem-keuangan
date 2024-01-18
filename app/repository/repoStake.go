package repository

import (
	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type StakeholderRepository interface {
	Store(Stakeholder *models.Stakeholder) error
	Update(id int, Stakeholder models.Stakeholder) error
	Delete(id int) error
	GetByID(id int) (*models.StakeholderResponse, error)
	GetList() ([]models.StakeholderResponse, error)
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

func (c *stakeholderRepository) GetList() ([]models.StakeholderResponse, error){
	var stake []models.Stakeholder

	err := c.db.Preload("Jabatan").Preload("Agama").Preload("Gender").Find(&stake).Error
	if err != nil {
		return nil, err
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
			TempatLahir: s.TanggalLahir,
			TanggalLahir: s.TanggalLahir,
			NomorTelepon: s.NomorTelepon,
			Email: s.Email,
			Alamat: s.Alamat,
		})
	}
	return stakeResponse, nil
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
		TempatLahir: stake.TanggalLahir,
		TanggalLahir: stake.TanggalLahir,
		NomorTelepon: stake.NomorTelepon,
		Email: stake.Email,
		Alamat: stake.Alamat,
	}
	return &stakeResponse, nil
}