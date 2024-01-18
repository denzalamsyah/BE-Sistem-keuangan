package services

import (
	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/repository"
)

type StakeholderServices interface {
	Store(Stakeholder *models.Stakeholder) error
	Update(id int, Stakeholder models.Stakeholder) error
	Delete(id int)error
	GetByID(id int) (*models.StakeholderResponse, error)
	GetList() ([]models.StakeholderResponse, error)
}

type stakeholderServices struct{
	stakeRepo repository.StakeholderRepository
}
func NewStakeService(stakeRepo repository.StakeholderRepository) StakeholderServices{
	return &stakeholderServices{stakeRepo}
}
func (c *stakeholderServices) Store(Stakeholder *models.Stakeholder) error{
	err := c.stakeRepo.Store(Stakeholder)
	if err != nil {
		return err
	}
	return nil
}

func (c *stakeholderServices) Update(id int, Stakeholder models.Stakeholder) error{
	err := c.stakeRepo.Update(id, Stakeholder)
	if err != nil {
		
		return err
	}
	return nil
}

func (c *stakeholderServices) Delete(id int) error{
	err := c.stakeRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *stakeholderServices) GetByID(id int) (*models.StakeholderResponse, error){
	stake, err := c.stakeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return stake, nil
}

func (c *stakeholderServices) GetList()([]models.StakeholderResponse, error){
	stake, err := c.stakeRepo.GetList()
	if err != nil {
		return nil, err
	}
	return stake, nil
}