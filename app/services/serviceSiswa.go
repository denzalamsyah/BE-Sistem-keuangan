package services

import (
	"github.com/denzalamsyah/simak/app/models"
	"github.com/denzalamsyah/simak/app/repository"
)

type SiswaServices interface {
	Store(Siswa *models.Siswa) error
	Update(id int, Siswa models.Siswa) error
	Delete(id int) error
	GetByID(id int) (*models.SiswaResponse, error)
	GetList() ([]models.SiswaResponse, error)
	HistoryPembayaranSiswa(siswaID int) ([]models.HistoryPembayaran, error)
}

type siswaServices struct {
	siswaRepo repository.SiswaRepository
}
func NewSiswaService(siswaRepo repository.SiswaRepository) SiswaServices {
	return &siswaServices{siswaRepo}
}


func (c *siswaServices) Store(siswa *models.Siswa) error {
	err := c.siswaRepo.Store(siswa)
	if err != nil {
		return err
	}
	return nil
}

func (c *siswaServices) Update(id int, siswa models.Siswa) error {
	err := c.siswaRepo.Update(id, siswa)
	if err != nil {
		
		return err
	}
	return nil
}

func (c *siswaServices) Delete(id int) error {
	err := c.siswaRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *siswaServices) GetByID(id int) (*models.SiswaResponse, error) {
	siswa, err := c.siswaRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return siswa, nil
}

func (c *siswaServices) GetList() ([]models.SiswaResponse, error) {
	siswas, err := c.siswaRepo.GetList()
	if err != nil {
		return nil, err
	}
	return siswas, nil
}

func (c *siswaServices) HistoryPembayaranSiswa(siswaID int) ([]models.HistoryPembayaran, error) {
	history, err := c.siswaRepo.HistoryPembayaranSiswa(siswaID)
	if err != nil {
		return nil, err
	}
	return history, nil
}
