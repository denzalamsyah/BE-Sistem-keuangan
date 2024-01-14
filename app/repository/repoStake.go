package repository

import (
	"github.com/denzalamsyah/simak/app/models"
	"gorm.io/gorm"
)

type StakeholderRepository interface {
	Store(Stakeholder *models.Stakeholder) error
	Update(id int, Stakeholder models.Stakeholder) error
	Delete(id int) error
	GetByID(id int) (*models.Stakeholder, error)
	GetList() ([]models.Stakeholder, error)
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

func (c *stakeholderRepository) GetList() ([]models.Stakeholder, error){
	var stake []models.Stakeholder

	err := c.db.Find(&stake).Error

	if err != nil{
		return nil, err
	}
	return stake, nil
}

func (c *stakeholderRepository) GetByID(id int) (*models.Stakeholder, error){
	var stake models.Stakeholder

	err := c.db.Where("id = ?", id).Find(&stake).Error

	if err != nil{
		return nil,  err
	}
	return &stake, nil
}