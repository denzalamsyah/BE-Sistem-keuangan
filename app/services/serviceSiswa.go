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
	GetList(page, pageSize int) ([]models.SiswaResponse, int, error)
	HistoryPembayaranSiswa(siswaID int) ([]models.HistoryPembayaran, error)
	GetTotalGenderCount() (int, int, error)
	Search(name, kelas, nisn, jurusan string) ([]models.SiswaResponse, error)

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

func (c *siswaServices) GetList(page, pageSize int) ([]models.SiswaResponse, int, error) {
    siswas, totalPage, err := c.siswaRepo.GetList(page, pageSize)
    if err != nil {
        return nil, 0, err
    }
    return siswas, totalPage, nil
}


func (c *siswaServices) HistoryPembayaranSiswa(siswaID int) ([]models.HistoryPembayaran, error) {
	history, err := c.siswaRepo.HistoryPembayaranSiswa(siswaID)
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (c *siswaServices) GetTotalGenderCount() (int, int, error) {
	countLakiLaki, countPerempuan, err := c.siswaRepo.GetTotalGenderCount()
	if err != nil {
		return 0, 0, err
	}
	return int(countLakiLaki), int(countPerempuan), nil
}

func(c *siswaServices) Search(name, nisn, kelas, jurusan string) ([]models.SiswaResponse, error){
	siswa, err := c.siswaRepo.Search(name, nisn, kelas, jurusan)
	if err != nil {
        return nil, err
    }
	return siswa, nil
}

