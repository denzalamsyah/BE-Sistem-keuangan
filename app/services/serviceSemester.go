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
	GetList(page, pageSize int) ([]models.PembayaranSemesterResponse, int, error)
	GetListByCategory(page, pageSize int, category string) ([]models.PembayaranSemesterResponse, int, error)
	Search(siswa, transaksi, semester, tanggal, nisn, kategori  string) ([]models.PembayaranSemesterResponse, error)
	GetLunasByNISN(nisn string) ([]models.PembayaranSemesterResponse, error)
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

func (c *semesterServices) GetList(page, pageSize int) ([]models.PembayaranSemesterResponse, int,error) {
	PembayaranSemesters, totalPage, err := c.semesterRepo.GetList(page, pageSize)
	if err != nil {
		return nil,0, err
	}
	return PembayaranSemesters, totalPage, nil
}

func (c *semesterServices) GetListByCategory(page, pageSize int, category string) ([]models.PembayaranSemesterResponse, int, error){
	PembayaranSemesters, totalPage, err := c.semesterRepo.GetListByCategory(page, pageSize, category)
	if err != nil {
		return nil,0, err
	}
	return PembayaranSemesters, totalPage, nil
}

func (c *semesterServices)  Search(siswa, transaksi, semester, tanggal, nisn, kategori string) ([]models.PembayaranSemesterResponse, error){
	pembeyaran, err := c.semesterRepo.Search(siswa, transaksi, semester, tanggal, nisn, kategori)

	if err != nil {
        return nil, err
    }
	return pembeyaran, nil
}
func (s *semesterServices) GetLunasByNISN(nisn string) ([]models.PembayaranSemesterResponse, error) {
	return s.semesterRepo.GetLunasByNISN(nisn)
}
