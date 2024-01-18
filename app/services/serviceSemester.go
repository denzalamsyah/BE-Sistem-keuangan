package services

import (
	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/repository"
)

type SemesterServices interface {
	Store(PembayaranSemester *models.PembayaranSemester) error
	Update(id int, PembayaranSemester models.PembayaranSemester) error
	Delete(id int) error
	GetByID(id int) (*models.PembayaranSemesterResponse, error)
	GetList() ([]models.PembayaranSemesterResponse, error)
}

type semesterServices struct{
	semesterRepo repository.SemesterRepository
}

func NewSemesterService(semesterRepo repository.SemesterRepository) SemesterServices {
	return &semesterServices{semesterRepo}
}

func (c *semesterServices) Store(PembayaranSemester *models.PembayaranSemester) error {
	err := c.semesterRepo.Store(PembayaranSemester)
	if err != nil {
		return err
	}
	return nil
}

func (c *semesterServices) Update(id int, PembayaranSemester models.PembayaranSemester) error {
	err := c.semesterRepo.Update(id, PembayaranSemester)
	if err != nil {
		
		return err
	}
	return nil
}

func (c *semesterServices) Delete(id int) error {
	err := c.semesterRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *semesterServices) GetByID(id int) (*models.PembayaranSemesterResponse, error) {
	PembayaranSemester, err := c.semesterRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return PembayaranSemester, nil
}

func (c *semesterServices) GetList() ([]models.PembayaranSemesterResponse, error) {
	PembayaranSemesters, err := c.semesterRepo.GetList()
	if err != nil {
		return nil, err
	}
	return PembayaranSemesters, nil
}
